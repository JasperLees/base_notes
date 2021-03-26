#### 快速排序

##### 解题思路

- 1：以最后一个结点为`target`;
- 2：从左边开始比较，往右走，到第一个比它大的数`numsleft]`
- 3：再从右边开始比较，往左走，找到第一个比他小的数`numsright]`
- 4：交换左边 `nums[left]、nums[right]`;
- 5: 重复2、3、4步，`left`和`right`重合，这是确定`target`在数组的最终位置；
- 6：分别将`target`左边数组和右边数组，递归去重复1~5步，最终获得排序

```go
func quickSort(nums []int, start, end int) {
	if start >= end {
		return
	}
	target := nums[end] // 去start 先找右边，取end先找左边
	left, right := start, end
	for left < right {
		// 找到从前往后的第一个比 target 大的 index
		for left < right && nums[left] <= target {
			left++
		}
		// 找到从后往前的第一个比 target 小的 index
		for left < right && nums[right] >= target {
			right--
		}
		// 交互
		nums[left], nums[right] = nums[right], nums[left]
	}
	// 找到target最终放的位置，交互target(end位) 和 nums[left]
	nums[end], nums[left] = nums[left], target
	// 递归排序 target 左右两边的数组
	quickSort(nums, start, left-1)
	quickSort(nums, left+1, end)
}

func main() {
	nums := []int{5, 2, 7, 4, 6, 3, 8, 1, 0, 9}
	quickSort(nums, 0, len(nums)-1)
	fmt.Println(nums)
}
```


#### 归并排序

##### 解题思路

- 1: 先将数据从中间位置开始拆分成`n`个数组
- 2: 递归的将数组两两合并
  - 2-1：因为单个数组是有序的，比较两个数组首个元素，小的先`push`，大的后`push`
  - 2-2: 将 `left` 或 `right` 剩余部分追加到后面
 
##### 类似题型

- 两个有序链表合并

```go
func mergeSort(arr []int) []int {
	if len(arr) < 2 {
		return arr
	}
	mid := len(arr) / 2
	return merge(mergeSort(arr[:mid]), mergeSort(arr[mid:]))
}

func merge(left, right []int) []int {
	arr := make([]int, 0)
	for len(left) > 0 && len(right) > 0 {
		if left[0] <= right[0] {
			arr = append(arr, left[0])
			left = left[1:]
		} else {
			arr = append(arr, right[0])
			right = right[1:]
		}
	}
	if len(left) > 0 {
		arr = append(arr, left...)
	}
	if len(right) > 0 {
		arr = append(arr, right...)
	}
	return arr
}

func main() {
	arr := []int{6, 3, 4, 7, 1, 5, 9, 8, 0}
	fmt.Println(mergeSort(arr))
}
```
