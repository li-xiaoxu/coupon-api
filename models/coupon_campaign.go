package models

import (
	"context"
	"hublabs/coupon-api/factory"
	"time"
)

type CampaignStatus string

const (
	CampaignStatusNormal CampaignStatus = "normal"
	CampaignStatusAbort  CampaignStatus = "abort"
)

type CouponCampaign struct {
	Id                int64              `json:"id" xorm:"pk"`
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
	user := string(getColleagueId(ctx))
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
