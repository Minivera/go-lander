package endToEnd_test

import (
	"testing"

	"github.com/playwright-community/playwright-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func stringPointer(val string) *string {
	return &val
}

func TestTodos(t *testing.T) {
	pw, err := playwright.Run()
	require.NoError(t, err)

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	require.NoError(t, err)

	context, err := browser.NewContext()
	require.NoError(t, err)

	type todoState struct {
		checked bool
		text    string
	}

	type action struct {
		writeTodo     *string
		clickOn       *string
		expectedTodos []todoState
	}

	tcs := []struct {
		scenario string
		actions  []action
	}{
		{
			scenario: "checking the default todo should work",
			actions: []action{
				{
					clickOn: stringPointer("#app ul li:first-of-type input"),
					expectedTodos: []todoState{
						{
							checked: true,
							text:    "write more examples",
						},
					},
				},
			},
		},
		{
			scenario: "deleting the default todo should work",
			actions: []action{
				{
					clickOn:       stringPointer("#app ul li:first-of-type button"),
					expectedTodos: []todoState{},
				},
			},
		},
		{
			scenario: "should be able to add a new todo, check it, then delete the default todo",
			actions: []action{
				{
					writeTodo: stringPointer("new todo"),
					expectedTodos: []todoState{

						{
							checked: false,
							text:    "write more examples",
						},
						{
							checked: false,
							text:    "new todo",
						},
					},
				},
				{
					clickOn: stringPointer("#app ul li:nth-of-type(2) input"),
					expectedTodos: []todoState{
						{
							checked: false,
							text:    "write more examples",
						},
						{
							checked: true,
							text:    "new todo",
						},
					},
				},
				{
					clickOn: stringPointer("#app ul li:nth-of-type(1) button"),
					expectedTodos: []todoState{
						{
							checked: true,
							text:    "new todo",
						},
					},
				},
			},
		},
		{
			scenario: "complex case, creating and deleting todos",
			actions: []action{
				{
					writeTodo: stringPointer("new todo 1"),
					expectedTodos: []todoState{
						{
							checked: false,
							text:    "write more examples",
						},
						{
							checked: false,
							text:    "new todo 1",
						},
					},
				},
				{
					writeTodo: stringPointer("new todo 2"),
					expectedTodos: []todoState{
						{
							checked: false,
							text:    "write more examples",
						},
						{
							checked: false,
							text:    "new todo 1",
						},
						{
							checked: false,
							text:    "new todo 2",
						},
					},
				},
				{
					clickOn: stringPointer("#app ul li:nth-of-type(2) input"),
					expectedTodos: []todoState{
						{
							checked: false,
							text:    "write more examples",
						},
						{
							checked: true,
							text:    "new todo 1",
						},
						{
							checked: false,
							text:    "new todo 2",
						},
					},
				},
				{
					writeTodo: stringPointer("new todo 3"),
					expectedTodos: []todoState{
						{
							checked: false,
							text:    "write more examples",
						},
						{
							checked: true,
							text:    "new todo 1",
						},
						{
							checked: false,
							text:    "new todo 2",
						},
						{
							checked: false,
							text:    "new todo 3",
						},
					},
				},
				{
					clickOn: stringPointer("#app ul li:nth-of-type(4) input"),
					expectedTodos: []todoState{
						{
							checked: false,
							text:    "write more examples",
						},
						{
							checked: true,
							text:    "new todo 1",
						},
						{
							checked: false,
							text:    "new todo 2",
						},
						{
							checked: true,
							text:    "new todo 3",
						},
					},
				},
				{
					clickOn: stringPointer("#app ul li:nth-of-type(3) button"),
					expectedTodos: []todoState{
						{
							checked: false,
							text:    "write more examples",
						},
						{
							checked: true,
							text:    "new todo 1",
						},
						{
							checked: true,
							text:    "new todo 3",
						},
					},
				},
				{
					clickOn: stringPointer("#app ul li:nth-of-type(1) button"),
					expectedTodos: []todoState{
						{
							checked: true,
							text:    "new todo 1",
						},
						{
							checked: true,
							text:    "new todo 3",
						},
					},
				},
				{
					clickOn: stringPointer("#app ul li:nth-of-type(1) input"),
					expectedTodos: []todoState{
						{
							checked: false,
							text:    "new todo 1",
						},
						{
							checked: true,
							text:    "new todo 3",
						},
					},
				},
				{
					clickOn: stringPointer("#app ul li:nth-of-type(1) button"),
					expectedTodos: []todoState{
						{
							checked: true,
							text:    "new todo 3",
						},
					},
				},
				{
					clickOn:       stringPointer("#app ul li:nth-of-type(1) button"),
					expectedTodos: []todoState{},
				},
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.scenario, func(t *testing.T) {
			page, err := context.NewPage()
			require.NoError(t, err)

			_, err = page.Goto("http://localhost:8080/todos/")
			require.NoError(t, err)

			title, err := page.Locator("#app h1")
			require.NoError(t, err)

			titleContent, err := title.TextContent()
			require.NoError(t, err)
			assert.Equal(t, "Sample todo app", titleContent)

			firstTodo, err := page.QuerySelector("#app ul li:first-of-type")
			require.NoError(t, err)

			firstTodoChecked, err := firstTodo.EvalOnSelector("input", "el => el.checked")
			require.NoError(t, err)

			assert.False(t, firstTodoChecked.(bool))

			firstTodoValue, err := firstTodo.EvalOnSelector("strong", "el => el.textContent")
			require.NoError(t, err)

			assert.Equal(t, firstTodoValue.(string), "write more examples")

			for _, action := range tc.actions {
				if action.clickOn != nil {
					element, err := page.QuerySelector(*action.clickOn)
					require.NoError(t, err)

					err = element.Click()
					require.NoError(t, err)
				}

				if action.writeTodo != nil {
					input, err := page.QuerySelector("#app div > div > input")
					require.NoError(t, err)

					err = input.Type(*action.writeTodo)
					require.NoError(t, err)

					inputUpdated, err := page.Locator("#app div > div > input")
					require.NoError(t, err)

					inputUpdatedContent, err := inputUpdated.InputValue()
					require.NoError(t, err)
					assert.Equal(t, *action.writeTodo, inputUpdatedContent)

					addButton, err := page.QuerySelector("#app div > div > button")
					require.NoError(t, err)

					err = addButton.Click()
					require.NoError(t, err)
				}

				todos, err := page.QuerySelectorAll("#app ul li")
				require.NoError(t, err)

				assert.Len(t, todos, len(action.expectedTodos))
				for i, todo := range todos {
					todoContent, err := todo.EvaluateHandle("el => el.querySelector('strong').textContent")
					require.NoError(t, err)

					todoChecked, err := todo.QuerySelector("input")
					require.NoError(t, err)

					todoCheckedValue, err := todoChecked.IsChecked()
					require.NoError(t, err)

					assert.Equal(t, action.expectedTodos[i].checked, todoCheckedValue)
					assert.Equal(t, action.expectedTodos[i].text, todoContent.String())
				}
			}

			err = page.Close()
			require.NoError(t, err)
		})
	}

	err = browser.Close()
	require.NoError(t, err)

	err = pw.Stop()
	require.NoError(t, err)
}
