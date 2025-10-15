package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"time"
)

// App struct
type App struct {
	ctx context.Context
}

// ClientInfo represents a discovered client
type ClientInfo struct {
	ID       string `json:"id"`
	HostName string `json:"hostname"`
	IP       string `json:"ip"`
	Status   string `json:"status"`
}

// SystemInfo represents system information from a client
type SystemInfo struct {
	Hostname string `json:"hostname"`
	User     string `json:"user"`
	OS       string `json:"os"`
	CPU      int    `json:"cpu"`
	MemMB    int    `json:"memMB"`
}

var discovered []ClientInfo

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// DiscoverClients broadcasts UDP to find available clients
func (a *App) DiscoverClients() []ClientInfo {
	fmt.Println("[DISCOVERY] Starting client discovery...")
	discovered = nil
	conn, err := net.ListenPacket("udp4", ":9998")
	if err != nil {
		fmt.Println("[ERROR] UDP listen error:", err)
		return nil
	}
	defer conn.Close()
	fmt.Println("[DISCOVERY] UDP listener started on :9998")

	// Send broadcast discovery request
	msg := []byte("DISCOVER_REQUEST")
	fmt.Printf("[DISCOVERY] Sending discovery request: %s\n", string(msg))

	// Send to network broadcast only
	bcast, _ := net.ResolveUDPAddr("udp4", "255.255.255.255:9999")
	n1, err1 := conn.WriteTo(msg, bcast)
	fmt.Printf("[DISCOVERY] Broadcast to 255.255.255.255:9999 - sent %d bytes, error: %v\n", n1, err1)

	conn.SetDeadline(time.Now().Add(3 * time.Second))
	buf := make([]byte, 2048)
	fmt.Println("[DISCOVERY] Waiting for responses for 3 seconds...")

	// Use a map to prevent duplicates (same client responding to both broadcasts)
	clientMap := make(map[string]ClientInfo)
	responseCount := 0

	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			fmt.Printf("[DISCOVERY] Read timeout or error: %v\n", err)
			break
		}
		responseCount++
		fmt.Printf("[DISCOVERY] Received response #%d from %s: %s\n", responseCount, addr, string(buf[:n]))

		var info ClientInfo
		if json.Unmarshal(buf[:n], &info) == nil {
			info.IP = addr.(*net.UDPAddr).IP.String()
			info.Status = "online"
			// Use hostname+IP as unique key to prevent duplicates
			key := info.HostName + "_" + info.IP
			clientMap[key] = info
			fmt.Printf("[DISCOVERY] Successfully parsed client: %s at %s\n", info.HostName, info.IP)
		} else {
			fmt.Printf("[ERROR] Failed to parse JSON response: %s\n", string(buf[:n]))
		}
	}

	// Convert map back to slice
	for _, client := range clientMap {
		discovered = append(discovered, client)
	}

	fmt.Printf("[DISCOVERY] Discovery complete. Found %d unique clients\n", len(discovered))
	return discovered
}

// GetClientInfo gets system information from a specific client
func (a *App) GetClientInfo(ip string) SystemInfo {
	conn, err := net.Dial("tcp", ip+":9000")
	if err != nil {
		return SystemInfo{Hostname: "Error: " + err.Error()}
	}
	defer conn.Close()

	io.WriteString(conn, `{"cmd":"get_info"}`+"\n")
	reply, _ := io.ReadAll(conn)

	var info SystemInfo
	json.Unmarshal(reply, &info)
	return info
}

// GetScreenshot captures a screenshot from a client and returns it as base64
func (a *App) GetScreenshot(ip string) string {
	conn, err := net.Dial("tcp", ip+":9000")
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	defer conn.Close()

	io.WriteString(conn, `{"cmd":"screenshot"}`+"\n")

	// Read the JPEG data directly
	imgData, err := io.ReadAll(conn)
	if err != nil {
		return fmt.Sprintf("Error reading screenshot: %v", err)
	}

	// Convert to base64 for frontend display
	return base64.StdEncoding.EncodeToString(imgData)
}

// BrowseForPC simulates browsing for a PC (placeholder)
func (a *App) BrowseForPC(name string) []ClientInfo {
	// For now, just return the discovered clients filtered by name
	var filtered []ClientInfo
	for _, client := range discovered {
		if name == "" || client.HostName == name {
			filtered = append(filtered, client)
		}
	}
	return filtered
}

// StartScreenStream initiates screen streaming from a client with specified framerate
func (a *App) StartScreenStream(ip string, framerate int) string {
	conn, err := net.Dial("tcp", ip+":9000")
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	defer conn.Close()

	cmd := fmt.Sprintf(`{"cmd":"start_stream","framerate":%d}`, framerate)
	io.WriteString(conn, cmd+"\n")
	fmt.Printf("[STREAM] Started stream for %s at %d FPS\n", ip, framerate)
	return "Stream started"
}

// StopScreenStream stops screen streaming from a client
func (a *App) StopScreenStream(ip string) string {
	conn, err := net.Dial("tcp", ip+":9000")
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	defer conn.Close()

	io.WriteString(conn, `{"cmd":"stop_stream"}`+"\n")
	fmt.Printf("[STREAM] Stopped stream for %s\n", ip)
	return "Stream stopped"
}

// GetStreamFrame gets a single frame from streaming client
func (a *App) GetStreamFrame(ip string) string {
	conn, err := net.Dial("tcp", ip+":9000")
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	defer conn.Close()

	io.WriteString(conn, `{"cmd":"stream_frame"}`+"\n")

	// Read the JPEG data directly
	imgData, err := io.ReadAll(conn)
	if err != nil {
		return fmt.Sprintf("Error reading stream frame: %v", err)
	}

	// Convert to base64 for frontend display
	return base64.StdEncoding.EncodeToString(imgData)
}
