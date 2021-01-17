package main

import (
	"errors"
	"image"
	"sync"

	"github.com/lehangajanayake/MissionImposible/frontend/models"
	"github.com/lehangajanayake/MissionImposible/frontend/network"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
)

//Update updates  the game every tick. The game logic runs here
func (g *Game) Update() error {
	g.Camera.Zoom = 1
	g.Player.IdleAnimation.Animate, g.Player.WalkingAnimation.Animate, g.Player.ShootingAnimation.Animate = false, false, false
	g.Camera.Move(models.Coordinates{X: g.ScreenWidth / 2, Y: g.ScreenHeight / 2})

	if g.Player.IsShooting() {
		g.Player.Shoot()
		return nil
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.Player.Walk("F")
	} else if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.Player.Walk("B")
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.Player.Walk("U")

	} else if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.Player.Walk("D")

	}
	collide := make(chan bool, len(g.Map.TransparentObstacles)+len(g.Map.RayObjects))
	var wg sync.WaitGroup
	for _, obj := range g.Map.TransparentObstacles {
		wg.Add(1)
		go func(obj *tiled.Object, wg *sync.WaitGroup) {
			defer wg.Done()
			if g.Player.Collides(image.Rect(int(obj.X), int(obj.Y), int(obj.X+obj.Width), int(obj.Y+obj.Height))) {
				collide <- true
			}
		}(obj, &wg)

	}
	for _, obj := range g.Map.RayObjects {
		wg.Add(1)
		go func(obj *tiled.Object, wg *sync.WaitGroup) {
			defer wg.Done()
			if g.Player.Collides(image.Rect(int(obj.X), int(obj.Y), int(obj.X+obj.Width), int(obj.Y+obj.Height))) {
				collide <- true
			}
			for i, v := range g.Bullets {
				if v.Collides(image.Rect(int(obj.X), int(obj.Y), int(obj.X+obj.Width), int(obj.Y+obj.Height))) {
					g.Bullets = append(g.Bullets[:i], g.Bullets[i+1:]...)
				}
			}
		}(obj, &wg)

	}
	wg.Wait()
	close(collide)

	for v := range collide {
		if v {
			g.Player.Coords = g.Player.LastPos
			break
		}
	}

	if g.Player.LastPos != g.Player.Coords {
		coords := &network.Coordinates{X: g.Player.Coords.X, Y: g.Player.Coords.Y}
		g.Client.UpdatePlayerCoordsWrite <- coords.String()
		g.Player.LastPos = g.Player.Coords
	}

	if g.Player.LastFacing != g.Player.FacingFront {
		g.Player.LastFacing = g.Player.FacingFront
		if g.Player.FacingFront {
			g.Client.UpdatePlayerFacingWrite <- "true"
		} else {
			g.Client.UpdatePlayerFacingWrite <- "false"
		}

	}

	switch g.Player.LastAnimation {
	case g.Player.IdleAnimation.Name:
		if !g.Player.IdleAnimation.Animate {
			g.Player.LastAnimation = g.Player.IdleAnimation.Name
			g.Client.UpdatePlayerAnimationWrite <- g.Player.IdleAnimation.Name
		}
	case g.Player.WalkingAnimation.Name:
		if !g.Player.WalkingAnimation.Animate {
			g.Player.LastAnimation = g.Player.WalkingAnimation.Name
			g.Client.UpdatePlayerAnimationWrite <- g.Player.WalkingAnimation.Name
		}
	case g.Player.ShootingAnimation.Name:
		if !g.Player.ShootingAnimation.Animate {
			g.Player.LastAnimation = g.Player.ShootingAnimation.Name
			g.Client.UpdatePlayerAnimationWrite <- g.Player.ShootingAnimation.Name
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.Player.Shoot()
		bullet := new(models.Bullet)
		g.Bullets = append(g.Bullets, bullet.New(g.Player.Coords, g.Player.FacingFront))
	}

	if g.Player.IsIdle() {
		g.Player.Idle()
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		return errors.New("Game Exited by the user")
	}
	if g.Player.Coords.X > g.ScreenWidth/2 && g.Player.Coords.Y > g.ScreenHeight/2 {
		for _, obj := range g.Map.BlindSpots {
			if g.Player.Collides(image.Rect(int(obj.X), int(obj.Y), int(obj.X+obj.Width), int(obj.Y+obj.Height))) {
				g.Camera.Zoom = 0.7
				g.Camera.Move(g.Player.Coords)
			}
		}
		g.Camera.Move(g.Player.Coords)
	} else if g.Player.Coords.X > g.ScreenWidth/2 {
		g.Camera.Move(models.Coordinates{X: g.Player.Coords.X, Y: g.ScreenHeight / 2})
	} else if g.Player.Coords.Y > g.ScreenHeight/2 {
		g.Camera.Move(models.Coordinates{X: g.ScreenWidth / 2, Y: g.Player.Coords.Y})
	}

	return nil
}
