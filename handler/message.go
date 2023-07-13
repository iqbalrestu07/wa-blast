package handler

import (
	"net/http"
	"wa-blast/request"
	"wa-blast/response"
	"wa-blast/services"
	"wa-blast/util"
)

// Sync ...
func Sync(r *http.Request) (*response.Success, error) {

	id := GetIDAuth(r)

	var req request.SyncAccount
	err := parseJSON(r, &req)
	if err != nil {
		return nil, util.NewError("400")
	}

	result, err := services.Message.SyncAccount(id, req)
	if err != nil {
		return nil, err
	}

	return result, err
}

// Message ...
func Message(r *http.Request) (*response.Success, error) {

	id := GetIDAuth(r)

	var req request.Message
	err := parseJSON(r, &req)
	if err != nil {
		return nil, util.NewError("400")
	}

	result, err := services.Message.Message(id, req)
	if err != nil {
		return nil, err
	}

	return result, err
}

// Message ...
func BlastMessage(r *http.Request) (*response.Success, error) {

	id := GetIDAuth(r)

	var req request.BlastMessage
	err := parseJSON(r, &req)
	if err != nil {
		return nil, util.NewError("400")
	}

	result, err := services.Message.Blast(id, req)
	if err != nil {
		return nil, err
	}

	return result, err
}
