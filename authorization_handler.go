package service

import (
	"context"
	"net"
	"net/http"
	"strings"
)

type AuthorizationHandler struct {
	GetAndVerifyToken func(authorization string, secret string) (bool, string, map[string]interface{}, int64, int64, error)
	Secret            string
	Ip                string
	Authorization     string
}

func NewAuthorizationHandler(verifyToken func(string, string) (bool, string, map[string]interface{}, int64, int64, error), secret string, options ...string) *AuthorizationHandler {
	return NewAuthorizationHandlerWithIp(verifyToken, secret, "", options...)
}

func NewAuthorizationHandlerWithIp(verifyToken func(string, string) (bool, string, map[string]interface{}, int64, int64, error), secret string, ip string, options ...string) *AuthorizationHandler {
	var authorization string
	if len(options) >= 1 {
		authorization = options[0]
	}
	return &AuthorizationHandler{Authorization: authorization, GetAndVerifyToken: verifyToken, Secret: secret, Ip: ip}
}

func (c *AuthorizationHandler) HandleAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		au := r.Header["Authorization"]
		authorization := au[0]
		isToken, _, data, _, _, err := c.GetAndVerifyToken(authorization, c.Secret)
		var ctx context.Context
		ctx = r.Context()
		if len(c.Ip) > 0 {
			ip := GetRemoteIp(r)
			ctx = context.WithValue(ctx, c.Ip, ip)
		}
		if !isToken {
			if len(c.Ip) == 0 {
				next.ServeHTTP(w, r)
			} else {
				next.ServeHTTP(w, r.WithContext(ctx))
			}
		} else {
			if err != nil {
				if len(c.Ip) == 0 {
					next.ServeHTTP(w, r)
				} else {
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			} else {
				if len(c.Authorization) > 0 {
					ctx := context.WithValue(ctx, c.Authorization, data)
					next.ServeHTTP(w, r.WithContext(ctx))
				} else {
					for k, e := range data {
						if len(k) > 0 {
							ctx = context.WithValue(ctx, k, e)
						}
					}
					next.ServeHTTP(w, r.WithContext(ctx))
				}
			}
		}
	})
}

func ParseBearerToken(data []string) string {
	if len(data) == 0 {
		return ""
	}
	authorization := data[0]
	if strings.HasPrefix(authorization, "Bearer ") != true {
		return ""
	}
	return authorization[7:]
}
func GetRemoteIp(r *http.Request) string {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr
	}
	return remoteIP
}