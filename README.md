# Blackfriday-LaTeX

Blackfriday-LaTeX is a LaTeX renderer for the
[Blackfriday](http://github.com/russross/blackfriday) library (v2 or above).

Warning: Both Blackfriday v2 and this renderer are a work-in-progress.

## Supported features

Among others:

- Fenced source code
- Title page
- Table of contents
- Tables

## Renderer parameters

- Author: The document author displayed by the `\maketitle` command. This will
only display if the `Titleblock` extension is on and a title is present.

- Languages: The languages to be used by the `babel` package. Languages must be
comma-spearated.

## Example

``` go
package main

import (
	"fmt"

	bflatex "bitbucket.org/ambrevar/blackfriday-latex"
	bf "gopkg.in/ambrevar/blackfriday.v2"
)

var input = `
% Sample input

# Section

Some _Markdown_ text.

## Subsection

Foobar.
`

func main() {
	extensions := bf.CommonExtensions | bf.TOC | bf.Titleblock

	ast := bf.Parse([]byte(input), bf.Options{Extensions: extensions})
	renderer := bflatex.Renderer{Author: "John Doe", Languages: "english,french", Extensions: extensions}
	fmt.Printf("%s\n", renderer.Render(ast))
}
```
