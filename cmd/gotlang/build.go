package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"gotlang/internal/transpile"
)

const buildCacheDirName = ".gotlang-build"

// prepareBuildDir translates every *.got file under srcDir into a mirrored
// *.go file inside srcDir/.gotlang-build, copies go.mod/go.sum alongside
// (so the generated package resolves imports correctly), and returns the
// build dir's absolute path.
func prepareBuildDir(srcDir string) (string, error) {
	srcDir, err := filepath.Abs(srcDir)
	if err != nil {
		return "", err
	}
	info, err := os.Stat(srcDir)
	if err != nil {
		return "", err
	}
	if !info.IsDir() {
		// Allow "gotlang run main.got" as shorthand for running the
		// single-file package containing it.
		srcDir = filepath.Dir(srcDir)
	}

	buildDir := filepath.Join(srcDir, buildCacheDirName)
	if err := os.RemoveAll(buildDir); err != nil {
		return "", fmt.Errorf("clearing build cache: %w", err)
	}
	if err := os.MkdirAll(buildDir, 0o755); err != nil {
		return "", err
	}

	var runtimeUsage transpile.RuntimeUsage
	pkgName := ""

	err = filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		if rel == buildCacheDirName {
			return filepath.SkipDir
		}
		if d.IsDir() {
			return nil
		}

		switch {
		case filepath.Ext(path) == ".got":
			src, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			runtimeUsage.Merge(transpile.NeedsRuntime(src))
			translated, err := transpile.Transpile(src)
			if err != nil {
				return fmt.Errorf("%s: %w", rel, err)
			}
			if pkgName == "" {
				if f, err := parser.ParseFile(token.NewFileSet(), "", translated, parser.PackageClauseOnly); err == nil {
					pkgName = f.Name.Name
				}
			}
			dst := filepath.Join(buildDir, rel[:len(rel)-len(".got")]+".go")
			return os.WriteFile(dst, translated, 0o644)
		case rel == "go.mod" || rel == "go.sum":
			return copyFile(path, filepath.Join(buildDir, rel))
		default:
			return nil
		}
	})
	if err != nil {
		return "", err
	}

	if runtimeUsage.Any() {
		runtimePath := filepath.Join(buildDir, "gotlang_runtime.go")
		if err := os.WriteFile(runtimePath, transpile.RuntimeSource(pkgName, runtimeUsage), 0o644); err != nil {
			return "", err
		}
	}

	return buildDir, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

// runGoToolchain runs "go <args...>" in dir, forwarding stdio and exit code.
func runGoToolchain(dir string, args ...string) error {
	cmd := exec.Command("go", args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		}
		return err
	}
	return nil
}
