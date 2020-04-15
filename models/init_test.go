package models

import (
	"context"
	"hublabs/coupon-api/config"
	"hublabs/coupon-api/factory"
	"os"
	"runtime"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/hublabs/common/auth"
	_ "github.com/mattn/go-sqlite3"

	"github.com/go-xorm/xorm"
	"github.com/pangpanglabs/goutils/behaviorlog"
	configutil "github.com/pangpanglabs/goutils/config"
	"github.com/pangpanglabs/goutils/echomiddleware"
	"github.com/pangpanglabs/goutils/jwtutil"
)

var ctx context.Context
var couponCampaigns []CouponCampaign
var prepareCoupons []PrepareCoupon
var coupons []Coupon

func init() {
	runtime.GOMAXPROCS(1)
	configutil.SetConfigPath("../")
	c := config.Init(os.Getenv("APP_ENV"))

	db, err := xorm.NewEngine(c.Database.Driver, c.Database.Connection)
	if err != nil {
		panic(err)
	}
	if err := DropTables(db); err != nil {
		panic(err)
	}
	if err := Init(db); err != nil {
		panic(err)
	}

	// db.ShowSQL(true)

	factory.Init(db)

	ctx = context.WithValue(context.Background(), echomiddleware.ContextDBName, db.NewSession())
	ctx = context.WithValue(ctx, "userClaim", auth.UserClaim{
		1,
		jwt.StandardClaims{
			Issuer: "pangpang",
		},
	})

	token, _ := jwtutil.NewToken(map[string]interface{}{"iss": "colleague", "tenantCode": "pangpang"})
	behaviorLogContext := behaviorlog.LogContext{AuthToken: token}
	ctx = behaviorLogContext.ToCtx(ctx)

	seed()
}

func seed() {
	couponInfo := CouponInfo{
		Id:     0,
		Title:  "测试couponCamagin_title",
		Desc:   "测试couponCamagin_desc",
		Notice: "测试couponCamagin_notice",
	}
	couponPeriod := CouponPeriod{
		Type:    PeriodTypeHandle,
		Count:   0,
		StartAt: time.Now().Add(time.Hour * -24),
		EndAt:   time.Now().Add(time.Hour * 24 * 30),
	}
	customerConditions := []CustomerCondition{
		CustomerCondition{
			SeqNo:    1,
			Comparer: "in",
			Type:     "member_id",
			Targets: []CustomerTarget{
				CustomerTarget{
					Value: "16",
				},
				CustomerTarget{
					Value: "100000001",
				},
			},
		},
		CustomerCondition{
			SeqNo:    2,
			Comparer: "in",
			Type:     "grade_id",
			Targets: []CustomerTarget{
				CustomerTarget{
					Value: "1",
				},
				CustomerTarget{
					Value: "2",
				},
			},
		},
		CustomerCondition{
			SeqNo:    2,
			Comparer: "in",
			Type:     "birthday_month",
			Targets: []CustomerTarget{
				CustomerTarget{
					Value: "01",
				},
				CustomerTarget{
					Value: "02",
				},
			},
		},
	}
	purchaseCondition := PurchaseCondition{
		PurchaseType: PurchaseTypeBoth,
		PriceAmt:     1000,
		QtyAmt:       2,
		PurchaseProducts: []PurchaseProduct{
			PurchaseProduct{
				Type:     "brand_id",
				Comparer: ComparerTypeIn,
				Targets: []PurchaseTarget{
					PurchaseTarget{
						Value: "16",
					},
				},
			},
		},
	}

	channelConditions := []ChannelCondition{
		ChannelCondition{
			Type:  "brand_id",
			Value: "16",
		},
		ChannelCondition{
			Type:  "store_id",
			Value: "2",
		},
	}

	prepareCoupons = append(prepareCoupons, PrepareCoupon{
		SeqNo:        1,
		OfferId:      1,
		SaleType:     "POS",
		CouponInfo:   couponInfo,
		CouponPeriod: couponPeriod,
		MaxPerQty:    1,
		MaxQty:       100,
		Percentage:   100,
		SendType:     "interface",
		Enable:       true,
		SendStatus:   "pending",
	})
	prepareCoupons = append(prepareCoupons, PrepareCoupon{
		SeqNo:              2,
		OfferId:            2,
		SaleType:           "POS",
		CouponInfo:         couponInfo,
		CouponPeriod:       couponPeriod,
		CustomerConditions: customerConditions,
		PurchaseConditions: []PurchaseCondition{purchaseCondition},
		MaxPerQty:          1,
		MaxQty:             100,
		Percentage:         100,
		SendType:           "interface",
		Enable:             true,
		SendStatus:         "pending",
	})
	prepareCoupons = append(prepareCoupons, PrepareCoupon{
		SeqNo:              3,
		OfferId:            2,
		SaleType:           "POS",
		CouponInfo:         couponInfo,
		CouponPeriod:       couponPeriod,
		CustomerConditions: nil,
		PurchaseConditions: nil,
		MaxPerQty:          2,
		MaxQty:             20,
		Percentage:         20,
		SendType:           "new",
		Enable:             true,
		SendStatus:         "pending",
	})
	prepareCoupons = append(prepareCoupons, PrepareCoupon{
		SeqNo:              3,
		OfferId:            2,
		SaleType:           "POS",
		CouponInfo:         couponInfo,
		CouponPeriod:       couponPeriod,
		CustomerConditions: nil,
		PurchaseConditions: nil,
		MaxPerQty:          2,
		MaxQty:             30,
		Percentage:         30,
		SendType:           "new",
		Enable:             true,
		SendStatus:         "pending",
	})
	prepareCoupons = append(prepareCoupons, PrepareCoupon{
		SeqNo:              3,
		OfferId:            2,
		SaleType:           "POS",
		CouponInfo:         couponInfo,
		CouponPeriod:       couponPeriod,
		CustomerConditions: nil,
		PurchaseConditions: nil,
		MaxPerQty:          2,
		MaxQty:             40,
		Percentage:         40,
		SendType:           "new",
		Enable:             true,
		SendStatus:         "pending",
	})
	prepareCoupons = append(prepareCoupons, PrepareCoupon{
		SeqNo:              3,
		OfferId:            2,
		SaleType:           "POS",
		CouponInfo:         couponInfo,
		CouponPeriod:       couponPeriod,
		CustomerConditions: nil,
		PurchaseConditions: nil,
		MaxPerQty:          2,
		MaxQty:             10,
		Percentage:         10,
		SendType:           "new",
		Enable:             true,
		SendStatus:         "pending",
	})
	prepareCoupons = append(prepareCoupons, PrepareCoupon{
		SeqNo:              4,
		OfferId:            1,
		SaleType:           "POS",
		CouponInfo:         couponInfo,
		CouponPeriod:       couponPeriod,
		CustomerConditions: nil,
		PurchaseConditions: nil,
		MaxPerQty:          1,
		MaxQty:             10,
		Percentage:         100,
		SendType:           "free",
		Enable:             true,
		SendStatus:         "pending",
	})

	couponCampaigns = append(couponCampaigns, CouponCampaign{
		Name:              "测试couponCamagin_1",
		Desc:              "测试couponCamagin_1_描述",
		PrepareCoupons:    prepareCoupons,
		ChannelConditions: channelConditions,
		StartAt:           time.Now().Add(time.Hour * -24),
		EndAt:             time.Now().Add(time.Hour * 24 * 30),
		Status:            CampaignStatusApprove,
	})
}
