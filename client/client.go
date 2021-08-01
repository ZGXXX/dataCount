package main

import (
	"bufio"
	"context"
	"fmt"
	pb "github.com/grpc-demo/datacount/protoc"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"log"
)

func main() {
	// 开始时间
	t1 := time.Now()
	conn, err := grpc.Dial("10.8.60.20:9002", grpc.WithInsecure(), grpc.WithBlock())
	//conn, err := grpc.Dial(":10086", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
		fmt.Println("连接失败")
	}
	defer conn.Close()

	//// 启动客户端
	//client := pb.NewStatisticServiceClient(conn)
	//request, _ := client.DealData(context.Background())
	//newArr := []*pb.User{
	//	{Name: "赵国鑫", Sex: "男", Age: "11", Province: "广东省"},
	//	{Name: "赵国鑫2", Sex: "男", Age: "22", Province: "广东省"},
	//	{Name: "赵国鑫3", Sex: "男", Age: "22", Province: "广东省"},
	//	{Name: "赵国鑫3", Sex: "女", Age: "22", Province: "广东省"},
	//	{Name: "赵国鑫4", Sex: "女", Age: "22", Province: "广东省"},
	//}
	//
	//err = request.Send(&pb.StreamRequest{User : newArr, Total : 1})
	//if err != nil {
	//	fmt.Println("发送失败")
	//}
	//res, err := request.CloseAndRecv()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println("调用gRPC方法成功，result = ", res.Result)
	//return

	// 本地待处理文件
	localFile, err := os.Open("test2.csv")
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

	t2 := time.Now()
	joinSqlTime := t2.Sub(t1)

	// 使用channel数组类型，将大数组分多次切割放入chan管道
	// 300万数据分100次传grpc，大小已经接近grpc默认接收字节最大值4M了，速率是最快的
	limit := 3
	ch := make(chan []*pb.User)
	allCount := len(requestArr) / limit
	go func() {
		for i:=0;i<=limit;i++ {
			var newArr []*pb.User
			if i == limit {
				newArr = requestArr[allCount * i : ]
			} else {
				newArr = requestArr[allCount * i : allCount * (i+1)]
			}
			ch <- newArr
		}
		close(ch)
	}()


	// 启动客户端
	client := pb.NewStatisticServiceClient(conn)
	// 从chan获取值进行发送请求
	timeNum := 0
	// 客户端流式请求
	request, _ := client.DealData(context.Background())
	for items := range ch {
		if len(items) > 0 {
			if err != nil {
				return
			}
			fmt.Println(items)
			clientStream(request, items, int64(limit))
			timeNum ++
		}
	}
	fmt.Println(timeNum)
	res, err := request.CloseAndRecv()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("调用gRPC方法成功，result = ", res.Result)

	// 结束时间
	t3 := time.Now()
	joinSqlTime2 := t3.Sub(t1)

	fmt.Printf("所有数据读取到数组花费：%f\n", joinSqlTime.Seconds())
	fmt.Printf("客户端发送花费时间：%f\n", joinSqlTime2.Seconds())
}

func clientStream (request pb.StatisticService_DealDataClient, items []*pb.User, total int64) {
	err := request.Send(&pb.StreamRequest{User : items, Total : total})
	if err != nil {
		fmt.Println("发送失败")
	}
}