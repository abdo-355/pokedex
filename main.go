package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("Hello, World!")
}

func cleanInput(text string) []string {
	cleaned := strings.TrimSpace(text)

	split := strings.Split(cleaned, " ")
	fmt.Println(split)

	var removedZero []string

	for _, v := range split {
		c := strings.TrimSpace(v)
		if c != "" {
			removedZero = append(removedZero, c)
		}
	}

	return removedZero
}
