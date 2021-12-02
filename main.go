package main

import (
	"context"
	"github.com/chromedp/chromedp"
	"github.com/pkg/errors"
	"log"
	"os"
	"sync"
)

var (
	commonContext *context.Context
	commonContextMtx sync.Mutex
)


func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	log.Println("app is running")
	if err := func() error {
		commonContextMtx.Lock()
		defer commonContextMtx.Unlock()

		if commonContext == nil {
			ctx0, _ := chromedp.NewContext(
				context.Background(),
				chromedp.WithLogf(log.Printf),
			)
			commonContext = &ctx0
		}
		const urlEnvKey = "URL"
		u := os.Getenv(urlEnvKey)
		if len(u) == 0 {
			return errors.Errorf("missing '%+v' env var", urlEnvKey)
		}
		selector := `title`
		log.Println("requesting", u)
		log.Println("selector", selector)
		var result string
		err := chromedp.Run(*commonContext,
			chromedp.Navigate(u),
			chromedp.WaitReady(selector),
			chromedp.OuterHTML(selector, &result),
		)
		if err != nil {
			return errors.WithStack(err)
		}
		log.Printf("result:\n%s", result)
		return nil
	}(); err != nil {
		log.Printf("error %+v \n", err)
	}
}