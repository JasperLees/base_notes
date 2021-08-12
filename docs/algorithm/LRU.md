### 简述 LRU 算法及其实现方式

##### 简述含义

- LRU 是一种比较常见的缓存算法，在内存满的时候，选择内存中最近最久未使用的页面予以淘汰
- 算法核心是，如果数据最近被访问过，那么将来被访问的概率会更高

##### 实现方式

- 链表 + 哈希表
- 如果某条数据被访问过，则把这条数据易到队尾，队首则是最近最少使用的数据，淘汰即可

- LRU (本例仅对value进行计算)

```go
type LRUCache struct {
    maxBytes int // 最大容量
    useBytes int // 已使用容量
    root     *ListNode // root结点作为根结点，仅用于指向第一个结点，不保存值
    cache map[string]*ListNode // cache保存node,用于判断key是否存在
}

// 设置值，如果该key在LRU中，则更新值，并将key移动到链表尾部
// 不存在，则直接插入队尾
// 设置完成，要检查缓存是否超过阈值，超过则删除队首元素
func (lc *LRUCache) Set(key string, value interface{})  {
    if node, ok := lc.cache[key]; ok {
        lc.root.MoveToBack(node) // 移动到队尾
        lc.useBytes -= lc.calcBytes(node.Val)
        node.Val = value // 更新值
    } else { // 如果不存在，队尾插入值，cache保存
        node := &ListNode{Key:key, Val: value}
        lc.cache[key] = node
        lc.root.Append(node)
    }   
    lc.useBytes += lc.calcBytes(value)
    // 检查是否超出容量
    lc.checkAndDel()
}

// 设置后，检查是否超过最大值，超过则移动至队首元素
func (lc *LRUCache) checkAndDel()  {
    for lc.maxBytes < lc.useBytes {
        if node := lc.root.Front(); node == nil {
            break
        } else {
            lc.Del(node.Key)
        }
    }
}

// 假设数据为普通的数据类型
func (lc *LRUCache) calcBytes(value interface{}) int {
	size := unsafe.SizeOf(vale)
	return int(size)
}

// 获取LRU值就是，从map中获取值
func (lc *LRUCache) Get(key string) interface{} {
    if node, ok := lc.cache[key]; ok {
        return node.Val
    }
    return nil
}

// 删除值，map和链表都要删除
func (lc *LRUCache) Del(key string) bool {
    if node, ok := lc.cache[key]; ok {
        delete(lc.cache, node.Key)
        lc.useBytes -= len(node.Val)
        lc.root.Remove(node)
        return true
    }
    return false
}

```

- 双向链表实现

```go
type ListNode struct {
	Key string
	Val interface{}
	Prev *ListNode
	Next *ListNode
}

// 将node移动到队尾
func (root *ListNode) MoveToBack(target *ListNode) {
	if target == root.Prev {
		return
	}
	// 将target从原位置是否，即删除
	root.Remove(target)
	// 将target添加至队尾
	root.Append(target)
}

// 增加元素到队尾
func (root *ListNode) Append(target *ListNode) {
	if target == root.Prev {
		return
	}
	// 处理target Prev
	root.Prev.Next = target
	target.Prev = target.Prev
	// 处理target Next
	root.Prev = target
	target.Next = root
}

// 删除元素
func (root *ListNode) Remove(target *ListNode) {
	if root.Prev == root {
		return
	}
	target.Prev.Next = target.Next
	target.Next.Prev = target.Prev
}

// 获取队首
func (root *ListNode) Front() *ListNode {
	if root.Next == root {
		return nil
	}
	return root.Next
}
```
