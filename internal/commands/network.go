package commands

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/MayR-Labs/mayrlabs-go/internal/utils"
	"github.com/spf13/cobra"
)

// MyIPCmd shows the user's public and local IP addresses
var MyIPCmd = &cobra.Command{
	Use:   "myip",
	Short: "Show your public and local IP addresses",
	Long:  "Display your public IP address (fetched from ipify.org) and your local IP address on the network.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔍 Fetching IP information...")

		// Get Public IP
		publicIP := "Unavailable"
		resp, err := http.Get("https://api.ipify.org")
		if err == nil {
			defer func() { _ = resp.Body.Close() }()
			body, err := io.ReadAll(resp.Body)
			if err == nil {
				publicIP = string(body)
			}
		}

		// Get Local IP
		localIP := "Unavailable"
		addrs, err := net.InterfaceAddrs()
		if err == nil {
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						localIP = ipnet.IP.String()
						break // Just take the first non-loopback IPv4
					}
				}
			}
		}

		fmt.Printf("\n🌍 Public IP: %s\n", publicIP)
		fmt.Printf("🏠 Local IP:  %s\n", localIP)

		return nil
	},
}

// PortCheckCmd checks if a port is open
var PortCheckCmd = &cobra.Command{
	Use:   "port-check [port]",
	Short: "Check if a port is open or in use on localhost",
	Long:  "Check if a specific TCP port is open or in use on the local machine.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var port string
		var err error

		if len(args) > 0 {
			port = args[0]
		} else {
			port, err = utils.PromptInput("Enter port to check: ")
			if err != nil {
				return err
			}
		}

		if port == "" {
			return fmt.Errorf("port cannot be empty")
		}

		target := fmt.Sprintf("localhost:%s", port)
		timeout := time.Second * 1

		conn, err := net.DialTimeout("tcp", target, timeout)
		if err != nil {
			// Connection refused usually means the port is closed/free
			// But strictly speaking, if we can't connect, it's not "open" for us to use?
			// Wait, usually "port check" means "is something running on it?".
			// If Dial succeeds, something IS listening.
			// If Dial fails, it's likely free.

			// Let's clarify the output for the user.
			if strings.Contains(err.Error(), "refused") {
				fmt.Printf("✅ Port %s appears to be FREE (connection refused).\n", port)
			} else {
				fmt.Printf("❓ Port %s status unknown: %v\n", port, err)
			}
			return nil
		}
		defer func() { _ = conn.Close() }()

		fmt.Printf("❌ Port %s is IN USE (connection successful).\n", port)
		return nil
	},
}
