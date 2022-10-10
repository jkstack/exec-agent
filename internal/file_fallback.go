//go:build !windows
// +build !windows

package internal

import (
	"fmt"
	"os"
	"os/user"
	"syscall"
)

func getFileUserGroup(fi os.FileInfo) (string, string) {
	var name, group string
	if stat, ok := fi.Sys().(*syscall.Stat_t); ok {
		u, err := user.LookupId(fmt.Sprintf("%d", stat.Uid))
		if err == nil {
			name = u.Name
		}
		g, err := user.LookupGroupId(fmt.Sprintf("%d", stat.Gid))
		if err == nil {
			group = g.Name
		}
	}
	return name, group
}
