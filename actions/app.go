package actions

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"buffalo_aws_oidc_keycloak/locales"
	"buffalo_aws_oidc_keycloak/public"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/middleware/csrf"
	"github.com/gobuffalo/middleware/forcessl"
	"github.com/gobuffalo/middleware/i18n"
	"github.com/gobuffalo/middleware/paramlogger"
	"github.com/unrolled/secure"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")

var (
	app     *buffalo.App
	appOnce sync.Once
	T       *i18n.Translator
)

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	appOnce.Do(func() {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_buffalo_aws_oidc_keycloak_session",
		})

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		app.Use(csrf.New)

		// Setup and use translations:
		app.Use(translations())

		app.GET("/", HomeHandler)

		app.GET("/login", LoginHandler)
		app.POST("/login", LoginHandler)

		auth := app.Group("/user")
		auth.Use(RequiresTokenMiddleware)
		auth.GET("/home", MainHandler)
		auth.GET("/logout", LogoutHandler)
		auth.GET("/sts", GetStsTokenHandler)
		auth.POST("/vm", GetVmHandler)

		app.ServeFiles("/", http.FS(public.FS())) // serve files from the public directory
	})

	return app
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(locales.FS(), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}

type UserData struct {
	Sub               string `json:"sub"`
	EmailVerified     bool   `json:"email_verified"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	Email             string `json:"email"`
}

func RequiresTokenMiddleware(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		bearerToken := c.Session().Get("access_token")
		if bearerToken == nil {
			c.Flash().Add("danger", "YOU ARD NOT AUTHORIZED!!")
			return c.Redirect(302, "/")
		}

		keycloakHost := os.Getenv("keycloakHost")
		realm := os.Getenv("realm")
		keycloakUrl := "https://" + keycloakHost + "/realms/" + realm + "/protocol/openid-connect/userinfo"

		req, _ := http.NewRequest("GET", keycloakUrl, nil)
		req.Header.Set("Authorization", "Bearer "+bearerToken.(string))

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			c.Flash().Add("danger", err.Error())
			return c.Redirect(302, "/")
		}

		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)
		if resp.Status != "200 OK" {
			c.Flash().Add("danger", "TOKEN EXPIRED")
			return c.Redirect(302, "/")
		}

		var userData UserData
		jsonerr := json.Unmarshal(body, &userData)
		if err != nil {
			c.Flash().Add("danger", "USER INFO ERR")
			c.Flash().Add("danger", jsonerr.Error())
			return c.Redirect(302, "/")
		}

		c.Set("Name", userData.Name)

		return next(c)
	}
}
