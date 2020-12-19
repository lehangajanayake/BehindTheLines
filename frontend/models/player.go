package models

import "github.com/hajimehoshi/ebiten/v2"

//Player model player
type Player struct{
	Img *ebiten.Image
	Coords Coordinates
	FacingFront bool
	Op *ebiten.DrawImageOptions
	WalkingAnimation Animation
	IdleAnimation Animation
}

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