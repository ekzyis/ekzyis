package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

type Post struct {
	Date        string
	Title       string
	ReadingTime string
	Sats        int
	Href        string
}

var (
	t     = template.Must(template.ParseGlob("html/template/*.html"))
	paths = map[string]any{
		"index.html":      nil,
		"404.html":        nil,
		"blog/index.html": nil,
		"blog/20230809-Demystifying-WireGuard-and-iptables.html": Post{
			Date:        "2023-08-09",
			Title:       "Demystifying WireGuard and iptables",
			ReadingTime: "15 minutes",
			Sats:        11623,
		},
	}
	dev bool
)

func init() {
	flag.BoolVar(&dev, "dev", false, "Specify if files should be built for development mode")
	flag.Parse()
}

func getPosts() []Post {
	var posts []Post
	for path, args := range paths {
		post, ok := args.(Post)
		if !ok {
			continue
		}
		post.Href = "/" + strings.ReplaceAll(strings.ToLower(path), ".html", "")
		posts = append(posts, post)
	}
	return posts
}

func buildFiles() {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	buildDate := time.Now().In(time.UTC).Format("2006-01-02 15:04:05.000000000 -0700")
	env := "production"
	if dev {
		env = "development"
	}
	for path, pathArgs := range paths {
		htmlTitle := "ekzyis"
		if path == "blog/index.html" {
			htmlTitle = "blog | ekzyis"
			pathArgs = map[string]any{"Posts": getPosts()}
		}
		if post, ok := pathArgs.(Post); ok {
			htmlTitle = post.Title
		}

		tmp, err := template.ParseFiles(fmt.Sprintf("html/pages/%s", path))
		if err != nil {
			panic(err)
		}
		buf := new(bytes.Buffer)
		tmp.Execute(buf, pathArgs)

		path = strings.ToLower(path)
		file, err := os.Create(fmt.Sprintf("public/%s", path))
		if err != nil {
			panic(err)
		}
		defer file.Close()

		rootArgs := map[string]any{
			"Title":     htmlTitle,
			"Body":      buf.String(),
			"BuildDate": buildDate,
			"Env":       env,
		}
		mw := m.Writer("text/html", file)
		defer mw.Close()
		err = t.ExecuteTemplate(mw, "layout.html", rootArgs)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	buildFiles()
}
