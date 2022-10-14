package file

import (
	"crypto/md5"
	"io"
	"os"
)

const blockSize = 32 * 1024

func Md5(dir string) ([md5.Size]byte, error) {
	var sum [md5.Size]byte
	f, err := os.Open(dir)
	if err != nil {
		return sum, err
	}
	defer f.Close()
	enc := md5.New()
	_, err = io.Copy(enc, f)
	if err != nil {
		return sum, err
	}
	copy(sum[:], enc.Sum(nil))
	return sum, nil
}

func Size(dir string) (int64, error) {
	fi, err := os.Stat(dir)
	if err != nil {
		return 0, err
	}
	return fi.Size(), nil
}

func Download(dir string, fn func(uint64, []byte)) error {
	f, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer f.Close()
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
			if err == io.EOF {
				err = nil
				more = false
			}
		}
		if err != nil {
			return err
		}
		if n > 0 {
			fn(offset, block[:n])
			offset += uint64(n)
		}
		if !more {
			break
		}
	}
	return nil
}
