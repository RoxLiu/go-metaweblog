package main

import (
	"runtime"
	"time"
	"fmt"
	"strings"
)

// 返回生成自然数序列的管道: 2, 3, 4, ...
func GenerateNatural() chan int {
	ch := make(chan int)
	go func() {
		for i := 2; ; i++ {
			ch <- i
		}
	}()
	return ch
}

// 管道过滤器: 删除能被素数整除的数
func PrimeFilter(in <-chan int, prime int) chan int {
	out := make(chan int)
	go func() {
		for {
			if i := <-in; i%prime != 0 {
				out <- i
			}
		}
	}()
	return out
}

func main() {
	fmt.Println(runtime.GOOS)

	//year, month, day := time.Now().Date()
	//date := fmt.Sprintf("%d-%02d-%02d", year, month, day)
	//fmt.Println(date)

	s := "2017-04-17-Lorem-Ipsum.html"
	array := strings.Split(s, "-")

	fmt.Println(strings.Join(array[:3], "-"))
	t, err := time.Parse("2006-01-02", strings.Join(array[:3], "-"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t)
	//ch := GenerateNatural() // 自然数序列: 2, 3, 4, ...
	//for i := 0; i < 1000; i++ {
	//	prime := <-ch // 新出现的素数
	//	fmt.Printf("%v: %v\n", i+1, prime)
	//	ch = PrimeFilter(ch, prime) // 基于新素数构造的过滤器
	//}
}
