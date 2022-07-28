title: Algorithms and Data Structures Series: Hash Maps
author: Geison Biazus
description:
image_path: /static/image/logo-small.png
time: 2022-07-28 09:00
--

This post is part of the algorithms and data structures series, a series of posts where I present the most common data structures and algorithms used in software engineering. In this post, I explain the basics of Hash maps, also known as hash tables or dictionaries.

The Hash map stores data in a key-value pair. Similar to arrays, accessing a piece of data from a position has a constant complexity, but the positions on hash maps behave differently. Instead of having sequential numeric indexes, the Hash map allows any type of data to be used as its key. This makes the Hash map one of best and most used data structures for optimizing algorithms.

## Contents

## Behavior

The main operations implemented by a hash map are setting a value into a key, retrieving the value from a key, and removing a value form a key. Additionaly, it can provide other operations like returning all keys from the hash map that can be used to traversing it.

Another peculiarity from hash maps comparing to arrays, is that it does not have order. Traversing a map will not bring the values in the same order they were added. There are some implementations that keep the order of the added keys, but with the cost of additional space.

## Hashing function

This data structure is called hash map due to how the keys and values are resolved and stored internally. It makes use of a hash funcion and an array of "buckets". The hash function receives the provided key and transforms it into a number in a way that every time the same key is provided, the same number is returned. This number is then used as the index of the internal array of "buckets". A bucket is another data structure that stores the key-value pair.

The process of inserting a value into a hash map works as follows:

1. The hash map receives the key and value to be stored
1. The key is given to the hash function and the bucket index is resolved
1. The key-value pair is inserted in the bucket of the previous index

Retrieving the value from a hash map follows a similar proccess:

1. The hash map receives a key
1. The key is given to the hash function and the bucket index is resolved
1. The value of the bucket is returned based on the previous index

Here is a visual representation of how values are stored:

![Hash map representation](/static/image/hash-map.png)

<center style="margin-top: -30px; margin-bottom: 30px;"><small>Source: <a href="https://en.wikipedia.org/wiki/Hash_table" target="_blank">https://en.wikipedia.org/wiki/Hash_table</a></small></center>

## Collision

Sometimes, two different keys can return the same result after been passed to the hash function. We call this collision. To solve this problem, the buckets don't contain the values directly, but lists of values of the same hashing result. This list can be implemented using other data structures like dynamic arrays or linked lists. There are other collision resolution techniques and you can check some of them on [Hash table - Wikipedia](https://en.wikipedia.org/wiki/Hash_table).

## Implementation

Now let's start with our hash map implementation. The code examples are written using the Go progamming language, but they are simple enough that cna be applied to any language.

### Insert

### Lookup

### Delete

### Traverse

## Complexity table

## Choosing hash maps over other data structures

### When to use

### When to avoid

## Final thoughts

## Full hash map implementation

## Sources

- [Master the Coding Interview: Data Structures + Algorithms](https://www.udemy.com/course/master-the-coding-interview-data-structures-algorithms/)
- [Hsh table - Wikipedia](https://en.wikipedia.org/wiki/Hash_table)
