package models

import (
	"context"
	"hublabs/coupon-api/factory"
)

type ComparerType string

const (
	ComparerTypeIn               ComparerType = "in"
	ComparerTypeNin              ComparerType = "nin"
	ComparerTypeGreaterThanEqual ComparerType = "gte"
	ComparerTypeLessThanEqual    ComparerType = "lte"
)

type CustomerCondition struct {
	PrepareCouponId int64            `json:"-" xorm:"pk index"`
	SeqNo           int              `json:"seqNo" xorm:"pk"`
	Comparer        ComparerType     `json:"comparer" xorm:"varchar(4) index(cnd)"`  // in、nin
	Type            string           `json:"type" xorm:"pk varchar(16) index(cnd)" ` // member_id、grade_id、birthday_month
	Cnt             int              `json:"-" xorm:"index"`
	Targets         []CustomerTarget `json:"targets" xorm:"-"`
}

type CustomerTarget struct {
	PrepareCouponId int64  `json:"-" xorm:"pk"`
	ConditionSeqNo  int    `json:"-" xorm:"pk"`
	ConditionType   string `json:"-" xorm:"varchar(16) pk"`
	Value           string `json:"value" xorm:"varchar(64) pk"`
}

func (c *CustomerCondition) Create(ctx context.Context) error {
	if _, err := factory.DB(ctx).Insert(c); err != nil {
		return err
	}

	for i := range c.Targets {
		c.Targets[i].PrepareCouponId = c.PrepareCouponId
		c.Targets[i].ConditionSeqNo = c.SeqNo
		c.Targets[i].ConditionType = c.Type
		if err := c.Targets[i].create(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (t *CustomerTarget) create(ctx context.Context) error {
	_, err := factory.DB(ctx).Insert(t)
	return err
}

func (CustomerCondition) getByPrepareCouponIds(ctx context.Context, ids ...int64) ([]CustomerCondition, error) {
	var conditions []CustomerCondition
	if err := factory.DB(ctx).In("prepare_coupon_id", ids).Find(conditions); err != nil {
		return nil, err
	}
	targets, err := CustomerTarget{}.getByPrepareCouponIds(ctx, ids...)
	if err != nil {
		return nil, err
	}
	for i := range conditions {
		for _, target := range targets {
			if conditions[i].PrepareCouponId == target.PrepareCouponId && conditions[i].Type == target.ConditionType && conditions[i].SeqNo == target.ConditionSeqNo {
				conditions[i].Targets = append(conditions[i].Targets, target)
			}
		}
	}

	return conditions, nil
}

func (CustomerTarget) getByPrepareCouponIds(ctx context.Context, ids ...int64) ([]CustomerTarget, error) {
	var targets []CustomerTarget
	err := factory.DB(ctx).In("prepare_coupon_id", ids).Find(targets)
	return targets, err
}
