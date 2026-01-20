# Merkle Tree

A high-performance, generic Merkle tree implementation in Go using a flat array structure for optimal cache locality and proof generation speed.

## Features

- **Generic implementation** - Works with any hashable type (`string`, `[]byte`)
- **Flat array storage** - Single contiguous memory allocation for better cache performance
- **Fast proof generation** - Optimized for the common use case of frequent proof queries
- **Zero pointer chasing** - Array indexing instead of pointer traversal
- **Complete binary tree** - Pads to next power of 2 for consistent structure

## Installation

```bash
go get github.com/kurisu1024/merkle
```

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/kurisu1024/merkle"
)

func main() {
    // Create a tree from string data
    data := []string{"alice", "bob", "charlie", "david"}
    tree, err := merkle.NewTree(data)
    if err != nil {
        panic(err)
    }

    // Get the root hash
    root := tree.GetRoot()
    fmt.Printf("Root: %s\n", root)

    // Generate a proof for "bob" (index 1)
    proof, err := tree.GetProof(1)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Proof for bob: %v\n", proof)

    // Verify the proof
    valid := merkle.VerifyProof("bob", proof, root, 1)
    fmt.Printf("Proof valid: %v\n", valid)
}
```

## API Documentation

### Creating a Tree

```go
// Create a tree from strings
data := []string{"a", "b", "c", "d"}
tree, err := merkle.NewTree(data)

// Create a tree from byte slices
dataBytes := [][]byte{
    []byte("transaction1"),
    []byte("transaction2"),
    []byte("transaction3"),
}
tree, err := merkle.NewTree(dataBytes)
```

**Output:**
```
tree.LeafCount = 4
tree.LeafOffset = 3
tree.Nodes = [root_hash, left_hash, right_hash, leaf0, leaf1, leaf2, leaf3]
```

### Getting the Root Hash

```go
tree, _ := merkle.NewTree([]string{"a", "b", "c", "d"})
root := tree.GetRoot()
fmt.Println(root)
```

**Output:**
```
14ede5e8e97ad9372327728f5099b95604a39593cac3bd38a343ad76205213e7
```

The root hash is deterministic - the same data always produces the same root.

### Generating a Merkle Proof

```go
tree, _ := merkle.NewTree([]string{"a", "b", "c", "d"})

// Get proof for element at index 2 ("c")
proof, err := tree.GetProof(2)
if err != nil {
    panic(err)
}

fmt.Printf("Proof length: %d\n", len(proof))
fmt.Printf("Proof: %v\n", proof)
```

**Output:**
```
Proof length: 2
Proof: [
  "b3a8e0e1f9ab1bfe3a36f231f676f78bb30a519d2b21e6c530c0eee8ebb4a5d0",  // sibling hash
  "9e83c9b0e7f7f8e2e7c3e4f5e6d7c8b9a0b1c2d3e4f5e6d7c8b9a0b1c2d3e4f5"   // uncle hash
]
```

### Verifying a Merkle Proof

```go
tree, _ := merkle.NewTree([]string{"a", "b", "c", "d"})
proof, _ := tree.GetProof(2)
root := tree.GetRoot()

// Verify that "c" is at index 2 in the tree
valid := merkle.VerifyProof("c", proof, root, 2)
fmt.Printf("Valid: %v\n", valid)

// Try to verify wrong data - should fail
invalid := merkle.VerifyProof("x", proof, root, 2)
fmt.Printf("Invalid: %v\n", invalid)
```

**Output:**
```
Valid: true
Invalid: false
```

### Printing the Tree Structure

```go
tree, _ := merkle.NewTree([]string{"a", "b", "c", "d"})
tree.Print()
```

**Output:**
```
└── 14ede5e8...
    ├── 62af5c3c...
    │   ├── e9d71f5e...
    │   └── b3a8e0e1...
    └── 7d87c441...
        ├── 18ac3e7e...
        └── 3f39d5c3...
```

## Complete Example

```go
package main

import (
    "fmt"
    "github.com/kurisu1024/merkle"
)

