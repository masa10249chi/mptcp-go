FROM golang:1.16

WORKDIR /go
RUN go get github.com/google/gopacket

RUN apt-get update && apt-get install -y vim libpcap-dev strace
RUN curl -fLo /root/.vim/autoload/plug.vim --create-dirs https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim
RUN git clone https://github.com/fatih/vim-go.git /root/.vim/pack/plugins/start/vim-go

COPY .vimrc /root/

WORKDIR /go/src
COPY src/client/ client/
COPY src/server/ server/

WORKDIR /go/src/client
RUN go build -o ../../bin/client_mptcp-tunneling client_mptcp-tunneling.go

WORKDIR /go/src/server
RUN go build -o ../../bin/server_mptcp-tunneling server_mptcp-tunneling.go
