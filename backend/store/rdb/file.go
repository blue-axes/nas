package rdb

import (
	stdErr "errors"
	"github.com/blue-axes/tmpl/pkg/constants"
	"github.com/blue-axes/tmpl/pkg/context"
	"github.com/blue-axes/tmpl/pkg/errors"
	"github.com/blue-axes/tmpl/types"
	"gorm.io/gorm"
	"path"
)

type (
	file struct {
		gorm.Model
		Name   string `gorm:"column:name; size:512; unique; not null; comment:the name"`
		Ext    string `gorm:"column:ext; size:50; comment: the file type ext"`
		Path   string `gorm:"column:path; size:1024; comment: real path"`
		Size   uint64 `gorm:"column:size; default: 0; comment: file size"`
		Md5Sum string `gorm:"column:md5_sum"`
	}
)

func (*file) TableName() string {
	return "file_object"
}

func (m *file) ToEntity() types.File {
	return types.File{
		ID:   m.ID,
		Name: m.Name,
		Ext:  m.Ext,
		Path: m.Path,
		Size: m.Size,
		Md5:  m.Md5Sum,
	}
}

func (m *file) FromEntity(e types.File) {
	m.ID = e.ID
	m.Name = e.Name
	m.Ext = e.Ext
	if m.Ext == "" {
		m.Ext = path.Ext(m.Name)
	}
	m.Path = e.Path
	m.Size = e.Size
	m.Md5Sum = e.Md5
}

func (s *txStore) ListFile(ctx *context.Context, cond *types.Condition, page *types.PageOrder) (res []types.File, err error) {
	var (
		mdls = make([]file, 0)
	)
	res = make([]types.File, 0)
	err = page.BuildPageOrder(cond.BuildCondition(s.db)).Find(&mdls).Error
	if err != nil {
		return nil, err
	}
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
