package main

import "testing"

func Test_writeCacheFile(t *testing.T) {
	var cache = []byte{0, 1, 2}
	t.Run(t.Name(), func(t *testing.T) {
		for i := 0; i < 10; i++ {
			t.Parallel()
			writeCacheFile(cache)
		}
	})
}
