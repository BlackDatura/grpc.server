package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/blackdatura/grpc.server/pb"
)

const PORT = ":5001"

type employeeService struct {
	pb.UnimplementedEmployeeServiceServer
}

// Unary 一元方法
func (emp *employeeService) GetByNo(ctx context.Context,
	req *pb.GetByNoRequest) (*pb.EmployeeResponse, error) {
	for _, e := range employees {
		if req.No == e.No {
			return &pb.EmployeeResponse{
				Employee: e,
			}, nil
		}
	}
	err := status.Error(codes.NotFound, "employee not found for No. "+strconv.Itoa(int(req.No)))
	return nil, err
}

// Server streaming 服务端流
func (emp *employeeService) GetAll(req *pb.GetAllRequest,
	stream pb.EmployeeService_GetAllServer) error {
	for _, e := range employees {
		err := stream.Send(&pb.EmployeeResponse{
			Employee: e,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Client Streaming 客户端流
func (emp *employeeService) AddPhoto(stream pb.EmployeeService_AddPhotoServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if ok {
		fmt.Printf("Employee: %s\n", md["no"])
	}
	bytesRecv := []byte{}
	for {
		data, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("File size: %d\n", len(bytesRecv))
			img, err := os.Create("code.bmp")
			if err != nil {
				return err
			}
			_, err = img.Write(bytesRecv)
			if err != nil {
				return err
			}
			img.Close()
			return stream.SendAndClose(&pb.AddPhotoResponse{IsOK: true})
		}
		if err != nil {
			return err
		}
		fmt.Printf("%d bytes data received: \n", len(data.Data))
		bytesRecv = append(bytesRecv, data.Data...)
	}
}

func (emp *employeeService) Save(ctx context.Context,
	req *pb.EmployeeRequest) (*pb.EmployeeResponse, error) {
	return nil, nil
}

// Bidirectional streaming 双向流
func (emp *employeeService) SaveAll(stream pb.EmployeeService_SaveAllServer) error {
	for {
		emp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		employees = append(employees, emp.Employee)
		err = stream.Send(&pb.EmployeeResponse{Employee: emp.Employee})
		if err != nil {
			return err
		}
	}
	for _, e := range employees {
		fmt.Println(e)
	}
	return nil
}

func main() {
	listen, err := net.Listen("tcp", PORT)
	if err != nil {
		log.Fatalf("fail to listen: %v", err)
	}
	creds, err := credentials.NewServerTLSFromFile("certs/cert.pem", "certs/key.pem")
	if err != nil {
		log.Fatalf("fail to load cert files: %v", err)
	}
	//不启用https
	//serverOptions:=[]grpc.ServerOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	serverOptions := []grpc.ServerOption{grpc.Creds(creds)}
	server := grpc.NewServer(serverOptions...)
	pb.RegisterEmployeeServiceServer(server, new(employeeService))
	log.Println("gRPC server started..." + PORT)
	server.Serve(listen)
}
