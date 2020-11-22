package latex_test

import (
	"fmt"


	// TODO: Update link on v2 release.
	bf "github.com/russross/blackfriday/v2"
	bflatex "github.com/bjmayor/blackfriday-latex"
)

func Example() {
	const input = `
# Section

Some _Markdown_ text.

## Subsection

Foobar.
`

	extensions := bf.CommonExtensions | bf.Titleblock
	renderer := &bflatex.Renderer{
		Author:    "John Doe",
		Languages: "english,french",
		Flags:     bflatex.TOC,
	}
	md := bf.New(bf.WithRenderer(renderer), bf.WithExtensions(extensions))

	ast := md.Parse([]byte(input))
	fmt.Printf("%s\n", renderer.Render(ast))
	// Output:
	// \section{Section}
	// Some \emph{Markdown} text.
	//
	// \subsection{Subsection}
	// Foobar.
}
