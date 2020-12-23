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
	Img *ebiten.Image
	Obstacles []*tiled.Object	
}


//LoadMap loads the map
func (m *Map) LoadMap(path string)error{
	var err error
	m.tile, err = tiled.LoadFromFile(path)
	if err != nil {
		return err
	}
	render, err := render.NewRenderer(m.tile)
	var buff []byte
	buffer := bytes.NewBuffer(buff)
	err = render.RenderVisibleLayers()
	if err != nil {
		return err
	}
	err = render.SaveAsPng(buffer)
	if err != nil{
		return err
	}
	img, err := png.Decode(buffer)
	if err != nil {
		return err
	}
	m.Img = ebiten.NewImageFromImage(img)
	return nil
}

//LoadObstacles loads the obstacles
func (m *Map) LoadObstacles(index int){
	m.Obstacles = m.tile.ObjectGroups[index].Objects
}
