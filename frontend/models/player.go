package models

import (
	"fmt"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

//Player model player
type Player struct {
	Img                                                *ebiten.Image
	Coords, LastPos                                    Coordinates
	FacingFront, LastFacing                            bool
	LastAnimation                                      string
	Op                                                 *ebiten.DrawImageOptions
	Gun                                                Gun
	WalkingAnimation, IdleAnimation, ShootingAnimation Animation
}

//Render reders a plyer in the screen
func (p *Player) Render(screen, img *ebiten.Image) {
	screen.DrawImage(img, p.Op)
}

//IsIdle returns true if the player is idle
func (p *Player) IsIdle() bool {
	return !p.WalkingAnimation.Animate && !p.ShootingAnimation.Animate
}

//Idle makes the player idle
func (p *Player) Idle() {
	p.IdleAnimation.Animate = true
	p.IdleAnimation.CurrentFrame++
}

//IsWalking returns true is the player is walking
func (p *Player) IsWalking() bool {
	return p.WalkingAnimation.Animate
}

//Walk make the player walk
func (p *Player) Walk(direction string) {
	switch direction {
	case "F": //Forward
		p.WalkingAnimation.Animate = true
		p.Coords.X += 2
		p.FacingFront = true
	case "B": //Backword
		p.WalkingAnimation.Animate = true
		p.Coords.X -= 2
		p.FacingFront = false
	case "U": //Up
		p.WalkingAnimation.Animate = true
		p.Coords.Y -= 2
	case "D": //Down
		p.WalkingAnimation.Animate = true
		p.Coords.Y += 2
	}

}

//Run make the player run
// func (p *Player) Run(direction string){
// 	switch direction{
// 	case "F": //Forward
// 		p.WalkingAnimation.Animate = true
// 		p.WalkingAnimation.CurrentFrame  +=2
// 		p.Coords.X +=2
// 		p.FacingFront = true
// 	case "B": //Backword
// 		p.WalkingAnimation.Animate = true
// 		p.WalkingAnimation.CurrentFrame +=2
// 		p.Coords.X -=2
// 		p.FacingFront = false
// 	case "U": //Up
// 		p.WalkingAnimation.Animate = true
// 		p.WalkingAnimation.CurrentFrame += 2
// 		p.Coords.Y -=2
// 	case "D": //Down
// 		p.WalkingAnimation.Animate = true
// 		p.WalkingAnimation.CurrentFrame += 2
// 		p.Coords.Y +=2
// 	}

// }

//Shoot shoots a bullet
func (p *Player) Shoot() {
	p.ShootingAnimation.Animate = true
	//p.Gun.Shoot()

}

//IsShooting returns true if the player is in shooting animation
func (p *Player) IsShooting() bool {
	return p.ShootingAnimation.Animate
}

//Collides telles whether the player is goona collide with te  object
func (p *Player) Collides(object image.Rectangle) bool {
	player := image.Rect(p.Coords.X-10, p.Coords.Y-10, p.Coords.X+10, p.Coords.Y+20)
	return player.Overlaps(object)
}

func (p *Player) String() string {
	return fmt.Sprintf("X: %v, Y: %v, Front: %v", p.Coords.X, p.Coords.Y, p.FacingFront)
}
