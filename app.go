package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/iftachsc/vmware"
	"github.com/iftachsc/zadara"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/vmware/govmomi"
)

var mySigningKey = []byte("secret")

type App struct {
	Router    *mux.Router
	DB        *mongo.Client
	VimClient *govmomi.Client
	//NetappClient *
	Zadata *zadara.Client
	ctx    context.Context
}

//Initialize this is great func
func (a *App) Initialize(user, password, dbname string) {

	ctx := context.Background()
	// connectionString :=
	// 	fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	//export COMPOSE_MONGODB_URL="mongodb://localhost:27017"
	//

	//
	//certpath, certavail := os.LookupEnv("PATH_TO_MONGODB_CERT")

	//
	var err error
	//var client *mongo.Client
	//if certavail {
	//	client, err := mongo.NewClientWithOptions(connectionString, clientopt.SSLOpt{CaFile: certpath})
	//} else {

	// client, err := mongo.NewClient(connectionString)
	// a.DB = client
	// //}
	// if err != nil {
	// 	log.Fatal(err)
	// }

	//initizalize vim client
	a.VimClient, err = vmware.NewClientFromEnv(ctx)
	if err != nil {
		log.Fatal(err)
	}
	a.ctx = ctx
	a.Router = mux.NewRouter()
	a.initializeRoutes()

}

func (a *App) initializeRoutes() {
	a.Router.Handle("/collection", jwtMiddleware.Handler(http.HandlerFunc(a.collect))).Methods("GET")
	a.Router.HandleFunc("/vms", a.getVms).Methods("GET")
	a.Router.Handle("/get-token", GetTokenHandler).Methods("GET")
	//a.Router.HandleFunc("/collections/{id:[0-9]+}", a.createCollection).Methods("GET")
}

//Run is
func (a *App) Run(addr string) {
	defer a.VimClient.Logout(a.ctx)

	println("Listening on", addr)
	log.Fatal(http.ListenAndServe(addr, handlers.LoggingHandler(os.Stdout, a.Router)))
}

func (a *App) createCollection(w http.ResponseWriter, r *http.Request) {
	err := a.DB.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	//collection := Collection()
}

func (a *App) collect(w http.ResponseWriter, r *http.Request) {

	//hosts, err := vmware.GetScsiLunDisks(a.VimClient, a.ctx)
	shuki, err := startCollection("sds", a.DB, nil, a.VimClient, a.ctx, "locationUuid")

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())

	} else {
		respondWithJSON(w, http.StatusOK, shuki)
	}
	// err := a.DB.Connect(context.Background())
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func (a *App) getVms(w http.ResponseWriter, r *http.Request) {

	hosts, err := vmware.GetEsxHost(a.VimClient, a.ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	} else {
		respondWithJSON(w, http.StatusOK, hosts)
	}
}

func (a *App) getHosts(w http.ResponseWriter, r *http.Request) {

	hosts, err := vmware.GetEsxHost(a.VimClient, a.ctx)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())

	} else {
		respondWithJSON(w, http.StatusOK, hosts)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)

}

var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	/* Create the token */
	token := jwt.New(jwt.SigningMethodHS256)

	//Create a map to store our claims
	claims := token.Claims.(jwt.MapClaims)

	/* Set token claims */
	claims["admin"] = true
	claims["name"] = "Ado Kukic"
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	/* Sign the token with our secret */
	tokenString, _ := token.SignedString(mySigningKey)

	/* Finally, write the token to the browser window */
	w.Write([]byte(tokenString))
})
