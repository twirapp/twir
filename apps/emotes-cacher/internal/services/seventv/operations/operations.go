package operations

// Op 	Name 	Type 	Description
// 0 	Dispatch 	⬇️ 	A standard event message, sent when a subscribed event is emitted
// 1 	Hello 	⬇️ 	Received upon connecting, presents info about the session
// 2 	Heartbeat 	⬇️ 	Ensures the connection is still alive
// 4 	Reconnect 	⬇️ 	Server wants the client to reconnect
// 5 	Ack 	⬇️ 	Server acknowledges an action by the client
// 6 	Error 	⬇️ 	An error occured, you should log this
// 7 	End of Stream 	⬇️ 	The server will send no further data and imminently end the connection
// 33 	Identify 	⬆️ 	Authenticate with an account
// 34 	Resume 	⬆️ 	Try to resume a previous session
// 35 	Subscribe 	⬆️ 	Watch for changes on specific objects or sources. Don't smash it!
// 36 	Unsubscribe 	⬆️ 	Stop listening for changes
// 37 	Signal 	⬆️

type IncomingOp int

const (
	// IncomingOpDispatch ⬇️ Operation codes as defined by the 7TV WebSocket API
	IncomingOpDispatch IncomingOp = 0
	// IncomingOpHello ⬇️ Received upon connecting, presents info about the session
	IncomingOpHello IncomingOp = 1
	// IncomingOpHeartbeat ⬇️ Ensures the connection is still alive
	IncomingOpHeartbeat IncomingOp = 2
	// IncomingOpReconnect ⬇️ Server wants the client to reconnect
	IncomingOpReconnect IncomingOp = 4
	// IncomingOpAck ⬇️ Server acknowledges an action by the client
	IncomingOpAck IncomingOp = 5
	// IncomingOpError ⬇️ An error occurred, you should log this
	IncomingOpError IncomingOp = 6
	// IncomingOpEndOfStream ⬇️ The server will send no further data and imminently end the connection
	IncomingOpEndOfStream IncomingOp = 7
)

type OutgoingOp int

const (
	// OutgoingOpIdentify ⬆️ Authenticate with an account
	OutgoingOpIdentify OutgoingOp = 33
	// OutgoingOpResume ⬆️ Try to resume a previous session
	OutgoingOpResume OutgoingOp = 34
	// OutgoingOpSubscribe ⬆️ Watch for changes on specific objects or sources. Don't smash it!
	OutgoingOpSubscribe OutgoingOp = 35
	// OutgoingOpUnsubscribe ⬆️ Stop listening for changes
	OutgoingOpUnsubscribe OutgoingOp = 36
	// OutgoingOpSignal ⬆️ Send a signal to the server
	OutgoingOpSignal OutgoingOp = 37
)
