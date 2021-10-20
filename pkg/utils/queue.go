package utils

import (
	"fmt"

	"github.com/aler9/goroslib"
	"github.com/aler9/goroslib/pkg/msgs/actionlib_msgs"
	"github.com/aler9/goroslib/pkg/msgs/geometry_msgs"
)

type PositionQueue struct {
	Positions []geometry_msgs.Pose
	Sub       *goroslib.Subscriber
	Item      chan geometry_msgs.Pose
}

func NewPositionQueue(node *goroslib.Node) *PositionQueue {
	pos := PositionQueue{}
	postions := make([]geometry_msgs.Pose, 0)
	sub, err := goroslib.NewSubscriber(goroslib.SubscriberConf{
		Node:     node,
		Topic:    "/move_base/result",
		Callback: pos.onMoveBaseResult,
	})
	if err != nil {
		fmt.Println(err)
	}
	pos.Item = make(chan geometry_msgs.Pose)
	pos.Positions = postions
	pos.Sub = sub
	return &pos
}

func (pq *PositionQueue) getFirst() (geometry_msgs.Pose, error) {
	if len(pq.Positions) == 0 {
		return geometry_msgs.Pose{}, fmt.Errorf("Empty queue")
	}
	item := pq.Positions[0]
	pq.Positions = pq.Positions[1:]
	return item, nil
}

func (pq *PositionQueue) Clean() {
	pq.Positions = make([]geometry_msgs.Pose, 0)
}

func (pq *PositionQueue) onMoveBaseResult(msg *MoveBaseActionResult) {
	fmt.Println("onMoveBaseResult")
	fmt.Printf("%+v\n", msg)
	switch msg.Status.Status {
	case actionlib_msgs.GoalStatus_SUCCEEDED:
		item, err := pq.getFirst()
		if err != nil {
			fmt.Println(err)
			return
		}
		pq.Item <- item
		fmt.Println("Send new position")
	case actionlib_msgs.GoalStatus_PREEMPTED:
		pq.Clean()
		fmt.Println("Clean all positions")
	default:
		fmt.Println("Unkown state ", msg.Status.Status)
	}

}

func (pq *PositionQueue) AddQuene(qu []geometry_msgs.Pose) {
	pq.Positions = qu
}
