package main

import (
	"bufio"
	bytes2 "bytes"
	"io"
	"log"
	"os"
)

func main() {
	dataFile, err := os.Open("source.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer dataFile.Close()

	data := bufio.NewReader(dataFile)
	var bt bytes2.Buffer
	localFile, err := os.OpenFile("本地文件.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	defer localFile.Close()

	for {
		line, _, err := data.ReadLine()
		if err != nil || err == io.EOF {
			break
		}
		str := string(line) + "\n"
		bt.WriteString(str)
		localFile.WriteString(bt.String())
	}
}