package safechange

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
	"reflect"
	"sort"

	"github.com/go-toolsmith/astcopy"
)

// AstFile is a wrapper for *ast.File which extends with additional methods.
type AstFile struct {
	*ast.File
}

// NewAstFile wraps an *ast.File to AstFile.
func NewAstFile(file *ast.File) *AstFile {
	return &AstFile{
		File: file,
	}
}

// Equals returns true if ASTs has exactly the same content.
func (file *AstFile) Equals(cmp *AstFile) bool {
	return reflect.DeepEqual(file, cmp)
}

// Copy returns a deep copy.
func (file *AstFile) Copy() *AstFile {
	return NewAstFile(astcopy.File(file.File))
}

// EquivalentTo returns true if ASTs are expected to be compiled into the same
// machine code.
func (file *AstFile) EquivalentTo(cmp *AstFile) bool {
	// copy ASTs, remove everything non-essential (like comments) and then compare

	fileCopy := file.Copy()
	fileCopy.stripNonEssentials()
	cmpCopy := cmp.Copy()
	cmpCopy.stripNonEssentials()

	// format ("go fmt")

	var outA bytes.Buffer
	{
		err := format.Node(&outA, token.NewFileSet(), fileCopy.File)
		assertNoError(err)
	}

	var outB bytes.Buffer
	{
		err := format.Node(&outB, token.NewFileSet(), cmpCopy.File)
		assertNoError(err)
	}

	return outA.String() == outB.String()
}

func assertNoError(err error) {
	if err != nil {
		panic(err)
	}
}

func (file *AstFile) stripNonEssentials() {

	// remove all comments
	ast.Inspect(file.File, func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.CommentGroup, *ast.Comment:
			panic("should not have happened")

		case *ast.Field:
			node.Doc = nil
			node.Comment = nil

		case *ast.ValueSpec:
			node.Doc = nil
			node.Comment = nil

		case *ast.TypeSpec:
			node.Doc = nil
			node.Comment = nil

		case *ast.GenDecl:
			node.Doc = nil

		case *ast.FuncDecl:
			node.Doc = nil

		case *ast.File:
			node.Doc = nil
			node.Comments = nil

		case *ast.ImportSpec:
			node.Doc = nil
		}

		return true
	})

	// group all imports together and sort
	var decls []ast.Decl
	importsDecl := &ast.GenDecl{
		Tok: token.IMPORT,
	}
	decls = append(decls, importsDecl)
	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			isImportOnly := true
			for _, spec := range decl.Specs {
				if spec, ok := spec.(*ast.ImportSpec); ok {
					importsDecl.Specs = append(importsDecl.Specs, spec)
				} else {
					isImportOnly = false
				}
			}
			if isImportOnly {
				continue
			}
		}
		decls = append(decls, decl)
	}
	sort.Slice(importsDecl.Specs, func(i, j int) bool {
		return importsDecl.Specs[i].(*ast.ImportSpec).Path.Value < importsDecl.Specs[j].(*ast.ImportSpec).Path.Value
	})
	file.Decls = decls
}
