package queue

import (
	"context"
	"fmt"
	"hublabs/coupon-api/models"
	"time"

	"github.com/hublabs/common/auth"

	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/labstack/gommon/log"
	"github.com/pangpanglabs/goutils/echomiddleware"
)

func SendTask(ctx context.Context, prepareId, lastCustId int64, seqNo int) error {
	pc, err := models.PrepareCoupon{}.Get(ctx, prepareId)
	if err != nil {
		return err
	}
	if pc.Enable == false || pc.SendStatus == models.SendStatusFinished {
		return nil
	}
	var eta time.Time
	if pc.SendType == models.SendTypeTimer {
		eta = pc.SendCondition.SendTime
	} else {
		//立即发放发券创建任务后5秒执行，以防API未执行完，campaign的status未改变
		eta = time.Now().Add(time.Second * 5)
	}
	signature := &tasks.Signature{
		Name: "send_coupon",
		Args: []tasks.Arg{
			{
				Type:  "int64",
				Value: prepareId,
			},
			{
				Type:  "int64",
				Value: lastCustId,
			},
			{
				Type:  "int",
				Value: seqNo,
			},
		},
		ETA: &eta,
	}
	asyncResult, err := Server.SendTask(signature)
	if err != nil {
		log.Fatal(err)
		return err
	}
	taskState := asyncResult.GetState()
	fmt.Println("=========", taskState)
	// if err := (models.PrepareCoupon{}).UpdateUuid(ctx, prepareId, taskState.TaskUUID); err != nil {
	// 	log.Fatal(err)
	// 	return err
	// }
	return nil
}

func SendCoupon(prepareId, lastCustId int64, seqNo int) error {
	//获取上下文
	session := db.NewSession()
	defer session.Close()
	ctx := context.WithValue(context.Background(), echomiddleware.ContextDBName, session)
	// 1.Get prepare coupon data
	pc, err := models.PrepareCoupon{}.Get(ctx, prepareId)
	if err != nil {
		return err
	}
	c, err := (models.CouponCampaign{}).Get(ctx, pc.CampaignId)
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, "userClaim", auth.UserClaim{
		ColleagueId: 1,
	})
	ctx = models.AddToCtx(ctx, c.TenantCode)

	if err != nil {
		return err
	}

	//修改prepareCoupon状态，准备发券
	// if err := pc.UpdateSendStatus(ctx, models.SendStatusSending); err != nil {
	// 	session.Rollback()
	// 	return err
	// }
	//添加发放记录数据(seqNo=-1,seqNo=0时是全部顾客)
	// if err := (models.SendRecord{}).CreateOrUpdate(ctx, pc.UUID, pc.Id, 0, -1, ""); err != nil {
	// 	return err
	// }
	//如果顾客条件未指定，不发券返回
	if len(pc.CustomerConditions) == 0 {
		return nil
	}
	//TODO:查询顾客并发券
	// if err := models.SendCouponByCustomer(ctx, pc, lastCustId, seqNo, c.TenantCode, c.Name, isTargetSpot); err != nil {
	// 	return err
	// }

	// if err := pc.UpdateSendStatus(ctx, models.SendStatusFinished); err != nil {
	// 	return err
	// }
	// if err := (models.SendRecord{}).Delete(ctx, pc.Id, pc.UUID); err != nil {
	// 	return err
	// }

	return nil
}
