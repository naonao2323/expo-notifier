package ring

import "testing"

func TestCap(t *testing.T) {
	tests := []struct {
		byteLimit     int
		byteThreshold int
		want          int
	}{
		{16, 4, 4},
		{20, 4, 8},
		{1 << 30, 1 << 20, 1024},
		{7, 2, 4},
		{0, 1, 1},
	}
	for _, tc := range tests {
		t.Run("", func(t *testing.T) {
			got := Cap(tc.byteLimit, tc.byteThreshold)
			if got != tc.want {
				t.Errorf("Cap(%d, %d) = %d, want %d", tc.byteLimit, tc.byteThreshold, got, tc.want)
			}
		})
	}
}

func TestRing(t *testing.T) {
	type op struct {
		kind    string // "push" or "pop"
		val     int    // value to push, or expected value on pop
		wantLen int
	}
	tests := []struct {
		name string
		cap  int
		ops  []op
	}{
		{
			name: "single push/pop",
			cap:  4,
			ops: []op{
				{"push", 1, 1},
				{"pop", 1, 0},
			},
		},
		{
			name: "FIFO order",
			cap:  4,
			ops: []op{
				{"push", 1, 1},
				{"push", 2, 2},
				{"push", 3, 3},
				{"pop", 1, 2},
				{"pop", 2, 1},
				{"pop", 3, 0},
			},
		},
		{
			name: "overwrite on overflow",
			cap:  4,
			ops: []op{
				{"push", 1, 1},
				{"push", 2, 2},
				{"push", 3, 3},
				{"push", 4, 4},
				{"push", 5, 4}, // tail wraps to 0, overwrites value 1; len stays at capacity
				{"pop", 5, 3},  // head=0 now holds 5, not 1
				{"pop", 2, 2},
				{"pop", 3, 1},
				{"pop", 4, 0},
			},
		},
		{
			name: "wrap-around",
			cap:  4,
			ops: []op{
				{"push", 1, 1},
				{"push", 2, 2},
				{"push", 3, 3},
				{"push", 4, 4},
				{"pop", 1, 3},
				{"pop", 2, 2},
				{"push", 5, 3}, // tail wraps to index 0
				{"push", 6, 4}, // tail wraps to index 1
				{"pop", 3, 3},
				{"pop", 4, 2},
				{"pop", 5, 1},
				{"pop", 6, 0},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			r := New[int](tc.cap)
			for i, o := range tc.ops {
				switch o.kind {
				case "push":
					r.Push(o.val)
				case "pop":
					got := r.Pop()
					if got != o.val {
						t.Errorf("step %d: Pop() = %d, want %d", i, got, o.val)
					}
				}
				if r.Len() != o.wantLen {
					t.Errorf("step %d (%s %d): Len() = %d, want %d", i, o.kind, o.val, r.Len(), o.wantLen)
				}
			}
		})
	}
}
