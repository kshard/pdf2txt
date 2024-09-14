//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/pdf2txt
//

// Package pdf2txt is a front-end to `pdftotext`. It uses pdf's abstract syntax
// tree produced by the utility to reconstruct the text.
package pdf2txt

import (
	"strings"
	"unicode"
)

// Configuration option for parser and renderer
type Option func(*Parser)

// Disables usage `pdftotext`, instead input stream is assumed to be tsv
func WithDirectStream() Option {
	return func(p *Parser) {
		p.useDirectStream = true
	}
}

type Stream func(*Node) error

//------------------------------------------------------------------------------

type Node struct {
	Level    int     `col:"0"`
	PageNum  int     `col:"1"`
	ParNum   int     `col:"2"`
	BlockNum int     `col:"3"`
	LineNum  int     `col:"4"`
	WordNum  int     `col:"5"`
	Left     float64 `col:"6"`
	Top      float64 `col:"7"`
	Width    float64 `col:"8"`
	Height   float64 `col:"9"`
	Conf     int     `col:"10"`
	Text     string  `col:"11"`
}

func (n *Node) IsCtrl() bool {
	return len(n.Text) > 6 && strings.HasPrefix(n.Text, "###") && strings.HasSuffix(n.Text, "###")
}

func (n *Node) IsPage() bool { return n.Text == "###PAGE###" }
func (n *Node) IsFlow() bool { return n.Text == "###FLOW###" }
func (n *Node) IsLine() bool { return n.Text == "###LINE###" }
func (n *Node) IsText() bool { return !n.IsCtrl() }

func (n *Node) UnicodeText() string {
	w := n.Text
	w = strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, w)

	w = strings.ReplaceAll(w, "ﬁ", "fi")
	w = strings.ReplaceAll(w, "ﬀ", "ff")
	w = strings.ReplaceAll(w, "ﬃ", "ffi")
	w = strings.ReplaceAll(w, "ﬄ", "ffl")

	return w
}

type Parser struct {
	Version         string
	useDirectStream bool
}
