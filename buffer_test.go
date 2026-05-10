package exponotifier

import "testing"

func TestBufferRingCap(t *testing.T) {
	tests := map[string]struct {
		byteLimit      int
		byteThreshold  int
		countThreshold int
		minItemBytes   int
		want           int
	}{
		"byte capacity dominates":              {16, 4, 10, 128, 4},
		"byte capacity dominates, non-power-2": {20, 4, 10, 128, 5},
		"count capacity dominates":             {1 << 20, 1 << 20, 10, 128, 819},
		"default values, count capacity wins":  {1 << 30, 1 << 20, 10, 128, 838860},
		"zero byteLimit returns minimum 1":     {0, 1, 10, 128, 1},
		"zero byteThreshold returns minimum 1": {16, 0, 10, 128, 1},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := bufferRingCap(tc.byteLimit, tc.byteThreshold, tc.countThreshold, tc.minItemBytes)
			if got != tc.want {
				t.Errorf("bufferRingCap(%d, %d, %d, %d) = %d, want %d",
					tc.byteLimit, tc.byteThreshold, tc.countThreshold, tc.minItemBytes, got, tc.want)
			}
		})
	}
}
