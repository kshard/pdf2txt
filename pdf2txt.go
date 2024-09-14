//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/kshard/pdf2txt
//

package pdf2txt

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/Masterminds/semver/v3"
)

const expectedVsn = ">= 22.05.0"

// Create new parser instance
func New(opts ...Option) (*Parser, error) {
	p := &Parser{}

	for _, opt := range opts {
		opt(p)
	}

	if !p.useDirectStream {
		vsn, err := p.checkVersion()
		if err != nil {
			return nil, err
		}
		p.Version = vsn
	}

	return p, nil
}

// check existence of pdftotext and return its version
func (p *Parser) version() (string, error) {
	// Expected output:
	//	pdftotext version 24.04.0
	//	Copyright 2005-2024 The Poppler Developers - http://poppler.freedesktop.org
	//	Copyright 1996-2011, 2022 Glyph & Cog, LLC
	buf := &bytes.Buffer{}
	cmd := exec.Command("pdftotext", "-v")
	cmd.Stderr = buf

	if err := cmd.Run(); err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(buf)
	if ok := scanner.Scan(); ok {
		vsn := scanner.Text()
		seq := strings.Fields(vsn)
		if len(seq) < 3 {
			return "", fmt.Errorf("pdftotext: unexpected output %s", vsn)
		}
		return seq[2], nil
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("pdftotext: unexpected output %s", buf.String())
}

// check existence of pdftotext and its version
func (p *Parser) checkVersion() (string, error) {
	vsn, err := p.version()
	if err != nil {
		return vsn, err
	}

	check, err := semver.NewConstraint(expectedVsn)
	if err != nil {
		return vsn, err
	}

	version, err := semver.NewVersion(vsn)
	if err != nil {
		return vsn, err
	}

	if !check.Check(version) {
		return vsn, fmt.Errorf("pdftotext version %s exists, version %s is expected", vsn, expectedVsn)
	}

	return vsn, nil
}

func (p *Parser) ToText(f io.Reader, w io.Writer) error {
	format := NewPlainText(w)
	return p.Stream(f, format.Visit)
}

func (p *Parser) ToMarkdown(f io.Reader, w io.Writer) error {
	format := NewMarkdown(w)
	return p.Stream(f, format.Visit)
}

// Stream PDF document
func (p *Parser) Stream(f io.Reader, stream Stream) (err error) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	r, w := io.Pipe()

	go func() {
		defer wg.Done()

		if p.useDirectStream {
			_, err = io.Copy(w, f)
		} else {
			cmd := exec.Command("pdftotext", "-tsv", "-", "-")
			cmd.Stdin = f
			cmd.Stdout = w
			cmd.Stderr = os.Stderr

			err = cmd.Run()
		}

		w.Close()
	}()

	go func() {
		defer wg.Done()

		scanner := bufio.NewScanner(r)
		scanner.Scan() // skip header
		for scanner.Scan() {
			txt := scanner.Text()
			seq := strings.Fields(txt)

			node, err := parseNode(seq)
			if err != nil {
				return
			}

			err = stream(node)
			if err != nil {
				return
			}
		}

		err = scanner.Err()
	}()

	wg.Wait()

	return
}

func parseNode(seq []string) (node *Node, err error) {
	node = &Node{}
	node.Level, err = strconv.Atoi(seq[0])
	if err != nil {
		return
	}

	node.PageNum, err = strconv.Atoi(seq[1])
	if err != nil {
		return
	}

	node.ParNum, err = strconv.Atoi(seq[2])
	if err != nil {
		return
	}

	node.BlockNum, err = strconv.Atoi(seq[3])
	if err != nil {
		return
	}

	node.LineNum, err = strconv.Atoi(seq[4])
	if err != nil {
		return
	}

	node.WordNum, err = strconv.Atoi(seq[5])
	if err != nil {
		return
	}

	node.Left, err = strconv.ParseFloat(seq[6], 64)
	if err != nil {
		return
	}

	node.Top, err = strconv.ParseFloat(seq[7], 64)
	if err != nil {
		return
	}

	node.Width, err = strconv.ParseFloat(seq[8], 64)
	if err != nil {
		return
	}

	node.Height, err = strconv.ParseFloat(seq[9], 64)
	if err != nil {
		return
	}

	node.Conf, err = strconv.Atoi(seq[10])
	if err != nil {
		return
	}

	node.Text = seq[11]

	return
}
