// Copyright © 2013-2016 Pierre Neidhardt <ambrevar@gmail.com>
// Use of this file is governed by the license that can be found in LICENSE.

// TODO: LaTeX-style "Quotes"? See v1.
// TODO: Add tests other TODOs are closed.
// TODO: Add flag to skip header.
// TODO: Add flag to use titleblock as chapter title.

// Package latex is a LaTeX renderer for the Blackfriday Markdown Processor
package latex

import (
	"bytes"
	"io"
	"path/filepath"

	bf "gopkg.in/russross/blackfriday.v2"
)

// Renderer is a type that implements the Renderer interface for LaTeX
// output.
//
// Do not create this directly, instead use the NewLatexRenderer function.
type Renderer struct {
	w          bytes.Buffer
	Extensions bf.Extensions

	Author    string
	Languages string
}

func cellAlignment(align bf.CellAlignFlags) string {
	switch align {
	case bf.TableAlignmentLeft:
		return "l"
	case bf.TableAlignmentRight:
		return "r"
	case bf.TableAlignmentCenter:
		return "c"
	default:
		return "l"
	}
}

func needsBackslash(c byte) bool {
	for _, r := range []byte("_{}%$&\\~#") {
		if c == r {
			return true
		}
	}
	return false
}

func (r *Renderer) esc(text []byte) {
	for i := 0; i < len(text); i++ {
		// directly copy normal characters
		org := i

		for i < len(text) && !needsBackslash(text[i]) {
			i++
		}
		if i > org {
			r.w.Write(text[org:i])
		}

		// escape a character
		if i >= len(text) {
			break
		}
		r.w.WriteByte('\\')
		r.w.WriteByte(text[i])
	}
}

// LaTeX preamble
// TODO: Color source code and links?
func (r *Renderer) writeDocumentHeader(title, author string, hasFigures bool) {
	titleDef := ""
	titleCommand := ""
	babel := ""

	if title != "" {
		titleDef = `
\title{` + title + `}
\author{` + author + `}
`
		titleCommand = `
\maketitle
`
		if r.Extensions&bf.TOC != 0 {
			titleCommand += `\vfill
\thispagestyle{empty}

\tableofcontents
`
			if hasFigures {
				titleCommand += `\listoffigures
`
			}
			titleCommand += `\clearpage
`
		}
	}

	if r.Languages != "" {
		babel = `
\usepackage[` + r.Languages + `]{babel}
`
	}

	r.out(`\documentclass{article}

\usepackage[utf8]{inputenc}
\usepackage[T1]{fontenc}
\usepackage{lmodern}
\usepackage{marvosym}
\usepackage{textcomp}
\DeclareUnicodeCharacter{20AC}{\EUR{}}
\DeclareUnicodeCharacter{2260}{\neq}
\DeclareUnicodeCharacter{2264}{\leq}
\DeclareUnicodeCharacter{2265}{\geq}
\DeclareUnicodeCharacter{22C5}{\cdot}
\DeclareUnicodeCharacter{A0}{~}
\DeclareUnicodeCharacter{B1}{\pm}
\DeclareUnicodeCharacter{D7}{\times}

\usepackage{amsmath}
\usepackage[export]{adjustbox} % loads also graphicx
\usepackage{listings}
\usepackage[margin=1in]{geometry}
\usepackage{verbatim}
\usepackage[normalem]{ulem}
\usepackage{hyperref}

\lstset{
	numbers=left,
	breaklines=true,
	xleftmargin=2\baselineskip,
	showstringspaces=false,
	basicstyle=\ttfamily,
	keywordstyle=\bfseries\color{green!40!black},
	commentstyle=\itshape\color{purple!40!black},
	stringstyle=\color{orange},
	numberstyle=\ttfamily,
	literate=
	{á}{{\'a}}1 {é}{{\'e}}1 {í}{{\'i}}1 {ó}{{\'o}}1 {ú}{{\'u}}1
	{Á}{{\'A}}1 {É}{{\'E}}1 {Í}{{\'I}}1 {Ó}{{\'O}}1 {Ú}{{\'U}}1
	` + "{à}{{\\`a}}1 {è}{{\\`e}}1 {ì}{{\\`i}}1 {ò}{{\\`o}}1 {ù}{{\\`u}}1" + `
	` + "{À}{{\\`A}}1 {È}{{\\'E}}1 {Ì}{{\\`I}}1 {Ò}{{\\`O}}1 {Ù}{{\\`U}}1" + `
	{ä}{{\"a}}1 {ë}{{\"e}}1 {ï}{{\"i}}1 {ö}{{\"o}}1 {ü}{{\"u}}1
	{Ä}{{\"A}}1 {Ë}{{\"E}}1 {Ï}{{\"I}}1 {Ö}{{\"O}}1 {Ü}{{\"U}}1
	{â}{{\^a}}1 {ê}{{\^e}}1 {î}{{\^i}}1 {ô}{{\^o}}1 {û}{{\^u}}1
	{Â}{{\^A}}1 {Ê}{{\^E}}1 {Î}{{\^I}}1 {Ô}{{\^O}}1 {Û}{{\^U}}1
	{œ}{{\oe}}1 {Œ}{{\OE}}1 {æ}{{\ae}}1 {Æ}{{\AE}}1 {ß}{{\ss}}1
	{ű}{{\H{u}}}1 {Ű}{{\H{U}}}1 {ő}{{\H{o}}}1 {Ő}{{\H{O}}}1
	{ç}{{\c c}}1 {Ç}{{\c C}}1 {ø}{{\o}}1 {å}{{\r a}}1 {Å}{{\r A}}1
	{€}{{\EUR}}1 {£}{{\pounds}}1
}
` + babel + `
\hypersetup{colorlinks,
	citecolor=black,
	filecolor=black,
	linkcolor=black,
	linktoc=page,
	urlcolor=black,
	pdfstartview=FitH,
	breaklinks=true,
	pdfauthor={Blackfriday Markdown Processor v` + bf.Version + `}
}

\newcommand{\HRule}{\rule{\linewidth}{0.5mm}}
\addtolength{\parskip}{0.5\baselineskip}
\parindent=0pt
` + titleDef + `
\begin{document}
` + titleCommand + `

`)
}

