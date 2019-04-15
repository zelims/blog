package utils

import "github.com/zelims/blog/app/models"

func Initialize() {
	models.InitializeRecaptcha()
}