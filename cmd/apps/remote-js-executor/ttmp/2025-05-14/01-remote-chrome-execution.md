One of the most bullet-proof ways is to drive your browser over the Chrome DevTools Protocol (CDP) or the equivalent in Firefox, have your CLI (“cursor”) connect to the browser’s debug port, and then simply send your snippet down as an “evaluate” call.

Below are two examples—one in Node (using puppeteer‐core), one in Go (using chromedp)—that all you have to do is start your browser with the remote-debugging port open, and then invoke your snippet on the active tab.

---

### 🔹 Option A: Node + puppeteer-core

1. Launch Chrome (or Edge) with

   ```bash
   /path/to/chrome --remote-debugging-port=9222
   ```

2. Then in your cursor tool (or a tiny helper script):

   ```js
   // inject.js
   const puppeteer = require('puppeteer-core');

   (async () => {
     // connect to the running browser
     const browser = await puppeteer.connect({
       browserURL: 'http://localhost:9222',
     });

     // grab the first open page (or find by URL)
     const [page] = await browser.pages();

     // your snippet from above as a single string
     const js = `
       (async function clickButtonsAndLogBodyText(rootId) {
         const subDiv = document.getElementById(rootId);
         if (!subDiv) { console.error('no id'); return; }
         const container = subDiv.parentElement.parentElement;
         const buttons = Array.from(container.querySelectorAll('button'));
         for (const btn of buttons) {
           btn.click();
           await new Promise(r => requestAnimationFrame(() => requestAnimationFrame(r)));
           console.log(document.body.innerText);
         }
       })('response-content-citations');
     `;

     // execute it in the page’s context
     await page.evaluate(js);

     // optional: keep the process alive to catch logs, or disconnect:
     // await browser.disconnect();
   })();
   ```

3. Run it with

   ```bash
   node inject.js
   ```

   you’ll see your `console.log()` lines streamed back to your terminal.

---

### 🔹 Option B: Go + chromedp

If your cursor tool is written in Go, you can use [chromedp](https://github.com/chromedp/chromedp):

```go
package main

import (
  "context"
  "log"
  "time"

  "github.com/chromedp/chromedp"
)

func main() {
  // connect to an existing Chrome with --remote-debugging-port=9222
  allocCtx, cancel := chromedp.NewRemoteAllocator(context.Background(), "ws://127.0.0.1:9222")
  defer cancel()

  ctx, cancel := chromedp.NewContext(allocCtx)
  defer cancel()

  // give Chrome a moment
  ctx, _ = context.WithTimeout(ctx, 10*time.Second)

  // your JS snippet
  js := `
    (async function clickButtonsAndLogBodyText(rootId) {
      const subDiv = document.getElementById(rootId);
      if (!subDiv) { console.error('no id'); return; }
      const container = subDiv.parentElement.parentElement;
      const buttons = Array.from(container.querySelectorAll('button'));
      for (const btn of buttons) {
        btn.click();
        await new Promise(r => requestAnimationFrame(() => requestAnimationFrame(r)));
        console.log(document.body.innerText);
      }
    })('response-content-citations');
  `

  // navigate if not already on target page
  if err := chromedp.Run(ctx,
    chromedp.Navigate(`https://your.target.page/`),
    chromedp.Sleep(1*time.Second), // let it settle
    chromedp.Evaluate(js, nil),
  ); err != nil {
    log.Fatal(err)
  }

  // keep running if you want to read the Chrome logs, or just exit:
  time.Sleep(2 * time.Second)
}
```

**Why this is “good”:**

* **Unattended** – no manual copy‐paste.
* **Scriptable** – integrate into your cursor CLI.
* **Reliable** – uses the same protocol DevTools uses under the hood.

---

#### Bonus: Bookmarklet‐style for one‐off tests

If you just want a super-quick manual route, wrap your function in a `javascript:` URL and paste that into the address bar:

```text
javascript:(function(){ /* ...your snippet... */ })();
```

…then hit ↵ and it runs on the current page. But for fully automated, go with CDP + Puppeteer or chromedp as shown above.
