package ring

import "testing"

func TestCap(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		byteLimit     int
		byteThreshold int
		want          int
	}{
		"exact power of 2":  {16, 4, 4},
		"rounds up":         {20, 4, 8},
		"large values":      {1 << 30, 1 << 20, 1024},
		"non-power divisor": {7, 2, 4},
		"zero limit":        {0, 1, 1},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := capacity(tc.byteLimit, tc.byteThreshold)
			if got != tc.want {
				t.Errorf("cap(%d, %d) = %d, want %d", tc.byteLimit, tc.byteThreshold, got, tc.want)
			}
		})
	}
}

func TestRing(t *testing.T) {
	t.Parallel()
	type op struct {
		kind    string // "push" or "pop"
		val     int    // value to push, or expected value on pop
		wantLen int
	}
	tests := map[string]struct {
		limit     int
		threshold int
		ops       []op
	}{
		"single push/pop": {
			limit: 16, threshold: 4,
			ops: []op{
				{"push", 1, 1},
				{"pop", 1, 0},
			},
		},
		"FIFO order": {
			limit: 16, threshold: 4,
			ops: []op{
				{"push", 1, 1},
				{"push", 2, 2},
				{"push", 3, 3},
				{"pop", 1, 2},
				{"pop", 2, 1},
				{"pop", 3, 0},
			},
		},
		"overwrite on overflow": {
			limit: 16, threshold: 4,
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
		"wrap-around": {
			limit: 16, threshold: 4,
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
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			r := NewRing[int](tc.limit, tc.threshold)
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
