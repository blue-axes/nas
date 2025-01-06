package postgres

import (
	"fmt"
	"github.com/blue-axes/tmpl/pkg/context"
	"github.com/blue-axes/tmpl/types"
	"strconv"
)

type (
	example struct {
		ID   int64  `gorm:"column:id;privateKey;autoIncrement"`
		Name string `gorm:"column:name;comment:the name"`
	}
)

func (*example) TableName() string {
	return "example"
}

func (m *example) ToExample() types.Example {
	return types.Example{
		ID:   fmt.Sprintf("%d", m.ID),
		Name: m.Name,
	}
}

func (m *example) FromExample(e types.Example) {
	m.ID, _ = strconv.ParseInt(e.ID, 10, 64)
	m.Name = e.Name
}

func (s *txStore) ListExamples(ctx *context.Context) (res []example, err error) {
	err = s.db.Find(&res).Error
	return res, err
}

func (s *txStore) CreateExamples(ctx *context.Context, e *types.Example) error {
	mdl := example{}
	mdl.FromExample(*e)
	err := s.db.Create(&mdl).Error
	if err != nil {
		return err
	}
	e.ID = fmt.Sprintf("%d", mdl.ID)
	return err
}
