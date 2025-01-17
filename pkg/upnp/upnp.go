package upnp

// UPNP code taken from Taipei Torrent license is below:
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
// Just enough UPnP to be able to forward ports
import (
	"bytes"
	"encoding/xml"
	"errors"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// NAT is an interface representing a NAT traversal options for example UPNP or NAT-PMP. It provides methods to query
// and manipulate this traversal to allow access to services.
type NAT interface {
	// GetExternalAddress - Get the external address from outside the NAT.
	GetExternalAddress() (addr net.IP, e error)
	// AddPortMapping - Add a port mapping for protocol (
	//  "udp" or "tcp") from external port to internal port with description lasting
	// for timeout.
	AddPortMapping(
		protocol string, externalPort, internalPort int,
		description string, timeout int,
	) (mappedExternalPort int, e error)
	// DeletePortMapping - Remove a previously added port mapping from external
	// port to internal port.
	DeletePortMapping(
		protocol string, externalPort,
		internalPort int,
	) (e error)
}
type upnpNAT struct {
	serviceURL string
	ourIP      string
}

// Discover searches the local network for a UPnP router returning a NAT for the network if so, nil if not.
func Discover() (nat NAT, e error) {
	var ssdp *net.UDPAddr
	ssdp, e = net.ResolveUDPAddr("udp4", "239.255.255.250:1900")
	if e != nil {
		E.Ln(e)
		return
	}
	var conn net.PacketConn
	conn, e = net.ListenPacket("udp4", ":0")
	if e != nil {
		E.Ln(e)
		return
	}
	socket := conn.(*net.UDPConn)
	defer func() {
		if e = socket.Close(); E.Chk(e) {
		}
	}()
	e = socket.SetDeadline(time.Now().Add(3 * time.Second))
	if E.Chk(e) {
		return
	}
	st := "ST: urn:schemas-upnp-org:device:InternetGatewayDevice:1\r\n"
	buf := bytes.NewBufferString(
		"M-SEARCH * HTTP/1.1\r\n" +
			"HOST: 239.255.255.250:1900\r\n" +
			st +
			"MAN: \"ssdp:discover\"\r\n" +
			"MX: 2\r\n\r\n",
	)
	message := buf.Bytes()
	answerBytes := make([]byte, 1024)
	for i := 0; i < 3; i++ {
		_, e = socket.WriteToUDP(message, ssdp)
		if e != nil {
			E.Ln(e)
			return
		}
		var n int
		n, _, e = socket.ReadFromUDP(answerBytes)
		if e != nil {
			E.Ln(e)
			continue
			// socket.Close()
			// return
		}
		answer := string(answerBytes[0:n])
		if !strings.Contains(answer, "\r\n"+st) {
			continue
		}
		// HTTP header field names are case-insensitive. http://www.w3.org/Protocols/rfc2616/rfc2616-sec4.html#sec4.2
		locString := "\r\nlocation: "
		locIndex := strings.Index(strings.ToLower(answer), locString)
		if locIndex < 0 {
			continue
		}
		loc := answer[locIndex+len(locString):]
		endIndex := strings.Index(loc, "\r\n")
		if endIndex < 0 {
			continue
		}
		locURL := loc[0:endIndex]
		var serviceURL string
		serviceURL, e = getServiceURL(locURL)
		if e != nil {
			E.Ln(e)
			return
		}
		var ourIP string
		ourIP, e = getOurIP()
		if e != nil {
			E.Ln(e)
			return
		}
		nat = &upnpNAT{serviceURL: serviceURL, ourIP: ourIP}
		return
	}
	e = errors.New("UPnP port routeable failed")
	return
}

// service represents the Service type in an UPnP xml description. Only the parts we care about are present and thus the
// xml may have more fields than present in the structure.
type service struct {
	ServiceType string `xml:"serviceType"`
	ControlURL  string `xml:"controlURL"`
}

// deviceList represents the deviceList type in an UPnP xml description. Only the parts we care about are present and
// thus the xml may have more fields than present in the structure.
type deviceList struct {
	XMLName xml.Name `xml:"deviceList"`
	Device  []device `xml:"device"`
}

// serviceList represents the serviceList type in an UPnP xml description. Only the parts we care about are present and
// thus the xml may have more fields than present in the structure.
type serviceList struct {
	XMLName xml.Name  `xml:"serviceList"`
	Service []service `xml:"service"`
}

// device represents the device type in an UPnP xml description. Only the parts we care about are present and thus the
// xml may have more fields than present in the structure.
type device struct {
	XMLName     xml.Name    `xml:"device"`
	DeviceType  string      `xml:"deviceType"`
	DeviceList  deviceList  `xml:"deviceList"`
	ServiceList serviceList `xml:"serviceList"`
}

// specVersion represents the specVersion in a UPnP xml description. Only the parts we care about are present and thus
// the xml may have more fields than present in the structure.
type specVersion struct {
	XMLName xml.Name `xml:"specVersion"`
	Major   int      `xml:"major"`
	Minor   int      `xml:"minor"`
}

// root represents the Root document for a UPnP xml description. Only the parts we care about are present and thus the
// xml may have more fields than present in the structure.
type root struct {
	XMLName     xml.Name `xml:"root"`
	SpecVersion specVersion
	Device      device
}

// getChildDevice searches the children of device for a device with the given type.
func getChildDevice(d *device, deviceType string) *device {
	for i := range d.DeviceList.Device {
		if d.DeviceList.Device[i].DeviceType == deviceType {
			return &d.DeviceList.Device[i]
		}
	}
	return nil
}

// getChildDevice searches the service list of device for a service with the given type.
func getChildService(d *device, serviceType string) *service {
	for i := range d.ServiceList.Service {
		if d.ServiceList.Service[i].ServiceType == serviceType {
			return &d.ServiceList.Service[i]
		}
	}
	return nil
}

// getOurIP returns a best guess at what the local IP is.
func getOurIP() (ip string, e error) {
	var hostname string
	hostname, e = os.Hostname()
	if e != nil {
		E.Ln(e)
		return
	}
	return net.LookupCNAME(hostname)
}

// getServiceURL parses the xml description at the given root url to find the url for the WANIPConnection service to be
// used for port forwarding.
func getServiceURL(rootURL string) (url string, e error) {
	var re *http.Response
	re, e = http.Get(rootURL)
	if e != nil {
		E.Ln(e)
		return
	}
	defer func() {
		if e = re.Body.Close(); E.Chk(e) {
		}
	}()
	if re.StatusCode >= 400 {
		e = errors.New(string(rune(re.StatusCode)))
		return
	}
	var r root
	e = xml.NewDecoder(re.Body).Decode(&r)
	if e != nil {
		E.Ln(e)
		return
	}
	a := &r.Device
	if a.DeviceType != "urn:schemas-upnp-org:device:InternetGatewayDevice:1" {
		e = errors.New("no InternetGatewayDevice")
		return
	}
	b := getChildDevice(a, "urn:schemas-upnp-org:device:WANDevice:1")
	if b == nil {
		e = errors.New("no WANDevice")
		return
	}
	c := getChildDevice(b, "urn:schemas-upnp-org:device:WANConnectionDevice:1")
	if c == nil {
		e = errors.New("no WANConnectionDevice")
		return
	}
	d := getChildService(c, "urn:schemas-upnp-org:service:WANIPConnection:1")
	if d == nil {
		e = errors.New("no WANIPConnection")
		return
	}
	url = combineURL(rootURL, d.ControlURL)
	return
}

// combineURL appends subURL onto rootURL.
func combineURL(rootURL, subURL string) string {
	protocolEnd := "://"
	protoEndIndex := strings.Index(rootURL, protocolEnd)
	a := rootURL[protoEndIndex+len(protocolEnd):]
	rootIndex := strings.Index(a, "/")
	return rootURL[0:protoEndIndex+len(protocolEnd)+rootIndex] + subURL
}

// soapBody represents the <s:Body> element in a SOAP reply. fields we don't care about are elided.
type soapBody struct {
	XMLName xml.Name `xml:"Body"`
	Data    []byte   `xml:",innerxml"`
}

// soapEnvelope represents the <s:Envelope> element in a SOAP reply. fields we don't care about are elided.
type soapEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    soapBody `xml:"Body"`
}

