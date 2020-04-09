package models

import (
	"context"
	"hublabs/coupon-api/factory"
	"time"
)

type CouponStatus string

const (
	CouponStatusNormal    CouponStatus = "normal"
	CouponStatusInvisible CouponStatus = "invisible"
	CouponStatusDisable   CouponStatus = "disable"
)

type Coupon struct {
	CouponNo        string       `json:"couponNo" xorm:"pk"`
	PrepareCouponId int64        `json:"prepareCouponId"`
	CustId          string       `json:"custId"` //能唯一标识顾客的字段
	Status          CouponStatus `json:"status"` //能见与否，能使用与否
	UseChk          bool         `json:"useChk"` //使用与否
	StartAt         time.Time    `json:"startAt"`
	EndAt           time.Time    `json:"endAt"`
	UseStore        string       `json:"useStore"` //能唯一标识卖场的字段
	Commit          Commit       `json:"commit" xorm:"extends"`
}

func (c *Coupon) Create(ctx context.Context) error {
	_, err := factory.DB(ctx).Insert(c)
	return err
}
