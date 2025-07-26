package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	embed "github.com/13rac1/goldmark-embed"
	"github.com/alecthomas/chroma/v2/styles"
	sn "github.com/ekzyis/snappy"
	figure "github.com/mangoumbrella/goldmark-figure"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
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
	Commit  string
}

type ErrorTemplateData struct {
	BaseURL string
	Commit  string
}

type PostTemplateData struct {
	Post    Post
	BaseURL string
	Commit  string
}

type MinifyWriter struct {
	w   io.Writer
	buf bytes.Buffer
	m   *minify.M
}

func NewHtmlMinifyWriter(w io.Writer) *MinifyWriter {
	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	return &MinifyWriter{
		w: w,
		m: m,
	}
}

func (mw *MinifyWriter) Write(p []byte) (n int, err error) {
	return mw.buf.Write(p)
}

func (mw *MinifyWriter) Close() error {
	return mw.m.Minify("text/html", mw.w, &mw.buf)
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
	commit  string
)

func main() {
	if os.Getenv("ENV") == "production" {
		baseUrl = "https://ekzy.is"
	}

	output, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		fmt.Printf("error getting commit: %v\n", err)
		os.Exit(1)
	}
	commit = strings.TrimRight(string(output), "\n")

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
	blogFiles, err := filepath.Glob(filepath.Join(contentDir, "*", "index.md"))
	if err != nil {
		return nil, err
	}

	journalFiles, err := filepath.Glob(filepath.Join(contentDir, "journal", "*", "index.md"))
	if err != nil {
		return nil, err
	}

	files := append(blogFiles, journalFiles...)

	return files, nil
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
	url := strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(dir, contentDir+"/")), "_", "-")

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
		"html/404.html",
		"html/template/head.html",
		"html/template/nav.html",
		"html/template/footer.html",
		"html/template/post/list.html",
		"html/template/post/single.html",
	)
	if err != nil {
		return fmt.Errorf("error parsing templates: %v", err)
	}

	var blogPosts []Post
	var journalPosts []Post
	for _, post := range posts {
		if !strings.Contains(post.Path, "/journal/") {
			blogPosts = append(blogPosts, post)
		} else {
			journalPosts = append(journalPosts, post)
		}
	}

	err = executeIndexTemplate(tmpl, "public/index.html", blogPosts)
	if err != nil {
		return fmt.Errorf("error executing index template: %v", err)
	}

	journalDir := filepath.Join("public", "journal")
	err = os.MkdirAll(journalDir, 0755)
	if err != nil {
		return fmt.Errorf("error creating directory for journal: %v", err)
	}
	err = executeIndexTemplate(tmpl, "public/journal/index.html", journalPosts)
	if err != nil {
		return fmt.Errorf("error executing journal index template: %v", err)
	}

	err = executeErrorTemplate(tmpl)
	if err != nil {
		return fmt.Errorf("error executing error template: %v", err)
	}

	err = executePostTemplates(tmpl, posts)
	if err != nil {
		return fmt.Errorf("error executing post templates: %v", err)
	}

	return nil
}

func executeIndexTemplate(tmpl *template.Template, outputPath string, posts []Post) error {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()

	mw := NewHtmlMinifyWriter(outputFile)
	defer mw.Close()
	err = tmpl.ExecuteTemplate(mw, "index.html", IndexTemplateData{
		Posts:   posts,
		BaseURL: baseUrl,
		Commit:  commit,
	})
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	fmt.Printf("written %s\n", outputPath)
	return nil
}

func executeErrorTemplate(tmpl *template.Template) error {
	outputPath := filepath.Join("public", "404.html")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()

	mw := NewHtmlMinifyWriter(outputFile)
	defer mw.Close()
	err = tmpl.ExecuteTemplate(mw, "404.html", ErrorTemplateData{
		BaseURL: baseUrl,
		Commit:  commit,
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
		mw := NewHtmlMinifyWriter(postFile)
		defer mw.Close()
		err = tmpl.ExecuteTemplate(mw, "post.html", PostTemplateData{
			Post:    post,
			BaseURL: baseUrl,
			Commit:  commit,
		})
		if err != nil {
			return fmt.Errorf("error executing single post template for %s: %v", post.Title, err)
		}

		// copy images as webp
		for _, image := range post.Images {
			dstImage := filepath.Join(postDir, filepath.Base(image))
			dstImage = strings.TrimSuffix(strings.TrimSuffix(dstImage, ".png"), ".jpg") + ".webp"
			err = toWebp(image, dstImage)
			if err != nil {
				return fmt.Errorf("error copying image %s: %v", image, err)
			}
		}

		fmt.Printf("written %s\n", postHTMLPath)
	}

	return nil
}

func toWebp(src, dst string) error {
	err := exec.Command("ffmpeg", "-i", src, "-c:v", "libwebp", dst).Run()
	if err != nil {
		return fmt.Errorf("error converting image to webp: %v", err)
	}

	return nil
}
