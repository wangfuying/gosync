package util

import (
	"io"
	"net"
)

func ReadData(conn net.Conn) ([]byte, error) {
	data := make([]byte, 0) //此处做一个输入缓冲以免数据过长读取到不完整的数据
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil && err != io.EOF {
			Error(err.Error())
			return nil, err
		}
		data = append(data, buf[:n]...)
		if n != 2048 {
			break
		}
	}
	return data, nil
}
