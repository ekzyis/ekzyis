package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

var (
	t     = template.Must(template.ParseGlob("html/template/*.html"))
	paths = []string{
		"index.html", "404.html",
		"blog/index.html",
		"blog/20230808-demystifying-wireguard-and-iptables.html",
	}
	dev bool
)

func init() {
	flag.BoolVar(&dev, "dev", false, "Specify if files should be built for development mode")
	flag.Parse()
}

func buildFiles() {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	buildDate := time.Now().In(time.UTC).Format("2006-01-02 15:04:05.000000000 -0700")
	for _, path := range paths {
		title := "ekzyis"
		if strings.Contains(path, "/") {
			title = strings.Split(path, "/")[0] + " | ekzyis"
		}
		env := "production"
		if dev {
			env = "development"
		}

		content, err := os.ReadFile(fmt.Sprintf("html/pages/%s", path))
		if err != nil {
			panic(err)
		}
		file, err := os.Create(fmt.Sprintf("public/%s", path))
		if err != nil {
			panic(err)
		}
		defer file.Close()

		data := map[string]string{
			"Title":     title,
			"Body":      string(content),
			"BuildDate": buildDate,
			"Env":       env,
		}
		mw := m.Writer("text/html", file)
		defer mw.Close()
		err = t.ExecuteTemplate(mw, "layout.html", data)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	buildFiles()
}
