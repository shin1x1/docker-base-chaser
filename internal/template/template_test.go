package template

import (
	"strings"
	"testing"
)

func TestTemplate_Exec(t *testing.T) {
	tmpl := `
Tag: {{.Tag}}
{{- if (matchSemVer "^7.4" .Tag) }}
^7.4
{{- else }}
other
{{- end }}

{{- if or (matchSemVer "^7.4" .Tag) (eq "7" .Version) }}
^7.4 or 7
{{- else }}
other
{{- end }}
`
	params := []struct {
		tag      string
		version  string
		expected string
	}{
		{
			tag:      "7.4.0-cli",
			version:  "7.4.0",
			expected: "Tag: 7.4.0-cli\n^7.4\n^7.4 or 7",
		},
		{
			tag:      "7.4-cli",
			version:  "7.4",
			expected: "Tag: 7.4-cli\n^7.4\n^7.4 or 7",
		},
		{
			tag:      "7-cli",
			version:  "7",
			expected: "Tag: 7-cli\nother\n^7.4 or 7",
		},
		{
			tag:      "7.3.0-cli",
			version:  "7.3.0",
			expected: "Tag: 7.3.0-cli\nother\nother",
		},
		{
			tag:      "7.5.0-cli",
			version:  "7.5.0",
			expected: "Tag: 7.5.0-cli\n^7.4\n^7.4 or 7",
		},
	}

	for k, v := range params {
		result, err := ExecString(tmpl, v.tag, v.version, "")
		if err != nil {
			t.Errorf("%+v", err)
		} else if strings.Trim(result, "\n") != v.expected {
			t.Errorf("key:%d expected:\n%s\n\nactual:\n%s", k, v.expected, result)
		}
	}
}
