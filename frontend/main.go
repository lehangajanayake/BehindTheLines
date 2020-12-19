package main

import (
	"errors"
	"image"
	_ "image/png"
	"log"
	"strconv"

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
	g.Player.WalkingAnimation.Animate = false
	g.Player.IdleAnimation.Animate = false
	g.Player.ShootingAnimation.Animate = false
	
	if !g.Player.ShootingAnimation.Done{
		g.Player.ShootingAnimation.CurrentFrame ++
		g.Player.ShootingAnimation.Animate = true
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeyW){
		g.Player.WalkingAnimation.Animate = true
		g.Player.WalkingAnimation.CurrentFrame ++
		g.Player.Coords.Y --
	}else if ebiten.IsKeyPressed(ebiten.KeyS){
		g.Player.WalkingAnimation.Animate = true
		g.Player.WalkingAnimation.CurrentFrame ++
		g.Player.Coords.Y ++ 
		
	}
	if ebiten.IsKeyPressed(ebiten.KeyD){
		g.Player.WalkingAnimation.Animate = true
		g.Player.WalkingAnimation.CurrentFrame ++
		g.Player.Coords.X ++
		g.Player.FacingFront = true
	}else if ebiten.IsKeyPressed(ebiten.KeyA){
		g.Player.WalkingAnimation.Animate = true
		g.Player.WalkingAnimation.CurrentFrame ++
		g.Player.Coords.X --
		g.Player.FacingFront = false	
	
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft){
		g.Player.ShootingAnimation.Animate = true
		g.Player.ShootingAnimation.Done = false
		g.Player.ShootingAnimation.CurrentFrame ++
	}

	if !g.Player.WalkingAnimation.Animate && !g.Player.ShootingAnimation.Animate{
		g.Player.IdleAnimation.CurrentFrame ++
		g.Player.IdleAnimation.Animate = true
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
	
	if g.Player.WalkingAnimation.Animate {
		f := (g.Player.WalkingAnimation.CurrentFrame / 10 ) % g.Player.WalkingAnimation.FrameNum
		x, y := g.Player.WalkingAnimation.FrameWidth*f, g.Player.WalkingAnimation.StartY
		screen.DrawImage(g.Player.Img.SubImage(image.Rect(x, y, x + g.Player.WalkingAnimation.FrameWidth, y + g.Player.WalkingAnimation.FrameHeight)).(*ebiten.Image), g.Player.Op)
	}else if g.Player.IdleAnimation.Animate{
		f := (g.Player.IdleAnimation.CurrentFrame / 20) % g.Player.IdleAnimation.FrameNum
		x, y := g.Player.IdleAnimation.FrameWidth*f, g.Player.IdleAnimation.StartY
		screen.DrawImage(g.Player.Img.SubImage(image.Rect(x, y, x + g.Player.IdleAnimation.FrameWidth, y + g.Player.IdleAnimation.FrameHeight)).(*ebiten.Image), g.Player.Op)
		
	}else if g.Player.ShootingAnimation.Animate{
		f := (g.Player.ShootingAnimation.CurrentFrame / 5) % g.Player.ShootingAnimation.FrameNum
		x, y := g.Player.ShootingAnimation.FrameWidth*f, g.Player.ShootingAnimation.StartY
		screen.DrawImage(g.Player.Img.SubImage(image.Rect(x, y, x + g.Player.ShootingAnimation.FrameWidth, y + g.Player.ShootingAnimation.FrameHeight)).(*ebiten.Image), g.Player.Op)
		if g.Player.ShootingAnimation.CurrentFrame == g.Player.ShootingAnimation.FrameNum*5{
			g.Player.ShootingAnimation.Done = true
			g.Player.ShootingAnimation.CurrentFrame = 1
		}
	}	
	ebitenutil.DebugPrint(screen, strconv.Itoa(int(ebiten.CurrentFPS())))
}

// Layout lays the screen
func (g *Game) Layout(outsideWidth, outsideHeight int)(int, int){
	return ebiten.ScreenSizeInFullscreen()
}



func main(){
	img,_, err := ebitenutil.NewImageFromFile("assets/hero_spritesheet.png")
	if err != nil {
		log.Fatalf("Cannot load the assets err : %w", err)
	}
	ebiten.SetFullscreen(true)
	w, h := ebiten.ScreenSizeInFullscreen()
	g := Game{
		Player: models.Player{
			Img: img,
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
				Done: true,
			},
			FacingFront: true,
		},
		Frames: 0,
	}
	if err := ebiten.RunGame(&g); err != nil{
		log.Fatal(err)
	}
	
}