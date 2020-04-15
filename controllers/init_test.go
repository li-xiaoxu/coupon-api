package controllers

import (
	"context"
	"hublabs/coupon-api/config"
	"hublabs/coupon-api/factory"
	"hublabs/coupon-api/models"
	"net/http"
	"os"
	"runtime"

	"github.com/hublabs/common/auth"

	"github.com/asaskevich/govalidator"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pangpanglabs/goutils/behaviorlog"
	"github.com/pangpanglabs/goutils/echomiddleware"
	"github.com/pangpanglabs/goutils/jwtutil"

	configutil "github.com/pangpanglabs/goutils/config"
)

var (
	echoApp          *echo.Echo
	handleWithFilter func(handlerFunc echo.HandlerFunc, c echo.Context) error
)

func init() {
	runtime.GOMAXPROCS(1)
	configutil.SetConfigPath("../")
	c := config.Init(os.Getenv("APP_ENV"))
	xormEngine, err := xorm.NewEngine(c.Database.Driver, c.Database.Connection)
	if err != nil {
		panic(err)
	}
	if err = models.DropTables(xormEngine); err != nil {
		panic(err)
	}
	if err = models.Init(xormEngine); err != nil {
		panic(err)
	}
	factory.Init(xormEngine)
	// xormEngine.ShowSQL()
	echoApp = echo.New()
	echoApp.Validator = &Validator{}

	db := echomiddleware.ContextDB("test", xormEngine, echomiddleware.KafkaConfig{})
	jwtSecret := middleware.JWT([]byte(os.Getenv("JWT_SECRET")))
	behaviorlogger := func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()
			c.SetRequest(req.WithContext(context.WithValue(req.Context(),
				behaviorlog.LogContextName, behaviorlog.New("test", req),
			)))
			return next(c)
		}
	}

	handleWithFilter = func(handlerFunc echo.HandlerFunc, c echo.Context) error {
		return behaviorlogger(jwtSecret(auth.UserClaimMiddleware()(db(handlerFunc))))(c)
	}
}

type Validator struct{}

func (v *Validator) Validate(i interface{}) error {
	_, err := govalidator.ValidateStruct(i)
	return err
}
func setHeader(r *http.Request) {
	token, _ := jwtutil.NewToken(map[string]interface{}{"aud": "colleague", "tenantCode": "pangpang", "iss": "pangpang"})
	r.Header.Set(echo.HeaderAuthorization, "Bearer "+token)
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
}
