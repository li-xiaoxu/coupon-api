package controllers

import (
	"hublabs/coupon-api/models"
	"hublabs/coupon-api/queue"
	"net/http"
	"strconv"

	"github.com/hublabs/common/api"
	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
	"github.com/pangpanglabs/goutils/converter"
)

type CouponCampaignController struct{}

func (c CouponCampaignController) Init(g echoswagger.ApiGroup) {
	g.SetSecurity("Authorization")

	g.GET("/:id", c.GetOne).
		AddParamPath(0, "id", "Id of ouponCampaign")
	g.GET("", c.GetAll).
		AddParamQueryNested(SearchInput{})
	g.POST("", c.Create).
		AddParamBody(models.CouponCampaign{}, "body", "CouponCampaign input", true)
	g.PUT("/:id", c.Update).
		AddParamPath(0, "id", "Id of couponCampaign").
		AddParamBody(models.CouponCampaign{}, "body", "couponCampaign input", true)

}

func (CouponCampaignController) GetOne(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}
	camapigns, err := models.CouponCampaign{}.Get(c.Request().Context(), id)
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	if camapigns == nil {
		return renderFail(c, api.ErrorNotFound.New(err))
	}
	return renderSucc(c, http.StatusOK, camapigns)
}

func (CouponCampaignController) GetAll(c echo.Context) error {
	var v SearchInput
	if err := c.Bind(&v); err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}
	ids := converter.StringToIntSlice(v.Ids)
	totalCount, campaigns, err := models.CouponCampaign{}.GetAll(c.Request().Context(), v.Q, v.Filter, ids, v.Offset, v.Limit, v.Sortby, v.Order, models.CampaignStatus(v.Status), v.Fields)
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	return renderSuccArray(c, false, false, totalCount, campaigns)
}

func (CouponCampaignController) Create(c echo.Context) error {
	var v models.CouponCampaign
	if err := c.Bind(&v); err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}
	if err := v.Create(c.Request().Context()); err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	if v.Status == models.CampaignStatusApprove {
		for i := range v.PrepareCoupons {
			if v.PrepareCoupons[i].SendType == models.SendTypeNow || v.PrepareCoupons[i].SendType == models.SendTypeTimer {
				if err := queue.SendTask(c.Request().Context(), v.PrepareCoupons[i].Id, 0, -1); err != nil {
					return renderFail(c, api.ErrorDB.New(err))
				}
			}
			if v.PrepareCoupons[i].SendType == models.SendTypeFree && v.PrepareCoupons[i].MaxQty > 0 {
				if err := models.CreateFreeCoupon(c.Request().Context(), v.PrepareCoupons[i]); err != nil {
					return renderFail(c, api.ErrorDB.New(err))
				}
			}
		}
	}
	return renderSucc(c, http.StatusCreated, v)
}

func (CouponCampaignController) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}
	var v models.CouponCampaign
	if err := c.Bind(&v); err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}

	campaign, err := models.CouponCampaign{}.Get(c.Request().Context(), id)
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	if campaign == nil {
		return renderFail(c, api.ErrorNotFound.New(err))
	}
	if campaign.Status != v.Status && v.Status == models.CampaignStatusApprove {
		colleagueId := models.GetColleagueId(c.Request().Context())
		if err != nil {
			return renderFail(c, api.ErrorDB.New(err))
		}
		v.Commit.UpdatedBy = string(colleagueId)
	}
	if err := v.Update(c.Request().Context(), id); err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	if v.PrepareCoupons == nil {
		v.PrepareCoupons = campaign.PrepareCoupons
	}
	if campaign.Status != v.Status && v.Status == models.CampaignStatusApprove {
		for i := range v.PrepareCoupons {
			if v.PrepareCoupons[i].SendType == models.SendTypeNow || v.PrepareCoupons[i].SendType == models.SendTypeTimer {
				if err := queue.SendTask(c.Request().Context(), v.PrepareCoupons[i].Id, 0, -1); err != nil {
					return renderFail(c, api.ErrorDB.New(err))
				}
			}
			if v.PrepareCoupons[i].SendType == models.SendTypeFree && v.PrepareCoupons[i].MaxQty > 0 {
				if err := models.CreateFreeCoupon(c.Request().Context(), v.PrepareCoupons[i]); err != nil {
					return renderFail(c, api.ErrorDB.New(err))
				}
			}
		}
	}
	return renderSucc(c, http.StatusNoContent, nil)
}
