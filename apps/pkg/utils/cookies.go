package utils

import (
	"forestwhisperGo/apps/internal/models/types"
	"time"

	"github.com/gofiber/fiber/v3"
)

const (
	cookieName = "forestwhisper"
	cookieTTL  = 7 * 24 * time.Hour
)

func SetCookieToken(c fiber.Ctx, token string) {
	c.Cookie(&fiber.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		MaxAge:   int(cookieTTL.Seconds()),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
	})
}

func GetCookieToken(c fiber.Ctx) *types.Payload {
	token := c.Cookies(cookieName)
	if token == "" {
		return nil
	}
	payload, err := ValidateToken(token)
	if err != nil {
		return nil
	}
	return payload
}
