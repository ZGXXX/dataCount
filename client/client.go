package main

import (
	"bufio"
	"context"
	"fmt"
	pb "github.com/grpc-demo/datacount/protoc"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
		fmt.Println("连接失败")
	}
	defer conn.Close()

	client := pb.NewCountServiceClient(conn)
	clientStream(client, "source_0.csv")
}

func clientStream (client pb.CountServiceClient, fileName string) error {
	// 客户端流式请求
	data, _ := client.ClientData(context.Background())
	// 本地待处理文件
	localFile, err := os.Open(fileName)
	if err != nil {
		log.Fatalln(err)
		fmt.Println("打开文件失败")
	}
	defer localFile.Close()
	bufferRead := bufio.NewReader(localFile)

	// 开始时间
	t1 := time.Now()

	 ch := make(chan string, 10000)
	// 遍历每一行
	for {
		line, err := bufferRead.ReadString('\n')
		ch <- line
		wg.Add(1)

		if err == io.EOF {
			go handle(data, line, ch)
			fmt.Println("文件读取完成")
			break
		}
		go handle(data, line, ch)
	}
	wg.Wait()
	res, err := data.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
	}
	// 结束时间
	t2 := time.Now()
	joinSqlTime := t2.Sub(t1)
	fmt.Printf("客户端发送花费时间：%f\n", joinSqlTime.Seconds())

	fmt.Println("客户端接收 Recv 次数:", res.Result)
	return nil
}

func handle(requestData pb.CountService_ClientDataClient,line string,ch chan string) {
	defer wg.Done()
	var lineName string
	var lineSex string
	var lineOld int64
	var lineProvince string
	var lineArray []string

	lineArray    = strings.Split(strings.Replace(line, " ", "", -1), ",")
	lineName     = lineArray[0]
	lineSex      = lineArray[1]
	lineOld, _   = strconv.ParseInt(lineArray[2], 10, 64)
	lineProvince = lineArray[3]

	fmt.Println("客户端开始发送 Stream :", line)
	err := requestData.Send(&pb.RowRequest{Name: lineName, Sex: lineSex, Old: lineOld, Province: lineProvince})
	if err != nil {
		fmt.Println("发送失败")
	}
	time.Sleep(time.Millisecond)
	<-ch
}