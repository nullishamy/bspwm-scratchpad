package main

import (
	"github.com/nullishamy/bspwm-scratchpad/transport"
)

type AddCommand struct {
}

func (a *AddCommand) Run(ctx *Context) error {
	msg := transport.Message{
		Ty: transport.MessageAddCurrentWindow,
	}

	_, err := ctx.client.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}

type RemoveCommand struct {
}

func (a *RemoveCommand) Run(ctx *Context) error {
	msg := transport.Message{
		Ty: transport.MessageRemoveCurrentWindow,
	}

	_, err := ctx.client.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}

type NextCommand struct {
}

func (a *NextCommand) Run(ctx *Context) error {
	msg := transport.Message{
		Ty: transport.MessageShowNextWindow,
	}

	_, err := ctx.client.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}

type PreviousCommand struct {
}

func (a *PreviousCommand) Run(ctx *Context) error {
	msg := transport.Message{
		Ty: transport.MessageShowPreviousWindow,
	}

	_, err := ctx.client.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}

type ShowCommand struct {
}

func (a *ShowCommand) Run(ctx *Context) error {
	msg := transport.Message{
		Ty: transport.MessageShowAllWindows,
	}

	_, err := ctx.client.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}
