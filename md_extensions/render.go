package mdextensions

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type GochanHTMLRenderer struct{}

func (r *GochanHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(KindPostRef, r.renderPostRef)
	reg.Register(KindGreentext, r.renderGreentext)
}

func (r *GochanHTMLRenderer) renderPostRef(
	w util.BufWriter,
	source []byte,
	node ast.Node,
	entering bool,
) (ast.WalkStatus, error) {

	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*PostRef)

	w.WriteString(`<a href="#`)
	w.WriteString(n.PostID)
	w.WriteString(`"><span class="post-ref">&gt;&gt;`)
	w.WriteString(n.PostID)
	w.WriteString(`</span></a>`)

	return ast.WalkSkipChildren, nil
}

func (r *GochanHTMLRenderer) renderGreentext(
	w util.BufWriter,
	source []byte,
	node ast.Node,
	entering bool,
) (ast.WalkStatus, error) {

	if !entering {
		return ast.WalkContinue, nil
	}

	n := node.(*Greentext)

	w.WriteString(`<span class="greentext">`)
	for _, line := range n.GreentextLines {
		w.Write(line)
	}
	w.WriteString(`</span>`)

	return ast.WalkSkipChildren, nil
}