// soapRequests performs a soap request with the given parameters and returns the xml replied stripped of the soap
// headers. in the case that the request is unsuccessful the an error is returned.
func soapRequest(url, function, message string) (replyXML []byte, e error) {
	fullMessage := "<?xml version=\"1.0\" ?>" +
		"<s:Envelope xmlns:s=\"http://schemas.xmlsoap." +
		"org/soap/envelope/\" s:encodingStyle=\"http://schemas.xmlsoap." +
		"org/soap/encoding/\">\r\n" +
		"<s:Body>" + message + "</s:Body></s:Envelope>"
	var req *http.Request
	req, e = http.NewRequest("POST", url, strings.NewReader(fullMessage))
	if e != nil {
		E.Ln(e)
		return nil, e
	}
	req.Header.Set("Content-Type", "text/xml ; charset=\"utf-8\"")
	req.Header.Set("User-Agent", "Darwin/10.0.0, UPnP/1.0, MiniUPnPc/1.3")
	// req.Header.Set("Transfer-Encoding", "chunked")
	req.Header.Set(
		"SOAPAction",
		"\"urn:schemas-upnp-org:service:WANIPConnection:1#"+function+"\"",
	)
	req.Header.Set("Connection", "Close")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Pragma", "no-cache")
	var r *http.Response
	r, e = http.DefaultClient.Do(req)
	if e != nil {
		E.Ln(e)
		return nil, e
	}
	if r.Body != nil {
		defer func() {
			if e = r.Body.Close(); E.Chk(e) {
			}
		}()
	}
	if r.StatusCode >= 400 {
		// L.Stderr(function, r.StatusCode)
		e = errors.New("Error " + strconv.Itoa(r.StatusCode) + " for " + function)
		r = nil
		return
	}
	var reply soapEnvelope
	e = xml.NewDecoder(r.Body).Decode(&reply)
	if e != nil {
		E.Ln(e)
		return nil, e
	}
	return reply.Body.Data, nil
}

