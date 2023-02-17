package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type AuthenticationSessionRequest struct {
	*AuthenticationSession

	// User *UserPayload `json:"user,omitempty"`

	// ProtectedID string `json:"id"` // override 'id' json to have more control
}

type AuthenticationSession struct {
	Subject      string                 `json:"subject"`
	Extra        map[string]interface{} `json:"extra"`
	Header       http.Header            `json:"header"`
	MatchContext MatchContext           `json:"match_context"`
}

type MatchContext struct {
	RegexpCaptureGroups []string    `json:"regexp_capture_groups"`
	URL                 *url.URL    `json:"url"`
	Method              string      `json:"method"`
	Header              http.Header `json:"header"`
}

func (a *AuthenticationSessionRequest) Bind(r *http.Request) error {
	// a.Article is nil if no Article fields are sent in the request. Return an
	// error to avoid a nil pointer dereference.
	if a.Subject == "" {
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

func (a *AuthenticationSessionRequest) Render(w http.ResponseWriter, r *http.Request) error {
	// Pre-processing before a response is marshalled and sent across the wire
	// a.Elapsed = 10
	return nil
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	// r.Post("/", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("welcome"))
	// 	log.Printf("Post body: %s", r.Body)
	// })
	r.Post("/", LogAuthenticationSession)
	http.ListenAndServe(":8080", r)
}

func LogAuthenticationSession(w http.ResponseWriter, r *http.Request) {
	data := &AuthenticationSessionRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	// article := data.Subject
	// dbNewArticle(article)

	json, _ := json.Marshal(data.AuthenticationSession)
	log.Printf("Post body: %v", string(json))

	render.Status(r, http.StatusCreated)
	render.Render(w, r, NewAuthenticationSessionResponse(data))
}

func NewAuthenticationSessionResponse(session *AuthenticationSessionRequest) *AuthenticationSessionRequest {
	return session
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request.",
		ErrorText:      err.Error(),
	}
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}
