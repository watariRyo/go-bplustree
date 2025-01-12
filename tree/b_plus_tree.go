package tree

type BPlusTree struct {
	root *Node
}

func NewBPlusTree() *BPlusTree {
	return &BPlusTree{
		root: NewNode(true),
	}
}

func (tree *BPlusTree) Insert(key int, value any) {
	root := tree.root
	if len(root.keys) == MaxKeys {
		// 根が分割される
		newRoot := NewNode(false)
		newRoot.children = append(root.children, root)
		tree.splitChild(newRoot, 0)
		tree.root = newRoot
	}
	tree.insertNonFull(tree.root, key, value)
}

func (tree *BPlusTree) insertNonFull(node *Node, key int, value any) {
	if node.isLeaf {
		// リーフノードにキー挿入
		idx := 0
		for idx < len(node.keys) && key > node.keys[idx] {
			idx++
		}
		node.keys = append(node.keys[:idx], append([]int{key}, node.keys[idx:]...)...)
		node.values = append(node.values[:idx], append([]any{value}, node.values[idx:]...)...)
	} else {
		// 内部の子ノードに挿入
		idx := 0
		for idx < len(node.keys) && key >= node.keys[idx] {
			idx++
		}
		if len(node.children[idx].keys) == MaxKeys {
			tree.splitChild(node, idx)
			if key > node.keys[idx] {
				idx++
			}
		}
		tree.insertNonFull(node.children[idx], key, value)
	}
}

func (tree *BPlusTree) splitChild(parent *Node, index int) {
	child := parent.children[index]
	mid := MaxKeys / 2

	// 分割基準のキーを取得
	midKey := child.keys[mid]

	// 新しいノードを作成
	newChild := NewNode(child.isLeaf)

	// リーフノードの場合
	if child.isLeaf {
		// midKey を新しいノードに含める
		newChild.keys = append(newChild.keys, child.keys[mid:]...)
		newChild.values = append(newChild.values, child.values[mid:]...)
		child.keys = child.keys[:mid]
		child.values = child.values[:mid]

		// リーフノードを連結
		newChild.next = child.next
		child.next = newChild
	} else {
		// 内部ノードの場合
		newChild.keys = append(newChild.keys, child.keys[mid+1:]...)
		newChild.children = append(newChild.children, child.children[mid+1:]...)
		child.keys = child.keys[:mid]
		child.children = child.children[:mid+1]
	}

	// 親ノードに midKey を挿入
	parent.keys = append(parent.keys[:index], append([]int{midKey}, parent.keys[index:]...)...)
	parent.children = append(parent.children[:index+1], append([]*Node{newChild}, parent.children[index+1:]...)...)
}

func (tree *BPlusTree) Search(key int) (interface{}, bool) {
	current := tree.root

	for current != nil {
		i := 0
		for i < len(current.keys) && key > current.keys[i] {
			i++
		}

		if i < len(current.keys) && key == current.keys[i] {
			if current.isLeaf {
				// リーフノードの場合は対応する値を返す
				return current.values[i], true
			}
			// 内部ノードの場合は次の子ノードへ進む
			current = current.children[i+1]
		} else if current.isLeaf {
			// リーフノードで見つからなかった場合
			return nil, false
		} else {
			// 内部ノードで適切な子ノードへ進む
			current = current.children[i]
		}
	}

	return nil, false
}
