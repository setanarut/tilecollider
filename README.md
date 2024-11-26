# tilecollider
Package tilecollider provides collision detection for tile-based games in Go language.

The method used is borrowed from [higherorderfun.com](http://higherorderfun.com/blog/2012/05/20/the-guide-to-implementing-2d-platformers/)

## Algorithm

For movement in each axis (x, y):
 1. We go from the tile coordinate of the leading edge of the movement to the tile coordinate of the destination of that edge.
 2. We go from the base tile coordinate of the opposite axis to the max tile coordinate of the opposite axis.
 3. For each of those tile coordinates `x` and `y`, we call the `OnCollision` callback. if it returns `true`, we've hit something and we stop checking that axis for collisions.
 4. Returns movement delta after collision (actual position change)

## Example

```Go
move := tilecollider.Collide(
	&Offset, // collision detection grid offset
	&GridSize, // collision detection grid size
	&TileSize,
	&Box, // player hit box
	&Vel, // velocity
	TileMap,
	func(axis, direction int, tileID uint8, tileCoords [2]int) bool {
		collidedTiles = append(collidedTiles, tileCoords)
		if axis == 1 && direction == 1 {
			IsOnFloor = true
		}
		if axis == 0 {
			Vel[0] = 0
		}
		if axis == 1 && direction == -1 {
			IsJumping = false
			Vel[1] = 0
		}
		return true
	})
Translate(&Box, move)
```