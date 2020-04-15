package models

import (
	"context"
	"hublabs/coupon-api/factory"
	"time"
)

type RecoverRecord struct {
	Id        int64     `json:"id"`
	CouponNo  string    `json:"couponNo" xorm:"index"`
	UserId    int64     `json:"userId"`
	UseStore  string    `json:"useStore"`
	UseAt     time.Time `json:"useAt"`
	CreatedAt time.Time `json:"createdAt" xorm:"created"`
}

func (r RecoverRecord) Create(ctx context.Context) error {
	_, err := factory.DB(ctx).Insert(r)
	return err
}
