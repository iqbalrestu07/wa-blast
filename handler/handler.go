// Package handler provides handler functions to handle request.
package handler

import (
	"fmt"
	"html"
	"io"
	"os"
	"path/filepath"
	"time"

	"encoding/json"
	"net/http"
	"net/url"

	"wa-blast/configs"
	"wa-blast/flags"
	"wa-blast/loggers"
	"wa-blast/query"
	"wa-blast/request"
	"wa-blast/response"
	"wa-blast/util"

	"github.com/gorilla/context"
	uuid "github.com/satori/go.uuid"
)

// Logger
var log = loggers.Get()

// RESTFunc is a handler function that handles error and writes response in JSON
type RESTFunc func(*http.Request) (*response.Success, error)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://infobankdata.octarine.id")
}

// ServeHTTP implement http.Handler interface to write success or error response in JSON
func (h RESTFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	// Get execution time
	start := time.Now()
	// Init http status
	var httpStatus int
	// Execute handler
	result, err := h(r)
	// If an error returned, return error
	if err != nil {
		httpStatus = sendErrorJSON(w, err)
	} else if result == nil {
		w.WriteHeader(http.StatusNoContent)
	} else {
		// if header exist, add header to response
		if result.Header != nil {
			for k, v := range result.Header {
				w.Header().Set(k, v)
			}
		}
		// send json success
		httpStatus = sendJSON(w, http.StatusOK, result)
	}
	// Log elapsed time
	log.Infof("HTTP Status: %d, Request: %s %s, Time elapsed: %s", httpStatus, r.Method, html.EscapeString(r.URL.Path), time.Since(start))
}

// sendJSON write response in JSON
func sendJSON(w http.ResponseWriter, httpStatus int, obj interface{}) int {
	// Add content type
	w.Header().Add(flags.HeaderKeyContentType, flags.ContentTypeJSON)
	// Write http status
	w.WriteHeader(httpStatus)
	// Send JSON response
	json.NewEncoder(w).Encode(obj)
	// Return httpStatus
	return httpStatus
}

// sendErrorJSON write error response in JSON
func sendErrorJSON(w http.ResponseWriter, err error) int {
	// Cast error to ApiError
	apiError := util.CastError(err)
	// Send error json
	return sendJSON(w, apiError.Status, apiError)
}

// parseJSON parse json request body to o (target) and returns error
func parseJSON(r *http.Request, o interface{}) error {
	d := json.NewDecoder(r.Body)
	if err := d.Decode(o); err != nil {
		return err
	}
	return nil
}

// NotFound Returns not found in json
func NotFound(_ *http.Request) (*response.Success, error) {
	return nil, util.NewError("404")
}

// MethodNotAllowed Returns not found in json
func MethodNotAllowed(_ *http.Request) (*response.Success, error) {
	return nil, util.NewError("404")
}

// Pagination ...
func Pagination(v url.Values) (limit int, skip int) {
	// Get skip
	skip = util.ParseInt(v.Get("skip"), 0)
	// Get limit
	limit = util.ParseInt(v.Get("limit"), 10)
	// Return pagination options
	return limit, skip
}

// FileUpload ...
func FileUpload(r *http.Request, field string, require bool) (string, error) {

	u2 := uuid.NewV4()

	uploadedFile, handler, err := r.FormFile(field)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}
	defer uploadedFile.Close()

	filename := u2.String() + handler.Filename

	fileLocation := filepath.Join(configs.MustGetString("file.images") + field + "/" + filename)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return "", err
	}

	return filename, err
}

// GetAuth ...
func GetAuth(r *http.Request) (string, error) {
	if fmt.Sprintf("%v", context.Get(r, "type")) == "jwt" {
		return request.JWTSubject(r), nil
	}

	if fmt.Sprintf("%v", context.Get(r, "type")) == "3digit" {
		return fmt.Sprintf("%v", context.Get(r, "id")), nil
	}

	return "", util.NewError("-1013")
}

// Filter ....
func Filter(reqQuery url.Values, filters map[string]query.Filter) map[string]query.Filter {
	// Get filter values from query
	for key, f := range filters {

		// If filter is available in query, add values
		if v := reqQuery.Get(key); v != "" {
			err := f.Set(v)
			// If error setting filter, delete from map
			if err != nil {
				delete(filters, key)
			}
		} else {
			delete(filters, key)
		}
	}
	return filters
}

func GetIDAuth(r *http.Request) string {
	return context.Get(r, "id").(string)
}
