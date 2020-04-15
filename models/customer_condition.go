package models

import (
	"context"
	"fmt"
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

type Member struct {
	Id         string `json:"id"`
	TenantCode string `json:"tenantCode"`
	MemberName string `json:"memberName"`
	Mobile     string `json:"mobile"`
	GradeId    int64  `json:"gradeId"`
	Birthday   string `json:"birthday"`
	MallId     int64  `json:"mallId"`
	UniqueId   string `json:"uniqueId"`
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
	if err := factory.DB(ctx).In("prepare_coupon_id", ids).Find(&conditions); err != nil {
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
	err := factory.DB(ctx).In("prepare_coupon_id", ids).Find(&targets)
	return targets, err
}

func (CustomerCondition) FilterPrepareCoupon(ctx context.Context, cust *Member) ([]int64, error) {
	var v []int64
	subQuery := `
	SELECT prepare_coupon_id, seq_no, cnt FROM customer_condition WHERE type = '%s' AND comparer = '%s' AND %s (
		SELECT 1 FROM customer_target WHERE prepare_coupon_id = customer_condition.prepare_coupon_id 
		AND condition_seq_no = customer_condition.seq_no
		AND condition_type = customer_condition.type
		AND value %s '%v'
	)`
	var query string
	var cndValue interface{}
	var index int
	for _, cnd := range customerConditionList {
		switch cnd.ConditionType {
		case ConditionTypeMemberId:
			cndValue = cust.Id
		case ConditionTypeBrandId:
			cndValue = cust.MallId
		case ConditionTypeGradeId:
			cndValue = cust.GradeId
		case ConditionTypeBirthdayMonth:
			if len(cust.Birthday) != 10 {
				continue
			}
			cndValue = cust.Birthday[5:7]
		}
		for _, cpr := range comparerList {
			if index != 0 {
				query += "\n        UNION ALL"
			}
			index++
			query += fmt.Sprintf(subQuery, cnd.ConditionType, cpr.Comparer, cpr.ExistsWord, cpr.Operator, cndValue)
		}
	}
	if err := factory.DBNewSession(ctx).SQL(fmt.Sprintf(`SELECT DISTINCT prepare_coupon_id FROM (%s
		) AS r
		GROUP BY prepare_coupon_id, seq_no, cnt
		HAVING COUNT(prepare_coupon_id) = cnt
		UNION
		SELECT id FROM prepare_coupon
		WHERE NOT EXISTS (
			SELECT 1 FROM customer_condition WHERE prepare_coupon.id = customer_condition.prepare_coupon_id
		) ORDER BY prepare_coupon_id ASC`, query)).Find(&v); err != nil {
		return nil, err
	}

	return v, nil
}
