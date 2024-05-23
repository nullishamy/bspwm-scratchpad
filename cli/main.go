package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/alecthomas/kong"
	"github.com/nullishamy/bspwm-scratchpad/transport"
)

type Client struct {
	con   net.Conn
	msgId int
}

func (c *Client) SendMessage(message transport.Message) (*transport.Message, error) {
	message.Id = c.msgId
	c.msgId += 1

	msgBytes, err := transport.EncodeMessage(message)
	if err != nil {
		return nil, err
	}

	_, err = c.con.Write(msgBytes)
	if err != nil {
		return nil, err
	}

	reply, err := transport.DecodeMessage(c.con)
	if err != nil {
		return nil, err
	}

	if reply.Ty == transport.MessageError {
		data := transport.ErrorMessage{}
		err := json.Unmarshal(reply.Data, &data)
		if err != nil {
			return nil, err
		}

		return nil, errors.New("error response: " + data.Details)
	}

	return reply, nil
}

type Globals struct {
	Socket string `help:"Location of daemon socket" default:"/tmp/scratch.sock" type:"path"`
}

type CLI struct {
	Globals
	Add    AddCommand    `cmd:"" help:"Add a window to the scratchpad."`
	Remove RemoveCommand `cmd:"" help:"Remove a window from the scratchpad."`
	Next NextCommand `cmd:"" help:"Show the next window."`
	Previous PreviousCommand `cmd:"" help:"Show previous window."`
	Show ShowCommand `cmd:"" help:"Show all windows."`
}

type Context struct {
	client  *Client
	globals *Globals
}

func main() {
	cli := CLI{}
	ctx := kong.Parse(&cli, kong.UsageOnError())

	c, err := net.Dial("unix", transport.DEFAULT_SOCK_PATH)
	if err != nil {
		// Exit nicely so we don't mess up scripts if the socket is down
		fmt.Println("error when dialing socket"+err.Error())
		return
	}

	client := Client{
		con:   c,
		msgId: 1,
	}

	err = ctx.Run(&Context{client: &client, globals: &cli.Globals})
	ctx.FatalIfErrorf(err)
}
