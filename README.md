# 🧭 TinyCache Design Roadmap

> TinyCache is not just a cache implementation.
> It is an exploration of **concurrency, eviction strategies, and memory-constrained system design**.

---

## 🎯 Goals

* Build a high-performance in-memory cache
* Explore tradeoffs in:

  * Concurrency
  * Memory usage
  * Eviction strategies
* Keep the system simple, but **intentionally designed**

---

## 🧱 Phase 1: Core Cache (Already Implemented ✅)

### Features

* Basic commands are-

    ***PING, GET, SET, EXISTS, DEL,***
    
    ***INCR, DECR, LPUSH, RPUSH, LPOP, RPOP,***

    ***LRANGE, EXPIRE, TTL, PERSIST, FLUSHALL***
* In-memory storage
* RESP-based communication
* TTL support
* LRU / LFU eviction
* Thread-safe operations

### Known Limitations

* Global locking may limit scalability
* Memory usage is not strictly bounded
* Eviction accuracy vs performance not deeply explored
* No benchmarking yet

---

## ⚙️ Phase 2: Concurrency Model (High Priority 🚀)

### Problem

Global locks create contention under high concurrency.

### Approach

Introduce **sharded cache**:

* Split cache into N shards
* Route keys using hash(key) % N
* Each shard has its own lock

### Expected Outcome

* Reduced lock contention
* Improved parallelism

### Tradeoffs

* Slight memory overhead (multiple maps)
* Uneven key distribution possible

---

## 🧠 Phase 3: Eviction Strategy Deep Dive

### Problem

Basic LRU / LFU may not perform well under real workloads.

### Experiments

#### 1. Segmented LRU

* Separate hot and cold data
* Frequently accessed keys promoted

#### 2. Approximate LFU

* Use probabilistic counters instead of exact frequency

### Goals

* Compare:

  * Accuracy
  * Memory overhead
  * Performance impact

---

## ⏳ Phase 4: Expiration Strategy

### Problem

TTL handling can introduce latency spikes.

### Approaches

* Lazy expiration (on access)
* Background cleanup (goroutine)
* Optional: time-based structures (min-heap / time wheel)

### Tradeoffs

* Lazy → low CPU, stale data risk
* Active → more CPU, cleaner memory

---

## 🧬 Phase 5: Memory Management

### Problem

Cache size must be bounded.

### Approach

* Introduce max memory limit
* Track approximate memory usage
* Trigger eviction when limit reached

### Challenges

* Estimating object size in Go
* Balancing accuracy vs overhead

---

## ⚡ Phase 6: Benchmarking & Performance

### Metrics

* Throughput (ops/sec)
* Latency (avg / p95 / p99)
* Memory usage

### Tools

* Go benchmarking (`testing.B`)
* Custom load testing script

### Goal

Quantify tradeoffs instead of guessing.

---

## 🔐 Phase 7: Concurrency Guarantees

### Define Clearly

* What is thread-safe?
* What consistency guarantees exist?

### Example

* Per-shard locking
* Read-heavy optimization using RWMutex

---

## 🧪 Phase 8: Edge Cases & Failure Handling

### Explore

* Concurrent read/write conflicts
* Expired key access
* Memory limit reached

### Goal

Document behavior, not just implement features.

---

## 💾 Phase 9: Persistence (Optional Advanced)

### Options

* Append-only file (AOF)
* Snapshot-based persistence

### Tradeoffs

* Durability vs performance
* Write amplification

---

### 📥 Installation
If you'd like to install TinyCache on Linux without cloning the repository, use the following command to install both the server and client:
```bash
curl -fsSL https://raw.githubusercontent.com/nahK994/TinyCache/master/install.sh | bash
```
This will download the binaries, install them to the appropriate locations, and set up the server as a systemd service.


### 🧹 Uninstallation
To uninstall from linex TinyCache, simply run the following command:
```bash
curl -fsSL https://raw.githubusercontent.com/nahK994/TinyCache/master/uninstall.sh | bash
```
This will stop the service, remove the binaries, and clean up all installed files.

## 📌 Final Vision

TinyCache aims to answer:

> “What happens when a simple cache is pushed toward real-world constraints?”
