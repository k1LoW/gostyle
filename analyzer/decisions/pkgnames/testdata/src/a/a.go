package a

import (
	renamed_log "log" // want "gostyle.pkgnames"

	_ "embed"
)

func f() {
	renamed_log.Println("hello")
}
