//go:build windows
// +build windows

package internal

import (
	"os"
)

func getFileUserGroup(fi os.FileInfo) (string, string) {
	return "", ""
}
