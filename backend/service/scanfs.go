package service

import (
	"path"
	"strings"

	"github.com/blue-axes/tmpl/pkg/context"
	"github.com/blue-axes/tmpl/types"
	log "github.com/sirupsen/logrus"
)

type ScanResult struct {
	Total   int `json:"Total"`
	New     int `json:"New"`
	Skipped int `json:"Skipped"`
}

func (svc *Service) ScanFiles(ctx *context.Context) (*ScanResult, error) {
	root := svc.cfg.Nas.SimpleUploadRoot
	result := &ScanResult{}
	err := svc.scanDir(ctx, root, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (svc *Service) scanDir(ctx *context.Context, dirPath string, result *ScanResult) error {
	entries, err := svc.vfs.ReadDir(dirPath)
	if err != nil {
		log.WithError(err).Warnf("scanfs: cannot read dir %s", dirPath)
		return nil
	}

	for _, entry := range entries {
		fullPath := path.Join(dirPath, entry.Name())

		// 跳过隐藏文件和临时文件
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			log.WithError(err).Warnf("scanfs: cannot stat %s", fullPath)
			continue
		}

		name := strings.TrimPrefix(fullPath, svc.cfg.Nas.SimpleUploadRoot)
		name = strings.TrimPrefix(name, "/")

		result.Total++

		// 检查数据库是否存在
		_, dbErr := svc.store.RDB().GetFileByPath(ctx, fullPath)
		if dbErr == nil {
			result.Skipped++
		} else {
			ext := ""
			if !info.IsDir() {
				ext = path.Ext(entry.Name())
			}
			fileRecord := &types.File{
				Name:  name,
				Ext:   ext,
				Path:  fullPath,
				Size:  uint64(info.Size()),
				IsDir: info.IsDir(),
			}
			if err := svc.store.RDB().CreateFile(ctx, fileRecord); err != nil {
				log.WithError(err).Warnf("scanfs: failed to insert %s", name)
				result.Skipped++
			} else {
				result.New++
			}
		}

		// 递归扫描子目录
		if info.IsDir() {
			svc.scanDir(ctx, fullPath, result)
		}
	}

	// 如果当前目录没有显式的目录记录且不是根目录，创建一个目录记录
	if dirPath != svc.cfg.Nas.SimpleUploadRoot {
		dirName := strings.TrimPrefix(dirPath, svc.cfg.Nas.SimpleUploadRoot)
		dirName = strings.TrimPrefix(dirName, "/")
		if dirName != "" {
			_, dbErr := svc.store.RDB().GetFileByName(ctx, dirName)
			if dbErr != nil {
				// 目录没有 path 字段也没关系，后续能通过 name 查
				dirRecord := &types.File{
					Name:  dirName,
					Path:  dirPath,
					IsDir: true,
				}
				if err := svc.store.RDB().CreateFile(ctx, dirRecord); err != nil {
					log.WithError(err).Debugf("scanfs: dir record exists or failed for %s, skip", dirName)
				}
			}
		}
	}
	return nil
}