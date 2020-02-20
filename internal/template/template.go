package template

import (
	"bytes"
	"io/ioutil"
	"os"
	"text/template"
)

type replacement struct {
	Tag                 string
	TemplateDestination string
}

func ExecFile(src, dest, tag string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	t := template.Must(template.New("tmpl").Parse(string(data)))

	r := replacement{Tag: tag}
	return t.Execute(f, r)
}

func Exec(tmpl, tag, dest string) (string, error) {
	t := template.Must(template.New("tmpl").Parse(tmpl))

	r := replacement{Tag: tag, TemplateDestination: dest}

	var b bytes.Buffer
	if err := t.Execute(&b, r); err != nil {
		return "", nil
	}

	return b.String(), nil
}
