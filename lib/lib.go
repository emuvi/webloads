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
	GetContents(resp.Body, file)
	PutReferences(file, input)
}

func GetContents(fromReader io.ReadCloser, file *os.File) {
	defer fromReader.Close()
	writer := Writer{file, false}
	tokens := html.NewTokenizer(fromReader)
	var blocks_stack = []Block{}
	var ignore_stack = []string{}
	for {
		typed := tokens.Next()
		if typed == html.ErrorToken {
			return
		}
		token := tokens.Token()
		switch {
		case typed == html.StartTagToken:
			if token.Data == "h1" {
				blocks_stack = append(blocks_stack, Block{
					tag:    token.Data,
					opener: "\n\n# ",
					closer: "\n\n",
				})
			} else if token.Data == "h2" {
				blocks_stack = append(blocks_stack, Block{
					tag:    token.Data,
					opener: "\n\n## ",
					closer: "\n\n",
				})
			} else if token.Data == "h3" {
				blocks_stack = append(blocks_stack, Block{
					tag:    token.Data,
					opener: "\n\n### ",
					closer: "\n\n",
				})
			} else if token.Data == "h4" {
				blocks_stack = append(blocks_stack, Block{
					tag:    token.Data,
					opener: "\n\n#### ",
					closer: "\n\n",
				})
			} else if token.Data == "h5" {
				blocks_stack = append(blocks_stack, Block{
					tag:    token.Data,
					opener: "\n\n##### ",
					closer: "\n\n",
				})
			} else if token.Data == "h6" {
				blocks_stack = append(blocks_stack, Block{
					tag:    token.Data,
					opener: "\n\n###### ",
					closer: "\n\n",
				})
			} else if token.Data == "pre" {
				blocks_stack = append(blocks_stack, Block{
					tag:    token.Data,
					opener: "\n```\n",
					closer: "\n```\n\n",
				})
			} else if token.Data == "p" {
				blocks_stack = append(blocks_stack, Block{
					tag:    token.Data,
					opener: "\n",
					closer: "\n",
				})
			} else if token.Data == "div" {
				blocks_stack = append(blocks_stack, Block{
					tag:    token.Data,
					opener: "\n",
					closer: "\n",
				})
			} else if token.Data == "span" {
				blocks_stack = append(blocks_stack, Block{
					tag:    token.Data,
					opener: "",
					closer: "",
				})
			} else if token.Data == "a" {
				blocks_stack = append(blocks_stack, Block{
					tag:    token.Data,
					opener: "",
					closer: "",
				})
			} else if token.Data == "style" {
				ignore_stack = append(ignore_stack, token.Data)
			} else if token.Data == "script" {
				ignore_stack = append(ignore_stack, token.Data)
			}
		case typed == html.TextToken:
			if len(blocks_stack) > 0 && len(ignore_stack) == 0 {
				stacked := &blocks_stack[len(blocks_stack)-1]
				text := strings.TrimSpace(token.Data)
				if len(text) > 0 {
					if stacked.texted != "" {
						stacked.texted += " " + text
					} else {
						stacked.texted = text
					}
				}
			}
		case typed == html.EndTagToken:
			if len(ignore_stack) > 0 {
				if token.Data == ignore_stack[len(ignore_stack)-1] {
					ignore_stack = ignore_stack[:len(ignore_stack)-1]
				}
			} else if len(blocks_stack) > 0 {
				stacked := &blocks_stack[len(blocks_stack)-1]
				if stacked.tag == token.Data {
					if stacked.texted != "" {
						for i := len(blocks_stack) - 2; i >= 0; i-- {
							parent := &blocks_stack[i]
							if parent.opened {
								break
							} else {
								if parent.opener != "" {
									writer.Write(parent.opener)
								}
								parent.opened = true
								if parent.texted != "" {
									writer.Write(parent.texted)
									parent.texted = ""
								}
							}
						}
						if stacked.opener != "" {
							writer.Write(stacked.opener)
						}
						stacked.opened = true
						writer.Write(stacked.texted)
					}
					if stacked.opened && stacked.closer != "" {
						writer.Write(stacked.closer)
					}
					blocks_stack = blocks_stack[:len(blocks_stack)-1]
				}
			}
		}
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

type Block struct {
	tag    string
	opener string
	closer string
	texted string
	opened bool
}

type Writer struct {
	file        *os.File
	on_new_line bool
}

func (w *Writer) Write(part string) {
	if !w.on_new_line {
		w.file.WriteString(" ")
	}
	w.file.WriteString(part)
	w.on_new_line = strings.HasSuffix(part, "\n")
}
