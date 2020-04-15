package controllers

import (
	"bytes"
	"encoding/json"
	"hublabs/coupon-api/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/pangpanglabs/goutils/test"
)

func TestCoupon(t *testing.T) {
	//审批代金券，进行发券
	campaigns := map[string]interface{}{
		"status": models.CampaignStatusApprove,
	}
	cc, _ := json.Marshal(campaigns)
	req := httptest.NewRequest(echo.PUT, "/v1/coupon-campaigns/3", bytes.NewReader(cc))
	setHeader(req)
	rec := httptest.NewRecorder()

	c := echoApp.NewContext(req, rec)
	c.SetPath("/v1/coupon-campaigns/:id")
	c.SetParamNames("id")
	c.SetParamValues("3")

	test.Ok(t, handleWithFilter(CouponCampaignController{}.Update, c))
	test.Equals(t, http.StatusNoContent, rec.Code)
	//查询券号
	req = httptest.NewRequest(echo.GET, "/v1/coupons?sendType=free", nil)
	setHeader(req)
	rec = httptest.NewRecorder()
	test.Ok(t, handleWithFilter(CouponController{}.GetAll, echoApp.NewContext(req, rec)))
	test.Equals(t, http.StatusOK, rec.Code)

	var v struct {
		Result struct {
			TotalCount int             `json:"totalCount"`
			Items      []models.Coupon `json:"items"`
		} `json:"result"`
		Success bool `json:"success"`
	}
	test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
	test.Equals(t, v.Result.TotalCount, 10)

	t.Run("Get", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/v1/coupons/"+v.Result.Items[0].CouponNo, nil)
		setHeader(req)
		rec := httptest.NewRecorder()
		c := echoApp.NewContext(req, rec)
		c.SetPath("/v1/coupons/:no")
		c.SetParamNames("no")
		c.SetParamValues(v.Result.Items[0].CouponNo)

		test.Ok(t, handleWithFilter(CouponController{}.GetOne, c))
		var result struct {
			Result  models.Coupon `json:"result"`
			Success bool          `json:"success"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &result))
		test.Equals(t, result.Result.CouponNo, v.Result.Items[0].CouponNo)
		test.Equals(t, result.Result.UseChk, false)
	})
	t.Run("UseCoupon", func(t *testing.T) {
		req := httptest.NewRequest(echo.POST, "/v1/coupons/"+v.Result.Items[1].CouponNo+"/use", nil)
		setHeader(req)
		rec := httptest.NewRecorder()
		c := echoApp.NewContext(req, rec)
		c.SetPath("/v1/coupons/:no/use")
		c.SetParamNames("no")
		c.SetParamValues(v.Result.Items[1].CouponNo)

		test.Ok(t, handleWithFilter(CouponController{}.Use, c))

		req = httptest.NewRequest(echo.GET, "/v1/coupons/"+v.Result.Items[1].CouponNo, nil)
		setHeader(req)
		rec = httptest.NewRecorder()
		c = echoApp.NewContext(req, rec)
		c.SetPath("/v1/coupons/:no")
		c.SetParamNames("no")
		c.SetParamValues(v.Result.Items[1].CouponNo)

		test.Ok(t, handleWithFilter(CouponController{}.GetOne, c))
		var result struct {
			Result  models.Coupon `json:"result"`
			Success bool          `json:"success"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &result))
		test.Equals(t, result.Result.CouponNo, v.Result.Items[1].CouponNo)
		test.Equals(t, result.Result.UseChk, true)

		req = httptest.NewRequest(echo.POST, "/v1/coupons/"+v.Result.Items[2].CouponNo+"/recover", nil)
		setHeader(req)
		rec = httptest.NewRecorder()
		c = echoApp.NewContext(req, rec)
		c.SetPath("/v1/coupons/:no/recover")
		c.SetParamNames("no")
		c.SetParamValues(v.Result.Items[1].CouponNo)

		test.Ok(t, handleWithFilter(CouponController{}.Recover, c))

		req = httptest.NewRequest(echo.GET, "/v1/coupons/"+v.Result.Items[1].CouponNo, nil)
		setHeader(req)
		rec = httptest.NewRecorder()
		c = echoApp.NewContext(req, rec)
		c.SetPath("/v1/coupons/:no")
		c.SetParamNames("no")
		c.SetParamValues(v.Result.Items[1].CouponNo)

		test.Ok(t, handleWithFilter(CouponController{}.GetOne, c))
		var coupon struct {
			Result  models.Coupon `json:"result"`
			Success bool          `json:"success"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &coupon))
		test.Equals(t, coupon.Result.CouponNo, v.Result.Items[1].CouponNo)
		test.Equals(t, coupon.Result.UseChk, false)
	})

	t.Run("Delete", func(t *testing.T) {
		req := httptest.NewRequest(echo.DELETE, "/v1/coupons/"+v.Result.Items[1].CouponNo, nil)
		setHeader(req)
		rec := httptest.NewRecorder()
		c := echoApp.NewContext(req, rec)
		c.SetPath("/v1/coupons/:no")
		c.SetParamNames("no")
		c.SetParamValues(v.Result.Items[2].CouponNo)

		test.Ok(t, handleWithFilter(CouponController{}.Delete, c))

		req = httptest.NewRequest(echo.GET, "/v1/coupons/"+v.Result.Items[2].CouponNo, nil)
		setHeader(req)
		rec = httptest.NewRecorder()
		c = echoApp.NewContext(req, rec)
		c.SetPath("/v1/coupons/:no")
		c.SetParamNames("no")
		c.SetParamValues(v.Result.Items[2].CouponNo)

		test.Ok(t, handleWithFilter(CouponController{}.GetOne, c))
		test.Equals(t, http.StatusNotFound, rec.Code)
	})
}
