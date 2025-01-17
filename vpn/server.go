package vpn

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/tun"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func SetupVpn() {
	interfaceName := "wg0"
	listenPort := 51820

	// Generate server private and public keys
	serverPrivateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}
	serverPublicKey := serverPrivateKey.PublicKey()

	fmt.Printf("Server Private Key: %s\n", serverPrivateKey.String())
	fmt.Printf("Server Public Key: %s\n", serverPublicKey.String())

	// Create a TUN
	tunDevice, err := tun.CreateTUN(interfaceName, 1500)
	if err != nil {
		log.Fatalf("Failed to create TUN interface: %v", err)
	}
	defer tunDevice.Close()

	tunDeviceName, err := tunDevice.Name()

	if err == nil {
		fmt.Printf("Created TUN device: %s\n", tunDeviceName)
	}

	lnk, err := netlink.LinkByName(tunDeviceName)
	if err != nil {
		log.Fatal(err)
	}

	ipConfig := &netlink.Addr{IPNet: &net.IPNet{
		IP:   net.ParseIP("10.0.0.1"),
		Mask: net.CIDRMask(24, 32),
	}}

	if err = netlink.AddrAdd(lnk, ipConfig); err != nil {
		log.Fatal(err)
	}

	if err = netlink.LinkSetUp(lnk); err != nil {
		log.Fatal(err)
	}

	// Listen on UDP port
	addr := net.UDPAddr{Port: listenPort, IP: net.IPv4zero}
	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		log.Fatalf("Failed to bind to port %d: %v", listenPort, err)
	}
	defer conn.Close()

	fmt.Printf("Listening on %s:%d\n", addr.IP, addr.Port)

	// Wait for interrupt signal
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	fmt.Println("Shutting down server...")
}

// setIPAddress sets the IP address and netmask for the TUN interface.
func setIPAddress(ifaceName, ip, mask string) error {
	_, err := exec.Command("ifconfig", "wg0", "10.0.0.1", "netmask", "255.255.255.0", "up").Output()
	if err != nil {
		return fmt.Errorf("failed to assign IP: %w", err)
	}
	// cmd = fmt.Sprintf("ip link set dev %s up", ifaceName)
	// _, err = exec.Command("sh", "-c", cmd).Output()
	// if err != nil {
	// 	return fmt.Errorf("failed to bring interface up: %w", err)
	// }
	return nil
}

// sudo ifconfig wg0 10.0.0.1 netmask 255.255.255.0 up
// sudo ifconfig wg0 10.0.0.1 netmask 255.255.255.0 down

// sudo ip addr add 10.1.0.10/24 dev O_O
// sudo ip link set dev O_O up
