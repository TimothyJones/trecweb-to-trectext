package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"os"
	"strings"
)

const (
	HEADER = iota
	HTTPHEADER
	HTMLCONTENT
)

func main() {
	file, err := os.Open("files.trec")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	current := HEADER

	htmlcontent := ""
	for scanner.Scan() {
		str := scanner.Text()
		switch current {
		case HEADER:
			fmt.Println(str)
		case HTTPHEADER:
		case HTMLCONTENT:
			if str != "</DOC>" {
				htmlcontent += str
			}
		}

		switch str {
		case "<DOC>":
			current = HEADER
		case "</DOCHDR>":
			fmt.Println("<TEXT>")
			current = HTTPHEADER
		case "":
			if current == HTTPHEADER {
				current = HTMLCONTENT
				htmlcontent = ""
			}
		case "</DOC>":

			doc, err := html.Parse(strings.NewReader(htmlcontent))
			if err != nil {
				log.Fatal(err)
			}
			var f func(*html.Node)
			f = func(n *html.Node) {

				if n.Type == html.ElementNode && (n.Data == "script" || n.Data == "style" || n.Data == "noscript") {
					return
				}
				if n.Type == html.TextNode {
					fmt.Println(" " + n.Data)
				}
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					f(c)
				}
			}
			f(doc)
			fmt.Println("</TEXT>")
			fmt.Println("</DOC>")
			current = HEADER
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
