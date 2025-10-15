package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"net"
	"os"
	"time"
)

type ClientInfo struct {
	ID       string `json:"id"`
	HostName string `json:"hostname"`
	IP       string `json:"ip"`
}

var discovered []ClientInfo

// broadcast a UDP packet asking for clients to respond
func DiscoverClients() []ClientInfo {
	discovered = nil
	conn, err := net.ListenPacket("udp4", ":9998")
	if err != nil {
		fmt.Println("UDP listen error:", err)
		return nil
	}
	defer conn.Close()

	bcast, _ := net.ResolveUDPAddr("udp4", "255.255.255.255:9999")
	msg := []byte("DISCOVER_REQUEST")
	conn.WriteTo(msg, bcast)

	conn.SetDeadline(time.Now().Add(3 * time.Second))
	buf := make([]byte, 2048)
	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			break
		}
		var info ClientInfo
		if json.Unmarshal(buf[:n], &info) == nil {
			info.IP = addr.(*net.UDPAddr).IP.String()
			discovered = append(discovered, info)
		}
	}
	return discovered
}

// get system info via TCP
func GetInfo(ip string) string {
	conn, err := net.Dial("tcp", ip+":9000")
	if err != nil {
		return fmt.Sprintf("connect error: %v", err)
	}
	defer conn.Close()
	io.WriteString(conn, `{"cmd":"get_info"}`+"\n")

	reply, _ := io.ReadAll(conn)
	return string(reply)
}

// request screenshot and save it
func GetScreenshot(ip string) string {
	conn, err := net.Dial("tcp", ip+":9000")
	if err != nil {
		return fmt.Sprintf("connect error: %v", err)
	}
	defer conn.Close()
	io.WriteString(conn, `{"cmd":"screenshot"}`+"\n")

	imgFile := fmt.Sprintf("screenshot_%s.jpg", ip)
	out, _ := os.Create(imgFile)
	defer out.Close()
	io.Copy(out, conn)
	return imgFile
}

// Decode JPEG helper (optional preview)
func DecodeScreenshot(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return jpeg.Decode(f)
}
