package main

import (
	"bufio"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strings"
	"time"
)

// This is the 1024-bit MODP group from RFC 5114, section 2.1:
const primeHex = "B10B8F96A080E01DDE92DE5EAE5D54EC52C99FBCFB06A3C69A6A9DCA52D23B616073E28675A23D189838EF1E2EE652C013ECB4AEA906112324975C3CD49B83BFACCBDD7D90C4BD7098488E9C219A73724EFFD6FAE5644738FAA31A4FF55BCCC0A151AF5F0DC8B4BD45BF37DF365C1A65E68CFDA76D4DA708DF1FB2BC2E4A4371"
const generatorHex = "A4D1CBD5C3FD34126765A442EFB99905F8104DD258AC507FD6406CFF14266D31266FEA1E5C41564B777E690F5504F213160217B4B01B886A5E91547F9E2749F4D7FBD7D3B9A92EE1909D0D2263F80A76A6A24C087A091F531DBF0A0169B6A28AD662A4D18E73AFA32D779D5918D08BC8858F4DCEF97C2A24855E6EEB22B3B2E5"

var encryptedData string

func clientSendMessageAndListen() {
	//开启客户端的本地监听（主要用来接收节点的reply信息）
	go clientTcpListen()
	fmt.Printf("客户端开启监听，地址：%s\n", clientAddr)

	fmt.Println(" ---------------------------------------------------------------------------------")
	fmt.Println("|  已进入PBFT测试Demo客户端，请启动全部节点后再发送消息！ :)  |")
	fmt.Println(" ---------------------------------------------------------------------------------")
	fmt.Println("请在下方输入要存入节点的信息：")
	//首先通过命令行获取用户输入
	stdReader := bufio.NewReader(os.Stdin)
	priv := &PrivateKey{
		PublicKey: PublicKey{
			G: fromHex(generatorHex),
			P: fromHex(primeHex),
		},
		X: fromHex("42"),
	}
	priv.Y = new(big.Int).Exp(priv.G, priv.X, priv.P)

	for {
		data, err := stdReader.ReadString('\n')
		values := []rune{}
		for _, value := range data {
			values = append(values, value)
		}
		fmt.Println("VALUES SLICE：", values)

		if err != nil {
			fmt.Println("Error reading from stdin")
			panic(err)
		}
		U.main_U(data)
		message := []byte(data)
		err = ioutil.WriteFile("output.txt", message, 0644)
		if err != nil {
			panic(err)
		}
		c1, c2, err := Encrypt(rand.Reader, &priv.PublicKey, message)
		fmt.Println("c1: ", c1)
		fmt.Println("c2: ", c2)
		if err != nil {
			//t.Errorf("error encrypting: %s", err)
			fmt.Println("error encrypting: %s", err)
		}

		fmt.Println("c1: ", c1)
		fmt.Println("c2: ", c2)
		encryptedData := c1.String() + "+" + c2.String()
		fmt.Println("c1+c2: ", encryptedData)
		//t1 := big.NewInt(c1)
		//t2 := big.NewInt(c2)
		//fmt.Println("t1: ", t1)
		//fmt.Println("t2: ", t2)

		r := new(Request)
		r.Timestamp = time.Now().UnixNano()
		r.ClientAddr = clientAddr
		r.Message.ID = getRandom()
		//消息内容就是用户的输入
		r.Message.Content = strings.TrimSpace(encryptedData)
		br, err := json.Marshal(r)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println(string(br))
		content := jointMessage(cRequest, br)
		//默认N0为主节点，直接把请求信息发送至N0
		tcpDial(content, nodeTable["N0"])
	}
}

//返回一个十位数的随机数，作为msgid
func getRandom() int {
	x := big.NewInt(10000000000)
	for {
		result, err := rand.Int(rand.Reader, x)
		if err != nil {
			log.Panic(err)
		}
		if result.Int64() > 1000000000 {
			return int(result.Int64())
		}
	}
}

func fromHex(hex string) *big.Int {
	n, ok := new(big.Int).SetString(hex, 16)
	if !ok {
		panic("failed to parse hex number")
	}
	return n
}
