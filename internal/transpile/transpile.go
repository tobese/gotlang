package transpile

import (
	"bytes"
	"go/scanner"
	"go/token"
)

// Transpile converts Gotlang source into Go source by swapping identifier
// tokens that match a Gotlang keyword for the real Go keyword they stand
// in for. String literals and comments are never inspected, since the
// underlying scanner only reports their content as STRING/COMMENT tokens,
// not IDENT — so a Gotlang keyword word appearing inside one is left as-is.
func Transpile(src []byte) ([]byte, error) {
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))

	var errs scanner.ErrorList
	var s scanner.Scanner
	s.Init(file, src, func(pos token.Position, msg string) {
		errs = append(errs, &scanner.Error{Pos: pos, Msg: msg})
	}, scanner.ScanComments)

	var out bytes.Buffer
	prev := 0
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		if tok != token.IDENT {
			continue
		}
		goKeyword, ok := Keywords[lit]
		if !ok || goKeyword == lit {
			continue
		}
		offset := fset.Position(pos).Offset
		out.Write(src[prev:offset])
		out.WriteString(goKeyword)
		prev = offset + len(lit)
	}
	if len(errs) > 0 {
		return nil, errs.Err()
	}

	out.Write(src[prev:])
	return out.Bytes(), nil
}
