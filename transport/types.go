package transport

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"math"
	"net"
)

const DEFAULT_SOCK_PATH = "/tmp/scratch.sock"

type MessageType int

const (
	MessageHello MessageType = iota
	MessageCurrentWindow MessageType = iota
	MessageAddCurrentWindow MessageType = iota
	MessageRemoveCurrentWindow MessageType = iota
	MessageShowNextWindow MessageType = iota
	MessageShowPreviousWindow MessageType = iota
	MessageShowAllWindows MessageType = iota
	MessageSetWindowVisibility = iota
	MessageError MessageType = iota
)

type Message struct {
	Ty   MessageType `json:"type"`
	Id   int   `json:"id"`
	Data json.RawMessage `json:"data"`
}

type CurrentWindowMessage struct {
	Window Window
}

type SetVisibilityMessage struct {
	ID int64
	NewVisibility bool
}

type ErrorMessage struct {
	Details string
}

type Window struct {
	ID          int64       `json:"id"`
	SplitType   string      `json:"splitType"`
	SplitRatio  float64     `json:"splitRatio"`
	Vacant      bool        `json:"vacant"`
	Hidden      bool        `json:"hidden"`
	Sticky      bool        `json:"sticky"`
	Private     bool        `json:"private"`
	Locked      bool        `json:"locked"`
	Marked      bool        `json:"marked"`
	Presel      interface{} `json:"presel"`
	Rectangle   Rectangle   `json:"rectangle"`
	Constraints Constraints `json:"constraints"`
	FirstChild  interface{} `json:"firstChild"`
	SecondChild interface{} `json:"secondChild"`
	Client      Client      `json:"client"`
}

type Client struct {
	ClassName         string    `json:"className"`
	InstanceName      string    `json:"instanceName"`
	BorderWidth       int64     `json:"borderWidth"`
	State             string    `json:"state"`
	LastState         string    `json:"lastState"`
	Layer             string    `json:"layer"`
	LastLayer         string    `json:"lastLayer"`
	Urgent            bool      `json:"urgent"`
	Shown             bool      `json:"shown"`
	TiledRectangle    Rectangle `json:"tiledRectangle"`
	FloatingRectangle Rectangle `json:"floatingRectangle"`
}

type Rectangle struct {
	X      int64 `json:"x"`
	Y      int64 `json:"y"`
	Width  int64 `json:"width"`
	Height int64 `json:"height"`
}

type Constraints struct {
	MinWidth  int64 `json:"min_width"`
	MinHeight int64 `json:"min_height"`
}

func EncodeMessage(message Message) ([]byte, error) {
	var buf bytes.Buffer

	msg, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	msgLen := int64(len(msg))
	if msgLen > math.MaxUint32 {
		return nil, errors.New("message too long to encode length in 4 bytes")
	}

	lenBuf := []byte{0, 0, 0, 0}
	binary.LittleEndian.PutUint32(lenBuf[:], uint32(msgLen))

	buf.Write(lenBuf[:])
	buf.Write(msg)

	return buf.Bytes(), nil
}

func DecodeMessage(con net.Conn) (*Message, error) {
	lenBuf := make([]byte, 4)
	nr, err := con.Read(lenBuf)
	if err != nil {
		return nil, err
	}

	if nr != 4 {
		return nil, errors.New("did not get the expected 4 bytes for len")
	}

	len := binary.LittleEndian.Uint32(lenBuf)
	jsonBuf := make([]byte, len)
	nr, err = con.Read(jsonBuf)

	if err != nil {
		return nil, err
	}

	for uint32(nr) != len {
		nr, err = con.Read(jsonBuf)

		if err != nil {
			return nil, err
		}
	}

	var message Message
	err = json.Unmarshal(jsonBuf, &message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}
