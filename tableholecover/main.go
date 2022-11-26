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
	outline, err := sdf.Box3D(sdf.V3{X: x + offset, Y: y + offset, Z: z}, 3)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	inline2D := sdf.Box2D(sdf.V2{X: x, Y: y}, 4)
	inline := sdf.Extrude3D(inline2D, z)
	bottom, err := sdf.Box3D(sdf.V3{X: x + offset, Y: y + offset, Z: z / 2}, 0)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
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

// bundle of wires is 2.0 mm approximately, do 2 as radius for 4mm length circle?
func cap(x, y, z, offset, radius float64) sdf.SDF3 {
	expansion := 2.0
	underside, err := sdf.Box3D(sdf.V3{X: x + offset, Y: y + offset, Z: z + 2}, 3)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	underside = sdf.Transform3D(underside, sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: -1}))
	topside, err := sdf.Box3D(sdf.V3{X: x + offset + expansion, Y: y + offset + expansion, Z: z + 3}, 3)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}

	bottom, err := sdf.Box3D(sdf.V3{X: x + offset + expansion + 0.4, Y: y + offset + expansion + 0.4, Z: z}, 0)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	bottom = sdf.Transform3D(bottom, sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: -z / 2}))

	cableHole2D, err := sdf.Circle2D(radius)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	cableHole2D = sdf.Transform2D(cableHole2D, sdf.Translate2d(sdf.V2{X: 0, Y: y / 3}))
	cableHole := sdf.Extrude3D(cableHole2D, z+4)

	fingerHole2D, err := sdf.Circle2D(radius / 2)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	fingerHole2D = sdf.Transform2D(fingerHole2D, sdf.Translate2d(sdf.V2{X: 0, Y: -y / 3}))
	fingerHole := sdf.Extrude3D(fingerHole2D, z+4)

	cover := sdf.Difference3D(topside, underside)
	cover = sdf.Difference3D(cover, bottom)
	cover = sdf.Difference3D(cover, cableHole)
	cover = sdf.Difference3D(cover, fingerHole)

	return cover
}

func main() {
	x, y, z := 151.3, 70.1, 30.6

	// liner
	outline := top(x, y, 3, 14)
	inner := inside(x, y, z, 5)
	outline = sdf.Transform3D(outline, sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: inner.BoundingBox().Max.Z - 1.5}))
	liner := sdf.Union3D(outline, inner)

	//cover
	cover := cap(x, y, 3, 14, 20)

	// account for monitor arm clamp
	obstruction := sdf.Box2D(sdf.V2{X: 79.2 + 2, Y: 4 + 16}, 0)
	obstruction = sdf.Transform2D(obstruction, sdf.Translate2d(sdf.V2{X: 0, Y: y/2 - 4/2}))
	clamp := sdf.Extrude3D(obstruction, 100)
	liner = sdf.Difference3D(liner, clamp)

	obstruction = sdf.Box2D(sdf.V2{X: 79.2 + 1, Y: 4 + 16}, 0)
	obstruction = sdf.Transform2D(obstruction, sdf.Translate2d(sdf.V2{X: 0, Y: y/2 - 4/2}))
	clamp = sdf.Extrude3D(obstruction, 100)
	cover = sdf.Difference3D(cover, clamp)

	// modify print orientation
	liner = sdf.Transform3D(liner, sdf.Rotate3d(sdf.V3{X: 0, Y: 1, Z: 0}, sdf.DtoR(180)))
	cover = sdf.Transform3D(cover, sdf.Rotate3d(sdf.V3{X: 0, Y: 1, Z: 0}, sdf.DtoR(180)))

	render.ToSTL(liner, 1200, "liner.stl", &render.MarchingCubesUniform{})
	render.ToSTL(cover, 1200, "cover.stl", &render.MarchingCubesUniform{})
	//	render.ToSTL(inline, 300, "inline.stl", &render.MarchingCubesUniform{})
}
