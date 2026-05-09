package rdb

import (
	stdErr "errors"
	"strings"

	"github.com/blue-axes/tmpl/pkg/constants"
	"github.com/blue-axes/tmpl/pkg/context"
	"github.com/blue-axes/tmpl/pkg/errors"
	"github.com/blue-axes/tmpl/types"
	"gorm.io/gorm"
	"path"
	"time"
)

type (
	file struct {
		ID        uint `gorm:"primarykey"`
		CreatedAt time.Time
		UpdatedAt time.Time
		Name      string `gorm:"column:name; size:512; unique; not null; comment:the name"`
		Ext       string `gorm:"column:ext; size:50; comment: the file type ext"`
		Path      string `gorm:"column:path; size:1024; comment: real path"`
		Size      uint64 `gorm:"column:size; default: 0; comment: file size"`
		Md5Sum    string `gorm:"column:md5_sum"`
		IsDir     bool   `gorm:"column:is_dir; default: false; comment: is directory"`
		Tags      string `gorm:"column:tags; type:text; comment: json array of tags"`
	}
)

func (*file) TableName() string {
	return "file_object"
}

func (m *file) ToEntity() types.File {
	tags := []string{}
	if m.Tags != "" {
		tags = strings.Split(strings.Trim(m.Tags, "[]"), ",")
		clean := make([]string, 0, len(tags))
		for _, t := range tags {
			t = strings.TrimSpace(t)
			t = strings.Trim(t, "\"")
			if t != "" {
				clean = append(clean, t)
			}
		}
		tags = clean
	}
	return types.File{
		ID:    m.ID,
		Name:  m.Name,
		Ext:   m.Ext,
		Path:  m.Path,
		Size:  m.Size,
		Md5:   m.Md5Sum,
		IsDir: m.IsDir,
		Tags:  tags,
	}
}

func (m *file) FromEntity(e types.File) {
	m.ID = e.ID
	m.Name = e.Name
	m.Ext = e.Ext
	if m.Ext == "" && !e.IsDir {
		m.Ext = path.Ext(m.Name)
	}
	m.Path = e.Path
	m.Size = e.Size
	m.Md5Sum = e.Md5
	m.IsDir = e.IsDir
	if len(e.Tags) > 0 {
		quoted := make([]string, len(e.Tags))
		for i, t := range e.Tags {
			quoted[i] = `"` + t + `"`
		}
		m.Tags = "[" + strings.Join(quoted, ",") + "]"
	} else {
		m.Tags = "[]"
	}
}

func (s *txStore) ListFile(ctx *context.Context, cond *types.Condition, page *types.PageOrder) (res []types.File, err error) {
	var (
		mdls = make([]file, 0)
	)
	res = make([]types.File, 0)
	q := cond.BuildCondition(s.db)
	if page != nil {
		q = page.BuildPageOrder(q)
	}
	err = q.Find(&mdls).Error
	if err != nil {
		return nil, err
	}
	for _, item := range mdls {
		res = append(res, item.ToEntity())
	}
	return res, err
}

func (s *txStore) SearchFiles(ctx *context.Context, keyword string, tag string) (res []types.File, err error) {
	var mdls []file
	q := s.db.Model(&file{})
	if keyword != "" {
		q = q.Where("name LIKE ?", "%"+keyword+"%")
	}
	if tag != "" {
		q = q.Where("tags LIKE ?", "%\""+tag+"\"%")
	}
	err = q.Find(&mdls).Error
	if err != nil {
		return nil, err
	}
	res = make([]types.File, 0, len(mdls))
	for _, item := range mdls {
		res = append(res, item.ToEntity())
	}
	return res, err
}

func (s *txStore) CreateFile(ctx *context.Context, e *types.File) error {
	mdl := file{}
	mdl.FromEntity(*e)
	err := s.db.Create(&mdl).Error
	if err != nil {
		return err
	}
	return err
}

func (s *txStore) UpdateFileByID(ctx *context.Context, id uint, e *types.File) error {
	mdl := file{}
	mdl.FromEntity(*e)
	mdl.ID = id
	err := s.db.Save(&mdl).Error
	if err != nil {
		return err
	}
	return err
}

func (s *txStore) UpdateFileTags(ctx *context.Context, name string, tags []string) error {
	var mdl file
	err := s.db.Where("name = ?", name).First(&mdl).Error
	if err != nil {
		if stdErr.Is(err, gorm.ErrRecordNotFound) {
			return errors.WithCode(constants.ErrCodeNotFound, name+" not found")
		}
		return err
	}
	e := mdl.ToEntity()
	e.Tags = tags
	mdl.FromEntity(e)
	return s.db.Save(&mdl).Error
}

func (s *txStore) CreateDir(ctx *context.Context, dirPath string) error {
	mdl := file{
		Name:  strings.TrimRight(dirPath, "/"),
		IsDir: true,
		Tags:  "[]",
	}
	return s.db.Create(&mdl).Error
}

func (s *txStore) GetFileByName(ctx *context.Context, name string) (res *types.File, err error) {
	mdl := file{}
	err = s.db.Where("name = ?", name).First(&mdl).Error
	if stdErr.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.WithCode(constants.ErrCodeNotFound, name+" not found")
	}
	f := mdl.ToEntity()
	return &f, nil
}

func (s *txStore) DeleteByName(ctx *context.Context, name string) error {
	mdl := file{}
	return s.db.Where("name = ?", name).Delete(&mdl).Error
}

func (s *txStore) DeleteDir(ctx *context.Context, dirPath string) error {
	dirPath = strings.TrimRight(dirPath, "/")
	return s.db.Where("name = ? OR name LIKE ?", dirPath, dirPath+"/%").Delete(&file{}).Error
}