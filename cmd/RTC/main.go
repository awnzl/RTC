package main

import (
	"log"
	"net/http"

	v1 "github.com/awnzl/RTC/api/v1"
	"github.com/awnzl/RTC/internal/hub"
	"github.com/awnzl/RTC/internal/session"
)

func main() {
	hb := hub.New()
	sm := session.NewService(hb)
	api := v1.New(sm)
	api.RegisterHandlers()

	s := http.Server{
		Addr:    "localhost:8080",
		Handler: nil,
	}
	log.Fatal(s.ListenAndServe())
}
