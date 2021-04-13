package protocol

import "encoding/xml"

type Room struct {
	XMLName xml.Name `xml:"room"`
	RoomID  string   `xml:"roomId,attr"`
	Data    Data     `xml:"data,omitEmpty"`
}

type Data struct {
	Class      string `xml:"class,attr,omitempty"`
	ColorAttr  string `xml:"color,attr,omitempty"`
	ColorField string `xml:"color,omitempty"`
	Piece      *Piece `xml:"piece,omitempty"`
	State      *State `xml:"state,omitempty"`
}

const DataClassState = "memento"
const DataClassResult = "result"
const DataClassMoveRequest = "sc.framework.plugins.protocol.MoveRequest"
const DataClassSetMove = "sc.plugin2021.SetMove"
const DataClassSkipMove = "sc.plugin2021.SkipMove"
