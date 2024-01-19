//go:build windows
// +build windows

package solstice

import (
	"fmt"
	"strings"
)

func HiddenFolderName(f string) string {
	return fmt.Sprintf("_%s", strings.ToLower(strings.TrimLeft(f, "._")))
}
