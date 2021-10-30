package safechange

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAstFileEquivalentTo(t *testing.T) {
	fileDataA := `package pkg
import (
	"fmt"

	"io"
)

import "io/ioutil"

import "errors"

func main() {
	// some comment
}
`

	fileDataB := `package pkg
import (
	"errors"
	"fmt"
	"io/ioutil"
	"io"
)

func main() {
}
`

	fileDataC := `package pkg
import (
	"errors"
	"fmt"
	"io/ioutil"
	"io"
)

func main() {
	fmt.Println("hello!")
}
`
	fSetA := token.NewFileSet()
	astA, err := parser.ParseFile(fSetA, "", fileDataA, 0)
	require.NoError(t, err)
	fileA := NewAstFile(astA)

	fSetB := token.NewFileSet()
	astB, err := parser.ParseFile(fSetB, "", fileDataB, 0)
	require.NoError(t, err)
	fileB := NewAstFile(astB)

	fSetC := token.NewFileSet()
	astC, err := parser.ParseFile(fSetC, "", fileDataC, 0)
	require.NoError(t, err)
	fileC := NewAstFile(astC)

	require.True(t, fileA.EquivalentTo(fileB))
	require.False(t, fileA.EquivalentTo(fileC))
}
