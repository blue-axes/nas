package service

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/blue-axes/tmpl/pkg/constants"
	"github.com/blue-axes/tmpl/pkg/errors"
	"github.com/blue-axes/tmpl/vfs"
	"github.com/google/uuid"
	"io"
	"io/fs"
	"os"
	"path"
)

func (svc *Service) SaveFile(name string, md5Sum string, r io.Reader) error {
	// 先写入临时文件，算出md5，如果一致则再转换为正式文件对象
	tmpFilename := svc.getTempFile()
	f, err := svc.vfs.OpenFile(tmpFilename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	sumReader := io.TeeReader(r, f)
	hash := md5.New()
	_, err = io.Copy(hash, sumReader)
	_ = f.Close()
	if err != nil {
		return nil
	}
	calcSum := hex.EncodeToString(hash.Sum(nil))
	if calcSum != md5Sum {
		// md5不对。 报错，删除文件
		_ = svc.vfs.Remove(tmpFilename)
		return errors.WithCode(constants.ErrFileCheckSumInvalid)
	}
	// 转为正式文件

	return nil
}

func (svc *Service) getTempFile() string {
	tempDir := svc.vfs.TempDir()
	id := uuid.New()
	return path.Join(tempDir, id.String())
}

func (svc *Service) FsStat(name string) (fs.FileInfo, error) {
	return svc.vfs.Stat(name)
}

func (svc *Service) FsOpenFile(name string, flag int) (vfs.File, error) {
	return svc.vfs.OpenFile(name, flag, 0666)
}

func (svc *Service) FsRemove(name string) error {
	return svc.vfs.Remove(name)
}

func (svc *Service) FsMkdirAll(name string) error {
	return svc.vfs.MkdirAll(name, 700)
}

func (svc *Service) FsReadDir(name string) ([]fs.DirEntry, error) {
	return svc.vfs.ReadDir(name)
}
