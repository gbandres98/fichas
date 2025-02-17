package fichas

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/anthonynsimon/bild/transform"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/sync/errgroup"
)

func getPicture(ctx context.Context, p *Pair, url string) error {
	err := os.MkdirAll("temp/dl", 0777)
	if err != nil {
		return err
	}

	id := getID(url)
	if id == "" {
		return fmt.Errorf("No se pudo obtener el id de la imagen %s", url)
	}

	name := "temp/dl/" + id + ".jpeg"
	processedName := "temp/process/" + id + ".jpeg"

	stat, err := os.Stat(name)
	if err == nil && stat.Size() > 0 {
		p.Img = processedName
		return nil
	}

	file, err := os.Create(name)
	if err != nil {
		return err
	}
	defer file.Close()

	res, err := driveService.Files.Get(id).Context(ctx).Download()
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return err
	}

	p.Img = processedName
	return nil
}

func ResizeImages(imageSize int, imageQuality int) error {
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	if err := c.Run(); err != nil {
		log.Printf("Error running command: %v\n", err)
	}

	_ = os.RemoveAll("temp/process")

	err := os.MkdirAll("temp/process", 0777)
	if err != nil {
		return err
	}

	files, err := os.ReadDir("temp/dl")
	if err != nil {
		return err
	}

	ctx := context.Background()

	wg, _ := errgroup.WithContext(ctx)
	wg.SetLimit(5)

	bar := progressbar.Default(int64(len(files)), "Procesando imágenes")

	for _, file := range files {
		go func(fileName string) {
			wg.Go(func() error {
				err := resizeImage("temp/dl/"+fileName, "temp/process/"+fileName, imageSize, imageQuality)
				if err != nil {
					return err
				}

				_ = bar.Add(1)
				return nil
			})
		}(file.Name())
	}

	return wg.Wait()
}

func resizeImage(input string, output string, imageSize int, imageQuality int) error {
	file, err := os.Create(output)
	if err != nil {
		return err
	}
	defer file.Close()

	img, err := os.Open(input)
	if err != nil {
		return err
	}
	defer img.Close()

	i, _, err := image.Decode(img)
	if err != nil {
		return err
	}

	height := imageSize
	width := i.Bounds().Dx() * height / i.Bounds().Dy()

	i = transform.Resize(i, width, height, transform.Linear)

	err = jpeg.Encode(file, i, &jpeg.Options{Quality: imageQuality})
	if err != nil {
		return err
	}

	return nil
}

func getID(url string) string {
	s := strings.Split(url, "id=")

	if len(s) == 2 {
		return s[1]
	}

	s = strings.Split(url, "/d/")

	if len(s) != 2 {
		return ""
	}

	id := s[1]

	s = strings.Split(id, "/")

	if len(s) == 0 {
		return ""
	}

	return s[0]
}
