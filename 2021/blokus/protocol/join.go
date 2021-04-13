package protocol

import (
	"encoding/xml"
	"fmt"
)

var ProtocolMessage = []byte("<protocol>\n")
var JoinMessage = []byte(fmt.Sprintf("<join gameType=%q />\n", GameTypeBlokus))

const GameTypeBlokus = "swc_2021_blokus"

func JoinPreparedMessage(reservationCode string) []byte {
	return []byte(fmt.Sprintf("<joinPrepared reservationCode=%q />\n", reservationCode))
}

type Joined struct {
	XMLName xml.Name `xml:"joined"`
	RoomID  string   `xml:"roomId,attr"`
}




