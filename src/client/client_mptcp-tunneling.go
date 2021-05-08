package main

import (
        "flag"
        "fmt"
        "net"
        "os"
        "sync"
        "syscall"
        "time"

        "github.com/google/gopacket"
        "github.com/google/gopacket/pcap"
)

const (
        MPTCP_SCHEDULER         = 43
        MPTCP_PATH_MANAGER      = 44
)

func handleMPTCPDecap(conn net.Conn, handle *pcap.Handle, wg *sync.WaitGroup) {
    defer wg.Done()

        conn.SetReadDeadline( time.Time{} /*time.Now().Add(10 * time.Second)*/ )
        messageBuf := make([]byte, 0xFFFF)
        for {
        //messageBuf := make([]byte, 1518)
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
                client_ip       = flag.String("client_ip", "", "Client IP address")
                pathmanager = flag.String("pathmanager", "ndiffports", "MPTCP pathmanager option")
                scheduler       = flag.String("scheduler", "default", "MPTCP scheduler option")
                server_ip   = flag.String("server_ip", "", "Server IP address")
                server_port = flag.String("server_port", "", "Server port number")

                pcap_device     string = "eth0"
                snapshot_len    int32  = 0xFFFF
                promiscuous     bool   = true
                // timeout              time.Duration = -1 * time.Second
                handle                  *pcap.Handle

                err     error
        )

        flag.Parse()

        my_addr, err := net.ResolveIPAddr("ip4", *client_ip)
        if *client_ip == "" || err != nil {
                fmt.Fprintf(os.Stderr, "Invalid client address.\n")
                os.Exit(1)
        }

        srv_addr, err := net.ResolveTCPAddr("tcp", *server_ip+ ":" + *server_port)
        if *server_ip == "" || *server_port == "" || err != nil {
        fmt.Fprintf(os.Stderr, "Invalid server address.\n")
        os.Exit(1)
        }
        srv_addrStr := *server_ip+ ":" + *server_port

        fmt.Printf("Client start ! (Client: %v <-> Server: %v:%v)\n", my_addr.IP, srv_addr.IP, srv_addr.Port)

        my_dialer := new(net.Dialer)
        my_dialer.Control = func(network, address string, c syscall.RawConn) error {
        var err error
        c.Control(func(fd uintptr) {
                err = syscall.SetsockoptString(int(fd), syscall.SOL_TCP, MPTCP_SCHEDULER, *scheduler)
                        if err == nil {
                                err = syscall.SetsockoptString(int(fd), syscall.SOL_TCP, MPTCP_PATH_MANAGER, *pathmanager)
                }
                })
        return err
    }

        conn, err := my_dialer.Dial("tcp", srv_addrStr)
    if err != nil {
                fmt.Fprintf(os.Stderr, "Failed to dial TCP.\n")
                os.Exit(1)
        }
        defer conn.Close()

        fmt.Printf("Hello, MPTCP Connection !\n")

        handle, err = pcap.OpenLive(pcap_device, snapshot_len, promiscuous, pcap.BlockForever /*timeout*/ )
        if err != nil {
                fmt.Fprintf(os.Stderr, "Failed to open pcap device.\n")
                os.Exit(1)
        }
        defer handle.Close()
        _ = handle.SetBPFFilter("src host 192.168.1.1")

        wg := &sync.WaitGroup{}
        wg.Add(2)

        go handleMPTCPDecap(conn, handle, wg)
        go handleMPTCPEncap(conn, handle, wg)

        wg.Wait()
}