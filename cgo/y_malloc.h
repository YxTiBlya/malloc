#ifndef Y_MALLOC_H
#define Y_MALLOC_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <assert.h>
#include <stdbool.h>

#define u8 uint8_t
#define u16 uint16_t

#define HEAP_SIZE 128
#define HEADER_SIZE 4

typedef struct
{
    u8 heap[HEAP_SIZE];
} vmem_t;

typedef struct node
{
    struct node* next;
    struct node* prev;
    u8*          data;
    size_t       size;
    bool         is_used;
} node_t;

typedef struct
{
    u8*     data;
    node_t* head;
    size_t  used;
} alloc_t;

void* y_malloc(size_t size);
void* y_calloc(size_t count, size_t size);
void  y_free(void* ptr);

void* y_allocate(alloc_t* alloc, size_t count, size_t size);
void  y_deallocate(alloc_t* alloc, void* ptr);

alloc_t* y_create_alloc(size_t size);

void y_create_node(alloc_t* alloc, size_t size, node_t** out_node, node_t* last_node, bool freed);
void y_find_best_node(alloc_t* alloc, size_t size, node_t** out_node, node_t** out_last_node);
void y_divide_node(alloc_t* alloc, size_t size, node_t** out_node);

#endif