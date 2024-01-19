//go:build ignore
// +build ignore

package main

import (
	"bufio"
	"fmt"
	"html/template"
	"io"
	"os"
	"strconv"
	"strings"
)

const templateString = `package {{.Package}}

// This file is generated !
// DO NOT EDIT

import "fyne.io/fyne/v2"

func GetMetadata() fyne.AppMetadata {
	md := App.Metadata()
	if md.ID == "" {
		md = fyne.AppMetadata {
			ID: "{{.ID}}",
			Name: "{{.Name}}",
			Version: "{{.Version}}",
			Build: {{.Build}},
			Icon: {{.Icon}},
			Release: false,
			Custom: map[string]string{},
		}
	}
	return md
}
`

type meta struct {
	Package string
	ID      string
	Name    string
	Icon    string
	Version string
	Build   int
}

func main() {
	inFile := os.Args[1]
	outFile := os.Args[2]
	pkg := os.Args[3]

	meta := parseInFile(inFile)
	meta.Package = pkg
	writeOutFile(outFile, meta)
}

func parseInFile(file string) meta {
	fh, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer fh.Close()

	meta := meta{
		Icon: "nil",
	}
	inDetails := false
	reader := bufio.NewReader(fh)
	for {
		str, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			panic(err)
		}
		str = strings.TrimRight(str, "\n")

		if !inDetails {
			if str == "[Details]" {
				inDetails = true
			}
			continue
		}

		if strings.HasPrefix(str, "[") {
			break
		}

		if strings.HasPrefix(str, "ID") {
			meta.ID = strings.Trim(strings.Split(str, "=")[1], " \"")
			continue
		}
		if strings.HasPrefix(str, "Name") {
			meta.Name = strings.Trim(strings.Split(str, "=")[1], " \"")
			continue
		}
		if strings.HasPrefix(str, "Icon") {
			icon := strings.Trim(strings.Split(str, "=")[1], " \"")
			if icon != "" {
				meta.Icon = icon
			}
			continue
		}
		if strings.HasPrefix(str, "Version") {
			meta.Version = strings.Trim(strings.Split(str, "=")[1], " \"")
			continue
		}
		if strings.HasPrefix(str, "Build") {
			meta.Build, _ = strconv.Atoi(strings.Trim(strings.Split(str, "=")[1], " "))
			continue
		}
	}
	return meta
}

func writeOutFile(file string, meta meta) {
	fh, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer fh.Close()

	writer := bufio.NewWriter(fh)
	defer writer.Flush()

	templ := template.Must(template.New("metadata").Parse(templateString))
	templ.Execute(writer, &meta)
}
