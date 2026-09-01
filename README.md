# makesite

A small static site generator written in Go. It reads a plain text post
file, renders it into an HTML template, and writes the result to disk.

## Requirements

- Go 1.21.1 or later
- No external Go dependencies (standard library only)

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

A post file is plain text. The first line is used as the title. The
rest of the file (after a blank line) is used as the body.

```
Post Title

Body text goes here. It can span multiple paragraphs.
```

## Project structure

- `makesite.go` - entry point; reads a post file, renders the template
- `template.tmpl` - HTML template with `Title` and `Body` placeholders
- `first-post.txt`, `latest-post.txt` - example post files
- `go.mod` - module definition

## License

MIT. See `LICENSE`.
