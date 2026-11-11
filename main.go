package main

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/typescript/typescript"
)

func identifierAtCursor(file string, row, column int) (*CallInfo, error) {
	source, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	parser := sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(typescript.GetLanguage())

	tree, err := parser.ParseCtx(context.Background(), nil, source)
	if err != nil {
		return nil, fmt.Errorf("failed to parse %s", file)
	}
	defer tree.Close()

	point := sitter.Point{
		Row:    uint32(row - 1),
		Column: uint32(column - 1),
	}

	node := tree.RootNode().NamedDescendantForPointRange(point, point)

	info := callAtCursor(node, source)

	return info, nil
}

type CallInfo struct {
	Object string
	Method string
}

func callAtCursor(node *sitter.Node, source []byte) *CallInfo {
	for node != nil {
		if node.Type() == "member_expression" {
			object := node.ChildByFieldName("object")
			property := node.ChildByFieldName("property")

			if object != nil && property != nil {
				return &CallInfo{
					Object: object.Content(source),
					Method: property.Content(source),
				}
			}
		}

		node = node.Parent()
	}

	return nil
}

func findGoMethod(root, receiver, method string) (string, int, error) {
	fset := token.NewFileSet()

	var foundFile string
	var foundLine int

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			if path != root && (d.Name() == "vendor" || d.Name() == "node_modules" || d.Name() == ".git") {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}

		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return nil
		}

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Recv == nil || fn.Name.Name != method {
				continue
			}

			recv := fn.Recv.List[0].Type

			if star, ok := recv.(*ast.StarExpr); ok {
				recv = star.X
			}

			ident, ok := recv.(*ast.Ident)
			if !ok || ident.Name != receiver {
				continue
			}

			pos := fset.Position(fn.Pos())
			foundFile = path
			foundLine = pos.Line
			return filepath.SkipAll
		}

		return nil
	})

	return foundFile, foundLine, err
}

func main() {
	row, _ := strconv.Atoi(os.Args[2])
	column, _ := strconv.Atoi(os.Args[3])

	info, err := identifierAtCursor(os.Args[1], row, column)
	if err != nil {
		panic(err)
	}

	projectRoot := os.Args[4]
	file, line, err := findGoMethod(projectRoot, info.Object, info.Method)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command("zed", fmt.Sprintf("%s:%d", file, line))

	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
