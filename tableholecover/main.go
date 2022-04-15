package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

/* Dimensions
holder width 78.2 mm
opening width 151.3 mm
opening height 70.1 mm
opening depth 30.6 mm
*/

func top(x, y, z, offset float64) sdf.SDF3 {
	log.Println("Hello world :)")
	outline, err := sdf.Box3D(sdf.V3{X: x + offset, Y: y + offset, Z: z}, 3)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	inline2D := sdf.Box2D(sdf.V2{X: x, Y: y}, 4)
	inline := sdf.Extrude3D(inline2D, z)
	bottom, _ := sdf.Box3D(sdf.V3{X: x + offset, Y: y + offset, Z: z / 2}, 0)
	bottom = sdf.Transform3D(bottom, sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: -z / 3}))
	outline = sdf.Difference3D(outline, inline)
	outline = sdf.Difference3D(outline, bottom)
	return outline
}

func inside(x, y, z, thickness float64) sdf.SDF3 {
	wall := sdf.Box2D(sdf.V2{X: x, Y: y}, 4)
	cavity := sdf.Box2D(sdf.V2{X: x - thickness, Y: y - thickness}, 4)
	wall = sdf.Difference2D(wall, cavity)
	return sdf.Extrude3D(wall, z)
}

func main() {
	x, y, z := 151.3, 70.1, 30.6
	outline := top(x, y, 3, 14)
	inner := inside(x, y, z, 5)
	obstruction := sdf.Box2D(sdf.V2{X: 79.2 + 2, Y: 4 + 14}, 0)
	obstruction = sdf.Transform2D(obstruction, sdf.Translate2d(sdf.V2{X: 0, Y: y/2 - 4/2}))
	outline = sdf.Transform3D(outline, sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: inner.BoundingBox().Max.Z - 1.5}))
	clamp := sdf.Extrude3D(obstruction, z)
	model := sdf.Union3D(outline, inner)
	model = sdf.Difference3D(model, clamp)
	model = sdf.Transform3D(model, sdf.Rotate3d(sdf.V3{X: 0, Y: 1, Z: 0}, sdf.DtoR(180)))
	render.ToSTL(model, 300, "liner.stl", &render.MarchingCubesUniform{})
	//	render.ToSTL(inline, 300, "inline.stl", &render.MarchingCubesUniform{})
}
