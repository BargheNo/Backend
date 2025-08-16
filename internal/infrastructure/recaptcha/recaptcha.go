package recaptcha

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/BargheNo/Backend/bootstrap"
)

type Recaptcha struct {
	Secret string
}

func NewRecaptcha(env *bootstrap.Recaptcha) *Recaptcha {
	return &Recaptcha{
		Secret: env.Secret,
	}
}

func (r *Recaptcha) Verify(token string) error {
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
		url.Values{"secret": {r.Secret}, "response": {token}})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var result struct {
		Success bool `json:"success"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}
	if !result.Success {
		return errors.New("invalid captcha")
	}
	return nil
}
