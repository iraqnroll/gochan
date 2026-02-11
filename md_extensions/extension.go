package mdextensions

import (
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
)

type GochanExtension struct {
}

func New() goldmark.Extender {
	return &GochanExtension{}
}

func (e *GochanExtension) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(
		parser.WithBlockParsers(
			util.Prioritized(&GochanGreentextParser{}, 70),
			util.Prioritized(parser.NewCodeBlockParser(), 100),
			util.Prioritized(parser.NewFencedCodeBlockParser(), 100),
		),
		parser.WithInlineParsers(
			util.Prioritized(&GochanInlineRefParser{}, 70),
		),
	)

	m.Renderer().AddOptions(
		renderer.WithNodeRenderers(
			util.Prioritized(&GochanHTMLRenderer{}, 500),
		),
	)
}
