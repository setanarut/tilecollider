package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/setanarut/tilecollider"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

var (
	PlayerX = 140.
	PlayerY = 130.
	PlayerW = 140.
	PlayerH = 20.

	TileMap = [][]uint8{
		{0, 0, 0, 0, 0, 0, 9, 1},
		{0, 0, 0, 0, 0, 6, 0, 1},
		{4, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 8, 0, 8, 3, 1},
		{2, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0, 0, 0, 0},
		{1, 0, 4, 0, 5, 0, 0, 0},
		{1, 4, 2, 8, 1, 88, 13, 1},
	}

	collider = tilecollider.NewTileCollider(TileMap, screenWidth/8, screenHeight/8)
)

type Game struct {
}

func (g *Game) Update() error {

	// Get axis input
	velX, velY := Axis()
	velY *= 4
	velX *= 4

	// Collide with tiles
	deltaX, deltaY := collider.Collide(
		PlayerX,
		PlayerY,
		PlayerW,
		PlayerH,
		velX,
		velY,
		nil,
	)

	// Update player position
	PlayerX += deltaX
	PlayerY += deltaY

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	// Draw tiles
	for y := 0; y < len(TileMap); y++ {
		for x := 0; x < len(TileMap[y]); x++ {
			if TileMap[y][x] != 0 {
				vector.DrawFilledRect(screen,
					float32(x*collider.TileSize[0]),
					float32(y*collider.TileSize[1]),
					float32(collider.TileSize[0]),
					float32(collider.TileSize[1]),
					color.Gray{Y: 128},
					true)
			}
		}
	}

	// Draw player
	vector.DrawFilledRect(screen,
		float32(PlayerX),
		float32(PlayerY),
		float32(PlayerW),
		float32(PlayerH),
		color.RGBA{47, 36, 254, 255},
		false)

	// Print collisions to the screen
	for i, c := range collider.Collisions {
		ebitenutil.DebugPrintAt(screen,
			fmt.Sprintf(
				"Tile ID: %d, Tile Coords: %v, Collision Normal: %v",
				c.TileID,
				c.TileCoords,
				c.Normal,
			), 20, 20+(i*20))
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}

func Axis() (axisX, axisY float64) {
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		axisY -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		axisY += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		axisX -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		axisX += 1
	}
	return axisX, axisY
}
