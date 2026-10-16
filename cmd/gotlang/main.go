// Command gotlang is the CLI for Gotlang, a Göteborg-humor keyword dialect
// of Go. It transpiles .got source to plain Go and hands off to the real
// go toolchain to build and run it.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"gotlang/internal/transpile"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	var err error
	switch os.Args[1] {
	case "translate":
		err = cmdTranslate(os.Args[2:])
	case "run":
		err = cmdRun(os.Args[2:])
	case "mecka":
		err = cmdMecka(os.Args[2:])
	case "-h", "-help", "--help", "help":
		usage()
		return
	default:
		fmt.Fprintf(os.Stderr, "gotlang: unknown command %q\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "gotlang: %v\n", err)
		os.Exit(1)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `gotlang is the Gotlang toolchain

Usage:

  gotlang translate <file.got> [-o out.go]   translate one file, print or write the Go source
  gotlang run <path>                         translate a package and run it
  gotlang mecka <path> [-o binary]           translate a package and compile it ("mecka" — GBG for tinkering a thing together)
`)
}

// cmdTranslate implements "gotlang translate".
func cmdTranslate(args []string) error {
	fs := flag.NewFlagSet("translate", flag.ExitOnError)
	out := fs.String("o", "", "write translated Go source to this file instead of stdout")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("translate: expected exactly one .got file, got %d args", fs.NArg())
	}

	src, err := os.ReadFile(fs.Arg(0))
	if err != nil {
		return err
	}
	translated, err := transpile.Transpile(src)
	if err != nil {
		return fmt.Errorf("translate %s: %w", fs.Arg(0), err)
	}

	if *out == "" {
		_, err = os.Stdout.Write(translated)
		return err
	}
	return os.WriteFile(*out, translated, 0o644)
}

// cmdRun implements "gotlang run": translate a package directory into a
// build cache, then "go run" it there, forwarding stdout/stderr/exit code.
func cmdRun(args []string) error {
	fs := flag.NewFlagSet("run", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		return err
	}
	path := "."
	if fs.NArg() == 1 {
		path = fs.Arg(0)
	} else if fs.NArg() > 1 {
		return fmt.Errorf("run: expected at most one path, got %d args", fs.NArg())
	}

	buildDir, err := prepareBuildDir(path)
	if err != nil {
		return err
	}
	return runGoToolchain(buildDir, "run", ".")
}

// cmdMecka implements "gotlang mecka": translate a package directory into a
// build cache, then "go build" it there.
func cmdMecka(args []string) error {
	fs := flag.NewFlagSet("mecka", flag.ExitOnError)
	out := fs.String("o", "", "output binary path")
	if err := fs.Parse(args); err != nil {
		return err
	}
	path := "."
	if fs.NArg() == 1 {
		path = fs.Arg(0)
	} else if fs.NArg() > 1 {
		return fmt.Errorf("mecka: expected at most one path, got %d args", fs.NArg())
	}

	buildDir, err := prepareBuildDir(path)
	if err != nil {
		return err
	}

	goArgs := []string{"build"}
	if *out != "" {
		abs, err := filepath.Abs(*out)
		if err != nil {
			return err
		}
		goArgs = append(goArgs, "-o", abs)
	}
	goArgs = append(goArgs, ".")
	return runGoToolchain(buildDir, goArgs...)
}
