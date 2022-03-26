package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

/*
13 mm - ridge height
10 mm - ridge length
7 mm - remote depth
22.4 mm - remote length
39.5 mm - remote width
lets use 1.2 mm for size of outline
main body needs to be a little bigger
*/

func grip() (sdf.SDF3, error) {
	gripDimensions := []sdf.V2{{X: 0, Y: 0}, {X: 40.7, Y: 0}, {X: 40.7, Y: 13.6}, {X: 0, Y: 13.6}}
	cutoutDimensions := []sdf.V2{{X: 0, Y: 1.2}, {X: 40.7, Y: 1.2}, {X: 40.7, Y: 12.4}, {X: 0, Y: 12.4}}
	grip, err := sdf.Polygon2D(gripDimensions)
	if err != nil {
		return nil, err
	}
	cutout, err := sdf.Polygon2D(cutoutDimensions)
	if err != nil {
		return nil, err
	}
	gripHands := sdf.Difference2D(grip, cutout)
	grip3D := sdf.Extrude3D(grip, 1.2)
	gripHands3D := sdf.Extrude3D(gripHands, 11.2)
	grip3D = sdf.Transform3D(grip3D, sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: -5}))
	grip3D = sdf.Union3D(grip3D, gripHands3D)
	return grip3D, err
}

func main() {
	grip, err := grip()
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	render.ToSTL(grip, 300, "grip.stl", &render.MarchingCubesOctree{})
}
