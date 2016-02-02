package rbt

import (
    "testing"
    "math/rand"
    "time"
)
var _,_ = rand.Seed, time.Now

// fill rbtree with random numbers as keys
func newtree(t *testing.T, iters int) *RbMap {
    rand.Seed(time.Now().UnixNano())
    // rb tree with integer keys
    r := NewRbMap(func(k1, k2 interface{}) bool { 
        return k1.(int) < (k2).(int)
    })
    if r == nil { 
        t.Fatalf("Can't create map of size %d", iters)
    }
    realsize := 0
    for i := 0; i < iters; i++ {
        num := rand.Intn(100000000)
        if r.Insert(num, i) {
            realsize++
        }
    }
    if r.Size() != realsize {
         t.Fatalf("Realsize mismatch: %d/%d", r.Size(), realsize);
    }
    r.verify()
    return r
}

func xTestFill(t *testing.T) {
    r := newtree(t, 100000)
    cnt_forward, cnt_backward := 0, 0
    for n := r.First(); n != nil; cnt_forward++ {
        n = n.Next()
    }
    for n := r.Last(); n != nil; cnt_backward++ {
        n = n.Prev()
    }
    if cnt_forward != r.Size() || cnt_backward != r.Size() { t.FailNow(); }
}

func TestFind(t *testing.T) {
    r := newtree(t, 1000000)
    kl := make(map[int]int)
    n := r.First()
    for i := 0; n != nil && i < r.Size()/11; i++ {
        adv := rand.Intn(9) + 1
        for j := 0; n != nil && j < adv; j++ {
            n = n.Next()
        }
        if n != nil {
            k := n.Key().(int)
            kl[k] = k
        }
    }
    for _, k := range kl {
        n := r.FindNode(k)
        if n == nil {
            t.Fatalf("Key %d not found")
        }
        if n.Key().(int) != k {
            t.Fatalf("Key mismatch: %d/%d", n.Key().(int), k)
        }  
        r.DeleteNode(n)
        v := r.Find(k)
        if v != nil { t.Fatalf("Key %d, dup key %d", k, v) }
    }
}

func TestDelete(t *testing.T) {
    r := newtree(t, 100000)
    i := 0
    for n := r.First(); nil != n; n = r.First() {
        r.DeleteNode(n)
        if i == 10000 || i == 70000 { 
            r.verify()
        }
        i++
    }
    if r.Size() != 0 { t.Fatalf("tree size non-null after delete") }
    r = newtree(t, 100000)
    for n := r.Last(); nil != n; n = r.Last() {
        r.DeleteNode(n)
    }
    if r.Size() != 0 { t.Fatalf("tree size non-null after delete") }
}