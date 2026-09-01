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

func renderPost(tmpl *template.Template, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	parts := strings.SplitN(string(content), "\n\n", 2)
	post := Post{
		Title: strings.TrimSpace(parts[0]),
		Body:  strings.TrimSpace(parts[1]),
	}

	if err := tmpl.Execute(os.Stdout, post); err != nil {
		return err
	}

	outputName := strings.TrimSuffix(path, ".txt") + ".html"
	out, err := os.Create(outputName)
	if err != nil {
		return err
	}
	defer out.Close()

	return tmpl.Execute(out, post)
}

func main() {
	file := flag.String("file", "first-post.txt", "name of the .txt file to render")
	dir := flag.String("dir", "", "directory to scan for .txt files")
	flag.Parse()

	tmpl, err := template.ParseFiles("template.tmpl")
	if err != nil {
		panic(err)
	}

	if *dir != "" {
		matches, err := filepath.Glob(filepath.Join(*dir, "*.txt"))
		if err != nil {
			panic(err)
		}
		for _, m := range matches {
			fmt.Println(m)
		}
		for _, m := range matches {
			if err := renderPost(tmpl, m); err != nil {
				panic(err)
			}
		}
		return
	}

	if err := renderPost(tmpl, *file); err != nil {
		panic(err)
	}
}