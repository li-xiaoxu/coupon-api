package controllers

import (
	"errors"
	"hublabs/coupon-api/models"
	"net/http"
	"time"

	"github.com/hublabs/common/api"
	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
)

type CouponController struct{}

func (c CouponController) Init(g echoswagger.ApiGroup) {
	g.SetSecurity("Authorization")

	g.GET("", c.GetAll).
		AddParamQueryNested(SearchInput{})
	g.GET("/:no", c.GetOne).
		AddParamPath("", "no", "No of Coupon").
		AddParamQueryNested(FieldInput{})
	g.POST("/:no/use", c.Use).
		AddParamPath("", "no", "Use Coupon")
	g.POST("/:no/recover", c.Recover).
		AddParamPath("", "no", "Recover Coupon")
	g.DELETE("/:no", c.Delete).
		AddParamPath("", "no", "Delete Coupon")
}

func (CouponController) GetAll(c echo.Context) error {
	var v CustSearchInput
	if err := c.Bind(&v); err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}
	if v.Limit == 0 {
		v.Limit = DefaultMaxResultCount
	}
	totalCount, items, err := models.Coupon{}.GetAll(c.Request().Context(), v.CustId, v.SaleType, v.Filter, v.Nos, v.Sortby, v.Order, v.Offset, v.Limit)
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	return renderSuccArray(c, false, false, totalCount, items)
}

func (CouponController) GetOne(c echo.Context) error {
	var v FieldInput
	if err := c.Bind(&v); err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}
	coupon, err := models.Coupon{}.Get(c.Request().Context(), c.Param("no"), v.Fields)
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	if coupon == nil {
		return renderFail(c, api.ErrorNotFound.New(errors.New("Coupon is not exist")))
	}
	return renderSucc(c, http.StatusOK, coupon)
}

func (CouponController) Use(c echo.Context) error {
	no := c.Param("no")
	coupon, err := models.Coupon{}.Get(c.Request().Context(), no, nil)
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	} else if coupon == nil || coupon.Status != models.CouponStatusNormal {
		return renderFail(c, api.ErrorParameter.New(errors.New("Coupon not found")))
	} else if coupon.UseChk {
		return renderFail(c, api.ErrorParameter.New(errors.New("Coupon has been used")))
	} else if coupon.EndAt.Before(time.Now()) {
		return renderFail(c, api.ErrorParameter.New(errors.New("Coupon expired")))
	} else if coupon.StartAt.After(time.Now()) {
		return renderFail(c, api.ErrorParameter.New(errors.New("Coupon has not yet taken effect")))
	}
	coupon.UseChk = true
	if err := coupon.Update(c.Request().Context(), coupon); err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	return renderSucc(c, http.StatusOK, nil)
}

func (CouponController) Recover(c echo.Context) error {
	no := c.Param("no")
	coupon, err := models.Coupon{}.Get(c.Request().Context(), no, nil)
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	} else if coupon == nil {
		return renderFail(c, api.ErrorParameter.New(errors.New("Coupon not found")))
	} else if !coupon.UseChk {
		return renderFail(c, api.ErrorParameter.New(errors.New("Coupon has not been used")))
	}
	colleagueId := models.GetColleagueId(c.Request().Context())
	if err != nil {
		return renderFail(c, api.ErrorParameter.New(errors.New("Colleague is not exist")))
	}
	r := &models.RecoverRecord{
		CouponNo: no,
		UserId:   colleagueId,
		UseStore: coupon.UseStore,
		UseAt:    coupon.Commit.UpdatedAt,
	}
	coupon.UseChk = false
	coupon.UseStore = ""
	if err := coupon.Update(c.Request().Context(), coupon); err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	if err := r.Create(c.Request().Context()); err != nil {
		return err
	}
	return renderSucc(c, http.StatusOK, nil)
}

func (CouponController) Delete(c echo.Context) error {
	no := c.Param("no")
	coupon, err := models.Coupon{}.Get(c.Request().Context(), no, nil)
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	} else if coupon == nil {
		return renderFail(c, api.ErrorParameter.New(errors.New("Coupon not found")))
	}
	if err := coupon.Delete(c.Request().Context(), coupon.CouponNo); err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	return renderSucc(c, http.StatusNoContent, nil)
}
