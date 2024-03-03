#include "y_malloc.h"

static u16 IN_USE;

entity_t LIST[ENT_COUNT];

void LOG()
{
	printf("OUR LIST\n");
	for (unsigned i = 0; i < IN_USE; i++)
	{
		printf("Data + HEADER. [%p]. Memory of our heap free: [%u]. List i = %d\n", LIST[i].ptr, LIST[i].size, i);
	}
	printf("Entities in use:[%d]\n", IN_USE);
}

entity_t* new_entity(size_t size)
{
    if (LIST[0].ptr == NULL && LIST[0].size == 0)
    {
		static vmemory_t vm;
        LIST[0].ptr = vm.heap;
        LIST[0].size = HEAP_SIZE;
        IN_USE++;
		LOG();
    }

    entity_t* best = LIST;

    for (unsigned i = 0; i < IN_USE; i++)
    {
        if (LIST[i].size >= size && LIST[i].size < best->size)
        {
            best = &LIST[i];
        }
    }

    return best;
}

void* y_malloc(size_t size)
{
    assert((size+HEADER) <= HEAP_SIZE);

    // examples
    // size = 4
    size += HEADER; // after size = 8

    entity_t* ent = new_entity(size);

    // example with new entity
    u8* start = ent->ptr; // start = ptr to first byte of heap
    u8* user_ptr = start + HEADER; // user_ptr = ptr to first byte after header, where header = metadata

    *start = size; // write to header size of allocation

    ent->ptr += size; // offset ptr to next free byte
    ent->size -= size; // write size after allocation

    assert(ent->size >= 0);

	LOG();
    return user_ptr;
}

// fragmentation
// |XXXX------|
// |--XX------|
// OR
// |XXXXXX----|
// |XX--XX----|
// how can i fix it?
void y_free(void* ptr)
{
    u8* start = (u8*)ptr - HEADER; // get first byte of ptr

	//u16 PREV_USE = IN_USE-1;
	LIST[IN_USE].ptr = start;
	LIST[IN_USE].size = *start;
	IN_USE++;
	LOG();
}

// void test()
// {
//     long *l1, *l2, *l3, *l4 , *l5, *l6;
//     l1 = y_malloc(sizeof(long));
//     *l1 = 1;

//     l2 = y_malloc(sizeof(long));
//     *l2 = 2;

// 	   y_free(l1);

//     l3 = y_malloc(sizeof(long));
//     *l3 = 3;

//     l4 = y_malloc(sizeof(long));
//     *l4 = 4;

//     l5 = y_malloc(sizeof(long));
//     *l5 = 5;

//     l6 = y_malloc(sizeof(long));
//     *l6 = 6;

// 	   printf("%ld, %ld, %ld, %ld, %ld, %ld\n", *l1, *l2, *l3, *l4, *l5, *l6);
// }

// int main(int argc, char** argv)
// {
//     test();
//     return 0;
// }
