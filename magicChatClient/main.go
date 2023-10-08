package main

import (
	"bufio"
	"crypto/aes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	//"encoding/hex"
)

const (
	MAIL_SERVER      = "www.test.com" // replace with ip or web addr
	MAILSRV_PORT     = "4875"
	MAIL_SERVER_TYPE = "tcp"
)

func main() {

	fmt.Println("Welcome to Magic Chat")
	var rawInputKey string
	correctSize := false
	var curSize int

	for !correctSize {

		fmt.Println("Enter your 5 Word Passphrase")
		fmt.Println("Format: one-two-three-four-five")
		fmt.Scanln(&rawInputKey) //get keyPhrase input

		//count number of "-" characters
		curSize = strings.Count(rawInputKey, "-")

		switch {
		case curSize < 4:
			continue
		case curSize == 4:
			correctSize = true
			continue
		case curSize > 4:
			continue
		}
	}
	//end of recieving and validating key and channel

	channel, passkeySalted, passKey := setupPass(rawInputKey) // get seperate passkey and channelkey

	//generate encryption key

	connectServer(channel, passKey, passkeySalted)
	//send key to server with set channel designator
	//get response
	//send open port and ip
	//recieve other client open port and ip
	//try to direct connect to user using that info

	//if cannot connect directly, relay through central server

}

func setupPass(rawKey string) (string, string, string) {

	keyArr := strings.SplitAfter(rawKey, "-")

	channel1 := keyArr[0]
	passkeySalted1 := keyArr[1] + keyArr[2] + keyArr[3] + keyArr[4]
	passKey1 := keyArr[1] + keyArr[2] + keyArr[3]

	//seperate channel and passphrase into two seperate values and return them seperately
	return channel1, passkeySalted1, passKey1
}

func getHash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
func connectServer(channel string, passPhrase string, saltedKey string) {
	//passPhrase needs to be hashed

	//connect to server:port on tcp (mailserver)
	conn, err := net.Dial(MAIL_SERVER_TYPE, MAIL_SERVER+":"+MAILSRV_PORT)
	if err != nil {
		panic(err)
	}

	tmPort := "14574" // to do in future, negotiate open port (not pre-programmed)
	tmpIP := getIp()
	portIP := tmpIP + ":" + tmPort
	hashIs := getHash(saltedKey)

	_, err = conn.Write([]byte(string(channel + "," + hashIs + "," + portIP)))
	if err != nil {
		panic(err)
	}
	//sends salted key and channel

	buffer := make([]byte, 1024)
	mLen, err := conn.Read(buffer)

	if err != nil {
		fmt.Println("Error: " + err.Error())
	}

	//regexp.Match( , string(buffer[:mLen]))

	//if server sends back ip and port
	if string(buffer[:mLen]) == "Error-NoMatch" {
		fmt.Println("Error - No Match Found")
		os.Exit(0)
		//if server sends back error message
	} else if matches, _ := regexp.MatchString("[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{3}:[0-9]{1,5}", string(buffer[:mLen])); matches {
		// if matches regex ip:port connect clients and chat
		fmt.Println("Connecting to Other Client")
		var encrKey string = keyGen(passPhrase)
		go chatDirect(string(buffer[:mLen]), encrKey)

		//start new thread reaching out to this port and ip

		//listen in channel for chatDirect to change flag, if changed break, else keep listening

	}

}

func chatDirect(connectTo string, encrkey string) {

	//listen
	go listen(connectTo, encrkey)

	//connect
	go send(connectTo, encrkey)

	//chat

}

func listen(listenTo string, encrKey1 string) {
	//parse port and use
	prtArr := strings.SplitAfter(listenTo, ":")
	ln, _ := net.Listen("tcp", prtArr[1])
	defer ln.Close()
	conn, _ := ln.Accept()
	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')

		decoded := decodeMessage(string(message), encrKey1)
		fmt.Println("RCVD> " + decoded)
	}
}

func send(conTo string, encrKey string) {
	for {
		conn, err := net.Dial("tcp", conTo)
		if err != nil {
			panic(err)
		}
		var yourMessage string
		fmt.Scanln(&yourMessage)
		var messageSend string = encodeMessage(yourMessage, encrKey)
		_, err = conn.Write([]byte(string(messageSend)))
		if err != nil {
			panic(err)
		}
	}
}

func getIp() string {
	var localAddress string

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress = conn.LocalAddr().(*net.UDPAddr).String()

	return localAddress
}

func keyGen(passPhrase string) string {
	var encrKey string

	shaer := sha256.New()
	shaer.Write([]byte(passPhrase))
	encrKey = hex.EncodeToString(shaer.Sum(nil))

	//keygen using hash of passPhrase
	return encrKey
}

func encodeMessage(theirMessage string, encryptionKey string) string {
	//encoded := theirMessage //placeholder
	aes, err := aes.NewCipher([]byte(encryptionKey))
	if err != nil {
		panic(err)
	}
	ciphertext := make([]byte, len(theirMessage))
	aes.Encrypt(ciphertext, []byte(theirMessage))
	return string(ciphertext)
}

func decodeMessage(message string, key1 string) string {
	//decoded := message //placeholder

	aes, err := aes.NewCipher([]byte(key1))
	if err != nil {
		panic(err)
	}
	pt := make([]byte, len(message))
	aes.Decrypt(pt, []byte(message))
	stringy := string(pt[:])
	return stringy
}
