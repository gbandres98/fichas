package fichas

import (
	"context"
	"io"
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

			p.Img, err = getPicture(values[len(values)-1])
			if err != nil {
				return data, err
			}

			data = append(data, p)
		}
	}

	return data, nil
}

func getPicture(url string) (string, error) {
	err := os.MkdirAll("test/dl", 777)
	if err != nil {
		return "", err
	}

	s := strings.Split(url, "id=")

	if len(s) < 2 {
		return "", nil
	}

	id := s[1]

	name := "dl/" + id + ".jpeg"

	file, err := os.Create("test/" + name)
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
