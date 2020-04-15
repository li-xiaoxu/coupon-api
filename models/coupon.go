package models

import (
	"context"
	"fmt"
	"hublabs/coupon-api/factory"
	"strings"
	"time"
)

type CouponStatus string

const (
	CouponStatusNormal    CouponStatus = "normal"
	CouponStatusInvisible CouponStatus = "invisible"
	CouponStatusDisable   CouponStatus = "disable"
)

type Coupon struct {
	CouponNo         string             `json:"couponNo" xorm:"pk"`
	PrepareCouponId  int64              `json:"prepareCouponId"`
	CustId           string             `json:"custId"` //能唯一标识顾客的字段
	Status           CouponStatus       `json:"status"` //能见与否，能使用与否
	UseChk           bool               `json:"useChk"` //使用与否
	StartAt          time.Time          `json:"startAt"`
	EndAt            time.Time          `json:"endAt"`
	OfferId          int64              `json:"offerId" xorm:"-"`
	SaleType         string             `json:"saleType" xorm:"-"`
	CouponInfo       CouponInfo         `json:"couponInfo" xorm:"-"`
	ChannelConditons []ChannelCondition `json:"channelConditions" xorm:"-"`
	SendType         SendType           `json:"sendType" xorm:"-"`
	UseStore         string             `json:"useStore"` //能唯一标识卖场的字段
	Commit           Commit             `json:"commit" xorm:"extends"`
}

func (c *Coupon) Create(ctx context.Context) error {
	_, err := factory.DB(ctx).Insert(c)
	return err
}

func (Coupon) Get(ctx context.Context, no string, fields FieldTypeList) (*Coupon, error) {
	var ce struct {
		Coupon     `xorm:"extends"`
		CouponInfo `xorm:"extends"`
		OfferId    int64
		Id         int64
		SaleType   string
		SendType   SendType
	}
	has, err := factory.DB(ctx).Table("coupon").Select("coupon.*, coupon_info.*, prepare_coupon.offer_id, prepare_coupon.sale_type, prepare_coupon.send_type,prepare_coupon.campaign_id as id").
		Join("INNER", "prepare_coupon", "prepare_coupon.id = coupon.prepare_coupon_id").
		Join("INNER", "coupon_info", "coupon_info.id = prepare_coupon.coupon_info_id").
		Where("coupon.coupon_no = ?", no).
		Get(&ce)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, nil
	}
	if fields.Contains(FieldTypeChannel) {
		channelConditions, err := ChannelCondition{}.getByCampaignIds(ctx, ce.Id)
		if err != nil {
			return nil, err
		}
		ce.Coupon.ChannelConditons = channelConditions
	}
	ce.Coupon.CouponInfo = ce.CouponInfo
	ce.Coupon.OfferId = ce.OfferId
	ce.Coupon.SaleType = ce.SaleType
	ce.Coupon.SendType = ce.SendType
	return &ce.Coupon, nil
}

func (d *Coupon) Update(ctx context.Context, dto *Coupon) error {
	_, err := factory.DB(ctx).Where("coupon_no = ?", d.CouponNo).AllCols().Update(dto)
	return err
}

func (Coupon) Delete(ctx context.Context, no string) (err error) {
	_, err = factory.DB(ctx).ID(no).Delete(&Coupon{})
	return
}

func (Coupon) GetAll(ctx context.Context, custId string, saleType, filter, nos string, sortby, order []string, offset, limit int) (totalCount int64, items []Coupon, err error) {
	var ces []struct {
		Coupon     `xorm:"extends"`
		CouponInfo `xorm:"extends"`
		Id         int64
		OfferId    int64
		SaleType   string
		SendType   SendType
	}
	query := factory.DB(ctx).Table("coupon").Select("coupon.*, coupon_info.*, prepare_coupon.offer_id, prepare_coupon.sale_type, prepare_coupon.send_type,prepare_coupon.campaign_id as id").
		Join("INNER", "prepare_coupon", "prepare_coupon.id = coupon.prepare_coupon_id").
		Join("INNER", "coupon_info", "coupon_info.id = prepare_coupon.coupon_info_id AND coupon_info.tenant_code = ?", getTenantCode(ctx)).
		Where("coupon.status = ?", CouponStatusNormal)

	if saleType != "" {
		query.Where("prepare_coupon.sale_type = ? OR prepare_coupon.sale_type = ''", saleType)
	}

	switch filter {
	case "unused":
		query.Where("coupon.use_chk = ?", false) // 没使用
	case "used":
		query.Where("coupon.use_chk = ?", true) // 使用了
	case "expired":
		query.Where("coupon.use_chk = ? AND coupon.end_at < ?", false, time.Now()) // 过期了
	case "valid":
		query.Where("coupon.use_chk = ? AND coupon.end_at >= ?", false, time.Now()) // 现在或未来可使用
	case "invalid":
		query.Where("coupon.use_chk = ? OR coupon.end_at < ?", true, time.Now()) // 使用了或过期了
	case "available":
		query.Where("coupon.use_chk = ? AND (? BETWEEN coupon.start_at AND coupon.end_at)", false, time.Now()) // 现在可使用
	case "unavailable":
		query.Where("coupon.use_chk = ? OR (? NOT BETWEEN coupon.start_at AND coupon.end_at)", true, time.Now()) // 现在不可用
	}

	var conditionQueries []string
	if custId != "" {
		conditionQueries = append(conditionQueries, fmt.Sprintf("coupon.cust_id = %s", custId))
	}
	if nos != "" {
		nos := strings.Split(nos, ",")
		conditionQueries = append(conditionQueries, fmt.Sprintf("coupon.coupon_no IN ('%s')", strings.Join(nos, "','")))
	}
	if len(conditionQueries) != 0 {
		query.Where(strings.Join(conditionQueries, " OR "))
	}

	if err = setSortOrder(query, sortby, order, "coupon"); err != nil {
		return
	}

	if limit == -1 {
		err = query.Find(&ces)
		totalCount = int64(len(ces))
	} else {
		totalCount, err = query.Limit(limit, offset).FindAndCount(&ces)
	}

	for _, cs := range ces {
		cs.Coupon.OfferId = cs.OfferId
		cs.Coupon.SaleType = cs.SaleType
		cs.Coupon.CouponInfo = cs.CouponInfo
		channelConditions, err := ChannelCondition{}.getByCampaignIds(ctx, cs.Id)
		if err != nil {
			return 0, nil, err
		}
		cs.Coupon.SendType = cs.SendType
		cs.Coupon.ChannelConditons = channelConditions
		items = append(items, cs.Coupon)
	}
	return
}

func (Coupon) GetByPrepareCouponIds(ctx context.Context, ids ...int64) ([]Coupon, error) {
	var coupons []Coupon
	if err := factory.DB(ctx).
		In("prepare_coupon_id", ids).Find(&coupons); err != nil {
		return nil, err
	}
	return coupons, nil
}

func (Coupon) CreateInArray(ctx context.Context, coupons []Coupon) error {
	if _, err := factory.DB(ctx).Insert(&coupons); err != nil {
		return err
	}
	return nil
}
