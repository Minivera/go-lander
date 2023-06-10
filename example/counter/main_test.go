package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/minivera/go-lander"
	"github.com/minivera/go-lander/nodes"
	"github.com/stretchr/testify/require"

	js "github.com/minivera/go-lander/go-wasm-dom"
)

func TestCounterApp(t *testing.T) {
	js.StartTestMode(t)
	js.EnableDebug(t)

	app := counterApp{}

	_, err := lander.RenderInto(
		lander.Component(app.render, nodes.Props{}, nodes.Children{}), "#app")
	require.NoError(t, err)

	assert.Equal(t, "test", "invalid")

	js.EndTestMode()
}
