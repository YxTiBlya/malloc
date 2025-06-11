package main

import (
	"fmt"
	"unsafe"
)

const (
	// heapSize uint  = 1 << 20 // 1 MB
	heapSize uint  = 1 << 8 // 256 bytes
	header   uint8 = 4
)

type vMemory struct {
	heap [heapSize]byte
}

// not "oop" because polymorphism slower imho (im not tested)
type allocator struct {
	head *node
	data unsafe.Pointer
	used uint
}

type node struct {
	next *node
	prev *node
	size uint
	data unsafe.Pointer
	used bool
}

//

func yCreateAllocator(size uint) *allocator {
	alloc = &allocator{
		data: unsafe.Pointer(&vmem.heap[0]),
	}

	yCreateNode(alloc, size+uint(header), nil, nil, false)

	return alloc
}

var (
	vmem  *vMemory
	alloc *allocator
)

func yCreateNode(alloc *allocator, size uint, outNode **node, lastNode *node, freed bool) {
	newNode := &node{size: size}

	if lastNode != nil {
		newNode.data = unsafe.Pointer((uintptr)(lastNode.data) + (uintptr)(lastNode.size))
		newNode.prev = lastNode
		lastNode.next = newNode
		*outNode = newNode
	} else {
		newNode.data = alloc.data
		alloc.head = newNode
	}

	if !freed {
		alloc.used += uint(size)
	}
}

func yFindBestNode(alloc *allocator, size uint, outNode, outLastNode **node) {
	current := alloc.head
	best := current
	var last *node

	for current != nil {
		if !current.used && current.size >= size && current.size < best.size {
			best = current
			if best.size == size {
				*outNode = best
				return
			}
		}

		last = current
		current = current.next
	}

	if !best.used && best.size >= size {
		*outNode = best
		return
	}

	*outNode = nil
	*outLastNode = last
}

func yDivideNode(alloc *allocator, size uint, outNode **node) {
	newNode := *outNode
	recidueNode := &node{}

	recidueNode.data = unsafe.Pointer(uintptr(newNode.data) + uintptr(size))
	recidueNode.prev = newNode
	recidueNode.next = newNode.next
	recidueNode.size = newNode.size - size

	newNode.next = recidueNode
	newNode.size = size
	alloc.used += uint(size)

	*outNode = newNode
}

func YAlloc(size uint8) unsafe.Pointer {
	if uint(size+header) > heapSize {
		return nil
	}

	if alloc == nil {
		alloc = yCreateAllocator(uint(size))
	}

	return allocate(alloc, 1, uint(size))
}

func YCalloc(count, size uint8) unsafe.Pointer {
	if uint(size*count+header) > heapSize {
		return nil
	}

	if alloc == nil {
		alloc = yCreateAllocator(uint(size * count))
	}

	return allocate(alloc, uint(count), uint(size))
}

// max size of allocation with var size - 255b and count size - 255 = 65025b
func allocate(alloc *allocator, count, size uint) unsafe.Pointer {
	totalSize := size*count + uint(header)

	if uint(totalSize) > heapSize-alloc.used {
		return nil
	}

	var node, lastNode *node
	yFindBestNode(alloc, totalSize, &node, &lastNode)
	if node == nil {
		yCreateNode(alloc, totalSize, &node, lastNode, false)
	}
	if node.size > totalSize {
		yDivideNode(alloc, totalSize, &node)
	}

	node.used = true
	sizeHeader := node.data
	countHeader := unsafe.Pointer(uintptr(sizeHeader) + uintptr(header/2))
	userPtr := unsafe.Pointer(uintptr(sizeHeader) + uintptr(header))

	*(*uint8)(sizeHeader) = uint8(size)
	*(*uint8)(countHeader) = uint8(count)

	return userPtr
}

func YFree(ptr unsafe.Pointer) {
	free(alloc, ptr)
}

func free(alloc *allocator, dataPtr unsafe.Pointer) {
	ptr := unsafe.Pointer(uintptr(dataPtr) - uintptr(header))
	current := alloc.head
	var targetNode, freedNode, prev, last *node

	for current != nil {
		if current.data == ptr {
			targetNode = current
			break
		}

		current = current.next
	}

	current = targetNode
	if current.used {
		alloc.used -= uint(current.size)
	}
	current.used = false

	for current != nil && !current.used {
		last = current
		current = current.next
	}

	current = last
	last = last.next

	var freedMem uint = 0
	for current != nil && !current.used {
		prev = current.prev
		freedMem += current.size
		current = prev
	}

	if current != nil {
		yCreateNode(alloc, freedMem, &freedNode, current, true)
		freedNode.next = last
	}
	if last != nil {
		if current != nil {
			last.prev = freedNode
		} else {
			yCreateNode(alloc, freedMem, nil, nil, true)
			alloc.head.next = last
			last.prev = alloc.head
		}
	}
}

func log() {
	curr := alloc.head
	for i := 1; curr != nil; i++ {
		fmt.Printf("%d) ptr: [%p], count: [%d], values: [", i, curr.data, *(*uint8)(unsafe.Pointer(uintptr(curr.data) + uintptr(header/2))))

		value := unsafe.Pointer(uintptr(curr.data) + uintptr(header))
		for j := 0; j < int(*(*uint8)(unsafe.Pointer(uintptr(curr.data) + uintptr(header/2)))); j++ {
			fmt.Printf(" %d ", *(*int64)(value))
			value = unsafe.Pointer(uintptr(value) + unsafe.Sizeof(int64(0)))
		}
		fmt.Printf("], total size: [%d bytes] used: [%v]\n", curr.size, curr.used)

		curr = curr.next
	}
	fmt.Printf("total memory: %d bytes\nfree memory: %d bytes\nused memory: %d bytes\n", heapSize, heapSize-alloc.used, alloc.used)

	// use heap size like 256 bytes = 1 << 8
	for i := 0; i < 244; i += 12 {
		fmt.Printf("%p | %3d %3d %3d %3d %3d %3d %3d %3d %3d %3d %3d %3d |\n",
			unsafe.Pointer(&vmem.heap[i]),
			vmem.heap[i],
			vmem.heap[i+1],
			vmem.heap[i+2],
			vmem.heap[i+3],
			vmem.heap[i+4],
			vmem.heap[i+5],
			vmem.heap[i+6],
			vmem.heap[i+7],
			vmem.heap[i+8],
			vmem.heap[i+9],
			vmem.heap[i+10],
			vmem.heap[i+11],
		)
	}
}

// todo: fix header, need to set values to the right byte
func main() {
	vmem = new(vMemory)

	ptr1 := YAlloc(uint8(unsafe.Sizeof(int64(0))))
	*(*int64)(ptr1) = 1

	ptr2 := YAlloc(uint8(unsafe.Sizeof(int64(0))))
	*(*int64)(ptr2) = 2

	ptr3 := YAlloc(uint8(unsafe.Sizeof(int64(0))))
	*(*int64)(ptr3) = 3

	YFree(ptr2)
	YFree(ptr1)

	ptr2 = YCalloc(25, uint8(unsafe.Sizeof(int64(0))))
	maxInt64 := int64(^uint(0) >> 1)
	for i := 0; i < 25; i++ {
		*(*int64)(ptr2) = maxInt64
		ptr2 = unsafe.Pointer(uintptr(ptr2) + unsafe.Sizeof(int64(0)))
	}
	*(*int64)(ptr2) = maxInt64

	ptr1 = YAlloc(uint8(unsafe.Sizeof(int64(0))))
	*(*int64)(ptr1) = 5

	log()
}
