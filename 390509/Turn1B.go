package main

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CommandHandler handles API commands.
type CommandHandler struct {
	commandBus CommandBus
}

// NewCommandHandler creates a new CommandHandler.
func NewCommandHandler(commandBus CommandBus) *CommandHandler {
	return &CommandHandler{commandBus: commandBus}
}

// RegisterRoutes registers the API routes.
func (h *CommandHandler) RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/posts", h.handleCreatePost).Methods("POST")
	r.HandleFunc("/posts/{postID}", h.handleGetPost).Methods("GET")
}

func (h *CommandHandler) handleCreatePost(w http.ResponseWriter, r *http.Request) {
	var cmd CreatePostCommand
	if err := decodeJSON(r.Body, &cmd); err != nil {
		writeErrorResponse(w, err)
		return
	}

	resp, err := h.commandBus.Execute(r.Context(), cmd)
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	writeJSONResponse(w, resp)
}

func (h *CommandHandler) handleGetPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postID, err := strconv.Atoi(vars["postID"])
	if err != nil {
		writeErrorResponse(w, err)
		return
	}

	cmd := GetPostCommand{PostID: postID}
	resp, err := h.commandBus.Execute(r.Context(), cmd)
	if err != nil {
		writeErrorResponse(w, err)
		return