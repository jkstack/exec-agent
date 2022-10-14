package internal

import (
	"bytes"
	"errors"
	"exec/internal/file"
	"os"
	"path/filepath"

	"github.com/jkstack/anet"
	"github.com/jkstack/jkframe/logging"
)

var errInvalidChecksum = errors.New("invalid checksum")

func (agent *Agent) Upload(taskID string, msg *anet.Msg) (ret error) {
	logging.Info("upload %s...", msg.Upload.Dir)
	msg.Upload.Dir = file.ReplaceDir(msg.Upload.Dir)

	os.MkdirAll(msg.Upload.Dir, 0755)
	logging.Info("mkdir %s...", msg.Upload.Dir)

	dir := filepath.Join(msg.Upload.Dir, msg.Upload.Name)

	defer func() {
		if ret != nil {
			os.Remove(dir)
			agent.uploadErr(taskID, dir, ret.Error())
			ret = nil
			return
		}
		agent.uploadOK(taskID, dir)
	}()

	var err error
	if msg.Upload.Size > 0 {
		if len(msg.Upload.Data) > 0 {
			err = file.WriteFile(dir, msg.Upload.Data)
		} else if len(msg.Upload.URI) > 0 {
			err = file.DownloadFrom(dir, agent.cfg.Server, msg.Upload.URI, msg.Upload.Token)
		}
		if err != nil {
			logging.Error("save(%s): %v", dir, err)
			return err
		}
		enc, err := file.Md5(dir)
		if err != nil {
			logging.Error("md5(%s): %v", dir, err)
			return err
		}
		if !bytes.Equal(enc[:], msg.Upload.MD5[:]) {
			logging.Error("invalid checksum(%s)", dir)
			return errInvalidChecksum
		}
	}

	logging.Info("chmod for file(%s): %s", dir, msg.Upload.Mod.String())
	err = file.Chmod(dir, msg.Upload.Mod)
	if err != nil {
		logging.Warning("chmod(%s): %v", dir, err)
		return nil
	}

	if len(msg.Upload.OwnUser) > 0 ||
		len(msg.Upload.OwnGroup) > 0 {
		err = file.Chown(dir, msg.Upload.OwnUser, msg.Upload.OwnGroup)
		if err != nil {
			logging.Warning("chown(%s): %v", dir, err)
			return nil
		}
	}

	return nil
}

func (agent *Agent) uploadErr(taskID, dir, msg string) {
	var m anet.Msg
	m.Type = anet.TypeUploadRep
	m.TaskID = taskID
	m.UploadRep = &anet.UploadRep{
		Dir:    dir,
		OK:     false,
		ErrMsg: msg,
	}
	agent.chWrite <- &m
}

func (agent *Agent) uploadOK(taskID, dir string) {
	var m anet.Msg
	m.Type = anet.TypeUploadRep
	m.TaskID = taskID
	m.UploadRep = &anet.UploadRep{
		Dir: dir,
		OK:  true,
	}
	agent.chWrite <- &m
}
