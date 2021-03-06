package lib

import (
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func Parse(input string, output string) {
	resp, err := http.Get(input)
	if err != nil {
		panic(err)
	}
	file, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	lines := GetContents(resp.Body)
	WriteLines(file, lines)
	PutReferences(file, input)
}

func GetContents(fromBody io.ReadCloser) []string {
	defer fromBody.Close()
	writer := Writer{}
	tokens := html.NewTokenizer(fromBody)
	var inContent = 0
	var linkHref = ""
	var inIgnored = []string{}
	for {
		kind := tokens.Next()
		if kind == html.ErrorToken {
			return writer.lines
		}
		token := tokens.Token()
		switch {
		case kind == html.StartTagToken:
			if token.Data == "h1" {
				writer.Write("\n\n# ")
				inContent++
			} else if token.Data == "h2" {
				writer.Write("\n\n## ")
				inContent++
			} else if token.Data == "h3" {
				writer.Write("\n### ")
				inContent++
			} else if token.Data == "h4" {
				writer.Write("\n#### ")
				inContent++
			} else if token.Data == "h5" {
				writer.Write("\n##### ")
				inContent++
			} else if token.Data == "h6" {
				writer.Write("\n###### ")
				inContent++
			} else if token.Data == "pre" {
				writer.Write("\n```\n")
				inContent++
			} else if token.Data == "div" {
				writer.Write("\n")
				inContent++
			} else if token.Data == "p" {
				writer.Write("\n")
				inContent++
			} else if token.Data == "code" {
				writer.Write("`")
				inContent++
			} else if token.Data == "span" {
				inContent++
			} else if token.Data == "a" {
				writer.Write("[")
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						linkHref = attr.Val
					}
				}
				inContent++
			} else if token.Data == "style" {
				inIgnored = append(inIgnored, token.Data)
			} else if token.Data == "script" {
				inIgnored = append(inIgnored, token.Data)
			}
		case kind == html.TextToken:
			if inContent > 0 {
				text := strings.TrimSpace(token.Data)
				if len(text) > 0 {
					writer.Write(text)
				}
			}
		case kind == html.EndTagToken:
			if len(inIgnored) > 0 {
				if token.Data == inIgnored[len(inIgnored)-1] {
					inIgnored = inIgnored[:len(inIgnored)-1]
				}
			} else {
				if token.Data == "h1" {
					writer.Write("\n\n")
					inContent--
				} else if token.Data == "h2" {
					writer.Write("\n\n")
					inContent--
				} else if token.Data == "h3" {
					writer.Write("\n")
					inContent--
				} else if token.Data == "h4" {
					writer.Write("\n")
					inContent--
				} else if token.Data == "h5" {
					writer.Write("\n")
					inContent--
				} else if token.Data == "h6" {
					writer.Write("\n")
					inContent--
				} else if token.Data == "pre" {
					writer.Write("\n```\n")
					inContent--
				} else if token.Data == "div" {
					writer.Write("\n")
					inContent--
				} else if token.Data == "p" {
					writer.Write("\n")
					inContent--
				} else if token.Data == "code" {
					writer.Write("`")
					inContent--
				} else if token.Data == "span" {
					inContent--
				} else if token.Data == "a" {
					writer.Write("]")
					if linkHref != "" {
						writer.Write("(")
						writer.Write(linkHref)
						writer.Write(")")
						linkHref = ""
					}
					inContent--
				}
			}
		}
	}
}

func WriteLines(file *os.File, lines []string) {
	for _, line := range lines {
		file.WriteString(line)
	}
}

func PutReferences(file *os.File, input string) {
	file.WriteString("\n")
	file.WriteString("\n###### WebLoads Reference")
	file.WriteString("\n")
	file.WriteString("\n- From: <")
	file.WriteString(input)
	file.WriteString(">")
	file.WriteString("\n- When: ")
	file.WriteString(time.Now().UTC().Format(time.RFC3339))
	file.WriteString("\n")
}

type Writer struct {
	lines []string
	ended string
}

func (w *Writer) Write(part string) {
	if w.shouldSpace(part) {
		part = " " + part
	}
	w.lines = append(w.lines, part)
	w.ended = part[len(part)-1:]
}

func (w *Writer) shouldSpace(part string) bool {
	if w.ended == "\n" {
		return false
	}
	if w.ended == "[" || w.ended == "]" || part == "]" {
		return false
	}
	if w.ended == "(" || part == ")" {
		return false
	}
	return true
}
