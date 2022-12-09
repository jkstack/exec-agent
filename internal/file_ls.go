package internal

import (
	"exec/internal/file"
	"os"
	"path"
	"path/filepath"

	"github.com/jkstack/anet"
)

const blockSize = 32 * 1024

// Ls handle ls command
func (agent *Agent) Ls(taskID, dir string) error {
	files, err := file.Ls(dir)
	if err != nil {
		agent.lsError(taskID, dir, err.Error())
		return nil
	}
	agent.lsOK(taskID, dir, files)
	return nil
}

func (agent *Agent) lsError(taskID, dir, msg string) {
	var m anet.Msg
	m.TaskID = taskID
	m.Type = anet.TypeLsRep
	m.LSRep = &anet.LsRep{
		Dir:    dir,
		OK:     false,
		ErrMsg: msg,
	}
	agent.chWrite <- &m
}

func (agent *Agent) lsOK(taskID, dir string, files []os.FileInfo) {
	fs := make([]anet.FileInfo, len(files))
	for i, file := range files {
		user, group := getFileUserGroup(file)
		fs[i] = anet.FileInfo{
			Name:    path.Base(file.Name()),
			Mod:     file.Mode(),
			User:    user,
			Group:   group,
			Size:    uint64(file.Size()),
			ModTime: file.ModTime(),
			IsLink:  file.Mode()&os.ModeSymlink == os.ModeSymlink,
		}
		if fs[i].IsLink {
			fi, err := os.Stat(path.Join(dir, file.Name()))
			if err != nil {
				continue
			}
			if fi.IsDir() {
				fs[i].Mod |= os.ModeDir
			}
			fs[i].LinkDir, _ = os.Readlink(filepath.Join(dir, file.Name()))
		}
	}

	var m anet.Msg
	m.Type = anet.TypeLsRep
	m.TaskID = taskID
	m.LSRep = &anet.LsRep{
		Dir:   dir,
		OK:    true,
		Files: fs,
	}
	agent.chWrite <- &m
}
