package common

import "net"

// GetOutboundIPAddressHack creates an outgoing connection, and finds the outgoing net host address from that
// connection. This is done because the listener's address is always localhost (127.0.0.1)
func GetOutboundIPAddressHack() string {
	connection, err := net.Dial("udp", "8.8.8.8:80") // Connect to Google services
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	localNetAddress := connection.LocalAddr().String()
	host, _, err := net.SplitHostPort(localNetAddress)
	if err != nil {
		panic(err)
	}
	return host
}

// getPortNumberFromListener parses the listener's address to obtain the port number it's running on as a string.
// This is necessary because we're using dynamic port allocation
func GetPortNumberFromListener(listener net.Listener) string {
	listenerAddress := listener.Addr().String()
	_, port, err := net.SplitHostPort(listenerAddress)
	if err != nil {
		panic(err)
	}
	return port
}
