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
		newRoot.children = append(newRoot.children, root)
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
		for idx < len(node.keys) && key > node.keys[idx] {
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
	parent.children = append(parent.children[:index+1], parent.children[index:]...)
	parent.children[index+1] = newChild
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

func (tree *BPlusTree) Delete(key int) bool {
	if tree.root == nil {
		return false // 木が空の場合
	}

	deleted := tree.deleteFromNode(tree.root, key)

	// ルートが空になった場合の処理
	if len(tree.root.keys) == 0 {
		if !tree.root.isLeaf {
			tree.root = tree.root.children[0] // 新しいルートに置き換え
		} else {
			tree.root = nil // 木全体が空
		}
	}

	return deleted
}

func (tree *BPlusTree) deleteFromNode(node *Node, key int) bool {
	idx := 0

	// キーの位置を特定
	for idx < len(node.keys) && key > node.keys[idx] {
		idx++
	}

	if node.isLeaf {
		// リーフノードの場合
		if idx < len(node.keys) && node.keys[idx] == key {
			// キーを削除
			node.keys = append(node.keys[:idx], node.keys[idx+1:]...)
			node.values = append(node.values[:idx], node.values[idx+1:]...)
			return true
		}
		return false // キーが見つからない
	}

	// 内部ノードの場合
	if idx < len(node.keys) && node.keys[idx] == key {
		// 内部ノードでキーが見つかった場合
		return tree.deleteInternalNode(node, idx)
	}

	// 子ノードに再帰的に削除を適用
	child := node.children[idx]
	deleted := tree.deleteFromNode(child, key)

	// 再平衡処理
	if len(child.keys) < MaxKeys/2 {
		tree.fixUnderflow(node, idx)
	}

	return deleted
}

func (tree *BPlusTree) deleteInternalNode(node *Node, idx int) bool {
	leftChild := node.children[idx]
	rightChild := node.children[idx+1]

	if len(rightChild.keys) >= MaxKeys/2 {
		// 右の子ノードの2番目を親に昇格
		raisingKey, _ := tree.getSecondMin(rightChild)
		node.keys[idx] = raisingKey

		successorKey, successorValue := tree.getMin(rightChild)
		if node.isLeaf {
			node.values[idx] = successorValue
		}
		// 右の子ノードから最小キーを取得して置き換え
		tree.deleteFromNode(rightChild, successorKey)
	} else {
		// 左右の子ノードをマージ
		node, idx = tree.mergeNodes(node, idx)
		tree.deleteFromNode(leftChild, node.keys[idx])
	}

	// 削除後に半分より小さくなる場合、マージ処理
	if len(rightChild.keys) < MaxKeys/2 {
		// 左右の子ノードをマージ
		tree.mergeNodes(node, idx)
	}

	return true
}

func (tree *BPlusTree) fixUnderflow(parent *Node, idx int) {
	child := parent.children[idx]
	println(idx)
	if idx > 0 {
		// 左隣の兄弟ノード
		leftSibling := parent.children[idx-1]
		if len(leftSibling.keys) > MaxKeys/2 {
			// 左の兄弟からキーを再分配
			child.keys = append([]int{leftSibling.keys[len(leftSibling.keys)-1]}, child.keys...)
			child.values = append([]any{leftSibling.values[len(leftSibling.values)-1]}, child.values...)
			parent.keys[idx-1] = leftSibling.keys[len(leftSibling.keys)-1]
			leftSibling.keys = leftSibling.keys[:len(leftSibling.keys)-1]
			leftSibling.values = leftSibling.values[:len(leftSibling.values)-1]
			return
		}
	}

	if idx < len(parent.children)-1 {
		// 右隣の兄弟ノード
		rightSibling := parent.children[idx+1]
		if len(rightSibling.keys) > MaxKeys/2 {
			// 右の兄弟からキーを再分配
			child.keys = append(child.keys, parent.keys[idx])
			child.values = append(child.values, rightSibling.values[0])
			parent.keys[idx] = rightSibling.keys[1] // 親keyを2番目の要素に移動
			rightSibling.keys = rightSibling.keys[1:]
			rightSibling.values = rightSibling.values[1:]
			return
		}
	}

	// 再分配ができない場合、兄弟ノードとマージ
	if idx > 0 {
		tree.mergeNodes(parent, idx-1)
	} else {
		tree.mergeNodes(parent, idx)
	}
}

func (tree *BPlusTree) mergeNodes(parent *Node, idx int) (*Node, int) {
	leftChild := parent.children[idx]
	rightChild := parent.children[idx+1]

	// 親キーを左子ノードに移動（内部ノードの場合）
	if !leftChild.isLeaf {
		leftChild.keys = append(leftChild.keys, parent.keys[idx])
	}

	// 右子ノードの内容を左子ノードにマージ
	leftChild.keys = append(leftChild.keys, rightChild.keys...)
	leftChild.values = append(leftChild.values, rightChild.values...)
	if !leftChild.isLeaf {
		leftChild.children = append(leftChild.children, rightChild.children...)
	}
	leftChild.next = rightChild.next

	// root置き換えに備え保持
	key := parent.keys[0]
	// 親ノードを更新
	parent.keys = append(parent.keys[:idx], parent.keys[idx+1:]...)
	parent.children = append(parent.children[:idx+1], parent.children[idx+2:]...)

	// マージ後にキー数が MaxKeys を超える場合、再分割
	if len(leftChild.keys) > MaxKeys {
		// 再分割を行う
		tree.splitChild(parent, idx)
	}

	// 更新でparentが0になった場合、唯一の子をルートにする
	if len(parent.keys) == 0 {
		i := 0
		// キーの位置を再特定
		for i < len(leftChild.keys) && key > leftChild.keys[i] {
			i++
		}
		tree.root = leftChild
		return tree.root, i
	} else {
		return parent, idx
	}
}

func (tree *BPlusTree) getMin(node *Node) (int, any) {
	current := node
	for !current.isLeaf {
		// 内部ノードの場合は左端の子ノードへ進む
		current = current.children[0]
	}
	// リーフノードの最初のキーと値を返す
	return current.keys[0], current.values[0]
}

func (tree *BPlusTree) getSecondMin(node *Node) (int, any) {
	current := node
	for !current.isLeaf {
		// 内部ノードの場合は左端の子ノードへ進む
		current = current.children[0]
	}
	// リーフノードの2番目のキーと値を返す（呼び出し元限定。エラーチェックしない）
	return current.keys[1], current.values[1]
}
