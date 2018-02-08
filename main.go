package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"
)

func main() {
	scope := gmail.GmailReadonlyScope
	cfg, err := google.ConfigFromJSON([]byte(os.Getenv("CLIENT_SECRET_JSON")), scope)
	if err != nil {
		log.Fatal(err)
	}
	cfg.ClientSecret = os.Getenv("CLIENT_SECRET")
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), &server{cfg: cfg}))
}

type server struct{ cfg *oauth2.Config }

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		// TODO: pick a better value for state?
		http.Redirect(w, r, s.cfg.AuthCodeURL("state"), http.StatusFound)
	case "/oauth2callback":
		// TODO: validate r.FormValue("state")?
		token, err := s.cfg.Exchange(r.Context(), r.FormValue("code"))
		if err != nil {
			log.Print(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, token.AccessToken)
	default:
		http.NotFound(w, r)
	}
}
