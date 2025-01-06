package service

import (
	"github.com/blue-axes/tmpl/pkg/context"
	"github.com/blue-axes/tmpl/types"
)

func (svc *Service) ListExample(ctx *context.Context) (res []types.Example, err error) {
	dbRes, err := svc.store.Postgres().ListExamples(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range dbRes {
		res = append(res, v.ToExample())
	}
	return res, nil
}

func (svc *Service) CreateExample(ctx *context.Context, e *types.Example) error {
	return svc.store.Postgres().CreateExamples(ctx, e)
}

func (svc *Service) MgListExample(ctx *context.Context) (res []types.Example, err error) {
	dbRes, err := svc.store.Mongo().ListExamples(ctx)
	if err != nil {
		return nil, err
	}
	for _, v := range dbRes {
		res = append(res, v.ToExample())
	}
	return res, nil
}

func (svc *Service) MgCreateExample(ctx *context.Context, e *types.Example) error {
	return svc.store.Mongo().CreateExamples(ctx, e)
}
