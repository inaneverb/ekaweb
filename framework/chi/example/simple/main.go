package main

import (
	"fmt"
	"net/http"

	"github.com/inaneverb/ekaweb/framework/chi/v2"
	"github.com/inaneverb/ekaweb/v2"
)

func main() {
	var r = ekaweb_chi.NewRouter(ekaweb.WithServerName("chi.example")).
		Get("/*", handler) // <-- catch all
	panic(http.ListenAndServe(":8081", r.Build()))
}

func handler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w,
		"Hello from:\n\tRegistered route: %s\n\tActual route: %s\n",
		ekaweb.RoutePath(r), r.RequestURI)
}
