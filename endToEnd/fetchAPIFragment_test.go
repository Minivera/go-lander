package endToEnd_test

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchAPIFragment(t *testing.T) {
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

	_, err = page.Goto("http://localhost:8080/fetchAPIFragment/")
	require.NoError(t, err)

	title, err := page.Locator("#app h1")
	require.NoError(t, err)

	titleContent, err := title.TextContent()
	require.NoError(t, err)
	assert.Equal(t, "Sample loading app with Fragments", titleContent)

	loading, err := page.QuerySelector("#app marquee")
	require.NoError(t, err)

	loadingContent, err := loading.TextContent()
	require.NoError(t, err)
	assert.Equal(t, "Loading...", loadingContent)

	// Wait for the loading to happen
	time.Sleep(2 * time.Second)

	for i := 1; i <= 30; i++ {
		idInput, err := page.Locator(fmt.Sprintf("#app div fieldset:nth-of-type(%d) input:first-of-type", i))
		require.NoError(t, err)

		inputValue, err := idInput.GetAttribute("value")
		require.NoError(t, err)
		assert.Equal(t, strconv.Itoa(i), inputValue)
	}

	err = browser.Close()
	require.NoError(t, err)

	err = pw.Stop()
	require.NoError(t, err)
}
