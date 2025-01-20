package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

const (
	RegisterClientCommand = "register-client"
	ListenTopicCommand    = "listen-topic"
	SendMessageCommand    = "send-message"
)

// Message represents a decoded protocol message
type Message struct {
	Command    string
	Topic      string
	ClientAddr string
	Payload    []byte
}

// decode follows the pattern of [total length] [cmd length] [cmd] ... other custom packages
func DecodeMessageFromConn(conn net.Conn) (*Message, error) {

	// Set a read deadline to avoid infinite blocking
	if err := conn.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		return nil, fmt.Errorf("setting read deadline: %w", err)
	}

	// Read total length (2 bytes)
	totalLenBuf := make([]byte, 2)
	if _, err := io.ReadFull(conn, totalLenBuf); err != nil {
		return nil, fmt.Errorf("reading total length: %w", err)
	}
	totalLen := binary.BigEndian.Uint16(totalLenBuf)

	// Read the complete message based on total length
	data := make([]byte, totalLen)
	if _, err := io.ReadFull(conn, data); err != nil {
		return nil, fmt.Errorf("reading message data: %w", err)
	}

	// Parse the message
	msg := &Message{}
	currentPos := 0

	// Read command
	if len(data) < 2 {
		return nil, fmt.Errorf("message too short for command length")
	}

	log.Println("cmd length bytes -> ", data[currentPos:])
	cmdLen := binary.BigEndian.Uint16(data[currentPos:])
	currentPos += 2

	if len(data) < currentPos+int(cmdLen) {
		return nil, fmt.Errorf("message too short for command")
	}
	log.Println("cmd bytes -> ", data[currentPos:currentPos+int(cmdLen)])
	msg.Command = string(data[currentPos : currentPos+int(cmdLen)])
	currentPos += int(cmdLen)

	// Parse rest based on command type
	switch msg.Command {
	case RegisterClientCommand:
		// Read client address
		if len(data) < currentPos+2 {
			return nil, fmt.Errorf("message too short for address length")
		}
		addrLen := binary.BigEndian.Uint16(data[currentPos:])
		currentPos += 2

		if len(data) < currentPos+int(addrLen) {
			return nil, fmt.Errorf("message too short for address")
		}
		msg.ClientAddr = string(data[currentPos : currentPos+int(addrLen)])

	case ListenTopicCommand, SendMessageCommand:
		// Read topic
		if len(data) < currentPos+2 {
			return nil, fmt.Errorf("message too short for topic length")
		}
		topicLen := binary.BigEndian.Uint16(data[currentPos:])
		currentPos += 2

		if len(data) < currentPos+int(topicLen) {
			return nil, fmt.Errorf("message too short for topic")
		}
		msg.Topic = string(data[currentPos : currentPos+int(topicLen)])
		currentPos += int(topicLen)

		// Read client address
		if len(data) < currentPos+2 {
			return nil, fmt.Errorf("message too short for address length")
		}
		addrLen := binary.BigEndian.Uint16(data[currentPos:])
		currentPos += 2

		if len(data) < currentPos+int(addrLen) {
			return nil, fmt.Errorf("message too short for address")
		}
		msg.ClientAddr = string(data[currentPos : currentPos+int(addrLen)])
		currentPos += int(addrLen)

		// For send-message, read payload
		if msg.Command == SendMessageCommand {
			if len(data) < currentPos+2 {
				return nil, fmt.Errorf("message too short for payload length")
			}
			payloadLen := binary.BigEndian.Uint16(data[currentPos:])
			currentPos += 2

			if len(data) < currentPos+int(payloadLen) {
				return nil, fmt.Errorf("message too short for payload")
			}
			msg.Payload = data[currentPos : currentPos+int(payloadLen)]
		}

	default:
		return nil, fmt.Errorf("unknown command: %s", msg.Command)
	}

	return msg, nil
}

// EncodeRegisterRequest encodes a register-client request
// [total length] [cmd length] [cmd] [client addr length] [client addr]
func EncodeRegisterRequest(clientAddr string) ([]byte, error) {
	cmd := RegisterClientCommand
	totalLen := 2 + 2 + len(cmd) + 2 + len(clientAddr)

	buf := make([]byte, totalLen)
	currentPos := 0

	// Write total length
	binary.BigEndian.PutUint16(buf[currentPos:], uint16(totalLen-2))
	currentPos += 2

	// Write command
	binary.BigEndian.PutUint16(buf[currentPos:], uint16(len(cmd)))
	currentPos += 2
	copy(buf[currentPos:], cmd)
	currentPos += len(cmd)

	// Write client address
	binary.BigEndian.PutUint16(buf[currentPos:], uint16(len(clientAddr)))
	currentPos += 2
	copy(buf[currentPos:], clientAddr)

	return buf, nil
}

// EncodeListenTopicRequest encodes a listen-topic request
// [total length] [cmd length] [cmd] [topic length] [topic] [client addr length] [client addr]
func EncodeListenTopicRequest(topic, clientAddr string) ([]byte, error) {
	cmd := ListenTopicCommand
	totalLen := 2 + 2 + len(cmd) + 2 + len(topic) + 2 + len(clientAddr)

	buf := make([]byte, totalLen)
	currentPos := 0

	// Write total length
	binary.BigEndian.PutUint16(buf[currentPos:], uint16(totalLen-2))
	currentPos += 2

	// Write command
	binary.BigEndian.PutUint16(buf[currentPos:], uint16(len(cmd)))
	currentPos += 2
	copy(buf[currentPos:], cmd)
	currentPos += len(cmd)

	// Write topic
	binary.BigEndian.PutUint16(buf[currentPos:], uint16(len(topic)))
	currentPos += 2
	copy(buf[currentPos:], topic)
	currentPos += len(topic)

	// Write client address
	binary.BigEndian.PutUint16(buf[currentPos:], uint16(len(clientAddr)))
	currentPos += 2
	copy(buf[currentPos:], clientAddr)

	return buf, nil
}

// EncodeMessage encodes a send-message request
// [total length] [cmd length] [cmd] [topic length] [topic] [client addr length] [client addr] [payload length] [payload]
func EncodeMessage(topic, clientAddr string, payload []byte) ([]byte, error) {
	cmd := SendMessageCommand
	totalLen := 2 + 2 + len(cmd) + 2 + len(topic) + 2 + len(clientAddr) + 2 + len(payload)

	buf := make([]byte, totalLen)
	currentPos := 0

	// Write total length
	binary.BigEndian.PutUint16(buf[currentPos:], uint16(totalLen-2))
	currentPos += 2

	// Write command
	binary.BigEndian.PutUint16(buf[currentPos:], uint16(len(cmd)))
	currentPos += 2
	copy(buf[currentPos:], cmd)
	currentPos += len(cmd)

	// Write topic
	binary.BigEndian.PutUint16(buf[currentPos:], uint16(len(topic)))
	currentPos += 2
	copy(buf[currentPos:], topic)
	currentPos += len(topic)

	// Write client address
	binary.BigEndian.PutUint16(buf[currentPos:], uint16(len(clientAddr)))
	currentPos += 2
	copy(buf[currentPos:], clientAddr)
	currentPos += len(clientAddr)

	// Write payload
	binary.BigEndian.PutUint16(buf[currentPos:], uint16(len(payload)))
	currentPos += 2
	copy(buf[currentPos:], payload)

	return buf, nil
}
