package hardware

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os/exec"
	"path"
	"sync"
	"syscall"

	"github.com/spf13/viper"
)

// System is used to control the controller hardware (Raspberry Pi).
type System struct {
	ethMode    string
	ethIP      net.IP
	ethMask    net.IPMask
	ethGateway net.IP
	ethDNS     net.IP

	ethInterface string
	ethMux       sync.Mutex
}

// EthernetConfig contains the ethernet configuration in a readable way
type EthernetConfig struct {
	// Mode can be: down, static, dhcp-client, dhcp-server
	Mode string `json:"mode"`

	// IP is the IP address and subnet in CIDR notation
	IP string `json:"ip"`

	// Gateway is the IP address of the gateway
	Gateway string `json:"gateway"`

	// DNS is the IP address of the gateway
	DNS string `json:"dns"`
}

const (
	// EthDown means that the link is down
	EthDown = "down"

	// EthStatic means that the IP address, gateway and DNS server are configured manually
	EthStatic = "static"

	// EthDhcpClient means that the network configuration is obtained using DHCP
	EthDhcpClient = "dhcp-client"

	// EthDhcpServer means that the controller is running a DHCP server
	EthDhcpServer = "dhcp-server"

	// EthUnmanaged means that the ethernet configuration is not managed by the controller
	EthUnmanaged = "unmanaged"
)

// NewSystem creates a new System struct.
func NewSystem() *System {
	system := &System{}

	// network configuration
	system.ethInterface = viper.GetString("ethernet")

	if system.ethInterface == "" {
		system.ethMode = EthUnmanaged
	} else {
		err := system.loadEthernetConfig()
		if err != nil {
			// loading the configuration failed -> take interface down
			system.ethMode = EthDown
			go system.reconfigureEthernet()
		}
	}

	return system
}

// Shutdown initiates a shutdown of the controller hardware.
func (system *System) Shutdown() {
	err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
	if err != nil {
		log.Print("Failed to poweroff controller: " + err.Error())
	}
}

// EthernetConfig returns the currently configured ethernet settings.
func (system *System) EthernetConfig() EthernetConfig {
	var ip, gateway, dns string

	if system.ethIP != nil {
		ip = system.ethIPNet()
	}

	if system.ethGateway != nil {
		gateway = system.ethGateway.String()
	}

	if system.ethDNS != nil {
		dns = system.ethDNS.String()
	}

	return EthernetConfig{system.ethMode, ip, gateway, dns}
}

// SetEthernetConfig stores and applies the new network configuration.
func (system *System) SetEthernetConfig(c EthernetConfig) error {
	// check if disabled
	if system.ethMode == EthUnmanaged {
		return errors.New("Ethernet configuration is unmanaged")
	}

	// first, validate all parameters
	// check mode
	if c.Mode != EthDown && c.Mode != EthStatic && c.Mode != EthDhcpClient && c.Mode != EthDhcpServer {
		return errors.New("Ethernet configuration has invalid mode")
	}

	var newIP net.IP
	var newMask net.IPMask
	var newGateway net.IP
	var newDNS net.IP

	// IP address, gateway and dns are only relevant for static configuration or dhcp server
	if c.Mode == EthStatic || c.Mode == EthDhcpServer {
		// check IP
		if c.Mode == EthStatic || c.Mode == EthDhcpServer {
			parsedIPAddr, parsedIPNet, err := net.ParseCIDR(c.IP)
			if err != nil {
				return errors.New("Ethernet configuration has invalid ip address or subnet")
			}
			if len(parsedIPAddr.To4()) != net.IPv4len {
				return errors.New("Ethernet configuration does not provide an IPv4 address")
			}
			newIP = parsedIPAddr
			newMask = parsedIPNet.Mask
		}

		// check gateway (optional)
		if c.Gateway != "" {
			parsedGateway := net.ParseIP(c.Gateway)
			if parsedGateway == nil {
				return errors.New("Ethernet configuration has invalid gateway")
			}
			newGateway = parsedGateway
		}

		// check dns (optional)
		parsedDNS := net.ParseIP(c.DNS)
		if parsedDNS == nil {
			return errors.New("Ethernet configuration has invalid DNS server")
		}
		newDNS = parsedDNS
	}

	// TODO
	if c.Mode == EthDhcpClient || c.Mode == EthDhcpServer {
		return errors.New("Ethernet configuration with DHCP not implemented yet")
	}

	system.ethMux.Lock()

	// next, check for changes. if we do not have a change, we do not reconfigure the network
	// in dhcp client mode, do not compare ip, dns and dns since they are set by the controller
	changed := true
	if system.ethMode == c.Mode && c.Mode == EthDhcpClient {
		changed = false
	} else if system.ethMode == c.Mode &&
		system.ethIP.Equal(newIP) &&
		bytes.Compare(system.ethMask, newMask) == 0 &&
		system.ethGateway.Equal(newGateway) &&
		system.ethDNS.Equal(newDNS) {
		changed = false
	}

	if changed {
		// now everything is validated and there was some change, so set the new configuration
		system.ethMode = c.Mode
		system.ethIP = newIP
		system.ethMask = newMask
		system.ethGateway = newGateway
		system.ethDNS = newDNS
	}

	system.ethMux.Unlock()

	if changed {
		system.saveEthernetConfig()
		go system.reconfigureEthernet()
	}

	return nil
}

