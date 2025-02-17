package fichas

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"text/template"
)

type Pair struct {
	Keys   []string
	Values []string
	Img    string
}

func Generate(path string, data []*Pair, pagesPerFile int) error {
	log.Println("Generando .tex")

	fileNo := 0

	for i := 0; i < len(data); i += pagesPerFile {
		fileNo++
		var err error

		if i >= len(data) {
			break
		}

		if i+pagesPerFile > len(data) {
			err = generateFile(path, data[i:], fileNo)
		} else {
			err = generateFile(path, data[i:i+pagesPerFile], fileNo)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func generateFile(path string, data []*Pair, fileIndex int) error {
	log.Printf("Generando pdf %d\n", fileIndex)

	funcMap := template.FuncMap{
		"inc": func(a int) int {
			return a + 1
		},
	}

	t, err := template.New("template.tex").
		Delims("[[", "]]").
		Funcs(funcMap).
		ParseFiles(path)

	if err != nil {
		return err
	}

	f, err := os.OpenFile("temp/res-"+strconv.Itoa(fileIndex)+".tex", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	err = t.Execute(f, data)
	if err != nil {
		return err
	}

	err = exec.Command("pdflatex", "-interaction=nonstopmode", f.Name()).Run()
	if err != nil {
		return err
	}

	log.Println("res-" + strconv.Itoa(fileIndex) + ".pdf")

	return nil
}
