package tree

const MaxKeys = 4 // 最大キー数（B+木の階層）

type Node struct {
	isLeaf   bool
	keys     []int
	children []*Node
	values   []any // リーフノードで使用
	next     *Node // 連結ポインタ
}

func NewNode(isLeaf bool) *Node {
	return &Node{
		isLeaf:   isLeaf,
		keys:     []int{},
		children: []*Node{},
		values:   []any{},
		next:     nil,
	}
}
