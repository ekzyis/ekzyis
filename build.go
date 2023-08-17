package main

import (
	"flag"
)

var (
	dev            bool
	BlogSrcDir     = "blog/"
	BlogDstDir     = "html/pages/blog/"
	HtmlSrcDirs    = []string{"html/pages/", "html/pages/blog/"}
	HtmlTargetDirs = []string{"public/", "public/blog/"}
)

func init() {
	flag.BoolVar(&dev, "dev", false, "Specify if files should be built for development mode")
	flag.Parse()
}

func main() {
	posts := GetPosts(BlogSrcDir)
	for _, post := range *posts {
		post.Render(BlogDstDir)
	}
	RenderBlogIndex(BlogSrcDir, BlogDstDir, posts)
	// Go does not support ** globstar ...
	// https://github.com/golang/go/issues/11862
	for i, srcDir := range HtmlSrcDirs {
		for _, source := range *GetHtmlSources(srcDir) {
			source.Render(HtmlTargetDirs[i])
		}
	}
}
