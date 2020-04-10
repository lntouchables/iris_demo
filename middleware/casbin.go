package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12/context"
	"github.com/casbin/casbin/v2"
	"iris/models"
	"net/http"
	"strconv"
	"time"
)

func New(e *casbin.Enforcer) *Casbin {
	return &Casbin{enforcer: e}
}

func (c *Casbin) ServeHTTP(ctx context.Context) {
	value := ctx.Values().Get("jwt").(*jwt.Token)
	token := models.OauthToken{}
	token.GetOauthTokenByToken(value.Raw)
	if token.Revoked || token.ExpressIn < time.Now().Unix(){
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.StopExecution()
		return
	}else if !c.Check(ctx.Request(), strconv.FormatUint(uint64(token.UserId), 10)) {
		ctx.StatusCode(http.StatusForbidden) // Status Forbidden
		ctx.StopExecution()
		return
	} else {
		ctx.Values().Set("auth_user_id", token.UserId)
	}

	ctx.Next()
}

type Casbin struct {
	enforcer *casbin.Enforcer
}


func (c *Casbin) Check(r *http.Request, userId string) bool {
	method := r.Method
	path := r.URL.Path
	ok, _ := c.enforcer.Enforce(userId, path, method)
	return ok
}