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

type Node interface {
	Pos() token.Pos
	End() token.Pos
}

type Token struct {
	Pos  token.Pos
	Tok  token.Token
	Text string
}

type Parser struct {
	scanner *scanner.Scanner
	err     scanner.ErrorList
	nodes   []Node
	current Node
	rewind  *Token
}

type State uint

const (
	OpenInterface = 1 << iota
	EndInterface
	OpenExtendedAttribute
	EndExtendedAttribute
)

func (p *Parser) Parse(fset *token.FileSet, filename string, src []byte) error {
	var s scanner.Scanner
	file := fset.AddFile(filename, fset.Base(), len(src))
	s.Init(file, src, p.handleScanError, scanner.ScanComments)
	p.scanner = &s
	var state State
	var nest int
	for {
		if p.err != nil {
			return p.err
		}
		tok := p.next()
		if tok.Tok == token.EOF {
			return nil
		}
		switch tok.Tok {
		case token.LBRACK:
			state = OpenExtendedAttribute
			nest++
			x := &ExtendedAttributeList{}
			x.pos = tok.Pos
			p.current = x
		case token.RBRACK:
			state = EndExtendedAttribute
			nest++
		case token.COMMA:
			if state == OpenExtendedAttribute {
				continue
			}
		}
		switch state {
		case OpenExtendedAttribute:
			p.parseAttribute()
		case EndExtendedAttribute:
			state = 0
			p.nodes = append(p.nodes, p.current)
			p.current = nil
		}
	}
}

func (p *Parser) peek() *Token {
	x := p.next()
	p.revind1(x)
	return x
}

func (p *Parser) parseAttribute() {
	x := p.next()
	switch x.Tok {
	case token.IDENT:
		next := p.peek()
		switch next.Tok {
		case token.COMMA, token.RBRACK:
			n := ExtendedAttributeNoArgs{
				Indent: x.Text,
			}
			n.pos = x.Pos
			n.end = next.Pos
			if v, ok := p.current.(*ExtendedAttributeList); ok {
				v.List = append(v.List, n)
			}
		}
	}
}

type position struct {
	pos token.Pos
	end token.Pos
}

func (e position) Pos() token.Pos {
	return e.pos
}
func (e position) End() token.Pos {
	return e.end
}

type ExtendedAttributeList struct {
	position
	List []Node
}

type ExtendedAttributeNoArgs struct {
	position
	Indent string
}

func (p *Parser) next() *Token {
	if p.rewind != nil {
		// If we callend rewind1, we return the rewind token and skip scanning. We are
		// setting rewind to nil so next p.next call will do a scan.
		x := p.rewind
		p.rewind = nil
		return x
	}
	pos, tok, lit := p.scanner.Scan()
	return &Token{Pos: pos, Tok: tok, Text: lit}
}

func (p *Parser) revind1(tok *Token) {
	p.rewind = tok
}

func (p *Parser) handleScanError(pos token.Position, msg string) {
	p.err = append(p.err, &scanner.Error{
		Pos: pos,
		Msg: msg,
	})
}
