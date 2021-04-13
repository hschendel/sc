package protocol

import "encoding/xml"

type State struct {
	Turn          uint     `xml:"turn,attr"`
	Round         uint     `xml:"round,attr"`
	StartPiece    string   `xml:"startPiece,attr"`
	StartTeam     StartTeam
	BlueShapes    []string     `xml:"blueShapes>shape"`
	YellowShapes  []string     `xml:"yellowShapes>shape"`
	RedShapes     []string     `xml:"redShapes>shape"`
	GreenShapes   []string     `xml:"greenShapes>shape"`
	ValidColors   []string     `xml:"validColors>color"`
	FirstTeam     Team         `xml:"first"`
	SecondTeam    Team         `xml:"second"`
	Board         []Field      `xml:"board>field"`
	LastMoveMono  []ColorEntry `xml:"lastMoveMono>entry,omitempty"`
}

type StartTeam struct {
	XMLName xml.Name `xml:"startTeam"`
	Class   string   `xml:"class,attr"`
	Name    string   `xml:",chardata"`
}

type Team struct {
	DisplayName string `xml:"displayName,attr"`
	Color       string `xml:"color"`
}

type Field struct {
	X       uint8  `xml:"x,attr"`
	Y       uint8  `xml:"y,attr"`
	Content string `xml:"content,attr"`
}

type ColorEntry struct {
	Color   string `xml:"color"`
	Boolean bool   `xml:"boolean"`
}
