package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
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
	goldmarkHtml "github.com/yuin/goldmark/renderer/html"
	"gopkg.in/yaml.v3"
)

var (
	quotes = []Quote{
		{
			// https://stacker.news/items/536315
			Quote:  "I'm a bitcoiner as much as I'm not.",
			Author: "ekzyis",
		},
		{
			// https://stacker.news/items/547681
			Quote:  "It's easier to give people the same grace you give yourself when you remember this is their first time in this life, too.",
			Author: "plebpoet",
		},
		{
			// https://morehouse.github.io/lightning/fake-channel-dos/
			Quote:  "Because in the end it doesn't matter how feature-rich and easy-to-use the Lightning Network is if it can't keep user funds safe.",
			Author: "Matt Morehouse",
		},
		{
			// https://www.goodreads.com/book/show/2761.The_Denial_of_Death
			Quote:  "The problem is that we all want to be more than a shitting and fucking creature that dies.",
			Author: "Review of The Denial of Death by Ernest Becker",
		},
		{
			// https://www.youtube.com/watch?v=l5NK8zdIK-g (17:53)
			Quote:  "Get enough sleep. Drink more water.",
			Author: "Gibi ASMR",
		},
		{
			// https://www.youtube.com/watch?v=l5NK8zdIK-g (34:57)
			Quote:  "We are insanely social creatures, as much as we pretend that we aren't.",
			Author: "Gibi ASMR, again",
		},
		{
			Quote:  "Wow, you made it to the end of these *very* inspiring quotes!",
			Author: "Sarcasm",
		},
	}
)

type Quote struct {
	Quote  string
	Author string
	URL    string
}

type Banner struct {
	Src    string
	Width  int
	Height int
}

type Post struct {
	Path      string
	Title     string
	Time      time.Time
	Publish   bool
	Frontpage bool
	URL       string
	Banner    *Banner
	Tags      []string
	SnId      int
	Comments  int
	Sats      int
	Markdown  string
	HTML      template.HTML
	Images    []string
}

type IndexTemplateData struct {
	Quotes []Quote
	Posts  []Post
	Commit string
}

type ErrorTemplateData struct {
	Commit string
}

type PostTemplateData struct {
	Post   Post
	Commit string
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
		// allow raw HTML in markdown, required for grid of images
		goldmark.WithRendererOptions(goldmarkHtml.WithUnsafe()),
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
	commit         string
	updateMetadata bool
)

func main() {
	fmt.Printf(
		`       ___.         .__.__       .___
___  __\_ |__  __ __|__|  |    __| _/
\  \/  /| __ \|  |  \  |  |   / __ |
 >    < | \_\ \  |  /  |  |__/ /_/ |
/__/\_ \|___  /____/|__|____/\____ |
      \/    \/                    \/
   xbuild: because it sounds cool` + "\n\n")
	flag.BoolVar(&updateMetadata, "M", false, "Update metadata from Stacker News API")
	flag.Parse()

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
		fmt.Printf("error parsing markdown: %v\n", err)
		os.Exit(1)
	}

	err = executeTemplates(posts)
	if err != nil {
		fmt.Printf("error executing templates: %v\n", err)
		os.Exit(1)
	}
}

func walkMarkdownContent() ([]string, error) {
	// if files are passed as arguments, use them instead of walking the content directory
	args := flag.Args()
	if len(args) > 0 {
		return args, nil
	}

	files, err := filepath.Glob(filepath.Join(contentDir, "*", "index.md"))
	if err != nil {
		return nil, err
	}

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

	// Should this post be published?
	publish, ok := frontmatter["publish"].(bool)
	if !ok {
		// publish by default
		publish = true
	}

	// can this post be on the frontpage?
	frontpage, ok := frontmatter["frontpage"].(bool)
	if !ok {
		// default to false if omitted
		frontpage = false
	}

	// banner
	banner := &Banner{}
	banner.Src, ok = frontmatter["banner"].(string)
	if ok {
		path := filepath.Join(filepath.Dir(path), banner.Src)
		banner.Width, banner.Height, err = imageDimensions(path)
		if err != nil {
			return nil, fmt.Errorf(">> failed to get image dimensions: %s: %v", path, err)
		}
		// TODO: fix awkward behavior wrt image paths in frontmatter vs content
		banner.Src = toWebpPath(banner.Src)
	} else {
		banner = nil
	}

	// stacker news
	snId, _ := frontmatter["sn_id"].(int)
	var comments, sats int
	if updateMetadata && snId > 0 {
		fmt.Printf("> updating metadata for %s\n", path)
		fmt.Printf(">> fetching item %d from Stacker News\n", snId)
		item, err := sn.NewClient().Item(snId)
		if err != nil {
			return nil, fmt.Errorf(">> failed to fetch item %d: %v", snId, err)
		}
		comments = item.NComments
		sats = item.Sats
		fmt.Printf(">> new metadata: %d comments, %d sats\n", comments, sats)
	}

	// url
	url, ok := frontmatter["url"].(string)
	if !ok {
		dir := filepath.Dir(path)
		// slash at the end is required so images can use relative links
		// without the slash at the end, src="./diff.webp" would be relative to root
		// TODO: load images even when there's no / at the end
		url = "/" + strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(dir, contentDir+"/")), "_", "-") + "/"
	}

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
			if !info.IsDir() && regexp.MustCompile(`\.(png|jpg|webp)$`).MatchString(path) {
				images = append(images, path)
			}
			return nil
		})
	if err != nil {
		return nil, fmt.Errorf("failed to find images for %s: %v", path, err)
	}

	return &Post{
		Path:      path,
		Title:     title,
		Time:      time,
		Publish:   publish,
		Frontpage: frontpage,
		Banner:    banner,
		SnId:      snId,
		Comments:  comments,
		Sats:      sats,
		URL:       url,
		Markdown:  markdown,
		HTML:      html,
		Images:    images,
	}, nil
}

