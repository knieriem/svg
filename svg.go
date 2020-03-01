// Package svg provides support for creating SVG documents.
package svg

import (
	"strconv"
	"strings"

	"encoding/xml"
)

const (
	nameSpace = "http://www.w3.org/2000/svg"
)

type Conf struct {
	// GenerateEmbeededStylesheet automatically copies
	// styles defined with UseStyle into a <style> tag
	// embedded into the SVG document,
	// with class names adjusted as needed.
	GenerateEmbeddedStylesheet bool

	// StylesheetUnifyStyles makes sure each style only occurs
	// once in the generated Stylesheet.
	StylesheetUnifyStyles bool

	// ScopeStyleDefinitions makes sure that classes
	// defined within the embedded stylesheet are
	// valid within the SVG document only, by inserting an ID selector
	// in front of each definition.
	// The purpose is to avoid side-effects when using multiple
	// generated SVG documents in one HTML document.
	// If set, Document.ID must be set to a value too.
	ScopeStyleDefinitions bool

	// Embedded, if set, makes sure that the SVG 'xmlns' attribute
	// is left out of the generated SVG.
	Embedded bool
}

// Document contains the SVG document.
type Document struct {
	XMLName xml.Name `xml:"svg"`

	ViewBox Ints   `xml:"viewBox,attr,omitempty"`
	Width   Length `xml:"width,attr,omitempty"`
	Height  Length `xml:"height,attr,omitempty"`

	Style string `xml:"style,omitempty"`

	Container

	styles struct {
		defMap    map[string]string
		classMap  map[string]string
		nConflict int
	}

	NameSpace string `xml:"xmlns,attr,omitempty"`
	conf      *Conf
}

// NewDocument creates an empty SVG document.
// Adjust ViewBox, Width and Height as needed.
func NewDocument(c *Conf) *Document {
	d := new(Document)
	if c == nil {
		c = &Conf{}
	}
	if !c.Embedded {
		d.NameSpace = nameSpace
	}
	d.conf = c
	return d
}

// MakeStyle returns a Styling that may be applied to stylable
// objects using the WithStyle method.
// If Conf.GenerateEmbeddedStylesheet is set, style
// definitions are appended to the document's Style field,
// and a Styling is returned specifying only a class name.
// Otherwise the returned Styling will result in an explicit
// style attribute value, if applied to an object, and the name
// won't be used.
func (d *Document) MakeStyle(name, style string) Styling {
	if !d.conf.GenerateEmbeddedStylesheet {
		if style != "" {
			return Styling{Style: style}
		}
		return Styling{Class: name}
	}

	s := &d.styles
	if s.defMap == nil {
		s.defMap = make(map[string]string, 16)
		s.classMap = make(map[string]string, 16)
	}
	class, styleExists := s.defMap[style]
	if !styleExists {
		if _, exists := s.classMap[name]; exists {
			s.nConflict++
			name += strconv.Itoa(s.nConflict)
		}
		if d.conf.StylesheetUnifyStyles {
			s.defMap[style] = name
		}
		s.classMap[name] = style
		class = name

		// update style
		if d.Style != "" {
			d.Style += " "
		}
		if d.conf.ScopeStyleDefinitions && d.ID != "" {
			d.Style += "#" + d.ID + " "
		}
		d.Style += "." + name + " {" + strings.TrimSuffix(style, ";") + "}"
	}
	return Styling{Class: class}
}

type Styling struct {
	Class string `xml:"class,attr,omitempty"`
	Style string `xml:"style,attr,omitempty"`
}

func (st *Styling) SetStyle(style string) *Styling {
	st.Style = strings.TrimSuffix(style, ";")
	return st
}

func (st *Styling) SetClass(class string) *Styling {
	st.Class = class
	return st
}

func (st *Styling) WithStyle(s Styling) *Styling {
	*st = s
	return st
}

type Stylable interface {
	SetClass(string) *Styling
	SetStyle(string) *Styling
	WithStyle(s Styling) *Styling
}

// ElemList is a slice of SVG elements embedded into the
// document container, or into group containers.
type ElemList []interface{}

func (el *ElemList) append(i interface{}) {
	*el = append(*el, i)
}

