package fichas

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"strings"

	"github.com/anthonynsimon/bild/transform"
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

	img, _, err := image.Decode(res.Body)
	if err != nil {
		return err
	}

	height := 2500
	width := img.Bounds().Dx() * height / img.Bounds().Dy()

	img = transform.Resize(img, width, height, transform.Linear)

	err = jpeg.Encode(file, img, &jpeg.Options{Quality: 75})
	if err != nil {
		return err
	}

	p.Img = name
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
