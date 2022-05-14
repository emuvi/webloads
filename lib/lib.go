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

func GetContents(fromReader io.ReadCloser, toFile *os.File) {
	defer fromReader.Close()
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
					opener: "# ",
					texted: "",
					closer: "\n\n",
				})
			} else if token.Data == "h2" {
				blocks_stack = append(blocks_stack, Block{
					tag:    token.Data,
					opener: "## ",
					texted: "",
					closer: "\n\n",
				})
			} else if token.Data == "h3" {
				blocks_stack = append(blocks_stack, Block{
					tag:    token.Data,
					opener: "### ",
					texted: "",
					closer: "\n\n",
				})
			} else if token.Data == "h4" {
				blocks_stack = append(blocks_stack, Block{
					tag:    token.Data,
					opener: "#### ",
					texted: "",
					closer: "\n\n",
				})
			} else if token.Data == "h5" {
				blocks_stack = append(blocks_stack, Block{
					tag:    token.Data,
					opener: "##### ",
					texted: "",
					closer: "\n\n",
				})
			} else if token.Data == "h6" {
				blocks_stack = append(blocks_stack, Block{
					tag:    token.Data,
					opener: "###### ",
					texted: "",
					closer: "\n\n",
				})
			} else if token.Data == "pre" {
				blocks_stack = append(blocks_stack, Block{
					tag:    token.Data,
					opener: "```\n",
					texted: "",
					closer: "\n```\n",
				})
			} else if token.Data == "p" {
				blocks_stack = append(blocks_stack, Block{
					tag:    token.Data,
					opener: "",
					texted: "",
					closer: "\n",
				})
			} else if token.Data == "div" {
				blocks_stack = append(blocks_stack, Block{
					tag:    token.Data,
					opener: "",
					texted: "",
					closer: "\n",
				})
			} else if token.Data == "span" {
				blocks_stack = append(blocks_stack, Block{
					tag:    token.Data,
					opener: "",
					texted: "",
					closer: "",
				})
			} else if token.Data == "a" {
				blocks_stack = append(blocks_stack, Block{
					tag:    token.Data,
					opener: "",
					texted: "",
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
						stacked.texted += " "
					}
					stacked.texted += text
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
						if stacked.opener != "" {
							toFile.WriteString(stacked.opener)
						}
						toFile.WriteString(stacked.texted)
						if stacked.closer != "" {
							toFile.WriteString(stacked.closer)
						}
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
	texted string
	closer string
}
