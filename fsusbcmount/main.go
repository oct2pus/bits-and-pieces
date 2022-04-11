package main

import (
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

func port(oX, oY, iX, iY, z float64) sdf.SDF3 {
	box := sdf.Box2D(sdf.V2{X: oX, Y: oY}, 0)
	porthole := sdf.Box2D(sdf.V2{X: iX, Y: iY}, 4)
	port := sdf.Difference2D(box, porthole)
	return sdf.Extrude3D(port, z)
}

func mount(x, y, z, diameter, spacing float64) sdf.SDF3 {
	plate := sdf.Box2D(sdf.V2{X: x, Y: y}, 0)
	hole, _ := sdf.Circle2D(diameter)
	hole2 := hole
	hole = sdf.Transform2D(hole, sdf.Translate2d(sdf.V2{X: 0, Y: -(spacing / 2)}))
	hole2 = sdf.Transform2D(hole2, sdf.Translate2d(sdf.V2{X: 0, Y: (spacing / 2)}))
	plate = sdf.Difference2D(plate, hole)
	plate = sdf.Difference2D(plate, hole2)
	return sdf.Extrude3D(plate, z)
}

func peg(x, y, z, diameter, spacing float64) sdf.SDF3 {
	plane := sdf.Box2D(sdf.V2{X: x, Y: y}, 0)
	hole, _ := sdf.Circle2D(diameter)
	hole2 := hole
	hole = sdf.Transform2D(hole, sdf.Translate2d(sdf.V2{X: 0, Y: -(spacing / 2)}))
	hole2 = sdf.Transform2D(hole2, sdf.Translate2d(sdf.V2{X: 0, Y: (spacing / 2)}))
	peg := sdf.Extrude3D(hole, z*2)
	peg2 := sdf.Extrude3D(hole2, z*2)
	peg = sdf.Transform3D(peg, sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: z}))
	peg2 = sdf.Transform3D(peg2, sdf.Translate3d(sdf.V3{X: 0, Y: 0, Z: z}))
	plate := sdf.Extrude3D(plane, z)
	plate = sdf.Union3D(plate, peg, peg2)
	return plate
}

func connectors(x, y, z float64) sdf.SDF3 {
	line := sdf.Box2D(sdf.V2{X: x, Y: y}, 0)
	line2 := line
	line2 = sdf.Transform2D(line2, sdf.Translate2d(sdf.V2{X: y, Y: 0}))
	line = sdf.Union2D(line, line2)
	return sdf.Extrude3D(line, z)
}

func combine(p, m, g, c sdf.SDF3, width, length, height float64) sdf.SDF3 {
	c = sdf.Transform3D(c, sdf.Rotate3d(sdf.V3{X: 0, Y: 0, Z: 1}, sdf.DtoR(90)))
	c = sdf.Transform3D(c, sdf.MirrorXY())
	g = sdf.Transform3D(g, sdf.Translate3d(sdf.V3{X: width * 1.15, Y: length / 1.9, Z: 0}))
	m = sdf.Transform3D(m, sdf.Translate3d(sdf.V3{X: -width / 2.5, Y: length / 1.9, Z: 0}))
	p = sdf.Transform3D(p, sdf.Rotate3d(sdf.V3{X: 0, Y: 1, Z: 0}, sdf.DtoR(90)))
	p = sdf.Transform3D(p, sdf.Translate3d(sdf.V3{X: -(width * 1.5), Y: length / 1.9, Z: height * 2.5}))
	return sdf.Union3D(p, m, g, c)
}

func main() {
	width, length, height, diameter, spacing := 10.0, 30.0, 2.0, 1.3, 22.0
	p := port(width+2, length, 7.8, 15.6, height)
	m := mount(width, length, height, diameter, spacing)
	g := mount(width/2, length, height, diameter+0.2, spacing)
	c := connectors(width/2, 32, height)
	part := combine(p, m, g, c, width, length, height)
	render.ToSTL(p, 300, "port.stl", &render.MarchingCubesOctree{})
	render.ToSTL(m, 300, "plate.stl", &render.MarchingCubesOctree{})
	render.ToSTL(g, 300, "peg.stl", &render.MarchingCubesOctree{})
	render.ToSTL(c, 300, "connectors.stl", &render.MarchingCubesOctree{})
	render.ToSTL(part, 300, "part.stl", &render.MarchingCubesOctree{})
}