func (system *System) reconfigureEthernet() {
	if system.ethMode == EthUnmanaged {
		return
	}

	system.ethMux.Lock()

	var err error
	// TODO: stop running dhcp server or client

	// flush the ip addresses
	err = exec.Command("ip", "addr", "flush", "dev", system.ethInterface).Run()
	if err != nil {
		log.Print("Failed to flush network interface: " + err.Error())
	}

	// remove the default gateway
	err = exec.Command("ip", "route", "del", "default").Run()
	if err != nil {
		log.Print("Failed to remove default gateway: " + err.Error())
	}

	// set interface state up or down
	state := "up"
	if system.ethMode == EthDown {
		state = "down"
	}
	err = exec.Command("ip", "link", "set", "dev", system.ethInterface, state).Run()
	if err != nil {
		log.Print("Failed to bring network interface " + state + ": " + err.Error())
	}

	// for static and dhcp server mode, set IP, gateway and dns
	if system.ethMode == EthStatic || system.ethMode == EthDhcpServer {
		// ip
		err = exec.Command("ip", "addr", "add", system.ethIPNet(), "dev", system.ethInterface).Run()
		if err != nil {
			log.Print("Failed to set IP address: " + err.Error())
		}

		// gateway
		err = exec.Command("ip", "route", "add", "default", "via", system.ethGateway.String(), "dev", system.ethInterface).Run()
		if err != nil {
			log.Print("Failed to set default gateway: " + err.Error())
		}

		// dns
		dnsConfig := []byte("nameserver " + system.ethDNS.String() + "\n")
		err = ioutil.WriteFile("/etc/resolv.conf", dnsConfig, 0644)
		if err != nil {
			log.Print("Failed to set DNS server: " + err.Error())
		}
	}

	// TODO: for dhcp client mode, run the dhcp client

	// TODO: for dhcp server mode, run the dhcp server

	system.ethMux.Unlock()
}

// saveEthernetConfig stores the ethernet configuration in a file
func (system *System) saveEthernetConfig() error {
	system.ethMux.Lock()

	data, err := json.MarshalIndent(system.EthernetConfig(), "", "    ")
	if err != nil {
		log.Print("Error while serializing JSON to store network configuration")
		return err
	}

	file := path.Join(viper.GetString("directories.config"), "ethernet.json")
	err = ioutil.WriteFile(file, data, 0644)
	if err != nil {
		log.Print("Failed to write the ethernet configuration to the config file: " + err.Error())
		return err
	}

	system.ethMux.Unlock()
	return nil
}

// loadEthernetConfig loads the ethernet configuration from a file
func (system *System) loadEthernetConfig() error {
	system.ethMux.Lock()

	file := path.Join(viper.GetString("directories.config"), "ethernet.json")
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Print("Failed to read ethernet configuration from config file: " + err.Error())
		return err
	}

	config := EthernetConfig{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Print("Malformed ethernet configuration in config file")
		return err
	}

	system.ethMux.Unlock()

	err = system.SetEthernetConfig(config)
	if err != nil {
		log.Print("Cannot load ethernet configuration from config file: " + err.Error())
		return err
	}

	return nil
}

// ethIPNet returns the ip address and subnet in CIDR notation
func (system *System) ethIPNet() string {
	if system.ethIP != nil {
		subnet, _ := system.ethMask.Size()
		return fmt.Sprintf("%s/%d", system.ethIP.String(), subnet)
	}

	return ""
}
