package internal

import (
	"bufio"
	"os"
	"text/template"
)

//go:generate go-bindata -nomemcopy -ignore=\.go$ -pkg internal ./...

func GenerateSiftMain(src string, dst string, data interface{}) error {
	tmpl := template.Must(template.New(src).Parse(string(MustAsset(src + ".tmpl"))))
	f, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	w := bufio.NewWriter(f)
	defer func() {
		_ = w.Flush()
		_ = f.Close()
	}()

	return tmpl.Execute(w, data)
}
