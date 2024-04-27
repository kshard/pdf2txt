//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/pdf2txt
//

package pdf2txt

import (
	"io"
	"math"
)

// PlainText converts stream of nodes into plain/text
type PlainText struct {
	w    io.Writer
	flow *Node
}

func NewPlainText(w io.Writer) *PlainText {
	f := &PlainText{w: w}
	return f
}

func (f *PlainText) Visit(node *Node) (err error) {
	if node.IsFlow() {
		f.flow = node
	}

	// Rule 1: H1 of the document is 1st paragraph
	if node.PageNum == 1 && node.ParNum == 0 && node.BlockNum == 0 {
		if node.IsText() {
			if _, err = f.w.Write([]byte(node.UnicodeText())); err != nil {
				return
			}
			_, err = f.w.Write([]byte(" "))
		}
		return
	}
	if node.PageNum == 1 && node.ParNum == 0 && node.BlockNum != 0 && node.IsFlow() {
		_, err = f.w.Write([]byte("\n"))
		return
	}

	// Rule 2: Paragraph has indent
	if node.IsLine() && math.Abs(node.Left-f.flow.Left) > 0.2 {
		_, err = f.w.Write([]byte("\n\n"))
		return
	}
	if node.IsLine() {
		_, err = f.w.Write([]byte("\n"))
		return
	}

	if node.IsText() {
		if _, err = f.w.Write([]byte(node.UnicodeText())); err != nil {
			return
		}
		_, err = f.w.Write([]byte(" "))
		return
	}

	return nil
}

//------------------------------------------------------------------------------

// Markdown converts stream of nodes into markdown
type Markdown struct {
	w    io.Writer
	flow *Node
}

func NewMarkdown(w io.Writer) *Markdown {
	f := &Markdown{w: w}
	return f
}

func (f *Markdown) Visit(node *Node) (err error) {
	if node.IsFlow() {
		f.flow = node
	}

	// Rule 1: H1 of the document is 1st paragraph
	if node.PageNum == 1 && node.ParNum == 0 && node.BlockNum == 0 {
		if node.IsFlow() {
			_, err = f.w.Write([]byte("# "))
			return
		}
		if node.IsText() {
			if _, err = f.w.Write([]byte(node.UnicodeText())); err != nil {
				return
			}
			_, err = f.w.Write([]byte(" "))
		}
		return
	}
	if node.PageNum == 1 && node.ParNum == 0 && node.BlockNum != 0 && node.IsFlow() {
		_, err = f.w.Write([]byte("\n"))
		return
	}

	// Rule 2: Paragraph has indent
	if node.IsLine() && math.Abs(node.Left-f.flow.Left) > 0.2 {
		_, err = f.w.Write([]byte("\n\n"))
		return
	}
	if node.IsLine() {
		_, err = f.w.Write([]byte("\n"))
		return
	}

	if node.IsText() {
		if _, err = f.w.Write([]byte(node.UnicodeText())); err != nil {
			return
		}
		_, err = f.w.Write([]byte(" "))
		return
	}

	return nil
}
