package main

import (
	"bytes"
	"os"
	"os/exec"
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
	t         = template.Must(template.ParseGlob("html/template/*.html"))
)

func init() {
	m = minify.New()
	m.AddFunc("text/html", html.Minify)
	BuildDate = time.Now().In(time.UTC).Format("2006-01-02 15:04:05.000000000 -0700")
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
	// FIXME
	//   this is just a temporary workaround.
	//   actual solution would be to use title from blog post markdown metadata
	switch source.FsPath {
	case "html/pages/blog/index.html":
		source.Title = "blog | ekzyis"
	case "html/pages/blog/20230809-demystifying-wireguard-and-iptables.html":
		source.Title = "Demystifying WireGuard and iptables | blog | ekzyis"
	case "html/pages/blog/20230821-wireguard-packet-forwarding.html":
		source.Title = "WireGuard Packet Forwarding | blog | ekzyis"
	case "html/pages/blog/20230925-wireguard-port-forwarding.html":
		source.Title = "WireGuard Port Forwarding | blog | ekzyis"
	default:
		source.Title = "ekzyis"
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

func GetVersion() string {
	cmd := exec.Command("git", "rev-parse", "--short", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	return out.String()
}

func (source *HtmlSource) Render() {
	args := map[string]any{
		"BuildDate": BuildDate,
		"Env":       Env,
		"Title":     source.Title,
		"Href":      source.Href,
		"Content":   string(source.Content),
		"Version":   GetVersion(),
	}
	ExecuteTemplate(args)
}

func ExecuteTemplate(args map[string]any) {
	mw := m.Writer("text/html", os.Stdout)
	defer mw.Close()
	err := t.ExecuteTemplate(mw, "layout.html", args)
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
