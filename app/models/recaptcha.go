package models

import (
	"github.com/haisum/recaptcha"
	"github.com/revel/revel"
	"net/http"
)

var captcha recaptcha.R

var secret = revel.Config.StringDefault("recaptcha.secret", "")

func InitializeRecaptcha() {
	captcha = recaptcha.R{ Secret: secret }
}

func ValidateRecaptcha(r *http.Request) bool {
	return captcha.Verify(*r)
}