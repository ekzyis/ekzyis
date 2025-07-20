package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type Markdown = string

type Post struct {
	Title  string
	Time   time.Time
	Hidden bool
	URL    string
	body   Markdown
}

type TemplateData struct {
	Posts []Post
}

var (
	contentDir = "content"
)

func walkMarkdownContent() ([]Markdown, error) {
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

func parseFrontmatter(content Markdown) (map[string]interface{}, error) {
	parts := strings.SplitN(content, "---", 3)
	if len(parts) < 3 {
		return nil, fmt.Errorf("no frontmatter found")
	}

	var frontmatter map[string]interface{}
	err := yaml.Unmarshal([]byte(parts[1]), &frontmatter)
	if err != nil {
		return nil, err
	}

	return frontmatter, nil
}

func formatDate(dateStr string) string {
	// Try different date formats
	formats := []string{
		"2006-01-02T15:04:05.000Z",
		"2006-01-02T15:04:05Z",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t.Format("2 Jan 2006")
		}
	}

	// If all parsing fails, return the original string
	return dateStr
}

func parseMarkdown(markdown []Markdown) ([]Post, error) {
	var posts []Post

	for _, path := range markdown {
		// Read the markdown file
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %v", path, err)
		}

		// Parse frontmatter
		frontmatter, err := parseFrontmatter(string(content))
		if err != nil {
			return nil, fmt.Errorf("failed to parse frontmatter for %s: %v", path, err)
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

		// url
		dir := filepath.Dir(path)
		url := strings.TrimPrefix(dir, contentDir+"/")

		// body
		body := Markdown(content)

		post := Post{
			Title:  title,
			Time:   time,
			Hidden: hidden,
			URL:    url,
			body:   body,
		}

		posts = append(posts, post)
	}

	// sort newest first
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Time.After(posts[j].Time)
	})

	return posts, nil
}

func executeTemplates(posts []Post) error {
	// define template functions
	tmpl := template.New("").Funcs(
		template.FuncMap{
			"formatTime": func(t time.Time, format string) string {
				return t.Format(format)
			},
		},
	)

	// parse templates
	tmpl, err := tmpl.ParseFiles(
		"html/index.html",
		"html/template/head.html",
		"html/template/nav.html",
		"html/template/footer.html",
		"html/template/post/list.html",
	)
	if err != nil {
		return fmt.Errorf("error parsing templates: %v", err)
	}

	// create output file
	outputPath := filepath.Join("public", "index.html")
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating output file: %v", err)
	}
	defer outputFile.Close()

	// execute template with post data
	err = tmpl.ExecuteTemplate(outputFile, "index.html", TemplateData{Posts: posts})
	if err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	fmt.Printf("generated %s with %d posts\n", outputPath, len(posts))

	return nil
}

func main() {
	// Step 1: Walk directory and collect all index.md files
	mdFiles, err := walkMarkdownContent()
	if err != nil {
		fmt.Printf("error walking markdown content: %v\n", err)
		os.Exit(1)
	}

	// Step 2: Parse frontmatter for each file
	posts, err := parseMarkdown(mdFiles)
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
