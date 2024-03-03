package main

import (
	"fmt"
	"unsafe"
)

const (
	heapSize = 128
	entCount = 40
	header   = 4
)

type vMemory struct {
	heap [heapSize]byte
}

type Entity struct {
	ptr  unsafe.Pointer
	size int32
}

var (
	inUse int = 0
	list  [entCount]Entity
	vm    vMemory
)

func newEntity(size uint8) *Entity {
	if list[0].ptr == nil && list[0].size == 0 {
		list[0].ptr = unsafe.Pointer(&vm.heap)
		list[0].size = heapSize
		inUse++
		log()
	}

	var best *Entity = &list[0]

	for i := 0; i < inUse; i++ {
		if list[i].size >= int32(size) && list[i].size < best.size {
			best = &list[i]
		}
	}

	return best
}

func y_malloc(size uint8) *uintptr {
	if size+header > heapSize {
		return nil
	}

	size += header

	ent := newEntity(size)

	start := ent.ptr
	userPtr := (*uintptr)(unsafe.Pointer(uintptr(start) + header))

	*(*uint8)(start) = size

	ent.ptr = unsafe.Pointer(uintptr(start) + uintptr(size))
	ent.size -= int32(size)

	if ent.size < 0 {
		panic("out of memory")
	}

	log()
	return userPtr
}

func y_free(ptr *uintptr) {
	start := unsafe.Pointer(uintptr(unsafe.Pointer(ptr)) - header)

	list[inUse].ptr = start
	list[inUse].size = int32(*(*uint8)(start))
	inUse++
	log()
}

func log() {
	fmt.Println("OUR LIST")
	for i := 0; i < inUse; i++ {
		fmt.Printf("Data + HEADER. [%p]. Memory of our heap free: [%d]. List i = %d\n", list[i].ptr, list[i].size, i)
	}
	fmt.Printf("Entities in use: [%d]\n", inUse)
}

func test() {
	ptr1 := y_malloc(uint8(unsafe.Sizeof(int(0))))
	*ptr1 = 1
	ptr2 := y_malloc(uint8(unsafe.Sizeof(float32(0))))
	*ptr2 = 2.0
	fmt.Printf("ptr1: [%v]. ptr2: [%v]\n", *ptr1, *ptr2)

	y_free(ptr1)

	ptr3 := y_malloc(uint8(unsafe.Sizeof(int(0))))
	*ptr3 = 3

	fmt.Printf("ptr1: [%v]. ptr2: [%v]. ptr3: [%v]\n", *ptr1, *ptr2, *ptr3)
}

func main() {
	test()
}
