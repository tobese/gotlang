// Package transpile converts Gotlang source (.got) into plain Go source by
// swapping Göteborg-humor words back to the real Go keyword or predeclared
// identifier they stand in for.
package transpile

// Keywords maps each Gotlang word to the real Go keyword or predeclared
// identifier it stands in for. This is the single source of truth for the
// language's vocabulary.
var Keywords = map[string]string{
	"paket":      "package",
	"nuffra":     "int",
	"tjöt":       "string",
	"knöin":      "import",
	"fixa":       "func",
	"edugo":      "if",
	"eller":      "else",
	"tjöta":      "for",
	"hemåt":      "return",
	"sticc":      "break",
	"körpå":      "continue",
	"kolla":      "switch",
	"typ":        "case",
	"annarsdå":   "default",
	"barakör":    "fallthrough",
	"byggsats":   "struct",
	"mellanhand": "interface",
	"register":   "map",
	"runtom":     "range",
	"plocka":     "select",
	"rör":        "chan",
	"sedan":      "defer",
	"schappa":    "goto",
	"fastpris":   "const",
	"grej":       "var",
	"sorten":     "type",
	"go":         "go",
}
