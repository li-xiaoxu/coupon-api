package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"hublabs/coupon-api/models"
	"net/http"
	"strconv"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/hublabs/common/api"
	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
)

type ExportController struct{}

func (c ExportController) Init(g echoswagger.ApiGroup) {
	g.GET("/free/:id", c.ExportFreeCoupon).
		AddParamPath(0, "id", "Id of PrepareCoupon")
}

func (ExportController) ExportFreeCoupon(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}
	campaign, err := models.CouponCampaign{}.GetByPrepareCouponId(c.Request().Context(), id)
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	if campaign == nil {
		return renderFail(c, api.ErrorParameter.New(errors.New("CouponCampaign is not exist")))
	}

	coupons, err := models.Coupon{}.GetByPrepareCouponIds(c.Request().Context(), id)
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}

	name := fmt.Sprintf("CouponList-%d.xlsx", id)
	xFile := excelize.NewFile()
	sheet := xFile.NewSheet("优惠券")

	xFile.SetCellValue("优惠券", "A1", "活动名称:")
	xFile.SetCellValue("优惠券", "B1", campaign.Name)
	xFile.SetCellValue("优惠券", "C1", "活动期间:")
	xFile.SetCellValue("优惠券", "D1", fmt.Sprintf("%s~%s", campaign.StartAt.Format("2006-01-02 15:04:05"), campaign.EndAt.Format("2006-01-02 15:04:05")))
	xFile.SetCellValue("优惠券", "A2", "优惠券号:")

	for i, c := range coupons {
		xFile.SetCellValue("优惠券", fmt.Sprintf("A%d", i+3), c.CouponNo)
	}
	xFile.SetActiveSheet(sheet)
	//将数据存入buff中
	var buff bytes.Buffer
	if err = xFile.Write(&buff); err != nil {
		panic(err)
	}
	//设置请求头  使用浏览器下载
	c.Response().Header().Set(echo.HeaderContentDisposition, "attachment; filename="+name)
	return c.Stream(http.StatusOK, echo.MIMEOctetStream, bytes.NewReader(buff.Bytes()))
}
