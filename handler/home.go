package handler

import (
	"net/http"

	"wa-blast/configs"
	"wa-blast/response"
)

// Health ...
func Health(r *http.Request) (*response.Success, error) {
	result := make(map[string]interface{}, 1)
	result["version"] = configs.MustGetString("server.version")
	return response.NewSuccess(result, nil), nil
}
