package main

import (
	"os"
	"strings"
	"webloads/lib"
)

func main() {
	var is_input = true
	var is_output = false
	var inputs = []string{}
	var outputs = []string{}
	index := 1
	length := len(os.Args)
	for index < length {
		if os.Args[index] == "-i" || os.Args[index] == "--input" {
			is_input = true
		} else if os.Args[index] == "-o" || os.Args[index] == "--output" {
			is_output = true
		} else if is_input {
			inputs = append(inputs, strings.ToLower(os.Args[index]))
		} else if is_output {
			outputs = append(outputs, strings.ToLower(os.Args[index]))
		}
		index++
	}
	for index, input := range inputs {
		output := "output.md"
		if index < len(outputs) {
			output = outputs[index]
		}
		lib.Parse(input, output)
	}
}
