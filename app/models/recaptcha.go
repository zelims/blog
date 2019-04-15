package models

import (
	"github.com/haisum/recaptcha"
	"net/http"
)

var captcha recaptcha.R

const secret = ""

func InitializeRecaptcha() {
	captcha = recaptcha.R{ Secret: secret }
}

func ValidateRecaptcha(r *http.Request) bool {
	return captcha.Verify(*r)
}