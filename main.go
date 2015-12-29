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
	DOCHEADER
	HTMLCONTENT
	HEADER_DOCID
)

func main() {
	sOut, err := os.Open("sOut")
	if err != nil {
		log.Fatal(err)
	}
	defer sOut.Close()

	file, err := os.Open("files.trec")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	current := HEADER

	var (
		htmlcontent, docno, doctitle, docurl string
	)

	for scanner.Scan() {
		str := scanner.Text()
		switch current {
		case HEADER_DOCID:
			fmt.Println(str)
			docno = str[7:][:25]
			current = HEADER
		case HEADER:
			fmt.Println(str)
		case DOCHEADER:
			fmt.Println(str)
			if str != "</DOCHDR>" {
				docurl = str
			}
		case HTTPHEADER:
		case HTMLCONTENT:
			if str != "</DOC>" {
				htmlcontent += str
			}
		}

		switch str {
		case "<DOC>":
			current = HEADER_DOCID
		case "<DOCHDR>":
			current = DOCHEADER
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
				if n.Type == html.ElementNode && n.Data == "title" {
					if n.FirstChild != nil {
						doctitle = n.FirstChild.Data
					}
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
			fmt.Fprintln(os.Stderr, docno)
			if doctitle == "" {
				fmt.Fprintln(os.Stderr, "(no title)")
			} else {
				fmt.Fprintln(os.Stderr, doctitle)
			}
			fmt.Fprintln(os.Stderr, docurl)
			docno = ""
			docurl = ""
			doctitle = ""
			fmt.Println("</TEXT>")
			fmt.Println("</DOC>")
			current = HEADER
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
