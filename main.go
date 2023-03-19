package main

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"time"
)

type list[T elem] struct {
	start *node[T]
	last  *node[T]
}
type elem interface {
	int | byte | string
}
type node[T elem] struct {
	value T
	next  *node[T]
}

func recursionView[T elem](el *node[T]) {
	if el != nil {
		fmt.Println(el.value)
		recursionView(el.next)
	}
}

func (el *node[T]) recursionViewWithReceiver() {
	if el != nil {
		fmt.Println(el.value)
		if el.next != nil {
			el.next.recursionViewWithReceiver()
		}
	}
}

func (l *list[T]) isEmpty() bool {
	return l.start != nil
}
func (l *list[T]) View() {

	for s := l.start; s != nil; s = s.next {
		fmt.Println(s.value)
	}
}

func (l *list[T]) add(value T) {

	if l.start == nil {
		l.start = &node[T]{value: value}
		l.last = l.start
	} else {
		l.last.next = &node[T]{value: value}
		l.last = l.last.next
	}
}

func addElems[T elem](l *list[T], start int, end int, c chan int) {
	for i := start; i < end; i++ {
		l.add(T(i))
	}
	c <- 1
}

func (l *list[T]) addAllElemsGOROUTINE(N int) {
	var numCPU = runtime.NumCPU() * 2
	c := make(chan int, numCPU)
	lists := make([]*list[T], numCPU)
	for i := 0; i < numCPU; i++ {
		lists[i] = new(list[T])
		lists[i].add(T(1))
	}
	var current int = 0
	for i := 0; i < numCPU; i++ {

		go addElems(lists[i], current, current+N/numCPU, c)
		current += N / numCPU
	}
	for i := 0; i < numCPU; i++ {
		<-c
	}

	l.add(T(1))
	for i := 0; i < numCPU; i++ {
		l.last.next = lists[i].start
	}

}

func (l *list[T]) addAllElems(N int) {
	for i := 0; i < N; i++ {
		l.add(T(i))
	}
}

func GetFunctionName(temp interface{}) string {
	strs := strings.Split((runtime.FuncForPC(reflect.ValueOf(temp).Pointer()).Name()), ".")
	return strs[len(strs)-1]
}

func timer(f func(), name string) {
	startTime := time.Now()
	f()
	endTime := time.Now().Sub(startTime).Seconds()
	fmt.Println(name, endTime)
}

func main() {
	var l1 list[int]
	var l2 list[int]
	const N int = 1000000000
	var numCPU = runtime.GOMAXPROCS(2)
	fmt.Println(numCPU)
	fmt.Println(numCPU)
	timer(func() { l1.addAllElemsGOROUTINE(N) }, "addAllElemsGOROUTINE")
	timer(func() { l2.addAllElems(N) }, "addAllElems")

	//l.View()
	fmt.Println("Готово")
	fmt.Scan(new(string))
}
