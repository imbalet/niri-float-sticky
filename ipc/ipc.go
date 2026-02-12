package ipc

import (
	"encoding/json"
	"io"
	"net"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type Command struct {
	Action   string `json:"action"`
	WindowID uint64 `json:"window_id"`
}

func SocketPath() string {
	return filepath.Join(os.Getenv("XDG_RUNTIME_DIR"), "niri-float-sticky.sock")
}

func StartIPC(cmdChan chan<- Command) {
	socketPath := SocketPath()
	_ = os.Remove(socketPath)

	l, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Panic(err)
	}

	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				continue
			}

			go handleConn(conn, cmdChan)
		}
	}()
}

func handleConn(c net.Conn, cmdChan chan<- Command) {
	defer c.Close()
	var req Command
	decoder := json.NewDecoder(io.LimitReader(c, 512))
	if err := decoder.Decode(&req); err != nil {
		log.Errorf("failed to handle connection: %v", err)
		return
	}

	cmdChan <- req
}
