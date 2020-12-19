package models

import (

	"github.com/hajimehoshi/ebiten/v2"
)

//Player model player
type Player struct{
	Img *ebiten.Image
	Coords Coordinates
	FacingFront bool
	Op *ebiten.DrawImageOptions
	Gun  Gun
	WalkingAnimation Animation
	IdleAnimation Animation
	ShootingAnimation Animation
}

//Render reders a plyer in the screen
func (p *Player) Render(screen, img *ebiten.Image){
	screen.DrawImage(img, p.Op)
}

//IsIdle returns true if the player is idle
func (p *Player) IsIdle()bool{
	return !p.WalkingAnimation.Animate && !p.ShootingAnimation.Animate 
}

//Idle makes the player idle
func (p *Player) Idle(){
	p.IdleAnimation.Animate = true
	p.IdleAnimation.CurrentFrame ++
}

//IsWalking returns true is the player is walking
func (p *Player) IsWalking()bool{
	return p.WalkingAnimation.Animate
}
//Walk make the player walk
func (p *Player) Walk(direction string){
	switch direction{
	case "F": //Forward
		p.WalkingAnimation.Animate = true
		p.WalkingAnimation.CurrentFrame ++
		p.Coords.X ++
		p.FacingFront = true
	case "B": //Backword
		p.WalkingAnimation.Animate = true
		p.WalkingAnimation.CurrentFrame ++
		p.Coords.X --
		p.FacingFront = false	
	case "U": //Up
		p.WalkingAnimation.Animate = true
		p.WalkingAnimation.CurrentFrame ++
		p.Coords.Y --
	case "D": //Down
		p.WalkingAnimation.Animate = true
		p.WalkingAnimation.CurrentFrame ++
		p.Coords.Y ++ 
	}
	
}

//Shoot shoots a bullet
func (p *Player) Shoot(shot bool){
	if !shot{
		p.Gun.Bullet.FacingFront = p.FacingFront
		p.Gun.Shoot()
	}
	p.ShootingAnimation.Animate = true
	p.ShootingAnimation.CurrentFrame ++
	
}

//IsShooting returns true if the player is in shooting animation
func (p *Player) IsShooting()bool{
	return p.ShootingAnimation.Animate
}
