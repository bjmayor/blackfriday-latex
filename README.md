# Blackfriday-LaTeX

Blackfriday-LaTeX is a LaTeX renderer for the [Blackfriday][] Markdown processor (v2).

[blackfriday]: http://github.com/russross/blackfriday

Warning: Both Blackfriday v2 and this renderer are a work-in-progress.

## Supported features

Among others:

- Optional preamble
- Title page
- Table of contents
- Footnotes
- Tables
- Fenced source code

## Math support

Since Markdown and CommonMark do not make any provision for parsing math, the
renderer uses the following rules to render math:

- Math blocks are introduced with code blocks having the `math` language specifier.

		``` math
		x+y=z
		```

- Inline math is introduce with inline code prefixed by `$$ ` (space matters).

		`$$ x+y=z`

## Documentation

See [godoc.org](https://godoc.org/github.com/ambrevar/blackfriday-latex).
