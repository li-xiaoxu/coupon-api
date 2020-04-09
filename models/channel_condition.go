package models

import (
	"context"
	"hublabs/coupon-api/factory"
)

type ChannelCondition struct {
	CampaignId int64  `json:"CampaignId"`
	Type       string `json:"type"` //brand_id|store_id
	Value      string `json:"value"`
}

func (c *ChannelCondition) Create(ctx context.Context) error {
	_, err := factory.DB(ctx).Insert(c)
	return err
}

func (ChannelCondition) getByCampaignIds(ctx context.Context, ids ...int64) ([]ChannelCondition, error) {
	var channelConditions []ChannelCondition
	err := factory.DB(ctx).In("campaign_id", ids).Find(&channelConditions)
	return channelConditions, err
}