func executeTemplates(posts []Post) error {
	tmpl := template.New("").Funcs(
		template.FuncMap{
			"formatTime": func(t time.Time, format string) string {
				return t.Format(format)
			},
			"hasPrefix": strings.HasPrefix,
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

	var frontpagePosts, otherPosts []Post
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Time.After(posts[j].Time)
	})
	for _, p := range posts {
		if !p.Publish {
			continue
		}
		if p.Frontpage {
			frontpagePosts = append(frontpagePosts, p)
		} else {
			otherPosts = append(otherPosts, p)
		}
	}

	err = executeIndexTemplate(tmpl, "public/index.html", frontpagePosts)
	if err != nil {
		return fmt.Errorf("error executing index template: %v", err)
	}

	err = executeIndexTemplate(tmpl, "public/other.html", otherPosts)
	if err != nil {
		return fmt.Errorf("error executing other index template: %v", err)
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
		Quotes: quotes,
		Posts:  posts,
		Commit: commit,
	})
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	fmt.Printf("> %s\n", outputPath)
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
		Commit: commit,
	})
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	fmt.Printf("> %s\n", outputPath)
	return nil
}

func executePostTemplates(tmpl *template.Template, posts []Post) error {
	for _, post := range posts {
		if strings.HasPrefix(post.URL, "https:") {
			// this post will link to an external site, no template execution needed
			continue
		}
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
			Post:   post,
			Commit: commit,
		})
		if err != nil {
			return fmt.Errorf("error executing single post template for %s: %v", post.Title, err)
		}

		// copy images as webp
		for _, image := range post.Images {
			dstImage := filepath.Join(postDir, filepath.Base(image))
			if strings.HasSuffix(image, ".webp") {
				// already webp, just copy it
				err = copyFile(image, dstImage)
				if err != nil {
					return err
				}
				fmt.Printf("> %s\n", dstImage)
				continue
			}
			dstImage = toWebpPath(dstImage)
			err = toWebp(image, dstImage)
			if err != nil {
				return fmt.Errorf("error copying image %s: %v", image, err)
			}
			fmt.Printf("> %s\n", dstImage)
		}

		fmt.Printf("> %s\n", postHTMLPath)
	}

	return nil
}

func copyFile(src, dst string) error {
	// copying a file using go is ridiculously more code
	err := exec.Command("cp", src, dst).Run()
	if err != nil {
		return fmt.Errorf("error copying file %s to %s: %v", src, dst, err)
	}
	return nil
}

func toWebpPath(path string) string {
	return strings.TrimSuffix(strings.TrimSuffix(path, ".png"), ".jpg") + ".webp"
}

func toWebp(src, dst string) error {
	err := exec.Command("ffmpeg", "-i", src, "-c:v", "libwebp", dst).Run()
	if err != nil {
		return fmt.Errorf("error converting image to webp: %v", err)
	}

	return nil
}

func imageDimensions(path string) (int, int, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "v:0",
		"-show_entries", "stream=width,height",
		"-of", "csv=p=0:s=x",
		path,
	)
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, err
	}
	s := strings.TrimSpace(string(out))
	parts := strings.Split(s, "x")
	if len(parts) != 2 {
		return 0, 0, err
	}
	w, err1 := strconv.Atoi(parts[0])
	h, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil || w <= 0 || h <= 0 {
		return 0, 0, err
	}
	return w, h, nil
}
