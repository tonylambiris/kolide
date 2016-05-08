package v1

import (
	"html/template"

	"github.com/mephux/kolide/controller/helpers"
	"github.com/mephux/kolide/version"
)

// HelperFunctions to use while rendering content
var HelperFunctions = template.FuncMap{
	"Name": func() string {
		return version.Name
	},
	"Version": func() string {
		return version.Version
	},
	"DateFormat": helpers.DateFormat,
}
