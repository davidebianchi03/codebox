package vpn

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

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
