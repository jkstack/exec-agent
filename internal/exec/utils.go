package exec

import (
	"encoding/base64"
	"errors"
	"os/exec"
	"runtime"
	"strings"

	"github.com/jkstack/anet"
)

func buildCommand(data *anet.ExecPayload) (*exec.Cmd, string, error) {
	var pass string
	switch {
	case strings.HasPrefix(data.Pass, "$0$"):
		pass = strings.TrimPrefix(data.Pass, "$0$")
	case strings.HasPrefix(data.Pass, "$1$"):
		src, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(data.Pass, "$1$"))
		if err != nil {
			return nil, "", err
		}
		dec, err := anet.Decrypt(src)
		if err != nil {
			return nil, "", err
		}
		pass = string(dec)
	default:
		pass = data.Pass
	}
	if len(pass) > 0 {
		pass += "\n"
	}
	if runtime.GOOS == "windows" {
		data.Auth = ""
	}
	switch strings.ToLower(data.Auth) {
	case "sudo":
		if len(data.User) == 0 {
			return nil, "", errors.New("no user set")
		}
		args := []string{"-S", "-u", data.User, "sh", "-c"}
		args = append(args, data.Cmd+" "+strings.Join(data.Args, " "))
		return exec.Command("sudo", args...), pass, nil
	case "su":
		if len(data.User) == 0 {
			return nil, "", errors.New("no user set")
		}
		args := []string{data.User, "-c",
			data.Cmd + " " + strings.Join(data.Args, " ")}
		return exec.Command("su", args...), pass, nil
	default:
		if runtime.GOOS == "windows" {
			return exec.Command(data.Cmd, data.Args...), "", nil
		}
		args := []string{"-c"}
		args = append(args, data.Cmd+" "+strings.Join(data.Args, " "))
		return exec.Command("sh", args...), "", nil
	}
}
