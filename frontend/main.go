package main

import (
	"errors"
	"fmt"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lehangajanayake/MissionImposible/frontend/models"
)

//Game the game
type Game struct{
	Player models.Player
	Frames int
}

//Update updates  the game 
func (g *Game) Update()error{
	g.Player.IdleAnimation.Reset()
	g.Player.WalkingAnimation.Reset()
	
	if g.Player.IsShooting(){
		g.Player.Shoot(true)
		g.Player.Gun.Bullet.Coords.Y = g.Player.Coords.Y
		g.Player.Gun.Bullet.Coords.X = g.Player.Coords.X
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeyD){
		g.Player.Walk("F")
	}else if ebiten.IsKeyPressed(ebiten.KeyA){
		g.Player.Walk("B")
	}
	if ebiten.IsKeyPressed(ebiten.KeyW){
		g.Player.Walk("U")
	}else if ebiten.IsKeyPressed(ebiten.KeyS){
		g.Player.Walk("D")
	}
	

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft){
		g.Player.Shoot(false)
	}

	if g.Player.IsIdle(){
		g.Player.Idle()
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ){
		return errors.New("Game Exited by the user")
	}

	return nil
}
// Draw draws to the screen every update
func (g *Game) Draw(screen *ebiten.Image){
	g.Player.Op.GeoM.Reset()
	g.Player.Op.GeoM.Translate(-float64(g.Player.WalkingAnimation.FrameWidth/2), -float64(g.Player.WalkingAnimation.FrameHeight/2)) //,ake the axiz of the player in teh midlle instead of the upper left conner
	if g.Player.FacingFront{
		g.Player.Op.GeoM.Scale(1,1)
	}else{
		g.Player.Op.GeoM.Scale(-1,1)
	}
	g.Player.Op.GeoM.Translate(float64(g.Player.Coords.X),float64(g.Player.Coords.Y))
	
	if g.Player.IsWalking(){
		f := (g.Player.WalkingAnimation.CurrentFrame / 10 ) % g.Player.WalkingAnimation.FrameNum
		x, y := g.Player.WalkingAnimation.FrameWidth*f, g.Player.WalkingAnimation.StartY
		g.Player.Render(screen, g.Player.Img.SubImage(image.Rect(x, y, x + g.Player.WalkingAnimation.FrameWidth, y + g.Player.WalkingAnimation.FrameHeight)).(*ebiten.Image))
	}else if g.Player.IsIdle(){
		f := (g.Player.IdleAnimation.CurrentFrame / 20) % g.Player.IdleAnimation.FrameNum
		x, y := g.Player.IdleAnimation.FrameWidth*f, g.Player.IdleAnimation.StartY
		g.Player.Render(screen, g.Player.Img.SubImage(image.Rect(x, y, x + g.Player.WalkingAnimation.FrameWidth, y + g.Player.WalkingAnimation.FrameHeight)).(*ebiten.Image))
		
	}else if g.Player.IsShooting(){
		f := (g.Player.ShootingAnimation.CurrentFrame / 5) % g.Player.ShootingAnimation.FrameNum
		x, y := g.Player.ShootingAnimation.FrameWidth*f, g.Player.ShootingAnimation.StartY
		g.Player.Render(screen, g.Player.Img.SubImage(image.Rect(x, y, x + g.Player.WalkingAnimation.FrameWidth, y + g.Player.WalkingAnimation.FrameHeight)).(*ebiten.Image))
		if g.Player.ShootingAnimation.CurrentFrame == g.Player.ShootingAnimation.FrameNum*5{ //done shooting
			g.Player.ShootingAnimation.Reset()
			g.Player.ShootingAnimation.CurrentFrame = 1
		}
	}
	if g.Player.Gun.Bullet.IsHit(){
		g.Player.Gun.Bullet.Move()
		g.Player.Gun.Bullet.Op.GeoM.Reset()
		//g.Player.Gun.Bullet.Op.GeoM.Translate(-float64(0.5), -float64(05)) //make the point centered
		//g.Player.Gun.Bullet.Op.GeoM.Scale(10,10)
		//println(g.Player.Gun.Bullet.Coords.X)
		g.Player.Gun.Bullet.Op.GeoM.Translate(float64(g.Player.Gun.Bullet.Coords.X), float64(g.Player.Gun.Bullet.Coords.Y))
		screen.DrawImage(g.Player.Gun.Bullet.Img, g.Player.Gun.Bullet.Op)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Bullets Left: %v", g.Player.Gun.Bullets))
}

// Layout lays the screen
func (g *Game) Layout(outsideWidth, outsideHeight int)(int, int){
	return ebiten.ScreenSizeInFullscreen()
}



func main(){
	player,_, err := ebitenutil.NewImageFromFile("assets/hero_spritesheet.png")
	if err != nil {
		log.Fatalf("Cannot load the assets err : %w", err)
	}
	bullet,_, err := ebitenutil.NewImageFromFile("assets/bullet.png")
	if err != nil {
		log.Fatalf("Cannot load the assets err : %w", err)
	}
	ebiten.SetFullscreen(true)
	w, h := ebiten.ScreenSizeInFullscreen()
	g := Game{
		Player: models.Player{
			Img: player,
			Op: &ebiten.DrawImageOptions{},
			Coords: models.Coordinates{
				X: w/2,
				Y: h/2,
			},
			WalkingAnimation: models.Animation{
				StartX: 0,
				StartY: 100,
				FrameNum: 6	,
				FrameWidth: 80,
				FrameHeight: 80,
				Animate: false,
			},
			IdleAnimation: models.Animation{
				StartX: 0,
				StartY: 0,
				FrameNum: 1,
				FrameHeight: 80,
				FrameWidth: 80,
				Animate: false,

			},
			ShootingAnimation: models.Animation{
				StartX: 0,
				StartY: 0,
				FrameNum: 8,
				FrameHeight: 80,
				FrameWidth: 80,
				Animate: false,
			},
			Gun: models.Gun{
				Bullets: 60,
				Bullet: models.Bullet{
					Img: bullet,
					Op: &ebiten.DrawImageOptions{},
					Hit: false,
				},
			},
			FacingFront: true,
		},
		Frames: 0,
	}
	if err := ebiten.RunGame(&g); err != nil{
		log.Fatal(err)
	}
	
}