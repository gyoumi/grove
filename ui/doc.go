// Package ui is grove's component library: themed building blocks (Button,
// Card, Input, Dialog, ...) styled with Tailwind utility classes on top of
// a CSS-variable design system, so one set of variables themes every
// component — including dark mode via the .dark class.
//
// The intended workflow is to own the source: `grove add button` copies a
// component's file(s) into your app's ui/ directory, where you edit them
// freely. Importing this package directly also works when you don't need
// to customize.
//
// These components assume the theme variables from `grove init`'s
// styles/input.css are present.
package ui
