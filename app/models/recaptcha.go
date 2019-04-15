package models

import (
	"github.com/haisum/recaptcha"
	"net/http"
)

var captcha recaptcha.R

const secret = "6LfMMJ4UAAAAABMBTYd_oiV2zUKqZBJ30daEe3zF"

func InitializeRecaptcha() {
	captcha = recaptcha.R{ Secret: secret }
}

func ValidateRecaptcha(r *http.Request) bool {
	return captcha.Verify(*r)
}