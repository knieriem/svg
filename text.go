package svg

import (
	"encoding/xml"
)

type TextAnchor string

type LengthAdjust string

const (
	AnchorMiddle TextAnchor = "middle"
	AnchorEnd    TextAnchor = "end"

	Spacing          LengthAdjust = "spacing"
	SpacingAndGlyphs LengthAdjust = "spacingAndGlyphs"
)

// TextInt places a text element using integer coordinates.
func (el *ElemList) TextInt(x, y int, content string) *TextObject {
	t := &text{TextObject: TextObject{X: float64(x), Y: float64(y)}}
	if content != "" {
		t.Data = append(t.Data, content)
	}
	el.append(t)
	return &t.TextObject
}

type text struct {
	XMLName xml.Name `xml:"text"`
	TextObject
}

// TextObject contains properties common to <text> and <tspan> elements.
type TextObject struct {
	X  float64 `xml:"x,attr,omitempty"`
	Y  float64 `xml:"y,attr,omitempty"`
	Dx Length  `xml:"dx,attr,omitempty"`
	Dy Length  `xml:"dy,attr,omitempty"`

	TextAnchor TextAnchor `xml:"text-anchor,attr,omitempty"`

	TextLength   Length       `xml:"textLength,attr,omitempty"`
	LengthAdjust LengthAdjust `xml:"lengthAdjust,attr,omitempty"`

	Rotate Floats64 `xml:"rotate,attr,omitempty"`

	Object
	Data TextData

	restorePrefix string
	restoreIndent string
}

func (t *TextObject) Anchor(a TextAnchor) *TextObject {
	t.TextAnchor = a
	return t
}

// AddSpan adds a <tspan> element to the parent <text> (or <tspan>) element.
func (t *TextObject) AddSpan(content string) *TextObject {
	ts := new(tspan)
	t.Data = append(t.Data, ts)
	if content != "" {
		ts.Data = append(ts.Data, content)
	}
	ts.restorePrefix = t.restorePrefix
	ts.restoreIndent = t.restoreIndent
	return &ts.TextObject
}

// XMLIndentHint allows the custom XML marshaler for <tspan> to
// temporarily deactivate indentation, to make sure there is no unintended
// white space between the <tspan> tag and the surrounding text.
func (t *TextObject) XMLIndentHint(prefix, indent string) *TextObject {
	t.restorePrefix = prefix
	t.restoreIndent = indent
	return t
}

// AddText adds more text (possibly after a <tspan> element) to a <text> object.
func (t *TextObject) AddText(content string) *TextObject {
	t.Data = append(t.Data, content)
	return t
}

type tspan struct {
	XMLName xml.Name `xml:"tspan"`
	TextObject
}

// TextData is a slice consisting of chardata, or <tspan> elements.
// It is a helper type that implements an xml.Marshaler for proper formatting.
type TextData []interface{}

func (list TextData) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	var err error
	for _, d := range list {
		switch x := d.(type) {
		case string:
			err = e.EncodeToken(xml.CharData(x))
		case *tspan:
			if x.restoreIndent != "" {
				e.Indent("", "")
			}
			err = e.Encode(d)
			if x.restoreIndent != "" {
				e.Indent(x.restorePrefix, x.restoreIndent)
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}
