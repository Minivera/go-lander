// Package state is an experimental package that add support for global state management to the library.
// The state is saved in a central store and components can hook into it by using the render props
// of the provided consumer components. The state store does not provide any optimizations beyond the
// base optimizations of the library, all components will be set to rerender is the state changes.
//
// It is necessary to create a global store to use the store capabilities, it should be created in a
// central location and reused throughout the lifecycle of the app.
//
// This package is even more unstable than the library itself, use at your own risk.
package state
