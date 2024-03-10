package actions

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
)

// HomeHandler is a default handler to serve up
// a home page.
func HomeHandler(c buffalo.Context) error {
	c.Set("pretitle", "Noting here")
	c.Set("title", "iamDash")
	return c.Render(http.StatusOK, tr.HTML("home/dash.html"))
}
