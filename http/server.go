package http

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Server struct {
	port int
}
type Request struct {
	method string
	uri string
	http string
	headers map[string]string
	body []byte
}

type Response struct {
	headers []string
	body []byte
}

func CreateServer(port int) Server {
	server := Server{
		port,
	}

	return server
}

func (server *Server) readRequest(c net.Conn) Request {
	fmt.Println("readRequest")

	request := Request{
		"",
		"",
		"",
		map[string]string{},
		[]byte{},
	}
	parsedStatus := false
	reader := bufio.NewReader(c)
	for {
			netData, err := reader.ReadString('\n')
			if err != nil {
					fmt.Println(err)
					return request
			}
			
			if !parsedStatus {
				x := strings.Split(strings.Trim(netData, "\r\n"), " ")
				request.method = x[0]
				request.uri = x[1]
				request.http = x[2]
				parsedStatus = true
				continue
			}

			if netData == "\r\n" {
				server.readBody(c, reader, &request)
				break
			}

			header := strings.SplitN(netData, ":", 2)
			
			request.headers[strings.Trim(header[0], " ")] = strings.Trim(header[1], " ")
	}

	return request
}

func (server *Server) readBody(c net.Conn, reader *bufio.Reader, request *Request) Request {
	fmt.Println("readBody")

	if length, ok := request.headers["Content-Length"]; ok {
		length, err := strconv.ParseInt(strings.Trim(length, "\r\n"), 10, 0);
		if (err != nil) {
			return *request
		}

		request.body = make([]byte,length) 
		reader.Read(request.body)
	}

	return *request
}

func (server *Server) sendResponse(c net.Conn, r Response) {
	for _, header := range r.headers {
		c.Write([]byte(header+"\n"))
	}
	c.Write([]byte("\r\n"))

	c.Write(r.body)
}

func (server *Server) handleConnection(c net.Conn) {
	fmt.Println("handleConnection")
	request := server.readRequest(c)
	fmt.Println(request)

	response := Response{
		[]string{"HTTP/1.1 200 OK"},
		request.body,
	}

	server.sendResponse(c, response)
	c.Close()
}


func (server *Server) Listen() {
	l, err := net.Listen("tcp4", ":"+strconv.FormatInt(int64(server.port), 10))
	if err != nil {
			fmt.Println(err)
			return
	}
	defer l.Close()

	for {
			c, err := l.Accept()
			if err != nil {
					fmt.Println(err)
					return
			}
			go server.handleConnection(c)
	}
}
