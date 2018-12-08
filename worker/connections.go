package worker

import (
	"encoding/gob"
	"net"
	"sync"
)

// Connections is a list of Connection objects
type Connections struct {
	Cons    []*Connection // The list of Connection objects
	mutex   *sync.Mutex   // For concurrency stuff
	counter int           // The roundRobin counter
}

// NewConnections creates a new empty list of connections
func NewConnections() *Connections {
	connections := new(Connections)
	connections.Cons = make([]*Connection, 0)
	connections.mutex = &sync.Mutex{}
	connections.counter = 0
	return connections
}

// AddConnection adds a new connection to the list of connections
func (connections *Connections) AddConnection(address string) {
	connection := NewConnection(address)
	connections.mutex.Lock()
	connections.Cons = append(connections.Cons, connection)
	connections.mutex.Unlock()
}

// Select uses round robin returns the next Connection's encoder along which to send the data
func (connections *Connections) Select() *gob.Encoder {
	for len(connections.Cons) == 0 {
		// Busy wait lol
	}
	connections.mutex.Lock()
	connections.counter++
	connections.counter %= len(connections.Cons)
	encoder := connections.Cons[connections.counter].Encoder
	connections.mutex.Unlock()
	return encoder
}

// Connection maintains a connection to the next node
type Connection struct {
	Address string       // The address of the next node
	Con     net.Conn     // The connection to the next node
	Encoder *gob.Encoder // The encoder for sending data along the connections
}

// NewConnection creates a new connection object
func NewConnection(address string) *Connection {
	connection := new(Connection)
	connection.Address = address
	con, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	connection.Con = con
	connection.Encoder = gob.NewEncoder(connection.Con)
	return connection
}

// Close closes the connection
func (connection *Connection) Close() {
	err := connection.Con.Close()
	if err != nil {
		panic(err)
	}
}
