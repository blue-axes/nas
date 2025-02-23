package cond

import "gorm.io/gorm"

type (
	conditionType     string
	conditionRelation string
	Condition         struct {
		condType conditionType
		sql      string
		value    []interface{}

		relation         conditionRelation
		subConditionList []Condition
	}

	PageOrder struct {
		pageNumber uint
		pageSize   uint
		order      []string
	}
)

const (
	condData conditionType = "data"
	condRel  conditionType = "rel"

	relAnd conditionRelation = "and"
	relOr  conditionRelation = "or"
	relNot conditionRelation = "not"
)

func ConditionNew(sql string, args ...interface{}) Condition {
	return Condition{
		condType: condData,
		sql:      sql,
		value:    args,
	}
}

func ConditionAnd(list ...Condition) Condition {
	var (
		res = Condition{
			condType:         condRel,
			relation:         relAnd,
			subConditionList: make([]Condition, 0, len(list)),
		}
	)

	// and条件可以合并，没必要嵌套太深
	for _, item := range list {
		if item.condType == condRel && item.relation == relAnd {
			res.subConditionList = append(res.subConditionList, item.subConditionList...)
		} else {
			res.subConditionList = append(res.subConditionList, item)
		}
	}
	return res
}

func ConditionOr(list ...Condition) Condition {
	var (
		res = Condition{
			condType:         condRel,
			relation:         relOr,
			subConditionList: make([]Condition, 0, len(list)),
		}
	)

	// or条件可以合并，没必要嵌套太深
	for _, item := range list {
		if item.condType == condRel && item.relation == relOr {
			res.subConditionList = append(res.subConditionList, item.subConditionList...)
		} else {
			res.subConditionList = append(res.subConditionList, item)
		}
	}
	return res
}

func ConditionNot(cond Condition) Condition {
	var (
		res = Condition{
			condType:         condRel,
			relation:         relNot,
			subConditionList: []Condition{cond},
		}
	)

	return res
}

func (cond *Condition) BuildCondition(db *gorm.DB) *gorm.DB {
	if cond == nil {
		return db
	}

	if cond.condType == condData {
		// 叶子节点
		return db.Where(cond.sql, cond.value...)
	}
	// 关系节点
	return db.Where(cond.buildRelationCondition(*db, cond.relation, cond.subConditionList...))
}

func (self *Condition) buildRelationCondition(db gorm.DB, rel conditionRelation, cond ...Condition) *gorm.DB {
	var resDb = &db
	for _, item := range cond {
		switch {
		case item.condType == condData && rel == relAnd:
			resDb = resDb.Where(item.sql, item.value...)
		case item.condType == condData && rel == relOr:
			resDb = resDb.Or(item.sql, item.value...)
		case item.condType == condData && rel == relNot:
			resDb = resDb.Not(item.sql, item.value...)
		case item.condType == condRel && rel == relAnd:
			resDb = resDb.Where(self.buildRelationCondition(db, item.relation, item.subConditionList...))
		case item.condType == condRel && rel == relOr:
			resDb = resDb.Or(self.buildRelationCondition(db, item.relation, item.subConditionList...))
		case item.condType == condRel && rel == relNot:
			resDb = resDb.Not(self.buildRelationCondition(db, item.relation, item.subConditionList...))
		}
	}
	return resDb
}

func PageOrderNew(pageNumber, pageSize uint, order []string) *PageOrder {
	return &PageOrder{
		pageNumber: pageNumber,
		pageSize:   pageSize,
		order:      order,
	}
}

func (p *PageOrder) BuildPageOrder(db *gorm.DB) *gorm.DB {
	if p == nil {
		return db
	}
	for _, v := range p.order {
		db = db.Order(v)
	}
	if p.pageSize == 0 {
		p.pageSize = 10
	}
	if p.pageNumber == 0 {
		return db
	}
	return db.Offset(int((p.pageNumber - 1) * p.pageSize)).Limit(int(p.pageSize))
}
