package bloompass

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type checkRequest struct {
	Password string `json:"password"`
}

type checkResponse struct {
	Success bool `json:"success"`
	Exists  int  `json:"exists"`
}

type ApiServer struct {
	filter *Bloom
	server *http.Server
}

func NewApiServer(host string, port string, bloom *Bloom) *ApiServer {
	srv := &ApiServer{
		filter: bloom,
	}
	srv.server = &http.Server{
		Addr:    host + port,
		Handler: srv.buildRouter(),
	}
	return srv
}

func (a *ApiServer) actionCheck(w http.ResponseWriter, r *http.Request) {
	var request checkRequest
	result := checkResponse{
		Success: false,
		Exists:  EXIST_NO,
	}
	reqBody, err := ioutil.ReadAll(r.Body)
	if err == nil {
		if err = json.Unmarshal(reqBody, &request); err == nil {
			result.Exists = a.filter.Exists(request.Password)
			result.Success = true
		}
	}
	json.NewEncoder(w).Encode(result)
}

func (a *ApiServer) buildRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	// we only accept SEARCH http action
	router.HandleFunc("/check", a.actionCheck).Methods("SEARCH")
	return router
}

func (a *ApiServer) Start() error {
	return a.server.ListenAndServe()
}

func (a *ApiServer) Stop(ctx context.Context) error {
	return a.server.Shutdown(ctx)
}
