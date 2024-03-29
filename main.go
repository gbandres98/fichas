package main

import (
	"github.com/gbandres98/fichas/internal/fichas"
)

func main() {
	data, err := fichas.Parse("fichas.xlsx")
	if err != nil {
		panic(err)
	}

	err = fichas.Generate("template.tex", data)
	if err != nil {
		panic(err)
	}
}
