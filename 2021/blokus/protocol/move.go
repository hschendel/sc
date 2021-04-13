package protocol

import "encoding/xml"

type Piece struct {
	Color     string   `xml:"color,attr"`
	Kind      string   `xml:"kind,attr"`
	Rotation  string   `xml:"rotation,attr"`
	IsFlipped bool     `xml:"isFlipped,attr"`
	Position  Position `xml:"position"`
}

type Position struct {
	X       uint8    `xml:"x,attr"`
	Y       uint8    `xml:"y,attr"`
}

type Hint struct {
	XMLName xml.Name `xml:"hint"`
	Content string   `xml:"content,attr"`
}