package main

import (
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

/*
13 mm - ridge height
10 mm - ridge length
7 mm - remote depth
22.4 mm - remote length
39.5 mm - remote width
*/

func grip(oX, oY, oZ, iX, iY, iZ float64) sdf.SDF3 {
	/*
		gripDimensions := []sdf.V2{{X: 0, Y: 0}, {X: 40.7, Y: 0}, {X: 40.7, Y: 14}, {X: 0, Y: 14}}
		cutoutDimensions := []sdf.V2{{X: 0, Y: 1}, {X: 40.7, Y: 1}, {X: 40.7, Y: 13}, {X: 0, Y: 13}}
		grip, err := sdf.Polygon2D(gripDimensions)
		if err != nil {
			return nil, err
		}
		cutout, err := sdf.Polygon2D(cutoutDimensions)
		if err != nil {
			return nil, err
		}
		gripHands := sdf.Difference2D(grip, cutout)
		grip3D := sdf.Extrude3D(grip, 1)
		gripHands3D := sdf.Extrude3D(gripHands, 12)
		grip3D = sdf.Transform3D(grip3D, sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: -6}))
		grip3D = sdf.Union3D(grip3D, gripHands3D)
		return grip3D, err*/
	outter := sdf.NewBox2(sdf.V2{X: 0, Y: 0}, sdf.V2{X: oX, Y: oY})
	inner := sdf.NewBox2(sdf.V2{X: 0, Y: 0}, sdf.V2{X: iX, Y: iY})
	bot := sdf.NewBox2(sdf.V2{X: 0, Y: 0}, sdf.V2{X: oX, Y: oY})
	o, _ := rect(outter)
	i, _ := rect(inner)
	b, _ := rect(bot)
	o = sdf.Difference2D(o, i)
	grip := sdf.Extrude3D(o, oZ)
	bottom := sdf.Extrude3D(b, iZ)
	bottom = sdf.Transform3D(bottom, sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: -(oZ / 2)}))
	return sdf.Union3D(grip, bottom)
}

func main() {
	body := grip(18.7, 42, 11, 12.7, 42, 3) //, err := grip()
	hole := grip(18.7, 42, 25, 12.7, 40, 3)
	body = sdf.Transform3D(body, sdf.Rotate3d(sdf.V3{X: 0, Y: 1, Z: 0}, sdf.DtoR(270)))
	body = sdf.Transform3D(body, sdf.Translate3d(sdf.V3{X: -(hole.BoundingBox().Max.X * 1.5), Y: 0, Z: (body.BoundingBox().Max.Z / 3) + 0.05}))
	holder := sdf.Union3D(body, hole)
	/*if err != nil {
		log.Fatalf("error: %v\n", err)
	}*/
	render.ToSTL(holder, 300, "holder.stl", &render.MarchingCubesOctree{})
}

func rect(box sdf.Box2) (sdf.SDF2, error) {
	dimensions := []sdf.V2{
		{box.Min.X, box.Min.Y},
		{box.Max.X, box.Min.Y},
		{box.Max.X, box.Max.Y},
		{box.Min.X, box.Max.Y},
	}
	p := sdf.NewPolygon()
	p.AddV2Set(dimensions)
	return sdf.Polygon2D(p.Vertices())
}
