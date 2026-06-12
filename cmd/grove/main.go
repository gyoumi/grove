// Command grove is the development CLI for the grove framework:
//
//	grove init <app>   scaffold a new app (Tailwind + shadcn theme included)
//	grove serve        dev server with rebuild-on-save and live reload
//	grove build        production build with size report
//	grove add <name>   copy a ui component's source into your project
package main

import (
	"flag"
	"fmt"
	"os"
)

// parseMixed parses flags and positional arguments in any order (the flag
// package alone stops at the first positional), returning the positionals.
func parseMixed(fs *flag.FlagSet, args []string) []string {
	var pos []string
	for {
		fs.Parse(args)
		args = fs.Args()
		if len(args) == 0 {
			return pos
		}
		pos = append(pos, args[0])
		args = args[1:]
	}
}

const usage = `grove — a React-style framework for Go and WebAssembly

Usage:
  grove init <app> [--grove <path>]   scaffold a new app
  grove serve [-port 8080] [-dir .]   run the dev server
  grove build [-dir .]                production build into dist/
  grove add <component> [-dir .]      copy a ui component into ./ui/
  grove add -list                     list available components
`

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "init":
		err = runInit(os.Args[2:])
	case "serve":
		err = runServe(os.Args[2:])
	case "build":
		err = runBuild(os.Args[2:])
	case "add":
		err = runAdd(os.Args[2:])
	case "help", "-h", "--help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "grove: unknown command %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "grove:", err)
		os.Exit(1)
	}
}
