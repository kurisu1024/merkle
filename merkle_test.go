package merkle_test

import (
	"fmt"
	"testing"

	"github.com/kurisu1024/merkle"
)

func TestNewTree(t *testing.T) {
	tests := []struct {
		name    string
		data    []string
		wantErr bool
	}{
		{
			name:    "empty data",
			data:    []string{},
			wantErr: true,
		},
		{
			name:    "single element",
			data:    []string{"a"},
			wantErr: false,
		},
		{
			name:    "two elements",
			data:    []string{"a", "b"},
			wantErr: false,
		},
		{
			name:    "three elements",
			data:    []string{"a", "b", "c"},
			wantErr: false,
		},
		{
			name:    "four elements (power of 2)",
			data:    []string{"a", "b", "c", "d"},
			wantErr: false,
		},
		{
			name:    "five elements",
			data:    []string{"a", "b", "c", "d", "e"},
			wantErr: false,
		},
		{
			name:    "eight elements (power of 2)",
			data:    []string{"a", "b", "c", "d", "e", "f", "g", "h"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree, err := merkle.NewTree(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewTree() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if tree == nil {
					t.Error("NewTree() returned nil tree")
					return
				}
				if tree.GetRoot() == "" {
					t.Error("NewTree() root hash is empty")
				}
			}
		})
	}
}

func TestNewTreeBytes(t *testing.T) {
	data := [][]byte{
		[]byte("test1"),
		[]byte("test2"),
		[]byte("test3"),
	}

	tree, err := merkle.NewTree(data)
	if err != nil {
		t.Fatalf("NewTree() error = %v", err)
	}

	if tree.GetRoot() == "" {
		t.Error("NewTree() root hash is empty")
	}
}

func TestGetRoot(t *testing.T) {
	tests := []struct {
		name string
		data []string
	}{
		{
			name: "single element",
			data: []string{"a"},
		},
		{
			name: "multiple elements",
			data: []string{"a", "b", "c", "d"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree, err := merkle.NewTree(tt.data)
			if err != nil {
				t.Fatalf("NewTree() error = %v", err)
			}

			root := tree.GetRoot()
			if root == "" {
				t.Error("GetRoot() returned empty string")
			}

			// Root should be deterministic
			tree2, _ := merkle.NewTree(tt.data)
			root2 := tree2.GetRoot()
			if root != root2 {
				t.Error("GetRoot() not deterministic")
			}
		})
	}
}

func TestGetRootDifferentData(t *testing.T) {
	tree1, _ := merkle.NewTree([]string{"a", "b", "c"})
	tree2, _ := merkle.NewTree([]string{"a", "b", "d"})

	if tree1.GetRoot() == tree2.GetRoot() {
		t.Error("Different data produced same root hash")
	}
}

func TestGetProof(t *testing.T) {
	tests := []struct {
		name     string
		data     []string
		index    int
		wantErr  bool
		proofLen int
	}{
		{
			name:     "single element, index 0",
			data:     []string{"a"},
			index:    0,
			wantErr:  false,
			proofLen: 0,
		},
		{
			name:     "two elements, index 0",
			data:     []string{"a", "b"},
			index:    0,
			wantErr:  false,
			proofLen: 1,
		},
		{
			name:     "two elements, index 1",
			data:     []string{"a", "b"},
			index:    1,
			wantErr:  false,
			proofLen: 1,
		},
		{
			name:     "four elements, index 0",
			data:     []string{"a", "b", "c", "d"},
			index:    0,
			wantErr:  false,
			proofLen: 2,
		},
		{
			name:     "four elements, index 3",
			data:     []string{"a", "b", "c", "d"},
			index:    3,
			wantErr:  false,
			proofLen: 2,
		},
		{
			name:     "three elements, index 2",
			data:     []string{"a", "b", "c"},
			index:    2,
			wantErr:  false,
			proofLen: 2,
		},
		{
			name:    "index out of range (negative)",
			data:    []string{"a", "b"},
			index:   -1,
			wantErr: true,
		},
		{
			name:    "index out of range (too large)",
			data:    []string{"a", "b"},
			index:   2,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree, err := merkle.NewTree(tt.data)
			if err != nil {
				t.Fatalf("NewTree() error = %v", err)
			}

			proof, err := tree.GetProof(tt.index)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProof() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(proof) != tt.proofLen {
					t.Errorf("GetProof() proof length = %v, want %v", len(proof), tt.proofLen)
				}
			}
		})
	}
}

