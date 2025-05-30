#include "y_malloc.h"

static vmem_t y_mem;
static alloc_t* y_alloc;

void* y_malloc(size_t size)
{
    assert((size+HEADER_SIZE) <= HEAP_SIZE);
    if (!y_alloc)
    {
        y_alloc = y_create_alloc(size);
    }
    
    return y_allocate(y_alloc, 1, size);
}

void* y_calloc(size_t count, size_t size)
{
    assert((size*count+HEADER_SIZE) <= HEAP_SIZE);
    
    if (!y_alloc)
    {
        y_alloc = y_create_alloc(count*size);
    }

    return y_allocate(y_alloc, count, size);
}

void y_free(void* ptr)
{
    assert(ptr);
    y_deallocate(y_alloc, ptr);
}

void* y_allocate(alloc_t* alloc, size_t count, size_t size)
{
    size_t total_size = size*count+HEADER_SIZE;
    assert(total_size <= HEAP_SIZE-y_alloc->used);

    node_t *node, *last_node;
    y_find_best_node(alloc, total_size, &node, &last_node);
    if (!node)
    {
        y_create_node(alloc, total_size, &node, last_node, false);
    }
    if (node->size > total_size)
    {
        y_divide_node(alloc, total_size, &node);
    }

    node->is_used = 1;
    u8* size_header = node->data;
    u8* count_header = size_header + HEADER_SIZE/2;
    u8* user_ptr = size_header + HEADER_SIZE;

    *size_header = size;
    *count_header = count;

    return user_ptr;
}

void y_deallocate(alloc_t* alloc, void* data_ptr)
{
    u8* ptr = (u8*)data_ptr - HEADER_SIZE;
    node_t* current = alloc->head;
    node_t *target_node, *freed_node, *prev, *last;

    while (current)
    {
        if (current->data == ptr)
        {
            target_node = current;
            break;
        }

        current = current->next;
    }

    current = target_node;
    if (current->is_used) { alloc->used -= current->size; }
    current->is_used = 0;

    while (current && !current->is_used)
    {
        last = current;
        current = current->next;
    }

    current = last;
    last = last->next;

    long freed_mem = 0;
    while (current && !current->is_used)
    {
        prev = current->prev;
        freed_mem += current->size;
        free(current);
        current = prev;
    }

    if (current) 
    {
        y_create_node(alloc, freed_mem, &freed_node, current, true);
        freed_node->next = last;
    }
    if (last)
    {
        if (current)
        {
            last->prev = freed_node;
        }
        else
        {
            y_create_node(alloc, freed_mem, NULL, NULL, true);
            alloc->head->next = last;
            last->prev = alloc->head;
        }
    }
}

alloc_t* y_create_alloc(size_t size)
{
    alloc_t* alloc = malloc(sizeof(alloc_t));
    memset(alloc, 0, sizeof(alloc_t));

    alloc->data = y_mem.heap;

    y_create_node(alloc, size+HEADER_SIZE, NULL, NULL, false);

    return alloc;
}

void y_create_node(alloc_t* alloc, size_t size, node_t** out_node, node_t* last_node, bool freed)
{
    node_t* new_node = malloc(sizeof(node_t));
    new_node->next = NULL;
    new_node->size = size;
    new_node->is_used = 0;

    if (last_node)
    {
        new_node->data = last_node->data + last_node->size;
        new_node->prev = last_node;
        last_node->next = new_node;
        *out_node = new_node;
    }
    else
    {
        new_node->data = alloc->data;
        new_node->prev = NULL;
        alloc->head = new_node;
    }

    if (!freed) { alloc->used += size; }
}

void y_find_best_node(alloc_t* alloc, size_t size, node_t** out_node, node_t** out_last_node)
{
    node_t* current = alloc->head;
    node_t* best = current;
    node_t* last;

    while (current)
    {
        if (current->is_used == 0 && current->size >= size && current->size < best->size)
        {
            best = current;
            if (best->size == size)
            {
                *out_node = best;
                return;
            }
        }

        last = current;
        current = current->next;
    }

    if (best->is_used == 0 && best->size >= size) 
    {
        *out_node = best;
        return;
    }

    *out_node = NULL;
    *out_last_node = last;
}

void y_divide_node(alloc_t* alloc, size_t size, node_t** out_node)
{
    node_t* new_node = *out_node;
    node_t* recidue_node = malloc(sizeof(node_t));

    recidue_node->data = new_node->data + size;
    recidue_node->prev = new_node;
    recidue_node->next = new_node->next;
    recidue_node->size = new_node->size - size;
    recidue_node->is_used = 0;
    
    new_node->next = recidue_node;
    new_node->size = size;
    alloc->used += size;
    
    *out_node = new_node;
}

void log_alloc()
{
    node_t* curr = y_alloc->head;

    int i = 1;
    while (curr)
    {
        printf("%d) ptr: [%p], value: [%d], size: [%zu bytes] used: [%d]\n", i, curr->data, *(curr->data+HEADER_SIZE), curr->size, curr->is_used);
        curr = curr->next;
        i++;
    }

    printf("total memory: %d bytes\nfree memory: %zu bytes\nused memory: %zu bytes\n", HEAP_SIZE, HEAP_SIZE-y_alloc->used, y_alloc->used);
}

// TODO: realloc
// int main(int argc, char** argv)
// {
//     long *ptr1, *ptr2, *ptr3;

//     ptr1 = y_malloc(sizeof(long));
//     *ptr1 = 1;

//     ptr2 = y_malloc(sizeof(long));
//     *ptr2 = 2;

//     ptr3 = y_malloc(sizeof(long));
//     *ptr3 = 3;

//     y_free(ptr2);
//     y_free(ptr1);

//     ptr2 = y_calloc(25, sizeof(long));
//     ptr2 -= sizeof(long);
//     for (int i = 0; i < 25; i++) {
//         ptr2 += sizeof(long);
//         *ptr2 = 4;
//     };

//     ptr1 = y_malloc(sizeof(long));
//     *ptr1 = 5;
    
//     log_alloc();

//     return 0;
// }