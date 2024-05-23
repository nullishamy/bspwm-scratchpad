package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"

	"github.com/nullishamy/bspwm-scratchpad/transport"
)

type Request struct {
	message *transport.Message
	server  *Server
	con     net.Conn
}

func (r Request) Id() int {
	return r.message.Id
}

type Response struct {
	request *Request
	message transport.Message
}

func errorReponse(req *Request, e error) Response {
	data := transport.ErrorMessage{
		Details: e.Error(),
	}

	bytes, err := json.Marshal(data)

	if err != nil {
		panic("Failed to marshal response message")
	}

	return Response{
		request: req,
		message: transport.Message{
			Ty:   transport.MessageError,
			Id:   req.Id(),
			Data: bytes,
		},
	}
}

type Handler interface {
	Execute(req Request) (Response, error)
}

type HelloHandler struct{}

func (h HelloHandler) Execute(req Request) (Response, error) {
	res := Response{
		request: &req,
		message: transport.Message{
			Ty: transport.MessageHello,
			Id: req.Id(),
		},
	}

	return res, nil
}

type CurrentWindowHandler struct{}

func (h CurrentWindowHandler) Execute(req Request) (Response, error) {
	currentWindow, err := GetCurrentWindow()
	if err != nil {
		return errorReponse(&req, errors.New("failed to get window "+err.Error())), err
	}

	bytes, err := json.Marshal(
		transport.CurrentWindowMessage{
			Window: *currentWindow,
		},
	)

	if err != nil {
		return errorReponse(&req, errors.New("failed to get marshal response "+err.Error())), err
	}

	res := Response{
		request: &req,
		message: transport.Message{
			Ty:   transport.MessageCurrentWindow,
			Id:   req.Id(),
			Data: bytes,
		},
	}

	return res, nil
}

type AddCurrentWindowHandler struct{}

func (h AddCurrentWindowHandler) Execute(req Request) (Response, error) {
	currentWindow, err := GetCurrentWindow()
	if err != nil {
		return errorReponse(&req, errors.New("failed to get window "+err.Error())), err
	}

	if !Contains(req.server.windows, currentWindow.ID) {
		req.server.windows = append(req.server.windows, currentWindow.ID)
	}

	fmt.Printf("windows: %+v\n", req.server.windows)

	// Only hide if we have other windows controlled, ie there's another window on display
	if len(req.server.windows) > 1 {
		if err := HideWindow(currentWindow); err != nil {
			return errorReponse(&req, errors.New("failed to hide window "+err.Error())), err
		}
	}

	res := Response{
		request: &req,
		message: transport.Message{
			Ty: transport.MessageAddCurrentWindow,
			Id: req.Id(),
		},
	}

	return res, nil
}

type RemoveCurrentWindowHandler struct{}

func (h RemoveCurrentWindowHandler) Execute(req Request) (Response, error) {
	windowIdx := req.server.currentWindow

	currentWindow, err := GetCurrentWindow()
	currentWindowId := currentWindow.ID

	if err != nil {
		return errorReponse(&req, errors.New("failed to get window "+err.Error())), err
	}

	res := Response{
		request: &req,
		message: transport.Message{
			Ty: transport.MessageRemoveCurrentWindow,
			Id: req.Id(),
		},
	}

	req.server.windows = Remove(req.server.windows, currentWindowId)
	fmt.Printf("windows: %+v\n", req.server.windows)
	fmt.Printf("window id: %+v\n", currentWindowId)
	windows := req.server.windows

	// The last window is being killed, do nothing
	if len(windows) == 0 {
		return res, nil
	}

	// We are now only tracking a single window, show it
	if len(windows) == 1 {
		newWindow, err := GetWindowDetails(windows[0])

		if err != nil {
			return errorReponse(&req, errors.New("failed to get window "+err.Error())), err
		}

		ShowWindow(newWindow)

		return res, nil
	}

	// Wrap around
	var nextWindowIdx int
	if windowIdx+1 >= len(windows) {
		nextWindowIdx = 0
	} else {
		nextWindowIdx = windowIdx + 1
	}

	nextWindowId := windows[nextWindowIdx]
	nextWindow, err := GetWindowDetails(nextWindowId)

	if err != nil {
		return errorReponse(&req, errors.New("failed to get window "+err.Error())), err
	}

	err = ShowWindow(nextWindow)
	if err != nil {
		return errorReponse(&req, errors.New("failed to show next window "+err.Error())), err
	}

	req.server.currentWindow = nextWindowIdx

	return res, nil
}

