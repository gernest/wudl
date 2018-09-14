package wudl

import (
	"go/token"
	"reflect"
	"testing"
)

func TestParse_extendedAttributes(t *testing.T) {
	parser := &Parser{}
	fset := token.NewFileSet()

	t.Run("forms", func(ts *testing.T) {
		ts.Run("takes no arguments", func(ts *testing.T) {
			nodes, err := parser.Parse(fset, "", []byte("[Replaceable]"))
			if err != nil {
				ts.Fatal(err)
			}
			ls := nodes[0].(*ExtendedAttributeList).List[0].(*ExtendedAttributeNoArgs)
			expect := "Replaceable"
			if ls.Ident != expect {
				t.Errorf("expected %s got %s", expect, ls.Ident)
			}
		})

		ts.Run("takes an argument list", func(ts *testing.T) {
			nodes, err := parser.Parse(fset, "", []byte("[Constructor(double x, double y)]"))
			if err != nil {
				ts.Fatal(err)
			}
			ls := nodes[0].(*ExtendedAttributeList).List[0].(*ExtendedAttributeArgList)
			expect := "Constructor"
			if ls.Ident != expect {
				t.Errorf("expected %s got %s", expect, ls.Ident)
			}
			args := [][]string{
				[]string{"double", "x"},
				[]string{"double", "y"},
			}
			if !reflect.DeepEqual(ls.Args, args) {
				t.Error("expected args to match")
			}

		})
	})

	// parser.Parse(fset, "", []byte("[NamedConstructor=Image(DOMString src)]"))
	// t.Error(pretty.Sprint(parser.nodes))

	// parser.Parse(fset, "", []byte("[PutForwards=name]"))
	// t.Error(pretty.Sprint(parser.nodes))

	// parser.Parse(fset, "", []byte("[Exposed=(Window,Worker)]"))
	// t.Error(pretty.Sprint(parser.nodes))
}
