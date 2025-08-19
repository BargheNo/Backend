package controller

import (
	"reflect"

	"github.com/BargheNo/Backend/internal/domain/exception"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate = validator.New()

func Validated[T any](ctx *gin.Context) T {
	var params T
	if err := ctx.ShouldBindUri(&params); err != nil {
		bindingError := exception.BindingError{Err: err}
		panic(bindingError)
	}

	if err := ctx.ShouldBindQuery(&params); err != nil {
		bindingError := exception.BindingError{Err: err}
		panic(bindingError)
	}

	if err := ctx.ShouldBind(&params); err != nil {
		bindingError := exception.BindingError{Err: err}
		panic(bindingError)
	}

	if err := validate.Struct(params); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			panic(formatValidationErrors[T](validationErrors))
		}
		panic(err)
	}
	return params
}

func formatValidationErrors[T any](errs validator.ValidationErrors) exception.ValidationErrors {
	var params T
	var validationErrors exception.ValidationErrors
	tagTypes := []string{"json", "uri", "form"}

	t := reflect.TypeOf(params)

	for _, e := range errs {
		field, _ := t.FieldByName(e.StructField())
		tagValue := getAnyTag(field, tagTypes...)
		validationErrors.Add(tagValue, e.Tag())
	}

	return validationErrors
}

func getAnyTag(field reflect.StructField, tagNames ...string) string {
	for _, tagName := range tagNames {
		if tag := field.Tag.Get(tagName); tag != "" {
			return tag
		}
	}
	return field.Name
}
