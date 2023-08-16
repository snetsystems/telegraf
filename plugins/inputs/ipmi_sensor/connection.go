package ipmi_sensor

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// Connection properties for a Client
type Connection struct {
	Hostname  string
	Username  string
	Password  string
	Port      int
	Interface string
	Privilege string
	HexKey    string
	IpmiIP    string
}

func NewConnection(server, privilege, hexKey string) *Connection {
	conn := &Connection{
		Privilege: privilege,
		HexKey:    hexKey,
	}
	inx1 := strings.LastIndex(server, "@")
	inx2 := strings.Index(server, "(")

	connstr := server

	if inx1 > 0 {
		security := server[0:inx1]
		connstr = server[inx1+1:]
		up := strings.SplitN(security, ":", 2)
		if len(up) == 2 {
			conn.Username = up[0]
			conn.Password = up[1]
		}
	}

	if inx2 > 0 {
		inx2 = strings.Index(connstr, "(")
		inx3 := strings.Index(connstr, ")")

		conn.Interface = connstr[0:inx2]
		conn.IpmiIP = connstr[inx2+1 : inx3]
	}

	return conn
}

func (c *Connection) options() []string {
	intf := c.Interface
	if intf == "" {
		intf = "lan"
	}

	options := []string{
		"-H", c.IpmiIP,
		"-U", c.Username,
		"-P", c.Password,
		"-I", intf,
	}

	if c.HexKey != "" {
		options = append(options, "-y", c.HexKey)
	}
	if c.Port != 0 {
		options = append(options, "-p", strconv.Itoa(c.Port))
	}
	if c.Privilege != "" {
		options = append(options, "-L", c.Privilege)
	}
	return options
}

// RemoteIP returns the remote (bmc) IP address of the Connection
func (c *Connection) RemoteIP() string {
	if net.ParseIP(c.IpmiIP) == nil {
		addrs, err := net.LookupHost(c.IpmiIP)
		if err != nil && len(addrs) > 0 {
			return addrs[0]
		}
	}
	return c.IpmiIP
}

// LocalIP returns the local (client) IP address of the Connection
func (c *Connection) LocalIP() string {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", c.IpmiIP, c.Port))
	if err != nil {
		// don't bother returning an error, since this value will never
		// make it to the bmc if we can't connect to it.
		return c.IpmiIP
	}
	_ = conn.Close()
	host, _, _ := net.SplitHostPort(conn.LocalAddr().String())
	return host
}
