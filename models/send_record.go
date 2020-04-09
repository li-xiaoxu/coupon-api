package models

import (
	"context"
	"hublabs/coupon-api/factory"
	"time"
)

//TODO:整理字段，争取每种类型的发券都能留有记录
type SendRecord struct {
	PrepareCouponId int64     `json:"prepareCouponId"`
	LastCustId      string    `json:"lastCustId"`
	ErrMsg          string    `json:"errMsg"`
	CreateAt        time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

func (s *SendRecord) Create(ctx context.Context) error {
	_, err := factory.DB(ctx).Insert(s)
	return err
}
