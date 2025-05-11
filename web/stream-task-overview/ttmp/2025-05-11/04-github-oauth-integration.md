Below is a minimal, production-ready\* GitHub SSO module for an Echo-based Go backend.
It uses:

* **labstack/echo v4** – HTTP framework
* **markbates/goth** – multi-provider OAuth2 helpers
* **agentstation/egothic** – thin adapter that makes Goth fit Echo idiomatically ([github.com][1])
* **echo-contrib/sessions** – cookie session storage

\*“Production-ready” here means the flow, namespacing, and cookie flags are correct. Feel free to swap in Redis, JWTs, or your own session store without changing the high-level flow.

---

## 1 . Register an OAuth App on GitHub

1. **Settings → Developer settings → OAuth Apps → “New OAuth App”.**
2. **Authorization callback URL:**

   ```
   http://localhost:1323/auth/github/callback
   ```
3. Note the **Client ID** and **Client Secret**.

Recommended scopes for a back-end-only sign-in:

```text
read:user user:email
```

---

## 2 . Create a small main.go

```go
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/agentstation/egothic"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo-contrib/sessions"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
)

func main() {
	// ------------------------------------------------------------------
	// 0. Runtime configuration
	// ------------------------------------------------------------------
	clientID     := mustEnv("GITHUB_CLIENT_ID")
	clientSecret := mustEnv("GITHUB_CLIENT_SECRET")
	sessionKey   := mustEnv("SESSION_SECRET")          // 32+ random bytes
	callbackURL  := "http://localhost:1323/auth/github/callback"

	// ------------------------------------------------------------------
	// 1. Bootstrap Echo & sessions
	// ------------------------------------------------------------------
	e := echo.New()
	store := sessions.NewCookieStore([]byte(sessionKey))
	store.Options = &sessions.Options{
		Path:     "/", MaxAge: 86400 * 7, HttpOnly: true, SameSite: http.SameSiteLaxMode,
	}
	e.Use(sessions.Sessions("gh-session", store))

	// ------------------------------------------------------------------
	// 2. Register GitHub provider with Goth
	// ------------------------------------------------------------------
	goth.UseProviders(
		github.New(clientID, clientSecret, callbackURL, "read:user", "user:email"),
	)

	// ------------------------------------------------------------------
	// 3. Auth routes
	// ------------------------------------------------------------------
	e.GET("/auth/github", egothic.BeginAuthHandler) // redirects to GitHub

	e.GET("/auth/github/callback", func(c echo.Context) error {
		user, err := egothic.CompleteUserAuth(c) // handles state, token exchange, etc.
		if err != nil {
			return err
		}
		// Persist the interesting bits in your session
		sess, _ := sessions.Get(c)
		sess.Values["uid"]   = user.UserID
		sess.Values["login"] = user.NickName
		sess.Save(c.Request(), c.Response())
		return c.Redirect(http.StatusFound, "/") // or frontend URL
	})

	// ------------------------------------------------------------------
	// 4. Protected API group
	// ------------------------------------------------------------------
	api := e.Group("/api", ensureAuth)
	api.GET("/me", func(c echo.Context) error {
		sess, _ := sessions.Get(c)
		return c.JSON(http.StatusOK, echo.Map{
			"user": sess.Values["login"],
		})
	})

	log.Fatal(e.Start(":1323"))
}

// ensureAuth makes sure a user is logged in before hitting protected endpoints.
func ensureAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := sessions.Get(c)
		if sess.Values["uid"] == nil {
			// not logged in -- kick to OAuth
			return c.Redirect(http.StatusFound, "/auth/github")
		}
		return next(c)
	}
}

func mustEnv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		log.Fatalf("%s must be set", k)
	}
	return v
}
```

### go.mod

```go
module github.com/yourcorp/echogh

go 1.22

require (
	github.com/agentstation/egothic v0.1.4 // latest as of 2025-05
	github.com/labstack/echo/v4     v4.11.2
	github.com/labstack/echo-contrib/sessions v0.3.0
	github.com/markbates/goth       v1.75.0
)
```

---

## 3 . Run it

```bash
export GITHUB_CLIENT_ID=...
export GITHUB_CLIENT_SECRET=...
export SESSION_SECRET=$(openssl rand -base64 32)

go run .
```

1. Navigate to `http://localhost:1323/auth/github` → you’ll be redirected to GitHub.
2. After authorizing, you return to `/auth/github/callback`, which stores the session.
3. Hitting `GET /api/me` now returns your GitHub handle.

---

## 4 . Integrating with a Front-End

If you need a token for SPA / mobile clients instead of a cookie:

1. Replace the `ensureAuth` middleware with one that expects an **Authorization: Bearer <jwt>**.
2. In the callback handler, mint a JWT:

```go
jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	"sub":  user.UserID,
	"login": user.NickName,
	"exp":  time.Now().Add(24*time.Hour).Unix(),
})
signed, _ := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
return c.JSON(http.StatusOK, echo.Map{"token": signed})
```

3. The rest of your API stays framework-agnostic.

---

### Why Goth + egothic instead of raw `golang.org/x/oauth2`?

* Handles CSRF-safe `state` param, token refresh, and provider quirks for you.
* Adding more SSO providers later (GitLab, Google, Azure AD, etc.) is one-liner.
* Egothic strips the `net/http` boilerplate so Echo’s context remains first-class.

(If you prefer to see the bare-metal `oauth2.Config` flow, the **README** in *go-simple-oauth2* shows the full exchange without abstraction.)&#x20;

Happy hacking!

[1]: https://github.com/agentstation/egothic "GitHub - agentstation/egothic: Gothic for Echo"
