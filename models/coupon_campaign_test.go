package models

import (
	"pangpangjan/kit/test"
	"testing"

	"github.com/renstrom/shortuuid"
)

func TestCouponCampaign(t *testing.T) {
	for _, campaign := range couponCampaigns {
		t.Run("Create CouponCampaign", func(t *testing.T) {
			err := campaign.Create(ctx)
			test.Ok(t, err)
			test.Equals(t, campaign.Name, "测试couponCamagin_1")
			test.Equals(t, campaign.Desc, "测试couponCamagin_1_描述")
			test.Equals(t, campaign.TenantCode, "pangpang")
			test.Equals(t, len(campaign.PrepareCoupons), 7)
			test.Equals(t, len(campaign.PrepareCoupons[0].PurchaseConditions), 0)
			test.Equals(t, len(campaign.PrepareCoupons[0].CustomerConditions), 0)
			test.Equals(t, campaign.PrepareCoupons[0].CouponInfo.Title, "测试couponCamagin_title")
			test.Equals(t, len(campaign.PrepareCoupons[1].PurchaseConditions), 1)
			test.Equals(t, campaign.PrepareCoupons[1].PurchaseConditions[0].PurchaseProducts[0].Type, TargetTypeBrand)
			test.Equals(t, len(campaign.PrepareCoupons[1].CustomerConditions), 3)
			test.Equals(t, campaign.PrepareCoupons[1].CustomerConditions[0].Type, "member_id")
			test.Equals(t, campaign.PrepareCoupons[1].CouponInfo.Title, "测试couponCamagin_title")
			for i := range campaign.PrepareCoupons {
				coupon := &Coupon{
					CouponNo:        shortuuid.New(),
					PrepareCouponId: prepareCoupons[i].Id,
					OfferId:         prepareCoupons[i].OfferId,
					CouponInfo:      prepareCoupons[i].CouponInfo,
					CustId:          "1",
					Status:          CouponStatusNormal,
				}
				err = coupon.Create(ctx)
				test.Ok(t, err)
			}
		})
	}

	t.Run("Get CouponCampaign By Id", func(t *testing.T) {
		campaign, err := CouponCampaign{}.Get(ctx, int64(1))
		test.Ok(t, err)
		test.Equals(t, campaign.Name, "测试couponCamagin_1")
		test.Equals(t, campaign.TenantCode, "pangpang")
		test.Equals(t, len(campaign.PrepareCoupons), 7)

		couponCampaign, err := CouponCampaign{}.Get(ctx, int64(1))
		test.Ok(t, err)
		test.Equals(t, couponCampaign.Name, "测试couponCamagin_1")
		test.Equals(t, campaign.Desc, "测试couponCamagin_1_描述")
		test.Equals(t, couponCampaign.TenantCode, "pangpang")
		test.Equals(t, len(couponCampaign.PrepareCoupons), 7)
		test.Equals(t, couponCampaign.PrepareCoupons[0].CouponInfo.Title, "测试couponCamagin_title")
		test.Equals(t, len(couponCampaign.PrepareCoupons[1].PurchaseConditions), 1)
		test.Equals(t, couponCampaign.PrepareCoupons[1].PurchaseConditions[0].PurchaseProducts[0].Type, TargetTypeBrand)
		test.Equals(t, len(couponCampaign.PrepareCoupons[1].CustomerConditions), 3)
		test.Equals(t, couponCampaign.PrepareCoupons[1].CustomerConditions[0].Type, "member_id")
		test.Equals(t, len(couponCampaign.PrepareCoupons[1].CustomerConditions[0].Targets), 2)
		test.Equals(t, len(couponCampaign.ChannelConditions), 2)
	})

	t.Run("Search CouponCampaigns", func(t *testing.T) {
		count, campaigns, err := CouponCampaign{}.GetAll(ctx, "", "", nil, 0, 10, []string{"id"}, []string{"asc"}, "", []FieldType{FieldTypePrepareCoupon})
		test.Ok(t, err)
		test.Equals(t, count, int64(1))
		test.Equals(t, len(campaigns), 1)
		test.Equals(t, campaigns[0].Name, "测试couponCamagin_1")
		test.Equals(t, campaigns[0].Desc, "测试couponCamagin_1_描述")
		test.Equals(t, campaigns[0].TenantCode, "pangpang")
		test.Equals(t, len(campaigns[0].PrepareCoupons), 7)
	})

	t.Run("Get PrepraCoupon by id And Update SendStatus", func(t *testing.T) {
		prepareCoupon, err := PrepareCoupon{}.Get(ctx, int64(1))
		test.Ok(t, err)
		test.Equals(t, prepareCoupon.CampaignId, int64(1))
		test.Equals(t, prepareCoupon.OfferId, int64(1))
		test.Equals(t, len(prepareCoupon.PurchaseConditions), 0)
		test.Equals(t, len(prepareCoupon.CustomerConditions), 0)

		prepareCoupon, err = PrepareCoupon{}.Get(ctx, int64(2))
		test.Ok(t, err)
		test.Equals(t, prepareCoupon.CampaignId, int64(1))
		test.Equals(t, prepareCoupon.OfferId, int64(2))
		test.Equals(t, len(prepareCoupon.PurchaseConditions), 1)
		test.Equals(t, len(prepareCoupon.CustomerConditions), 3)
		prepareCoupon.SendStatus = SendStatusFinished
		err = prepareCoupon.Update(ctx)
		test.Ok(t, err)
		test.Equals(t, prepareCoupon.SendStatus, SendStatusFinished)
	})

	t.Run("Search PrepareCoupons", func(t *testing.T) {
		count, prepareCoupons, err := PrepareCoupon{}.GetAll(ctx, "", "true", "100000001", "new", "", []string{"id"}, []string{"asc"}, 0, 10)
		test.Ok(t, err)
		test.Equals(t, count, int64(4))
		test.Equals(t, len(prepareCoupons), 4)
		test.Equals(t, prepareCoupons[0].CampaignId, int64(1))
		test.Equals(t, prepareCoupons[0].OfferId, int64(2))
		test.Equals(t, len(prepareCoupons[0].PurchaseConditions), 0)
		test.Equals(t, len(prepareCoupons[0].CustomerConditions), 0)
		test.Equals(t, prepareCoupons[1].CampaignId, int64(1))
		test.Equals(t, prepareCoupons[1].OfferId, int64(2))
		test.Equals(t, len(prepareCoupons[1].PurchaseConditions), 0)
		test.Equals(t, len(prepareCoupons[1].CustomerConditions), 0)
	})

	t.Run("Update CouponCampaign", func(t *testing.T) {
		campaign := couponCampaigns[0]
		campaign.Id = 1
		campaign.Name = "测试couponCamagin"
		campaign.Desc = "测试couponCamagin描述"
		campaign.Status = "abort"
		err := campaign.Update(ctx, campaign.Id)
		test.Ok(t, err)
		test.Equals(t, campaign.Name, "测试couponCamagin")
		test.Equals(t, campaign.Desc, "测试couponCamagin描述")
		test.Equals(t, campaign.Status, CampaignStatusAbort)
	})

	t.Run("Test filter customer condition", func(t *testing.T) {
		cust1 := &Member{
			Id:       "100000001",
			GradeId:  1,
			Birthday: "2019-06-01",
		}
		cust2 := &Member{
			Id:       "6",
			GradeId:  1,
			Birthday: "2019-06-01",
		}
		cust3 := &Member{
			Id:       "2",
			GradeId:  1,
			Birthday: "2019-02-01",
		}

		list1, err := CustomerCondition{}.FilterPrepareCoupon(ctx, cust1)
		test.Ok(t, err)
		test.Equals(t, len(list1), 7)
		list2, err := CustomerCondition{}.FilterPrepareCoupon(ctx, cust2)
		test.Ok(t, err)
		test.Equals(t, len(list2), 6)
		list3, err := CustomerCondition{}.FilterPrepareCoupon(ctx, cust3)
		test.Ok(t, err)
		test.Equals(t, list3, []int64{1, 2, 3, 4, 5, 6, 7})
		test.Equals(t, len(list3), 7)
	})

	t.Run("Test find prepareCoupon", func(t *testing.T) {
		count, prepareCoupons, err := PrepareCoupon{}.GetAll(ctx, "", "true", "100000001", "new", "", []string{"id"}, []string{"asc"}, 0, 10)
		test.Ok(t, err)
		test.Equals(t, count, int64(4))
		test.Equals(t, len(prepareCoupons), 4)
		var arr []float64
		for i := range prepareCoupons {
			arr = append(arr, prepareCoupons[i].Percentage)
		}
		index, err := FindPrepareCoupon(prepareCoupons)
		test.Ok(t, err)
		test.Equals(t, index <= 3, true)
		test.Equals(t, index >= 0, true)
		coupon := &Coupon{
			CouponNo:        shortuuid.New(),
			PrepareCouponId: prepareCoupons[index].Id,
			OfferId:         prepareCoupons[index].OfferId,
			CouponInfo:      prepareCoupons[index].CouponInfo,
			CustId:          "1",
			Status:          CouponStatusNormal,
		}
		err = coupon.Create(ctx)
		test.Ok(t, err)
	})

	t.Run("Test coupon search", func(t *testing.T) {
		count, coupons, err := Coupon{}.GetAll(ctx, "", "POS", "", "", []string{"prepare_coupon_id"}, []string{"asc"}, 0, 100)
		test.Ok(t, err)
		test.Equals(t, count, int64(8))
		test.Equals(t, len(coupons), 8)
	})

	t.Run("Delete CouponCampaign", func(t *testing.T) {
		err := CouponCampaign{}.Delete(ctx, int64(1))
		test.Ok(t, err)
		count, campaigns, err := CouponCampaign{}.GetAll(ctx, "", "", nil, 0, 10, []string{"id"}, []string{"asc"}, "", []FieldType{FieldTypePrepareCoupon})
		test.Ok(t, err)
		test.Equals(t, count, int64(0))
		test.Equals(t, len(campaigns), 0)
	})
}
