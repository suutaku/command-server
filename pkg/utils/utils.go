package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/aler9/goroslib"
	"github.com/aler9/goroslib/pkg/msgs/actionlib_msgs"
	"github.com/aler9/goroslib/pkg/msgs/geometry_msgs"
	"github.com/aler9/goroslib/pkg/msgs/nav_msgs"
	"github.com/aler9/goroslib/pkg/msgs/std_msgs"
)

type Service struct {
	name            string
	host            string
	sub             *goroslib.Subscriber
	currentLocation nav_msgs.Odometry
	node            *goroslib.Node
	pub             *goroslib.Publisher
	pubc            *goroslib.Publisher
	pq              *PositionQueue
}

func NewService(name string, host string) *Service {
	ip := os.Getenv("ROS_IP")
	if host == "" {
		ip = "127.0.0.1"
	}
	n, err := goroslib.NewNode(goroslib.NodeConf{
		Name:          name,
		MasterAddress: host,
		Host:          ip,
	})
	if err != nil {
		panic(err)
	}
	p, err := goroslib.NewPublisher(goroslib.PublisherConf{
		Node:  n,
		Topic: "/move_base_simple/goal",
		Msg:   &geometry_msgs.PoseStamped{},
		Latch: true,
	})
	if err != nil {
		panic(err)
	}

	p2, err := goroslib.NewPublisher(goroslib.PublisherConf{
		Node:  n,
		Topic: "/move_base/cancel",
		Msg:   &actionlib_msgs.GoalID{},
		Latch: true,
	})
	if err != nil {
		panic(err)
	}
	pq := NewPositionQueue(n)
	sv := Service{
		name: name,
		host: host,
		node: n,
		pub:  p,
		pubc: p2,
		pq:   pq,
	}
	go func() {
		for {
			log.Println("wait pose sequence")
			pose := <-pq.Item
			sv.PublishPoseStamped(pose)
		}
	}()
	return &sv
}

func (pub *Service) PublishPoseStamped(input geometry_msgs.Pose) error {
	input.Orientation = pub.currentLocation.Pose.Pose.Orientation
	data := geometry_msgs.PoseStamped{
		Pose: input,
		Header: std_msgs.Header{
			FrameId: "map",
		},
	}

	if pub.pub == nil {
		return fmt.Errorf("publisher is empty")
	}
	pub.pub.Write(&data)
	// p.Close()
	return nil
}

func (svc *Service) onMessage(msg *nav_msgs.Odometry) {
	//log.Printf("Incoming: %+v\n", msg.Pose)
	svc.currentLocation = *msg
	//log.Printf("Incoming: %+v\n", svc.currentLocation.Pose)
}

func (svc *Service) onGoalMessage(msg *geometry_msgs.PoseStamped) {
	//log.Printf("Incoming: %+v\n", msg.Pose)
	log.Printf("Incoming: %+v\n", svc.currentLocation.Pose)
}

func (svc *Service) Subscribe() error {
	// create a subscriber
	sub, err := goroslib.NewSubscriber(goroslib.SubscriberConf{
		Node:     svc.node,
		Topic:    "/odom",
		Callback: svc.onMessage,
	})
	if err != nil {
		return err
	}
	svc.sub = sub
	return nil
}

func (svc *Service) SubscribeTopic(topic string) error {
	// create a subscriber
	sub, err := goroslib.NewSubscriber(goroslib.SubscriberConf{
		Node:     svc.node,
		Topic:    topic,
		Callback: svc.onGoalMessage,
	})
	if err != nil {
		return err
	}
	svc.sub = sub
	return nil
}

func (svc *Service) StopTest() {
	svc.pubc.Write(&actionlib_msgs.GoalID{})
	log.Println("Write stop command")
}

func (svc *Service) GetCurrentLocation() nav_msgs.Odometry {
	return svc.currentLocation
}

func (svc *Service) Close() {
	svc.sub.Close()
}

func (svc *Service) AddQuene(qu []geometry_msgs.Pose) {
	svc.pq.AddQuene(qu)
}

func (svc *Service) CleanQueue() {
	svc.pq.Clean()
}
