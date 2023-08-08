package main

import (
	"github.com/gbandres98/fichas/internal/fichas"
)

func main() {
	data, err := fichas.Parse("test/fichas.xlsx")
	if err != nil {
		panic(err)
	}

	err = fichas.Generate("test/t.tex", data)
	if err != nil {
		panic(err)
	}
}
