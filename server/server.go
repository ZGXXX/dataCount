package main

import (
	"bufio"
	"fmt"
	pb "github.com/grpc-demo/datacount/protoc"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"os"
)

type statisticService struct {
	pb.UnimplementedStatisticServiceServer
}

func main() {
	listen, err := net.Listen("tcp", ":10086")
	if err != nil {
		log.Fatalln(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterStatisticServiceServer(grpcServer, &statisticService{})

	grpcServer.Serve(listen)
}

func (s *statisticService) DealData (stream pb.StatisticService_DealDataServer) error {
	localFile, err := os.OpenFile("统计文件.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer localFile.Close()
	newWriter := bufio.NewWriter(localFile)
	fmt.Println("开始接收")
	for {
		row, err := stream.Recv()
		result := "处理成功"
		if err == io.EOF {
			return stream.SendAndClose(&pb.StreamResponse{Result: result})
		}
		if err != nil {
			log.Fatalln(err)
		}
		for _, item := range row.User {
			newWriter.WriteString(item.Name + item.Sex + item.Age + item.Province + "\n")
		}
		newWriter.Flush()
	}
	return nil
}