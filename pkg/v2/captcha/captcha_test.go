package captcha

import (
	"testing"
)

func TestCaptcha(t *testing.T) {
	t.Run("Turnstile", func(t *testing.T) {
		captcha := NewCaptchaTurnstile("sitekey", "secret")
		if captcha.Type != "Turnstile" {
			t.Errorf("Expected Turnstile, got %s", captcha.Type)
		}
	})

	t.Run("Testing", func(t *testing.T) {
		captcha := NewCaptchaTesting("sitekey", "secret")
		if captcha.Type != "Testing" {
			t.Errorf("Expected Testing, got %s", captcha.Type)
		}
		err := captcha.Verify("token", "ip")
		if err != nil {
			t.Errorf("Expected nil, got %s", err)
		}
	})
}
