package mdextensions

import (
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// TODO: Separate parsers at least by block/inline
type GochanInlineRefParser struct{}

type GochanGreentextParser struct{}

type GochanBlockRefParser struct{}

func (p *GochanInlineRefParser) Trigger() []byte {
	return []byte{'>'}
}

func (p *GochanGreentextParser) Trigger() []byte {
	return []byte{'>'}
}

func (p *GochanBlockRefParser) Trigger() []byte {
	return []byte{'>'}
}

func (p *GochanBlockRefParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()
	//fmt.Printf("[BlockRef] Checking line: %q\n", line)
	// Only handle standalone >> references at start of line (block level)
	if len(line) < 3 || !(line[0] == '>' && line[1] == '>') {
		//fmt.Printf("[BlockRef] Rejecting line: %q\n", line)
		return nil, parser.NoChildren
	}

	// Check if this is just digits after >> (simple post reference)
	// Let inline parser handle these instead
	i := 2
	for i < len(line) && line[i] >= '0' && line[i] <= '9' {
		i++
	}
	if i > 2 && (i == len(line) || line[i] == ' ' || line[i] == '\n') {
		//fmt.Printf("[BlockRef] Rejecting line: %q\n", line)
		return nil, parser.NoChildren // Let inline parser handle simple >>123
	}

	node := &PostRefBlock{}
	node.PostRefLines = [][]byte{append([]byte{}, line...)}
	reader.Advance(len(line))
	return node, parser.Close
}
func (p *GochanBlockRefParser) Continue(node ast.Node, reader text.Reader, pc parser.Context) parser.State {
	// single-line greentext only
	return parser.Close
}
func (p *GochanBlockRefParser) Close(node ast.Node, reader text.Reader, pc parser.Context) {}
func (p *GochanBlockRefParser) CanInterruptParagraph() bool                                { return true }
func (p *GochanBlockRefParser) CanAcceptIndentedLine() bool                                { return false }

func (p *GochanGreentextParser) Open(parent ast.Node, reader text.Reader, pc parser.Context) (ast.Node, parser.State) {
	line, _ := reader.PeekLine()
	//fmt.Printf("[Greentext] Checking line: %q\n", line)
	if len(line) == 0 || line[0] != '>' || (len(line) > 1 && line[1] == '>') {
		//fmt.Printf("[Greentext] Rejecting line: %q\n", line)
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
	//fmt.Printf("[InlineRef] Checking line: %q\n", line)
	// Post reference: (eg. >>123), can be anyuwhere in the
	if len(line) >= 3 && line[0] == '>' && line[1] == '>' {
		i := 2
		for i < len(line) && line[i] >= '0' && line[i] <= '9' {
			i++
		}
		if i == 2 {
			//fmt.Printf("[InlineRef] rejecting line: %q\n", line)
			return nil
		}

		block.Advance(i)

		return &PostRef{
			PostID: string(line[2:i]),
		}
	}

	return nil
}
