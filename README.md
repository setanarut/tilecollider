[![GoDoc](https://godoc.org/github.com/setanarut/tilecollider?status.svg)](https://pkg.go.dev/github.com/setanarut/tilecollider)

# tilecollider

A simple 2D tile-based collision detection package for Go.

![screen](https://github.com/user-attachments/assets/c2c3c53c-a4f5-4afb-b3a7-4812e0aeb9ab)

## Features

- Fast tile-based collision detection
- Easy integration with game engines like Ebitengine
- Generic tile map support with any Integer type `[y][x]T`
- Support for non-square tiles (different width and height values)
- Adaptive iteration count based on movement speed (anti-tunneling)

## Installation

```sh
go get github.com/setanarut/tilecollider
```

## Usage

See the [examples](./examples) directory for usage example.

## Run Ebitengine example on your machine

```sh
go run github.com/setanarut/tilecollider/examples/demo@latest
```
