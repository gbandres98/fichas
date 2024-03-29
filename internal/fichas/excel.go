package fichas

import (
	"context"
	"io"
	"log"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var driveService *drive.Service

func init() {
	s, err := drive.NewService(context.Background(), option.WithServiceAccountFile("sa.json"))
	if err != nil {
		panic(err)
	}

	driveService = s

	_, err = s.About.Get().Fields("user").Do()
	if err != nil {
		panic(err)
	}
}

func Parse(path string) ([]pair, error) {
	log.Println("Leyendo excel...")

	data := []pair{}

	f, err := excelize.OpenFile(path)
	if err != nil {
		return data, err
	}
	defer f.Close()

	for _, sheet := range f.GetSheetList() {
		keys := []string{}

		rows, err := f.Rows(sheet)
		if err != nil {
			return data, err
		}

		for rows.Next() {
			values, err := rows.Columns()
			if err != nil {
				return data, err
			}

			if len(values) < 2 {
				continue
			}

			if len(keys) == 0 {
				keys = values
				continue
			}

			p := pair{
				Keys:   keys[:len(keys)-1],
				Values: values[:len(values)-1],
			}

			pictureFile, err := getPicture(values[len(values)-1])
			if err != nil {
				log.Printf("Error al descargar imagen %s: %v\n", values[len(values)-1], err)
			} else {
				p.Img = pictureFile
			}

			data = append(data, p)
		}
	}

	return data, nil
}

func getPicture(url string) (string, error) {
	err := os.MkdirAll("temp/dl", 0777)
	if err != nil {
		return "", err
	}

	s := strings.Split(url, "id=")

	if len(s) < 2 {
		return "", nil
	}

	log.Printf("Descargando imagen %s\n", url)

	id := s[1]

	name := "temp/dl/" + id + ".jpeg"

	file, err := os.Create(name)
	if err != nil {
		return "", err
	}
	defer file.Close()

	res, err := driveService.Files.Get(id).Download()
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	_, err = io.Copy(file, res.Body)
	if err != nil {
		return "", err
	}

	return name, nil
}
