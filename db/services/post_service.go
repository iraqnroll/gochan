package services

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/iraqnroll/gochan/config"
	"github.com/iraqnroll/gochan/db/models"
	mdextensions "github.com/iraqnroll/gochan/md_extensions"
	"github.com/iraqnroll/gochan/rand"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

type PostRepository interface {
	GetAllByThread(thread_id int, for_mod bool) ([]models.PostDto, error)
	CreateNew(thread_id int, identifier, content, fingerprint, tripcode string, is_op bool) (models.PostDto, error)
	GetMostRecent(num_of_posts int) ([]models.RecentPostsDto, error)
	UpdateAttachedMedia(post_id int, attached_media, original_media string) error
	SoftDeletePost(post_id int) error
	RemoveSoftDeleteFromPost(post_id int) error
}

type PostService struct {
	PostRepo        PostRepository
	PostPolicy      *bluemonday.Policy
	MdParser        goldmark.Markdown
	FingerprintSalt string
	TripcodeSalt    string
}

func initPostPolicies() *bluemonday.Policy {
	policy := bluemonday.UGCPolicy()
	policy.AllowStandardURLs()
	policy.AllowElements("a", "pre", "code")
	policy.AllowAttrs("class").OnElements("span")
	policy.AllowAttrs("href").OnElements("a")

	return policy
}

func initGoldmarkParser() goldmark.Markdown {
	p := parser.NewParser(
		parser.WithBlockParsers(
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
}

func NewPostService(repo PostRepository, fprintSalt, tripcodeSalt string) *PostService {
	return &PostService{
		PostRepo:        repo,
		FingerprintSalt: fprintSalt,
		TripcodeSalt:    tripcodeSalt,
		PostPolicy:      initPostPolicies(),
		MdParser:        initGoldmarkParser(),
	}
}

// Creates a post in a specific thread
func (ps *PostService) CreatePost(thread_id int, identifier, content, fingerprint, tripcode string, is_op bool) (models.PostDto, error) {
	post, err := ps.PostRepo.CreateNew(thread_id, identifier, content, fingerprint, tripcode, is_op)
	if err != nil {
		return post, fmt.Errorf("PostService.CreateReply failed : %w", err)
	}

	return post, nil
}

// Fetches all posts of the specified thread
func (ps *PostService) GetThreadPosts(thread_id int, for_mod bool) ([]models.PostDto, error) {
	posts, err := ps.PostRepo.GetAllByThread(thread_id, for_mod)
	if err != nil {
		return nil, fmt.Errorf("PostService.GetThreadPosts failed : %w", err)
	}

	return posts, nil
}

func (ps *PostService) GetMostRecent(num_of_posts int) ([]models.RecentPostsDto, error) {
	posts, err := ps.PostRepo.GetMostRecent(num_of_posts)
	if err != nil {
		return nil, fmt.Errorf("PostService.GetMostRecent failed : %s", err)
	}

	return posts, nil
}

// TODO: Deal with attached post media, if we hide the post, we want to hide the attached content from /static as well
func (ps *PostService) SoftDeletePost(post_id int) error {
	return ps.PostRepo.SoftDeletePost(post_id)
}

func (ps *PostService) RemoveSoftDeleteFromPost(post_id int) error {
	return ps.PostRepo.RemoveSoftDeleteFromPost(post_id)
}

func (ps *PostService) UpdateAttachedMedia(post_id int, attached_media, original_media string) error {
	return ps.PostRepo.UpdateAttachedMedia(post_id, attached_media, original_media)
}

func (ps *PostService) GenerateFingerprint(ip string) string {
	return rand.GenerateSha256Hash(ip, ps.FingerprintSalt)
}

func (ps *PostService) RenderSafeMarkdown(md string) (string, error) {
	var buf bytes.Buffer
	if err := ps.MdParser.Convert([]byte(md), &buf); err != nil {
		return "", err
	}
	safe := ps.PostPolicy.SanitizeBytes(buf.Bytes())

	return string(safe), nil
}

// TODO: make the tripcode prefix configurable
// TODO: add validation to tripcodes (min. length, allowed chars..)
func (ps *PostService) GetTripcodeHash(name string, authenticated bool) (string, string) {
	ident, pw, found := strings.Cut(name, "#")
	if !found || pw == "" {
		return ident, ""
	}
	tripcode := rand.GenerateSha256Hash(pw, ps.TripcodeSalt)

	if authenticated {
		if auth := config.AuthenticatedTripcode(pw); auth != "" {
			tripcode = auth
		}
	}
	return ident, tripcode
}
