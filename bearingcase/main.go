package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

const LAYERHEIGHT = 0.16
const CIRCLERADIUS = 11.5 // radius of a single bearing
const CIRCLEHEIGHT = 7.7  // height of a single bearing

func main() {
	tube := createHollowCylinder(CIRCLERADIUS, CIRCLEHEIGHT*10)
	cap := createHollowCylinder(CIRCLERADIUS+0.5, CIRCLEHEIGHT+(CIRCLEHEIGHT*0.5))
	render.ToSTL(tube, 500, "tube.stl", &render.MarchingCubesOctree{})
	render.ToSTL(cap, 500, "cap.stl", &render.MarchingCubesOctree{})
}

func createHollowCylinder(radius, height float64) sdf.SDF3 {
	circle, err := sdf.Circle2D(radius + (LAYERHEIGHT * 3))
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	circle2, err := sdf.Circle2D(radius)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	rim := sdf.Extrude3D(circle, height)
	inside := sdf.Extrude3D(circle2, height-(LAYERHEIGHT*3))
	inside = sdf.Transform3D(inside, sdf.Translate3d(sdf.V3{0, 0, (LAYERHEIGHT * 3)}))
	return sdf.Difference3D(rim, inside)
}
