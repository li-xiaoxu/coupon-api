package models

import (
	"context"
	"hublabs/coupon-api/factory"
)

type CouponInfo struct {
	Id         int64  `json:"id"`
	TenantCode string `json:"tenantCode"`
	Title      string `json:"title"`
	Notice     string `json:"notice"`
	Desc       string `json:"desc"`
	Headline   string `json:"headline"` //预留字段
}

func (c *CouponInfo) Create(ctx context.Context) error {
	c.TenantCode = getTenantCode(ctx)
	_, err := factory.DB(ctx).Insert(c)
	return err
}

func (CouponInfo) Get(ctx context.Context, id int64) (couponInfo *CouponInfo, err error) {
	_, err = factory.DB(ctx).ID(id).Get(couponInfo)
	return
}

func (CouponInfo) getByIds(ctx context.Context, ids ...int64) ([]CouponInfo, error) {
	var couponInfos []CouponInfo
	if err := factory.DB(ctx).In("id", ids).Find(&couponInfos); err != nil {
		return nil, err
	}
	return couponInfos, nil
}

func (CouponInfo) GetAll(ctx context.Context, q string, sortby, order []string, skipCount, maxResultCount int) (items []CouponInfo, totalCount int64, err error) {
	query := factory.DB(ctx).Where("tenant_code = ?", getTenantCode(ctx))
	if err = setSortOrder(query, sortby, order); err != nil {
		return
	}
	if q != "" {
		query = query.Where("title like ?", q+"%")
	}

	totalCount, err = query.Limit(maxResultCount, skipCount).FindAndCount(&items)
	if len(items) == 0 {
		return
	}

	return
}
