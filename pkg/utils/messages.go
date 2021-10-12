package utils

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
