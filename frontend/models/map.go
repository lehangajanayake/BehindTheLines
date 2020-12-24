package models

import (
	"bytes"
	"image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lafriks/go-tiled"
	"github.com/lafriks/go-tiled/render"
)

//Map model map
type Map struct {
	tile *tiled.Map
	Op *ebiten.DrawImageOptions
	World *ebiten.Image
	Trees *ebiten.Image
	Obstacles []*tiled.Object	
}


//LoadMap loads the map
func (m *Map) LoadMap(path string)error{
	var err error
	m.tile, err = tiled.LoadFromFile(path)
	if err != nil {
		return err
	}
	render1, err := render.NewRenderer(m.tile)
	var buff []byte
	buffer := bytes.NewBuffer(buff)
	err = render1.RenderLayer(0)
	if err != nil {
		return err
	}
	err = render1.SaveAsPng(buffer)
	if err != nil{
		return err
	}
	img, err := png.Decode(buffer)
	if err != nil {
		return err
	}
	m.World = ebiten.NewImageFromImage(img)
	buffer.Reset()
	render2, err := render.NewRenderer(m.tile)
	err = render2.RenderLayer(1)
	if err != nil {
		return err
	}
	err = render2.SaveAsPng(buffer)
	if err != nil{
		return err
	}
	img, err = png.Decode(buffer)
	if err != nil {
		return err
	}
	m.Trees = ebiten.NewImageFromImage(img)
	return nil
}

//LoadObstacles loads the obstacles
func (m *Map) LoadObstacles(index int){
	m.Obstacles = m.tile.ObjectGroups[index].Objects
	
}
