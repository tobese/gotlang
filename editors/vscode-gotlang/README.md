# Gotlang for VS Code

Syntax highlighting for `.got` files, plus a run command — no build step,
not published to the Marketplace.

## What it does

- Registers the `gotlang` language for `.got` files, with highlighting for
  the full Gotlang keyword table (`edugo`, `eller`, `tjöta`, `fixa`,
  `byggsats`, ...), strings, comments, numbers, and the `hörru()`/`Tjöta()`
  builtins.
- Adds a **Gotlang: Run File** command (and a ▶ button in the editor
  title bar for `.got` files) that runs the current file through
  `gotlang run` in an integrated terminal.
  - If the file lives inside a checkout of the gotlang repo (i.e. a
    `cmd/gotlang/main.go` is found by walking up from the file), it runs
    the CLI straight from source via `go run`, so nothing needs to be
    built or installed first.
  - Otherwise it falls back to assuming a `gotlang` binary is on `PATH`.

## Trying it

No npm install needed — it's plain JS with no dependencies.

1. Open this `editors/vscode-gotlang` folder in VS Code.
2. Press `F5` (Run → Start Debugging) to launch an Extension Development
   Host window with the extension loaded.
3. In that window, open a `.got` file (e.g. `examples/hello/hello.got` from
   the repo root) and confirm it's syntax-highlighted.
4. Run it via the Command Palette (**Gotlang: Run File**) or the ▶ button.

## Known limitations

This covers highlighting and running a file — it does not provide
IntelliSense (autocomplete, hover docs, go-to-definition, diagnostics).
That would mean either a dedicated language server or forwarding to
`gopls` against the transpiled `.go` output with position-mapping back to
the `.got` source — a meaningfully bigger lift, not attempted here.
