#ifndef Y_MALLOC_H
#define Y_MALLOC_H

#include <stdio.h>
#include <string.h>
#include <stdint.h>
#include <assert.h>

#define u8 uint8_t
#define u16 uint16_t
#define HEAP_SIZE 128
#define ENT_COUNT 40
#define HEADER 4

// memory emulation
typedef struct vmemory
{
    u8 heap[HEAP_SIZE];
} vmemory_t;

// emulation lists of os
typedef struct entity
{
    u8* ptr;
    int size;
} entity_t;

entity_t* new_entity(size_t size);
void* y_malloc(size_t size);
void y_free(void* ptr);

#endif