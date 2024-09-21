package fichas

import (
	"log"
	"os"
	"os/exec"
	"text/template"
)

type Pair struct {
	Keys   []string
	Values []string
	Img    string
}

func Generate(path string, data []*Pair) error {
	log.Println("Generando .tex")
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

	f, err := os.OpenFile("temp/res.tex", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	err = t.Execute(f, data)
	if err != nil {
		return err
	}

	log.Println("Generando pdf")
	exec.Command("pdflatex", "-interaction=nonstopmode", f.Name()).Run()

	return nil
}
