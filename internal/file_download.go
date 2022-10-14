package internal

import (
	"crypto/md5"
	"exec/internal/file"

	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/compress"
	"github.com/jkstack/jkframe/logging"
)

func (agent *Agent) Download(taskID, dir string) error {
	logging.Info("download %s...", dir)
	sum, err := file.Md5(dir)
	if err != nil {
		logging.Error("md5sum(%s): %v", dir, err)
		agent.downloadErr(taskID, dir, err.Error())
		return nil
	}
	size, err := file.Size(dir)
	if err != nil {
		logging.Error("file size(%s): %v", dir, err)
		agent.downloadErr(taskID, dir, err.Error())
		return nil
	}
	agent.downloadOK(taskID, dir, sum, uint64(size))
	err = file.Download(dir, func(offset uint64, data []byte) {
		agent.downloadData(taskID, offset, data)
	})
	if err != nil {
		logging.Error("download(%s): %v", dir, err)
		agent.downloadDataError(taskID, err.Error())
		return nil
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
		Data:   compress.Compress(data),
	}
	agent.chWrite <- &msg
}
