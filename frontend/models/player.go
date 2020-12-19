package models

import (
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

//Player model player
type Player struct{
	Img *ebiten.Image
	Coords Coordinates
	FacingFront bool
	Op *ebiten.DrawImageOptions
	WalkingAnimation Animation
	IdleAnimation Animation
}

