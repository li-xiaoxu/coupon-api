package models

import (
	"sync"

	"github.com/go-xorm/xorm"
)

var (
	db   *xorm.Engine
	once sync.Once
)

func Init(e *xorm.Engine) error {
	err := e.Sync(
		new(CouponCampaign),
		new(ChannelCondition),
		new(PrepareCoupon),
		new(CustomerCondition),
		new(PurchaseCondition),
		new(PurchaseProduct),
		new(PurchaseTarget),
		new(Coupon),
		new(RecoverRecord),
		new(SendRecord),
	)
	if err != nil {
		return err
	}
	return nil
}

func DropTables(e *xorm.Engine) error {
	return e.DropTables(
		new(CouponCampaign),
		new(ChannelCondition),
		new(PrepareCoupon),
		new(CustomerCondition),
		new(PurchaseCondition),
		new(PurchaseProduct),
		new(PurchaseTarget),
		new(Coupon),
		new(RecoverRecord),
		new(SendRecord),
	)
}
