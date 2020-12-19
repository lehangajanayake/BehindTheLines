package models

//Coordinates contains the x and y value
type Coordinates struct{
	X, Y int
}

//Animation model for animation
type Animation struct{
	Animate bool
	FrameNum, CurrentFrame, FrameHeight, FrameWidth int
	AnimationArray []Coordinates
}