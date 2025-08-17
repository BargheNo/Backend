package recaptcha

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/BargheNo/Backend/bootstrap"
	"github.com/BargheNo/Backend/internal/domain/exception"
)

type Recaptcha struct {
	constants *bootstrap.Constants
	Secret    string
}

func NewRecaptcha(
	constants *bootstrap.Constants,
	env *bootstrap.Recaptcha,
) *Recaptcha {
	return &Recaptcha{
		constants: constants,
		Secret:    env.Secret,
	}
}

func (r *Recaptcha) Verify(token string) error {
	resp, err := http.PostForm("https://www.google.com/recaptcha/api/siteverify",
		url.Values{"secret": {r.Secret}, "response": {token}})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
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
		var validationErrors exception.ValidationErrors
		validationErrors.Add(r.constants.Field.Recaptcha, r.constants.Tag.InvalidRecaptcha)
		return validationErrors
	}
	return nil
}
