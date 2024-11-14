package single_cache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

var db = map[string]string{
	"Tom":  "123",
	"Jcak": "456",
	"Lucy": "789",
}

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Fatalf("callback failed")
	}
}

func TestGroup_Get(t *testing.T) {
	localCount := make(map[string]int, len(db))
	g := NewGroup("scores", GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				if _, ok := localCount[key]; !ok {
					localCount[key] = 0
				}
				localCount[key]++
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}), 2<<10)
	for k, v := range db {
		if view, err := g.Get(k); err != nil || view.String() != v {
			t.Fatalf("Get cache failed")
		}
		if _, err := g.Get(k); err != nil || localCount[k] > 1 {
			t.Fatalf("cache miss(%s)", k)
		}
	}
	if view, err := g.Get("unknow"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
