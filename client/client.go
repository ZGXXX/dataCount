package main

import (
	"bufio"
	"context"
	"fmt"
	pb "github.com/grpc-demo/datacount/protoc"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
)

func main() {
	// 开始时间
	t1 := time.Now()
	// 开启grpc客户端连接
	//conn, err := grpc.Dial("10.8.60.20:9002", grpc.WithInsecure(), grpc.WithBlock())
	conn, err := grpc.Dial(":10086", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Println("连接失败")
	}
	defer conn.Close()
	// 处理文件数组
	limit := 100
	ch := make(chan []*pb.User)

	handleArray("source", limit, ch)
	// 启动客户端
	client := pb.NewStatisticServiceClient(conn)
	// 从chan获取值进行发送请求
	// 客户端流式请求
	request, _ := client.DealData(context.Background())
	for items := range ch {
		if len(items) > 0 {
			if err != nil {
				return
			}
			clientStream(request, items, int64(limit))
		}
	}
	res, err := request.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("调用gRPC方法成功，result = ", res.Result)

	// 结束时间
	t3 := time.Now()
	joinSqlTime := t3.Sub(t1)
	fmt.Printf("客户端发送花费时间：%f\n", joinSqlTime.Seconds())
}

func clientStream (request pb.StatisticService_DealDataClient, items []*pb.User, total int64) {
	err := request.Send(&pb.StreamRequest{User : items, Total : total})
	if err != nil {
		fmt.Println("发送失败")
	}
}

func handleArray (filePath string, limit int, ch chan []*pb.User) {
	files, _ := ioutil.ReadDir(filePath)
	go func() {
		for _, f := range files{
			localFile, err := os.Open(filePath + "/" + f.Name())
			if err != nil {
				fmt.Println("打开文件失败")
			}
			defer localFile.Close()

			var requestArr []*pb.User
			var lineName, lineSex, lineAge, lineProvince string
			var lineArray []string

			// 遍历每一行，存到一个大数组中
			fileScan := bufio.NewScanner(localFile)
			for fileScan.Scan() {
				line := fileScan.Text()
				lineArray    = strings.Split(strings.Replace(line, " ", "", -1), ",")
				lineName     = lineArray[0]
				lineSex      = lineArray[1]
				lineAge      = lineArray[2]
				lineProvince = lineArray[3]
				newArr := &pb.User{
					Name:     lineName,
					Age:      lineAge,
					Province: lineProvince,
					Sex:      lineSex,
				}
				requestArr = append(requestArr, newArr)
			}

			allCount := len(requestArr) / limit
			for i:=0;i<=limit;i++ {
				var newArr []*pb.User
				if i == limit {
					newArr = requestArr[allCount * i : ]
				} else {
					newArr = requestArr[allCount * i : allCount * (i+1)]
				}
				ch <- newArr
			}
		}
		close(ch)
	}()
}

func test () {
	// 启动客户端
	conn, err := grpc.Dial("10.8.60.20:9002", grpc.WithInsecure(), grpc.WithBlock())

	client := pb.NewStatisticServiceClient(conn)
	request, _ := client.DealData(context.Background())
	newArr := []*pb.User{
		{Name: "赵国鑫", Sex: "男", Age: "11", Province: "广东省"},
		{Name: "赵国鑫2", Sex: "男", Age: "22", Province: "广东省"},
		{Name: "赵国鑫3", Sex: "男", Age: "22", Province: "广东省"},
		{Name: "赵国鑫3", Sex: "女", Age: "22", Province: "广东省"},
		{Name: "赵国鑫4", Sex: "女", Age: "22", Province: "广东省"},
	}

	err = request.Send(&pb.StreamRequest{User : newArr, Total : 1})
	if err != nil {
		fmt.Println("发送失败")
	}
	res, err := request.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("调用gRPC方法成功，result = ", res.Result)
	return
}