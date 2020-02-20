package version

import "fmt"

var (
	version = "%VERSION%"
	commit  = "%COMMIT%"
)

func Text(name string) string {
	return fmt.Sprintf("%s version %s (%s)", name, version, commit)
}
