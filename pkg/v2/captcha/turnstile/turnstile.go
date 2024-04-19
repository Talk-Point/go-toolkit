package turnstile

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

func VerifyRequest(secret string, token string, ip string) error {
	formData := url.Values{}
	formData.Set("secret", secret)
	formData.Set("response", token)
	formData.Set("remoteip", ip)

	const url = "https://challenges.cloudflare.com/turnstile/v0/siteverify"

	req, err := http.NewRequest("POST", url, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var outcome map[string]interface{}
	if err := json.Unmarshal(body, &outcome); err != nil {
		return err
	}

	if success, ok := outcome["success"].(bool); ok && success {
		return nil
	}
	return errors.New("verification failed")
}
