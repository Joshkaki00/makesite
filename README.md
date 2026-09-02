# makesite

A small static site generator written in Go. It reads a plain text post
file, renders it into an HTML template, and writes the result to disk.

## Requirements

- Go 1.22 or later
- [goldmark](https://github.com/yuin/goldmark) for Markdown rendering

## Third-party library

I will use the **goldmark** library. The documentation is located at
https://github.com/yuin/goldmark. My goal is to use it to parse
Markdown (`.md`) post files and transform them into HTML, so that
`#` through `######` headings become `<h1>` through `<h6>` elements
(along with other Markdown formatting like bold, italics, lists, and
links).

## Usage

Build the binary:

```
go build -o makesite .
```

Run it against a text file:

```
./makesite --file=first-post.txt
```

This prints the rendered HTML to stdout and writes an HTML file named
after the input file (for example, `first-post.txt` produces
`first-post.html`).

If `--file` is omitted, it defaults to `first-post.txt`.

## Post file format

Posts can be plain text (`.txt`) or Markdown (`.md`).

For `.txt` files, the first line is used as the title, and the rest of
the file (after a blank line) is used as the body verbatim:

```
Post Title

Body text goes here. It can span multiple paragraphs.
```

For `.md` files, the title comes from the leading `# ` heading, and the
rest of the file (after a blank line) is parsed as Markdown and
rendered to HTML using [goldmark](https://github.com/yuin/goldmark).
Headings `#` through `######` become `<h1>` through `<h6>` elements:

```
# Post Title

## A Subheading

Body text with **bold**, _italic_, lists, links, and more.
```

`--dir` scans for both `*.txt` and `*.md` files.

## Project structure

- `makesite.go` - entry point; reads a post file, renders the template
- `template.tmpl` - HTML template with `Title` and `Body` placeholders
- `first-post.txt`, `latest-post.txt` - example plain-text post files
- `fifth-post.md` - example Markdown post file
- `go.mod` / `go.sum` - module definition and dependency checksums

## License

MIT. See `LICENSE`.
