package pokecache

import (
	"fmt"
	"testing"
	"time"
)

func TestAddGet(t *testing.T) {
	const interval = 5 * time.Second
	cases := []struct {
		key string
		val []byte
	}{
		{
			key: "https://example.com",
			val: []byte("testdata"),
		},
		{
			key: "https://example.com/path",
			val: []byte("moretestdata"),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Test case %v", i), func(t *testing.T) {
			cache := NewCache(interval)
			cache.Add(c.key, c.val)
			val, ok := cache.Get(c.key)
			if !ok {
				t.Errorf("expected to find key")
				return
			}
			if string(val) != string(c.val) {
				t.Errorf("expected to find value")
				return
			}
		})
	}
}

func TestReapLoop(t *testing.T) {
	const baseTime = 5 * time.Millisecond
	const waitTime = baseTime + 5*time.Millisecond
	cache := NewCache(baseTime)
	cache.Add("https://example.com", []byte("testdata"))

	_, ok := cache.Get("https://example.com")
	if !ok {
		t.Errorf("expected to find key")
		return
	}

	time.Sleep(waitTime)

	_, ok = cache.Get("https://example.com")
	if ok {
		t.Errorf("expected to not find key")
		return
	}
}

func TestNewCache(t *testing.T) {
	const interval = 1 * time.Second
	cache := NewCache(interval)
	
	if cache == nil {
		t.Errorf("expected cache to be created")
		return
	}
	
	if cache.entries == nil {
		t.Errorf("expected entries map to be initialized")
		return
	}
}

func TestAdd(t *testing.T) {
	cache := NewCache(5 * time.Second)
	key := "test-key"
	val := []byte("test-value")
	
	cache.Add(key, val)
	
	// Verify the entry was added
	retrievedVal, ok := cache.Get(key)
	if !ok {
		t.Errorf("expected to find key after adding")
		return
	}
	
	if string(retrievedVal) != string(val) {
		t.Errorf("expected retrieved value to match added value")
		return
	}
}

func TestGetNonExistent(t *testing.T) {
	cache := NewCache(5 * time.Second)
	key := "non-existent-key"
	
	val, ok := cache.Get(key)
	if ok {
		t.Errorf("expected to not find non-existent key")
		return
	}
	
	if val != nil {
		t.Errorf("expected nil value for non-existent key")
		return
	}
}

func TestConcurrentAccess(t *testing.T) {
	cache := NewCache(5 * time.Second)
	
	// Test concurrent writes
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			key := fmt.Sprintf("key-%d", i)
			val := []byte(fmt.Sprintf("value-%d", i))
			cache.Add(key, val)
			done <- true
		}(i)
	}
	
	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// Verify all entries were added
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("key-%d", i)
		expectedVal := fmt.Sprintf("value-%d", i)
		
		val, ok := cache.Get(key)
		if !ok {
			t.Errorf("expected to find key %s", key)
			return
		}
		
		if string(val) != expectedVal {
			t.Errorf("expected value %s, got %s", expectedVal, string(val))
			return
		}
	}
}

func TestCacheExpiration(t *testing.T) {
	const shortInterval = 10 * time.Millisecond
	cache := NewCache(shortInterval)
	
	// Add an entry
	key := "test-key"
	val := []byte("test-value")
	cache.Add(key, val)
	
	// Verify it exists
	_, ok := cache.Get(key)
	if !ok {
		t.Errorf("expected to find key immediately after adding")
		return
	}
	
	// Wait for expiration
	time.Sleep(shortInterval + 5*time.Millisecond)
	
	// Verify it's gone
	_, ok = cache.Get(key)
	if ok {
		t.Errorf("expected key to be expired and removed")
		return
	}
}

func TestMultipleEntriesExpiration(t *testing.T) {
	const shortInterval = 10 * time.Millisecond
	cache := NewCache(shortInterval)
	
	// Add multiple entries
	keys := []string{"key1", "key2", "key3"}
	vals := [][]byte{[]byte("val1"), []byte("val2"), []byte("val3")}
	
	for i, key := range keys {
		cache.Add(key, vals[i])
	}
	
	// Verify all exist
	for _, key := range keys {
		_, ok := cache.Get(key)
		if !ok {
			t.Errorf("expected to find key %s", key)
			return
		}
	}
	
	// Wait for expiration
	time.Sleep(shortInterval + 5*time.Millisecond)
	
	// Verify all are gone
	for _, key := range keys {
		_, ok := cache.Get(key)
		if ok {
			t.Errorf("expected key %s to be expired and removed", key)
			return
		}
	}
}

func TestOverwriteEntry(t *testing.T) {
	cache := NewCache(5 * time.Second)
	key := "test-key"
	originalVal := []byte("original-value")
	newVal := []byte("new-value")
	
	// Add original entry
	cache.Add(key, originalVal)
	
	// Verify original value
	val, ok := cache.Get(key)
	if !ok {
		t.Errorf("expected to find key after adding original value")
		return
	}
	if string(val) != string(originalVal) {
		t.Errorf("expected original value")
		return
	}
	
	// Overwrite with new value
	cache.Add(key, newVal)
	
	// Verify new value
	val, ok = cache.Get(key)
	if !ok {
		t.Errorf("expected to find key after overwriting")
		return
	}
	if string(val) != string(newVal) {
		t.Errorf("expected new value after overwriting")
		return
	}
}
