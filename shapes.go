package svg

import (
	"encoding/xml"
	"strconv"
)

// ShapeObject embeds Object and provides a PathLength attribute
// field that is common to all basic shapes
type ShapeObject struct {
	Object
	PathLength float64 `xml:"pathLength,attr,omitempty"`
}

// LineInt draws a line specified by integer coordinates.
func (el *ElemList) LineInt(x1, y1, x2, y2 int) *ShapeObject {
	l := &line{X1: float64(x1), Y1: float64(y1), X2: float64(x2), Y2: float64(y2)}
	el.append(l)
	return &l.ShapeObject
}

type line struct {
	XMLName xml.Name `xml:"line"`
	X1      float64  `xml:"x1,attr"`
	Y1      float64  `xml:"y1,attr"`
	X2      float64  `xml:"x2,attr"`
	Y2      float64  `xml:"y2,attr"`
	ShapeObject
}

// RectInt draws a rectangle based on integer coordinates.
func (el *ElemList) RectInt(x, y, w, h int) *Rect {
	r := &Rect{X: float64(x), Y: float64(y), Width: float64(w), Height: float64(h)}
	el.append(r)
	return r
}

type Rect struct {
	XMLName     xml.Name `xml:"rect"`
	X           float64  `xml:"x,attr,omitempty"`
	Y           float64  `xml:"y,attr,omitempty"`
	Width       float64  `xml:"width,attr"`
	Height      float64  `xml:"height,attr"`
	Rx          float64  `xml:"rx,attr,omitempty"`
	Ry          float64  `xml:"ry,attr,omitempty"`
	ShapeObject `xml:"x,attr,omitempty"`
}

// CircleInt draws a circle based on integer coordinates.
func (el *ElemList) CircleInt(cx, cy, r int) *ShapeObject {
	c := &circle{X: float64(cx), Y: float64(cy), R: float64(r)}
	el.append(c)
	return &c.ShapeObject
}

type circle struct {
	XMLName xml.Name `xml:"circle"`
	X       float64  `xml:"cx,attr"`
	Y       float64  `xml:"cy,attr"`
	R       float64  `xml:"r,attr"`
	ShapeObject
}

// EllipseInt draws an ellipse based on integer coordinates.
func (el *ElemList) EllipseInt(cx, cy, rx, ry int) *ShapeObject {
	e := &ellipse{X: float64(cx), Y: float64(cy), Rx: float64(rx), Ry: float64(ry)}
	el.append(e)
	return &e.ShapeObject
}

type ellipse struct {
	XMLName xml.Name `xml:"circle"`
	X       float64  `xml:"cx,attr"`
	Y       float64  `xml:"cy,attr"`
	Rx      float64  `xml:"rx,attr"`
	Ry      float64  `xml:"ry,attr"`
	ShapeObject
}

// Polyline adds an empty polyline element to the ElemList.
// Points may be added using the AddInt method of the returned
// object.
func (el *ElemList) PolyLine() *PolyLine {
	line := &PolyLine{}
	el.append(line)
	return line
}

type PolyLine struct {
	XMLName xml.Name `xml:"polyline"`
	Points  `xml:"points,attr"`
	ShapeObject
}

func (line *PolyLine) PreAlloc(n int) *PolyLine {
	if line.Points == nil {
		line.Points = make(Points, 0, n)
	}
	return line
}

// Polygon adds an empty polygon element to the ElemList.
// Points may be added using the AddInt method of the returned
// object.
func (el *ElemList) Polygon() *PolyLine {
	p := &polygon{}
	el.append(p)
	return &p.PolyLine
}

type polygon struct {
	XMLName xml.Name `xml:"polygon"`
	PolyLine
}

// Points is a slice of 2D coordinates that marshals,
// if used in an XML attribute, into a  list of space separated pairs
// of comma separated numbers.
type Points [][2]float64

func (pts Points) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	s := make([]string, len(pts))
	ff := func(f float64) string {
		return strconv.FormatFloat(f, 'g', -1, 64)
	}
	for i, pt := range pts {
		s[i] = ff(pt[0]) + "," + ff(pt[1])
	}
	return makeListAttr(name, s)
}

// AddInt adds a point specified by integer coordinates.
func (pts *Points) AddInt(x, y int) {
	*pts = append(*pts, [2]float64{float64(x), float64(y)})
}

// Path adds a <path> element.
func (el *ElemList) Path(d string) *ShapeObject {
	p := &path{D: d}
	el.append(p)
	return &p.ShapeObject
}

type path struct {
	XMLName xml.Name `xml:"path"`
	D       string   `xml:"d,attr,omitempty"`
	ShapeObject
}
