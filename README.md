# pdf2txt

The library converts pdf to text formats (e.g. plain text, markdown). The library is built over [poppler-utils](https://poppler.freedesktop.org) (version >=22.05.0). `pdftotext` is required and available in the path.

## Quick Start

```go
import (
  "github.com/kshard/pdf2txt"
)

// Create parser
parser, err := pdf2txt.New()
if err != nil {
  panic(err)
}

// Open input stream (io.Reader) to PDF 
fd, err := os.Open(/* path to file */)

// Open output stream (io.Writer) to destination
buf := &bytes.Buffer{}

// Convert
if err := parser.ToText(fd, buf); err != nil {
  panic(err)
}
```

