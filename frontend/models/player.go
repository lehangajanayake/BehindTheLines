package models

import (
	"fmt"
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

//Player model player
type Player struct {
	Img                                                *ebiten.Image
	Coords, LastPos                                    Coordinates
	FacingFront, LastFacing, Guard                     bool
	LastAnimation                                      string
	Op                                                 *ebiten.DrawImageOptions
	Gun                                                Gun
	//Animations["Walking"], IdleAnimation, ShootingAnimation Animation
	State string
	Animations map[string]*Animation
}

//Render reders a plyer in the screen
func (p *Player) Render(screen, img *ebiten.Image) {
	screen.DrawImage(img, p.Op)
}


//Walk make the player walk
func (p *Player) Walk(direction string) {
	p.State = "Walking"
	switch direction {
	case "F": //Forward
		p.Coords.X += 2
		p.FacingFront = true
	case "B": //Backword
		p.Coords.X -= 2
		p.FacingFront = false
	case "U": //Up
		p.Coords.Y -= 2
	case "D": //Down
		p.Coords.Y += 2
	}

}

//Run make the player run
// func (p *Player) Run(direction string){
// 	switch direction{
// 	case "F": //Forward
// 		p.State =
// 		p.Animations["Walking"].CurrentFrame  +=2
// 		p.Coords.X +=2
// 		p.FacingFront = true
// 	case "B": //Backword
// 		p.State =
// 		p.Animations["Walking"].CurrentFrame +=2
// 		p.Coords.X -=2
// 		p.FacingFront = false
// 	case "U": //Up
// 		p.State =
// 		p.Animations["Walking"].CurrentFrame += 2
// 		p.Coords.Y -=2
// 	case "D": //Down
// 		p.State =
// 		p.Animations["Walking"].CurrentFrame += 2
// 		p.Coords.Y +=2
// 	}

// }

//Shoot shoots a bullet
func (p *Player) Shoot() {
	p.State = "Shooting"
	//p.Gun.Shoot()

}

//Collides telles whether the player is goona collide with te  object
func (p *Player) Collides(object image.Rectangle) bool {
	player := image.Rect(p.Coords.X-10, p.Coords.Y-10, p.Coords.X+10, p.Coords.Y+20)
	return player.Overlaps(object)
}

func (p *Player) String() string {
	return fmt.Sprintf("X: %v, Y: %v, Front: %v", p.Coords.X, p.Coords.Y, p.FacingFront)
}


func (p *Player) Animate(can *ebiten.Image){
	switch p.State{
	case "Idle":
		if _, ok := p.Animations["Idle"]; !ok {
			log.Fatal("Animation cant be found")
		}
		p.Animations["Idle"].CurrentFrame++
		f := (p.Animations["Idle"].CurrentFrame / 20) % p.Animations["Idle"].FrameNum
		x, y := p.Animations["Idle"].FrameWidth*f, p.Animations["Idle"].StartY
		p.Render(can, p.Animations["Idle"].Img.SubImage(image.Rect(x, y, x+p.Animations["Idle"].FrameWidth, y+p.Animations["Idle"].FrameHeight)).(*ebiten.Image))
		if p.Animations["Idle"].CurrentFrame == p.Animations["Idle"].FrameNum*20 { //done ideling
			p.Animations["Idle"].Reset()
		}
	case "Walking":
		if _, ok := p.Animations["Walking"]; !ok {
			log.Fatal("Animation cant be found")
		}
		p.Animations["Walking"].CurrentFrame++
		f := (p.Animations["Walking"].CurrentFrame / 10) % p.Animations["Walking"].FrameNum
		x, y := p.Animations["Walking"].FrameWidth*f, p.Animations["Walking"].StartY
		p.Render(can, p.Animations["Walking"].Img.SubImage(image.Rect(x, y, x+p.Animations["Walking"].FrameWidth, y+p.Animations["Walking"].FrameHeight)).(*ebiten.Image))
		if p.Animations["Walking"].CurrentFrame == p.Animations["Walking"].FrameNum*10 { //done shooting
			p.Animations["Walking"].Reset()
		}
		
	case "Attacking":
		if _, ok := p.Animations["Attacking"]; !ok {
			log.Fatal("Animation cant be found")
		}
		p.Animations["Attack"].CurrentFrame++
		f := (p.Animations["Attack"].CurrentFrame / 3) % p.Animations["Attack"].FrameNum
		x, y := p.Animations["Attack"].FrameWidth*f, p.Animations["Attack"].StartY
		p.Render(can, p.Animations["Attack"].Img.SubImage(image.Rect(x, y, x+p.Animations["Shooting"].FrameWidth, y+p.Animations["Shooting"].FrameHeight)).(*ebiten.Image))

		if p.Animations["Attack"].CurrentFrame == p.Animations["Attack"].FrameNum*3 { //done shooting
			p.Animations["Attack"].Reset()
			p.Gun.Shoot()
		}
	}
}