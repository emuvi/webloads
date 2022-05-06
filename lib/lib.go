package lib

import (
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func Parse(input string, output string) {
	resp, err := http.Get(input)
	if err != nil {
		panic(err)
	}
	work, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	Catch(resp.Body, work)
}

func Catch(fromReader io.ReadCloser, toWriter *os.File) {
	defer fromReader.Close()
	defer toWriter.Close()
	tkns := html.NewTokenizer(fromReader)
	var open_by = []string{}
	var blocked_by = []string{}
	for {
		tknType := tkns.Next()
		if tknType == html.ErrorToken {
			return
		}
		tkn := tkns.Token()
		switch {
		case tknType == html.StartTagToken:
			if tkn.Data == "h1" {
				toWriter.WriteString("\n# ")
				open_by = append(open_by, tkn.Data)
			} else if tkn.Data == "h2" {
				toWriter.WriteString("\n## ")
				open_by = append(open_by, tkn.Data)
			} else if tkn.Data == "h3" {
				toWriter.WriteString("\n### ")
				open_by = append(open_by, tkn.Data)
			} else if tkn.Data == "h4" {
				toWriter.WriteString("\n#### ")
				open_by = append(open_by, tkn.Data)
			} else if tkn.Data == "h5" {
				toWriter.WriteString("\n##### ")
				open_by = append(open_by, tkn.Data)
			} else if tkn.Data == "h6" {
				toWriter.WriteString("\n###### ")
				open_by = append(open_by, tkn.Data)
			} else if tkn.Data == "p" {
				open_by = append(open_by, tkn.Data)
			} else if tkn.Data == "div" {
				open_by = append(open_by, tkn.Data)
			} else if tkn.Data == "span" {
				open_by = append(open_by, tkn.Data)
			} else if tkn.Data == "a" {
				open_by = append(open_by, tkn.Data)
			} else if tkn.Data == "style" {
				blocked_by = append(blocked_by, tkn.Data)
			} else if tkn.Data == "script" {
				blocked_by = append(blocked_by, tkn.Data)
			}
		case tknType == html.TextToken:
			if len(open_by) > 0 && len(blocked_by) == 0 {
				text := strings.TrimSpace(tkn.Data)
				if len(text) > 0 {
					toWriter.WriteString(text)
					toWriter.WriteString(" ")
				}
			}
		case tknType == html.EndTagToken:
			if len(open_by) > 0 && open_by[len(open_by)-1] == tkn.Data {
				if strings.HasPrefix(tkn.Data, "h") {
					toWriter.WriteString("\n\n")
				} else if tkn.Data == "p" {
					toWriter.WriteString("\n")
				} else if tkn.Data == "div" {
					toWriter.WriteString("\n")
				}
				open_by = open_by[:len(open_by)-1]
			} else if len(blocked_by) > 0 && blocked_by[len(blocked_by)-1] == tkn.Data {
				blocked_by = blocked_by[:len(blocked_by)-1]
			}
		}
	}
}
