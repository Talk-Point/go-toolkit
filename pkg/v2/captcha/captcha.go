// Package captcha provides a client for Cloudflare Turnstile Captcha.
//
// The package allows you to create and verify captchas of different types.
// Currently, it supports two types of captchas: Turnstile and Testing.
// The Turnstile captcha is a real captcha that requires verification,
// while the Testing captcha is a dummy captcha used for testing purposes.
//
// Each captcha is represented by a Captcha struct, which contains the
// necessary information for captcha verification, such as the site key,
// secret, and the type of captcha.
//
// The package provides two functions for creating captchas: NewCaptchaTurnstile
// and NewCaptchaTesting. These functions return a pointer to a Captcha struct
// with the specified site key, secret, and type.
//
// The Captcha struct has a Verify method that takes a token and an IP address,
// and verifies the captcha based on its type. If the captcha type is not
// supported, the method returns an error.
//
// Example Usage:
//
//	cfToken := c.FormValue("cf-turnstile-response")
//	ip := c.RealIP()
//	cap := captcha.NewCaptchaTurnstile(h.CaptchaConfig.SiteKey, h.CaptchaConfig.SecretKey)
//	if cfToken == "" {
//	    return Render(c, http.StatusOK, login.LoginForm(cap, h.AuthenticatorId, login.LoginErrors{
//	        Email:              email,
//	        InvalidCredentials: "Verification failed",
//	    }))
//	}
package captcha

import (
	"fmt"

	"github.com/Talk-Point/go-toolkit/pkg/v2/captcha/turnstile"
)

type CaptchaType int

const (
	Turnstile CaptchaType = iota
	Testing
)

func (ct CaptchaType) String() string {
	return [...]string{"Turnstile", "Testing"}[ct]
}

type Captcha struct {
	IsActive bool
	SiteKey  string
	Secret   string
	Type     string
}

func (c *Captcha) Verify(token string, ip string) error {
	if c.Type == Turnstile.String() {
		return turnstile.VerifyRequest(c.Secret, token, ip)
	} else if c.Type == Testing.String() {
		return nil
	}

	return fmt.Errorf("Captcha type not supported: %s", c.Type)
}

func NewCaptchaTurnstile(siteKey string, secret string) *Captcha {
	return &Captcha{
		IsActive: true,
		SiteKey:  siteKey,
		Secret:   secret,
		Type:     Turnstile.String(),
	}
}

func NewCaptchaTesting(siteKey string, secret string) *Captcha {
	return &Captcha{
		IsActive: true,
		SiteKey:  siteKey,
		Secret:   secret,
		Type:     Testing.String(),
	}
}
