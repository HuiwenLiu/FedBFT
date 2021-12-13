package util

import "fmt"

func main(){
	done := make(chan int, 2)
	i := 0
	var add = func(){
		for j:=0; j < 10000; j++ {
			i++
		}
		done <- 1
	}
	go add()
	go add()
	<- done
	<- done
	fmt.Printf("%d", i)
}