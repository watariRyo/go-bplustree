package tree

import (
	"reflect"
	"testing"
)

func Contains(list interface{}, elem interface{}) bool {
	listV := reflect.ValueOf(list)

	if listV.Kind() == reflect.Slice {
		for i := 0; i < listV.Len(); i++ {
			item := listV.Index(i).Interface()
			// 型変換可能か確認する
			if !reflect.TypeOf(elem).ConvertibleTo(reflect.TypeOf(item)) {
				continue
			}
			// 型変換する
			target := reflect.ValueOf(elem).Convert(reflect.TypeOf(item)).Interface()
			// 等価判定をする
			if ok := reflect.DeepEqual(item, target); ok {
				return true
			}
		}
	}
	return false
}

func Test_Insert(t *testing.T) {
	// keyの最大数は4
	t.Run("Test Only Reaf Node.", func(t *testing.T) {
		tree := NewBPlusTree()

		tree.Insert(10, "A")
		tree.Insert(20, "B")
		tree.Insert(5, "C")
		tree.Insert(15, "D")

		keyExpected := []int{10, 20, 5, 15}
		valueExpected := []string{"C", "A", "D", "B"}

		keys := tree.root.keys
		values := tree.root.values

		for idx, key := range keys {
			if !Contains(keyExpected, key) {
				t.Errorf("Could not find key %v in tree %v", key, tree)
			}
			got := values[idx]
			want := valueExpected[idx]
			if got != want {
				t.Errorf("Key missmatched value. Got %v, want %v", got, want)
			}
		}
	})

	t.Run("Test Inner Node And Reaf Node KeyValue And Next.", func(t *testing.T) {
		tree := NewBPlusTree()

		tree.Insert(10, "A")
		tree.Insert(20, "B")
		tree.Insert(5, "C")
		tree.Insert(15, "D")
		tree.Insert(1, "E")

		// 1, 5, 10, 15, 20
		// key = 15
		// child1 1, 5, 10 next is exist
		// child2 15, 20 next is nil

		// root key is 15
		if tree.root.keys[0] != 15 {
			t.Errorf("root key is not set correctly. set %v", tree.root.keys)
		}

		// root value is nil
		if len(tree.root.values) != 0 {
			t.Errorf("root value is not set correctly. set %v", tree.root.values)
		}

		child1ExpectedKeys := []int{1, 10, 5}
		child1ExpectedValues := []any{"E", "C", "A"}

		child1keys := tree.root.children[0].keys
		child1values := tree.root.children[0].values

		child2ExpectedKeys := []int{15, 20}
		child2ExpectedValues := []any{"D", "B"}

		child2keys := tree.root.children[1].keys
		child2values := tree.root.children[1].values

		child2 := tree.root.children[1]
		child1Next := tree.root.children[0].next

		testLeafNode(t, child1keys, child1ExpectedKeys, child1values, child1ExpectedValues, child1Next, child2)
		testLeafNode(t, child2keys, child2ExpectedKeys, child2values, child2ExpectedValues, nil, nil)

		if child2.next != nil {
			t.Errorf("last leaf should not have next.")
		}
	})

	t.Run("Test Inner Node Split Reach MaxKeys.", func(t *testing.T) {
		tree := NewBPlusTree()

		tree.Insert(10, "A")
		tree.Insert(20, "B")
		tree.Insert(5, "C")
		tree.Insert(15, "D")
		tree.Insert(1, "E")
		tree.Insert(2, "F")
		tree.Insert(3, "G")
		tree.Insert(4, "H")
		tree.Insert(9, "I")
		tree.Insert(8, "J")
		tree.Insert(7, "K")
		tree.Insert(6, "L")
		tree.Insert(11, "M")
		tree.Insert(12, "N")
		tree.Insert(13, "O")

		// rootが分割される境界

		// 5, 9, 11, 15
		// 1, 2, 3, 4
		// 5, 6, 7, 8
		// 9, 10,
		// 11, 12, 13,
		// 15, 20

		// set 14
		tree.Insert(14, "P")

		// 11
		// child: 5, 9
		// 5----
		// 1, 2, 3, 4
		// 5, 6, 7, 8
		// ----5
		// 9----
		// 9, 10
		// ----9
		// child: 15
		// 11, 12, 13, 14 <- set reaf node
		// 15, 20

		// root key is 11
		if tree.root.keys[0] != 11 {
			t.Errorf("root key is not set correctly. set %v", tree.root.keys)
		}

		// root value is nil
		if len(tree.root.values) != 0 {
			t.Errorf("root value is not set correctly. set %v", tree.root.values)
		}

		// inner node 1 key is 5, 9
		child1ExpectedKeys := []int{5, 9}
		child1keys := tree.root.children[0].keys

		for _, key := range child1keys {
			if !Contains(child1ExpectedKeys, key) {
				t.Errorf("Could not find key %v in tree %v", key, tree)
			}
		}
		if len(tree.root.children[0].values) != 0 {
			t.Errorf("inner node value shuold be nil.")
		}

		// inner node 2 key is 15
		if tree.root.children[1].keys[0] != 15 {
			t.Errorf("inner node 2 key is not set correctly. set %v", tree.root.keys)
		}

		if len(tree.root.children[1].values) != 0 {
			t.Errorf("root value is not set correctly. set %v", tree.root.values)
		}

		// leaf node 1
		leaf1ExpectedKeys := []int{1, 2, 3, 4}
		leaf1ExpectedValues := []any{"E", "F", "G", "H"}

		leaf1keys := tree.root.children[0].children[0].keys
		leaf1values := tree.root.children[0].children[0].values

		leaf2 := tree.root.children[0].children[1]
		leaf1Next := tree.root.children[0].children[0].next

		testLeafNode(t, leaf1keys, leaf1ExpectedKeys, leaf1values, leaf1ExpectedValues, leaf1Next, leaf2)

		// leaf node 2
		leaf2ExpectedKeys := []int{8, 7, 6, 5}
		leaf2ExpectedValues := []any{"C", "L", "K", "J"}

		leaf2keys := tree.root.children[0].children[1].keys
		leaf2values := tree.root.children[0].children[1].values

		leaf3 := tree.root.children[0].children[2]
		leaf2Next := tree.root.children[0].children[1].next

		testLeafNode(t, leaf2keys, leaf2ExpectedKeys, leaf2values, leaf2ExpectedValues, leaf2Next, leaf3)

		// leaf node 3
		leaf3ExpectedKeys := []int{9, 10}
		leaf3ExpectedValues := []any{"I", "A"}

		leaf3keys := tree.root.children[0].children[2].keys
		leaf3values := tree.root.children[0].children[2].values

		leaf4 := tree.root.children[1].children[0]
		leaf3Next := tree.root.children[0].children[2].next

		testLeafNode(t, leaf3keys, leaf3ExpectedKeys, leaf3values, leaf3ExpectedValues, leaf3Next, leaf4)

		// leaf node 4
		leaf4ExpectedKeys := []int{11, 12, 13, 14}
		leaf4ExpectedValues := []any{"M", "N", "O", "P"}

		leaf4keys := tree.root.children[1].children[0].keys
		leaf4values := tree.root.children[1].children[0].values

		leaf5 := tree.root.children[1].children[1]
		leaf4Next := tree.root.children[1].children[0].next

		testLeafNode(t, leaf4keys, leaf4ExpectedKeys, leaf4values, leaf4ExpectedValues, leaf4Next, leaf5)

		// leaf node 5
		leaf5ExpectedKeys := []int{15, 20}
		leaf5ExpectedValues := []any{"D", "B"}

		leaf5keys := tree.root.children[1].children[1].keys
		leaf5values := tree.root.children[1].children[1].values

		if leaf5.next != nil {
			t.Errorf("last leaf should not have next.")
		}

		testLeafNode(t, leaf5keys, leaf5ExpectedKeys, leaf5values, leaf5ExpectedValues, nil, nil)
	})
}

