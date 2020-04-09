package models

import (
	"context"
	"hublabs/coupon-api/factory"
	"time"
)

type SendType string
type SendStatus string
type PeriodType string

const (
	SendTypeNow       SendType = "now"       // 立即发放
	SendTypeTimer     SendType = "timer"     // 定时发放
	SendTypeCycle     SendType = "cycle"     // 周期发放(现只实现了每日发放)
	SendTypePurchase  SendType = "purchase"  // 购买触发
	SendTypeInterface SendType = "interface" // 接口触发
	SendTypeNew       SendType = "new"       // 新会员注册触发
	SendTypeFree      SendType = "free"      // 代金券发放
	SendTypeBirth     SendType = "birth"     // 生日券发放
)

const (
	SendStatusPending  SendStatus = "pending"  //准备中
	SendStatusSending  SendStatus = "sending"  //发放中
	SendStatusFinished SendStatus = "finished" //发放完毕
)

const (
	PeriodTypeHandle     PeriodType = "handle"
	PeriodTypeDay        PeriodType = "day"        //发券当天
	PeriodTypeMonth      PeriodType = "month"      //发券当月
	PeriodTypeBirthDay   PeriodType = "birthday"   //生日当天
	PeriodTypeBirthMonth PeriodType = "birthMonth" //生日当月
)

//TODO:CustomerCondition结构
type PrepareCoupon struct {
	Id                 int64               `json:"id" xorm:"pk"`
	CampaignId         int64               `json:"campaignId"`
	OfferId            int64               `json:"offerId"`
	CustomerConditions []CustomerCondition `json:"customerConditions"`
	SeqNo              int                 `json:"seqNo"`    //相同seqNo表示同一发券的不同批次
	SaleType           string              `json:"saleType"` //线上，线下
	CouponInfo         CouponInfo          `json:"couponInfo" xorm:"-"`
	CouponInfoId       int64               `json:"couponInfoId"`
	CouponPeriod       CouponPeriod        `json:"couponPeriod" xorm:"json"`
	MaxPerQty          int64               `json:"maxPerQty"`
	MaxQty             int64               `json:"maxQty"`
	Percentage         float64             `json:"percentage"`
	SendType           SendType            `json:"sendType"`
	SendCondition      string              `json:"sendCondition" xorm:"json"`
	PurchaseConditions []PurchaseCondition `json:"purchaseConditions" xorm:"-"`
	SmsContent         string              `json:"smsContent"`
	Alert              Alert               `json:"alter" xorm:"json"`
	Enable             bool                `json:"enable"`
	SendStatus         SendStatus          `json:"sendStatus"`
	Commit             Commit              `json:"commit" xorm:"extends"`
}

type CouponPeriod struct {
	Type    PeriodType `json:"type"`
	Count   int        `json:"count"` //根据条件确定有效期后，往后延长几天，用于发券当天和发券当月
	StartAt time.Time  `json:"startAt"`
	EndAt   time.Time  `json:"endAt"`
}

type SendCondition struct {
	SendTime time.Timer `json:"sendTime,omitempty"` //定时发放的时间
	Period   int64      `json:"period,omitempty"`   //周期发放的周期
}

type Alert struct {
	SmsAlert    bool `json:"smsAlert"`
	WechatAlert bool `json:"wechatAlert"`
}

func (p *PrepareCoupon) Create(ctx context.Context) error {
	if p.CouponInfo.Id == 0 {
		if err := p.CouponInfo.Create(ctx); err != nil {
			return err
		}
	}
	p.CouponInfoId = p.CouponInfo.Id
	p.SendStatus = SendStatusPending
	if _, err := factory.DB(ctx).Insert(p); err != nil {
		return nil
	}

	seqCnt := make(map[int]int)
	for _, c := range p.CustomerConditions {
		if _, ok := seqCnt[c.SeqNo]; ok {
			seqCnt[c.SeqNo] = seqCnt[c.SeqNo] + 1
		} else {
			seqCnt[c.SeqNo] = 1
		}
	}
	for i := range p.CustomerConditions {
		p.CustomerConditions[i].PrepareCouponId = p.Id
		p.CustomerConditions[i].Cnt = seqCnt[p.CustomerConditions[i].SeqNo]
		if err := p.CustomerConditions[i].Create(ctx); err != nil {
			return err
		}
	}

	for j := range p.PurchaseConditions {
		p.PurchaseConditions[j].PrepareCouponId = p.Id
		if err := p.PurchaseConditions[j].Create(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (PrepareCoupon) Get(ctx context.Context, id int64) (*PrepareCoupon, error) {
	var pc PrepareCoupon
	if _, err := factory.DB(ctx).ID(id).Get(pc); err != nil {
		return nil, err
	}
	pcs, err := PrepareCoupon{}.getRelatedInfos(ctx, pc)
	if err != nil {
		return nil, err
	}
	return &pcs[0], nil
}

func (PrepareCoupon) getByCampaignIds(ctx context.Context, campaignIds ...int64) ([]PrepareCoupon, error) {
	var pcs []PrepareCoupon
	if err := factory.DB(ctx).In("campaign_id", campaignIds).Find(pcs); err != nil {
		return nil, err
	}

	pcs, err := PrepareCoupon{}.getRelatedInfos(ctx, pcs...)
	if err != nil {
		return nil, err
	}

	return pcs, nil
}

func (PrepareCoupon) getRelatedInfos(ctx context.Context, pcs ...PrepareCoupon) ([]PrepareCoupon, error) {
	var (
		ids           []int64
		couponInfoIds []int64
	)
	for _, pc := range pcs {
		ids = append(ids, pc.Id)
		couponInfoIds = append(couponInfoIds, pc.CouponInfoId)
	}

	ccs, err := CustomerCondition{}.getByPrepareCouponIds(ctx, ids...)
	if err != nil {
		return nil, err
	}

	ps, err := PurchaseCondition{}.getByPrepareCouponIds(ctx, ids...)
	if err != nil {
		return nil, err
	}

	cs, err := CouponInfo{}.getByIds(ctx, couponInfoIds...)
	if err != nil {
		return nil, err
	}

	for i := range pcs {
		for _, cc := range ccs {
			if pcs[i].Id == cc.PrepareCouponId {
				pcs[i].CustomerConditions = append(pcs[i].CustomerConditions, cc)
			}
		}

		for _, pc := range ps {
			if pcs[i].Id == pc.PrepareCouponId {
				pcs[i].PurchaseConditions = append(pcs[i].PurchaseConditions, pc)
			}
		}

		for _, c := range cs {
			if pcs[i].CouponInfoId == c.Id {
				pcs[i].CouponInfo = c
			}
		}
	}
	return pcs, nil
}
