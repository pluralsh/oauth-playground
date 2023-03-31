package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/render"
)

type BootstrapRequest struct {
	UserID string `json:"userId,omitempty"`
}

func (b *BootstrapRequest) Bind(r *http.Request) error {
	// a.Article is nil if no Article fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if b.UserID == "" {
		return errors.New("missing required Article fields.")
	}

	// a.User is nil if no Userpayload fields are sent in the request. In this app
	// this won't cause a panic, but checks in this Bind method may be required if
	// a.User or futher nested fields like a.User.Name are accessed elsewhere.

	// just a post-process after a decode..
	// a.ProtectedID = ""                                 // unset the protected ID
	// a.Article.Title = strings.ToLower(a.Article.Title) // as an example, we down-case
	return nil
}

func (b *BootstrapRequest) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	// a.Elapsed = 10
	return nil
}

func (h *Handler) BootstrapAdmin(w http.ResponseWriter, r *http.Request) {
	log := h.Log.WithName("BootstrapAdmin")
	data := &BootstrapRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	currentAdmins, err := h.C.GetOrganizationAdmins(context.Background(), "main") // TODO: decide if we want to harcode this
	if err != nil {
		log.Error(err, "failed to get current admins")
		render.Render(w, r, ErrFailedToGetAdmins(err))
		return
	}

	if len(currentAdmins) == 0 {
		log.Info("No current admins, adding initial user as admin", "user", data.UserID)
		err = h.C.AddAdminToOrganization(context.Background(), "main", data.UserID)
		if err != nil {
			log.Error(err, "failed to add bootstrap user")
			render.Render(w, r, ErrFailedToSetInitialAdmin(err))
			return
		}
	}

	json, _ := json.Marshal(data.UserID)
	log.Info("Post", "body", string(json))

	render.Status(r, http.StatusOK) // TODO: we should probably do error checking above so we can return a 400 if something goes wrong
	render.Render(w, r, h.BootstrapUser(data))
}

func (h *Handler) BootstrapUser(b *BootstrapRequest) *BootstrapRequest {
	return b
}
