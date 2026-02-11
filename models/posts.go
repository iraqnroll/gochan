package models

import (
	"bytes"
	"sync"

	mdextensions "github.com/iraqnroll/gochan/md_extensions"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

type PostDto struct {
	Id            int    `json:"id"`
	ThreadId      int    `json:"thread_id"`
	Identifier    string `json:"identifier"`
	Content       string `json:"content"`
	PostTimestamp string `json:"post_timestamp"`
	IsOP          bool   `json:"is_op"`
	HasMedia      string
}

type RecentPostsDto struct {
	Board_uri      string `json:"board_uri"`
	Board_name     string `json:"board_name"`
	Thread_id      int    `json:"thread_id"`
	Thread_topic   string `json:"thread_topic"`
	Post_id        int    `json:"post_id"`
	Post_ident     string `json:"post_ident"`
	Post_content   string `json:"post_content"`
	Post_timestamp string `json:"post_timestamp"`
	HasMedia       string
}

var getMarkdownParser = sync.OnceValue(func() goldmark.Markdown {
	p := parser.NewParser(
		parser.WithBlockParsers(
			util.Prioritized(&mdextensions.GochanBlockRefParser{}, 50),
			util.Prioritized(&mdextensions.GochanGreentextParser{}, 50),
			util.Prioritized(parser.NewCodeBlockParser(), 60),
			util.Prioritized(parser.NewFencedCodeBlockParser(), 60),
			util.Prioritized(parser.NewParagraphParser(), 100),
		),
		parser.WithInlineParsers(
			util.Prioritized(parser.NewCodeSpanParser(), 70),
			util.Prioritized(&mdextensions.GochanInlineRefParser{}, 100),
		),
	)

	r := renderer.NewRenderer(
		renderer.WithNodeRenderers(
			util.Prioritized(html.NewRenderer(), 100),
			util.Prioritized(&mdextensions.GochanHTMLRenderer{}, 500),
		),
	)

	return goldmark.New(
		goldmark.WithParser(p),
		goldmark.WithRenderer(r),
		goldmark.WithExtensions(mdextensions.New()),
	)
})

func RenderSafeMarkdown(md string, pPol *bluemonday.Policy) (string, error) {
	var buf bytes.Buffer
	if err := getMarkdownParser().Convert([]byte(md), &buf); err != nil {
		return "", err
	}
	safe := pPol.SanitizeBytes(buf.Bytes())

	return string(safe), nil
}
