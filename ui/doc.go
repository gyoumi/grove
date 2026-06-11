// Package ui ports shadcn/ui components to grove: the same Tailwind class
// recipes and CSS-variable theme, expressed as Go functions. The visual
// design and class strings follow shadcn/ui (MIT) closely so apps styled
// for shadcn look identical here.
//
// Like shadcn, the intended workflow is to own the source: `grove add
// button` copies a component's file(s) into your app's ui/ directory, where
// you edit them freely. Importing this package directly also works when
// you don't need to customize.
//
// These components assume the shadcn theme variables from `grove init`'s
// styles/input.css are present.
package ui
