// Red-Black tree implementation. Each tree node (entry) contains Key and Value.
// Entries in the RbTree are always ordered according to the key value, so
// it can be used as ordered set (std::set in C++) or ordered map (std::map).
// Note: all methods are not goroutine-safe.
package rbt

type RbTree struct {
    less       LessFunc
    root       *RbTreeNode
    size       int
}

// Red-black tree node, contains key and value. It is safe to overwrite Value
// in-place.
type RbTreeNode struct {
    left, right, parent *RbTreeNode
    // key
    key          interface{}
    Value        interface{}
    isred        bool         // true == red, false == black
}

// LessFunc is a key comparsion function. 
// Must return true if k1 < k2, false otherwise
type LessFunc func(k1, k2 interface{}) bool

// Create new RbTree with provided comparsion function. The latter must
// return true if k1 < k2, false otherwise.
func NewRbTree(lessFunc LessFunc) *RbTree {
    return &RbTree{ less: lessFunc }
}

// Find node by key and return its Value, returns nil if key not found.
func (t *RbTree) Find(key interface{}) interface{} {
    n := t.FindNode(key)
    if n != nil {
        return n.Value
    }
    return nil
}

// Find a node by key, returns nil if not found
func (t *RbTree) FindNode(key interface{}) *RbTreeNode {
    x := t.root
    for x != nil {
        if t.less(x.key, key) {
            x = x.left
        } else {
            if t.less(key, x.key) {
                x = x.right
            } else {
                return x
            }
        }
    }
    return nil
}

// Get last node in the tree (with highest key value).
func (t *RbTree) Last() *RbTreeNode {
    if nil == t.root {
        return nil
    }
    return t.root.min()
}

// Get first node in the tree (with lowest key value).
func (t *RbTree) First() *RbTreeNode {
    if nil == t.root {
        return nil
    }
    return t.root.max()
}

