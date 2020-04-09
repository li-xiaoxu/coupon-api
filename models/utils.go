package models

import (
	"context"

	"github.com/hublabs/common/auth"
)

func getTenantCode(ctx context.Context) string {
	user := auth.UserClaim{}.FromCtx(ctx)
	return user.Issuer
}

func getColleagueId(ctx context.Context) int64 {
	user := auth.UserClaim{}.FromCtx(ctx)
	return user.ColleagueId
}
