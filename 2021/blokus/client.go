package blokus

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/hschendel/sc"
	"github.com/hschendel/sc/2021/blokus/protocol"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

type Client struct {
	Conn            net.Conn
	Player          Player
	MoveTimeout     time.Duration
	State           MutableState
	ReservationCode string
	DebugTo         *os.File
}

func OpenClient(address net.Addr, player Player, emptyState MutableState) (cl *Client, err error) {
	if player == nil {
		panic(errors.New("player is nil"))
	}
	if emptyState == nil {
		panic(errors.New("emptyState is nil"))
	}
	var conn net.Conn
	conn, err = net.Dial(address.Network(), address.String())
	if err != nil {
		return
	}
	cl = &Client{
		Conn:            conn,
		Player:          player,
		MoveTimeout:     DefaultMoveTimeout,
		State:           emptyState,
		ReservationCode: "",
	}
	return
}

var DefaultServerAddress net.Addr = &net.TCPAddr{
	IP:   net.IPv4(127, 0, 0, 1),
	Port: DefaultServerPort,
}

const DefaultServerPort = 13050
const DefaultMoveTimeout = 2 * time.Second

func (c *Client) Run() (err error) {
	defer c.Player.End()
	moveTimeout := c.MoveTimeout
	if moveTimeout == 0 {
		moveTimeout = DefaultMoveTimeout
	}
	var xc xmlConn
	xc.init(c.Conn, c.DebugTo)
	if err = xc.sendBytes(protocol.ProtocolMessage); err != nil {
		return
	}
	var roomID string
	var isFirstPlayer bool
	if roomID, isFirstPlayer, err = xc.join(c.ReservationCode, c.State); err != nil {
		err = fmt.Errorf("cannot join game: %s", err)
		return
	}
	log.Printf("joined room %q, firstPlayer=%v", roomID, isFirstPlayer)
	var colors [2]Color
	if isFirstPlayer {
		colors[0] = ColorBlue
		colors[1] = ColorRed
	} else {
		colors[0] = ColorYellow
		colors[1] = ColorGreen
	}
	var gameEnded bool
	var validColors map[Color]bool

	for colorIdx := 0; ; colorIdx = (colorIdx + 1) % len(colors) {
		if gameEnded, validColors, err = xc.waitForMoveRequest(roomID, c.State); err != nil || gameEnded {
			return
		}
		playerColor := colors[colorIdx]
		if !validColors[playerColor] {
			// advance to next color if playerColor is not valid
			colorIdx = (colorIdx + 1) % len(colors)
			playerColor = colors[colorIdx]
			if !validColors[playerColor] {
				// wait for next MoveRequest if no color is valid
				continue
			}
		}
		move := c.Player.NextMove(c.State, playerColor, sc.NewTimeout(moveTimeout))
		if err = xc.sendMove(roomID, colors[colorIdx], move); err != nil {
			err = fmt.Errorf("cannot send move: %s", err)
			return
		}
	}
}

type xmlConn struct {
	enc *xml.Encoder
	dec *xml.Decoder
	w   io.Writer
	r   io.Reader
}

func (x *xmlConn) init(conn net.Conn, debugTo *os.File) {
	if debugTo != nil {
		x.w = &loggingWriter{
			Out: debugTo,
			W:   conn,
		}
		x.r = &loggingReader{
			Out: debugTo,
			R:   conn,
		}
	} else {
		x.w = conn
		x.r = conn
	}
	x.enc = xml.NewEncoder(x.w)
	x.dec = xml.NewDecoder(x.r)
}

func (x *xmlConn) send(v interface{}) error {
	if err := x.enc.Encode(v); err != nil {
		return err
	}
	return nil
}

func (x *xmlConn) sendBytes(p []byte) error {
	_, err := x.w.Write(p)
	if err == nil {
		err = x.enc.Flush()
	}
	return err
}

func (x *xmlConn) receive(v interface{}) error {
	return x.dec.Decode(v)
}

func (x *xmlConn) join(reservationCode string, intoState MutableState) (roomID string, isFirstPlayer bool, err error) {
	if reservationCode != "" {
		if err = x.sendBytes(protocol.JoinPreparedMessage(reservationCode)); err != nil {
			return
		}
	} else {
		if err = x.sendBytes(protocol.JoinMessage); err != nil {
			return
		}
	}
	if err = x.expectProtocol(); err != nil {
		return
	}
	var joined protocol.Joined
	if err = x.receive(&joined); err != nil {
		return
	}
	roomID = joined.RoomID

	var welcomeMessage protocol.Room
	if err = x.receive(&welcomeMessage); err != nil {
		return
	}
	isFirstPlayer = welcomeMessage.Data.ColorAttr == "ONE"

	var stateInRoom protocol.Room
	if err = x.receive(&stateInRoom); err != nil {
		return
	}
	if _, err = fillState(intoState, stateInRoom.Data.State); err != nil {
		return
	}
	return
}

