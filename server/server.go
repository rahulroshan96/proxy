package server

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type server struct {
	proxyServer *proxyServer
	muxServer   *muxServer
}

type muxServer struct {
	sharedConfig *sharedConfig
}

func (s *muxServer) SetSharedConfig(sharedConfig *sharedConfig) {
	s.sharedConfig = sharedConfig
}

func NewServer() *server {
	return &server{
		proxyServer: newProxyServer(),
		muxServer:   &muxServer{},
	}
}
func (s *server) getAllConfigs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.muxServer.sharedConfig.configuration)
}

func (s *server) createConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	insertID := uuid.New()
	pp := []*ProxyConfiguration{}
	var p ProxyConfiguration
	x := json.NewDecoder(r.Body).Decode(&p)
	pp = append(pp, &p)
	if x != nil {
		panic("Error while receiving request")
	}

	s.muxServer.sharedConfig.Lock()
	logrus.Info("inserting config")
	s.muxServer.sharedConfig.configuration[insertID.String()] = pp
	s.muxServer.sharedConfig.Unlock()
	json.NewEncoder(w).Encode(s.muxServer.sharedConfig.configuration)
}

func (s *server) deleteConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	if _, ok := s.muxServer.sharedConfig.configuration[params["id"]]; ok {
		s.muxServer.sharedConfig.Lock()
		logrus.Info("deleting config")
		delete(s.muxServer.sharedConfig.configuration, params["id"])
		s.muxServer.sharedConfig.Unlock()
		json.NewEncoder(w).Encode(s.muxServer.sharedConfig.configuration)
		return
	}
	json.NewEncoder(w).Encode("")
}
func (s *server) getConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	if _, ok := s.muxServer.sharedConfig.configuration[params["id"]]; ok {
		logrus.Info("deleting config")
		json.NewEncoder(w).Encode(s.muxServer.sharedConfig.configuration[params["id"]])
		return
	}
	json.NewEncoder(w).Encode("")
}

func (s *server) resetConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	s.muxServer.sharedConfig.Lock()
	logrus.Info("resetting config")
	s.muxServer.sharedConfig.configuration = make(map[string][]*ProxyConfiguration, 0)
	s.muxServer.sharedConfig.Unlock()
}

func (s *server) Run() {
	// Intialize the proxy server
	// then mux server to handle the configuration.
	sharedConfig := &sharedConfig{
		configuration: make(map[string][]*ProxyConfiguration, 0),
	}
	handler := &proxyUpdateHandler{
		sharedConfig: sharedConfig,
	}
	s.proxyServer.AddHandler(handler)
	s.muxServer.SetSharedConfig(sharedConfig)
	// run the proxy server
	go s.proxyServer.Run()
	router := mux.NewRouter()
	router.HandleFunc("/config", s.getAllConfigs).Methods("GET")
	router.HandleFunc("/config", s.createConfig).Methods("POST")
	router.HandleFunc("/config", s.resetConfig).Methods("DELETE")
	router.HandleFunc("/config/{id}", s.deleteConfig).Methods("DELETE")
	router.HandleFunc("/config/{id}", s.getConfig).Methods("GET")
	http.ListenAndServe(":4996", router)
}
