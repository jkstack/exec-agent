package file

import (
	"os"
	"runtime"
	"time"

	"github.com/jkstack/jkframe/logging"
	"github.com/shirou/gopsutil/v3/disk"
)

type windowsFileInfo struct {
	dir string
}

func (i windowsFileInfo) Name() string {
	return i.dir
}

func (i windowsFileInfo) Size() int64 {
	return 0
}

func (i windowsFileInfo) Mode() os.FileMode {
	return os.ModeDir
}

func (i windowsFileInfo) ModTime() time.Time {
	return time.Time{}
}

func (i windowsFileInfo) IsDir() bool {
	return true
}

func (i windowsFileInfo) Sys() interface{} {
	return nil
}

// Ls handle ls command
func Ls(dir string) ([]os.FileInfo, error) {
	logging.Info("ls %s", dir)
	if runtime.GOOS == "windows" && dir == "/" {
		parts, err := disk.Partitions(false)
		if err != nil {
			return nil, err
		}
		files := make([]os.FileInfo, len(parts))
		for i, part := range parts {
			files[i] = windowsFileInfo{part.Mountpoint}
		}
		return files, nil
	}
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	ret := make([]os.FileInfo, len(files))
	for i, file := range files {
		ret[i], err = file.Info()
		if err != nil {
			continue
		}
	}
	return ret, nil
}
