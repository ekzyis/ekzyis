package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sourcegraph/syntaxhighlight"
)

func SyntaxHighlighting(element *goquery.Selection) {
	if element.HasClass("language-diff") {
		// syntaxhighlight does not support diff so we run our custom code in that case
		text := strings.Split(element.Text(), "\n")
		p1 := regexp.MustCompile(`^\+ `)
		p2 := regexp.MustCompile(`^- `)
		for i, line := range text {
			if p1.MatchString(line) {
				text[i] = fmt.Sprintf("<span class=\"diff-add\">%s</span>", line)
			}
			if p2.MatchString(line) {
				text[i] = fmt.Sprintf("<span class=\"diff-remove\">%s</span>", line)
			}
		}
		element.SetHtml(strings.Join(text, "\n"))
		return
	}
	formatted, err := syntaxhighlight.AsHTML([]byte(element.Text()))
	if err != nil {
		panic(err)
	}
	element.SetHtml(string(formatted))
}
