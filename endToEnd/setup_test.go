package endToEnd_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/playwright-community/playwright-go"
)

var browserContext playwright.BrowserContext

func TestMain(m *testing.M) {
	pw, err := playwright.Run()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	context, err := browser.NewContext()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	browserContext = context

	code := m.Run()

	err = browser.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = pw.Stop()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(code)
}
