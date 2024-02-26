package grifts

import (
	"buffalo_aws_oidc_keycloak/actions"

	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
