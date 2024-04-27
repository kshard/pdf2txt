//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/pdf2txt
//

package pdf2txt_test

import (
	"bytes"
	"compress/gzip"
	"os"
	"testing"

	"github.com/fogfish/it/v2"
	"github.com/kshard/pdf2txt"
)

func TestText(t *testing.T) {
	for file, title := range map[string]string{
		"./testdata/spec_ieee_article_2_col.tsv.gz":  "NEW FACULTY 101: AN ORIENTATION TO THE PROFESSION",
		"./testdata/spec_arxiv_article_1_col.tsv.gz": "CateCom: a practical data-centric approach to categorization of computational models.",
	} {
		t.Run(file, func(t *testing.T) {
			fd, err := os.Open(file)
			it.Then(t).Should(it.Nil(err))

			r, err := gzip.NewReader(fd)
			it.Then(t).Should(it.Nil(err))

			b := &bytes.Buffer{}

			p, err := pdf2txt.New(pdf2txt.WithDirectStream())
			it.Then(t).Should(it.Nil(err))

			err = p.ToText(r, b)
			txt := b.String()

			it.Then(t).Should(
				it.Nil(err),
				it.Equal(txt[:len(title)], title),
			)
		})
	}
}

func TestMarkdown(t *testing.T) {
	for file, title := range map[string]string{
		"./testdata/spec_ieee_article_2_col.tsv.gz":  "# NEW FACULTY 101: AN ORIENTATION TO THE PROFESSION",
		"./testdata/spec_arxiv_article_1_col.tsv.gz": "# CateCom: a practical data-centric approach to categorization of computational models.",
	} {
		t.Run(file, func(t *testing.T) {
			fd, err := os.Open(file)
			it.Then(t).Should(it.Nil(err))

			r, err := gzip.NewReader(fd)
			it.Then(t).Should(it.Nil(err))

			b := &bytes.Buffer{}

			p, err := pdf2txt.New(pdf2txt.WithDirectStream())
			it.Then(t).Should(it.Nil(err))

			err = p.ToMarkdown(r, b)
			txt := b.String()

			it.Then(t).Should(
				it.Nil(err),
				it.Equal(txt[:len(title)], title),
			)
		})
	}
}
