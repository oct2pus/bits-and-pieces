package main

import (
	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

const (
	BRACE_LENGTH   = WOOD_THICKNESS + (BRACE_GIRTH * 2) + TOLERANCE
	BRACE_HEIGHT   = 50.0
	BRACE_GIRTH    = 12.5
	HOLE_SPACING   = 46.2
	M4_DIAMETER    = 4.0
	MOUNT_HEIGHT   = 65.0
	TOLERANCE      = 0.5
	WOOD_THICKNESS = 15.0
)

// render

func main() {
	render.ToDXF(brace2D(), "brace.dxf", render.NewDualContouring2D(400))
	render.ToSTL(brace(), "brace.stl", render.NewMarchingCubesOctree(600))
	// render.ToSTL(screwHoles(), "screwHoles.stl", render.NewMarchingCubesOctree(600))
}

// 3D

func brace() sdf.SDF3 {
	brace := sdf.Extrude3D(brace2D(), HOLE_SPACING+BRACE_GIRTH)
	mount := mount()
	//braceTop = sdf.Transform2D(braceTop, sdf.Translate2d(v2.Vec{X: 0, Y: height/2 + topHeight/2 - BRACE_GIRTH/8}))

	mount = sdf.Transform3D(mount, sdf.Translate3d(v3.Vec{X: 0, Y: BRACE_HEIGHT/2 + MOUNT_HEIGHT/2 - BRACE_GIRTH/8, Z: 0}))

	brace = sdf.Union3D(brace, mount)

	return brace
}

func mount() sdf.SDF3 {
	mount := sdf.Extrude3D(mount2D(), HOLE_SPACING+BRACE_GIRTH)
	screwHoles := sdf.Transform3D(screwHoles(), sdf.RotateY(sdf.DtoR(270)))

	screwHoles = sdf.Transform3D(screwHoles, sdf.Translate3d(v3.Vec{X: -BRACE_HEIGHT/2 - (-BRACE_GIRTH / 1.15), Y: 0, Z: 0}))

	return sdf.Difference3D(mount, screwHoles)
}

func screwHoles() sdf.SDF3 {
	coneOuterDiameter := 7.0
	coneHeight := 2.3 + TOLERANCE

	screwHoles := sdf.Extrude3D(screwHoles2D(), BRACE_GIRTH)
	cone, _ := sdf.Cone3D(coneHeight, M4_DIAMETER/2, coneOuterDiameter/2, 0)
	cone = sdf.Transform3D(cone, sdf.Translate3d(v3.Vec{X: 0, Y: 0, Z: BRACE_GIRTH/2 - coneHeight/2}))

	return sdf.Union3D(
		screwHoles,
		sdf.Transform3D(cone, sdf.Translate3d(v3.Vec{X: HOLE_SPACING / 2, Y: 0, Z: 0})),
		sdf.Transform3D(cone, sdf.Translate3d(v3.Vec{X: -HOLE_SPACING / 2, Y: 0, Z: 0})),
	)
}

// 2D

func brace2D() sdf.SDF2 {

	braceBody := sdf.Box2D(v2.Vec{X: BRACE_LENGTH, Y: BRACE_GIRTH}, 0)
	braceBack := sdf.Box2D(v2.Vec{X: BRACE_GIRTH, Y: BRACE_HEIGHT}, 0.5)
	braceFront := braceBack

	braceBody = sdf.Transform2D(braceBody, sdf.Translate2d(v2.Vec{X: 0, Y: BRACE_HEIGHT/2 - (BRACE_GIRTH / 2)}))
	braceBack = sdf.Transform2D(braceBack, sdf.Translate2d(v2.Vec{X: -BRACE_LENGTH/2 - (-BRACE_GIRTH / 2), Y: 0}))
	braceFront = sdf.Transform2D(braceFront, sdf.Translate2d(v2.Vec{X: BRACE_LENGTH/2 - (BRACE_GIRTH / 2), Y: 0}))

	brace := sdf.Union2D(braceBody, braceBack, braceFront)
	return brace
}

func mount2D() sdf.SDF2 {
	circleDiameter := 62.0

	mount := sdf.Box2D(v2.Vec{X: BRACE_LENGTH, Y: MOUNT_HEIGHT}, 0.5)
	circleCutout, _ := sdf.Circle2D(circleDiameter / 2)

	circleCutout = sdf.Transform2D(circleCutout, sdf.Translate2d(v2.Vec{X: circleDiameter / 3, Y: 0}))

	mount = sdf.Difference2D(mount, circleCutout)
	mount = sdf.Cut2D(mount, v2.Vec{X: 0, Y: BRACE_LENGTH / 3}, v2.Vec{X: 1, Y: 0})
	return mount
}

func screwHoles2D() sdf.SDF2 {
	m4, _ := sdf.Circle2D(M4_DIAMETER / 2)

	return sdf.Union2D(
		sdf.Transform2D(m4, sdf.Translate2d(v2.Vec{X: HOLE_SPACING / 2, Y: 0})),
		sdf.Transform2D(m4, sdf.Translate2d(v2.Vec{X: -HOLE_SPACING / 2, Y: 0})),
	)
}

/*
wood thickness - 15.0mm
wohome height - 71.0mm
screw spacing - 46.2
screw size = m4
*/
