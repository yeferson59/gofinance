package invariants

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"
)

// TestErrorReturningFunctionsAreGuarded keeps the error contract from eroding:
// every exported function or method in finance/* that returns an error must
// open with `defer guard.Recover(&err)`, so an out-of-range input that trips
// one of the decimal engine's panicking helpers becomes that function's error
// rather than a crash. A new function that forgets the guard fails here,
// before a fuzzer or a user finds the input that panics it.
func TestErrorReturningFunctionsAreGuarded(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("..", "finance", "*", "*.go"))
	if err != nil {
		t.Fatal(err)
	}

	fset := token.NewFileSet()
	checked := 0

	for _, path := range files {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}

		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			t.Fatal(err)
		}

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil || !fn.Name.IsExported() || !returnsError(fn) {
				continue
			}

			checked++

			if !opensWithGuard(fn) {
				t.Errorf("%s: %s returns an error but does not start with `defer guard.Recover(&err)`",
					fset.Position(fn.Pos()), fn.Name.Name)
			}
		}
	}

	if checked == 0 {
		t.Fatal("no error-returning functions found; is the glob still right?")
	}
}

func returnsError(fn *ast.FuncDecl) bool {
	results := fn.Type.Results
	if results == nil || len(results.List) == 0 {
		return false
	}

	id, ok := results.List[len(results.List)-1].Type.(*ast.Ident)

	return ok && id.Name == "error"
}

func opensWithGuard(fn *ast.FuncDecl) bool {
	if len(fn.Body.List) == 0 {
		return false
	}

	deferStmt, ok := fn.Body.List[0].(*ast.DeferStmt)
	if !ok {
		return false
	}

	sel, ok := deferStmt.Call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	pkg, ok := sel.X.(*ast.Ident)

	return ok && pkg.Name == "guard" && sel.Sel.Name == "Recover"
}
