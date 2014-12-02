// Copyright 2013-3014 Adam Presley. All rights reserved
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package websocket

import (
	ws "github.com/gorilla/websocket"
	"github.com/mailslurper/libmailslurper/model/mailitem"
)

// Structure for tracking and working with websockets
type WebSocketConnection struct {
	// Websocket connection handle
	WS *ws.Conn

	// Buffered channel for outbound messages
	SendChannel chan mailitem.MailItem
}

func ActivateSocket(connection *WebSocketConnection) {
	WebSocketConnections[connection] = true
}

/*
This function takes a MailItem and sends it to all open websockets.
*/
func BroadcastMessageToWebSockets(message mailitem.MailItem) {
	for connection := range WebSocketConnections {
		connection.SendChannel <- message
	}
}

func DestroyConnection(connection *WebSocketConnection) {
	// Remove the connection from our map, and close its channel
	delete(WebSocketConnections, connection)
	close(connection.SendChannel)
}
