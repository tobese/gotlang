package transpile

import (
	"go/format"
	"strings"
	"testing"
)

func TestTranspileKeywords(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "if/else",
			in:   "edugo x > 0 { } eller { }",
			want: "if x > 0 { } else { }",
		},
		{
			name: "for/break/continue",
			in:   "tjöta { sticc; körpå }",
			want: "for { break; continue }",
		},
		{
			name: "switch/case/default/fallthrough",
			in:   "kolla x { typ 1: barakör; annarsdå: }",
			want: "switch x { case 1: fallthrough; default: }",
		},
		{
			name: "func/return/var/const/type",
			in:   "fixa f() { grej x sorten int; fastpris y = 1; hemåt x }",
			want: "func f() { var x type int; const y = 1; return x }",
		},
		{
			name: "struct/interface/map/chan/select/range",
			in:   "byggsats S { }; mellanhand I { }; register[string]int; rör int; plocka { }; tjöta runtom m { }",
			want: "struct S { }; interface I { }; map[string]int; chan int; select { }; for range m { }",
		},
		{
			name: "defer/goto/package/import",
			in:   "paket main\nknöin \"fmt\"\nsedan f()\nschappa L",
			want: "package main\nimport \"fmt\"\ndefer f()\ngoto L",
		},
		{
			name: "go keyword is a no-op passthrough",
			in:   "go f()",
			want: "go f()",
		},
		{
			name: "int/string predeclared identifiers",
			in:   "fixa f(n nuffra, s tjöt) nuffra { hemåt n }",
			want: "func f(n int, s string) int { return n }",
		},
		{
			name: "string literal containing a keyword word is untouched",
			in:   `fixa f() { hemåt "edugo tjöta fixa" }`,
			want: `func f() { return "edugo tjöta fixa" }`,
		},
		{
			name: "comment containing a keyword word is untouched",
			in:   "fixa f() { } // edugo tjöta fixa",
			want: "func f() { } // edugo tjöta fixa",
		},
		{
			name: "identifier that merely contains a keyword word is untouched",
			in:   "grej edugoSomething = 1",
			want: "var edugoSomething = 1",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Transpile([]byte(tc.in))
			if err != nil {
				t.Fatalf("Transpile() error = %v", err)
			}
			if strings.TrimSpace(string(got)) != strings.TrimSpace(tc.want) {
				t.Errorf("Transpile() =\n%q\nwant\n%q", got, tc.want)
			}
		})
	}
}

func TestTranspileFizzbuzzIsValidGo(t *testing.T) {
	const src = `paket main

knöin "fmt"

fixa main() {
	tjöta i := 1; i <= 15; i++ {
		kolla {
		typ i%15 == 0:
			fmt.Println("FizzBuzz")
		typ i%3 == 0:
			fmt.Println("Fizz")
		typ i%5 == 0:
			fmt.Println("Buzz")
		annarsdå:
			fmt.Println(i)
		}
	}
}
`
	out, err := Transpile([]byte(src))
	if err != nil {
		t.Fatalf("Transpile() error = %v", err)
	}
	if _, err := format.Source(out); err != nil {
		t.Fatalf("transpiled output is not valid Go: %v\n---\n%s", err, out)
	}
}

func TestNeedsRuntime(t *testing.T) {
	cases := []struct {
		name string
		src  string
		want RuntimeUsage
	}{
		{"calls hörru", `fixa main() { namn := hörru() }`, RuntimeUsage{Input: true}},
		{"calls Tjöta", `fixa main() { Tjöta("hej") }`, RuntimeUsage{Print: true}},
		{"calls both", `fixa main() { Tjöta(hörru()) }`, RuntimeUsage{Input: true, Print: true}},
		{"mentions hörru only in a string", `fixa main() { hemåt "hörru" }`, RuntimeUsage{}},
		{"mentions hörru only in a comment", "fixa main() { } // hörru", RuntimeUsage{}},
		{"identifier merely containing the word", `grej hörrugrej = 1`, RuntimeUsage{}},
		{"lowercase tjöta keyword is not the Tjöta helper", `tjöta { }`, RuntimeUsage{}},
		{"no mention at all", `fixa main() {}`, RuntimeUsage{}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := NeedsRuntime([]byte(tc.src)); got != tc.want {
				t.Errorf("NeedsRuntime(%q) = %+v, want %+v", tc.src, got, tc.want)
			}
		})
	}
}

func TestRuntimeSourceIsValidGo(t *testing.T) {
	cases := []RuntimeUsage{
		{Input: true},
		{Print: true},
		{Input: true, Print: true},
	}
	for _, usage := range cases {
		out := RuntimeSource("main", usage)
		if _, err := format.Source(out); err != nil {
			t.Fatalf("RuntimeSource(%+v) output is not valid Go: %v\n---\n%s", usage, err, out)
		}
	}
}
