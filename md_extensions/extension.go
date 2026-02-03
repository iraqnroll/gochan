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
			util.Prioritized(&GochanBlockRefParser{}, 50),
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

// func (e *GochanExtension) Extend(m goldmark.Markdown) {
// 	blockParsers := []util.PrioritizedValue{
// 		util.Prioritized(parser.NewParagraphParser(), 400),
// 		util.Prioritized(&GochanBlockRefParser{}, 200),
// 		util.Prioritized(&GochanGreentextParser{}, 100),
// 	}

// 	m.Parser().AddOptions(
// 		parser.WithBlockParsers(blockParsers...),
// 		parser.WithInlineParsers(
// 			util.Prioritized(&GochanInlineRefParser{}, 300),
// 		),
// 	)

// 	m.Renderer().AddOptions(
// 		renderer.WithNodeRenderers(
// 			util.Prioritized(&GochanHTMLRenderer{}, 500),
// 		),
// 	)
// }
