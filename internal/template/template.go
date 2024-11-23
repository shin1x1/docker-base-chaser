package template

import (
	"bytes"
	"github.com/Masterminds/semver/v3"
	"io"
	"io/ioutil"
	"os"
	"text/template"
)

type replacement struct {
	Tag                 string
	Version             string
	TemplateDestination string
}

func ExecFile(src, dest, tag, version string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}

	return Exec(f, string(data), tag, version, dest)
}

func ExecString(tmpl, tag, version, dest string) (string, error) {
	var b bytes.Buffer
	if err := Exec(&b, tmpl, tag, version, dest); err != nil {
		return "", err
	}

	return b.String(), nil
}

func Exec(writer io.Writer, tmpl, tag, version, dest string) error {
	t := template.Must(template.New("tmpl").Funcs(funcMap()).Parse(tmpl))

	r := replacement{Tag: tag, Version: version, TemplateDestination: dest}

	return t.Execute(writer, r)
}

func funcMap() template.FuncMap {
	return template.FuncMap{
		"matchSemVer": matchSemVer,
	}
}

func matchSemVer(constraint, version string) (bool, error) {
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return false, err
	}

	v, err := semver.NewVersion(version)
	if err != nil {
		return false, err
	}

	v1, _ := v.SetMetadata("")
	v2, _ := v1.SetPrerelease("")
	return c.Check(&v2), nil
}
