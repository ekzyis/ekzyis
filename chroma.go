package main

import (
	"fmt"

	"github.com/alecthomas/chroma/v2"
	chromalexers "github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

// termLexer is a minimal Chroma lexer for ```term code blocks. It is very simple: it only
// categorizes lines starting with # as comments. Everything else is text.
var termLexer = chroma.MustNewLexer(
	&chroma.Config{
		Name:     "term",
		Aliases:  []string{"term"},
		EnsureNL: true,
	},
	func() chroma.Rules {
		return chroma.Rules{
			"root": {
				{Pattern: `#[^\n]*`, Type: chroma.CommentSingle}, // comment line
				{Pattern: `[^\n]+`, Type: chroma.Text},           // command ($ prompt) / output line
				{Pattern: `\n`, Type: chroma.Text},
			},
		}
	},
)
var _ = chromalexers.Register(termLexer)

func codeStyle() *chroma.Style {
	style, err := styles.Get("catppuccin-mocha").Builder().
	  // custom comment color; applies to all languages
		Add(chroma.CommentSingle, "#7f849c").
		Build()
	if err != nil {
		panic(fmt.Sprintf("failed to build code style: %v", err))
	}
	return style
}
