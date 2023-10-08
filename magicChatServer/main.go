package main

import (
	"fmt"
	"net"
)

func main() {

	fmt.Println("Starting Server")

	ln, _ := net.listen("tcp", ":4875")
	defer ln.Close()
	c := make(chan string)

	var arra []string
	var i int = 0
	fl := "true"

	for {
		ln.Accept()

		go connections(ln, c chan string)

		outputs := <-c

		uno := outputs[0]
		dos := outputs[1]

		if uno[0] == dos[0] {
			if uno[1] == dos[1] {
				conn, err := net.Dial("tcp", uno[2])
				if err != nil {
					panic(err)
				}
				var messageSend = dos[2]
				_, err = conn.Write([]byte(string(messageSend)))
				if err != nil {
					panic(err)
				}
				conn, err := net.Dial("tcp", dos[2])
				if err != nil {
					panic(err)
				}
				var messageSend = uno[2]
				_, err = conn.Write([]byte(string(messageSend)))
				if err != nil {
					panic(err)
				}

			}
		} else {
			conn, err := net.Dial("tcp", uno[2])
			if err != nil {
				panic(err)
			}
			var messageSend = "Error-NoMatch"
			_, err = conn.Write([]byte(string(messageSend)))
			if err != nil {
				panic(err)
			}
			conn, err := net.Dial("tcp", dos[2])
			if err != nil {
				panic(err)
			}
			var messageSend = "Error-NoMatch"
			_, err = conn.Write([]byte(string(messageSend)))
			if err != nil {
				panic(err)
			}
		}
	}

	//recieve channel and salted hash, after you get 2

	//compare channel and salted hash
	//return flag based on output

}

//compare channel and hash from threads
//if both match, return success flag,

func connections(conn net.Conn, c chan string) {
	// recieve string
	buf := make([]byte, 1024)
	var recvString string
	//recvString, _ = conn.Read(string([]byte(buf))) // Change the size as needed
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading from connection:", err)
	}
	recvString = string(buf)

	//parse channel, saltedhash, IP:PORT
	//info := strings.SplitAfter(recvString, ",")
	//channelw := info[0]
	//hash := info[1]
	//ipPort := info[2]

	//send channelw and hash
	c <- recvString

	//if hashes and channels match send ip:port

}
