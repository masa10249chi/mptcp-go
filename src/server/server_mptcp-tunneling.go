package main

import (
        "flag"
        "fmt"
        "net"
        "os"
        "sync"
        "time"

        "github.com/google/gopacket"
        "github.com/google/gopacket/pcap"
)

func handleMPTCPDecap(conn net.Conn, handle *pcap.Handle, wg *sync.WaitGroup) {
        defer wg.Done()

        conn.SetReadDeadline( time.Time{} /*time.Now().Add(10 * time.Second)*/ )
        messageBuf := make([]byte, 0xFFFF)
        for {
                // messageBuf := make([]byte, 0xFFFF)
                messageLen, err := conn.Read(messageBuf)
                if err != nil {
                fmt.Println(err)
                        fmt.Fprintf(os.Stderr, "Wrong message.\n")
                os.Exit(1)
                }
                err = handle.WritePacketData(messageBuf[:messageLen])
        if err != nil {
                        fmt.Fprintf(os.Stderr, "Failed to decap message.\n")
                        os.Exit(1)
        }
        }
}

func handleMPTCPEncap(conn net.Conn, handle *pcap.Handle, wg *sync.WaitGroup) {
        defer wg.Done()

        conn.SetWriteDeadline( time.Time{} /*time.Now().Add(10 * time.Second)*/ )

        packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
        for packet := range packetSource.Packets() {
                _, err := conn.Write(packet.Data())
        if err != nil {
                fmt.Fprintf(os.Stderr, "Failed to encap message.\n")
                os.Exit(1)
        }
        }
}

func main() {
    var  (
                server_ip   = flag.String("server_ip", "", "Server IP address")
                server_port = flag.String("server_port", "", "Server port number")

                pcap_device     string = "eth0"
                snapshot_len    int32  = 0xFFFF
                promiscuous     bool   = true
                // timeout         time.Duration = -1 * time.Second
                handle          *pcap.Handle

                err error
    )

        flag.Parse()

        srv_addr, err := net.ResolveTCPAddr("tcp", *server_ip+ ":" + *server_port)
        if *server_ip == "" || *server_port == "" || err != nil {
        fmt.Fprintf(os.Stderr, "Invalid server address.\n")
        os.Exit(1)
        }

        fmt.Printf("Server start ! (Server: %v:%v)\n", srv_addr.IP, srv_addr.Port)

        listner, err := net.ListenTCP("tcp", srv_addr)
        if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to listen TCP.\n")
        os.Exit(1)
        }

        handle, err = pcap.OpenLive(pcap_device, snapshot_len, promiscuous, pcap.BlockForever /*timeout*/ )
        if err != nil {
                fmt.Fprintf(os.Stderr, "Failed to open pcap device.\n")
                os.Exit(1)
        }
        defer handle.Close()
        _ = handle.SetBPFFilter("src host 192.168.1.2")

        for {
                conn, err := listner.Accept()
                if err != nil {
                        continue
                }

                go handleConnection(conn, handle)
        }
}

func handleConnection(conn net.Conn, handle *pcap.Handle) {
        defer conn.Close()

        fmt.Println("Client accept!")

        wg := &sync.WaitGroup{}
        wg.Add(2)

        go handleMPTCPDecap(conn, handle, wg)
        go handleMPTCPEncap(conn, handle, wg)

        wg.Wait()
}