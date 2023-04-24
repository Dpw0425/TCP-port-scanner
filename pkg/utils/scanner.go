package utils

import (
	"fmt"
	"net"
	"time"
)

func Scanner(host string, port int, timeout time.Duration) bool {
	time.Sleep(time.Millisecond * 1)
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if err != nil {
		return false
	}

	conn.Close()
	return true
}
