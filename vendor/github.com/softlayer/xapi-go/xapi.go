// Package xapi is a generic binding of the Citrix XenServer API.
// Read more here: http://docs.vmd.citrix.com/XenServer/6.0.0/1.0/en_gb/api/
package xapi

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/kolo/xmlrpc"
)

// Client is the Xen API client
type Client struct {
	Session  string
	URI      string
	Username string
	Password string
	Version  string
	rpc      *xmlrpc.Client
}

// NewClient stand up a new XapiClient. Version should probably be "1.2" unless you know what you are doing.
func NewClient(uri, username, password, version string) (client Client) {
	client.URI = uri
	client.Username = username
	client.Password = password
	client.Version = version
	client.rpc, _ = xmlrpc.NewClient(
		client.URI,
		&http.Transport{
			Dial: TimeoutDialer(),
		})

	return
}

type Struct map[string]interface{}

// Login will authenticate the current session
func (client *Client) Login() (err error) {
	resp := StringResponse{}
	err = client.Call(
		&resp,
		"session.login_with_password",
		client.Username,
		client.Password,
		client.Version)
	if err == nil {
		client.Session = resp.Value
	}
	return
}

func (client *Client) GetSession() (session Session, err error) {
	var resp SessionResponse
	err = client.SessionCall(&resp, "session.get_record",
		client.Session)
	return resp.Value, err
}

func (client *Client) RegisterEvent(classes ...interface{}) error {
	return client.SessionCall(nil, "event.register", classes)
}

func (client *Client) UnregisterAllEvents() error {
	return client.SessionCall(nil, "event.unregister", []string{"*"})
}

func (client *Client) NextEvent() (ev []Event, err error) {
	var resp EventsResponse
	err = client.SessionCall(&resp, "event.next")
	return resp.Value, err
}

func (client *Client) GetVMs() (machines []string, err error) {
	var resp StringsResponse
	err = client.SessionCall(&resp, "VM.get_all")
	return resp.Value, err
}

func (client *Client) GetVM(opref string) (vm VM, err error) {
	var resp VMResponse
	err = client.SessionCall(&resp, "VM.get_record", opref)
	return resp.Value, err
}

func (client *Client) GetVIF(opref string) (vif VIF, err error) {
	var resp VIFResponse
	err = client.SessionCall(&resp, "VIF.get_record", opref)
	return resp.Value, err
}

func (client *Client) GetPIF(opref string) (pif PIF, err error) {
	var resp PIFResponse
	err = client.SessionCall(&resp, "PIF.get_record", opref)
	return resp.Value, err
}

func (client *Client) GetHost(opref string) (host Host, err error) {
	var resp HostResponse
	err = client.SessionCall(&resp, "host.get_record", opref)
	return resp.Value, err
}

// GetHostname gets the hostname of a Host.  Useful in combination with GetSession and session.This_host
func (client *Client) GetHostname(opref string) (hostname string, err error) {
	var resp StringResponse
	err = client.SessionCall(&resp, "host.get_hostname", opref)
	return resp.Value, err
}

// SessionCall is a useful for making multiple calls that require the session ID.  Automatically prepends the existing
// session OpaqRef to the beginning of the API call.  You can see the session ID by looking at
// Client.Session.
func (client *Client) SessionCall(result Response, call string, params ...interface{}) (err error) {
	if client.Session == "" {
		return fmt.Errorf("no session")
	}

	params = append([]interface{}{client.Session}, params...)
	return client.Call(result, call, params...)
}

// TimeoutDialer is a custom Dialer for HTTP so that the initial connection only lasts for 1 minute
// and that the lifetime of the connection is only 1 minute as well. See http://golang.org/pkg/net/#Conn
// You shouldn't need to use this directly.
func TimeoutDialer() func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, time.Minute)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(time.Minute))
		return conn, nil
	}
}

// Call makes a generic RPC call passing in a pointer to a struct (or f). The call parameter
// is a combination of class.message.  For example: VIF.get_record, host.evacuate, pool.eject.
// Any time the XAPI specifies a `type ref` it's really an OpaqueReference, which is a UUID, and
// as far as xmlrpc and like library are concerned, a string.
//		x := xapi.NewClient("http://localhost/", "username", "password", "1.2")
//		host := xapi.HostResult{}
//		err := x.Call(&host, "host.get_record", "324c2264-d86f-4a42-a971-bb5fd6203877")
//		if err != nil {
//			fmt.Printf("%v", host)
//		}
func (client *Client) Call(result Response, call string, params ...interface{}) error {

	err := client.rpc.Call(call, params, result)
	if err != nil {
		return err
	}

	// TODO: The XMLRPC library isn't smart enough to handle composition so the Status and ErrorDescription values are never set
	// if result.Status() != "Success" {
	// 	return fmt.Errorf("XenServer Failed: %s", result.ErrorDescription())
	// }

	return nil
}
