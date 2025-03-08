package controller

import (
	"github.com/BargheNo/Backend/internal/domain/localization"
	"github.com/gin-gonic/gin"
)

func GetTranslator(c *gin.Context, key string) localization.TranslatorInstance {
	translator, exists := c.Get(key)
	if !exists {
		panic("translator not registered!")
	}

	return translator.(localization.TranslatorInstance)
}
