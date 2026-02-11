package mdextensions

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
)

var (
	KindPostRef   = ast.NewNodeKind("PostRef")
	KindGreentext = ast.NewNodeKind("Greentext")
)

type PostRef struct {
	ast.BaseInline
	PostID string
}

type Greentext struct {
	ast.BaseBlock
	GreentextLines [][]byte
}

func (n *PostRef) Kind() ast.NodeKind { return KindPostRef }

func (n *Greentext) Kind() ast.NodeKind { return KindGreentext }

// Dump dumps the contents of Node to stdout for debugging.
func (n *PostRef) Dump(src []byte, level int) {
	ast.DumpHelper(n, src, level, map[string]string{
		"PostID": string(n.PostID),
	}, nil)
}

func (n *Greentext) Dump(src []byte, level int) {
	joined := bytes.Join(n.GreentextLines, []byte("\n"))
	ast.DumpHelper(n, src, level, map[string]string{
		"GreentextLines": string(joined),
	}, nil)
}
