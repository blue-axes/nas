package service

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/blue-axes/tmpl/pkg/context"
	"github.com/blue-axes/tmpl/store/rdb"
	"github.com/blue-axes/tmpl/types"
	"github.com/google/uuid"
	"io"
	"os"
	"path"
	"strings"
)

func (svc *Service) SimpleSaveFile(ctx *context.Context, name string, r io.Reader, overwrite bool) error {
	var (
		hash     = md5.New()
		reader   = io.TeeReader(r, hash)
		realName = svc.simpleGetRealFilename(name)
	)
	err := svc.vfs.MkdirAll(path.Dir(realName), 0700)
	if err != nil {
		return err
	}
	// 查询记录是否存在
	fileInfo, err := svc.store.RDB().GetFileByName(ctx, name)
	if err == nil && overwrite {
		// 文件存在。就重写
		f, err := svc.vfs.OpenFile(realName, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		size, err := io.Copy(f, reader)
		_ = f.Close()
		if err != nil {
			return err
		}
		md5sum := hex.EncodeToString(hash.Sum(nil))

		err = svc.store.RDB().UpdateFileByID(ctx, fileInfo.ID, &types.File{
			ID:   fileInfo.ID,
			Name: name,
			Ext:  path.Ext(name),
			Path: realName,
			Size: uint64(size),
			Md5:  md5sum,
		})
		if err != nil {
			return err
		}
	}

	f, err := svc.vfs.OpenFile(realName, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	size, err := io.Copy(f, reader)
	_ = f.Close()
	if err != nil {
		return err
	}
	md5sum := hex.EncodeToString(hash.Sum(nil))
	// 数据库中插入记录
	err = svc.store.RDB().CreateFile(ctx, &types.File{
		Name: name,
		Ext:  path.Ext(name),
		Path: realName,
		Size: uint64(size),
		Md5:  md5sum,
	})
	if err != nil {
		return err
	}
	return nil
}

func (svc *Service) simpleGetRealFilename(name string) string {
	name = path.Clean(name)
	baseName := path.Base(name)
	baseDir := path.Dir(name)
	filename := ""
	switch svc.cfg.Nas.RealFilenamePolicy {
	case types.RFNP_UUID:
		id, err := uuid.NewUUID()
		if err != nil {
			filename = ""
		} else {
			filename = id.String()
		}
	case types.RFNP_Origin:
		filename = baseName
	case types.RFNP_Underline:
		baseDir = ""
		filename = strings.ReplaceAll(name, "/", "_")
	}

	return path.Join(svc.cfg.Nas.SimpleUploadRoot, baseDir, filename)
}

func (svc *Service) SimpleGetFileInfo(ctx *context.Context, name string) (*types.File, error) {
	return svc.store.RDB().GetFileByName(ctx, name)
}

func (svc *Service) SimpleGetFileContent(ctx *context.Context, name string) (io.ReadSeekCloser, *types.File, error) {
	fileInfo, err := svc.store.RDB().GetFileByName(ctx, name)
	if err != nil {
		return nil, nil, err
	}
	rc, err := svc.vfs.OpenFile(fileInfo.Path, os.O_RDONLY, 0600)
	if err != nil {
		return nil, nil, err
	}
	return rc, fileInfo, nil
}

func (svc *Service) SimpleDeleteFile(ctx *context.Context, name string) error {
	info, err := svc.store.RDB().GetFileByName(ctx, name)
	if err != nil {
		return err
	}
	err = svc.store.RDB().Transaction(func(store rdb.TxStore) error {
		err = svc.vfs.Remove(info.Path)
		if err != nil {
			return err
		}
		return store.DeleteByName(ctx, name)
	})
	return err
}

func (svc *Service) SimpleListFiles(ctx *context.Context, filePath string) ([]types.File, error) {
	var cond *types.Condition = nil
	if filePath != "" {
		filePath = strings.TrimRight(filePath, "/")
		tmp := types.ConditionOr(types.ConditionNew("name = ?", filePath),
			types.ConditionNew("name LIKE ?", filePath+"/%"))
		cond = &tmp
	}
	return svc.store.RDB().ListFile(ctx, cond, nil)
}
