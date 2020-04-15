package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hublabs/coupon-api/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pangpanglabs/goutils/test"

	"github.com/labstack/echo"
)

func TestCampaignCRUD(t *testing.T) {
	couponInfo := map[string]interface{}{
		"id":         0,
		"title":      "测试couponCampagin_title",
		"desc":       "测试couponCampagin_desc",
		"notice":     "测试couponCampagin_notice",
		"allowExtra": true,
	}
	couponPeriod := map[string]interface{}{
		"type":    models.PeriodTypeHandle,
		"count":   0,
		"startAt": time.Now().Add(time.Hour * -24),
		"endAt":   time.Now().Add(time.Hour * 24 * 30),
	}
	inputs := []map[string]interface{}{
		{
			"name": "测试接口触发",
			"desc": "测试接口触发_描述",
			"prepareCoupons": []map[string]interface{}{
				{
					"seqNo":        1,
					"offerId":      1,
					"saleType":     "POS",
					"couponInfo":   couponInfo,
					"couponPeriod": couponPeriod,
					"maxPerQty":    1,
					"maxQty":       100,
					"percentage":   100,
					"sendType":     "interface",
					"enable":       true,
					"sendStatus":   "pending",
				},
				{
					"seqNo":        2,
					"offerId":      2,
					"saleType":     "POS",
					"couponInfo":   couponInfo,
					"couponPeriod": couponPeriod,
					"customerConditions": []map[string]interface{}{
						{
							"seqNo":    1,
							"comparer": "in",
							"type":     "brand_id",
							"targets": []map[string]interface{}{
								{
									"value": "16",
								},
							},
						},
					},
					"maxPerQty":  1,
					"maxQty":     100,
					"percentage": 100,
					"sendType":   "interface",
					"enable":     true,
					"sendStatus": "pending",
				},
			},
			"channelConditions": []map[string]interface{}{
				{
					"type":  "brand_id",
					"value": "16",
				},
			},
			"startAt": time.Now().Add(time.Hour * -24),
			"endAt":   time.Now().Add(time.Hour * 24 * 30),
			"status":  models.CampaignStatusPending,
		},
		{
			"name": "测试新会员发放",
			"desc": "测试新会员发放_描述",
			"prepareCoupons": []map[string]interface{}{
				{
					"seqNo":        3,
					"offerId":      2,
					"saleType":     "POS",
					"couponInfo":   couponInfo,
					"couponPeriod": couponPeriod,
					// "customerConditions": []map[string]interface{}{
					// 	{
					// 		"seqNo":    1,
					// 		"comparer": "in",
					// 		"type":     "brand_id",
					// 		"targets": []map[string]interface{}{
					// 			{
					// 				"value": "16",
					// 			},
					// 		},
					// 	},
					// },
					"maxPerQty":  2,
					"maxQty":     20,
					"percentage": 20,
					"sendType":   "new",
					"enable":     true,
					"sendStatus": "pending",
				},
				{
					"seqNo":        3,
					"offerId":      2,
					"saleType":     "POS",
					"couponInfo":   couponInfo,
					"couponPeriod": couponPeriod,
					// "customerConditions": []map[string]interface{}{
					// 	{
					// 		"seqNo":    1,
					// 		"comparer": "in",
					// 		"type":     "brand_id",
					// 		"targets": []map[string]interface{}{
					// 			{
					// 				"value": "16",
					// 			},
					// 		},
					// 	},
					// },
					"maxPerQty":  2,
					"maxQty":     30,
					"percentage": 30,
					"sendType":   "new",
					"enable":     true,
					"sendStatus": "pending",
				},
				{
					"seqNo":        3,
					"offerId":      2,
					"saleType":     "POS",
					"couponInfo":   couponInfo,
					"couponPeriod": couponPeriod,
					// "customerConditions": []map[string]interface{}{
					// 	{
					// 		"seqNo":    1,
					// 		"comparer": "in",
					// 		"type":     "brand_id",
					// 		"targets": []map[string]interface{}{
					// 			{
					// 				"value": "16",
					// 			},
					// 		},
					// 	},
					// },
					"maxPerQty":  2,
					"maxQty":     40,
					"percentage": 40,
					"sendType":   "new",
					"enable":     true,
					"sendStatus": "pending",
				},
				{
					"seqNo":        3,
					"offerId":      2,
					"saleType":     "POS",
					"couponInfo":   couponInfo,
					"couponPeriod": couponPeriod,
					"customerConditions": []map[string]interface{}{
						{
							"seqNo":    1,
							"comparer": "in",
							"type":     "brand_id",
							"targets": []map[string]interface{}{
								{
									"value": "16",
								},
							},
						},
					},
					"maxPerQty":  2,
					"maxQty":     10,
					"percentage": 10,
					"sendType":   "new",
					"enable":     true,
					"sendStatus": "pending",
				},
			},
			"channelConditions": []map[string]interface{}{
				{
					"type":  "brand_id",
					"value": "16",
				},
			},
			"startAt": time.Now().Add(time.Hour * -24),
			"endAt":   time.Now().Add(time.Hour * 24 * 30),
			"status":  models.CampaignStatusPending,
		},
		{
			"name": "测试代金券发放",
			"desc": "测试代金券发放_描述",
			"prepareCoupons": []map[string]interface{}{
				{
					"seqNo":        4,
					"offerId":      3,
					"saleType":     "POS",
					"couponInfo":   couponInfo,
					"couponPeriod": couponPeriod,
					"customerConditions": []map[string]interface{}{
						{
							"seqNo":    1,
							"comparer": "in",
							"type":     "brand_id",
							"targets": []map[string]interface{}{
								{
									"value": "16",
								},
							},
						},
					},
					"maxPerQty":  1,
					"maxQty":     10,
					"percentage": 100,
					"sendType":   "free",
					"enable":     true,
					"sendStatus": "pending",
				},
			},
			"channelConditions": []map[string]interface{}{
				{
					"type":  "brand_id",
					"value": "16",
				},
			},
			"startAt": time.Now().Add(time.Hour * -24),
			"endAt":   time.Now().Add(time.Hour * 24 * 30),
			"status":  models.CampaignStatusPending,
		},
		{
			"name": "测试新会员发放",
			"desc": "测试新会员发放_描述",
			"prepareCoupons": []map[string]interface{}{
				{
					"seqNo":        1,
					"offerId":      3,
					"saleType":     "POS",
					"couponInfo":   couponInfo,
					"couponPeriod": couponPeriod,
					// "customerConditions": []map[string]interface{}{
					// 	{
					// 		"seqNo":    1,
					// 		"comparer": "in",
					// 		"type":     "birth",
					// 		"targets": []map[string]interface{}{
					// 			{
					// 				"value": time.Now().Format("2006-01-02"),
					// 			},
					// 		},
					// 	},
					// },
					"maxPerQty":  1,
					"maxQty":     0,
					"percentage": 100,
					"sendType":   "birth",
					"enable":     true,
					"sendStatus": "pending",
				},
			},
			"channelConditions": []map[string]interface{}{
				{
					"type":  "brand_id",
					"value": "16",
				},
			},
			"startAt": time.Now().Add(time.Hour * -24),
			"endAt":   time.Now().Add(time.Hour * 24 * 30),
			"status":  models.CampaignStatusPending,
		},
	}
	for i, p := range inputs {
		pb, _ := json.Marshal(p)
		t.Run(fmt.Sprint("Create#", i+1), func(t *testing.T) {
			req := httptest.NewRequest(echo.POST, "/v1/coupon-campaigns", bytes.NewReader(pb))
			setHeader(req)
			rec := httptest.NewRecorder()
			test.Ok(t, handleWithFilter(CouponCampaignController{}.Create, echoApp.NewContext(req, rec)))
			test.Equals(t, http.StatusCreated, rec.Code)
		})
	}
	t.Run("GetAll", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/v1/coupon-campaigns", nil)
		setHeader(req)
		rec := httptest.NewRecorder()
		test.Ok(t, handleWithFilter(CouponCampaignController{}.GetAll, echoApp.NewContext(req, rec)))
		test.Equals(t, http.StatusOK, rec.Code)

		var v struct {
			Result struct {
				TotalCount int                     `json:"totalCount"`
				Items      []models.CouponCampaign `json:"items"`
			} `json:"result"`
			Success bool `json:"success"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, v.Result.TotalCount, 4)
		test.Equals(t, v.Result.Items[0].Name, "测试接口触发")
		test.Equals(t, v.Result.Items[1].Name, "测试新会员发放")
		test.Equals(t, v.Result.Items[2].Name, "测试代金券发放")
	})

	t.Run("Get", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/v1/coupon-campaigns/1", nil)
		setHeader(req)
		rec := httptest.NewRecorder()
		c := echoApp.NewContext(req, rec)
		c.SetPath("/v1/coupon-campaigns/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		test.Ok(t, handleWithFilter(CouponCampaignController{}.GetOne, c))
		test.Equals(t, http.StatusOK, rec.Code)
		var v struct {
			Result  models.CouponCampaign `json:"result"`
			Success bool                  `json:"success"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, v.Result.Name, "测试接口触发")
	})

	t.Run("Update", func(t *testing.T) {
		//审批接口触发发券
		campaigns := map[string]interface{}{
			"status": models.CampaignStatusApprove,
		}
		cc, _ := json.Marshal(campaigns)
		req := httptest.NewRequest(echo.PUT, "/v1/coupon-campaigns/1", bytes.NewReader(cc))
		setHeader(req)
		rec := httptest.NewRecorder()

		c := echoApp.NewContext(req, rec)
		c.SetPath("/v1/coupon-campaigns/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		test.Ok(t, handleWithFilter(CouponCampaignController{}.Update, c))
		test.Equals(t, http.StatusNoContent, rec.Code)

		//审批新会员注册发券
		req = httptest.NewRequest(echo.PUT, "/v1/coupon-campaigns/2", bytes.NewReader(cc))
		setHeader(req)
		rec = httptest.NewRecorder()

		c = echoApp.NewContext(req, rec)
		c.SetPath("/v1/coupon-campaigns/:id")
		c.SetParamNames("id")
		c.SetParamValues("2")

		test.Ok(t, handleWithFilter(CouponCampaignController{}.Update, c))
		test.Equals(t, http.StatusNoContent, rec.Code)

		//审批生日券
		req = httptest.NewRequest(echo.PUT, "/v1/coupon-campaigns/4", bytes.NewReader(cc))
		setHeader(req)
		rec = httptest.NewRecorder()

		c = echoApp.NewContext(req, rec)
		c.SetPath("/v1/coupon-campaigns/:id")
		c.SetParamNames("id")
		c.SetParamValues("4")

		test.Ok(t, handleWithFilter(CouponCampaignController{}.Update, c))
		test.Equals(t, http.StatusNoContent, rec.Code)
	})
}
