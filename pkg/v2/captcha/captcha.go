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
