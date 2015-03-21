package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/convox/kernel/web/Godeps/_workspace/src/github.com/codegangsta/negroni"
	"github.com/convox/kernel/web/Godeps/_workspace/src/github.com/ddollar/nlogger"
	"github.com/convox/kernel/web/Godeps/_workspace/src/github.com/gorilla/mux"

	"github.com/convox/kernel/web/controllers"
)

var port string = "5000"

func redirect(path string) func(http.ResponseWriter, *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		http.Redirect(rw, r, path, http.StatusFound)
	}
}

func parseForm(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	r.ParseMultipartForm(2048)
	next(rw, r)
}

func main() {
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	router := mux.NewRouter()

	router.HandleFunc("/", redirect("/apps")).Methods("GET")

	router.HandleFunc("/apps", controllers.AppList).Methods("GET")
	router.HandleFunc("/apps", controllers.AppCreate).Methods("POST")
	router.HandleFunc("/apps/{app}", controllers.AppShow).Methods("GET")
	router.HandleFunc("/apps/{app}", controllers.AppDelete).Methods("DELETE")
	router.HandleFunc("/apps/{app}/builds", controllers.AppBuilds).Methods("GET")
	router.HandleFunc("/apps/{app}/build", controllers.AppBuild).Methods("POST")
	router.HandleFunc("/apps/{app}/changes", controllers.AppChanges).Methods("GET")
	router.HandleFunc("/apps/{app}/logs", controllers.AppLogs)
	router.HandleFunc("/apps/{app}/logs/stream", controllers.AppLogStream)
	router.HandleFunc("/apps/{app}/processes/{process}", controllers.ProcessShow).Methods("GET")
	router.HandleFunc("/apps/{app}/processes/{process}/logs", controllers.ProcessLogs).Methods("GET")
	router.HandleFunc("/apps/{app}/processes/{process}/logs/stream", controllers.ProcessLogStream)
	router.HandleFunc("/apps/{app}/processes/{process}/resources", controllers.ProcessResources).Methods("GET")
	router.HandleFunc("/apps/{app}/promote", controllers.AppPromote).Methods("POST")
	router.HandleFunc("/apps/{app}/releases", controllers.AppReleases).Methods("GET")
	router.HandleFunc("/apps/{app}/resources", controllers.AppResources).Methods("GET")
	router.HandleFunc("/apps/{app}/services", controllers.AppServices).Methods("GET")
	router.HandleFunc("/apps/{app}/status", controllers.AppStatus).Methods("GET")

	n := negroni.New(
		negroni.NewRecovery(),
		nlogger.New("ns=kernel", nil),
		negroni.NewStatic(http.Dir("public")),
	)

	n.Use(negroni.HandlerFunc(parseForm))
	n.UseHandler(router)
	n.Run(fmt.Sprintf(":%s", port))
}
