// Package tilecollider provides collision detection for tile-based games.
//
// -----------------------------------------------------------------------------
//
// http://higherorderfun.com/blog/2012/05/20/the-guide-to-implementing-2d-platformers/
//
// https://github.com/hackergrrl/collide-2d-aabb-tilemap.
package tilecollider

import (
	"math"

	"golang.org/x/exp/constraints"
)

// Integer is an alias for constraints.Integer
type Integer = constraints.Integer

// OnCollide is a callback function that handles collisions.
//
// If your callback returns true, it is assumed that you're done checking
// tiles along this access and it will move to the next axis if any.
//
// Parameters:
//   - axis: collision axis (0 for X, 1 for Y)
//   - dir: collision direction (-1 for negative, 1 for positive)
//   - tile: tile ID from the tilemap at collision point
//   - coords: [x,y] grid coordinates of the collided tile
type OnCollide[T Integer] func(axis, dir int, tile T, coords [2]int) bool

// Collide handles collision detection and response.
//
// # Returns:
//
//   - Movement delta after collision (actual position change)
//
// # Parameters:
//
//   - o: Collision detection grid offset
//   - g: Collision detection grid size
//   - t: Tilemap tile size [w and h]
//   - b: Player box [x, y, w, h] (xy is upper-left)
//   - m: Move delta vector
//   - tm: 2D slice representing the tile grid, where each cell contains a tile ID of type T. 0 is dedicated to non-solid tile
//   - on: Collision callback function
//
// # Example:
//
//		 move := tilecollider.Collide(
//			&Offset,
//			&GridSize,
//			&TileSize,
//			&Box,
//			&Vel,
//			TileMap,
//			func(axis, direction int, tileID uint8, tileCoords [2]int) bool {
//				collidedTiles = append(collidedTiles, tileCoords)
//				if axis == 1 && direction == 1 {
//					IsOnFloor = true
//				}
//				if axis == 0 {
//					Vel[0] = 0
//				}
//				if axis == 1 && direction == -1 {
//					IsJumping = false
//					Vel[1] = 0
//				}
//				return true
//			})
//
//	 Translate(&Box, move)
func Collide[T Integer](o, g, t *[2]int, b *[4]float64, m *[2]float64, tm [][]T, on OnCollide[T]) [2]float64 {
	move := [2]float64{m[0], m[1]}
	if m[0] == 0 && m[1] == 0 {
		return move
	}
	move[0] = CollideAxis[T](0, o, g, t, b, m, tm, on)
	move[1] = CollideAxis[T](1, o, g, t, b, m, tm, on)
	return move
}

// CollideAxis handles collision detection along a single axis.
//
// Returns:
//
//   - movement delta after collision (actual position change)
//
// Parameters:
//
//   - a: Collision detection axis (0 for X, 1 for Y)
//   - o: Collision detection grid offset
//   - g: Collision detection grid size
//   - t: Tilemap tile size [w and h]
//   - b: Player box [x, y, w, h] (xy is upper-left)
//   - m: move delta vector
//   - tm: 2D slice representing the tile grid, where each cell contains a tile ID of type T
//   - on: Collision callback function
func CollideAxis[T Integer](a int, o, g, t *[2]int, b *[4]float64, m *[2]float64, tm [][]T, on OnCollide[T]) float64 {
	// l: location - coordinates of the tile being checked for collision
	l := [2]int{0, 0}

	// p: positive - determines direction of movement
	p := m[a] > 0

	// e: edge - the edge of the box to check based on movement direction
	var e float64
	if p {
		e = b[a] + b[a+2] // positive movement: right/bottom edge
	} else {
		e = b[a] // negative movement: left/top edge
	}

	// d: direction - sign of movement (+1 or -1)
	d := 1
	if !p {
		d = -1
	}

	// x: cross axis - opposite of the checked axis (if 0 then 1, if 1 then 0)
	x := 1 - a

	// s: start - first tile to check on cross axis
	s := int(math.Floor(b[x] / float64(t[x])))

	// n: end - last tile to check on cross axis
	n := int(math.Ceil((b[x] + b[x+2]) / float64(t[x])))

	// c: current - current tile position on movement axis
	c := int(math.Floor(e / float64(t[a])))

	// r: reach - target tile position on movement axis
	r := int(math.Floor((e+m[a])/float64(t[a]))) + d

	// Check collision for each tile
	for i := c; i != r; i += d {
		// Check map boundaries
		if i < o[a] || i >= g[a] {
			continue
		}
		// Check tiles on cross axis
		for j := s; j < n; j++ {
			if j < o[x] || j >= g[x] {
				continue
			}
			l[a] = i
			l[x] = j
			tileID := tm[l[1]][l[0]]
			if tileID == 0 {
				continue
			}

			// Calculate collision point
			var collisionEdge float64
			if d > 0 {
				collisionEdge = float64(i * t[a])
			} else {
				collisionEdge = float64((i + 1) * t[a])
			}

			// Distance to collision
			delta := collisionEdge - e

			// Collision occurred, return new position
			if on(a, d, tileID, l) {
				return delta
			}
		}
	}

	// No collision, return full movement
	return m[a]
}
