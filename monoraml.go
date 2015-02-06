package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func includeblock(path string, name string, count int) string {

	incfile, err := os.Open(path + name)
	check(err)
	defer incfile.Close()

	indent := strings.Repeat(" ", count+1)
	var block []string
	scanner := bufio.NewScanner(incfile)
	for scanner.Scan() {
		inline := scanner.Text()
		block = append(block, indent, inline, "\n")
	}

	if err := scanner.Err(); err != nil {
		check(err)
	}
	return strings.Join(block, "")
}

func main() {
	argsWithProg := os.Args
	path := ""
	filename := ""
	if len(argsWithProg) > 1 {
		filename = os.Args[1]
	} else {
		fmt.Printf("Usage:\n%v filename <path>\n\nFilename without extension, always uses .raml\nIf no path is given the current dir is assumed\n\n", os.Args[0])
		os.Exit(1)
	}
	if len(argsWithProg) > 2 {
		path = os.Args[2]
	}
	ramlfile, err := os.Open(path + filename + ".raml")
	check(err)
	defer ramlfile.Close()
	outfile, err := os.Create(path + filename + "-mono.raml")
	check(err)
	defer outfile.Close()

	scanner := bufio.NewScanner(ramlfile)
	for scanner.Scan() {
		inline := scanner.Text()
		if strings.Contains(inline, "!include") {
			if strings.Contains(inline, "example:") {
				count := 0
				where := 0
				spl := strings.Split(inline, " ")
				for _, s := range spl {
					if len(s) == 0 {
						count++
					} else {
						where++
					}

					if where == 1 {
						if s == "#" {
							break
						}
						indent := strings.Repeat(" ", count)
						out := fmt.Sprintf("%v%v |\n", indent, s)
						_, err = outfile.WriteString(string(out))
						check(err)
					}
					if where == 3 {
						out := fmt.Sprintf("%v", includeblock(path, s, count))
						_, err = outfile.WriteString(string(out))
						check(err)
					}
				}
			}
		} else {
			out := fmt.Sprintf("%v\n", inline)
			_, err = outfile.WriteString(string(out))
			check(err)
		}
	}

	if err := scanner.Err(); err != nil {
		check(err)
	}
}
