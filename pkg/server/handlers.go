package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aler9/goroslib/pkg/msgs/geometry_msgs"
	"github.com/suutaku/command-server/pkg/utils"
	"gopkg.in/yaml.v2"
)

func (server *Server) targetHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Map-Resolution,Origin-X,Origin-Y")

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
		return
	}
	var checkMap = make(map[string]interface{})
	json.Unmarshal(b, &checkMap)
	if checkMap["positions"] != nil {
		fmt.Println("positions come")
		ps, ok := checkMap["positions"].([]geometry_msgs.Pose)
		if ok {
			err = server.service.PublishPoseStamped(ps[0])
			if err != nil {
				log.Println(err)
				w.Write([]byte(err.Error()))
			}
			server.AddQuene(ps[1:])
		}
	} else {
		var pose geometry_msgs.Pose
		err = json.Unmarshal(b, &pose)
		if err != nil {
			log.Println(err)
			w.Write([]byte(err.Error()))
			return
		}
		server.service.StopTest()
		server.CleanQueue()
		log.Println(string(b))
		err = server.service.PublishPoseStamped(pose)
		if err != nil {
			log.Println(err)
			w.Write([]byte(err.Error()))
		}
	}

}

func (server *Server) runshHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	server.service.StopTest()
	log.Println("runsh request")
}

func (server *Server) locationHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Expose-Headers", "Map-Resolution,Origin-X,Origin-Y")
	log.Println("location request")

	w.Header().Set("Access-Control-Allow-Origin", "*")

	b, err := json.Marshal(server.service.GetCurrentLocation().Pose.Pose)
	if err != nil {
		w.Write([]byte("location not available"))
	}
	var ret utils.Pose
	json.Unmarshal(b, &ret)
	nb, err := json.Marshal(ret)
	if err != nil {
		w.Write([]byte("location not available"))
	}
	//log.Println(server.service.GetCurrentLocation().Pose.Pose)
	w.Write(nb)
}

func (server *Server) mapHandler(w http.ResponseWriter, r *http.Request) {

	// load yaml
	//log.Println("map request")
	b, err := ioutil.ReadFile("/root/maps/map.yaml")
	if err != nil {
		log.Println(err)
		return
	}
	var tmp utils.MapMeta
	err = yaml.Unmarshal(b, &tmp)
	if err != nil {
		log.Println(err)
		return
	}
	w.Header().Set("Map-Resolution", fmt.Sprintf("%f", tmp.Resolution))
	w.Header().Set("Origin-X", fmt.Sprintf("%f", tmp.Origin[0]))
	w.Header().Set("Origin-Y", fmt.Sprintf("%f", tmp.Origin[1]))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Map-Resolution,Origin-X,Origin-Y")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	bb, err := ioutil.ReadFile(tmp.Image)
	if err != nil {
		w.Write([]byte("map not available"))
		log.Println(err)
		return
	}
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(bb)))
	w.Write(bb)
}

func (server *Server) restartHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	server.service.StopTest()
}
