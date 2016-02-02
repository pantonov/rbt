// Red-Black tree implementation. Each tree node (entry) contains Key and Value.
// Entries in the RbMap are always ordered according to the key value, so
// it can be used as ordered set (std::set in C++) or ordered map (std::map).
// Note: all methods are not goroutine-safe.
package rbt

// import ( "strings" ; "fmt" )

type RbMap struct {
    less       LessFunc
    root       *RbMapNode
    size       int
}

// Red-black tree node, contains key and value. It is safe to overwrite Value
// in-place.
type RbMapNode struct {
    left, right, parent *RbMapNode
    // key
    key          interface{}
    Value        interface{}
    isred        bool         // true == red, false == black
}

// LessFunc is a key comparsion function. 
// Must return true if k1 < k2, false otherwise.
type LessFunc func(k1, k2 interface{}) bool

// Create new RbMap with provided key comparsion function. 
func NewRbMap(lessFunc LessFunc) *RbMap {
    return &RbMap{ less: lessFunc }
}

// Find node by key and return its Value, returns nil if key not found.
func (t *RbMap) Find(key interface{}) interface{} {
    n := t.FindNode(key)
    if n != nil {
        return n.Value
    }
    return nil
}

// Find a node by key, returns nil if not found.
func (t *RbMap) FindNode(key interface{}) *RbMapNode {
    x := t.root
    for x != nil {
        if t.less(x.key, key) {
            x = x.right
        } else {
            if t.less(key, x.key) {
                x = x.left
            } else {
                return x
            }
        }
    }
    return nil
}

// Get last node in the tree (with highest key value).
func (t *RbMap) Last() *RbMapNode {
    if nil == t.root {
        return nil
    }
    return t.root.max()
}

// Get first node in the tree (with lowest key value).
func (t *RbMap) First() *RbMapNode {
    if nil == t.root {
        return nil
    }
    return t.root.min()
}

// Get next node, in ascending key value order.
func (x *RbMapNode) Next() *RbMapNode {
    if x.right != nil {
        return x.right.min()
    }
    y := x.parent
    for y != nil && x == y.right {
        x = y
        y = y.parent
    }
    return y
}

// Returns key associated with tree node.
func (x *RbMapNode) Key() interface{} {
    return x.key
}

// Get previous node, in descending key value order.
func (x *RbMapNode) Prev() *RbMapNode {
    if x.left != nil {
        return x.left.max()
    }
    y := x.parent
    for y != nil && x == y.left {
        x = y
        y = y.parent
    }
    return y
}

// Returns number of entries in the tree. This function returns internal
// counter, therefore it is fast and safe to use in loops.
func (t *RbMap) Size() int {
    return t.size
}

// Remove all entries in the tree.
func (t *RbMap) Clear() {
    t.root = nil
    t.size = 0
}

// Insert key and value into the tree. If new entry is created, returns true.
// If key already exists, value gets replaced and Insert returns false.
func (t *RbMap) Insert(key interface{}, value interface{}) bool {
    x := t.root
    var y *RbMapNode

    for x != nil {
        y = x
        if t.less(x.key, key) {
            x = x.right
        } else if t.less(key, x.key) {
            x = x.left
        } else {
            x.Value = value
            return false // overwrite value
        }
    }
    z := &RbMapNode{parent: y, isred: true, key: key, Value: value}
    if y == nil {
        t.root = z
    } else {
        if t.less(key, y.key) {
            y.left = z
        } else {
            y.right = z
        }
    }
    t.rb_insert_fixup(z)
    t.size++
    return true
}

// Delete tree node by key. Returns true if key was found and deleted.
func (t *RbMap) Delete(key interface{}) bool {
    if z := t.FindNode(key); z != nil {
        t.DeleteNode(z)
        return true
    }
    return false
}

// Delete tree node.
func (t *RbMap) DeleteNode(n *RbMapNode) {
    var x *RbMapNode
    if nil != n.left && nil != n.right {
        x = n.left.max()
        n.key, n.Value = x.key, x.Value
        n = x
    }
    if nil == n.right {
        x = n.left
    } else {
        x = n.right
    }
    if isBlack(n) {
        n.isred = isRed(x)
        if nil != n.parent {
            t.rb_delete_fixup(n)
        }
    }
    t.rbreplace(n, x)
    if isRed(t.root) {
        t.root.isred = false
    }
    t.size--
}

