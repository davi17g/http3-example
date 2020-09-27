# http3-example
Simple http/3 Client/Server implementation example written in golang

## Getting Started

### Prerequisite
In order to run this example you need to Generate SSL Certificate files. SSL file generation example can be found in [quic-go repository](https://github.com/lucas-clemente/quic-go/blob/master/internal/testdata/generate_key.sh).

### Run
#### Server:
- Navigate to server directory: `cd server`
- Build: `go build` 
- Execute: `./server --addr <ip-address> --port <port> --cert <certificate-file-path> --key <private-key-file-path>`

#### Client:
- Navigate to client directory: `cd client`
- Build: `go build`
- Execute: `./client --addr <server-address> --port <server-port> --keys <certificate-authority-file-path>`
