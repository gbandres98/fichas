package fichas

import (
	"context"
	"log"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"sync/atomic"
	"time"

	"github.com/xuri/excelize/v2"
	"golang.org/x/sync/errgroup"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var driveService *drive.Service

func init() {
	s, err := drive.NewService(context.Background(), option.WithCredentialsFile("sa.json"))
	if err != nil {
		panic(err)
	}

	driveService = s

	_, err = s.About.Get().Fields("user").Do()
	if err != nil {
		panic(err)
	}
}

func Parse(path string) ([]*Pair, error) {
	log.Println("Leyendo excel...")

	data := []*Pair{}

	f, err := excelize.OpenFile(path)
	if err != nil {
		return data, err
	}
	defer f.Close()

	ctx := context.Background()

	total := atomic.Int32{}
	downloaded := atomic.Int32{}
	wg, ctx := errgroup.WithContext(ctx)
	wg.SetLimit(20)

	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)

		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				c := exec.Command("clear")
				c.Stdout = os.Stdout
				if err := c.Run(); err != nil {
					log.Printf("Error running command: %v\n", err)
				}
				log.Printf("Descargando imágenes... %d/%d\n", downloaded.Load(), total.Load())
			}
		}
	}()

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

			p := &Pair{
				Keys:   keys[:len(keys)-1],
				Values: values[:len(values)-1],
			}

			go func(p *Pair, url string) {
				total.Add(1)

				wg.Go(func() error {
					err = getPicture(ctx, p, values[len(values)-1])
					if err != nil {
						return err
					}

					downloaded.Add(1)
					return nil
				})
			}(p, values[len(values)-1])

			data = append(data, escapeRow(p))
		}

		if err := wg.Wait(); err != nil {
			return nil, err
		}

		return data, nil
	}

	return data, nil
}

func escapeRow(row *Pair) *Pair {
	for i, v := range row.Values {
		row.Values[i] = escapeString(v)
	}

	return row
}

func escapeString(s string) string {
	s = strings.ReplaceAll(s, `\`, `\textbackslash`)
	s = strings.ReplaceAll(s, `&`, `\&`)
	s = strings.ReplaceAll(s, `%`, `\%`)
	s = strings.ReplaceAll(s, `$`, `\$`)
	s = strings.ReplaceAll(s, `#`, `\#`)
	s = strings.ReplaceAll(s, `_`, `\_`)
	s = strings.ReplaceAll(s, `{`, `\{`)
	s = strings.ReplaceAll(s, `}`, `\}`)
	s = strings.ReplaceAll(s, `~`, `\textasciitilde`)
	s = strings.ReplaceAll(s, `^`, `\textasciicircum`)
	s = replaceQuotes(s)
	s = strings.ReplaceAll(s, `<`, `\textit{`)
	s = strings.ReplaceAll(s, `>`, `}`)

	if _, err := url.ParseRequestURI(s); err == nil {
		s = `\url{` + s + `}`
	}

	return s
}

func replaceQuotes(s string) string {
	open := false
	for _, r := range s {
		if r == '"' {
			if open {
				s = strings.Replace(s, `"`, "''", 1)
			} else {
				s = strings.Replace(s, `"`, "``", 1)
			}

			open = !open
		}
	}

	return s
}
