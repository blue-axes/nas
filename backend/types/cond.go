package types

import "github.com/blue-axes/tmpl/store/rdb/cond"

type (
	Condition = cond.Condition
	PageOrder = cond.PageOrder
)

var (
	ConditionNew = cond.ConditionNew
	ConditionAnd = cond.ConditionAnd
	ConditionOr  = cond.ConditionOr
	ConditionNot = cond.ConditionNot

	PageOrderNew = cond.PageOrderNew
)
