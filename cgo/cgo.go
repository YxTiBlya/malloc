package main

//#include "y_malloc.h"
import "C"
import (
	"fmt"
	"unsafe"
)

func y_malloc(size uint8) *uintptr {
	return (*uintptr)(C.y_malloc(C.ulonglong(size)))
}

func y_free(ptr *uintptr) {
	C.y_free(unsafe.Pointer(ptr))
}

func main() {
	ptr1 := y_malloc(uint8(unsafe.Sizeof(int(0))))
	*ptr1 = 1
	ptr2 := y_malloc(uint8(unsafe.Sizeof(int(0))))
	*ptr2 = 2
	ptr3 := y_malloc(uint8(unsafe.Sizeof(int(0))))
	*ptr3 = 3

	fmt.Printf("ptr1: [%d]. ptr2: [%d]. ptr3: [%d]\n", *ptr1, *ptr2, *ptr3)

	y_free(ptr2)
	y_free(ptr1)

	ptr1 = y_malloc(uint8(unsafe.Sizeof(int(0))))
	*ptr1 = 4

	fmt.Printf("ptr1: [%v]. ptr2: [%v]. ptr3: [%v]\n", *ptr1, *ptr2, *ptr3)
}
