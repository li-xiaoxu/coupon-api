package controllers

import (
	"context"
	"errors"
	"hublabs/coupon-api/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/hublabs/common/api"
	"github.com/labstack/echo"
	"github.com/pangpanglabs/echoswagger"
	"github.com/renstrom/shortuuid"
	"github.com/sirupsen/logrus"
)

type PrepareCouponController struct{}

func (c PrepareCouponController) Init(g echoswagger.ApiGroup) {
	g.SetSecurity("Authorization")

	g.GET("/:id", c.Get).
		AddParamPath(0, "id", "Id of prepareCoupon")
	g.POST("", c.Create).
		AddParamBody(models.PrepareCoupon{}, "body", "models of prepareCoupon", true)
	g.POST("/:id/send", c.SendCoupon).
		SetDescription("Send coupon to customer").
		AddParamPath(0, "id", "Id of PrepareCoupon").
		AddParamBody(SendInput{}, "body", "SendInput model", true)
	g.POST("/new", c.SendNewCoupon).
		AddParamBody(SendInput{}, "body", "SendInput model", true)
	g.POST("/birth", c.SendBirthCoupon).
		AddParamBody(SendInput{}, "body", "SendInput model", true)
	// g.POST("/handle/:sendType", c.HandleCoupon)
}

func (PrepareCouponController) Get(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}
	pc, err := models.PrepareCoupon{}.Get(c.Request().Context(), id)
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	if pc == nil {
		return renderFail(c, api.ErrorNotFound.New(err))
	}

	return renderSucc(c, http.StatusOK, pc)
}

func (PrepareCouponController) Create(c echo.Context) error {
	var v models.PrepareCoupon
	if err := c.Bind(&v); err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}
	if err := v.Create(c.Request().Context()); err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	return renderSucc(c, http.StatusCreated, nil)
}

func (PrepareCouponController) SendCoupon(c echo.Context) error {
	var v SendInput
	if err := c.Bind(&v); err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}
	if err := c.Validate(&v); err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}
	pc, err := models.PrepareCoupon{}.Get(c.Request().Context(), id)
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	} else if pc == nil {
		return renderFail(c, api.ErrorNotFound.New(errors.New("Invalid PrepareCoupon")))
	}

	coupon, err := sendCouponToCust(c.Request().Context(), *pc, v.CustId, v.BrandId)
	if err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}

	return renderSucc(c, http.StatusOK, coupon)
}

func (PrepareCouponController) SendNewCoupon(c echo.Context) error {
	var v SendInput
	if err := c.Bind(&v); err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}
	if err := c.Validate(&v); err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}
	//查询有效期内符合顾客条件的新会员触发的prepareCoupon
	_, items, err := models.PrepareCoupon{}.GetAll(c.Request().Context(), "", "true", v.CustId, string(models.SendTypeNew), string(models.CampaignStatusApprove), []string{"id"}, []string{"asc"}, 0, 0)
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	//如果没有新会员发券检查是否有生日券
	// if len(items) == 0 {
	// 	return ReturnApiSucc(c, http.StatusOK, 0)
	// }

	// 判断是否存在分批发券
	var (
		alreadyPrepare []int64 //已经循环过的prepareCoupon的ID
		couponCount    int
	)

	isExistPre := func(id int64) bool {
		for i := range alreadyPrepare {
			if alreadyPrepare[i] == id {
				return true
			}
		}
		return false
	}

	for _, item := range items {
		//list为相同批次的prepareCoupon
		var list []models.PrepareCoupon

		if isExistPre(item.Id) {
			continue
		}
		for _, pc := range items {
			if pc.CampaignId == item.CampaignId && pc.SeqNo == item.SeqNo {
				list = append(list, pc)
				alreadyPrepare = append(alreadyPrepare, pc.Id)
			}
		}

		index := 0
		if len(list) > 1 {
			pIndex, err := models.FindPrepareCoupon(list)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"error": err,
				}).Error("New customer find prepareCoupon error")
				continue
			}
			index = pIndex
		}
		coupon, err := sendCouponToCust(c.Request().Context(), list[index], v.CustId, v.BrandId)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("New customer send error")
			continue
		}
		if coupon != nil {
			couponCount++
		}
	}

	//生日券发放
	if v.Birthday == "" || time.Now().Format("01") != strings.Split(v.Birthday, "-")[1] {
		return renderSucc(c, http.StatusOK, couponCount)
	}
	_, pcs, err := models.PrepareCoupon{}.GetAll(c.Request().Context(), "", "true", v.CustId, string(models.SendTypeBirth), string(models.CampaignStatusApprove), []string{"id"}, []string{"asc"}, 0, 0)
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	for _, p := range pcs {
		_, err := sendCouponToCust(c.Request().Context(), p, v.CustId, v.BrandId)
		if err != nil {
			return renderFail(c, api.ErrorParameter.New(err))
		}
		couponCount++
	}

	return renderSucc(c, http.StatusOK, couponCount)
}

