package main

import "fmt"

type BTreeNode struct {
	keys []int
	t    int
	c    []*BTreeNode
	n    int
	leaf bool
}

type BTree struct {
	root *BTreeNode
	t    int
}

func NewBTreeNode(t int, leaf bool) *BTreeNode {
	return &BTreeNode{
		keys: make([]int, t<<1-1),
		t:    t,
		c:    make([]*BTreeNode, t<<1),
		leaf: leaf,
	}
}

func (n *BTreeNode) search(k int) *BTreeNode {
	i := findGE(n.keys, 0, n.n-1, k)

	if n.keys[i] == k {
		return n
	}
	if n.leaf {
		return nil
	}
	return n.c[i].search(k)
}

func (n *BTreeNode) traverse() {
	i := 0
	for ; i < n.n; i++ {
		if !n.leaf {
			n.c[i].traverse()
		}
		fmt.Printf("%d ", n.keys[i])
	}

	if !n.leaf {
		n.c[i].traverse()
	}
}

func (n *BTreeNode) splitChild(i int, y *BTreeNode) {
	z := NewBTreeNode(y.t, y.leaf)
	z.n = n.t - 1

	for j := 0; j < n.t-1; j++ {
		z.keys[j] = y.keys[n.t+j]
	}
	if !y.leaf {
		for j := 0; j < n.t; j++ {
			z.c[j] = y.c[n.t+j]
		}
	}

	y.n = n.t - 1
	for j := n.n; j > i; j-- {
		n.c[j+1] = n.c[j]
	}
	n.c[i+1] = z

	for j := n.n - 1; j >= i; j-- {
		n.keys[j+1] = n.keys[j]
	}
	n.keys[i] = y.keys[n.t-1]
	n.n = n.n + 1
}

func (n *BTreeNode) insertNonFull(k int) {
	if n.leaf {
		i := n.n - 1

		for ; i >= 0 && n.keys[i] > k; i-- {
			n.keys[i+1] = n.keys[i]
		}
		n.keys[i+1] = k
		n.n = n.n + 1
	} else {
		i := findGE(n.keys, 0, n.n-1, k)

		if n.c[i].n == (n.t)<<1-1 {
			n.splitChild(i, n.c[i])
			if k > n.keys[i] {
				i++
			}
		}
		n.c[i].insertNonFull(k)
	}
}

func (n *BTreeNode) getPred(i int) int {
	pred := n.c[i]

	for !pred.leaf {
		pred = pred.c[pred.n]
	}
	return pred.keys[pred.n-1]
}

func (n *BTreeNode) getSucc(i int) int {
	succ := n.c[i+1]

	for !succ.leaf {
		succ = succ.c[0]
	}
	return succ.keys[0]
}

func (n *BTreeNode) removeFromLeaf(i int) {
	for j := i + 1; j < n.n; j++ {
		n.keys[j-1] = n.keys[j]
	}

	(n.n)--
}

func (n *BTreeNode) borrowFromPrev(i int) {
	child := n.c[i]
	sibling := n.c[i-1]

	for j := child.n - 1; j >= 0; j-- {
		child.keys[j+1] = child.keys[j]
	}
	if !child.leaf {
		for j := child.n; j >= 0; j-- {
			child.c[j+1] = child.c[j]
		}
	}
	child.keys[0] = n.keys[i-1]
	if !child.leaf {
		child.c[0] = sibling.c[sibling.n]
	}

	n.keys[i-1] = sibling.keys[sibling.n-1]
	(child.n)++
	(sibling.n)--
}

func (n *BTreeNode) borrowFromNext(i int) {
	child := n.c[i]
	sibling := n.c[i+1]

	child.keys[child.n] = n.keys[i]
	if !child.leaf {
		child.c[child.n+1] = sibling.c[0]
	}
	n.keys[i] = sibling.keys[0]
	for j := 1; j < sibling.n; j++ {
		sibling.keys[j-1] = sibling.keys[j]
	}
	if !sibling.leaf {
		for j := 1; j <= sibling.n; j++ {
			sibling.c[j-1] = sibling.c[j]
		}
	}
	(child.n)++
	(sibling.n)--
}

func (n *BTreeNode) merge(i int) {
	child := n.c[i]
	sibling := n.c[i+1]

	child.keys[n.t-1] = n.keys[i]
	for j := 0; j < sibling.n; j++ {
		child.keys[n.t+j] = sibling.keys[j]
	}
	if !child.leaf {
		for j := 0; j <= sibling.n; j++ {
			child.c[n.t+j] = sibling.c[j]
		}
	}
	for j := i + 1; j < n.n; j++ {
		n.keys[j-1] = n.keys[j]
	}
	for j := i + 2; j <= n.n; j++ {
		n.c[j-1] = n.c[j]
	}
	child.n += (sibling.n + 1)
	(n.n)--
}

