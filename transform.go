package svg

import (
	"encoding/xml"
	"strconv"
	"strings"
)

// TransformList is a slice of SVG transformations,
// that marshals into a list of transformation specifications
// to be used in the transform attribute of group containers
type TransformList []Transform

func (tl *TransformList) append(t Transform) *TransformList {
	*tl = append(*tl, t)
	return tl
}

func (tl TransformList) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	s := make([]string, len(tl))
	args := make([]string, 8)
	for i, t := range tl {
		args = args[:0]
		for ia := range t.Args {
			args = append(args, t.Args[ia].String())
		}
		s[i] = t.Name + "(" + strings.Join(args, ",") + ")"
	}
	return makeListAttr(name, s)
}

type Transform struct {
	Name string
	Args []TransformArg
}

type TransformArg interface {
	String() string
}

func (tl *TransformList) TranslateInt(x, y int) *TransformList {
	return tl.append(translateInt(x, y))
}

func translateInt(x, y int) Transform {
	return Transform{Name: "translate", Args: []TransformArg{intArg(x), intArg(y)}}
}

// RotateOrig adds a rotation by the specified number of degrees around
// the origin of the current coordinate system.
func (tl *TransformList) RotateOrig(degrees float64) *TransformList {
	return tl.append(ftrans("rotate", degrees))
}

// SkewX performs a skew transformation along the x axis by the specified angle.
func (tl *TransformList) SkewX(degrees float64) *TransformList {
	return tl.append(ftrans("skewX", degrees))
}

// SkewY performs a skew transformation along the y axis by the specified angle.
func (tl *TransformList) SkewY(degrees float64) *TransformList {
	return tl.append(ftrans("skewY", degrees))
}

func ftrans(name string, f float64) Transform {
	return Transform{Name: name, Args: []TransformArg{floatArg(f)}}
}

type intArg int

func (i intArg) String() string { return strconv.Itoa(int(i)) }

type floatArg float64

func (f floatArg) String() string { return strconv.FormatFloat(float64(f), 'g', -1, 64) }