type ShowPreviousWindowHandler struct{}

func (h ShowPreviousWindowHandler) Execute(req Request) (Response, error) {
	windows := req.server.windows
	windowIdx := req.server.currentWindow

	res := Response{
		request: &req,
		message: transport.Message{
			Ty: transport.MessageShowNextWindow,
			Id: req.Id(),
		},
	}

	// No tracked windows, noop
	if len(windows) == 0 {
		return res, nil
	}

	currentWindowId := windows[windowIdx]
	currentWindow, err := GetWindowDetails(currentWindowId)

	if err != nil {
		return errorReponse(&req, errors.New("failed to get window "+err.Error())), err
	}

	// Single tracked window, always show it
	if len(windows) == 1 {
		ShowWindow(currentWindow)
		return res, nil
	}

	// Wrap around
	var previousWindowIdx int
	if windowIdx-1 < 0 {
		previousWindowIdx = len(windows) - 1
	} else {
		previousWindowIdx = windowIdx - 1
	}

	previousWindowId := windows[previousWindowIdx]
	previousWindow, err := GetWindowDetails(previousWindowId)

	if err != nil {
		return errorReponse(&req, errors.New("failed to get window "+err.Error())), err
	}

	err = HideWindow(currentWindow)
	if err != nil {
		return errorReponse(&req, errors.New("failed to hide current window "+err.Error())), err
	}

	err = ShowWindow(previousWindow)
	if err != nil {
		return errorReponse(&req, errors.New("failed to show previous window "+err.Error())), err
	}

	req.server.currentWindow = previousWindowIdx

	return res, nil
}

type ShowNextWindowHandler struct{}

func (h ShowNextWindowHandler) Execute(req Request) (Response, error) {
	windows := req.server.windows
	windowIdx := req.server.currentWindow

	res := Response{
		request: &req,
		message: transport.Message{
			Ty: transport.MessageShowNextWindow,
			Id: req.Id(),
		},
	}

	// No tracked windows, noop
	if len(windows) == 0 {
		return res, nil
	}

	currentWindowId := windows[windowIdx]
	currentWindow, err := GetWindowDetails(currentWindowId)

	if err != nil {
		return errorReponse(&req, errors.New("failed to get window "+err.Error())), err
	}

	// Single tracked window, always show it
	if len(windows) == 1 {
		ShowWindow(currentWindow)
		return res, nil
	}

	// Wrap around
	var nextWindowIdx int
	if windowIdx+1 >= len(windows) {
		nextWindowIdx = 0
	} else {
		nextWindowIdx = windowIdx + 1
	}

	nextWindowId := windows[nextWindowIdx]
	nextWindow, err := GetWindowDetails(nextWindowId)

	if err != nil {
		return errorReponse(&req, errors.New("failed to get window "+err.Error())), err
	}

	err = HideWindow(currentWindow)
	if err != nil {
		return errorReponse(&req, errors.New("failed to hide current window "+err.Error())), err
	}

	err = ShowWindow(nextWindow)
	if err != nil {
		return errorReponse(&req, errors.New("failed to show next window "+err.Error())), err
	}

	req.server.currentWindow = nextWindowIdx

	return res, nil
}

type VisibilityHandler struct{}

func (h VisibilityHandler) Execute(req Request) (Response, error) {
	data := transport.SetVisibilityMessage{}
	err := json.Unmarshal(req.message.Data, &data)

	if err != nil {
		return errorReponse(&req, errors.New("failed to unmarshal data "+err.Error())), err
	}

	window, err := GetWindowDetails(data.ID)
	if err != nil {
		return errorReponse(&req, errors.New("failed to get window "+err.Error())), err
	}

	if data.NewVisibility {
		ShowWindow(window)
	} else {
		HideWindow(window)
	}

	res := Response{
		request: &req,
		message: transport.Message{
			Ty: transport.MessageSetWindowVisibility,
			Id: req.Id(),
		},
	}

	return res, nil
}

type ShowAllHandler struct{}

func (h ShowAllHandler) Execute(req Request) (Response, error) {
	for _, id := range req.server.windows {
		window, err := GetWindowDetails(id)
		if err != nil {
			return errorReponse(&req, errors.New("failed to get window "+err.Error())), err
		}

		err = ShowWindow(window)
		if err != nil {
			return errorReponse(&req, errors.New("failed to show window "+err.Error())), err
		}
	}

	res := Response{
		request: &req,
		message: transport.Message{
			Ty: transport.MessageShowAllWindows,
			Id: req.Id(),
		},
	}

	return res, nil
}
