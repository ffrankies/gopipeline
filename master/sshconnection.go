package master

import (
	"io/ioutil"
	"net"
	"os/user"
	"strconv"
	"strings"

	"golang.org/x/crypto/ssh"
)

// The SSHConnection contains information about an SSH connection, and methods for running commands over the SSH connection
// @see: https://github.com/jilieryuyi/ssh-simple-client/blob/master/main.go
type SSHConnection struct {
	Address        string       // The IP address of the server
	User           string       // The username on the server
	PrivateKeyPath string       // Path to private key
	Port           int          // The port on which to connect to the server
	session        *ssh.Session // The session for running commands
	client         *ssh.Client  // The client for creating sessions
}

// NewSSHConnection creates a new SSHConnection object
func NewSSHConnection(address string, remoteUser string, port int) *SSHConnection {
	sshConnection := new(SSHConnection)
	sshConnection.Address = address
	sshConnection.User = remoteUser
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	sshConnection.PrivateKeyPath = usr.HomeDir + "/.ssh/id_rsa"
	sshConnection.Port = 22
	sshConnection.client = createClient(sshConnection.User, sshConnection.Address, sshConnection.Port)
	return sshConnection
}

// RunCommand a single command through the SSH Connection
func (conn *SSHConnection) RunCommand(command string) (output string, err error) {
	session, err := conn.client.NewSession()
	if err != nil {
		return
	}
	defer session.Close()
	outputBytes, err := session.Output(command)
	if err != nil {
		return
	}
	output = string(outputBytes)
	return
}

// Close closes connection
func (conn *SSHConnection) Close() error {
	err := conn.client.Close()
	return err
}

// Creates a client connection to the given address with the given user
func createClient(remoteUser string, address string, port int) *ssh.Client {
	publicKeyConfig := getPublicKeyConfig(remoteUser)
	connectionType := "tcp"
	addressAndPort := address + ":" + strconv.Itoa(port)
	if strings.Count(address, ":") > 0 {
		connectionType = "tcp6"
		addressAndPort = "[" + address + "]:" + strconv.Itoa(port)
	}
	sshClient, err := ssh.Dial(connectionType, addressAndPort, publicKeyConfig)
	if err != nil {
		panic(err)
	}
	return sshClient
}

// Creates an ssh config based on the private key of the current user
// @see: https://golang-basic.blogspot.com/2014/06/step-by-step-guide-to-ssh-using-go.html
// Big Question: Why is it using the private key???? I thought it needed the public key??
func getPublicKeyConfig(remoteUser string) (publicKeyConfig *ssh.ClientConfig) {
	// Get private key signer
	privateKeySigner, err := getPrivateKeySigner()
	if err != nil {
		panic(err)
	}
	// Create the ssh config
	publicKeyConfig = &ssh.ClientConfig{
		User: remoteUser,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(privateKeySigner),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	return
}

// Retrieves the public key signer from the current user's home directory
// @see: https://golang-basic.blogspot.com/2014/06/step-by-step-guide-to-ssh-using-go.html
func getPrivateKeySigner() (privateKeySigner ssh.Signer, err error) {
	usr, _ := user.Current()
	privateKeyFile := usr.HomeDir + "/.ssh/id_rsa"
	privateKeyBuffer, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return
	}
	privateKeySigner, err = ssh.ParsePrivateKey(privateKeyBuffer)
	if err != nil {
		return
	}
	return
}
