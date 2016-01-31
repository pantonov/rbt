package rbt

import (
    "testing"
    "math/rand"
    "time"
)

// fill rbtree with random numbers as keys
func newtree(t *testing.T, iters int) *RbTree {
    rand.Seed(time.Now().UnixNano())
    // rb tree with integer keys
    r := NewRbTree(func(k1, k2 interface{}) bool { 
        return k1.(int) < (k2).(int)
    })
    if r == nil { t.FailNow() }
    realsize := 0
    for i := 0; i < iters; i++ {
        num := rand.Intn(100000000)
        if r.Insert(num, i) {
            realsize++
        }
    }
    if r.Size() != realsize { t.FailNow(); }
    return r
}

func TestFill(t *testing.T) {
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
    for i := 0; n != nil && i < 90000; i++ {
        adv := rand.Intn(10)
        for j := 0; n != nil && j < adv; j++ {
            n = n.Next()
        }
        k := n.Key().(int)
        kl[k] = k
    }
    for _, k := range kl {
        v := r.Find(k)
        if v == nil { t.FailNow(); }
        r.Delete(k)
        v = r.Find(k)
        if v != nil { t.Fatalf("Key %d, dup key %d", k, v) }
    }
}

func TestDelete(t *testing.T) {
    r := newtree(t, 100000)
    for n := r.First(); nil != n; n = r.First() {
        r.DeleteNode(n)
    }
    if r.Size() != 0 { t.Fail(); }
    
    r = newtree(t, 100000)
    for n := r.Last(); nil != n; n = r.Last() {
        r.DeleteNode(n)
    }
    if r.Size() != 0 { t.FailNow(); }
}