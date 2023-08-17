package main

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/PuerkitoBio/goquery"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
)

var (
	MarkdownToHtmlFlags = html.CommonFlags | html.HrefTargetBlank
)

type MarkdownPost struct {
	// file system path
	FsPath string
	// markdown content
	Content []byte
	// args parsed from markdown
	Title       string
	Date        string
	ReadingTime string
	Sats        int
}

func NewMarkdownPost(path string) *MarkdownPost {
	post := MarkdownPost{FsPath: path}
	post.LoadContent()
	return &post
}

func (post *MarkdownPost) LoadContent() {
	f, err := os.OpenFile(post.FsPath, os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}
	sc := bufio.NewScanner(f)
	post.ParseArgs(sc)
	var content []byte
	for sc.Scan() {
		read := sc.Bytes()
		content = append(content, read...)
		content = append(content, '\n')
	}
	err = sc.Err()
	if err != nil {
		panic(err)
	}
	post.Content = content
}

func (post *MarkdownPost) ParseArgs(sc *bufio.Scanner) {
	for sc.Scan() {
		line := sc.Text()
		parts := strings.Split(line, ":")
		if len(parts) < 2 {
			break
		}
		parts[1] = strings.Trim(parts[1], " \n")
		switch parts[0] {
		case "Title":
			post.Title = parts[1]
		case "Date":
			post.Date = parts[1]
		case "ReadingTime":
			post.ReadingTime = parts[1]
		case "Sats":
			sats, err := strconv.Atoi(parts[1])
			if err != nil {
				panic(err)
			}
			post.Sats = int(sats)
		}
	}
	err := sc.Err()
	if err != nil {
		panic(err)
	}
}

func (post *MarkdownPost) InsertHeader(htmlContent *[]byte) {
	header := []byte("" +
		"<code class=\"bg-transparent\"><strong><pre class=\"bg-transparent text-center\">\n" +
		" _     _             \n" +
		"| |__ | | ___   __ _ \n" +
		"| '_ \\| |/ _ \\ / _` |\n" +
		"| |_) | | (_) | (_| |\n" +
		"|_.__/|_|\\___/ \\__, |\n" +
		"                |___/ </pre></strong></code>\n" +
		"<div><div class=\"font-mono mb-1 text-center\">\n" +
		"<strong>{{- .Title }}</strong><br />\n" +
		"<small>{{- .Date }} | {{ .ReadingTime }} | {{ .Sats }} sats</small>\n" +
		"</div>\n")
	*htmlContent = append(header, *htmlContent...)
}

func (post *MarkdownPost) StyleHtml(htmlContent *[]byte) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(*htmlContent))
	if err != nil {
		panic(err)
	}
	doc.Find("img").Each(func(index int, element *goquery.Selection) {
		element.AddClass("flex m-auto")
	})
	doc.Find("pre, code").Each(func(index int, element *goquery.Selection) {
		element.AddClass("code")
	})
	doc.Find("code[class*=\"language-\"]").Each(func(index int, element *goquery.Selection) {
		SyntaxHighlighting(element)
	})
	htmlS, err := doc.Html()
	if err != nil {
		panic(err)
	}
	*htmlContent = []byte(htmlS)
}

func GetPosts(srcDir string) *[]MarkdownPost {
	paths, err := filepath.Glob(srcDir + "*.md")
	if err != nil {
		panic(err)
	}
	var posts []MarkdownPost
	for _, path := range paths {
		post := NewMarkdownPost(path)
		posts = append(posts, *post)
	}
	return &posts
}

func (post *MarkdownPost) Render(destDir string) {
	destPath := strings.ReplaceAll(destDir+filepath.Base(post.FsPath), ".md", ".html")
	f, err := os.Create(destPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	opts := html.RendererOptions{Flags: MarkdownToHtmlFlags}
	renderer := html.NewRenderer(opts)
	html := markdown.ToHTML(post.Content, nil, renderer)
	post.InsertHeader(&html)
	post.StyleHtml(&html)
	t, err := template.New("post").Parse(string(html))
	if err != nil {
		panic(err)
	}
	t.Execute(f, *post)
}

func RenderBlogIndex(srcDir string, destDir string, posts *[]MarkdownPost) {
	srcPath := srcDir + "index.html"
	destPath := destDir + "index.html"
	f, err := os.Create(destPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	t := template.New(filepath.Base(srcPath))
	t = t.Funcs(template.FuncMap{
		"ToHref": func(fsPath string) string {
			return "/" + strings.ReplaceAll(fsPath, ".md", ".html")
		},
	})
	t, err = t.ParseFiles(srcPath)
	if err != nil {
		panic(err)
	}
	err = t.Execute(f, map[string][]MarkdownPost{"Posts": *posts})
	if err != nil {
		panic(err)
	}
}
