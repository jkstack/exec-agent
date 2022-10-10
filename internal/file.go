package internal

import (
	"crypto/md5"
	"exec/internal/file"
	"exec/internal/utils"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
)

const blockSize = 32 * 1024

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

func (agent *Agent) Download(taskID, dir string) error {
	logging.Info("download %s...", dir)
	f, err := os.Open(dir)
	if err != nil {
		logging.Error("open(%s): %v", dir, err)
		agent.downloadErr(taskID, dir, err.Error())
		return nil
	}
	defer f.Close()
	enc := md5.New()
	_, err = io.Copy(enc, f)
	if err != nil {
		logging.Error("md5(%s): %v", dir, err)
		agent.downloadErr(taskID, dir, err.Error())
		return nil
	}
	fi, err := f.Stat()
	if err != nil {
		logging.Error("stat(%s): %v", dir, err)
		agent.downloadErr(taskID, dir, err.Error())
		return nil
	}
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		logging.Error("seek(%s): %v", dir, err)
		agent.downloadErr(taskID, dir, err.Error())
		return nil
	}
	var sum [md5.Size]byte
	copy(sum[:], enc.Sum(nil))
	agent.downloadOK(taskID, dir, sum, uint64(fi.Size()))
	block := make([]byte, blockSize)
	var offset uint64
	for {
		more := true
		n, err := io.ReadFull(f, block)
		if err != nil {
			if err == io.ErrUnexpectedEOF {
				err = nil
				more = false
			}
		}
		if err != nil {
			logging.Error("read(%s): %v", dir, err)
			agent.downloadDataError(taskID, err.Error())
			return nil
		}
		if n > 0 {
			agent.downloadData(taskID, offset, block[:n])
			offset += uint64(n)
		}
		if !more {
			break
		}
	}
	return nil
}

func (agent *Agent) downloadErr(taskID, dir, msg string) {
	var m anet.Msg
	m.Type = anet.TypeDownloadRep
	m.TaskID = taskID
	m.DownloadRep = &anet.DownloadRep{
		Dir:    dir,
		OK:     false,
		ErrMsg: msg,
	}
	agent.chWrite <- &m
}

func (agent *Agent) downloadOK(taskID, dir string, md5 [md5.Size]byte, size uint64) {
	var msg anet.Msg
	msg.Type = anet.TypeDownloadRep
	msg.TaskID = taskID
	msg.DownloadRep = &anet.DownloadRep{
		Dir:       dir,
		OK:        true,
		Size:      size,
		BlockSize: blockSize,
		MD5:       md5,
	}
	agent.chWrite <- &msg
}

func (agent *Agent) downloadDataError(taskID, msg string) {
	var m anet.Msg
	m.Type = anet.TypeDownloadError
	m.TaskID = taskID
	m.DownloadError = &anet.DownloadError{Msg: msg}
	agent.chWrite <- &m
}

func (agent *Agent) downloadData(taskID string, offset uint64, data []byte) {
	var msg anet.Msg
	msg.Type = anet.TypeDownloadData
	msg.TaskID = taskID
	msg.DownloadData = &anet.DownloadData{
		Offset: offset,
		Data:   utils.EncodeData(data),
	}
	agent.chWrite <- &msg
}
