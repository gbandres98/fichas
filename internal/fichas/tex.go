package fichas

import (
	"os"
	"os/exec"
	"text/template"
)

type pair struct {
	Keys   []string
	Values []string
	Img    string
}

func Generate(path string, data []pair) error {
	funcMap := template.FuncMap{
		"inc": func(a int) int {
			return a + 1
		},
	}

	t, err := template.New("t.tex").
		Delims("[[", "]]").
		Funcs(funcMap).
		ParseFiles(path)

	if err != nil {
		return err
	}

	f, err := os.OpenFile("test/res.tex", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	err = t.Execute(f, data)
	if err != nil {
		return err
	}

	// b, err := ioutil.ReadFile(f.Name())
	// if err != nil {
	// 	return err
	// }

	// b = bytes.ReplaceAll(b, []byte("<"), []byte("\\textit{"))
	// b = bytes.ReplaceAll(b, []byte(">"), []byte("}"))

	// _, err = f.Write(b)
	// if err != nil {
	// 	return err
	// }

	exec.Command("pdflatex", "-interaction=nonstopmode", f.Name()).Run()

	return nil
}
