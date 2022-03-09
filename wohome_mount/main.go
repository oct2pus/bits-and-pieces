package main

import (
	"log"

	"github.com/deadsy/sdfx/render"
	"github.com/deadsy/sdfx/sdf"
)

type Brace struct {
	faces   []sdf.SDF2
	heights []float64
}

func newBrace(frontDimensions, backDimensions, topDimensions []sdf.V2, frontHeight, backHeight, topHeight float64) (Brace, error) {
	var b Brace

	faces := make([]sdf.Polygon, 0)
	faces = append(faces, *sdf.NewPolygon())
	faces = append(faces, *sdf.NewPolygon())
	faces = append(faces, *sdf.NewPolygon())
	faces[0].AddV2Set(frontDimensions)
	faces[1].AddV2Set(backDimensions)
	faces[2].AddV2Set(topDimensions)

	b.faces = make([]sdf.SDF2, 0)
	frontFace, err := sdf.Polygon2D(faces[0].Vertices())
	if err != nil {
		return b, err
	}
	backFace, err := sdf.Polygon2D(faces[1].Vertices())
	if err != nil {
		return b, err
	}
	topFace, err := sdf.Polygon2D(faces[2].Vertices())
	if err != nil {
		return b, err
	}

	b.faces = append(b.faces, frontFace)
	b.faces = append(b.faces, backFace)
	b.faces = append(b.faces, topFace)

	b.heights = make([]float64, 0)
	b.heights = append(b.heights, frontHeight)
	b.heights = append(b.heights, backHeight)
	b.heights = append(b.heights, topHeight)

	return b, nil
}

func (b Brace) Extrude() sdf.SDF3 {
	parts := make([]sdf.SDF3, 0)
	for x := 0; x < len(b.faces); x++ {
		parts = append(parts, sdf.Extrude3D(b.faces[x], b.heights[x]))
	}
	// handle front
	parts[0] = sdf.Transform3D(parts[0], sdf.Translate3d(sdf.V3{X: 0, Y: -b.faces[0].BoundingBox().Max.Y, Z: b.faces[0].BoundingBox().Max.X - 0.2}))
	// handle back
	parts[1] = sdf.Transform3D(parts[1], sdf.Translate3d(sdf.V3{X: 0, Y: b.faces[2].BoundingBox().Max.Y, Z: b.faces[1].BoundingBox().Max.Y}))

	return sdf.Union3D(parts[0], parts[1], parts[2])
}

type Facet struct {
	peg    sdf.SDF2
	line   sdf.SDF2
	offset sdf.SDF2
	height float64
}

func newFacet(radius float64, lineDimensions []sdf.V2, height float64) (Facet, error) {
	var f Facet
	var err error

	f.peg, err = sdf.Circle2D(radius)
	if err != nil {
		return f, err
	}

	f.line, err = sdf.Polygon2D(lineDimensions)
	if err != nil {
		return f, err
	}

	offsets := make([]sdf.SDF2, 0)
	offsetDimensions := []sdf.V2{
		{X: f.line.BoundingBox().Min.X, Y: f.line.BoundingBox().Min.Y},
		{X: f.line.BoundingBox().Max.X / 4, Y: f.line.BoundingBox().Min.Y},
		{X: f.line.BoundingBox().Max.X / 4, Y: f.line.BoundingBox().Max.Y},
		{X: f.line.BoundingBox().Min.X, Y: f.line.BoundingBox().Max.Y},
	}
	offset, err := sdf.Polygon2D(offsetDimensions)
	if err != nil {
		return f, err
	}
	offsets = append(offsets, offset)
	offsets = append(offsets, offset)

	// move one offset to the other side
	offsets[1] = sdf.Transform2D(offsets[1], sdf.Translate2d(sdf.V2{X: (f.line.BoundingBox().Max.X * 0.75), Y: 0}))

	f.offset = sdf.Union2D(offsets[0], offsets[1])

	f.height = height

	return f, nil
}

func (f Facet) Extrude() sdf.SDF3 {
	peg := sdf.Extrude3D(f.peg, f.height)

	line := sdf.Extrude3D(f.line, f.height)

	offset := sdf.Extrude3D(f.offset, f.height-(f.height/4))

	tab := sdf.Difference3D(line, offset)
	tab = sdf.Transform3D(tab, sdf.Translate3d(sdf.V3{X: -(peg.BoundingBox().Max.X / 1.75), Y: 0, Z: 0}))

	return sdf.Union3D(peg, tab)
}

func main() {
	brace, err := newBrace(
		[]sdf.V2{{X: 0, Y: 0}, {X: 15, Y: 0}, {X: 15, Y: 4}, {X: 0, Y: 4}},
		[]sdf.V2{{X: 0, Y: 0}, {X: 15, Y: 0}, {X: 15, Y: 4}, {X: 0, Y: 4}},
		[]sdf.V2{{X: 0, Y: 0}, {X: 15, Y: 0}, {X: 15, Y: 23}, {X: 0, Y: 23}},
		33.5,
		12,
		4)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	facet, err := newFacet(
		3.8,
		[]sdf.V2{{X: 0, Y: 0}, {X: 4.3, Y: 0}, {X: 4.3, Y: 13.5}, {X: 0, Y: 13.5}},
		5)
	if err != nil {
		log.Fatalf("error: %v\n", err)
	}
	model := combine(brace.Extrude(), facet.Extrude())
	render.ToSTL(model, 300, "brace.stl", &render.MarchingCubesOctree{})
}

func combine(b, f sdf.SDF3) sdf.SDF3 {
	f = sdf.Transform3D(f, sdf.Rotate3d(sdf.V3{X: 1, Y: 0, Z: 0}, sdf.DtoR(270)))
	f = sdf.Transform3D(f, sdf.Translate3d(sdf.V3{X: b.BoundingBox().Max.X / 2, Y: -f.BoundingBox().Max.Y * 2.625, Z: b.BoundingBox().Max.Z * 0.75}))
	return sdf.Union3D(b, f)
}

/*func brace() {
	// X: 15, Y: 33.5, Z: 4
	frontDimensions := []sdf.V2{{X: 0, Y: 0}, {X: 15, Y: 0}, {X: 15, Y: 4}, {X: 0, Y: 4}}
	// X: 15, Y: 12, Z: 4
	backDimensions := []sdf.V2{{X: 0, Y: 0}, {X: 15, Y: 0}, {X: 15, Y: 4}, {X: 0, Y: 4}}
	// X: 15, Y: 4, Z: 23
	topDimensions := []sdf.V2{{X: 0, Y: 0}, {X: 15, Y: 0}, {X: 15, Y: 23}, {X: 0, Y: 23}}
	faces := make([]sdf.Polygon, 3)
	faces = append(faces, *sdf.NewPolygon())
	faces = append(faces, *sdf.NewPolygon())
	faces = append(faces, *sdf.NewPolygon())
	faces[0].AddV2Set(frontDimensions)
	faces[1].AddV2Set(backDimensions)
	faces[2].AddV2Set(topDimensions)

	front_ex := sdf.Extrude3D(, 33.5)

}*/