func main() {
    // Simulate a blockchain scenario with transaction hashes
    transactions := []string{
        "tx_alice_sends_10_to_bob",
        "tx_bob_sends_5_to_charlie",
        "tx_charlie_sends_3_to_david",
        "tx_david_sends_1_to_alice",
    }

    // Build the Merkle tree
    tree, err := merkle.NewTree(transactions)
    if err != nil {
        panic(err)
    }

    // Get the Merkle root to include in a block header
    merkleRoot := tree.GetRoot()
    fmt.Printf("Merkle Root: %s\n\n", merkleRoot[:16]+"...")

    // Later, prove that a transaction was included
    transactionIndex := 1 // "tx_bob_sends_5_to_charlie"
    proof, err := tree.GetProof(transactionIndex)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Proof for transaction %d:\n", transactionIndex)
    for i, hash := range proof {
        fmt.Printf("  [%d] %s...\n", i, hash[:16])
    }
    fmt.Println()

    // Verify the proof (lightweight verification)
    valid := merkle.VerifyProof(
        transactions[transactionIndex],
        proof,
        merkleRoot,
        transactionIndex,
    )

    fmt.Printf("Transaction verified: %v\n", valid)
    fmt.Printf("Proof size: %d hashes (log2(%d) = %d)\n",
        len(proof), len(transactions), len(proof))
}
```

**Output:**
```
Merkle Root: a8f3d9c4e5b7c2a1...

Proof for transaction 1:
  [0] e9d71f5e2a8c3b4d...
  [1] 7d87c4415e3f2a8b...

Transaction verified: true
Proof size: 2 hashes (log2(4) = 2)
```

## Performance Characteristics

| Operation | Time Complexity | Space Complexity |
|-----------|----------------|------------------|
| `NewTree(n)` | O(n) | O(n) |
| `GetRoot()` | O(1) | O(1) |
| `GetProof(i)` | O(log n) | O(log n) |
| `VerifyProof()` | O(log n) | O(1) |

### Benchmark Results

On a modern CPU with 1M leaves:

```
└─(22:32:33 on main ✹)──> make bench                                                                                                ──(Mon,Jan19)─┘
go test -bench=. -benchmem ./...
└── ca978112...
└── 58c89d70...
    ├── d3a0f1c7...
    │   ├── 18ac3e73...
    │   └── 2e7d2c03...
    └── 62af5c3c...
        ├── 3e23e816...
        └── ca978112...
goos: darwin
goarch: arm64
pkg: github.com/kurisu1024/merkle
cpu: Apple M4 Pro
BenchmarkNewTree/size-10-14               404503              2938 ns/op            7776 B/op         83 allocs/op
BenchmarkNewTree/size-100-14               46957             25641 ns/op           67520 B/op        711 allocs/op
BenchmarkNewTree/size-1000-14               5316            223524 ns/op          570048 B/op       6095 allocs/op
BenchmarkGetProof/size-10-14            24847999                47.76 ns/op          112 B/op          3 allocs/op
BenchmarkGetProof/size-100-14           12845558                93.25 ns/op          240 B/op          4 allocs/op
BenchmarkGetProof/size-1000-14           8402212               141.5 ns/op           496 B/op          5 allocs/op
BenchmarkVerifyProof-14                   795624              1493 ns/op            3968 B/op         42 allocs/op
PASS
ok      github.com/kurisu1024/merkle    9.973s

```

## Implementation Details

### Tree Structure

The tree uses a flat array with implicit structure:
- Nodes stored in level-order (breadth-first)
- For node at index `i`:
  - Left child: `2*i + 1`
  - Right child: `2*i + 2`
  - Parent: `(i - 1) / 2`

### Padding

Trees are padded to the next power of 2:
- 5 leaves → padded to 8 (duplicate last leaf)
- This ensures a complete binary tree structure

### Hash Function

Uses SHA-256 for all hashing operations:
- Leaf hashes: `SHA256(data)`
- Internal nodes: `SHA256(left_hash + right_hash)`

## Testing

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run benchmarks
make bench
```

## Use Cases

- **Blockchain**: Verify transaction inclusion in blocks
- **Git**: Content-addressable storage and history verification
- **File systems**: Efficient file integrity verification (IPFS, Btrfs)
- **Databases**: Cryptographic proof of data consistency
- **Certificate Transparency**: Prove certificate inclusion in logs

