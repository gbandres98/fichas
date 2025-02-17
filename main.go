package main

import (
	"log"
	"os"

	"github.com/gbandres98/fichas/internal/fichas"
	"github.com/urfave/cli/v2"
)

var (
	imageSize    int
	imageQuality int
	pagesPerFile int
)

func main() {
	app := &cli.App{
		Name:   "fichas",
		Usage:  "Genera fichas a partir de un fichero Excel",
		Action: app,
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "image-size",
				Aliases:     []string{"s"},
				Usage:       "Tamaño de las imágenes",
				Value:       2500,
				Destination: &imageSize,
			},
			&cli.IntFlag{
				Name:        "image-quality",
				Aliases:     []string{"q"},
				Usage:       "Tamaño de las imágenes",
				Value:       75,
				Destination: &imageQuality,
			},
			&cli.IntFlag{
				Name:        "pages-per-file",
				Aliases:     []string{"p"},
				Usage:       "Número de páginas por fichero",
				Value:       100,
				Destination: &pagesPerFile,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func app(ctx *cli.Context) error {
	file := "fichas.xlsx"
	if ctx.Args().Present() {
		file = ctx.Args().First()
	}

	data, err := fichas.Parse(file)
	if err != nil {
		return err
	}

	err = fichas.ResizeImages(imageSize, imageQuality)
	if err != nil {
		return err
	}

	err = fichas.Generate("template.tex", data, pagesPerFile)
	if err != nil {
		return err
	}

	return nil
}
