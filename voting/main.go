package main

import (
	"fmt"
	"log"
	"os"
	"time"

	C "./client"
)

const nodeCount = 4

//客户端的监听地址
var clientAddr = "127.0.0.1:8888"

//节点池，主要用来存储监听地址
var nodeTable map[string]string

func main() {
	//为四个节点生成公私钥
	// C.Main_c("x")
	genRsaKeys()
	nodeTable = map[string]string{
		"N0": "127.0.0.1:8000",
		"N1": "127.0.0.1:8001",
		"N2": "127.0.0.1:8002",
		"N3": "127.0.0.1:8003",
	}
	if len(os.Args) != 2 {
		log.Panic("Incorrect input parameters")
	}

	nodeID := os.Args[1]
	if nodeID == "client" {
		user_p := C.Main_c(nodeID)
		AsyncServiceOut(user_p, 2e9)
		clientSendMessageAndListen() //启动客户端程序
	} else if addr, ok := nodeTable[nodeID]; ok {
		t12_p := C.Main_c(nodeID)
		AsyncServiceOut(t12_p, 2e9)
		p := NewPBFT(nodeID, addr)
		go p.tcpListen() //启动节点
	} else {
		log.Fatal("无此节点编号！")
	}
	select {}
}

func AsyncService(service func() string) chan string {
	retCh := make(chan string, 1)
	go func() {
		ret := service()
		fmt.Println("service()执行结束.")
		retCh <- ret
		fmt.Println("service()返回值塞进通道.")
	}()
	return retCh
}

func AsyncServiceOut(service func() string, duration time.Duration) {
	select {
	case ret := <-AsyncService(service):
		fmt.Println("====", ret)
	case <-time.After(duration):
		fmt.Println("time out")
	}
}
