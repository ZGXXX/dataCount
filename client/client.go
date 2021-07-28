package main

import (
	"bufio"
	"context"
	"fmt"
	pb "github.com/grpc-demo/datacount/protoc"
	"google.golang.org/grpc"
	"log"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	// 开始时间
	t1 := time.Now()
	var requestArr []string

	conn, err := grpc.Dial(":50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
		fmt.Println("连接失败")
	}
	defer conn.Close()


	// 本地待处理文件
	localFile, err := os.Open("test.csv")
	if err != nil {
		fmt.Println("打开文件失败")
	}
	defer localFile.Close()
	// 遍历每一行
	fileScan := bufio.NewScanner(localFile)
	for fileScan.Scan() {
		requestArr = append(requestArr, fileScan.Text())
	}

	limit := 100000
	ch := make(chan []string, 2000)
	go func() {
		allCount := len(requestArr) / limit
		for i:=0;i<=limit;i++ {
			var newArr []string
			if i == limit {
				newArr = requestArr[allCount * i : ]
			} else {
				newArr = requestArr[allCount * i : allCount * (i+1)]
			}
			ch <- newArr
		}
		close(ch)
	}()

	client := pb.NewCountServiceClient(conn)
	for item:= range ch{
		clientStream(client, item)
	}


	// 结束时间
	t2 := time.Now()
	joinSqlTime := t2.Sub(t1)
	fmt.Printf("客户端发送花费时间：%f\n", joinSqlTime.Seconds())

	//test()
}

func clientStream (client pb.CountServiceClient, data []string) {
	// 客户端流式请求
	request, _ := client.ClientData(context.Background())
	err := request.Send(&pb.RowRequest{Line: data})
	if err != nil {
		fmt.Println("发送失败")
	}

	res, err := request.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("调用gRPC方法成功，result = ", res.Result)
}


func test () {
	//本地待处理文件
	localFile, err := os.Open("test.csv")
	if err != nil {
		log.Fatalln(err)
		fmt.Println("打开文件失败")
	}
	defer localFile.Close()

	// 开始时间
	t1 := time.Now()

	var requestArr []string
	ch := make(chan []string, 2000)

	//遍历每一行
	fileScan := bufio.NewScanner(localFile)
	i := 0
	for fileScan.Scan() {
		if len(requestArr) < 100 {
			requestArr = append(requestArr, fileScan.Text())
		}
		if len(requestArr) == 100 {
			i = i + 1
			fmt.Println(i)
			ch <- requestArr
			requestArr = []string{}
		}
	}

	go func() {
		for {
			fileData, ok := <-ch
			if !ok {
				break
			}
			fmt.Println(fileData)
		}
	}()

	// 结束时间
	t2 := time.Now()
	joinSqlTime := t2.Sub(t1)
	fmt.Printf("客户端发送花费时间：%f\n", joinSqlTime.Seconds())

}