func (t* RbMap) rb_delete_fixup(n *RbMapNode) {
    var s, p *RbMapNode
    for {
        s, p = n.sibling(), n.parent
        if isRed(s) {
            p.isred, s.isred = true, false
            if n == p.left {
                t.left_rotate(p)
                s = p.right
            } else {
                t.right_rotate(p)
                s = p.left
            }
        }
        if isBlack(p) && isBlack(s) && isBlack(s.left) && isBlack(s.right) {
            s.isred = true
            if nil != p.parent {
                n = p
                continue
            }
            return
        } else {
            break
        }
    }
    if isRed(n.parent) && isBlack(s) && isBlack(s.left) && isBlack(s.right) {
        s.isred, n.parent.isred = true, false
    } else {
        if isBlack(s) {
            if n == n.parent.left && isRed(s.left) && isBlack(s.right) {
                s.isred, s.left.isred = true, false
                t.right_rotate(s)
                s = n.parent.right
            } else if n == n.parent.right && isRed(s.right) && isBlack(s.left) {
                s.isred, s.right.isred = true, false
                t.left_rotate(s)
                s = n.parent.left
            }
        }
        s.isred = n.parent.isred
        n.parent.isred = false
        if n == n.parent.left {
            s.right.isred = false
            t.left_rotate(n.parent)
        } else {
            s.left.isred = false
            t.right_rotate(n.parent)
        }
    }
}

func (t *RbMap) rb_insert_fixup(x *RbMapNode) {
    var y *RbMapNode
    for isRed(x.parent) {
        if x.parent == x.parent.parent.left {
            y = x.parent.parent.right
            if isRed(y) {
                x.parent.isred, y.isred = false, false
                x.parent.parent.isred = true
                x = x.parent.parent
            } else {
                if x == x.parent.right {
                    x = x.parent
                    t.left_rotate(x)
                }
                x.parent.isred, x.parent.parent.isred = false, true
                t.right_rotate(x.parent.parent)
            }
        } else {
            y = x.parent.parent.left
            if isRed(y) {
                x.parent.isred, y.isred = false, false
                x.parent.parent.isred = true
                x = x.parent.parent
            } else {
                if x == x.parent.left {
                    x = x.parent
                    t.right_rotate(x)
                }
                x.parent.isred, x.parent.parent.isred = false, true
                t.left_rotate(x.parent.parent)
            }
        }
    }
    t.root.isred = false
}

func (n *RbMapNode) sibling() *RbMapNode {
    if n == n.parent.left {
        return n.parent.right
    } else {
        return n.parent.left
    }
}

func (n *RbMapNode) min() *RbMapNode {
    for n.left != nil {
        n = n.left
    }
    return n
}

func (n *RbMapNode) max() *RbMapNode {
    for n.right != nil {
        n = n.right
    }
    return n
}

func (t *RbMap) left_rotate(n *RbMapNode) {
    r := n.right
    t.rbreplace(n, r)
    n.right = r.left
    if nil != r.left {
        r.left.parent = n
    } 
    r.left, n.parent = n, r
}

func (t *RbMap) right_rotate(n *RbMapNode) {
    l := n.left
    t.rbreplace(n, l)
    n.left = l.right
    if nil != l.right {
        l.right.parent = n
    }
    l.right, n.parent = n, l
}

func (t *RbMap) rbreplace(u, v *RbMapNode) {
    parent := u.parent
    if parent == nil {
        t.root = v
    } else if u == parent.left {
        parent.left = v
    } else {
        parent.right = v
    }
    if v != nil {
        v.parent = parent
    }
}

func isBlack(n *RbMapNode) bool {
    return nil == n || !n.isred
}

func isRed(n *RbMapNode) bool {
    return nil != n && n.isred
}

/*

// uncomment for debugging only, because this adds dependency on strings, fmt

func (n *RbMapNode) Dump(indent int, tag string) {
    idn := strings.Repeat(" ", indent*4)
    c := 'B'
    if n.isred { c = 'R' }
    fmt.Printf("%s%s[%v:%v]%c\n", idn, tag, n.Key(), n.Value, c)
    if n.left != nil {
        n.left.Dump(indent + 1, "L:")
    }
    if n.right != nil {
        n.right.Dump(indent + 1, "R:")
    }
}

func (t *RbMap) Dump() {
    if t.root == nil {
        fmt.Printf("<NULL TREE>\n")
    } else {
        t.root.Dump(0, "*")
    }
}
*/

// Internal tree consistency check used by tests. 
func (t *RbMap) verify() {
    if nil == t.root { return }
    if isRed(t.root) { panic("root is red") }
    verify1(t.root)
    verify2(t.root)
}

func verify1(n *RbMapNode) {
    if isRed(n) {
        if !isBlack(n.left)   { panic("left is not black") }
        if !isBlack(n.right)  { panic("right is not black") }
        if !isBlack(n.parent) { panic("parent is not black") }
    } 
    if nil == n { return }
    verify1(n.left)
    verify1(n.right)
}

func verify2(n *RbMapNode) {
    black_count_path := -1
    verify2h(n, 0, &black_count_path)
}

func verify2h(n *RbMapNode, black_count int, path_black_count *int) {
    if isBlack(n) {
        black_count++
    }
    if n == nil {
        if *path_black_count == -1 {
            *path_black_count = black_count;
        } else {
            if black_count != *path_black_count {
                panic("black count")
            }
        }
        return
    }
    verify2h(n.left,  black_count, path_black_count)
    verify2h(n.right, black_count, path_black_count)
}

