package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

/*
outer diameter:
76.0mm
inner diameter:
70.5 mm
expected length:
30.0 mm
thickness of lamp:
2.6~2.9mm, use 3.0 mm for tolerance

this model is a piece of shit
i hate my light fixtures
*/
func main() {
	render.ToSTL(modelThinner(0.2), "model02.stl", render.NewMarchingCubesUniform(800))
	render.ToSTL(model(0.5), "model05.stl", render.NewMarchingCubesUniform(800))
	render.ToSTL(model(0.9), "model09.stl", render.NewMarchingCubesUniform(800))
}

func model(tol float64) sdf.SDF3 {
	wDia, pDia, oDia, iDia := 77.0, 69.5, 76.0, 70.5
	//tol := 0.9 // far right
	//tol := 0.5 // far left
	//tol := 0.2 //center two
	//	h := 30.0
	thick := 3.0
	stopperLength := thick / 3
	length := 30.0

	//grip
	outer, err := sdf.Circle2D((oDia + tol) / 2)
	if err != nil {
		log.Fatalln(err)
	}
	inner, err := sdf.Circle2D((iDia - tol) / 2)
	if err != nil {
		log.Fatalln(err)
	}
	lips2D, err := sdf.Circle2D((wDia + tol) / 2)
	if err != nil {
		log.Fatalln(err)
	}
	cutout, err := sdf.Circle2D((pDia - tol) / 2)
	if err != nil {
		log.Fatalln(err)
	}

	walls := sdf.Difference2D(lips2D, inner)
	grip := sdf.Difference2D(lips2D, outer)
	inner = sdf.Difference2D(inner, cutout)
	grip = sdf.Union2D(grip, inner)
	walls = sdf.Union2D(walls, inner)

	lips := sdf.Extrude3D(grip, thick)
	stopper := sdf.Extrude3D(walls, stopperLength)

	supportCone, _ := sdf.Cone3D(thick, (wDia)/2, (wDia-tol)/2, 0)
	supportConeNegative, _ := sdf.Cone3D(thick, (pDia-tol*2)/2, (wDia-tol)/2, 0)
	supportCone = sdf.Difference3D(supportCone, supportConeNegative)
	//	top := sdf.Extrude3D(lips2D, thick/6)
	cone, _ := sdf.Cone3D(length, (wDia+tol)/2, (iDia-tol)/2, 0)
	cone2, _ := sdf.Cone3D(length-stopperLength, (wDia+tol-2)/2, (iDia-tol-2)/2, 0)
	cone = sdf.Difference3D(
		cone,
		sdf.Transform3D(cone2, sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: -thick / 6})),
	)
	model := sdf.Union3D(
		lips,
		sdf.Transform3D(stopper, sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: thick/2 + stopperLength/2})),
		sdf.Transform3D(supportCone, sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: thick + stopperLength})),
		sdf.Transform3D(cone, sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: thick/2 + stopperLength + length/2})),
	)

	return model
}

// modelThinner is for thinner tolerance versions of the model
func modelThinner(tol float64) sdf.SDF3 {
	wDia, pDia, oDia, iDia := 77.0, 69.5, 76.0, 70.5
	//tol := 0.9 // far right
	//tol := 0.5 // far left
	//tol := 0.2 //center two
	//	h := 30.0
	thick := 3.0
	stopperLength := thick / 3
	length := 30.0

	//grip
	outer, err := sdf.Circle2D((oDia + tol) / 2)
	if err != nil {
		log.Fatalln(err)
	}
	inner, err := sdf.Circle2D((iDia - tol) / 2)
	if err != nil {
		log.Fatalln(err)
	}
	lips2D, err := sdf.Circle2D((wDia + tol) / 2)
	if err != nil {
		log.Fatalln(err)
	}
	cutout, err := sdf.Circle2D((pDia - tol) / 2)
	if err != nil {
		log.Fatalln(err)
	}

	walls := sdf.Difference2D(lips2D, inner)
	grip := sdf.Difference2D(lips2D, outer)
	inner = sdf.Difference2D(inner, cutout)
	grip = sdf.Union2D(grip, inner)
	walls = sdf.Union2D(walls, inner)

	lips := sdf.Extrude3D(grip, thick)
	stopper := sdf.Extrude3D(walls, stopperLength)

	supportCone, _ := sdf.Cone3D(thick, (wDia-tol*2)/2, (wDia-tol*3)/2, 0)
	supportConeNegative, _ := sdf.Cone3D(thick, (pDia-tol*2)/2, (wDia-tol)/2, 0)
	supportCone = sdf.Difference3D(supportCone, supportConeNegative)
	//	top := sdf.Extrude3D(lips2D, thick/6)
	cone, _ := sdf.Cone3D(length, (wDia+tol)/2, (iDia-tol)/2, 0)
	cone2, _ := sdf.Cone3D(length-stopperLength, (wDia+tol-2)/2, (iDia-tol-2)/2, 0)
	cone = sdf.Difference3D(
		cone,
		sdf.Transform3D(cone2, sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: -thick / 6})),
	)
	model := sdf.Union3D(
		lips,
		sdf.Transform3D(stopper, sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: thick/2 + stopperLength/2})),
		sdf.Transform3D(supportCone, sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: thick + stopperLength})),
		sdf.Transform3D(cone, sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: thick/2 + stopperLength + length/2})),
	)

	return model
}
