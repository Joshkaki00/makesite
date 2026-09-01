package main

import (
	"os"
	"strings"
	"text/template"
)

type Post struct {
	Title string
	Body  string
}

func main() {
	content, err := os.ReadFile("first-post.txt")
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
}