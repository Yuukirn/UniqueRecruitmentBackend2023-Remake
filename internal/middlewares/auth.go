package middlewares

import (
	"UniqueRecruitmentBackend/global"
	"UniqueRecruitmentBackend/internal/common"
	"UniqueRecruitmentBackend/internal/constants"
	error2 "UniqueRecruitmentBackend/internal/error"
	"UniqueRecruitmentBackend/internal/tracer"
	"context"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"net/http"
)

func ctxWithUID(ctx context.Context, uid string) context.Context {
	return context.WithValue(ctx, "X-UID", uid)
}

func AuthMiddleware(c *gin.Context) {
	apmCtx, span := tracer.Tracer.Start(c.Request.Context(), "Authentication")
	defer span.End()

	cookie, err := c.Cookie("uid")
	if errors.Is(err, http.ErrNoCookie) {
		c.Abort()
		common.Error(c, error2.UnauthorizedError)
		return
	}
	s := sessions.Default(c)
	u := s.Get(cookie)
	if u == nil {
		c.Abort()
		common.Error(c, error2.UnauthorizedError)
		return
	}
	uid, ok := u.(string)
	if !ok {
		c.Abort()
		common.Error(c, error2.UnauthorizedError)
		return
	}
	c.Request = c.Request.WithContext(ctxWithUID(apmCtx, uid))

	span.SetAttributes(attribute.String("UID", uid))
	c.Next()
}

func RoleMiddleware(c *gin.Context, role constants.Role) {
	apmCtx, span := tracer.Tracer.Start(c.Request.Context(), "Authentication")
	defer span.End()

	uid := common.GetUID(c)
	client := global.GetSSOClient()
	ok, err := client.CheckPermissionByRole(apmCtx, uid, string(role))
	if err != nil || !ok {
		c.Abort()
		common.Error(c, error2.CheckPermissionError)
		return
	}

	c.Next()
}
