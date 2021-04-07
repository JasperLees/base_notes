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

#### 堆排序

##### 定义

- 大顶堆定义：
  - 每个结点的值都大于或等于其左右孩子结点的值
  - `nums[i] >= nums[2i+1] && nums[i] >= nums[2i+2]`
- 小顶堆定义：
  - 每个结点的值都小于或等于其左右孩子结点的值
  - `nums[i] <= nums[2i+1] && nums[i] <= nums[2i+2]`
  
##### 解题思想

- 1:将待排序序列构造一个大顶堆，此时整个数组的最大值就是堆顶的根结点
- 2:将其与末尾元素进行交换，此时末尾为最大值
- 3:将剩下`n-1`个元素重新构建一个大顶堆，得到次小元素
- 4:重复`n-1`次，得到排序数组

##### 时空复杂度

- `Time:O(nlogn)`
- `Space:O(1)`

```go
func heapSort(nums []int) {
	// 构建大顶堆，从第一个非叶子结点(len(nums)/2 - 1)从下至上，从左至右调整结构
	for i := len(nums)/2 - 1; i >= 0; i-- {
		adjustHeap(nums, i, len(nums))
	}
	// 调整堆结构+交换堆顶元素与末尾元素
	for i := len(nums) - 1; i > 0; i-- {
		nums[0], nums[i] = nums[i], nums[0]
		adjustHeap(nums, 0, i)
	}
}

// 调整大顶堆
func adjustHeap(nums []int, start, end int) {
	target := nums[start]
	// 从start结点的左子结点开始(2start+1)
	for i := start*2 + 1; i < end; i = i*2 + 1 {
		// 如果左子结点小于右结点，k指向右子结点(目标是找出最大值)
		if i+1 < end && nums[i] < nums[i+1] {
			i++
		}
		// 如果子结点大于父结点，将子结点的值赋值给父结点(父结点的值已经保留)
		if nums[i] > target {
			nums[start] = nums[i]
			start = i
		} else {
			break
		}
	}
	nums[start] = target
}

func main() {
	nums := []int{4, 9, 0, 5, 6, 3, 7, 2, 8, 1}
	heapSort(nums)
	fmt.Println(nums)
}
```

#### 冒泡排序

##### 解题思路

- 1：第一次循环，选择数组`(0~n-1)`最大的元素放在最后一位`(n-1)`
- 2: 第`k`次循环，选择数组`(0~n-k)`最大的元素放到`n-k`位

##### 时空复杂度

- `Time:O(n^2)`
  - 在循环中加个`flag=false`，如果在`j`的排序改变值则`flag=true`;
  - 循环一次判断`flag == false`，没有改变数据，结束循环，能达到最好的时间复杂度是`O(n)` 
- `Space:O(1)`

```go
func bubbleSort(nums []int) {
	for i := 0; i < len(nums); i++ {
		for j := 0; j < len(nums)-i-1; j++ {
			if nums[j] > nums[j+1] {
				nums[j], nums[j+1] = nums[j+1], nums[j]
			}
		}
		fmt.Println(nums)
	}
}

func main() {
	nums := []int{4, 9, 0, 5, 6, 3, 7, 2, 8, 1}
	bubbleSort(nums)
	fmt.Println(nums)
}
```

#### 选择排序

##### 解题思路

- 1: 第一次遍历，从`(1~n-1)`中找到最小的元素，放到第`0`位
- 2: 第`k`次遍历，从`(k+1~n-1)`中找到最小的元素，放到第`k`位

##### 时空复杂度

- `Time:O(n^2)`
- `Space:O(1)`

```go
func selectSort(nums []int)  {
	for i := 0; i < len(nums); i++ {
		min := nums[i]
		for j := i; j < len(nums); j++ {
			if min > nums[j] {
				min, nums[j] = nums[j], min
			}
		}
		nums[i] = min
	}
}

func main() {
	nums := []int{4, 9, 0, 5, 6, 3, 7, 2, 8, 1}
	selectSort(nums)
	fmt.Println(nums)
}
```

#### 计数排序

##### 使用场景

- 最快的排序，比快速排序`O(nlogn)`更快
- 适用于一定范围的整数排序(例如：取值是0~9)：
  - 不适用于最大最小值差距过大的数组
  - 不适用于不是整数的数组
- 利用数组下标来确定元素的正确位置

##### 解题思路

- 1:先找到最大最小值，以确定计数数组`countArr`的长度
- 2:遍历`nums`，统计个数值出现的次数
- 3:遍历`countArr`，将其中的出现次数，累计变更为，最后一次出现的下标
- 4:逆序遍历`nums`，取出 `nums[i]` 对于的位置，放进结果数组

##### 时空复杂度:

- `Time:O(n+k)`
- `Space:O(k)`

```go
func countSort(nums []int) []int {
	if len(nums) < 2 {
		return nums
	}
	// 找到最大最小值，以觉得计数数组的长度 Time:O(n)
	max, min := nums[0], nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] > max {
			max = nums[i]
		}
		if nums[i] < min {
			min = nums[i]
		}
	}
	// 计数，统计数组是 countArr[数值]=出现次数 Time:O(n) Space:(k)
	countArr := make([]int, max-min+1)
	for i := 0; i < len(nums); i++ {
		countArr[nums[i]-min]++
	}
	// 相加，是统计数组变成 countArr[数值]=排序后的下标，Time:O(k)
	var sum int
	for i := 0; i < len(countArr); i++ {
		sum += countArr[i]
		countArr[i] = sum
	}
	// 从后往前遍历，存放有序的元素 Time:O(n)
	results := make([]int, len(nums))
	for i := len(nums)-1; i >= 0; i-- {
		val := nums[i]-min
		results[countArr[val]-1] = nums[i] // 计算从1开始，所以要-1
		countArr[val]--
	}
	return results
}

func main() {
	nums := []int{4, 9, 0, 5, 6, 3, 7, 2, 8, 1}
	nums = countSort(nums)
	fmt.Println(nums)
}
```

#### 桶排序

##### 使用场景

- 当数列取值范围过大，或者不是整数是，不能使用计数排序，可以使用桶排序

##### 解题思路

##### 时空复杂度

```go
func bucketSort(nums []float64) []float64 {
	if len(nums) < 2 {
		return nums
	}
	// 1.获取数组的max和min，计算出统计数组长度k
	max, min := nums[0], nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] > max {
			max = nums[i]
		}
		if nums[i] < min {
			min = nums[i]
		}
	}
	// 2. 初始化桶
	bucketNum := len(nums)
	bucketArr := make([][]float64, bucketNum)
	for i := 0; i < bucketNum; i++ {
		bucketArr[i] = make([]float64, 0)
	}
	// 3. 遍历原始数组，将每个元素放入桶中
	for i := 0; i < len(nums); i++ {
		num := (nums[i]-min) * float64(bucketNum - 1) / (max-min)
		bucketArr[int(num)] = append(bucketArr[int(num)], nums[i])
	}
	// 4. 对每个桶内部进行排序
	for i := 0; i < len(bucketArr); i++ {
		sort.Float64s(bucketArr[i])
	}
	// 5. 输出元素
	results := make([]float64, 0)
	for i := 0; i < len(bucketArr); i++ {
		results = append(results, bucketArr[i]...)
	}
	return results
}

func main() {
	nums := []float64{4.1, 9.2, 0.4, 0.8, 5.1, 6.1, 3.1, 7.1, 2.1, 8.1, 1.1}
	nums = bucketSort(nums)
	fmt.Println(nums)
}
```
