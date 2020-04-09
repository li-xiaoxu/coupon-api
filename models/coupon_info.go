package models

import (
	"context"
	"hublabs/coupon-api/factory"
)

type CouponInfo struct {
	Id         int64  `json:"id" xorm:"pk"`
	TenantCode string `json:"tenantCode"`
	Title      string `json:"title"`
	Notice     string `json:"notice"`
	Desc       string `json:"desc"`
	Headline   string `json:"headline"` //预留字段
}

func (c *CouponInfo) Create(ctx context.Context) error {
	_, err := factory.DB(ctx).Insert(c)
	return err
}

func (CouponInfo) Get(ctx context.Context, id int64) (couponInfo *CouponInfo, err error) {
	_, err = factory.DB(ctx).ID(id).Get(couponInfo)
	return
}

func (CouponInfo) getByIds(ctx context.Context, ids ...int64) ([]CouponInfo, error) {
	var couponInfos []CouponInfo
	if err := factory.DB(ctx).In("id", ids).Find(couponInfos); err != nil {
		return nil, err
	}
	return couponInfos, nil
}
