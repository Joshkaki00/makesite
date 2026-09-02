package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/yuin/goldmark"
)

type Post struct {
	Title string
	Body  string
}

// renderMarkdown converts Markdown source (including # through ######
// headings) into HTML using goldmark.
func renderMarkdown(source []byte) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert(source, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func renderPost(tmpl *template.Template, path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	isMarkdown := strings.HasSuffix(path, ".md")

	parts := strings.SplitN(string(content), "\n\n", 2)
	title := strings.TrimSpace(parts[0])
	if isMarkdown {
		// A Markdown post's title comes from a leading "# " heading.
		title = strings.TrimSpace(strings.TrimPrefix(title, "#"))
	}

	body := ""
	if len(parts) > 1 {
		body = strings.TrimSpace(parts[1])
	}
	if isMarkdown {
		html, err := renderMarkdown([]byte(body))
		if err != nil {
			return err
		}
		body = html
	}

	post := Post{
		Title: title,
		Body:  body,
	}

	if err := tmpl.Execute(os.Stdout, post); err != nil {
		return err
	}

	outputName := strings.TrimSuffix(strings.TrimSuffix(path, ".txt"), ".md") + ".html"
	out, err := os.Create(outputName)
	if err != nil {
		return err
	}
	defer out.Close()

	return tmpl.Execute(out, post)
}

func main() {
	file := flag.String("file", "first-post.txt", "name of the .txt or .md file to render")
	dir := flag.String("dir", "", "directory to scan for .txt and .md files")
	flag.Parse()

	tmpl, err := template.ParseFiles("template.tmpl")
	if err != nil {
		panic(err)
	}

	if *dir != "" {
		var matches []string
		err := filepath.WalkDir(*dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, ".txt") {
				matches = append(matches, path)
			}
			return nil
		})
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