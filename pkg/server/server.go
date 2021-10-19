package server

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aler9/goroslib/pkg/msgs/geometry_msgs"
	"github.com/rs/cors"
	"github.com/suutaku/command-server/pkg/utils"
)

type Server struct {
	port    string
	service *utils.Service
	handler http.Handler
}

func NewServer(port string) *Server {
	host := os.Getenv("ROS_MASTER_URI")
	if host == "" {
		host = "172.16.2.182:34107"
	}
	log.Println("Get master uri:", host)
	return &Server{
		port:    port,
		service: utils.NewService("command-server", host),
	}
}

func (server *Server) InitServer() {

	r := http.NewServeMux()
	r.HandleFunc("/target", server.targetHandler)
	r.HandleFunc("/runsh", server.runshHandler)
	r.HandleFunc("/restart", server.restartHandler)
	r.HandleFunc("/location", server.locationHandler)
	r.HandleFunc("/map", server.mapHandler)
	server.handler = cors.Default().Handler(r)

	// cors.Default() setup the middleware with default options being
	// all origins accepted with simple methods (GET, POST). See
	// documentation below for more options.
	server.handler = cors.Default().Handler(r)

}

func (server *Server) Serve() {
	log.Println("Subscribe")
SUB:
	err := server.service.Subscribe()
	if err != nil {
		log.Println(err)
		log.Println("Retry subscribing")
		time.Sleep(time.Second * 5)
		goto SUB
	} else {
		log.Println("Subscribe success")
	}

	log.Println("Listen commands")
	// srv := &http.Server{
	// 	ReadTimeout:  5 * time.Second,
	// 	WriteTimeout: 10 * time.Second,
	// 	Handler:      server.handler,
	// 	Addr:         ":" + server.port,
	// }

	http.ListenAndServe(":"+server.port, server.handler)
}

func (server *Server) AddQuene(qu []geometry_msgs.Pose) {
	server.service.AddQuene(qu)
}

func (server *Server) CleanQueue() {
	server.service.CleanQueue()
}
