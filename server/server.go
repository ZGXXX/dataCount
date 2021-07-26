package main

import (
	"bufio"
	"bytes"
	"fmt"
	pb "github.com/grpc-demo/datacount/protoc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
)

type countService struct {
	pb.UnimplementedCountServiceServer
}

var wg sync.WaitGroup

func main() {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalln(err)
	}
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	pb.RegisterCountServiceServer(grpcServer, &countService{})

	grpcServer.Serve(listen)
}

func (s *countService) ClientData (stream pb.CountService_ClientDataServer) error {
	localFile, err := os.OpenFile("统计文件.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer localFile.Close()
	newWriter := bufio.NewWriter(localFile)
	result := int64(0)
	fmt.Println("开始接收")
	for {
		row, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.CountResponse{Result: result})
		}
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println(row)

		var buffer bytes.Buffer
		buffer.WriteString(row.Name)
		buffer.WriteString("，")
		buffer.WriteString(row.Sex)
		buffer.WriteString("，今年")
		buffer.WriteString(strconv.FormatInt(row.Old, 10))
		buffer.WriteString("岁，")
		buffer.WriteString("来自")
		buffer.WriteString(row.Province)
		tmp := buffer.String()
		fmt.Println(tmp)
		newWriter.WriteString(tmp)
		newWriter.Flush()
		result += 1
	}
	return nil
}