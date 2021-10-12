package test

import (
	"log"
	"testing"
	"time"

	"github.com/aler9/goroslib/pkg/msgs/geometry_msgs"
	"github.com/suutaku/command-server/pkg/utils"
)

func TestPublish(t *testing.T) {
	log.SetFlags(log.Llongfile)
	svc := utils.NewService("test", "http://172.16.0.76:34107")
	if svc == nil {
		panic("bad")
	}
	pose := geometry_msgs.Pose{
		Position: geometry_msgs.Point{
			X: 0.0,
			Y: -3.36888885498,
			Z: -0.293332672119,
		},
		Orientation: geometry_msgs.Quaternion{
			X: 0.0,
			Y: 0.0,
			Z: -0.937883963896,
			W: 0.3469490,
		},
	}

	// err := svc.Subscribe()
	// if err != nil {
	// 	panic(err)
	// }

	// err := svc.SubscribeTopic("/move_base_simple/goal")
	// if err != nil {
	// 	panic(err)
	// }
	for {
		err := svc.PublishPoseStamped(pose)
		if err != nil {
			log.Println(err)
		}
		time.Sleep(5 * time.Second)
	}
	// select {}
}
