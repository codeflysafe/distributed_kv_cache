package lru

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	value := "1234"
	getTests := [3]string{"123", "32", "3"}
	for _, key := range getTests {
		lru := New(0)
		lru.Add(key, value)
		val, ok := lru.Get(key)
		if !ok {
			t.Fatalf("%s : cache hit = %v; want: %v", key, ok, !ok)
		} else if val != value {
			t.Fatalf("%s expected get to return %s but got %v", key, value, val)
		}
	}
}

func TestRemove(t *testing.T) {
	lru := New(0)
	key := "mykey"
	lru.Add(key, 1234)
	// print(lru.cache["mykey"].Value)
	if val, ok := lru.Get(key); !ok {
		t.Fatal("TestRemove returned no match")
	} else if val != 1234 {
		t.Fatalf("TestRemove failed.  Expected %d, got %v", 1234, val)
	}

	lru.Remove(key)
	if _, ok := lru.Get(key); ok {
		t.Fatal("TestRemove returned a removed entry")
	}
}

func TestEvict(t *testing.T) {
	evictedKeys := make([]string, 0)
	onEvictedFun := func(key string, value interface{}) {
		evictedKeys = append(evictedKeys, key)
	}

	lru := New(20)
	lru.OnEvicted = onEvictedFun
	for i := 0; i < 22; i++ {
		lru.Add(fmt.Sprintf("myKey%d", i), 1234)
	}

	if len(evictedKeys) != 2 {
		t.Fatalf("got %d evicted keys; want 2", len(evictedKeys))
	}
	if evictedKeys[0] != "myKey0" {
		t.Fatalf("got %v in first evicted key; want %s", evictedKeys[0], "myKey0")
	}
	if evictedKeys[1] != "myKey1" {
		t.Fatalf("got %v in second evicted key; want %s", evictedKeys[1], "myKey1")
	}
}
