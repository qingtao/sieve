package filter

import (
	"strings"
	"unicode"
)

const (
	// 通配符
	symbolStar = '*'
)

// 节点
type node struct {
	// 是否结束
	IsEnd bool
	// 标签
	Tag uint8
	// 替换
	AutoReplace bool
	// 联想字符
	Children map[rune]*node
}

// 添加关键词
func (root *node) AddWord(word string, tag uint8, autoReplace bool) bool {
	word = strings.TrimSpace(word)
	if len(word) == 0 {
		return true
	}

	node := root
	var i, j int
	var w rune
	for i, w = range word {
		// 首字符为#是注释，忽略
		if i == 0 && w == '#' {
			return true
		}

		w = trans(w)
		// 首字符不能以符号或者通配符开始
		if i == 0 && (w < 0 || w == symbolStar) {
			break
		}

		if w > 0 {
			node = node.addChild(w)
			j++
		}
	}

	// 禁止单个字+符号的形式，导致单个字被错杀
	if i > 1 && j < 2 && len([]rune(word)) > 1 {
		root.RemoveWord(word)
		return false
	}

	// 非根节点才修改，防止无效关键词修改根节点
	if node != root {
		node.IsEnd = true
		node.Tag = tag
		node.AutoReplace = autoReplace
		return true
	}

	return false
}

// 删除关键词
func (root *node) RemoveWord(word string) bool {
	path := []rune(word)
	ptrs := make([]*node, len(path))
	node := root

	ok := false
	// 正向检验关键词是否存在
	for i, w := range path {
		node, ok = node.Children[w]
		if !ok {
			return false
		}
		ptrs[i] = node
	}

	node.IsEnd = false
	for i := len(path) - 1; i > 0; i-- {
		if ptrs[i].IsEnd || len(ptrs[i].Children) > 0 {
			break
		}
		delete(ptrs[i-1].Children, path[i])
	}

	return true
}

// 添加单个字符
func (n *node) addChild(w rune) *node {
	if n.Children == nil {
		n.Children = make(map[rune]*node)
	} else {
		child, ok := n.Children[w]
		if ok {
			return child
		}
	}

	child := &node{}
	n.Children[w] = child
	return child
}

// 获取子字符节点
func (n *node) getChild(w rune) *node {
	child, ok := n.Children[w]
	if ok {
		return child
	}
	return nil
}

func (root *node) Search(ws []rune) (start int, end int, tag uint8, autoReplace bool) {
	if len(ws) == 0 {
		return
	}

	node := root
	start = -1

	length := len(ws)
	for i := 0; i < length; i++ {
		w := trans(ws[i])
		if w <= 0 {
			continue
		}
		// fmt.Println("当前字符是:", i, string(w))

		// 查询是否存在该字符，如果不存在尝试查找通配符
		temp := node.getChild(w)
		if node != root && temp == nil {
			node = node.getChild(symbolStar)
		} else {
			node = temp
		}

		// 举例 「苹果」和「苹果**本」是关键词
		if node == nil {
			// 苹果笔记
			if end > 0 {
				break
			}
			// 当前未匹配，回退到根节点
			node = root

			// 苹方
			if start >= 0 {
				start = -1
				// 当前字符可能不是关键字的中间字符，但是可能是起始字符，重新判断
				i--
			}
		} else {
			// 苹
			if start == -1 {
				start = i
			}

			// 苹果
			if node.IsEnd {
				end = i
				tag = node.Tag
				autoReplace = node.AutoReplace
				if len(node.Children) == 0 {
					break
				}
			}
		}

	}

	// 匹配成功时，适配数组左开右闭把end+1
	if end == 0 {
		start = 0
		end = 0
	} else {
		end += 1
	}

	// fmt.Println("index end:", string(ws), start, end)
	return
}

func trans(w rune) rune {
	if w > 255 {
		// 判断是否为符号
		if unicode.IsPunct(w) {
			return -1
		}
		// 其余文字
		return w
	}

	if w == symbolStar || (w >= 'a' && w <= 'z') || (w >= '0' && w <= '9') {
		return w
	}

	if w >= 'A' && w <= 'Z' {
		return w + 32
	}

	return -1
}