// getExternalIPAddressResponse represents the XML response to a GetExternalIPAddress SOAP request.
type getExternalIPAddressResponse struct {
	XMLName           xml.Name `xml:"GetExternalIPAddressResponse"`
	ExternalIPAddress string   `xml:"NewExternalIPAddress"`
}

// GetExternalAddress implements the NAT interface by fetching the external IP from the UPnP router.
func (n *upnpNAT) GetExternalAddress() (addr net.IP, e error) {
	message := "<u:GetExternalIPAddress xmlns:u=\"urn:schemas-upnp-org:service:WANIPConnection:1\"/>\r\n"
	var response []byte
	response, e = soapRequest(n.serviceURL, "GetExternalIPAddress", message)
	if e != nil {
		E.Ln(e)
		return nil, e
	}
	var reply getExternalIPAddressResponse
	e = xml.Unmarshal(response, &reply)
	if e != nil {
		E.Ln(e)
		return nil, e
	}
	addr = net.ParseIP(reply.ExternalIPAddress)
	if addr == nil {
		return nil, errors.New("unable to parse ip address")
	}
	return addr, nil
}

// AddPortMapping implements the NAT interface by setting up a port forwarding from the UPnP router to the local machine
// with the given ports and protocol.
func (n *upnpNAT) AddPortMapping(
	protocol string,
	externalPort, internalPort int,
	description string,
	timeout int,
) (mappedExternalPort int, e error) {
	// A single concatenation would break ARM compilation.
	message := "<u:AddPortMapping xmlns:u=\"urn:schemas-upnp-org:service:WANIPConnection:1\">\r\n" +
		"<NewRemoteHost></NewRemoteHost><NewExternalPort>" + strconv.Itoa(externalPort)
	message += "</NewExternalPort><NewProtocol>" + strings.ToUpper(protocol) + "</NewProtocol>"
	message += "<NewInternalPort>" + strconv.Itoa(internalPort) + "</NewInternalPort>" +
		"<NewInternalClient>" + n.ourIP + "</NewInternalClient>" +
		"<NewEnabled>1</NewEnabled><NewPortMappingDescription>"
	message += description +
		"</NewPortMappingDescription><NewLeaseDuration>" + strconv.Itoa(timeout) +
		"</NewLeaseDuration></u:AddPortMapping>"
	var response []byte
	response, e = soapRequest(n.serviceURL, "AddPortMapping", message)
	if e != nil {
		E.Ln(e)
		return
	}
	// TODO: check response to see if the port was forwarded
	//
	// If the port was not wildcard we don't get an reply with the port in it. Not sure about wildcard yet. miniupnpc
	// just checks for error codes here.
	mappedExternalPort = externalPort
	_ = response
	return
}

// DeletePortMapping implements the NAT interface by removing up a port forwarding from the UPnP router to the local
// machine with the given ports and.
func (n *upnpNAT) DeletePortMapping(
	protocol string,
	externalPort, internalPort int,
) (e error) {
	message := "<u:DeletePortMapping xmlns:u=\"urn:schemas-upnp-org:service:WANIPConnection:1\">\r\n" +
		"<NewRemoteHost></NewRemoteHost><NewExternalPort>" + strconv.Itoa(externalPort) +
		"</NewExternalPort><NewProtocol>" + strings.ToUpper(protocol) + "</NewProtocol>" +
		"</u:DeletePortMapping>"
	var response []byte
	response, e = soapRequest(n.serviceURL, "DeletePortMapping", message)
	if e != nil {
		E.Ln(e)
		return
	}
	// TODO: check response to see if the port was deleted
	_ = response
	return
}
