package request

import (
	"net/http"

	"wa-blast/flags"
)

func JWTSubject(r *http.Request) string {
	return r.Header.Get(flags.HeaderKeyCOBRASubject)
}