func testLeafNode(t *testing.T, keys, expectedKeys []int, values, expectedValues []any, next *Node, expectedNext *Node) {
	for idx, key := range keys {
		if !Contains(expectedKeys, key) {
			t.Errorf("Could not find key %v.", key)
		}
		got := values[idx]
		want := expectedValues[idx]
		if got != want {
			t.Errorf("Key missmatched value. Got %v, want %v", got, want)
		}
	}

	// last node想定の場合はnextの一致確認をしない
	if expectedNext == nil {
		return
	}

	if expectedNext != next {
		t.Errorf("next node is not set correctly. expectedNext: %p next:%p", expectedNext, next)
	}
}

func Test_Search(t *testing.T) {
	tree := NewBPlusTree()

	tree.Insert(10, "A")
	tree.Insert(20, "B")
	tree.Insert(5, "C")
	tree.Insert(15, "D")
	tree.Insert(25, "E")

	expectedList := []struct {
		key   int
		value any
	}{
		{10, "A"},
		{20, "B"},
		{5, "C"},
		{15, "D"},
		{25, "E"},
	}

	for _, expected := range expectedList {
		value, found := tree.Search(expected.key)
		if !found {
			t.Errorf("Could not search by key: %d.", expected.key)
		}
		if value != expected.value {
			t.Errorf("Value is not set coorectly. got=%v want=%v", expected.value, value)
		}
	}
}

