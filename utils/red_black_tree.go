package utils

type TreeNode struct {
	ID       int
	ParentID int
	Name     string
	Left     *TreeNode
	Right    *TreeNode
	red      bool
}

// 红黑树的插入操作
func Insert(root *TreeNode, id, parentID int, name string) *TreeNode {
	if root == nil {
		return &TreeNode{ID: id, ParentID: parentID, Name: name, red: true}
	}

	if id < root.ID {
		root.Left = Insert(root.Left, id, parentID, name)
	} else if id > root.ID {
		root.Right = Insert(root.Right, id, parentID, name)
	} else {
		// Handle duplicate ID (if any)
	}

	// 检查并修复红黑树的性质
	root = fixUp(root)

	return root
}

// 红黑树的修复操作
func fixUp(node *TreeNode) *TreeNode {
	if isRed(node.Right) && !isRed(node.Left) {
		node = rotateLeft(node)
	}
	if isRed(node.Left) && isRed(node.Left.Left) {
		node = rotateRight(node)
	}
	if isRed(node.Left) && isRed(node.Right) {
		colorFlip(node)
	}

	return node
}

// 左旋转操作
func rotateLeft(node *TreeNode) *TreeNode {
	x := node.Right
	node.Right = x.Left
	x.Left = node
	x.red = node.red
	node.red = true
	return x
}

// 右旋转操作
func rotateRight(node *TreeNode) *TreeNode {
	x := node.Left
	node.Left = x.Right
	x.Right = node
	x.red = node.red
	node.red = true
	return x
}

// 颜色翻转操作
func colorFlip(node *TreeNode) {
	node.red = !node.red
	node.Left.red = !node.Left.red
	node.Right.red = !node.Right.red
}

// 判断节点是否为红色
func isRed(node *TreeNode) bool {
	if node == nil {
		return false
	}
	return node.red
}

// 通过 ID 查找父 ID
func FindParentID(root *TreeNode, id int) int {
	for root != nil {
		if id == root.ID {
			return root.ParentID
		} else if id < root.ID {
			root = root.Left
		} else {
			root = root.Right
		}
	}
	return -1 // 如果找不到对应的 ID，则返回 -1
}
