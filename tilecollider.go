// Package tilecollide provides collision detection for rectangles and 2D tilemaps.
package tilecollider

import (
	"math"

	"golang.org/x/exp/constraints"
)

// Integer is a constraint that matches any integer type.
type Integer = constraints.Integer

// CollisionInfo stores information about a collision with a tile
type CollisionInfo[T Integer] struct {
	TileID     T      // ID of the collided tile
	TileCoords [2]int // X,Y coordinates of the tile in the tilemap
	Normal     [2]int // Normal vector of the collision (-1/0/1)
}

// Collider handles collision detection between rectangles and a 2D tilemap
type Collider[T Integer] struct {
	Collisions     []CollisionInfo[T] // List of collisions from last check
	TileSize       [2]int             // Width and height of tiles
	TileMap        [][]T              // 2D grid of tile IDs
	NonSolidTileID T                  // Sets the ID of non-solid tiles. Defaults to 0.
}

// NewCollider creates a new tile collider with the given tilemap and tile dimensions
func NewCollider[T Integer](tileMap [][]T, tileWidth, tileHeight int) *Collider[T] {
	return &Collider[T]{
		TileMap:  tileMap,
		TileSize: [2]int{tileWidth, tileHeight},
	}
}

// CollisionCallback is called when collisions occur, receiving collision info and final movement
type CollisionCallback[T Integer] func([]CollisionInfo[T], float64, float64)

// Collide checks for collisions when moving a rectangle and returns the allowed movement
func (c *Collider[T]) Collide(rectX, rectY, rectW, rectH, moveX, moveY float64, onCollide CollisionCallback[T]) (float64, float64) {

	c.Collisions = c.Collisions[:0]

	if moveX == 0 && moveY == 0 {
		return moveX, moveY
	}

	if math.Abs(moveX) > math.Abs(moveY) {
		if moveX != 0 {
			moveX = c.collideX(rectX, rectY, rectW, rectH, moveX)
		}
		if moveY != 0 {
			moveY = c.collideY(rectX+moveX, rectY, rectW, rectH, moveY)
		}
	} else {
		if moveY != 0 {
			moveY = c.collideY(rectX, rectY, rectW, rectH, moveY)
		}
		if moveX != 0 {
			moveX = c.collideX(rectX, rectY+moveY, rectW, rectH, moveX)
		}
	}

	if onCollide != nil {
		onCollide(c.Collisions, moveX, moveY)
	}

	return moveX, moveY
}

// collideX checks for collisions along the X axis and returns the allowed X movement
func (c *Collider[T]) collideX(rectX, rectY, rectW, rectH, moveX float64) float64 {

	checkLimit := max(1, int(math.Ceil(math.Abs(moveX)/float64(c.TileSize[0])))+1)

	playerTop := int(math.Floor(rectY / float64(c.TileSize[1])))
	playerBottom := int(math.Ceil((rectY+rectH)/float64(c.TileSize[1]))) - 1

	if moveX > 0 {
		startX := int(math.Floor((rectX + rectW) / float64(c.TileSize[0])))
		endX := startX + checkLimit
		endX = min(endX, len(c.TileMap[0]))

		for y := playerTop; y <= playerBottom; y++ {
			if y < 0 || y >= len(c.TileMap) {
				continue
			}
			for x := startX; x < endX; x++ {
				if x < 0 || x >= len(c.TileMap[0]) {
					continue
				}
				if c.TileMap[y][x] != c.NonSolidTileID {
					tileLeft := float64(x * c.TileSize[0])
					collision := tileLeft - (rectX + rectW)
					if collision <= moveX {
						moveX = collision
						c.Collisions = append(c.Collisions, CollisionInfo[T]{
							TileID:     c.TileMap[y][x],
							TileCoords: [2]int{x, y},
							Normal:     [2]int{-1, 0},
						})
					}
				}
			}
		}
	}

	if moveX < 0 {
		endX := int(math.Floor(rectX / float64(c.TileSize[0])))
		startX := endX - checkLimit
		startX = max(startX, 0)

		for y := playerTop; y <= playerBottom; y++ {
			if y < 0 || y >= len(c.TileMap) {
				continue
			}
			for x := startX; x <= endX; x++ {
				if x < 0 || x >= len(c.TileMap[0]) {
					continue
				}
				if c.TileMap[y][x] != c.NonSolidTileID {
					tileRight := float64((x + 1) * c.TileSize[0])
					collision := tileRight - rectX
					if collision >= moveX {
						moveX = collision
						c.Collisions = append(c.Collisions, CollisionInfo[T]{
							TileID:     c.TileMap[y][x],
							TileCoords: [2]int{x, y},
							Normal:     [2]int{1, 0},
						})
					}
				}
			}
		}
	}

	return moveX
}

// collideY checks for collisions along the Y axis and returns the allowed Y movement
func (c *Collider[T]) collideY(rectX, rectY, rectW, rectH, moveY float64) float64 {

	checkLimit := max(1, int(math.Ceil(math.Abs(moveY)/float64(c.TileSize[1])))+1)

	playerLeft := int(math.Floor(rectX / float64(c.TileSize[0])))
	playerRight := int(math.Ceil((rectX+rectW)/float64(c.TileSize[0]))) - 1

	if moveY > 0 {
		startY := int(math.Floor((rectY + rectH) / float64(c.TileSize[1])))
		endY := startY + checkLimit
		endY = min(endY, len(c.TileMap))

		for x := playerLeft; x <= playerRight; x++ {
			if x < 0 || x >= len(c.TileMap[0]) {
				continue
			}
			for y := startY; y < endY; y++ {
				if y < 0 || y >= len(c.TileMap) {
					continue
				}
				if c.TileMap[y][x] != c.NonSolidTileID {
					tileTop := float64(y * c.TileSize[1])
					collision := tileTop - (rectY + rectH)
					if collision <= moveY {
						moveY = collision
						c.Collisions = append(c.Collisions, CollisionInfo[T]{
							TileID:     c.TileMap[y][x],
							TileCoords: [2]int{x, y},
							Normal:     [2]int{0, -1},
						})
					}
				}
			}
		}
	}

	if moveY < 0 {
		endY := int(math.Floor(rectY / float64(c.TileSize[1])))
		startY := endY - checkLimit
		startY = max(startY, 0)

		for x := playerLeft; x <= playerRight; x++ {
			if x < 0 || x >= len(c.TileMap[0]) {
				continue
			}
			for y := startY; y <= endY; y++ {
				if y < 0 || y >= len(c.TileMap) {
					continue
				}
				if c.TileMap[y][x] != c.NonSolidTileID {
					tileBottom := float64((y + 1) * c.TileSize[1])
					collision := tileBottom - rectY
					if collision >= moveY {
						moveY = collision
						c.Collisions = append(c.Collisions, CollisionInfo[T]{
							TileID:     c.TileMap[y][x],
							TileCoords: [2]int{x, y},
							Normal:     [2]int{0, 1},
						})
					}
				}
			}
		}
	}

	return moveY
}
