package goapimaker

import (
	"fmt"
)

func WriteGoMod(projectName string) string {
	return fmt.Sprintf(`
module %s

go 1.25.6

require github.com/lib/pq v1.12.3 // indirect
	`, projectName)
}

func WriteSum() string {
	return `
github.com/lib/pq v1.12.3 h1:tTWxr2YLKwIvK90ZXEw8GP7UFHtcbTtty8zsI+YjrfQ=
github.com/lib/pq v1.12.3/go.mod h1:/p+8NSbOcwzAEI7wiMXFlgydTwcgTr3OSKMsD2BitpA=
	`
}
