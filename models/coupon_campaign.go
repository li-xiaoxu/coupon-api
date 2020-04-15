package models

import (
	"context"
	"hublabs/coupon-api/factory"
	"time"
)

type CampaignStatus string

const (
	CampaignStatusPending CampaignStatus = "pending"
	CampaignStatusApprove CampaignStatus = "approve"
	CampaignStatusReject  CampaignStatus = "reject"
	CampaignStatusAbort   CampaignStatus = "abort"
)

type CouponCampaign struct {
	Id                int64              `json:"id"`
	TenantCode        string             `json:"tenantCode"`
	ChannelConditions []ChannelCondition `json:"channelConditions" xorm:"-"`
	PrepareCoupons    []PrepareCoupon    `json:"prepareCoupons" xorm:"-"`
	Name              string             `json:"name"`
	Desc              string             `json:"desc"`
	StartAt           time.Time          `json:"startAt"`
	EndAt             time.Time          `json:"endAt"`
	Status            CampaignStatus     `json:"status"`
	Commit            Commit             `json:"-" xorm:"extends"`
}

type Commit struct {
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt" xorm:"created"`
	UpdatedBy string    `json:"updatedBy"`
	UpdatedAt time.Time `json:"updatedAt" xorm:"updated"`
}

func (c *CouponCampaign) Create(ctx context.Context) error {
	c.TenantCode = getTenantCode(ctx)
	user := string(GetColleagueId(ctx))
	c.Commit.CreatedBy = user
	c.Commit.UpdatedBy = user
	if _, err := factory.DB(ctx).Insert(c); err != nil {
		return err
	}

	for i := range c.ChannelConditions {
		c.ChannelConditions[i].CampaignId = c.Id
		if err := c.ChannelConditions[i].Create(ctx); err != nil {
			return err
		}
	}

	for j := range c.PrepareCoupons {
		c.PrepareCoupons[j].CampaignId = c.Id
		if err := c.PrepareCoupons[j].Create(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (CouponCampaign) Get(ctx context.Context, id int64) (*CouponCampaign, error) {
	var c CouponCampaign
	if _, err := factory.DB(ctx).ID(id).Get(&c); err != nil {
		return nil, err
	}

	cs, err := ChannelCondition{}.getByCampaignIds(ctx, id)
	if err != nil {
		return nil, err
	}
	pcs, err := PrepareCoupon{}.getByCampaignIds(ctx, id)
	if err != nil {
		return nil, err
	}
	c.ChannelConditions = cs
	c.PrepareCoupons = pcs

	return &c, nil
}

func (CouponCampaign) GetAll(ctx context.Context, q, filter string, ids []int64, offset, limit int, sortby, order []string, status CampaignStatus, fields FieldTypeList) (totalCount int64, items []CouponCampaign, err error) {
	query := factory.DB(ctx).Where("tenant_code = ?", getTenantCode(ctx))
	if q != "" {
		query.And("name like ?", q+"%")
	}
	if len(ids) != 0 {
		query.In("id", ids)
	}
	if status != "" {
		query.And("status = ?", status)
	}

	if filter != "" {
		switch filter {
		//等待中
		case "wait":
			query.And("start_at > ?", time.Now())
		//现在或未来进行
		case "valid":
			query.And("end_at > ?", time.Now())
		//进行中
		case "available":
			query.And("? BETWEEN start_at AND end_at", time.Now())
		//已过期
		case "expired":
			query.And("final_at < ?", time.Now())
		}
	}
	if err = setSortOrder(query, sortby, order); err != nil {
		return
	}

	if limit == -1 {
		err = query.Find(&items)
		totalCount = int64(len(items))
	} else {
		totalCount, err = query.Limit(limit, offset).FindAndCount(&items)
	}
	if len(items) == 0 || len(fields) == 0 {
		return
	}
	ids = nil
	var (
		pcs []PrepareCoupon
		ccs []ChannelCondition
	)

	for _, item := range items {
		ids = append(ids, item.Id)
	}

	if fields.Contains(FieldTypePrepareCoupon) {
		pcs, err = PrepareCoupon{}.getByCampaignIds(ctx, ids...)
		if err != nil {
			return
		}
	}
	if fields.Contains(FieldTypeChannel) {
		ccs, err = ChannelCondition{}.getByCampaignIds(ctx, ids...)
		if err != nil {
			return
		}
	}

	for i := range items {
		for _, p := range pcs {
			if items[i].Id == p.CampaignId {
				items[i].PrepareCoupons = append(items[i].PrepareCoupons, p)
			}
		}
		for _, c := range ccs {
			if items[i].Id == c.CampaignId {
				items[i].ChannelConditions = append(items[i].ChannelConditions, c)
			}
		}
	}

	return
}

//未更新子表
func (c *CouponCampaign) Update(ctx context.Context, id int64) error {
	_, err := factory.DB(ctx).ID(id).Update(c)
	return err
}

func (CouponCampaign) GetByPrepareCouponId(ctx context.Context, id int64) (*CouponCampaign, error) {
	var c CouponCampaign
	_, err := factory.DB(ctx).Table("coupon_campaign").Select("coupon_campaign.*").
		Join("INNER", "prepare_coupon", "coupon_campaign.id = prepare_coupon.campaign_id").
		Where("prepare_coupon.id = ?", id).Get(&c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (CouponCampaign) Delete(ctx context.Context, id int64) error {
	_, err := factory.DB(ctx).ID(id).Delete(&CouponCampaign{})
	return err
}
