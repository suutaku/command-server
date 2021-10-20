package utils

import (
	"github.com/aler9/goroslib/pkg/msg"
	"github.com/aler9/goroslib/pkg/msgs/actionlib_msgs"
	"github.com/aler9/goroslib/pkg/msgs/std_msgs"
)

type MapMeta struct {
	Image          string
	Resolution     float32
	Origin         []float32
	Negate         float32
	Occupiedthresh float32
	FreeThresh     float32
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}

type Orientation struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
	W float64 `json:"w"`
}

type Pose struct {
	Position    Position    `json:"position"`
	Orientation Orientation `json:"orientation"`
}

type MoveBaseResult struct {
	msg.Package `ros:"move_base_msgs"`
}

type MoveBaseActionResult struct {
	msg.Package `ros:"move_base_msgs"`
	Header      std_msgs.Header
	Status      actionlib_msgs.GoalStatus
	Result      MoveBaseResult
}