func Test_Delete(t *testing.T) {
	t.Run("Delete node. root is leafnode.", func(t *testing.T) {
		tree := NewBPlusTree()

		tree.Insert(10, "A")
		tree.Insert(20, "B")
		tree.Insert(5, "C")
		tree.Insert(15, "D")

		keyExpected := []int{10, 20}
		valueExpected := []any{"A", "B"}

		tree.Delete(15)
		tree.Delete(5)

		keys := tree.root.keys
		values := tree.root.values

		if len(tree.root.keys) != 2 {
			t.Errorf("delete node went wrong.")
		}

		testLeafNode(t, keys, keyExpected, values, valueExpected, nil, nil)
	})

	t.Run("Delete Node Child And Merge. Raise root", func(t *testing.T) {
		tree := NewBPlusTree()

		tree.Insert(10, "A")
		tree.Insert(20, "B")
		tree.Insert(5, "C")
		tree.Insert(15, "D")
		tree.Insert(1, "E")

		// 1, 5, 10, 15, 20
		// key = 15
		// child1 1, 5, 10 next is exist
		// child2 15, 20 next is nil

		tree.Delete(1)
		tree.Delete(5)

		// 10, 15, 20

		keyExpected := []int{10, 15, 20}
		valueExpected := []any{"A", "D", "B"}

		keys := tree.root.keys
		values := tree.root.values

		if len(keys) != 3 {
			t.Errorf("delete node went wrong.")
		}

		testLeafNode(t, keys, keyExpected, values, valueExpected, nil, nil)
	})

	t.Run("Root remove", func(t *testing.T) {
		tree := NewBPlusTree()

		tree.Insert(10, "A")
		tree.Delete(10)

		if tree.root != nil {
			t.Errorf("delete one root node went wrong.")
		}

		got := tree.Delete(10)
		if got {
			t.Errorf("delete nothing went wrong.")
		}

		tree = NewBPlusTree()

		tree.Insert(10, "A")
		tree.Insert(20, "B")
		tree.Insert(5, "C")
		tree.Insert(15, "D")
		tree.Insert(1, "E")

		// 1, 5, 10, 15, 20
		// root = 15
		// child1 1, 5, 10 next is exist
		// child2 15, 20 next is nil

		tree.Delete(15)
		
		// 1, 5, 10, 20
		keyExpected := []int{1, 5, 10, 20}
		valueExpected := []any{"E", "C", "A", "B"}

		keys := tree.root.keys
		values := tree.root.values

		testLeafNode(t, keys, keyExpected, values, valueExpected, nil, nil)
	})

	t.Run("Change Parent Right Leaf Merge", func(t *testing.T) {
		tree := NewBPlusTree()

		tree.Insert(10, "A")
		tree.Insert(20, "B")
		tree.Insert(5, "C")
		tree.Insert(15, "D")
		tree.Insert(1, "E")
		tree.Insert(2, "J")
		tree.Insert(3, "F")
		tree.Insert(4, "G")
		tree.Insert(7, "H")
		tree.Insert(6, "I")
		
		tree.Delete(5)

		if len(tree.root.keys) != 2 || tree.root.keys[0] != 6 || tree.root.keys[1] != 15 {
			t.Errorf("delete root node or change parent went wrong.")
		}

		tree.Delete(1)
		tree.Delete(2)
		tree.Delete(3)

		if len(tree.root.keys) != 2 || tree.root.keys[0] != 7 || tree.root.keys[1] != 15 {
			t.Errorf("change parent went wrong.")
		}

		keyExpected := []int{4, 6}
		valueExpected := []any{"G", "I"}

		keys := tree.root.children[0].keys
		values := tree.root.children[0].values

		leaf1Next := tree.root.children[0].next
		leaf2 := tree.root.children[1]

		testLeafNode(t, keys, keyExpected, values, valueExpected, leaf1Next, leaf2)

		keyExpected = []int{7, 10}
		valueExpected = []any{"H", "A"}

		keys = tree.root.children[1].keys
		values = tree.root.children[1].values

		leaf2Next := tree.root.children[1].next
		leaf3 := tree.root.children[2]

		testLeafNode(t, keys, keyExpected, values, valueExpected, leaf2Next, leaf3)

		keyExpected = []int{15, 20}
		valueExpected = []any{"D", "B"}

		keys = tree.root.children[2].keys
		values = tree.root.children[2].values

		if tree.root.children[2].next != nil {
			t.Errorf("last leaf node should not has next leaf pointer.")
		}

		testLeafNode(t, keys, keyExpected, values, valueExpected, nil, nil)
	})

	t.Run("Change Parent Left Leaf Merge", func(t *testing.T) {
		tree := NewBPlusTree()

		tree.Insert(10, "A")
		tree.Insert(20, "B")
		tree.Insert(5, "C")
		tree.Insert(15, "D")
		tree.Insert(1, "E")
		tree.Insert(2, "J")
		tree.Insert(3, "F")
		tree.Insert(4, "G")
		tree.Insert(7, "H")
		tree.Insert(6, "I")

		tree.Delete(5)

		if len(tree.root.keys) != 2 || tree.root.keys[0] != 6 || tree.root.keys[1] != 15 {
			t.Errorf("delete root node or change parent went wrong.")
		}

		tree.Delete(6)

		if len(tree.root.keys) != 2 || tree.root.keys[0] != 7 || tree.root.keys[1] != 15 {
			t.Errorf("delete root node or change parent went wrong.")
		}

		// 内部ノード削除とリーフノード削除のマージ
		tree.Delete(7)

		if len(tree.root.keys) != 2 || tree.root.keys[0] != 3 || tree.root.keys[1] != 15 {
			t.Errorf("delete root node or change parent went wrong.")
		}

		// merge発生
		keyExpected := []int{1, 2}
		valueExpected := []any{"E", "J"}

		keys := tree.root.children[0].keys
		values := tree.root.children[0].values

		leaf1Next := tree.root.children[0].next
		leaf2 := tree.root.children[1]

		testLeafNode(t, keys, keyExpected, values, valueExpected, leaf1Next, leaf2)

		keyExpected = []int{3, 4, 10}
		valueExpected = []any{"F", "G", "A"}

		keys = tree.root.children[1].keys
		values = tree.root.children[1].values

		leaf2Next := tree.root.children[1].next
		leaf3 := tree.root.children[2]


		testLeafNode(t, keys, keyExpected, values, valueExpected, leaf2Next, leaf3)

		keyExpected = []int{15, 20}
		valueExpected = []any{"D", "B"}

		keys = tree.root.children[2].keys
		values = tree.root.children[2].values

		leaf2Next = tree.root.children[2].next
		if leaf2Next != nil {
			t.Errorf("last leaf node next should be null")
		}

		testLeafNode(t, keys, keyExpected, values, valueExpected, nil, nil)

		// リーフノードの削除によるマージ
		tree.Delete(20)

		if len(tree.root.keys) != 2 || tree.root.keys[0] != 3 || tree.root.keys[1] != 10 {
			t.Errorf("delete root node or change parent went wrong.")
		}

		keyExpected = []int{1, 2}
		valueExpected = []any{"E", "J"}

		keys = tree.root.children[0].keys
		values = tree.root.children[0].values

		leaf1Next = tree.root.children[0].next
		leaf2 = tree.root.children[1]

		testLeafNode(t, keys, keyExpected, values, valueExpected, leaf1Next, leaf2)

		keyExpected = []int{3, 4}
		valueExpected = []any{"F", "G"}

		keys = tree.root.children[1].keys
		values = tree.root.children[1].values

		leaf2Next = tree.root.children[1].next
		leaf3 = tree.root.children[2]


		testLeafNode(t, keys, keyExpected, values, valueExpected, leaf2Next, leaf3)

		keyExpected = []int{10, 15}
		valueExpected = []any{"A", "D"}

		keys = tree.root.children[2].keys
		values = tree.root.children[2].values

		leaf2Next = tree.root.children[2].next
		if leaf2Next != nil {
			t.Errorf("last leaf node next should be null")
		}

		testLeafNode(t, keys, keyExpected, values, valueExpected, nil, nil)
	})

	t.Run("Delete Deep Node", func(t *testing.T) {
		// test tree（目視で問題ないこと一通り確認。11のみ実装）
		// test 11 root
		// test 5 inner
		// test 1 left leaf
		// test 7 center left leaf
		// test 9 inner
		// test 10 center right leaf
		// test 13 right leaf
		// test 15 inner
		// test 20 last leaf

		// 11
		// child: 5, 9
		// 5----
		// 1, 2, 3, 4
		// 5, 6, 7, 8
		// ----5
		// 9----
		// 9, 10
		// ----9
		// child: 15
		// 11, 12, 13, 14 <- set reaf node
		// 15, 20

		t.Run("delete root node", func(t *testing.T) {
			tree := NewBPlusTree()

			tree.Insert(10, "A")
			tree.Insert(20, "B")
			tree.Insert(5, "C")
			tree.Insert(15, "D")
			tree.Insert(1, "E")
			tree.Insert(2, "F")
			tree.Insert(3, "G")
			tree.Insert(4, "H")
			tree.Insert(9, "I")
			tree.Insert(8, "J")
			tree.Insert(7, "K")
			tree.Insert(6, "L")
			tree.Insert(11, "M")
			tree.Insert(12, "N")
			tree.Insert(13, "O")
			tree.Insert(14, "P")

			tree.Delete(11)

			// leaf node 1
			leaf1ExpectedKeys := []int{1, 2, 3, 4}
			leaf1ExpectedValues := []any{"E", "F", "G", "H"}
	
			leaf1keys := tree.root.children[0].keys
			leaf1values := tree.root.children[0].values
	
			leaf2 := tree.root.children[1]
			leaf1Next := tree.root.children[0].next
	
			testLeafNode(t, leaf1keys, leaf1ExpectedKeys, leaf1values, leaf1ExpectedValues, leaf1Next, leaf2)
	
			// leaf node 2
			leaf2ExpectedKeys := []int{8, 7, 6, 5}
			leaf2ExpectedValues := []any{"C", "L", "K", "J"}
	
			leaf2keys := tree.root.children[1].keys
			leaf2values := tree.root.children[1].values
	
			leaf3 := tree.root.children[2]
			leaf2Next := tree.root.children[1].next
	
			testLeafNode(t, leaf2keys, leaf2ExpectedKeys, leaf2values, leaf2ExpectedValues, leaf2Next, leaf3)
	
			// leaf node 3
			leaf3ExpectedKeys := []int{9, 10}
			leaf3ExpectedValues := []any{"I", "A"}
	
			leaf3keys := tree.root.children[2].keys
			leaf3values := tree.root.children[2].values
	
			leaf4 := tree.root.children[3]
			leaf3Next := tree.root.children[2].next
	
			testLeafNode(t, leaf3keys, leaf3ExpectedKeys, leaf3values, leaf3ExpectedValues, leaf3Next, leaf4)
	
			// leaf node 4
			leaf4ExpectedKeys := []int{12, 13, 14}
			leaf4ExpectedValues := []any{"N", "O", "P"}
	
			leaf4keys := tree.root.children[3].keys
			leaf4values := tree.root.children[3].values
	
			leaf5 := tree.root.children[4]
			leaf4Next := tree.root.children[3].next
	
			testLeafNode(t, leaf4keys, leaf4ExpectedKeys, leaf4values, leaf4ExpectedValues, leaf4Next, leaf5)
	
			// leaf node 5
			leaf5ExpectedKeys := []int{15, 20}
			leaf5ExpectedValues := []any{"D", "B"}
	
			leaf5keys := tree.root.children[4].keys
			leaf5values := tree.root.children[4].values
	
			if leaf5.next != nil {
				t.Errorf("last leaf should not have next.")
			}
	
			testLeafNode(t, leaf5keys, leaf5ExpectedKeys, leaf5values, leaf5ExpectedValues, nil, nil)
		})

		// 目視確認用
		// t.Logf("%v", tree.root)
		// t.Logf("%v", tree.root.children[0])
		// t.Logf("%v", tree.root.children[0].children[0])
		// t.Logf("%v", tree.root.children[0].children[1])
		// t.Logf("%v", tree.root.children[0].children[2])
		// t.Logf("%v", tree.root.children[1])
		// t.Logf("%v", tree.root.children[1].children[0])
		// t.Logf("%v", tree.root.children[1].children[1])
		// t.Logf("%v", tree.root.children[2])
		// t.Logf("%v", tree.root.children[3])
		// t.Logf("%v", tree.root.children[4])
	})
}