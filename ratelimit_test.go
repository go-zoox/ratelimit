package ratelimit

import (
	"fmt"
	"testing"
	"time"
)

func TestRateLimit(t *testing.T) {
	id := "127.0.0.1"
	r := NewMemory("web", 2*time.Second, 3)

	if remaining := r.Remaining(id); remaining != 3 {
		t.Fatal("remaining should be 3, but got", remaining)
	}
	if isExceeded := r.IsExceeded(id); isExceeded {
		t.Fatal("isExceeded should be false, but got", isExceeded)
	}

	if err := r.Inc(id); err != nil {
		t.Fatal(err)
	}
	if remaining := r.Remaining(id); remaining != 2 {
		t.Fatal("remaining should be 2, but got", remaining)
	}
	if isExceeded := r.IsExceeded(id); isExceeded {
		t.Fatal("isExceeded should be false, but got", isExceeded)
	}
	fmt.Println(r.Status(id))

	if err := r.Inc(id); err != nil {
		t.Fatal(err)
	}
	if remaining := r.Remaining(id); remaining != 1 {
		t.Fatal("remaining should be 1, but got", remaining)
	}
	if isExceeded := r.IsExceeded(id); isExceeded {
		t.Fatal("isExceeded should be false, but got", isExceeded)
	}
	fmt.Println(r.Status(id))

	if err := r.Inc(id); err != nil {
		t.Fatal(err)
	}
	if remaining := r.Remaining(id); remaining != 0 {
		t.Fatal("remaining should be 0, but got", remaining)
	}
	if isExceeded := r.IsExceeded(id); isExceeded {
		t.Fatal("isExceeded should be false, but got", isExceeded)
	}
	fmt.Println(r.Status(id))

	if err := r.Inc(id); err != nil {
		t.Fatal(err)
	}
	if remaining := r.Remaining(id); remaining != -1 {
		t.Fatal("remaining should be 0, but got", remaining)
	}
	if isExceeded := r.IsExceeded(id); !isExceeded {
		t.Fatal("isExceeded should be true, but got", isExceeded)
	}
	fmt.Println(r.Status(id))

	time.Sleep(1 * time.Second)
	if remaining := r.Remaining(id); remaining != -1 {
		t.Fatal("remaining should be -1, but got", remaining)
	}

	time.Sleep(1100 * time.Millisecond)
	if remaining := r.Remaining(id); remaining != 3 {
		t.Fatal("remaining should be 3, but got", remaining)
	}
	if isExceeded := r.IsExceeded(id); isExceeded {
		t.Fatal("isExceeded should be false, but got", isExceeded)
	}
}
