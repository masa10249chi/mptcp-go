# mptcp-go
Golang based Multipath TCP v0 tunneling works on Docker.
 
# Requirement
* Multipath TCP Linux Kernel (latest: https://github.com/multipath-tcp/mptcp/releases/tag/v0.95.1)
* Docker-CE

# Installation
```bash
git clone https://github.com/masa10249chi/mptcp-go.git
cd mptcp-go
docker build -t mptcp-go .
```

# Usage
Client side:
```bash
docker run -it --name mptcp-go-client --privileged mptcp-go \
           /go/bin/client_mptcp-tunneling -client_ip xx.xx.xx.xx -server_ip xx.xx.xx.xx -server_port xxxx \
           -pathmanager {default|fullmesh|ndiffports|binder} -scheduler {default|roundrobin|ndiffports|redundant}
```
Server side:
```bash
docker run -it --name mptcp-go-server --privileged mptcp-go \
           /go/bin/server_mptcp-tunneling -server_ip xx.xx.xx.xx -server_port xxxx
```
