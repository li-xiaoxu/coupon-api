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

	"github.com/labstack/echo"
	"github.com/pangpanglabs/goutils/test"
)

func TestPrepareCouponCURD(t *testing.T) {
	couponPeriod := map[string]interface{}{
		"type":    models.PeriodTypeHandle,
		"count":   0,
		"startAt": time.Now().Add(time.Hour * -24),
		"endAt":   time.Now().Add(time.Hour * 24 * 30),
	}
	prepareCoupons := []map[string]interface{}{
		{
			"seqNo":        1,
			"offerId":      1,
			"saleType":     "POS",
			"couponInfoId": 1,
			"couponPeriod": couponPeriod,
			"maxPerQty":    1,
			"maxQty":       100,
			"percentage":   100,
			"sendType":     "new",
			"enable":       true,
			"sendStatus":   "pending",
		},
		{
			"seqNo":        2,
			"offerId":      2,
			"saleType":     "POS",
			"couponInfoId": 1,
			"couponPeriod": couponPeriod,
			"customerConditions": []map[string]interface{}{
				{
					"brand_id": 16,
				},
			},
			"maxPerQty":  1,
			"maxQty":     100,
			"percentage": 100,
			"sendType":   "internal",
			"enable":     true,
			"sendStatus": "pending",
		},
	}
	for i, p := range prepareCoupons {
		pb, _ := json.Marshal(p)
		t.Run(fmt.Sprint("Create#", i+1), func(t *testing.T) {
			req := httptest.NewRequest(echo.POST, "/v1/prepare-coupons", bytes.NewReader(pb))
			setHeader(req)
			rec := httptest.NewRecorder()
			test.Ok(t, handleWithFilter(PrepareCouponController{}.Create, echoApp.NewContext(req, rec)))
			test.Equals(t, http.StatusCreated, rec.Code)
		})
	}
	t.Run("Get", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/v1/prepare-coupons/1", nil)
		setHeader(req)
		rec := httptest.NewRecorder()
		c := echoApp.NewContext(req, rec)
		c.SetPath("/v1/prepare-coupons/:id")
		c.SetParamNames("id")
		c.SetParamValues("1")

		test.Ok(t, handleWithFilter(PrepareCouponController{}.Get, c))
		var v struct {
			Result  models.PrepareCoupon `json:"result"`
			Success bool                 `json:"success"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, v.Result.SendType, models.SendTypeInterface)
	})

	t.Run("SendCoupon", func(t *testing.T) {
		p := map[string]interface{}{
			"custId":  "100000000",
			"brandId": 16,
		}
		body, _ := json.Marshal(p)
		req := httptest.NewRequest(echo.POST, "/v1/prepare-coupons/1/send", bytes.NewReader(body))
		setHeader(req)
		rec := httptest.NewRecorder()
		c := echoApp.NewContext(req, rec)
		c.SetPath("/v1/prepare-coupons/:id/send")
		c.SetParamNames("id")
		c.SetParamValues("1")
		test.Ok(t, handleWithFilter(PrepareCouponController{}.SendCoupon, c))
		test.Equals(t, http.StatusOK, rec.Code)

		var v struct {
			Result  models.Coupon `json:"result"`
			Success bool          `json:"success"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, v.Result.PrepareCouponId, int64(1))
	})

	t.Run("SendNewCoupon", func(t *testing.T) {
		p := map[string]interface{}{
			"custId":   "100000000",
			"brandId":  16,
			"birthday": "1990-09-10",
		}
		body, _ := json.Marshal(p)
		req := httptest.NewRequest(echo.POST, "/v1/prepare-coupons/new", bytes.NewReader(body))
		setHeader(req)
		rec := httptest.NewRecorder()
		test.Ok(t, handleWithFilter(PrepareCouponController{}.SendNewCoupon, echoApp.NewContext(req, rec)))
		test.Equals(t, http.StatusOK, rec.Code)

		var v struct {
			Result  int64 `json:"result"`
			Success bool  `json:"success"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, v.Result, int64(1))
	})

	t.Run("SendBirthCoupon", func(t *testing.T) {
		c := map[string]interface{}{
			"custId":   "100000000",
			"brandId":  16,
			"birthday": time.Now().Format("2006-01-02"),
		}
		body, _ := json.Marshal(c)
		req := httptest.NewRequest(echo.POST, "/v1/prepare-coupons/birth", bytes.NewReader(body))
		setHeader(req)
		rec := httptest.NewRecorder()
		test.Ok(t, handleWithFilter(PrepareCouponController{}.SendBirthCoupon, echoApp.NewContext(req, rec)))
		test.Equals(t, http.StatusOK, rec.Code)

		var v struct {
			Result  int64 `json:"result"`
			Success bool  `json:"success"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, v.Result, int64(1))
	})
}
