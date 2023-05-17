package endToEnd_test

import (
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCounter(t *testing.T) {
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

	_, err = page.Goto("http://localhost:8080/counter/")
	require.NoError(t, err)

	title, err := page.Locator("#app h1")
	require.NoError(t, err)

	titleContent, err := title.TextContent()
	require.NoError(t, err)
	assert.Equal(t, "Sample counter app", titleContent)

	app, err := page.Locator("#app > div > div:nth-of-type(1)")
	require.NoError(t, err)

	minusButton, err := app.Locator("button:nth-of-type(1)")
	require.NoError(t, err)
	plusButton, err := app.Locator("button:nth-of-type(2)")
	require.NoError(t, err)

	content, err := app.Locator("div")
	require.NoError(t, err)

	contentText, err := content.TextContent()
	require.NoError(t, err)
	assert.Equal(t, "Counter is at: 0", contentText)

	// Press the plus button once, counter should be at 1
	err = plusButton.Click()
	require.NoError(t, err)

	contentText, err = content.TextContent()
	require.NoError(t, err)
	assert.Equal(t, "Counter is at: 1", contentText)

	// Press the minus button once, counter should be at 0
	err = minusButton.Click()
	require.NoError(t, err)

	contentText, err = content.TextContent()
	require.NoError(t, err)
	assert.Equal(t, "Counter is at: 0", contentText)

	// Press the minus button again, counter should get to -1
	err = minusButton.Click()
	require.NoError(t, err)

	contentText, err = content.TextContent()
	require.NoError(t, err)
	assert.Equal(t, "Counter is at: -1", contentText)

	// Press the plus button 100 times, content should be at 100
	for i := 0; i <= 100; i++ {
		err = plusButton.Click()
		require.NoError(t, err)
	}

	contentText, err = content.TextContent()
	require.NoError(t, err)
	assert.Equal(t, "Counter is at: 100", contentText)

	err = browser.Close()
	require.NoError(t, err)

	err = pw.Stop()
	require.NoError(t, err)
}
