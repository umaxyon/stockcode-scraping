package lib

import "testing"

func TestUtils(t *testing.T) {
	t.Run("EquallyDivide", func(t *testing.T) {
		actual := EquallyDivide(40, 15)
		if len(actual) != 15 {
			t.Fatal("div size error")
		}
		var cnt = 0
		for r := range actual {
			cnt += actual[r]
		}
		if cnt != 40 {
			t.Fatal("total error")
		}
	})

	t.Run("ListChunk", func(t *testing.T) {
		data := []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p"}
		actual := ListChunk(data, 5)
		if len(actual) != 5 {
			t.Fatal("chunk size error")
		}
	})
}