func (n *BTreeNode) fill(i int) {
	if i != 0 && n.c[i-1].n >= n.t {
		n.borrowFromPrev(i)
	} else if i != n.n && n.c[i+1].n >= n.t {
		n.borrowFromNext(i)
	} else {
		if i == n.n {
			n.merge(i - 1)
		} else {
			n.merge(i)
		}
	}
}

func (n *BTreeNode) remove(k int) {
	i := findGE(n.keys, 0, n.n-1, k)
	if n.keys[i] == k {
		if n.leaf {
			n.removeFromLeaf(i)
		} else {
			n.removeFromNonLeaf(i)
		}
	} else {
		if n.leaf {
			fmt.Printf("The key %d is does not exist in the tree\n", k)
			return
		}
		flag := false
		if i == n.n {
			flag = true
		}
		if n.c[i].n < n.t {
			n.fill(i)
		}

		if flag && i > n.n {
			n.c[i-1].remove(k)
		} else {
			n.c[i].remove(k)
		}
	}
}

func (n *BTreeNode) removeFromNonLeaf(i int) {
	k := n.keys[i]

	if n.c[i].n >= n.t {
		pred := n.getPred(i)
		n.keys[i] = pred
		n.c[i].remove(pred)
	} else if n.c[i+1].n >= n.t {
		succ := n.getSucc(i)
		n.keys[i] = succ
		n.c[i+1].remove(succ)
	} else {
		n.merge(i)
		n.c[i].remove(k)
	}
}

func NewBTree(t int) *BTree {
	return &BTree{
		t: t,
	}
}

func (t *BTree) Search(k int) *BTreeNode {
	if t.root != nil {
		return t.root.search(k)
	}

	return nil
}

func (t *BTree) Traverse() {
	if t.root != nil {
		t.root.traverse()
	}
}

func (t *BTree) Insert(k int) {
	if t.root == nil {
		t.root = NewBTreeNode(t.t, true)
		t.root.keys[0] = k
		t.root.n = 1
	} else {
		if t.root.n == (t.t)<<1-1 {
			newRoot := NewBTreeNode(t.t, false)
			newRoot.c[0] = t.root
			newRoot.insertNonFull(k)
			t.root = newRoot
		} else {
			t.root.insertNonFull(k)
		}
	}
}

func (t *BTree) Remove(k int) {
	if t.root == nil {
		fmt.Println("The tree is empty")
		return
	}

	t.root.remove(k)

	if t.root.n == 0 {
		if t.root.leaf {
			t.root = nil
		} else {
			t.root = t.root.c[0]
		}
	}
}

func findGE(s []int, left, right, k int) int {
	if left <= right {
		mid := left + (right-left)>>1

		if k == s[mid] {
			return mid
		} else if k > s[mid] {
			return findGE(s, mid+1, right, k)
		} else {
			return findGE(s, left, mid-1, k)
		}
	}

	return left
}

func main() {
	t := NewBTree(3)
	t.Insert(10)
	t.Insert(20)
	t.Insert(5)
	t.Insert(6)
	t.Insert(12)
	t.Insert(30)
	t.Insert(7)
	t.Insert(17)

	fmt.Println("Traversal of the constructed tree is ")
	t.Traverse()
	fmt.Println()

	t.Remove(5)
	fmt.Println("Traversal of the constructed tree is after remove 5")
	t.Traverse()
	fmt.Println()

	t.Remove(7)
	fmt.Println("Traversal of the constructed tree is after remove 7")
	t.Traverse()
	fmt.Println()

	t.Remove(10)
	fmt.Println("Traversal of the constructed tree is after remove 10")
	t.Traverse()
	fmt.Println()

	t.Remove(20)
	fmt.Println("Traversal of the constructed tree is after remove 20")
	t.Traverse()
	fmt.Println()

	t.Remove(6)
	fmt.Println("Traversal of the constructed tree is after remove 6")
	t.Traverse()
	fmt.Println()

	t.Remove(12)
	fmt.Println("Traversal of the constructed tree is after remove 12")
	t.Traverse()
	fmt.Println()

	t.Remove(17)
	fmt.Println("Traversal of the constructed tree is after remove 17")
	t.Traverse()
	fmt.Println()

	t.Remove(99)
	fmt.Println("Traversal of the constructed tree is after remove 99")
	t.Traverse()
	fmt.Println()

	t.Remove(30)
	fmt.Println("Traversal of the constructed tree is after remove 30")
	t.Traverse()
	fmt.Println()

	t.Remove(50)
	fmt.Println("Traversal of the constructed tree is after remove 50")
	t.Traverse()
	fmt.Println()
}
