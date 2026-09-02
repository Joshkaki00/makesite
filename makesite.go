package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/yuin/goldmark"
)

type Post struct {
	Title string
	Body  string
}

const (
	ansiBold      = "\033[1m"
	ansiBoldGreen = "\033[1;32m"
	ansiReset     = "\033[0m"
)

// renderMarkdown converts Markdown source (including # through ######
// headings) into HTML using goldmark.
func renderMarkdown(source []byte) (string, error) {
	var buf bytes.Buffer
	if err := goldmark.Convert(source, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func renderPost(tmpl *template.Template, path string) (int64, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return 0, err
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
			return 0, err
		}
		body = html
	}

	post := Post{
		Title: title,
		Body:  body,
	}

	// Rendering into an in-memory buffer first lets us write the file in a
	// single syscall and read the size back from the buffer, avoiding an
	// extra fstat() round trip.
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, post); err != nil {
		return 0, err
	}

	outputName := strings.TrimSuffix(strings.TrimSuffix(path, ".txt"), ".md") + ".html"
	if err := os.WriteFile(outputName, buf.Bytes(), 0o644); err != nil {
		return 0, err
	}

	return int64(buf.Len()), nil
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
		start := time.Now()

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

		// Rendering is I/O-bound (reading source files, writing HTML), so
		// overlapping many of these operations via a bounded worker pool
		// cuts wall-clock time well below what a purely CPU-bound workload
		// would gain from parallelism.
		maxWorkers := runtime.NumCPU() * 4
		sem := make(chan struct{}, maxWorkers)
		var wg sync.WaitGroup
		var mu sync.Mutex
		var totalBytes int64
		var firstErr error

		for _, m := range matches {
			wg.Add(1)
			sem <- struct{}{}
			go func(path string) {
				defer wg.Done()
				defer func() { <-sem }()

				size, err := renderPost(tmpl, path)

				mu.Lock()
				defer mu.Unlock()
				if err != nil {
					if firstErr == nil {
						firstErr = err
					}
					return
				}
				totalBytes += size
			}(m)
		}
		wg.Wait()

		if firstErr != nil {
			panic(firstErr)
		}

		kb := float64(totalBytes) / 1000.0
		elapsed := time.Since(start).Seconds()
		fmt.Printf("%sSuccess!%s Generated %s%d%s pages (%.1fkB total) in %.2f seconds\n",
			ansiBoldGreen, ansiReset, ansiBold, len(matches), ansiReset, kb, elapsed)
		return
	}

	if _, err := renderPost(tmpl, *file); err != nil {
		panic(err)
	}
}
