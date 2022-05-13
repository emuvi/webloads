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
	var open_by = []string{}
	var blocked_by = []string{}
	for {
		typed := tokens.Next()
		if typed == html.ErrorToken {
			return
		}
		token := tokens.Token()
		switch {
		case typed == html.StartTagToken:
			if token.Data == "h1" {
				toFile.WriteString("\n# ")
				open_by = append(open_by, token.Data)
			} else if token.Data == "h2" {
				toFile.WriteString("\n## ")
				open_by = append(open_by, token.Data)
			} else if token.Data == "h3" {
				toFile.WriteString("\n### ")
				open_by = append(open_by, token.Data)
			} else if token.Data == "h4" {
				toFile.WriteString("\n#### ")
				open_by = append(open_by, token.Data)
			} else if token.Data == "h5" {
				toFile.WriteString("\n##### ")
				open_by = append(open_by, token.Data)
			} else if token.Data == "h6" {
				toFile.WriteString("\n###### ")
				open_by = append(open_by, token.Data)
			} else if token.Data == "p" {
				open_by = append(open_by, token.Data)
			} else if token.Data == "div" {
				open_by = append(open_by, token.Data)
			} else if token.Data == "span" {
				open_by = append(open_by, token.Data)
			} else if token.Data == "a" {
				open_by = append(open_by, token.Data)
			} else if token.Data == "style" {
				blocked_by = append(blocked_by, token.Data)
			} else if token.Data == "script" {
				blocked_by = append(blocked_by, token.Data)
			}
		case typed == html.TextToken:
			if len(open_by) > 0 && len(blocked_by) == 0 {
				text := strings.TrimSpace(token.Data)
				if len(text) > 0 {
					toFile.WriteString(text)
					toFile.WriteString(" ")
				}
			}
		case typed == html.EndTagToken:
			if len(open_by) > 0 && open_by[len(open_by)-1] == token.Data {
				if strings.HasPrefix(token.Data, "h") {
					toFile.WriteString("\n\n")
				} else if token.Data == "p" {
					toFile.WriteString("\n")
				} else if token.Data == "div" {
					toFile.WriteString("\n")
				}
				open_by = open_by[:len(open_by)-1]
			} else if len(blocked_by) > 0 && blocked_by[len(blocked_by)-1] == token.Data {
				blocked_by = blocked_by[:len(blocked_by)-1]
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
	file.WriteString("\n- Time: ")
	file.WriteString(time.Now().UTC().Format(time.RFC3339))
	file.WriteString("\n")
}
