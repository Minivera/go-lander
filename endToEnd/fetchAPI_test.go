package endToEnd_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFetchAPI(t *testing.T) {
	page, err := browserContext.NewPage()
	require.NoError(t, err)

	_, err = page.Goto("http://localhost:8080/fetchAPI/")
	require.NoError(t, err)

	title, err := page.Locator("#app h1")
	require.NoError(t, err)

	titleContent, err := title.TextContent()
	require.NoError(t, err)
	assert.Equal(t, "Sample loading app", titleContent)

	loading, err := page.QuerySelector("#app marquee")
	require.NoError(t, err)

	loadingContent, err := loading.TextContent()
	require.NoError(t, err)
	assert.Equal(t, "Loading...", loadingContent)

	// Wait for the loading to happen
	time.Sleep(2 * time.Second)

	idInput, err := page.Locator("#app div input:first-of-type")
	require.NoError(t, err)

	inputValue, err := idInput.GetAttribute("value")
	require.NoError(t, err)
	assert.Equal(t, "1", inputValue)

	err = page.Close()
	require.NoError(t, err)
}
