Below is a **Go (Echo v4) “Twitch module”** that plugs into the GitHub-style skeleton you already have.
It covers every *MVP* feature listed in the integration spec (viewer count, stream info, basic chat commands, secure OAuth) ([Twitch Developers][1])

---

## 1 Register a Twitch application

1. Go to [https://dev.twitch.tv/console/apps](https://dev.twitch.tv/console/apps) → **“Register Your App”**.
2. Set an **OAuth Redirect URL** e.g.

```
http://localhost:1323/auth/twitch/callback
```

3. Note the **Client ID** and **Client Secret**.
4. Add the scopes we need:

```
user:read:email           # identify the broadcaster
channel:read:stream_key   # (future) if you want uptime/health
chat:read chat:edit       # read & send chat messages
```

---

## 2 Dependencies

```bash
go get github.com/markbates/goth
go get github.com/markbates/goth/providers/twitch
go get github.com/gempir/go-twitch-irc/v4        # IRC / chat
go get github.com/labstack/echo/v4
go get github.com/labstack/echo-contrib/sessions
```

---

## 3 Main wiring (main.go)

```go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo-contrib/sessions"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/twitch"
	"github.com/agentstation/egothic"      // same tiny adapter we used for GitHub
)

func main() {
	// ---------- 0. ENV ----------
	cid    := mustEnv("TWITCH_CLIENT_ID")
	secret := mustEnv("TWITCH_CLIENT_SECRET")
	sessionKey := mustEnv("SESSION_SECRET")
	redirect   := "http://localhost:1323/auth/twitch/callback"

	// ---------- 1. Echo ----------
	e := echo.New()
	store := sessions.NewCookieStore([]byte(sessionKey))
	store.Options = &sessions.Options{Path: "/", MaxAge: 86400*3, HttpOnly: true, SameSite: http.SameSiteLaxMode}
	e.Use(sessions.Sessions("tw-session", store))

	// ---------- 2. Goth provider ----------
	goth.UseProviders(
		twitch.New(cid, secret, redirect,
			"user:read:email", "chat:read", "chat:edit"),
	)

	// ---------- 3. OAuth endpoints ----------
	e.GET("/auth/twitch", egothic.BeginAuthHandler)

	e.GET("/auth/twitch/callback", func(c echo.Context) error {
		user, err := egothic.CompleteUserAuth(c) // exchanges code→token
		if err != nil { return err }

		// persist tokens in session (for API + chat)
		sess, _ := sessions.Get(c)
		sess.Values["uid"]   = user.UserID
		sess.Values["login"] = user.NickName
		sess.Values["token"] = user.AccessToken
		sess.Save(c.Request(), c.Response())

		go ensureIrcClient(user.AccessToken, user.NickName) // chat bot
		return c.Redirect(http.StatusFound, "/")
	})

	// ---------- 4. Protected API ----------
	api := e.Group("/api/twitch", ensureAuth)
	api.GET("/stream", getStreamInfo)   // viewer count etc.
	api.POST("/chat", sendChatMessage)  // ?msg=hello

	log.Fatal(e.Start(":1323"))
}

func mustEnv(k string) string { v := os.Getenv(k); if v=="" { log.Fatalf("%s missing", k)}; return v }

// --- auth guard (same pattern as GitHub) ---
func ensureAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := sessions.Get(c)
		if sess.Values["token"] == nil { return c.Redirect(http.StatusFound, "/auth/twitch") }
		return next(c)
	}
}
```

---

## 4 Viewer analytics endpoint

```go
type helixStreams struct {
	Data []struct {
		ID            string `json:"id"`
		UserLogin     string `json:"user_login"`
		GameName      string `json:"game_name"`
		Title         string `json:"title"`
		ViewerCount   int    `json:"viewer_count"`
		StartedAt     string `json:"started_at"`
	} `json:"data"`
}

func getStreamInfo(c echo.Context) error {
	sess, _ := sessions.Get(c)
	token   := sess.Values["token"].(string)
	clientID:= os.Getenv("TWITCH_CLIENT_ID")
	userID  := sess.Values["uid"].(string)

	req, _ := http.NewRequest("GET",
		"https://api.twitch.tv/helix/streams?user_id="+userID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Client-Id", clientID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil { return err }
	defer resp.Body.Close()

	var payload helixStreams
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return err
	}
	if len(payload.Data)==0 {
		return c.JSON(http.StatusOK, echo.Map{"live": false})
	}
	s := payload.Data[0]
	return c.JSON(http.StatusOK, echo.Map{
		"live":        true,
		"title":       s.Title,
		"game":        s.GameName,
		"viewers":     s.ViewerCount,
		"started_at":  s.StartedAt,
	})
}
```

*Endpoint used:* `GET /helix/streams` ([Twitch Developers][2])

---

## 5 Chat client (read + command-driven replies)

```go
var once sync.Once
func ensureIrcClient(token, login string) {
	once.Do(func() {
		client := twitchirc.NewClient(login, "oauth:"+token)

		// echo viewer commands back (very simple MVP)
		client.OnPrivateMessage(func(m twitchirc.PrivateMessage) {
			switch strings.TrimSpace(m.Message) {
			case "!task":
				client.Say(m.Channel, "Current task: implementing Twitch integration in Go!")
			case "!uptime":
				// call /api/twitch/stream internally or cache uptime
			}
		})

		if err := client.Connect(); err != nil {
			log.Println("IRC error:", err)
		}
	})
}

func sendChatMessage(c echo.Context) error {
	sess, _ := sessions.Get(c)
	token := sess.Values["token"].(string)
	login := sess.Values["login"].(string)

	client := twitchirc.NewClient(login, "oauth:"+token)
	if err := client.Connect(); err != nil { return err }
	defer client.Disconnect()

	msg := c.QueryParam("msg")
	client.Say(login, msg)
	return c.NoContent(http.StatusNoContent)
}
```

*Chat transport:* Twitch IRC (WebSocket optional) ([Twitch Developers][3], [Twitch Developer Forums][4])

---

## 6 Refresh tokens

Twitch tokens last four hours. When any request returns **401**, exchange your **refresh token**:

````bash
POST https://id.twitch.tv/oauth2/token
     ?grant_type=refresh_token
     &refresh_token=<refresh>
     &client_id=<cid>
     &client_secret=<secret>
``` :contentReference[oaicite:4]{index=4}

Store the new pair (`access_token`, `refresh_token`) back in the session.

---

## 7 Security & CORS

| Case | Recommendation |
|------|----------------|
| SPA served **same origin** as Echo | Continue using the secure, `HttpOnly` cookie (`tw-session`) you already set up. |
| SPA served from **`app.example.com`** → API on **`api.example.com`** | Keep the cookie but add `SameSite=None; Secure` and `e.Use(middleware.CORSWithConfig(... AllowCredentials=true))`. |
| Native mobile / CLI | Return a JWT (or raw Twitch access token) in JSON instead of a cookie. |

(See previous cookie-vs-token table for trade-offs.)

---

### What’s covered

* **OAuth 2.0 Authorization-Code flow** with secure session storage.  
* **Viewer analytics**: live/offline, title, game, viewer count.  
* **Chat interaction**: reads messages & responds to `!command`, plus an API to send messages.  
* **Extensible**: add any Helix call (`/channels`, `/analytics`) with the same bearer-token pattern.

You can now wire these endpoints into the dashboard/overlay that the integration spec describes and iterate toward Phase 2 features (FAQ bot, highlights, etc.) when ready.
::contentReference[oaicite:5]{index=5}
````

[1]: https://dev.twitch.tv/docs/api/reference?utm_source=chatgpt.com "Reference - Twitch Developers"
[2]: https://dev.twitch.tv/docs/chat/irc?utm_source=chatgpt.com "IRC Concepts - Twitch Developers"
[3]: https://dev.twitch.tv/docs/authentication/getting-tokens-oauth/?utm_source=chatgpt.com "Getting OAuth Access Tokens | Twitch Developers"
[4]: https://discuss.dev.twitch.com/t/oauth-authorization-flow-refreshing-access-token/62408?utm_source=chatgpt.com "OAuth authorization flow - refreshing access token - API"
