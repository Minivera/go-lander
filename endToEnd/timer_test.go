package endToEnd_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimer(t *testing.T) {
	page, err := browserContext.NewPage()
	require.NoError(t, err)

	_, err = page.Goto("http://localhost:8080/timer/")
	require.NoError(t, err)

	title, err := page.Locator("#app h1")
	require.NoError(t, err)

	titleContent, err := title.TextContent()
	require.NoError(t, err)
	assert.Equal(t, "Sample timer app", titleContent)

	app, err := page.Locator("#app > div > div:nth-of-type(1)")
	require.NoError(t, err)

	startButton, err := app.Locator("button:nth-of-type(1)")
	require.NoError(t, err)
	stopButton, err := app.Locator("button:nth-of-type(2)")
	require.NoError(t, err)
	resetButton, err := app.Locator("button:nth-of-type(3)")
	require.NoError(t, err)

	content, err := app.Locator("div")
	require.NoError(t, err)

	contentText, err := content.TextContent()
	require.NoError(t, err)
	assert.Equal(t, "0.00s", contentText)

	// Press the start button, wait for 5 seconds, then stop
	err = startButton.Click()
	require.NoError(t, err)

	time.Sleep(5 * time.Second)
	err = stopButton.Click()
	require.NoError(t, err)

	contentText, err = content.TextContent()
	require.NoError(t, err)
	assert.Contains(t, contentText, "5.")

	// Press again, wait for 1 second, then stop
	err = startButton.Click()
	require.NoError(t, err)

	time.Sleep(1 * time.Second)
	err = stopButton.Click()
	require.NoError(t, err)

	contentText, err = content.TextContent()
	require.NoError(t, err)
	assert.Contains(t, contentText, "6.")

	// Press reset, check that the timer has reset
	err = resetButton.Click()
	require.NoError(t, err)

	contentText, err = content.TextContent()
	require.NoError(t, err)
	assert.Equal(t, "0.00s", contentText)

	// Press start three times, wait for 1 second, then stop. Should not have gone faster
	err = startButton.Click()
	require.NoError(t, err)
	err = startButton.Click()
	require.NoError(t, err)
	err = startButton.Click()
	require.NoError(t, err)

	time.Sleep(1 * time.Second)
	err = stopButton.Click()
	require.NoError(t, err)

	contentText, err = content.TextContent()
	require.NoError(t, err)
	assert.Contains(t, contentText, "1.")

	err = page.Close()
	require.NoError(t, err)
}
