package main

import (
	"hublabs/coupon-api/config"
	"hublabs/coupon-api/controllers"
	"hublabs/coupon-api/factory"
	"hublabs/coupon-api/models"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/hublabs/common/auth"

	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pangpanglabs/echoswagger"
	"github.com/pangpanglabs/goutils/echomiddleware"
)

func main() {
	c := config.Init(os.Getenv("APP_ENV"))
	db := initDB(c.Database.Driver, c.Database.Connection, c.Debug)
	defer db.Close()

	factory.Init(db)
	// queue.Init(db)

	e := echo.New()
	r := echoswagger.New(e, "/doc", &echoswagger.Info{
		Title:       "Coupon API",
		Description: "This is docs for coupon-api service",
		Version:     "1.0.0",
	})

	r.AddSecurityAPIKey("Authorization", "JWT Token", echoswagger.SecurityInHeader)

	controllers.CouponCampaignController{}.Init(r.Group("CouponCamapign", "/v1/coupon-camapigns"))
	controllers.PrepareCouponController{}.Init(r.Group("PrepareCoupon", "/v1/prepare-coupons"))
	controllers.CouponController{}.Init(r.Group("Coupons", "/v1/coupons"))
	controllers.CouponInfoController{}.Init(r.Group("CouponInfos", "/v1/coupon-infos"))
	controllers.ExportController{}.Init(r.Group("Export", "/v1/exports"))

	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	e.Pre(middleware.RemoveTrailingSlash())
	e.Pre(echomiddleware.ContextBase())
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(echomiddleware.ContextLogger())
	e.Use(echomiddleware.ContextDB(c.ServiceName, db, echomiddleware.KafkaConfig(c.Database.Logger.Kafka)))
	e.Use(echomiddleware.BehaviorLogger(c.ServiceName, echomiddleware.KafkaConfig(c.BehaviorLog.Kafka)))
	e.Use(auth.UserClaimMiddleware("/ping", "/doc"))

	if err := e.Start(":5000"); err != nil {
		panic(err)
	}
}

func initDB(driver, connection string, debug bool) *xorm.Engine {
	db, err := xorm.NewEngine(driver, connection)
	if err != nil {
		panic(err)
	}

	if driver == "sqlite3" {
		runtime.GOMAXPROCS(1)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(time.Minute * 10)

	db.ShowSQL(debug)

	models.Init(db)
	return db
}
