package service

import (
	"encoding/json"
	"net/http"

	"github.com/clockworksoul/cog2/data/rest"
	"github.com/clockworksoul/cog2/dataaccess"
	"github.com/gorilla/mux"
)

// handleGetGroups handles "GET /v2/group"
func handleGetGroups(w http.ResponseWriter, r *http.Request) {
	group, err := dataaccess.GroupList()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(group)
}

// handleGetGroup handles "GET /v2/group/{groupname}"
func handleGetGroup(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	exists, err := dataaccess.GroupExists(params["groupname"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		http.Error(w, "No such group", http.StatusNotFound)
		return
	}

	group, err := dataaccess.GroupGet(params["groupname"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(group)
}

// handlePostGroup handles "POST /v2/group/{groupname}"
func handlePostGroup(w http.ResponseWriter, r *http.Request) {
	var group rest.Group
	var err error

	params := mux.Vars(r)

	_ = json.NewDecoder(r.Body).Decode(&group)

	group.Name = params["groupname"]

	exists, err := dataaccess.GroupExists(group.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// NOTE: Should we just make "update" create groups that don't exist?
	if exists {
		err = dataaccess.GroupUpdate(group)
	} else {
		err = dataaccess.GroupCreate(group)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// handlePostGroup handles "DELETE /v2/group/{groupname}"
func handleDeleteGroup(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	err := dataaccess.GroupDelete(params["groupname"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func addGroupMethodsToRouter(router *mux.Router) {
	router.HandleFunc("/v2/group", handleGetGroups).Methods("GET")
	router.HandleFunc("/v2/group/{groupname}", handleGetGroup).Methods("GET")
	router.HandleFunc("/v2/group/{groupname}", handlePostGroup).Methods("POST")
	router.HandleFunc("/v2/group/{groupname}", handleDeleteGroup).Methods("DELETE")
}