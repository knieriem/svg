package svg

import (
	"encoding/xml"
	"strconv"
)

// LineInt draws a line specified by integer coordinates.
func (el *ElemList) LineInt(x1, y1, x2, y2 int) *Object {
	l := &line{X1: float64(x1), Y1: float64(y1), X2: float64(x2), Y2: float64(y2)}
	el.append(l)
	return &l.Object
}

type line struct {
	XMLName xml.Name `xml:"line"`
	X1      float64  `xml:"x1,attr"`
	Y1      float64  `xml:"y1,attr"`
	X2      float64  `xml:"x2,attr"`
	Y2      float64  `xml:"y2,attr"`
	Object
}

// RectInt draws a rectangle based on integer coordinates.
func (el *ElemList) RectInt(x, y, w, h int) *Object {
	r := &rect{X: float64(x), Y: float64(y), Width: float64(w), Height: float64(h)}
	el.append(r)
	return &r.Object
}

type rect struct {
	XMLName xml.Name `xml:"rect"`
	X       float64  `xml:"x,attr"`
	Y       float64  `xml:"y,attr"`
	Width   float64  `xml:"width,attr"`
	Height  float64  `xml:"height,attr"`
	Object
}

// CircleInt draws a circle based on integer coordinates.
func (el *ElemList) CircleInt(cx, cy, r int) *Object {
	c := &circle{X: float64(cx), Y: float64(cy), R: float64(r)}
	el.append(c)
	return &c.Object
}

type circle struct {
	XMLName xml.Name `xml:"circle"`
	X       float64  `xml:"cx,attr"`
	Y       float64  `xml:"cy,attr"`
	R       float64  `xml:"r,attr"`
	Object
}

// EllipseInt draws an ellipse based on integer coordinates.
func (el *ElemList) EllipseInt(cx, cy, rx, ry int) *Object {
	e := &ellipse{X: float64(cx), Y: float64(cy), Rx: float64(rx), Ry: float64(ry)}
	el.append(e)
	return &e.Object
}

type ellipse struct {
	XMLName xml.Name `xml:"circle"`
	X       float64  `xml:"cx,attr"`
	Y       float64  `xml:"cy,attr"`
	Rx      float64  `xml:"rx,attr"`
	Ry      float64  `xml:"ry,attr"`
	Object
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
	Object
}

func (line *PolyLine) PreAlloc(n int) *PolyLine {
	if line.Points == nil {
		line.Points = make(Points, 0, n)
	}
	return line
}

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
