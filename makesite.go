package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Post struct {
	Title string
	Body  string
}

func main() {
	file := flag.String("file", "first-post.txt", "name of the .txt file to render")
	dir := flag.String("dir", "", "directory to scan for .txt files")
	flag.Parse()

	if *dir != "" {
		matches, err := filepath.Glob(filepath.Join(*dir, "*.txt"))
		if err != nil {
			panic(err)
		}
		for _, m := range matches {
			fmt.Println(m)
		}
		return
	}

	content, err := os.ReadFile(*file)
	if err != nil {
		panic(err)
	}

	parts := strings.SplitN(string(content), "\n\n", 2)
	post := Post{
		Title: strings.TrimSpace(parts[0]),
		Body:  strings.TrimSpace(parts[1]),
	}

	tmpl, err := template.ParseFiles("template.tmpl")
	if err != nil {
		panic(err)
	}

	if err := tmpl.Execute(os.Stdout, post); err != nil {
		panic(err)
	}

	outputName := strings.TrimSuffix(*file, ".txt") + ".html"
	out, err := os.Create(outputName)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	if err := tmpl.Execute(out, post); err != nil {
		panic(err)
	}
}