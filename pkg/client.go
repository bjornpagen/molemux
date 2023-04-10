package molemux

import (
	"net"
	"sync"

	"github.com/bjornpagen/ezdb"
	swgp "github.com/database64128/swgp-go/service"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun"
)

type Option func(option *options) error

type options struct {
	targetConnections *int
}

func WithTargetConnections(n int) Option {
	return func(option *options) error {
		option.targetConnections = &n
		return nil
	}
}

type Client struct {
	dbPath  string
	options *options

	clientState
}

func NewClient(dbPath string, opts ...Option) (*Client, error) {
	options := &options{}
	for _, opt := range opts {
		if err := opt(options); err != nil {
			return nil, err
		}
	}

	if options.targetConnections == nil {
		n := 5
		options.targetConnections = &n
	}

	return &Client{
		dbPath:  dbPath,
		options: options,
	}, nil
}

type ModeMuxNode string

type clientState struct {
	// State for the WireGuard server.
	localWgServer *LocalWireGuardServer

	// State for the MoleMux client.
	localIPsDB       ezdb.DBRef[string, []ModeMuxNode]
	wgConnections    map[ModeMuxNode]WireGuardConnection
	throttleDetector ThrottleDetector
	mutex            sync.Mutex
}

type WireGuardConnection struct {
	swgpClient  *swgp.ClientConfig
	moleMuxConn net.Conn
	healthy     bool
}

type ThrottleDetector struct {
	// TODO: Implement
}

type LocalWireGuardServer struct {
	device *device.Device
	tun    tun.Device
	ipc    net.Listener
}

type LocalServerConfig struct {
	privateKey string
	listenPort int
}
