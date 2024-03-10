package actions

import (
	"buffalo_aws_oidc_keycloak/public"
	"buffalo_aws_oidc_keycloak/templates"

	"github.com/gobuffalo/buffalo/render"
)

var tr *render.Engine

func init() {
	tr = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "tabler.plush.html",

		// fs.FS containing templates
		TemplatesFS: templates.FS(),

		// fs.FS containing assets
		AssetsFS: public.FS(),

		// Add template helpers here:
		Helpers: render.Helpers{
			// for non-bootstrap form helpers uncomment the lines
			// below and import "github.com/gobuffalo/helpers/forms"
			// forms.FormKey:     forms.Form,
			// forms.FormForKey:  forms.FormFor,
		},
	})
}
