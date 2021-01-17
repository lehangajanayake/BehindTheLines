package main

import (
	"fmt"
	"image"
	"image/color"

	//"githib.com/lehangajanayake/MissionImposible/frontend/models"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/lehangajanayake/MissionImposible/frontend/ray"
)

// Draw draws to the screen every update
func (g *Game) Draw(screen *ebiten.Image) {

	g.Camera.View.DrawImage(g.Map.World, g.Map.Op)

	g.Player.Op.GeoM.Reset()
	g.Player.Op.GeoM.Translate(-float64(g.Player.WalkingAnimation.FrameWidth/2), -float64(g.Player.WalkingAnimation.FrameHeight/2)) //,ake the axis of the player in teh middle instead of the upper left conner
	if g.Player.FacingFront {
		g.Player.Op.GeoM.Scale(0.5, 0.5)
	} else {
		g.Player.Op.GeoM.Scale(-0.5, 0.5)
	}
	g.Player.Op.GeoM.Translate(float64(g.Player.Coords.X), float64(g.Player.Coords.Y))

	if g.Player.IsWalking() {
		g.Player.WalkingAnimation.CurrentFrame++
		f := (g.Player.WalkingAnimation.CurrentFrame / 10) % g.Player.WalkingAnimation.FrameNum
		x, y := g.Player.WalkingAnimation.FrameWidth*f, g.Player.WalkingAnimation.StartY
		g.Player.Render(g.Camera.View, g.Player.Img.SubImage(image.Rect(x, y, x+g.Player.WalkingAnimation.FrameWidth, y+g.Player.WalkingAnimation.FrameHeight)).(*ebiten.Image))
		if g.Player.WalkingAnimation.CurrentFrame == g.Player.WalkingAnimation.FrameNum*10 || !g.Player.IsWalking() { //done shooting
			g.Player.WalkingAnimation.Reset()
		}

	} else if g.Player.IsIdle() {
		g.Player.IdleAnimation.CurrentFrame++
		f := (g.Player.IdleAnimation.CurrentFrame / 20) % g.Player.IdleAnimation.FrameNum
		x, y := g.Player.IdleAnimation.FrameWidth*f, g.Player.IdleAnimation.StartY
		g.Player.Render(g.Camera.View, g.Player.Img.SubImage(image.Rect(x, y, x+g.Player.WalkingAnimation.FrameWidth, y+g.Player.WalkingAnimation.FrameHeight)).(*ebiten.Image))
		if g.Player.IdleAnimation.CurrentFrame == g.Player.IdleAnimation.FrameNum*20 || !g.Player.IsIdle() { //done ideling
			g.Player.IdleAnimation.Reset()
		}

	} else if g.Player.IsShooting() {
		g.Player.IdleAnimation.CurrentFrame++
		f := (g.Player.ShootingAnimation.CurrentFrame / 3) % g.Player.ShootingAnimation.FrameNum
		x, y := g.Player.ShootingAnimation.FrameWidth*f, g.Player.ShootingAnimation.StartY
		g.Player.Render(g.Camera.View, g.Player.Img.SubImage(image.Rect(x, y, x+g.Player.WalkingAnimation.FrameWidth, y+g.Player.WalkingAnimation.FrameHeight)).(*ebiten.Image))

		if g.Player.ShootingAnimation.CurrentFrame == g.Player.ShootingAnimation.FrameNum*3 { //done shooting
			g.Player.ShootingAnimation.Reset()
			g.Player.Gun.Shoot()
		}
	}
	for _, v := range g.Client.Players {
		v.CurrentFrame++
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(g.Player.WalkingAnimation.FrameWidth/2), -float64(g.Player.WalkingAnimation.FrameHeight/2)) //,ake the axis of the player in teh middle instead of the upper left conner
		if v.FacingFront {
			op.GeoM.Scale(0.5, 0.5)
		} else {
			op.GeoM.Scale(-0.5, 0.5)
		}
		op.GeoM.Translate(float64(v.Coords.X), float64(v.Coords.Y))
		switch v.Animation {
		case g.Player.IdleAnimation.Name:
			f := (v.CurrentFrame / 20) % g.Player.IdleAnimation.FrameNum
			x, y := g.Player.IdleAnimation.FrameWidth*f, g.Player.IdleAnimation.StartY
			g.Camera.View.DrawImage(g.Player.Img.SubImage(image.Rect(x, y, x+g.Player.IdleAnimation.FrameWidth, y+g.Player.IdleAnimation.FrameHeight)).(*ebiten.Image), op)
		case g.Player.WalkingAnimation.Name:
			f := (v.CurrentFrame / 10) % g.Player.WalkingAnimation.FrameNum
			x, y := g.Player.WalkingAnimation.FrameWidth*f, g.Player.WalkingAnimation.StartY
			g.Camera.View.DrawImage(g.Player.Img.SubImage(image.Rect(x, y, x+g.Player.WalkingAnimation.FrameWidth, y+g.Player.WalkingAnimation.FrameHeight)).(*ebiten.Image), op)
		case g.Player.ShootingAnimation.Name:
			f := (v.CurrentFrame / 3) % g.Player.ShootingAnimation.FrameNum
			x, y := g.Player.ShootingAnimation.FrameWidth*f, g.Player.ShootingAnimation.StartY
			g.Camera.View.DrawImage(g.Player.Img.SubImage(image.Rect(x, y, x+g.Player.ShootingAnimation.FrameWidth, y+g.Player.ShootingAnimation.FrameHeight)).(*ebiten.Image), op)
		}
		//g.Camera.View.DrawImage(g.Player.Img.SubImage(image.Rect(0, 0, g.Player.WalkingAnimation.FrameWidth,g.Player.WalkingAnimation.FrameHeight)).(*ebiten.Image), op)
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