func (r *Renderer) writeDocumentFooter() {
	r.out(`\end{document}`)
	r.cr()
}

func languageAttr(info []byte) string {
	infoWords := bytes.Split(info, []byte("\t "))
	if len(infoWords) > 0 {
		return string(infoWords[0])
	}
	return ""
}

func (r *Renderer) out(s ...string) {
	for _, v := range s {
		r.w.WriteString(v)
	}
}

func (r *Renderer) cr() {
	r.out("\n")
}

func (r *Renderer) env(environment string, entering bool, args ...string) {
	if entering {
		r.out(`\begin{`, environment, `}`)
		if len(args) > 0 {
			r.out("[")
			for _, v := range args {
				r.out(v)
				r.out(",")
			}
			r.out("]")
		}

		r.cr()
	} else {
		r.out(`\end{`, environment, `}`)
		r.cr()
		r.cr()
	}
}

func (r *Renderer) envLiteral(environment string, literal []byte, args ...string) {
	r.env(environment, true, args...)
	r.w.Write(literal)
	r.env(environment, false)
}

func (r *Renderer) cmd(command string, entering bool) {
	if entering {
		r.out(`\`, command, `{`)
	} else {
		r.out(`}`)
	}
}

// Return the first ASCII character that is not in 'text'.
// The resulting delimiter cannot be '*' nor space.
func getDelimiter(text []byte) byte {
	delimiters := make([]bool, 128)
	for _, v := range text {
		if v < 128 {
			delimiters[v] = true
		}
	}
	// '!' is the character after space in the ASCII encoding.
	for k := byte('!'); k < byte('*'); k++ {
		if !delimiters[k] {
			return k
		}
	}
	// '+' is the character after '*' in the ASCII encoding.
	for k := byte('+'); k < 128; k++ {
		if !delimiters[k] {
			return k
		}
	}
	return 0
}

// RenderNode renders a single node.
// As a rule for consistency, each node is responsible for appending the needed
// line breaks. Line breaks are never prepended.
func (r *Renderer) RenderNode(w io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	switch node.Type {

	case bf.BlockQuote:
		r.env("quotation", entering)

	case bf.Code:
		// TODO: Reach a consensus for math syntax.
		if bytes.HasPrefix(node.Literal, []byte("$$ ")) {
			// Inline math
			r.w.WriteByte('$')
			r.w.Write(node.Literal[3:])
			r.w.WriteByte('$')
			break
		}
		delimiter := getDelimiter(node.Literal)
		r.out(`\lstinline`)
		if delimiter != 0 {
			r.w.WriteByte(delimiter)
			r.w.Write(node.Literal)
			r.w.WriteByte(delimiter)
		} else {
			r.out(`@`)
			r.out("<RENDERING ERROR: no delimiter found>")
			r.out(`@`)
		}

	case bf.CodeBlock:
		lang := languageAttr(node.Info)
		if lang == "math" {
			r.out("\\[\n")
			r.w.Write(node.Literal)
			r.out("\\]")
			r.cr()
			r.cr()
			break
		}
		r.envLiteral("lstlisting", node.Literal, `language=`+lang)

	case bf.Del:
		r.cmd("sout", entering)

	case bf.Document:
		break

	case bf.Emph:
		r.cmd("emph", entering)

	case bf.Hardbreak:
		// TODO: How many cr()? Or should it be a manual space? Like \vspace{\baselineskip}.
		r.cr()

	case bf.Header:
		if node.IsTitleblock {
			// Nothing to print but its children.
			break
		}
		if entering {
			switch node.Level {
			case 1:
				r.out(`\section{`)
			case 2:
				r.out(`\subsection{`)
			case 3:
				r.out(`\subsubsection{`)
			case 4:
				r.out(`\paragraph{`)
			case 5:
				r.out(`\subparagraph{`)
			case 6:
				r.out(`\textbf{`)
			}
		} else {
			r.out(`}`)
			switch node.Level {
			// Paragraph need no newline.
			case 1, 2, 3:
				r.cr()
			default:
				r.out(" ")
			}
		}

	case bf.HTMLBlock:
		break

	case bf.HTMLSpan:
		// HTML code makes no sense in LaTeX.
		break

	case bf.HorizontalRule:
		r.out(`\HRule{}`)
		r.cr()

	case bf.Image:
		if entering {
			dest := node.LinkData.Destination
			if bytes.HasPrefix(dest, []byte("http://")) || bytes.HasPrefix(dest, []byte("https://")) {
				r.out(`\url{`, string(dest), `}`)
				return bf.SkipChildren
			}
			if node.LinkData.Title != nil {
				r.out(`\begin{figure}[!ht]`)
				r.cr()
			}
			r.out(`\begin{center}`)
			r.cr()
			ext := filepath.Ext(string(dest))
			if len(ext) > 0 {
				// Trim extension so that LaTeX loads the most appropriate file.
				dest = dest[:len(dest)-len(ext)]
			}
			r.out(`\includegraphics[max width=\textwidth, max height=\textheight]{`, string(dest), `}`)
			r.cr()
			r.out(`\end{center}`)
			r.cr()
			if node.LinkData.Title != nil {
				r.out(`\caption{`, string(node.LinkData.Title), `}`)
				r.cr()
				r.out(`\end{figure}`)
				r.cr()
			}
		}
		return bf.SkipChildren

	case bf.Item:
		if entering {
			if node.ListFlags&bf.ListTypeTerm != 0 {
				r.out(`\item [`)
			} else if node.ListFlags&bf.ListTypeDefinition == 0 {
				r.out(`\item `)
			}
		} else {
			if node.ListFlags&bf.ListTypeTerm != 0 {
				r.out("] ")
			}
		}

	case bf.Link:
		// TODO: Links: What about safety? See HTML renderer.
		// TODO: Links: Add relative link support?
		// TODO: Links: Add e-mail support?
		// if kind == bf.LinkTypeEmail {
		// 	r.w.WriteString("mailto:")
		// }
		if node.NoteID != 0 {
			// TODO: Implement footnotes. Wait for upstream fix.
			if entering {
				r.out(`\footnotemark[`, string(node.LinkData.Destination), `]`)
			}
			break
		}
		r.cmd("href{"+string(node.LinkData.Destination)+"}", entering)

	case bf.List:
		listType := "itemize"
		if node.ListFlags&bf.ListTypeOrdered != 0 {
			listType = "enumerate"
		}
		if node.ListFlags&bf.ListTypeDefinition != 0 {
			listType = "description"
		}
		r.env(listType, entering)

	case bf.Paragraph:
		if !entering {
			// If paragraph is the term of a definition list, don't insert new lines.
			if node.Parent.Type != bf.Item || node.Parent.ListFlags&bf.ListTypeTerm == 0 {
				r.cr()
				// Don't insert an additional linebreak after last node of an item, a quote, etc.
				if node.Next != nil {
					r.cr()
				}
			}
		}

	case bf.Softbreak:
		// r.cr()
		// TODO: Make it configurable via out(renderer.softbreak)
		break

	case bf.Strong:
		r.cmd("textbf", entering)

	case bf.Table:
		if entering {
			r.out(`\begin{center}`)
			r.cr()
			r.out(`\begin{tabular}{`)
			node.Walk(func(c *bf.Node, entering bool) bf.WalkStatus {
				if c.Type == bf.TableCell && entering {
					for cell := c; cell != nil; cell = cell.Next {
						r.out(cellAlignment(cell.Align))
					}
					return bf.Terminate
				}
				return bf.GoToNext
			})
			r.out("}")
			r.cr()

		} else {
			r.out(`\end{tabular}`)
			r.cr()
			r.out(`\end{center}`)
			r.cr()
			r.cr()
		}

	case bf.TableBody:
		// Nothing to do here.
		break

	case bf.TableCell:
		if node.IsHeader {
			r.cmd("textbf", entering)
		}
		if !entering && node.Next != nil {
			r.out(" & ")
		}

	case bf.TableHead:
		if !entering {
			r.out(`\hline`)
			r.cr()
		}

	case bf.TableRow:
		if !entering {
			r.out(` \\`)
			r.cr()
		}

	case bf.Text:
		r.esc(node.Literal)
		break

	default:
		panic("Unknown node type " + node.Type.String())
	}
	return bf.GoToNext
}

// Get title: concatenate all Text children of Titleblock.
func title(ast *bf.Node) []byte {
	titleRenderer := Renderer{}

	ast.Walk(func(node *bf.Node, entering bool) bf.WalkStatus {
		if node.Type == bf.Header && node.HeaderData.IsTitleblock && entering {
			node.Walk(func(c *bf.Node, entering bool) bf.WalkStatus {
				return titleRenderer.RenderNode(&titleRenderer.w, c, entering)
			})
			return bf.Terminate
		}
		return bf.GoToNext
	})
	return titleRenderer.w.Bytes()
}

func hasFigures(ast *bf.Node) bool {
	result := false
	ast.Walk(func(node *bf.Node, entering bool) bf.WalkStatus {
		if node.Type == bf.Image && node.LinkData.Title != nil {
			result = true
			return bf.Terminate
		}
		return bf.GoToNext
	})
	return result
}

// Render prints out the header, 'ast' and its children recursively, and finally
// a footer.
func (r *Renderer) Render(ast *bf.Node) []byte {
	r.writeDocumentHeader(string(title(ast)), r.Author, hasFigures(ast))
	ast.Walk(func(node *bf.Node, entering bool) bf.WalkStatus {
		if node.Type == bf.Header && node.HeaderData.IsTitleblock {
			return bf.SkipChildren
		}
		return r.RenderNode(&r.w, node, entering)
	})
	r.writeDocumentFooter()
	return r.w.Bytes()
}
