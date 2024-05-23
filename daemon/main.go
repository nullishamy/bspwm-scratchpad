package main

import (
	"fmt"
	"github.com/nullishamy/bspwm-scratchpad/transport"
	"net"
    "os"
    "os/signal"
    "syscall"
)

type Server struct {
	handlers map[transport.MessageType]Handler
	windows  []int64

	// Index into `windows`
	currentWindow int
}

func (s Server) Install(ty transport.MessageType, h Handler) {
	s.handlers[ty] = h
}

func handleConnection(c net.Conn, s *Server) {
	for {
		message, err := transport.DecodeMessage(c)
		if err != nil {
			return
		}

		fmt.Printf("Received: %+v\n", message)

		handler := s.handlers[message.Ty]

		if handler == nil {
			panic("Missing handler for type " + fmt.Sprint(message.Ty))
		}

		req := Request{
			message: message,
			con:     c,
			server:  s,
		}

		reply, err := handler.Execute(req)
		if err != nil {
			panic("execute: " + err.Error())
		}

		replyBytes, err := transport.EncodeMessage(reply.message)
		if err != nil {
			panic("encode: " + err.Error())
		}

		_, err = c.Write(replyBytes)
		if err != nil {
			panic("write: " + err.Error())
		}
	}
}

func main() {
	sockPath := transport.DEFAULT_SOCK_PATH

	l, err := net.Listen("unix", sockPath)
	if err != nil {
		println("listen error", err.Error())
		return
	}

	fmt.Printf("listening on %s\n", sockPath)

	server := Server{
		handlers:      make(map[transport.MessageType]Handler),
		windows:       []int64{},
		currentWindow: 0,
	}

	server.Install(transport.MessageHello, HelloHandler{})
	server.Install(transport.MessageCurrentWindow, CurrentWindowHandler{})
	server.Install(transport.MessageSetWindowVisibility, VisibilityHandler{})
	server.Install(transport.MessageAddCurrentWindow, AddCurrentWindowHandler{})
	server.Install(transport.MessageRemoveCurrentWindow, RemoveCurrentWindowHandler{})
	server.Install(transport.MessageShowNextWindow, ShowNextWindowHandler{})
	server.Install(transport.MessageShowPreviousWindow, ShowPreviousWindowHandler{})
	server.Install(transport.MessageShowAllWindows, ShowAllHandler{})

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	go func() {
		sig := <-sigs
        if sig == syscall.SIGINT {
            l.Close()
            os.Remove(sockPath)
        }
		done <- true
	}()

	for {
		con, err := l.Accept()
		if err != nil {
			println("\nerror when accepting connection:", err.Error())
			return
		}

		go handleConnection(con, &server)
	}
}
