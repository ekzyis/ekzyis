package main

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

var (
	m         *minify.M
	BuildDate string
	Env       string
	t         = template.Must(template.ParseGlob("html/template/*.html"))
)

func init() {
	m = minify.New()
	m.AddFunc("text/html", html.Minify)
	BuildDate = time.Now().In(time.UTC).Format("2006-01-02 15:04:05.000000000 -0700")
	Env = "production"
	if dev {
		Env = "development"
	}
}

type HtmlSource struct {
	FsPath  string
	Title   string
	Href    string
	Content []byte
}

func NewHtmlSource(path string) *HtmlSource {
	source := HtmlSource{FsPath: path}
	source.ComputeTitle()
	source.ComputeHref()
	source.LoadContent()
	return &source
}

func (source *HtmlSource) ComputeTitle() {
	source.Title = "ekzyis"
	if source.FsPath == "blog/index.html" {
		source.Title = "blog | ekzyis"
	}
}

func (source *HtmlSource) ComputeHref() {
	source.Href = strings.ReplaceAll(strings.ToLower(source.FsPath), ".html", "")
	if source.Href == "index" {
		source.Href = "/"
	}
}

func (source *HtmlSource) LoadContent() {
	content, err := os.ReadFile(source.FsPath)
	if err != nil {
		panic(err)
	}
	source.Content = content
}

func (source *HtmlSource) Render(destDir string) {
	destPath := destDir + filepath.Base(source.FsPath)
	args := map[string]any{
		"BuildDate": BuildDate,
		"Env":       Env,
		"Title":     source.Title,
		"Href":      source.Href,
		"Content":   string(source.Content),
	}
	ExecuteTemplate(destPath, args)
}

func ExecuteTemplate(destPath string, args map[string]any) {
	file, err := os.Create(destPath)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	mw := m.Writer("text/html", file)
	defer mw.Close()
	err = t.ExecuteTemplate(mw, "layout.html", args)
	if err != nil {
		panic(err)
	}
}

func GetHtmlSources(srcDir string) *[]HtmlSource {
	var sources []HtmlSource
	paths, err := filepath.Glob(srcDir + "*.html")
	if err != nil {
		panic(err)
	}
	for _, path := range paths {
		sources = append(sources, *NewHtmlSource(path))
	}
	return &sources
}
