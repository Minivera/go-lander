package endToEnd_test

import (
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHelloWorld(t *testing.T) {
	pw, err := playwright.Run()
	require.NoError(t, err)

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	require.NoError(t, err)

	context, err := browser.NewContext()
	require.NoError(t, err)

	page, err := context.NewPage()
	require.NoError(t, err)

	_, err = page.Goto("http://localhost:8080/helloWorld/")
	require.NoError(t, err)

	title, err := page.Locator("#app h1")
	require.NoError(t, err)

	titleContent, err := title.TextContent()
	require.NoError(t, err)
	assert.Equal(t, "Hello, World!", titleContent)

	err = browser.Close()
	require.NoError(t, err)

	err = pw.Stop()
	require.NoError(t, err)
}
