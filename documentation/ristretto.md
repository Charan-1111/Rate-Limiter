# Why Ristretto? A Deep Dive into Our In-Memory Caching Strategy

In the earlier versions of `rateLimiter`, caching policies locally was done using standard Go primitives: a `map[string]*PolicySchema` protected by a `sync.RWMutex` (or `sync.Mutex`). While simple and effective for small-scale applications, building a high-throughput, low-latency microservice like a rate limiter quickly exposed the limitations of traditional map implementations.

To solve these architectural bottlenecks, we migrated our in-memory policy cache to [**Dgraph's Ristretto**](https://github.com/dgraph-io/ristretto).

This document outlines the problems with Native Maps and why Ristretto is the superior choice for our use case.

---

## 1. Concurrency and Lock Contention

### The Problem with `sync.RWMutex`
A native `map` in Go is not thread-safe. To prevent race conditions and panics during concurrent reads and writes, it must be protected by a Mutex. While a `sync.RWMutex` allows multiple simultaneous readers, any write operation (like adding a new policy to the cache after a database fetch) completely blocks **all** readers. Under high load, thousands of goroutines end up waiting on a single lock, severely degrading throughput.

### The Ristretto Solution
Ristretto is built from the ground up to be heavily concurrent. It completely eliminates the need for global locks. 
- It uses lock-free, asynchronous ring buffers to track reads and writes. 
- Writes are deferred to a background goroutine, meaning your HTTP request handler is never blocked waiting to insert something into the cache. 

## 2. Memory Bounds and Eviction Policies

### The Problem with Native Maps
By default, a Go `map` will grow indefinitely. If you temporarily cache rate limit policies for millions of unique users, those entries will stay in memory forever unless you manually delete them. Implementing an LFU (Least Frequently Used) or LRU (Least Recently Used) eviction policy on top of a map is extremely complex and introduces massive computational overhead.

### The Ristretto Solution
Ristretto integrates a highly optimized admission policy known as **TinyLFU**.
- We can easily set boundaries like `MaxCost` (maximum memory allocation) or `NumCounters`.
- Once the cache reaches its limit, Ristretto intelligently evicts the least valuable/frequently used items to make room for new ones. Memory usage remains perfectly stable regardless of the traffic volume.

## 3. Time-To-Live (TTL)

### The Problem with Native Maps
If a rate-limiting policy drastically changes in the central PostgreSQL database, the microservice nodes need a way to realize their local cache is stale. With a standard map, you must either manually build background polling loops to invalidate keys, or wrap every map value in an expiry struct and check it on every read. Both add clutter and performance penalties to the codebase.

### The Ristretto Solution
Ristretto natively supports TTL (Time-To-Live). 
- We simply call `c.data.SetWithTTL(key, policy, 1, constants.PolicyCacheDuration)`.
- Ristretto automatically purges the item from memory when the time expires. The next request naturally results in a cache-miss, smoothly falling back to the database to fetch the freshest policy data.

---

## Summary

Migrating to Ristretto elevates our application from a "proof of concept" to a production-ready system. 

By offloading concurrency management, admission policies (TinyLFU), and eviction loops (TTL) to a battlefield-tested library, the Rate Limiter microservice achieves higher throughput, tighter bounded memory usage, and infinitely cleaner business logic.
