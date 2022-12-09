package file

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"strconv"

	"github.com/jkstack/jkframe/compress"
	"github.com/jkstack/jkframe/logging"
)

var cli http.Client

// ReplaceDir replace dir from $$TMP$$ to temporary directory
func ReplaceDir(dir string) string {
	if dir == "$$TMP$$" {
		return os.TempDir()
	}
	return dir
}

// WriteFile write data to file
func WriteFile(dir, data string) error {
	var dec []byte
	if len(data) > 0 {
		var err error
		dec, err = compress.Decompress(data)
		if err != nil {
			return err
		}
	}
	f, err := os.Create(dir)
	if err != nil {
		return err
	}
	defer f.Close()
	n, err := io.Copy(f, bytes.NewReader(dec))
	if err != nil {
		return err
	}
	logging.Info("save %s: %d bytes written", dir, n)
	return nil
}

// DownloadFrom download file from uri
func DownloadFrom(dir, server, uri, token string) error {
	logging.Info("download file from uri %s token: %s", uri, token)
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s%s", server, uri), nil)
	if err != nil {
		return err
	}
	if len(token) > 0 {
		req.Header.Set("X-Token", token)
	}
	rep, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer rep.Body.Close()
	if rep.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(rep.Body)
		return fmt.Errorf("http_code=%d, data=%s", rep.StatusCode, string(data))
	}
	f, err := os.Create(dir)
	if err != nil {
		return err
	}
	defer f.Close()
	n, err := io.Copy(f, rep.Body)
	if err != nil {
		return err
	}
	logging.Info("save %s: %d bytes written", dir, n)
	return nil
}

// Chmod chmod
func Chmod(dir string, mode os.FileMode) error {
	return os.Chmod(dir, mode)
}

// Chown chown
func Chown(dir, u, g string) error {
	uid := -1
	gid := -1
	if len(u) > 0 {
		u, err := user.Lookup(u)
		if err == nil {
			n, err := strconv.ParseInt(u.Uid, 10, 64)
			if err == nil {
				uid = int(n)
			}
		}
	}
	if len(g) > 0 {
		g, err := user.LookupGroup(g)
		if err == nil {
			n, err := strconv.ParseInt(g.Gid, 10, 64)
			if err == nil {
				gid = int(n)
			}
		}
	}
	return os.Chown(dir, uid, gid)
}
