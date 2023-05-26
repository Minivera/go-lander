package endToEnd_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRouterWithHelmet(t *testing.T) {
	type expect struct {
		expectedPath    string
		expectedTitle   string
		selector        string
		expectedContent string
	}

	type action struct {
		clickOn string
		expect  expect
	}

	tcs := []struct {
		scenario string
		actions  []action
	}{
		{
			scenario: "HelloWorld should navigate to /hello, then back to home",
			actions: []action{
				{
					clickOn: "To /hello",
					expect: expect{
						expectedPath:    "/hello",
						expectedTitle:   "Sample routing app - Hello",
						selector:        "#app div h2",
						expectedContent: "Hello, world!",
					},
				},
				{
					clickOn: "Go back to Home",
					expect: expect{
						expectedPath:    "/",
						expectedTitle:   "Sample routing app",
						selector:        "#app div h2",
						expectedContent: "Home page",
					},
				},
			},
		},
		{
			scenario: "App should navigate to /app, then back to home",
			actions: []action{
				{
					clickOn: "To /app",
					expect: expect{
						expectedPath:    "/app",
						expectedTitle:   "Sample routing app - App",
						selector:        "#app div h2",
						expectedContent: "Welcome to the app",
					},
				},
				{
					clickOn: "Go back to Home",
					expect: expect{
						expectedPath:    "/",
						expectedTitle:   "Sample routing app",
						selector:        "#app div h2",
						expectedContent: "Home page",
					},
				},
			},
		},
		{
			scenario: "App with pattern matching should work as expected",
			actions: []action{
				{
					clickOn: "To /app",
					expect: expect{
						expectedPath:    "/app",
						expectedTitle:   "Sample routing app - App",
						selector:        "#app div h2",
						expectedContent: "Welcome to the app",
					},
				},
				{
					clickOn: "Test the pattern matching",
					expect: expect{
						expectedPath:    "/app/something/other",
						expectedTitle:   "Sample routing app - App",
						selector:        "#app div div div",
						expectedContent: "Matched:Pathname: http://localhost:8080/app/something/otherPath somethingSubpath other",
					},
				},
				// Click twice just to be extra safe
				{
					clickOn: "Test the pattern matching",
					expect: expect{
						expectedPath:    "/app/something/other",
						expectedTitle:   "Sample routing app - App",
						selector:        "#app div div div",
						expectedContent: "Matched:Pathname: http://localhost:8080/app/something/otherPath somethingSubpath other",
					},
				},
				{
					clickOn: "Go back to Home",
					expect: expect{
						expectedPath:    "/",
						expectedTitle:   "Sample routing app",
						selector:        "#app div h2",
						expectedContent: "Home page",
					},
				},
			},
		},
		{
			scenario: "App should be able to cycle between home and app",
			actions: []action{
				{
					clickOn: "To /app",
					expect: expect{
						expectedPath:    "/app",
						expectedTitle:   "Sample routing app - App",
						selector:        "#app div h2",
						expectedContent: "Welcome to the app",
					},
				},
				{
					clickOn: "Go back to Home",
					expect: expect{
						expectedPath:    "/",
						expectedTitle:   "Sample routing app",
						selector:        "#app div h2",
						expectedContent: "Home page",
					},
				},
				{
					clickOn: "To /app",
					expect: expect{
						expectedPath:    "/app",
						expectedTitle:   "Sample routing app - App",
						selector:        "#app div h2",
						expectedContent: "Welcome to the app",
					},
				},
				{
					clickOn: "Go back to Home",
					expect: expect{
						expectedPath:    "/",
						expectedTitle:   "Sample routing app",
						selector:        "#app div h2",
						expectedContent: "Home page",
					},
				},
				{
					clickOn: "To /app",
					expect: expect{
						expectedPath:    "/app",
						expectedTitle:   "Sample routing app - App",
						selector:        "#app div h2",
						expectedContent: "Welcome to the app",
					},
				},
				{
					clickOn: "Go back to Home",
					expect: expect{
						expectedPath:    "/",
						expectedTitle:   "Sample routing app",
						selector:        "#app div h2",
						expectedContent: "Home page",
					},
				},
			},
		},
		{
			scenario: "Redirect should keep the user in the app",
			actions: []action{
				{
					clickOn: "To /redirect, which will send us back here",
					expect: expect{
						expectedPath:    "/",
						expectedTitle:   "Sample routing app",
						selector:        "#app div h2",
						expectedContent: "Home page",
					},
				},
			},
		},
		{
			scenario: "Any other path should send to 404",
			actions: []action{
				{
					clickOn: "To the 404 page",
					expect: expect{
						expectedPath:    "/notfound",
						expectedTitle:   "Sample routing app - Not found",
						selector:        "#app div h2",
						expectedContent: "404! `http://localhost:8080/notfound` was not found",
					},
				},
				{
					clickOn: "Go back to Home",
					expect: expect{
						expectedPath:    "/",
						expectedTitle:   "Sample routing app",
						selector:        "#app div h2",
						expectedContent: "Home page",
					},
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			page, err := browserContext.NewPage()
			require.NoError(t, err)

			_, err = page.Goto("http://localhost:8080/routerWithHelmet/")
			require.NoError(t, err)

			title, err := page.Locator("#app h1")
			require.NoError(t, err)

			titleContent, err := title.TextContent()
			require.NoError(t, err)
			assert.Equal(t, "Sample routing app", titleContent)

			for _, action := range tc.actions {
				anchor, err := page.EvaluateHandle(
					fmt.Sprintf("() => [...document.querySelectorAll('#app div a')].find(el => el.innerText === '%s')", action.clickOn),
				)
				require.NoError(t, err)

				err = anchor.AsElement().Click()
				require.NoError(t, err)

				selectedText, err := page.Locator(action.expect.selector)
				require.NoError(t, err)

				selectedTextContent, err := selectedText.TextContent()
				require.NoError(t, err)
				assert.Equal(t, action.expect.expectedContent, selectedTextContent)

				assert.Contains(t, page.URL(), action.expect.expectedPath)

				title, err := page.Title()
				require.NoError(t, err)
				assert.Contains(t, title, action.expect.expectedTitle)
			}

			err = page.Close()
			require.NoError(t, err)
		})
	}
}
