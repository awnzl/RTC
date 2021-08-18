package main

import (
	"log"

	"github.com/awnzl/RTC/internal/server"
)


func main() {
	s := server.New()
	log.Fatal(s.ListenAndServe())
}
