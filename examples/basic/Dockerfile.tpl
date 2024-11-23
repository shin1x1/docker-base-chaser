FROM php:{{.Tag}}

{{- if (or (matchSemVer "^7.4" .Tag) (eq "7" .Version)) }}
  7.4
{{- else }}
  other
{{- end }}