func (x *xmlConn) expectProtocol() error {
	buf := make([]byte, len(protocol.ProtocolMessage))
	n := 0
	for n < len(buf) {
		nn, err := x.r.Read(buf[n:])
		if err != nil {
			return err
		}
		n += nn
	}
	if !bytes.Equal(buf, protocol.ProtocolMessage) {
		return fmt.Errorf("expected %q but got %q", string(protocol.ProtocolMessage), string(buf))
	}
	return nil
}

func (x *xmlConn) waitForMoveRequest(roomID string, intoState MutableState) (gameEnded bool, validColors map[Color]bool, err error) {
	for {
		var room protocol.Room
		if err = x.receive(&room); err != nil {
			return
		}
		if room.RoomID != roomID {
			err = fmt.Errorf("expected roomId to be %q but got %q", roomID, room.RoomID)
			return
		}
		switch room.Data.Class {
		case protocol.DataClassState:
			validColors, err = fillState(intoState, room.Data.State)
			if err != nil {
				return
			}
		case protocol.DataClassResult:
			gameEnded = true
			return
		case protocol.DataClassMoveRequest:
			return
		}
	}
}

func (x *xmlConn) sendMove(roomID string, color Color, move Move) error {
	var room protocol.Room
	room.RoomID = roomID
	if move.IsEmpty() {
		room.Data.Class = protocol.DataClassSkipMove
		room.Data.ColorField = color.String()
	} else {
		room.Data.Class = protocol.DataClassSetMove
		room.Data.Piece = &protocol.Piece{
			Color:     color.String(),
			Kind:      move.Transformation.Piece().String(),
			Rotation:  move.Transformation.Rotation().String(),
			IsFlipped: move.Transformation.Flipped(),
			Position:  protocol.Position{X: move.X, Y: move.Y},
		}
	}
	return x.send(&room)
}

func fillState(s MutableState, xs *protocol.State) (validColors map[Color]bool, err error) {
	if xs == nil {
		panic(errors.New("xs is nil (state not set)"))
	}
	validColors = make(map[Color]bool, 4)
	s.Reset()
	if startPiece, parseErr := ParsePiece(xs.StartPiece); parseErr != nil {
		err = fmt.Errorf("cannot parse startPiece value %q: %s", xs.StartPiece, parseErr)
		return
	} else {
		s.SetStartPiece(startPiece)
	}
	if err = setPlayedPieces(s, xs); err != nil {
		return
	}
	for _, entry := range xs.LastMoveMono {
		c, parseErr := ParseColor(entry.Color)
		if parseErr != nil {
			err = fmt.Errorf("cannot parse lastMoveMono entry color value %q: %s", entry.Color, parseErr)
			return
		}
		s.SetLastMoveMono(c, entry.Boolean)
	}
	for _, field := range xs.Board {
		c, parseErr := ParseColor(field.Content)
		if parseErr != nil {
			err = fmt.Errorf("cannot parse color content attribute value %q of board field x=%d, y=%d: %s", field.Content, field.X, field.Y, parseErr)
			return
		}
		s.Set(field.X, field.Y, c, true)
	}
	for _, colorStr := range xs.ValidColors {
		c, parseErr := ParseColor(colorStr)
		if parseErr != nil {
			err = fmt.Errorf("cannot parse color content value %q of validColors: %s", colorStr, parseErr)
			return
		}
		validColors[c] = true
	}
	return
}

func setPlayedPieces(s MutableState, xs *protocol.State) error {
	if err := setPlayedPiecesForColor(s, ColorBlue, xs.BlueShapes); err != nil {
		return err
	}
	if err := setPlayedPiecesForColor(s, ColorYellow, xs.YellowShapes); err != nil {
		return err
	}
	if err := setPlayedPiecesForColor(s, ColorRed, xs.RedShapes); err != nil {
		return err
	}
	if err := setPlayedPiecesForColor(s, ColorGreen, xs.GreenShapes); err != nil {
		return err
	}
	return nil
}

func setPlayedPiecesForColor(s MutableState, c Color, sl []string) error {
	pieces, err := parsePieces(sl)
	if err != nil {
		return err
	}
	s.SetNotPlayedPiecesFor(c, pieces)
	return nil
}

func parsePieces(sl []string) (pieces []Piece, err error) {
	for _, s := range sl {
		var p Piece
		p, err = ParsePiece(s)
		if err != nil {
			return
		}
		pieces = append(pieces, p)
	}
	return
}

// ClientMain can be used to implement a simple player main function.
// It uses os.Args
func ClientMain(player Player) {
	addr := DefaultServerAddress
	if len(os.Args) > 1 {
		port, err := strconv.Atoi(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid value for port: %q", os.Args[1])
		}
		addr = &net.TCPAddr{
			IP:   net.IPv4(127, 0, 0, 1),
			Port: port,
		}
	}
	var state BasicState
	client, err := OpenClient(addr, player, &state)
	if err != nil {
		fmt.Fprintf(os.Stderr, "cannot connect to server: %s", err)
		os.Exit(1)
	}
	client.DebugTo = os.Stderr
	err = client.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while running: %s", err)
		os.Exit(2)
	}
}
