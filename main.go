package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	embed "github.com/13rac1/goldmark-embed"
	"github.com/alecthomas/chroma/v2/styles"
	sn "github.com/ekzyis/snappy"
	figure "github.com/mangoumbrella/goldmark-figure"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"gopkg.in/yaml.v3"
)

type Post struct {
	Path     string
	Title    string
	Time     time.Time
	Hidden   bool
	URL      string
	Banner   string
	SnId     int
	Comments int
	Sats     int
	Markdown string
	HTML     template.HTML
	Images   []string
}

type IndexTemplateData struct {
	Posts   []Post
	BaseURL string
}

type PostTemplateData struct {
	Post    Post
	BaseURL string
}

var (
	contentDir = "content"
	mdParser   = goldmark.New(
		goldmark.WithExtensions(
			extension.Footnote,
			figure.Figure,
			embed.New(),
			highlighting.NewHighlighting(
				// https://swapoff.org/chroma/playground/
				highlighting.WithCustomStyle(styles.Get("catppuccin-mocha")),
			),
		),
	)
	baseUrl = "http://localhost:8080/"
)

func main() {
	if os.Getenv("ENV") == "production" {
		baseUrl = "https://ekzy.is"
	}

	mdFiles, err := walkMarkdownContent()
	if err != nil {
		fmt.Printf("error walking markdown content: %v\n", err)
		os.Exit(1)
	}

	posts, err := parsePosts(mdFiles)
	if err != nil {
		fmt.Printf("error parsing markdown: %v", err)
		os.Exit(1)
	}

	err = executeTemplates(posts)
	if err != nil {
		fmt.Printf("error executing templates: %v", err)
		os.Exit(1)
	}
}

func walkMarkdownContent() ([]string, error) {
	pattern := filepath.Join(contentDir, "*", "index.md")
	markdownFiles, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	var filteredFiles []string
	for _, path := range markdownFiles {
		// journal posts have their own page
		if !strings.Contains(path, "/journal/") {
			filteredFiles = append(filteredFiles, path)
		}
	}

	return filteredFiles, nil
}

func parsePosts(paths []string) ([]Post, error) {
	var posts []Post

	for _, path := range paths {
		post, err := parsePost(path)
		if err != nil {
			return nil, fmt.Errorf("failed to parse post %s: %v", path, err)
		}

		posts = append(posts, *post)
	}

	// sort newest first
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Time.After(posts[j].Time)
	})

	return posts, nil
}

func parsePost(path string) (*Post, error) {

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", path, err)
	}

	parts := strings.SplitN(string(content), "---", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("no frontmatter found")
	}

	var frontmatter map[string]interface{}
	err = yaml.Unmarshal([]byte(parts[1]), &frontmatter)
	if err != nil {
		return nil, err
	}

	// title
	title, _ := frontmatter["title"].(string)

	// time
	time, ok := frontmatter["date"].(time.Time)
	if !ok {
		return nil, fmt.Errorf("failed to parse date for %s: %v", path, frontmatter["date"])
	}

	// hidden
	hidden, _ := frontmatter["hidden"].(bool)

	// banner
	banner, _ := frontmatter["banner"].(string)

	// stacker news
	snId, _ := frontmatter["sn_id"].(int)
	var comments, sats int
	if snId > 0 {
		item, err := sn.NewClient().Item(snId)
		if err != nil {
			return nil, fmt.Errorf("failed to get sn item %d: %v", snId, err)
		}
		comments = item.NComments
		sats = item.Sats
	}

	// url
	dir := filepath.Dir(path)
	url := strings.TrimPrefix(dir, contentDir+"/")

	// markdown
	markdown := string(parts[2])

	// html
	var buf bytes.Buffer
	if err := mdParser.Convert([]byte(markdown), &buf); err != nil {
		return nil, fmt.Errorf("failed to convert markdown for %s: %v", path, err)
	}
	html := template.HTML(buf.String())

	// images
	var images []string
	err = filepath.Walk(
		filepath.Dir(path),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && regexp.MustCompile(`\.(png|jpg)$`).MatchString(path) {
				images = append(images, path)
			}
			return nil
		})
	if err != nil {
		return nil, fmt.Errorf("failed to find images for %s: %v", path, err)
	}

	return &Post{
		Path:     path,
		Title:    title,
		Time:     time,
		Hidden:   hidden,
		Banner:   banner,
		SnId:     snId,
		Comments: comments,
		Sats:     sats,
		URL:      url,
		Markdown: markdown,
		HTML:     html,
		Images:   images,
	}, nil
}

func executeTemplates(posts []Post) error {
	tmpl := template.New("").Funcs(
		template.FuncMap{
			"formatTime": func(t time.Time, format string) string {
				return t.Format(format)
			},
		},
	)

	tmpl, err := tmpl.ParseFiles(
		"html/index.html",
		"html/post.html",
		"html/template/head.html",
		"html/template/nav.html",
		"html/template/footer.html",
		"html/template/post/list.html",
		"html/template/post/single.html",
	)
	if err != nil {
		return fmt.Errorf("error parsing templates: %v", err)
	}

	err = executeIndexTemplate(tmpl, posts)
	if err != nil {
		return fmt.Errorf("error executing index template: %v", err)
	}

	err = executePostTemplates(tmpl, posts)
	if err != nil {
		return fmt.Errorf("error executing post templates: %v", err)
	}

	return nil
}

func executeIndexTemplate(tmpl *template.Template, posts []Post) error {
	outputPath := filepath.Join("public", "index.html")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()

	err = tmpl.ExecuteTemplate(outputFile, "index.html", IndexTemplateData{
		Posts:   posts,
		BaseURL: baseUrl,
	})
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	fmt.Printf("written %s\n", outputPath)
	return nil
}

func executePostTemplates(tmpl *template.Template, posts []Post) error {
	for _, post := range posts {
		postDir := filepath.Join("public", post.URL)
		err := os.MkdirAll(postDir, 0755)
		if err != nil {
			return fmt.Errorf("error creating directory for post %s: %v", post.Title, err)
		}

		// Create HTML file
		postHTMLPath := filepath.Join(postDir, "index.html")
		postFile, err := os.Create(postHTMLPath)
		if err != nil {
			return fmt.Errorf("error creating HTML file for post %s: %v", post.Title, err)
		}
		defer postFile.Close()

		// Execute the single post template
		err = tmpl.ExecuteTemplate(postFile, "post.html", PostTemplateData{
			Post:    post,
			BaseURL: baseUrl,
		})
		if err != nil {
			return fmt.Errorf("error executing single post template for %s: %v", post.Title, err)
		}

		// copy images
		for _, image := range post.Images {
			dstImage := filepath.Join(postDir, filepath.Base(image))
			err = copyFile(image, dstImage)
			if err != nil {
				return fmt.Errorf("error copying image %s: %v", image, err)
			}
		}

		fmt.Printf("written %s\n", postHTMLPath)
	}

	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error opening source file %s: %v", src, err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("error creating destination file %s: %v", dst, err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("error copying file %s: %v", src, err)
	}

	return nil
}
