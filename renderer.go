package main

import (
	"log"
	"path"

	"github.com/namsral/flag"
)

var (
	Source string
	Env    string
)

func init() {
	flag.StringVar(&Source, "src", "", "Source file")
	flag.StringVar(&Env, "env", "development", "Specify for which environment files should be built")
	flag.Parse()
}

func RenderExtension(path string, ext string) {
	switch ext {
	case ".md":
		NewMarkdownPost(path).Render()
	case ".html":
		NewHtmlSource(path).Render()
	default:
		log.Fatalf("unknown extension: %s", ext)
	}
}

func main() {
	if Source == "" {
		log.Fatal("no source given")
	}

	if Source == "blog/index.html" {
		RenderBlogIndex("blog/")
		return
	}

	if Source != "" {
		ext := path.Ext(Source)
		if ext == "" {
			log.Fatal("file has no extension")
		}
		RenderExtension(Source, ext)
		return
	}
}
