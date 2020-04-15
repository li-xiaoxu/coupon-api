package models

import (
	"context"
	"hublabs/coupon-api/factory"
)

type PurchaseType string
type TargetType string

const (
	PurchaseTypePrice  PurchaseType = "price"
	PurchaseTypeQty    PurchaseType = "qty"
	PurchaseTypeBoth   PurchaseType = "both"
	PurchaseTypeEither PurchaseType = "either"
)

const (
	TargetTypeBrand          TargetType = "brand_id"
	TargetTypeProduct        TargetType = "product_id"
	TargetTypeSku            TargetType = "sku_id"
	TargetTypeSkuIndentifier TargetType = "sku_uid"
)

type PurchaseCondition struct {
	PrepareCouponId  int64             `json:"prepareCouponId" xorm:"index"`
	PurchaseType     PurchaseType      `json:"purchaseType"`
	PriceAmt         float64           `json:"priceAmt"`
	QtyAmt           int64             `json:"qtyAmt"`
	PurchaseProducts []PurchaseProduct `json:"productConditions" xorm:"-"`
}

type PurchaseProduct struct {
	PrepareCouponId int64            `json:"prepareCouponId" xorm:"pk index"`
	Type            TargetType       `json:"type" xorm:"pk varchar(16) index(cnd)"`
	Comparer        ComparerType     `json:"comparer" xorm:"varchar(4) index(cnd)"`
	Targets         []PurchaseTarget `json:"targets" xorm:"-"`
}

type PurchaseTarget struct {
	PrepareCouponId     int64      `json:"prepareCouponId" xorm:"pk"`
	ConditionTargetType TargetType `json:"conditionTargetType" xorm:"varchar(16) pk"`
	Value               string     `json:"value" xorm:"varhcar(16) pk"`
}

func (c *PurchaseCondition) Create(ctx context.Context) error {
	if _, err := factory.DB(ctx).Insert(c); err != nil {
		return err
	}
	for i := range c.PurchaseProducts {
		c.PurchaseProducts[i].PrepareCouponId = c.PrepareCouponId
		if err := c.PurchaseProducts[i].create(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (PurchaseCondition) getByPrepareCouponIds(ctx context.Context, ids ...int64) ([]PurchaseCondition, error) {
	var conditions []PurchaseCondition
	if err := factory.DB(ctx).In("prepare_coupon_id", ids).Find(&conditions); err != nil {
		return nil, err
	}

	products, err := PurchaseProduct{}.getByPrepareCouponIds(ctx, ids...)
	if err != nil {
		return nil, err
	}

	for i := range conditions {
		for _, product := range products {
			if conditions[i].PrepareCouponId == product.PrepareCouponId {
				conditions[i].PurchaseProducts = append(conditions[i].PurchaseProducts, product)
			}
		}
	}
	return conditions, nil
}

func (p *PurchaseProduct) create(ctx context.Context) error {
	if _, err := factory.DB(ctx).Insert(p); err != nil {
		return err
	}
	for i := range p.Targets {
		p.Targets[i].PrepareCouponId = p.PrepareCouponId
		p.Targets[i].ConditionTargetType = p.Type
		if err := p.Targets[i].create(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (PurchaseProduct) getByPrepareCouponIds(ctx context.Context, ids ...int64) ([]PurchaseProduct, error) {
	var products []PurchaseProduct
	if err := factory.DB(ctx).In("prepare_coupon_id", ids).Find(&products); err != nil {
		return nil, err
	}

	targets, err := PurchaseTarget{}.getByPrepareCouponIds(ctx, ids...)
	if err != nil {
		return nil, err
	}

	for i := range products {
		for _, target := range targets {
			if products[i].PrepareCouponId == target.PrepareCouponId && products[i].Type == target.ConditionTargetType {
				products[i].Targets = append(products[i].Targets, target)
			}
		}
	}
	return products, nil
}

func (t *PurchaseTarget) create(ctx context.Context) error {
	_, err := factory.DB(ctx).Insert(t)
	return err
}

func (PurchaseTarget) getByPrepareCouponIds(ctx context.Context, ids ...int64) ([]PurchaseTarget, error) {
	var targets []PurchaseTarget
	if err := factory.DB(ctx).In("prepare_coupon_id", ids).Find(&targets); err != nil {
		return nil, err
	}
	return targets, nil
}
