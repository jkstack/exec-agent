//go:build windows || aix
// +build windows aix

package exec

import (
	"io"
	"os/exec"
	"strings"
)

func createPty(cmd *exec.Cmd, stdin string) (io.ReadCloser, func(), error) {
	cmd.Stdin = strings.NewReader(stdin)
	pr, pw := io.Pipe()
	cmd.Stdout = pw
	cmd.Stderr = pw
	err := cmd.Start()
	if err != nil {
		pr.Close()
		pw.Close()
		return nil, nil, err
	}
	return pr, func() {
		pr.Close()
		pw.Close()
	}, nil
}
