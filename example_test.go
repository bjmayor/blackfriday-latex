package latex_test

import (
	"fmt"

	bflatex "bitbucket.org/ambrevar/blackfriday-latex"

	// TODO: Update link on v2 release.
	bf "gopkg.in/russross/blackfriday.v2"
)

func Example() {
	const input = `
# Section

Some _Markdown_ text.

## Subsection

Foobar.
`

	extensions := bf.CommonExtensions | bf.TOC | bf.Titleblock

	ast := bf.Parse([]byte(input), bf.Options{Extensions: extensions})

	renderer := bflatex.Renderer{
		Author:     "John Doe",
		Languages:  "english,french",
		Extensions: extensions,
	}

	fmt.Printf("%s\n", renderer.Render(ast))
	// Output:
	// \section{Section}
	// Some \emph{Markdown} text.
	//
	// \subsection{Subsection}
	// Foobar.
}