// Returns previous node in the tree, in descending key value order.
func (x *RbTreeNode) Prev() *RbTreeNode {
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
func (x *RbTreeNode) Key() interface{} {
    return x.key
}

// Returns next node in the tree, in ascending key value order.
func (x *RbTreeNode) Next() *RbTreeNode {
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
func (t *RbTree) Size() int {
    return t.size
}

// Remove all entries in the tree.
func (t *RbTree) Clear() {
    t.root = nil
    t.size = 0
}

// Insert key and value into the tree. If new entry is created, returns true.
// If key already exists, value gets replaced and Insert returns false.
func (t *RbTree) Insert(key interface{}, value interface{}) bool {
    x := t.root
    var y *RbTreeNode

    for x != nil {
        y = x
        if t.less(x.key, key) {
            x = x.left
        } else if t.less(key, x.key) {
            x = x.right
        } else {
            x.Value = value
            return false // overwrite value
        }
    }
    z := &RbTreeNode{parent: y, isred: true, key: key, Value: value}
    if y == nil {
        t.root = z
    } else {
        if t.less(y.key, key) {
            y.left = z
        } else {
            y.right = z
        }
    }
    t.rb_insert_fixup(z)
    t.size++
    return true
}

// Delete tree node by key.
func (t *RbTree) Delete(key interface{}) {
    if z := t.FindNode(key); z != nil {
        t.DeleteNode(z)
    }
}

// Delete tree node.
func (t *RbTree) DeleteNode(z *RbTreeNode) {
    var x, y, parent *RbTreeNode
    y, y_original_isred, parent := z, z.isred, z.parent
    if z.left == nil {
        x = z.right
        t.rbtransplant(z, z.right)
    } else if z.right == nil {
        x = z.left
        t.rbtransplant(z, z.left)
    } else {
        y = z.right.min()
        y_original_isred, x = y.isred, y.right
        if y.parent == z {
            if x == nil {
                parent = y
            } else {
                x.parent = y
            }
        } else {
            t.rbtransplant(y, y.right)
            y.right = z.right
            y.right.parent = y
        }
        t.rbtransplant(z, y)
        y.left = z.left
        y.left.parent, y.isred = y, z.isred
    }
    if !y_original_isred && x != nil {
        t.rb_delete_fixup(x, parent)
    }
    t.size--
}

func (t *RbTree) rb_insert_fixup(x *RbTreeNode) {
    var y *RbTreeNode
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

func (t *RbTree) rb_delete_fixup(x, parent *RbTreeNode) {
    var w *RbTreeNode
    for x != t.root && isBlack(x) {
        if x == parent.left {
            w = parent.right
            if isRed(w) {
                w.isred, parent.isred = false, true
                t.left_rotate(parent)
                w = parent.right
            }
            if w != nil && isBlack(w.left) && isBlack(w.right) {
                w.isred = true 
                x = parent
                parent = x.parent
            } else {
                if w != nil && isBlack(w.right) {
                    if w.left != nil {
                        w.left.isred = false
                    }
                    w.isred = true
                    t.right_rotate(w)
                    w = parent.right
                }
                if w != nil {
                    w.isred = parent.isred
                    if w.right != nil {
                        w.right.isred = false
                    }
                }
                parent.isred = false
                t.left_rotate(parent)
                x = t.root
            }
        } else {
            w := parent.left
            if isRed(w) {
                w.isred, parent.isred = false, true
                t.right_rotate(parent)
                w = parent.left
            }
            if w != nil && isBlack(w.left) && isBlack(w.right) {
                w.isred = true
                x = parent
                parent = x.parent
            } else {
                if w != nil && isBlack(w.left) {
                    w.isred = true
                    if w.right != nil {
                        w.right.isred = false
                    }
                    t.left_rotate(w)
                    w = parent.left
                }
                if w != nil {
                    w.isred = parent.isred
                    if w.left != nil {
                        w.left.isred = false
                    }
                }
                parent.isred = false
                t.right_rotate(parent)
                x = t.root
            }
        }
    }
    if x != nil {
        x.isred = false
    }
}

func (n *RbTreeNode) min() *RbTreeNode {
    for n.left != nil {
        n = n.left
    }
    return n
}

func (n *RbTreeNode) max() *RbTreeNode {
    for n.right != nil {
        n = n.right
    }
    return n
}

func (t *RbTree) left_rotate(x *RbTreeNode) {
    y := x.right
    x.right = y.left
    if y.left != nil {
        y.left.parent = x
    }
    y.parent = x.parent
    if x.parent == nil {
        t.root = y
    } else {
        if x == x.parent.left {
          x.parent.left = y
       } else {
          x.parent.right = y
       }
    }
    y.left, x.parent = x, y
}

func (t *RbTree) right_rotate(y *RbTreeNode) {
    x := y.left
    y.left = x.right
    if x.right != nil {
        x.right.parent = y
    }
    x.parent = y.parent
    if y.parent == nil {
        t.root = x
    } else { 
        if y == y.parent.left {
          y.parent.left = x
       } else {
          y.parent.right = x
       }
    }
    x.right, y.parent = y, x
}

func (t *RbTree) rbtransplant(u, v *RbTreeNode) {
    if u.parent == nil {
        t.root = v
    } else if u == u.parent.left {
        u.parent.left = v
    } else {
        u.parent.right = v
    }
    if v == nil {
        return
    }
    v.parent = u.parent
}

func isBlack(n *RbTreeNode) bool {
    return nil == n || !n.isred
}

func isRed(n *RbTreeNode) bool {
    return nil != n && n.isred
}

/*
// uncomment for debugging only
func (n *RbTreeNode) Dump(indent int, tag string) {
    idn := strings.Repeat(" ", indent*4)
    fmt.Printf("%s%s[%v:%v]%d\n", idn, tag, n.Key, n.Value, n.isred)
    if n.left != nil {
        n.left.Dump(indent + 1, "L:")
    }
    if n.right != nil {
        n.right.Dump(indent + 1, "R:")
    }
}

func (t *RbTree) Dump() {
    if t.root == nil {
        fmt.Printf("<NULL TREE>\n")
    } else {
        t.root.Dump(0, "*")
    }
}
*/