//go:build !windows && !aix
// +build !windows,!aix

package exec

import (
	"io"
	"os/exec"
	"time"

	"github.com/creack/pty"
)

func createPty(cmd *exec.Cmd, stdin string) (io.ReadCloser, func(), error) {
	tty, err := pty.Start(cmd)
	if err != nil {
		return nil, nil, err
	}
	if len(stdin) > 0 {
		go func() {
			time.Sleep(time.Second)
			tty.WriteString(stdin)
		}()
	}
	return tty, func() {
		tty.Close()
	}, nil
}
