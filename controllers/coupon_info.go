package controllers

import (
	"hublabs/coupon-api/models"
	"net/http"
	"strconv"

	"github.com/hublabs/common/api"
	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
)

type CouponInfoController struct{}

func (c CouponInfoController) Init(g echoswagger.ApiGroup) {
	g.SetSecurity("Authorization")

	g.GET("/:id", c.GetOne).
		AddParamPath("0", "id", "id of couponInfo")
	g.GET("", c.GetAll).
		AddParamQueryNested(SearchInput{})
	g.POST("", c.Create).
		AddParamBody(models.CouponInfo{}, "body", "input of couponInfo", true)
}

func (CouponInfoController) GetOne(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}
	couponInfo, err := models.CouponInfo{}.Get(c.Request().Context(), id)
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	if couponInfo == nil {
		return renderFail(c, api.ErrorNotFound.New(err))
	}
	return renderSucc(c, http.StatusOK, couponInfo)
}

func (CouponInfoController) GetAll(c echo.Context) error {
	var v SearchInput
	if err := c.Bind(&v); err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}
	items, totalCount, err := models.CouponInfo{}.GetAll(c.Request().Context(), v.Q, v.Sortby, v.Order, v.Offset, v.Limit)
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	return renderSuccArray(c, false, false, totalCount, items)
}

func (CouponInfoController) Create(c echo.Context) error {
	var v models.CouponInfo
	if err := c.Bind(&v); err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}
	if err := v.Create(c.Request().Context()); err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	return renderSucc(c, http.StatusCreated, v)
}
