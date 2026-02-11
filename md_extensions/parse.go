package mdextensions

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// TODO: Separate parsers at least by block/inline
type GochanInlineRefParser struct{}

type GochanGreentextParser struct{}

func (p *GochanInlineRefParser) Trigger() []byte {
	return []byte{'>'}
}

func (p *GochanGreentextParser) Trigger() []byte {
	return []byte{'>'}
}

func (p *GochanGreentextParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()
	if len(line) == 0 || line[0] != '>' || (len(line) > 1 && line[1] == '>') {
		return nil, parser.NoChildren
	}

	node := &Greentext{}
	node.GreentextLines = [][]byte{append([]byte{}, line...)}
	reader.Advance(len(line))
	return node, parser.HasChildren
}
func (p *GochanGreentextParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	// single-line greentext only
	return parser.Close
}
func (p *GochanGreentextParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {}
func (p *GochanGreentextParser) CanInterruptParagraph() bool                                { return true }
func (p *GochanGreentextParser) CanAcceptIndentedLine() bool                                { return false }

func (p *GochanInlineRefParser) Parse(parent ast.Node, block text.Reader, pc parser.Context) ast.Node {
	line, _ := block.PeekLine()
	if len(line) >= 3 && line[0] == '>' && line[1] == '>' {
		i := 2
		for i < len(line) && line[i] >= '0' && line[i] <= '9' {
			i++
		}
		if i == 2 {
			return nil
		}

		block.Advance(i)

		return &PostRef{
			PostID: string(line[2:i]),
		}
	}

	return nil
}
