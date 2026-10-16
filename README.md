# Gotlang

Gotlang is Go with its reserved keywords renamed to Göteborg slang and dry
puns. Nothing about the language itself changes — `edugo x > 0 { }` is just
`if x > 0 { }` wearing a funnier hat. Under the hood, `gotlang` is a lexical
preprocessor: it swaps Gotlang keyword tokens back to real Go keywords using
Go's own `go/scanner`, then hands the result to the real `go` toolchain.

Source files use the `.got` extension.

## Keyword table

| Gotlang | Go | Why |
|---|---|---|
| `paket` | `package` | direct Swedish cognate |
| `knöin` | `import` | |
| `fixa` | `func` | "to fix/sort/make" |
| `edugo` | `if` | "e du go?" — "you good?" |
| `eller` | `else` | Swedish "or" |
| `tjöta` | `for` | "to go on and on" — fits looping |
| `hemåt` | `return` | "homeward" |
| `sticc` | `break` | "jag sticker" — I'm bailing |
| `körpå` | `continue` | "köra på" — carry on |
| `kolla` | `switch` | "take a look" |
| `typ` | `case` | GBG filler slang "like/such as" |
| `annarsdå` | `default` | "otherwise then" |
| `barakör` | `fallthrough` | "bara kör" — just keep going |
| `byggsats` | `struct` | "build kit" |
| `mellanhand` | `interface` | "middleman" |
| `register` | `map` | |
| `runtom` | `range` | "all around" |
| `plocka` | `select` | "to pick" |
| `rör` | `chan` | "pipe" |
| `sedan` | `defer` | "then/later" |
| `schappa` | `goto` | GBG slang for taking off fast |
| `fastpris` | `const` | "fixed price" |
| `grej` | `var` | "a thing" |
| `sorten` | `type` | "the sort/kind" |
| `go` | `go` | unchanged |
| `nuffra` | `int` | |
| `tjöt` | `string` | |

Beyond the ~25 reserved keywords, two predeclared type identifiers are
renamed too (`nuffra`, `tjöt`). Everything else in builtins and stdlib
(`fmt.Println`, `true`, `false`, `make`, `len`, ...) is untouched.

## hörru() / Tjöta() — runtime helpers

Two words in Gotlang aren't token swaps like the rest of the table — they're
real helper functions that `gotlang run`/`gotlang build` generate and inject
into the build (as `gotlang_runtime.go`) whenever a `.got` file calls them.
`gotlang translate` only ever does token swapping on a single file, so a
file calling either will translate to Go that references an undefined
function — use `run`/`build` for a complete, runnable program.

- `hörru()` reads a line from stdin, the way Python's `input()` does —
  something Go has no builtin for.
- `Tjöta(...)` prints its arguments followed by a newline, the way
  `fmt.Println` does — "tjöta" being GBG slang for yapping on, fitting for
  a function whose whole job is to talk.

```
grej namn tjöt = hörru()
Tjöta("tjena", namn)
```

See `examples/samtal` for a full program that reads a name and a number,
and `src/hallå.got` for a minimal `Tjöta()` hello-world.

## Usage

```sh
go build ./cmd/gotlang   # or: go run ./cmd/gotlang ...

gotlang translate <file.got> [-o out.go]   # emit the translated Go source
gotlang run <path>                         # translate a package and run it
gotlang build <path> [-o binary]           # translate a package and build it
```

Flags must come before the positional path, same as the `go` tool itself.

`run`/`build` translate every `*.got` file under `<path>` into a
`.gotlang-build/` cache directory next to it (copying `go.mod`/`go.sum`
alongside so imports resolve), then invoke the real `go` toolchain there.
Each program needs its own directory with its own `go.mod`, same as any Go
package.

Try it:

```sh
go run ./cmd/gotlang run examples/hello/hello.got
go run ./cmd/gotlang run examples/fizzbuzz/fizzbuzz.got
go run ./cmd/gotlang run examples/samtal/samtal.got
go run ./cmd/gotlang run src/hallå.got
```

`src/` holds standalone `.got` programs (outside the `examples/` demos) —
start with `src/hallå.got`, a hello-world written in Gotlang:

```go
paket main

fixa main() {
	Tjöta("Hallå, Göteborg!")
}
```

## Editor support

`editors/vscode-gotlang` is a local (unpublished) VS Code extension: syntax
highlighting for `.got` files plus a **Gotlang: Run File** command. See its
own README for how to load it in an Extension Development Host. It doesn't
provide IntelliSense (autocomplete/diagnostics/go-to-definition).

## Known limitation

This is pure lexical substitution, not a scope-aware rename: a Gotlang
keyword word is swapped wherever it appears as a bare identifier — including
if you name a struct field or method the same as a keyword word (e.g. a
field literally called `fixa`). This mirrors how real reserved words behave
in Go (you can't name a variable `if` either), so it's expected rather than
a bug.
