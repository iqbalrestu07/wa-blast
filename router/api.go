package router

import (
	"wa-blast/flags"
	"wa-blast/handler"
)

// routeAPI configure request routing in API. Handlers must be defined in handler package
func routeAPI(r Router) {

	r.HandleREST("/", handler.Health, flags.ACLEveryone).Methods("GET")
	r.HandleREST("/account/sync", handler.Sync, flags.ACLAuthenticatedUser).Methods("POST")
	r.HandleREST("/account/message", handler.Message, flags.ACLAuthenticatedUser).Methods("POST")
	r.HandleREST("/blast/message", handler.BlastMessage, flags.ACLAuthenticatedUser).Methods("POST")
}
