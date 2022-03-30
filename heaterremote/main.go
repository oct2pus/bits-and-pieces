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
lets use 1.6 mm for size of outline
main body needs to be a little bigger
*/

func grip() sdf.SDF3 { //, error) {
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
	outter := sdf.NewBox2(sdf.V2{X: 0, Y: 0}, sdf.V2{X: 18, Y: 40})
	inner := sdf.NewBox2(sdf.V2{X: 0, Y: 0}, sdf.V2{X: 12.7, Y: 40})
	bot := sdf.NewBox2(sdf.V2{X: 0, Y: 0}, sdf.V2{X: 18, Y: 40})
	o, _ := rect(outter)
	i, _ := rect(inner)
	b, _ := rect(bot)
	o = sdf.Difference2D(o, i)
	grip := sdf.Extrude3D(o, 10)
	bottom := sdf.Extrude3D(b, 3)
	bottom = sdf.Transform3D(bottom, sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: -4.25}))
	return sdf.Union3D(grip, bottom)
}

func main() {
	grip := grip() //, err := grip()
	/*if err != nil {
		log.Fatalf("error: %v\n", err)
	}*/
	render.ToSTL(grip, 300, "grip.stl", &render.MarchingCubesOctree{})
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
