package main

import (
	"fmt"
	"html/template"
	"os"
)

func main() {
	tmpl, err := template.ParseFiles(
		"html/index.html",
		"html/template/head.html",
		"html/template/nav.html",
		"html/template/footer.html",
	)
	if err != nil {
		fmt.Printf("Error parsing templates: %v\n", err)
		os.Exit(1)
	}

	outputFile, err := os.Create("public/index.html")
	if err != nil {
		fmt.Printf("Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	err = tmpl.Execute(outputFile, nil)
	if err != nil {
		fmt.Printf("Error executing template: %v\n", err)
		os.Exit(1)
	}
}
