package endToEnd_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginForm(t *testing.T) {
	page, err := browserContext.NewPage()
	require.NoError(t, err)

	_, err = page.Goto("http://localhost:8080/loginForm/")
	require.NoError(t, err)

	title, err := page.Locator("#app h1")
	require.NoError(t, err)

	titleContent, err := title.TextContent()
	require.NoError(t, err)
	assert.Equal(t, "Log into our app", titleContent)

	err = page.Close()
	require.NoError(t, err)
}
