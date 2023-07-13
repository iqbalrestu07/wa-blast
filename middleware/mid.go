package middleware

import (
	"net/http"

	"wa-blast/configs"
	"wa-blast/handler"
	"wa-blast/response"
	"wa-blast/services"
	"wa-blast/util"

	"github.com/gorilla/context"
)

// Middleware ...
type Middleware func(handler.RESTFunc, ...string) handler.RESTFunc

// Auth ...
func Auth() Middleware {
	// Create middleware
	m := func(next handler.RESTFunc, args ...string) handler.RESTFunc {
		// Define new handler
		h := func(r *http.Request) (*response.Success, error) {

			user, pass, ok := r.BasicAuth()

			if !ok {
				return nil, util.NewError("401")
			}

			ok = validateBasicAuth(user, pass)
			if !ok {
				return nil, util.NewError("401")
			}

			// company := models.CompaniesAuth{
			// 	ID:        "123",
			// 	SecretKey: "123",
			// }

			company, err := services.User.ValidateAuth(r.Header.Values("key")[0])
			if err != nil {
				return nil, err
			}

			context.Set(r, "id", company.ID)

			return next(r)

		}
		return h
	}
	// Return middleware
	return m
}

func validateBasicAuth(user, pass string) bool {

	if user != configs.MustGetString("admin.user") {
		return false
	}

	if pass != configs.MustGetString("admin.pass") {
		return false
	}

	return true
}
