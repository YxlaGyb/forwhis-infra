package utils

import (
	"forestwhisperGo/magic/models/types"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	cookieName = "forestwhisper"
	cookieTTL  = 7 * 24 * time.Hour
)

func SetCookieToken(c *gin.Context, token string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(cookieTTL.Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
}

func GetCookieToken(c *gin.Context) *types.Payload {
	token, err := c.Cookie(cookieName)
	if err != nil || token == "" {
		return nil
	}
	payload, err := ValidateToken(token)
	if err != nil {
		return nil
	}
	return payload
}