func TestVerifyProof(t *testing.T) {
	tests := []struct {
		name  string
		data  []string
		index int
	}{
		{
			name:  "single element",
			data:  []string{"a"},
			index: 0,
		},
		{
			name:  "two elements, first",
			data:  []string{"a", "b"},
			index: 0,
		},
		{
			name:  "two elements, second",
			data:  []string{"a", "b"},
			index: 1,
		},
		{
			name:  "four elements, first",
			data:  []string{"a", "b", "c", "d"},
			index: 0,
		},
		{
			name:  "four elements, middle",
			data:  []string{"a", "b", "c", "d"},
			index: 2,
		},
		{
			name:  "four elements, last",
			data:  []string{"a", "b", "c", "d"},
			index: 3,
		},
		{
			name:  "odd number of elements",
			data:  []string{"a", "b", "c", "d", "e"},
			index: 4,
		},
		{
			name:  "seven elements",
			data:  []string{"a", "b", "c", "d", "e", "f", "g"},
			index: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree, err := merkle.NewTree(tt.data)
			if err != nil {
				t.Fatalf("NewTree() error = %v", err)
			}

			proof, err := tree.GetProof(tt.index)
			if err != nil {
				t.Fatalf("GetProof() error = %v", err)
			}

			rootHash := tree.GetRoot()
			valid := merkle.VerifyProof(tt.data[tt.index], proof, rootHash, tt.index)

			if !valid {
				t.Errorf("VerifyProof() = false, want true for valid proof")
			}
		})
	}
}

func TestVerifyProofInvalid(t *testing.T) {
	data := []string{"a", "b", "c", "d"}
	tree, _ := merkle.NewTree(data)

	tests := []struct {
		name     string
		data     string
		proof    []string
		rootHash string
		index    int
	}{
		{
			name:     "wrong data",
			data:     "wrong",
			proof:    func() []string { p, _ := tree.GetProof(0); return p }(),
			rootHash: tree.GetRoot(),
			index:    0,
		},
		{
			name:     "wrong root hash",
			data:     data[0],
			proof:    func() []string { p, _ := tree.GetProof(0); return p }(),
			rootHash: "wrong_hash",
			index:    0,
		},
		{
			name:     "wrong index",
			data:     data[0],
			proof:    func() []string { p, _ := tree.GetProof(0); return p }(),
			rootHash: tree.GetRoot(),
			index:    1,
		},
		{
			name:     "tampered proof",
			data:     data[0],
			proof:    []string{"tampered_hash"},
			rootHash: tree.GetRoot(),
			index:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid := merkle.VerifyProof(tt.data, tt.proof, tt.rootHash, tt.index)
			if valid {
				t.Errorf("VerifyProof() = true, want false for invalid proof")
			}
		})
	}
}

func TestVerifyProofAllLeaves(t *testing.T) {
	data := []string{"leaf0", "leaf1", "leaf2", "leaf3", "leaf4", "leaf5", "leaf6", "leaf7"}
	tree, err := merkle.NewTree(data)
	if err != nil {
		t.Fatalf("NewTree() error = %v", err)
	}

	rootHash := tree.GetRoot()

	// Verify proof for every leaf
	for i, d := range data {
		proof, err := tree.GetProof(i)
		if err != nil {
			t.Fatalf("GetProof(%d) error = %v", i, err)
		}

		if !merkle.VerifyProof(d, proof, rootHash, i) {
			t.Errorf("VerifyProof() failed for leaf %d", i)
		}
	}
}

func TestPrint(t *testing.T) {
	// Test that Print doesn't panic
	tests := []struct {
		name string
		data []string
	}{
		{
			name: "single element",
			data: []string{"a"},
		},
		{
			name: "multiple elements",
			data: []string{"a", "b", "c", "d"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree, err := merkle.NewTree(tt.data)
			if err != nil {
				t.Fatalf("NewTree() error = %v", err)
			}

			// Just verify it doesn't panic
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Print() panicked: %v", r)
				}
			}()

			tree.Print()
		})
	}
}

func BenchmarkNewTree(b *testing.B) {
	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		data := make([]string, size)
		for i := range data {
			data[i] = string(rune('a' + i%26))
		}

		b.ResetTimer()
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				merkle.NewTree(data)
			}
		})
	}
}

func BenchmarkGetProof(b *testing.B) {
	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		data := make([]string, size)
		for i := range data {
			data[i] = string(rune('a' + i%26))
		}

		tree, _ := merkle.NewTree(data)

		b.ResetTimer()
		b.Run(fmt.Sprintf("size-%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tree.GetProof(i % size)
			}
		})
	}
}

func BenchmarkVerifyProof(b *testing.B) {
	data := make([]string, 1000)
	for i := range data {
		data[i] = string(rune('a' + i%26))
	}

	tree, _ := merkle.NewTree(data)
	proof, _ := tree.GetProof(500)
	rootHash := tree.GetRoot()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		merkle.VerifyProof(data[500], proof, rootHash, 500)
	}
}