func (PrepareCouponController) SendBirthCoupon(c echo.Context) error {
	var v SendInput
	if err := c.Bind(&v); err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}
	if err := c.Validate(&v); err != nil {
		return renderFail(c, api.ErrorParameter.New(err))
	}
	if time.Now().Format("01") != strings.Split(v.Birthday, "-")[1] {
		return renderSucc(c, http.StatusOK, 0)
	}
	_, items, err := models.PrepareCoupon{}.GetAll(c.Request().Context(), "", "true", v.CustId, string(models.SendTypeBirth), string(models.CampaignStatusApprove), []string{"id"}, []string{"asc"}, 0, 0)
	if err != nil {
		return renderFail(c, api.ErrorDB.New(err))
	}
	var couponCount int
	for _, p := range items {
		_, err := sendCouponToCust(c.Request().Context(), p, v.CustId, v.BrandId)
		if err != nil {
			return renderFail(c, api.ErrorDB.New(err))
		}
		couponCount++
	}
	return renderSucc(c, http.StatusOK, couponCount)
}

func sendCouponToCust(ctx context.Context, pc models.PrepareCoupon, custId string, brandId int64) (*models.Coupon, error) {
	cp, err := models.CouponCampaign{}.Get(ctx, pc.CampaignId)
	if err != nil {
		return nil, err
	} else if cp == nil {
		return nil, errors.New("Invalid PrepareCoupon")
	}

	if pc.MaxPerQty > 0 && pc.ReceivedInfo.CustQty >= pc.MaxPerQty {
		return nil, errors.New("More than per cust maximum limit")
	} else if pc.MaxQty > 0 && pc.ReceivedInfo.Qty >= pc.MaxQty {
		return nil, errors.New("More than total maximum limit")
	} else if cp.EndAt.Before(time.Now()) {
		return nil, errors.New("Campaign expired")
	} else if cp.StartAt.After(time.Now()) {
		return nil, errors.New("Campaign has not yet taken effect")
	}
	p := pc.CouponPeriod.GetCouponPeriod(time.Now())
	coupon := &models.Coupon{
		CouponNo:        shortuuid.New(),
		PrepareCouponId: pc.Id,
		OfferId:         pc.OfferId,
		CouponInfo:      pc.CouponInfo,
		CustId:          custId,
		StartAt:         p.StartAt,
		EndAt:           p.EndAt,
		Status:          models.CouponStatusNormal,
	}
	if err := coupon.Create(ctx); err != nil {
		return nil, err
	}
	// if pc.Alert.WeChatAlert || pc.Alert.SmsAlert {
	// 	if err := models.SendAlertToCust(ctx, brandId, *cp, pc, *coupon); err != nil {
	// 		return nil, errors.New("Send alert error")
	// 	}
	// }

	return coupon, nil
}
