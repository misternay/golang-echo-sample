package handler

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
)

func TransValidator() ut.Translator {
	var universalTraslator *ut.UniversalTranslator
	e := en.New()
	universalTraslator = ut.New(e, e)
	trans, _ := universalTraslator.GetTranslator("en")
	trans.Add("required", "{0} is required", false)
	trans.Add("min", "{0} must be greater than {1}", true)
	trans.Add("max", "{0} must be less than {1}", true)
	trans.Add("eqfield", "{0} is not match with {1}", false)
	return trans
}
