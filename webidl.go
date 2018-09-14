package wudl

import (
	"go/scanner"
	"go/token"
)

var keywordList = []string{
	"any", "attribute", "ArrayBuffer", "boolean", "byte", "ByteString", "callback", "const", "creator", "DataView",
	"Date", "deleter", "dictionary", "DOMString", "double", "enum", "Error", "exception", "false", "float",
	"Float32Array", "Float64Array", "FrozenArray", "getter", "implements", "includes", "Infinity", "-Infinity", "inherit", "Int8Array",
	"Int16Array", "Int32Array", "interface", "iterable", "legacycaller", "legacyiterable", "long", "maplike", "mixin",
	"namespace", "NaN", "null", "object", "octet", "optional", "or", "partial", "Promise", "readonly", "record", "RegExp", "required",
	"sequence", "setlike", "setter", "short", "static", "stringifier", "true", "typedef",
	"Uint8Array", "Uint16Array", "Uint32Array", "Uint8ClampedArray", "unrestricted", "unsigned", "USVString",
	"void",
}

var keywords map[string]bool

func init() {
	keywords = make(map[string]bool)
	for _, v := range keywordList {
		keywords[v] = true
	}
}

// IsKeyword return true if the given identifier is a reserved webidl keyword.
func IsKeyword(ident string) bool {
	return keywords[ident]
}

type Parser struct {
	scanner *scanner.Scanner
}

func (p *Parser) Parse(fset *token.FileSet, filename string, src []byte) {
	var s scanner.Scanner
	file := fset.AddFile(filename, fset.Base(), len(src))
	s.Init(file, src, nil, scanner.ScanComments)
	p.scanner = &s
}