func (el *ElemList) UseObjectInt(x, y int, id string) *Object {
	u := &use{X: float64(x), Y: float64(y), Href: "#" + id}
	el.append(u)
	return &u.Object
}

type use struct {
	XMLName xml.Name `xml:"use"`
	X       float64  `xml:"x,attr,omitempty"`
	Y       float64  `xml:"y,attr,omitempty"`
	Href    string   `xml:"href,attr,omitempty"`
	Object
}

// Container contains child elements. It may be styled and transformed.
type Container struct {
	Object
	ElemList `xml:",omitempty"`
}

// Group is used as a container to group other SVG elements.
type Group struct {
	XMLName xml.Name `xml:"g"`
	Container
}

// Defs is used as a container to group other SVG elements
// that won't be displayed initially.
type Defs struct {
	XMLName xml.Name `xml:"defs"`
	Container
}

// Defs appends a defs element.
func (el *ElemList) Defs() *Container {
	g := new(Defs)
	el.append(g)
	return &g.Container
}

// Group appends a group element.
func (el *ElemList) Group() *Container {
	g := new(Group)
	el.append(g)
	return &g.Container
}

// PreAlloc preallocates memory for the given number of elements.
func (c *Container) PreAlloc(n int) *Container {
	if c.ElemList == nil {
		c.ElemList = make(ElemList, 0, n)
	}
	return c
}

// SetID sets the ID field of the embedded Object
func (c *Container) SetID(id string) *Container {
	c.Object.SetID(id)
	return c
}

// An Object may be styled and transformed.
type Object struct {
	ID            string `xml:"id,attr,omitempty"`
	TransformList `xml:"transform,attr,omitempty"`
	Styling
}

func (o *Object) SetID(id string) *Object {
	o.ID = id
	return o
}

// Title appends a title element.
func (el *ElemList) Title(content string) {
	t := &title{Data: content}
	el.append(t)
}

type title struct {
	XMLName xml.Name `xml:"title"`
	Data    string   `xml:",chardata"`
}

// Ints is a slice of integers that marshals, if used as an XML
// attribute value, into a list of space separated decimal string
// representations of these integers.
type Ints []int

func (ints Ints) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	s := make([]string, len(ints))
	for i, v := range ints {
		s[i] = strconv.Itoa(v)
	}
	return makeListAttr(name, s)
}

// Floats64 is a slice of float64 values that marshals, if used as an XML
// attribute value, into a list of space separated string representations
// of these float64 values.
type Floats64 []float64

func (f Floats64) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	s := make([]string, len(f))
	for i, v := range f {
		s[i] = strconv.FormatFloat(v, 'g', -1, 64)
	}
	return makeListAttr(name, s)
}

func makeListAttr(name xml.Name, values []string) (xml.Attr, error) {
	var a xml.Attr
	a.Name = name
	a.Value = strings.Join(values, " ")
	return a, nil
}

// Length may be a value with a unit, a percentage, or a number.
type Length interface {
	xml.MarshalerAttr
}

// Number returns a value that will be marshaled without a unit.
func Number(f float64) Length {
	return number(f)
}

type number float64

func (n number) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return marshalLengthAttr(name, float64(n), "")
}

// EmUnits returns a Length that will be marshaled with an "em" suffix.
func EmUnits(f float64) Length {
	return emUnits(f)
}

type emUnits float64

func (u emUnits) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return marshalLengthAttr(name, float64(u), "em")
}

// ExUnits returns a Length that will be marshaled with an "ex" suffix.
func ExUnits(f float64) Length {
	return exUnits(f)
}

type exUnits float64

func (u exUnits) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return marshalLengthAttr(name, float64(u), "ex")
}

// Percentage returns a Length that will be marshaled with a "%" suffix.
func Percentage(f float64) Length {
	return percentage(f)
}

type percentage float64

func (p percentage) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return marshalLengthAttr(name, float64(p), "%")
}

func marshalLengthAttr(name xml.Name, f float64, unit string) (xml.Attr, error) {
	var a xml.Attr
	a.Name = name
	a.Value = strconv.FormatFloat(f, 'g', -1, 64) + unit
	return a, nil
}
