package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	//"log"

	//"github.com/lehangajanayake/MissionImposible/frontend/models"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lehangajanayake/MissionImposible/frontend/ray"
)

// Draw draws to the screen every update
func (g *Game) Draw(screen *ebiten.Image) {

	g.Camera.View.DrawImage(g.Map.World, g.Map.Op)

	g.Player.Op.GeoM.Reset()
	g.Player.Op.GeoM.Translate(-float64(g.Player.Animations["Idle"].FrameWidth/2), -float64(g.Player.Animations["Idle"].FrameHeight/2)) //,ake the axis of the player in teh middle instead of the upper left conner
	if !g.Player.FacingFront {
		g.Player.Op.GeoM.Scale(-1, 1)
	}
	g.Player.Op.GeoM.Translate(float64(g.Player.Coords.X), float64(g.Player.Coords.Y))
	g.Player.Animate(g.Camera.View)

	for _, v := range g.Client.Players {
		op := &ebiten.DrawImageOptions{}
		
		if v.Guard{
			op.GeoM.Translate(-float64(g.GuardAnimation["Idle"].FrameWidth/2), -float64(g.GuardAnimation["Idle"].FrameHeight/2)) //,ake the axis of the player in teh middle instead of the upper left conner
			op.GeoM.Translate(float64(v.Coords.X), float64(v.Coords.Y))
			if !v.FacingFront {
				op.GeoM.Scale(-1, 1)
			}
			switch v.State {
			case "Idle":
				if _, ok := g.GuardAnimation["Idle"]; !ok {
					log.Fatal("Animation cant be found")
				}
				g.GuardAnimation["Idle"].CurrentFrame++
				f := (g.GuardAnimation["Idle"].CurrentFrame / 20) % g.GuardAnimation["Idle"].FrameNum
				x, y := g.GuardAnimation["Idle"].FrameWidth*f, g.GuardAnimation["Idle"].StartY
				g.Camera.View.DrawImage( g.GuardAnimation["Idle"].Img.SubImage(image.Rect(x, y, x+g.GuardAnimation["Idle"].FrameWidth, y+g.GuardAnimation["Idle"].FrameHeight)).(*ebiten.Image), op)
				if g.GuardAnimation["Idle"].CurrentFrame == g.GuardAnimation["Idle"].FrameNum*20 { //done ideling
					g.GuardAnimation["Idle"].Reset()
				}
			case "Walking":
				if _, ok := g.GuardAnimation["Walking"]; !ok {
					log.Fatal("Animation cant be found")
				}
				g.GuardAnimation["Walking"].CurrentFrame++
				f := (g.GuardAnimation["Walking"].CurrentFrame / 10) % g.GuardAnimation["Walking"].FrameNum
				x, y := g.GuardAnimation["Walking"].FrameWidth*f, g.GuardAnimation["Walking"].StartY
				g.Camera.View.DrawImage( g.GuardAnimation["Walking"].Img.SubImage(image.Rect(x, y, x+g.GuardAnimation["Walking"].FrameWidth, y+g.GuardAnimation["Walking"].FrameHeight)).(*ebiten.Image), op)
				if g.GuardAnimation["Walking"].CurrentFrame == g.GuardAnimation["Walking"].FrameNum*10 { //done shooting
					g.GuardAnimation["Walking"].Reset()
				}

			case "Attacking":
				if _, ok := g.GuardAnimation["Attacking"]; !ok {
					log.Fatal("Animation cant be found")
				}
				g.GuardAnimation["Attack"].CurrentFrame++
				f := (g.GuardAnimation["Attack"].CurrentFrame / 3) % g.GuardAnimation["Attack"].FrameNum
				x, y := g.GuardAnimation["Attack"].FrameWidth*f, g.GuardAnimation["Attack"].StartY
				g.Camera.View.DrawImage( g.GuardAnimation["Attack"].Img.SubImage(image.Rect(x, y, x+g.GuardAnimation["Shooting"].FrameWidth, y+g.GuardAnimation["Shooting"].FrameHeight)).(*ebiten.Image), op)
			default:
				log.Println("wft")

			}
		}else{
			op.GeoM.Translate(-float64(g.GuardAnimation["Idle"].FrameWidth/2), -float64(g.GuardAnimation["Idle"].FrameHeight/2)) //,ake the axis of the player in teh middle instead of the upper left conner
			op.GeoM.Translate(float64(v.Coords.X), float64(v.Coords.Y))
			if !v.FacingFront {
				op.GeoM.Scale(-1, 1)
			}
			switch v.State {
			case "Idle":
				if _, ok := g.NinjaAnimation["Idle"]; !ok {
					log.Fatal("Animation cant be found")
				}
				g.NinjaAnimation["Idle"].CurrentFrame++
				f := (g.NinjaAnimation["Idle"].CurrentFrame / 20) % g.NinjaAnimation["Idle"].FrameNum
				x, y := g.NinjaAnimation["Idle"].FrameWidth*f, g.NinjaAnimation["Idle"].StartY
				g.Camera.View.DrawImage( g.NinjaAnimation["Idle"].Img.SubImage(image.Rect(x, y, x+g.NinjaAnimation["Idle"].FrameWidth, y+g.NinjaAnimation["Idle"].FrameHeight)).(*ebiten.Image), op)
				if g.NinjaAnimation["Idle"].CurrentFrame == g.NinjaAnimation["Idle"].FrameNum*20 { //done ideling
					g.NinjaAnimation["Idle"].Reset()
				}
			case "Walking":
				if _, ok := g.NinjaAnimation["Walking"]; !ok {
					log.Fatal("Animation cant be found")
				}
				g.NinjaAnimation["Walking"].CurrentFrame++
				f := (g.NinjaAnimation["Walking"].CurrentFrame / 10) % g.NinjaAnimation["Walking"].FrameNum
				x, y := g.NinjaAnimation["Walking"].FrameWidth*f, g.NinjaAnimation["Walking"].StartY
				g.Camera.View.DrawImage( g.NinjaAnimation["Walking"].Img.SubImage(image.Rect(x, y, x+g.NinjaAnimation["Walking"].FrameWidth, y+g.NinjaAnimation["Walking"].FrameHeight)).(*ebiten.Image), op)
				if g.NinjaAnimation["Walking"].CurrentFrame == g.NinjaAnimation["Walking"].FrameNum*10 { //done shooting
					g.NinjaAnimation["Walking"].Reset()
				}

			case "Attacking":
				if _, ok := g.NinjaAnimation["Attacking"]; !ok {
					log.Fatal("Animation cant be found")
				}
				g.NinjaAnimation["Attack"].CurrentFrame++
				f := (g.NinjaAnimation["Attack"].CurrentFrame / 3) % g.NinjaAnimation["Attack"].FrameNum
				x, y := g.NinjaAnimation["Attack"].FrameWidth*f, g.NinjaAnimation["Attack"].StartY
				g.Camera.View.DrawImage( g.NinjaAnimation["Attack"].Img.SubImage(image.Rect(x, y, x+g.NinjaAnimation["Shooting"].FrameWidth, y+g.NinjaAnimation["Shooting"].FrameHeight)).(*ebiten.Image), op)

			}
		}
		
	}

	if len(g.Bullets) != 0 {
		for _, v := range g.Bullets {
			v.Move()
			v.Op.GeoM.Reset()
			v.Op.GeoM.Scale(4, 2)
			v.Op.GeoM.Translate(float64(v.Coords.X), float64(v.Coords.Y))
			v.Render(g.Camera.View, g.BulletImg)

		}
	}

	//render the trees
	g.Camera.View.DrawImage(g.Map.Trees, g.Map.Op)
	//render teh shadow
	g.ShadowImg.Fill(color.Black)
	//log.Println(g.RayObjects)
	rays := ray.Cast(float64(g.Player.Coords.X), float64(g.Player.Coords.Y), g.RayObjects)

	// Subtract ray triangles from shadow
	opt := &ebiten.DrawTrianglesOptions{}
	opt.Address = ebiten.AddressRepeat
	opt.CompositeMode = ebiten.CompositeModeSourceOut
	for i, line := range rays {
		nextLine := rays[(i+1)%len(rays)]

		// Draw triangle of area between rays
		v := ray.Vertices(float64(g.Player.Coords.X), float64(g.Player.Coords.Y), nextLine.X2, nextLine.Y2, line.X2, line.Y2)
		g.ShadowImg.DrawTriangles(v, []uint16{0, 1, 2}, g.TriangleImg, opt)
	}
	op := &ebiten.DrawImageOptions{}
	op.ColorM.Scale(1, 1, 1, 0.7)
	g.Camera.View.DrawImage(g.ShadowImg, op)
	//render using the camera
	g.Camera.Render(screen)
	//clear the camera
	g.Camera.View.Clear()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Bullets Left: %v , Current TPS: %0.2f, Current FPS: %0.2f", g.Player.Gun.Bullets, ebiten.CurrentTPS(), ebiten.CurrentFPS()))
	//ebitenutil.DebugPrint(screen, g.Player.String())
}
