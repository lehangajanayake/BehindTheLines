package models

import "github.com/hajimehoshi/ebiten/v2"

//Coordinates contains the x and y value
type Coordinates struct {
	X, Y int
}

//Animation model for animation
type Animation struct {
	Img                                                             *ebiten.Image
	Animate                                                         bool
	FrameNum, CurrentFrame, FrameHeight, FrameWidth, StartX, StartY int
}

//Reset the animation
func (a *Animation) Reset() {
	a.Animate = false
	a.CurrentFrame = 1
}
