//go:generate protoc -I=./proto --go_out=./proto ./proto/problem.proto --go-grpc_out=./proto ./proto/problem.proto
//go:generate openssl req -new -x509 -sha256 -key problserv.key -out problserv.crt -days 3650 -subj "/CN=problemservice" -addext "subjectAltName = DNS:problemservice"
//go:generate openssl req -new -x509 -sha256 -key problserv.key -out probllocal.crt -days 3650 -subj "/CN=localhost" -addext "subjectAltName = DNS:localhost"

package main

import (
	"ProblemMicro/configs"
	"ProblemMicro/proto"
	"ProblemMicro/service"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

const MicroName = "problem_service"

func main() {
	connectionDB := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		configs.PG_HOST,
		configs.PG_PORT,
		configs.POSTGRES_USER,
		configs.POSTGRES_PASSWORD,
		configs.POSTGRES_DB)

	db, err := sql.Open("postgres", connectionDB)
	if err != nil {
		log.Panicf("%s: failed to open db connection - %v", MicroName, err)
	}
	defer db.Close()

	service := service.NewProblemService(db)

	listener, err := net.Listen("tcp", net.JoinHostPort("", configs.GRPC_PORT))
	if err != nil {
		log.Panicf("%s: failed to listen on port - %v", MicroName, err)
	}

	creds, err := credentials.NewServerTLSFromFile(configs.CERTIFICATE, configs.KEY_PRIVATE)
	if err != nil {
		log.Panicf("%s: can't load TLS keys : %v", MicroName, err)
	}

	server := grpc.NewServer(grpc.Creds(creds))
	defer server.GracefulStop()
	proto.RegisterProblemServiceServer(server, service)
	reflection.Register(server)

	if err := server.Serve(listener); err != nil {
		log.Panicf("%s: failed to start grpc - %v", MicroName, err)
	}
}
