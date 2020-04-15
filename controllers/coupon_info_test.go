package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hublabs/coupon-api/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/pangpanglabs/goutils/test"
)

func CouponInfoTest(t *testing.T) {
	inputs := []map[string]interface{}{
		{
			"title":      "测试couponInfo_title",
			"desc":       "测试couponInfo_desc",
			"notice":     "测试couponInfo_notice",
			"allowExtra": true,
		},
		{
			"title":      "测试title",
			"desc":       "测试desc",
			"notice":     "测试notice",
			"allowExtra": true,
		},
		{
			"title":      "couponInfo_title",
			"desc":       "couponInfo_desc",
			"notice":     "couponInfo_notice",
			"allowExtra": true,
		},
		{
			"title":      "测试一下couponInfo_title",
			"desc":       "测试一下couponInfo_desc",
			"notice":     "测试一下couponInfo_notice",
			"allowExtra": true,
		},
		{
			"title":      "测试测试couponInfo_title",
			"desc":       "测试测试couponInfo_desc",
			"notice":     "测试测试couponInfo_notice",
			"allowExtra": true,
		},
	}
	for i, c := range inputs {
		pb, _ := json.Marshal(c)
		t.Run(fmt.Sprint("Create#", i+1), func(t *testing.T) {
			req := httptest.NewRequest(echo.POST, "/v1/coupon-infos", bytes.NewReader(pb))
			setHeader(req)
			rec := httptest.NewRecorder()
			test.Ok(t, handleWithFilter(CouponInfoController{}.Create, echoApp.NewContext(req, rec)))
			test.Equals(t, http.StatusCreated, rec.Code)
		})
	}
	t.Run("GetAll", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/v1/coupon-infos", nil)
		setHeader(req)
		rec := httptest.NewRecorder()
		test.Ok(t, handleWithFilter(CouponInfoController{}.GetAll, echoApp.NewContext(req, rec)))
		test.Equals(t, http.StatusOK, rec.Code)

		var v struct {
			Result struct {
				TotalCount int                 `json:"totalCount"`
				Items      []models.CouponInfo `json:"items"`
			} `json:"result"`
			Success bool `json:"success"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, v.Result.TotalCount, 5)
		test.Equals(t, v.Result.Items[0].Title, "测试couponInfo_title")
		test.Equals(t, v.Result.Items[1].Title, "测试title")
		test.Equals(t, v.Result.Items[2].Title, "couponInfo_title")
		test.Equals(t, v.Result.Items[3].Title, "测试一下couponInfo_title")
		test.Equals(t, v.Result.Items[4].Title, "测试测试couponInfo_title")
	})
	t.Run("GetAllByTitle", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/v1/coupon-infos?q=测试", nil)
		setHeader(req)
		rec := httptest.NewRecorder()
		test.Ok(t, handleWithFilter(CouponInfoController{}.GetAll, echoApp.NewContext(req, rec)))
		test.Equals(t, http.StatusOK, rec.Code)

		var v struct {
			Result struct {
				TotalCount int                 `json:"totalCount"`
				Items      []models.CouponInfo `json:"items"`
			} `json:"result"`
			Success bool `json:"success"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, v.Result.TotalCount, 4)
		test.Equals(t, v.Result.Items[0].Title, "测试couponInfo_title")
		test.Equals(t, v.Result.Items[1].Title, "测试title")
		test.Equals(t, v.Result.Items[2].Title, "测试一下couponInfo_title")
		test.Equals(t, v.Result.Items[3].Title, "测试测试couponInfo_title")
	})
	t.Run("GetAllByTitle#2", func(t *testing.T) {
		req := httptest.NewRequest(echo.GET, "/v1/coupon-infos?q=测试couponInfo", nil)
		setHeader(req)
		rec := httptest.NewRecorder()
		test.Ok(t, handleWithFilter(CouponInfoController{}.GetAll, echoApp.NewContext(req, rec)))
		test.Equals(t, http.StatusOK, rec.Code)

		var v struct {
			Result struct {
				TotalCount int                 `json:"totalCount"`
				Items      []models.CouponInfo `json:"items"`
			} `json:"result"`
			Success bool `json:"success"`
		}
		test.Ok(t, json.Unmarshal(rec.Body.Bytes(), &v))
		test.Equals(t, v.Result.TotalCount, 1)
		test.Equals(t, v.Result.Items[0].Title, "测试couponInfo_title")
	})
}
