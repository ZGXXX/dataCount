package main

import (
	"bufio"
	"context"
	"fmt"
	pb "github.com/grpc-demo/datacount/protoc"
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

	//conn, err := grpc.Dial(":50051", grpc.WithInsecure(), grpc.WithBlock())
	//if err != nil {
	//	log.Fatalln(err)
	//	fmt.Println("连接失败")
	//}
	//defer conn.Close()


	// 本地待处理文件
	localFile, err := os.Open("source_0.csv")
	if err != nil {
		fmt.Println("打开文件失败")
	}
	defer localFile.Close()
	// 遍历每一行
	fileScan := bufio.NewScanner(localFile)
	for fileScan.Scan() {
		requestArr = append(requestArr, fileScan.Text())
	}
	//client := pb.NewCountServiceClient(conn)
	t2 := time.Now()
	joinSqlTime := t2.Sub(t1)

	limit := 5
	ch := make(chan []string, 2000)
	allCount := len(requestArr) / limit

	go func() {
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

	saveFile, err := os.OpenFile("统计文件.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	defer localFile.Close()
	newWriter := bufio.NewWriter(saveFile)

	countNum := 0
	for items := range ch {
		for _, item := range items {
			newWriter.WriteString(item + "\n")
			countNum ++
		}
	}
	newWriter.Flush()


	// 结束时间
	t3 := time.Now()
	joinSqlTime2 := t3.Sub(t1)

	fmt.Printf("所有数据读取到数组花费：%f\n", joinSqlTime.Seconds())
	fmt.Printf("客户端发送花费时间：%f\n", joinSqlTime2.Seconds())
	fmt.Println(countNum)
	//test()
}

func clientStream (client pb.CountServiceClient, ch chan []string) {
	defer wg.Done()
	// 客户端流式请求
	request, _ := client.ClientData(context.Background())
	data := <-ch
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