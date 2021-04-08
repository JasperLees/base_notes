## 时间复杂度

|            |   Best    |  Average  |   Worst   | Worst(Space) | 稳定性 |
| ---------- | :-------: | :-------: | :-------: | :----------: | :----: |
| QuickSort  | O(n logn) | O(n logn) |  O(n^2)   |   O(logn)    | 不稳定 |
| MergeSort  | O(n log)  | O(n log)  | O(n log)  |     O(n)     |  稳定  |
| HeapSort   | O(n logn) | O(n logn) | O(n logn) |     O(1)     | 不稳定 |
| BubbleSort |   O(n)    |  O(n^2)   |  O(n^2)   |     O(1)     |  稳定  |
| SelectSort |  O(n^2)   |  O(n^2)   |  O(n^2)   |     O(1)     | 不稳定 |
| CountSort  |  O(n+k)   |  O(n+k)   |  O(n+k)   |     O(k)     |  稳定  |
| BucketSort |  O(n+k)   |  O(n+k)   |  O(n^2)   |     O(n)     |  稳定  |
| RadixSort  |  O(nk)    |  O(nk)    |  O(nk)    |    O(n+k)    |  稳定  |
| InsertSort |  O(n)     |  O(n^2)   |  O(n^2)   |    O(1)      |  稳定  |
| ShellSort  |  O(nlogn) |  O(n(logn)^2)   |  O(n(logn)^2)   |    O(1)      |  不稳定  |

##### 稳定性说明

- 两个值相同的元素在排序前后是否有顺序的变化

- 假如 `a`原本在`b` 前面，且`a == b`
- 稳定：排序后`a`仍然在`b`前面
- 不稳定：排序后`a`可以出现在`b`后面
- 选择排序的不稳定原因：其根本逻辑，是选择后面的一个最小数和当前数交换
  - 如果当前元素比一个元素小，且该小元素在和当前相等的元素后面，交换后稳定性被破坏
  -  数组`[A1, A2, B1](A1=A2>B1)`交换后变成 `[B1, A2, A1]`，破坏稳定性

##### 稳定的排序

- 归并排序、冒泡排序、插入排序、计算排序、桶排序、基数排序

##### 不稳定的排序

- 快速排序、堆排序、选择排序、希尔排序

## 排序算法

#### 快速排序

##### 解题思路

- 1: 选择目标(通常最后一个结点)
- 2: 通过从前往后比较，找到第一个比`target`大的元素；从后往前比较，找到第一个比`target`小的元素
- 3：然后交换两个元素，重复第2步，直到找到target在数组的最终位置(left)
- 4：用`(start, left-1)`和`(left+1, end)`去递归，重复1、2、3步直到数组有序

##### 时空复杂度

- `Time`：`O(nlogn)`[最坏：`O(n^2)`，平均每次遍历+移动的元素都是`N`]
- `Space`：`O(log n)`递归栈

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
		// 交换
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

------

#### 归并排序

##### 解题思路

- 1: 先将数据从中间位置开始拆分成`n`个数组
- 2: 递归的将数组两两合并
    - 2-1：因为单个数组是有序的，比较两个数组首个元素，小的先`push`，大的后`push`
    - 2-2: 将 `left` 或 `right` 剩余部分追加到后面

##### 类似题型

- 两个有序链表合并

##### 时空复杂度

- `Time`：`O(n logn)`
- `Space`：`O(n)`

```go
func mergeSort(nums []int) []int {
	if len(nums) < 2 {
		return nums
	}
	mid := len(nums) / 2
	return merge(mergeSort(nums[mid:]), mergeSort(nums[:mid]))
}

func merge(left, right []int) []int {
	nums := []int{}
	for len(left) > 0 && len(right) > 0 {
		if left[0] <= right[0] {
			nums = append(nums, left[0])
			left = left[1:]
		} else {
			nums = append(nums, right[0])
			right = right[1:]
		}
	}
	if len(left) > 0 {
		nums = append(nums, left...)
	}
	if len(right) > 0 {
		nums = append(nums, right...)
	}
	return nums
}

func main() {
	nums := []int{5, 2, 7, 4, 6, 3, 8, 1, 0, 9}
	fmt.Println(mergeSort(nums))
}
```

------

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

------

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
	var flag bool
	for i := 0; i < len(nums); i++ {
		flag = false
		for j := 0; j < len(nums)-i-1; j++ {
			if nums[j] > nums[j+1] {
				nums[j], nums[j+1] = nums[j+1], nums[j]
				flag = true
			}
		}
		if !flag {
			break
		}
	}
}
```

------

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

------

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

------

#### 桶排序

##### 使用场景

- 当数列取值范围过大，或者不是整数是，不能使用计数排序，可以使用桶排序

##### 解题思路

- 1.获取数组的`max`和`min`，每个桶的区间 = 两者差 / 桶的个数(数组长度)
- 2.初始化桶，二位数组
- 3.遍历原始数组，将每个元素放入桶中
- 4.对每个桶内部进行排序
- 5.将元素数据放回数组

##### 时空复杂度

- `Time:O(n+k)`
- `Space：O(n)`，所有桶的元素`==len(nums)`

```go
func bucketSort(nums []float64) []float64 {
	if len(nums) < 2 {
		return nums
	}
	// 获取 max 和 min Time:O(n)
	max, min := nums[0], nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] > max {
			max = nums[i]
		}
		if nums[i] < min {
			min = nums[i]
		}
	}
	diff := (max - min) / float64(len(nums)) // 桶区间

	// 初始化桶 Space:O(n)
	bucketArr := make([][]float64, len(nums))

	// 数据入桶 Time:O(n)
	var idx int
	for i := 0; i < len(nums); i++ {
		if idx =int(nums[i]/diff) - 1; idx < 0 {
			idx = 0
		}
		bucketArr[idx] = append(bucketArr[idx], nums[i])
	}

	// 单个桶排序 Time： n/k * log(n/k) * k => nlog(n/k)
	for i := 0; i < len(bucketArr); i++ {
		sort.Float64s(bucketArr[i])
	}

	// 输出排序 Time:O(n)
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

------

#### 基数排序

##### 适用场景

- 数值排序：最大最小值间隔比较大，取最大值位数为排序的次数，最小值前面补零
- 字符串排序：取最长的字符串的长度为排序次数，短的字符串在前面补零

##### 解题思路：
- 1:获取最大的值，以确定循环的次数
- 2:创建10个桶(有负数创建20个，字符串创建128个,`ASCII`码的取值范围)
- 3:从个位开始，以个位为`idx`，放入对应的桶，然后将所有桶按顺序放回 `nums`
- 4:同样使用十位、百位进行相同的操作，直到最高位为止，得到有序数组

##### 时空复杂度

- `Time:O(kn)` 
- `Space:O(n)`

```go
func radixSort(nums []int) []int {
	if len(nums) < 2 {
		return nums
	}
	// 排序的次数由最大值的位数确定，获取最大值 Time:O(n)
	max := nums[0]
	for i := 1; i < len(nums); i++ {
		if nums[i] > max {
			max = nums[i]
		}
	}
	// 对每个位数排序，从个位开始 Time:O(k*n) Space:O(n+k), k 表示一位数组个数，n表示元素个数
	for i := 1; max/i > 0; i = i*10  {
		bucketArr := make([][]int, 10)
		for j := 0; j < len(nums); j++ {
			idx := (nums[j] / i) % 10
			bucketArr[idx] = append(bucketArr[idx], nums[j])
		}
		nums = make([]int, 0)
		for j := 0; j < len(bucketArr); j++ {
			nums = append(nums, bucketArr[j]...)
		}
	}
	return nums
}

func main()  {
	nums := []int{40, 2, 4, 5, 15, 201, 304, 8999, 800, 7, 33}
	nums = radixSort(nums)
	fmt.Println(nums)
}
```

##### 计数排序 `VS` 桶排序 `VS` 基数排序

- 计数排序：单个桶保存单一数值(可能相同)
- 桶排序：  单个桶保存一定范文的数值
- 基数排序：根据数值的每位数字来分配桶

------

#### 插入排序

##### 解题思路

- 1:前半部分是有序的，后半部分是无需的
- 2:`nums[i]` 从后往前比较 [0~i-1]之间的数，找到第一个比他小的数值，插入它的后面

##### 时空复杂度

- `Time:O(n^2)`，最好是 `O(n)`
- `Space:O(1)`

```go
func insertSort(nums []int) {
	for i := 1; i < len(nums); i++ {
		for j := i; j > 0; j-- {
			if nums[j] < nums[j-1] {
				nums[j], nums[j-1] = nums[j-1], nums[j]
			} else {
				break
			}
		}
	}
}

func main() {
	nums := []int{4, 9, 0, 5, 6, 3, 7, 2, 8, 1}
	insertSort(nums)
	fmt.Println(nums)
}
```

------

#### 希尔排序

##### 解题思路

- 1:希尔排序是插入排序的一个优化，普通插入排序增量为1
- 2:第一次，以增量为 `len(nums)/2` 的一次插入排序
- 3:第`N`次，以增量为 `len(nums)/(2^N)` 的一次插入排序
- 4:直到增量为1，完成排序

##### 时空复杂度

- `Time:O(n*(logn)^2)`, 最好能达到`O(nlogn)` 
- `Space:O(1)`

```go
func shellSort(nums []int) {
	for d := len(nums) / 2; d > 0; d /= 2 { // O(logn)
		for i := 0; i < len(nums); i++ {    // O(n)
			for j := i - d; j >= 0; j -= d { // O(logn)
				if nums[j] > nums[j+d] {
					nums[j], nums[j+d] = nums[j+d], nums[j]
				} else {
					break
				}
			}
		}
	}
}

func main() {
	nums := []int{4, 9, 0, 5, 6, 3, 7, 2, 8, 1}
	shellSort(nums)
	fmt.Println(nums)
}
```

------

## `Go`的排序

##### `Go`语言的排序主要基于两种排序方法：

- 1：整型(`sort.Ints()`)、浮点型(`sort.float64s()`)、字符串`sort.Strings()`底层都是基于`sort.Sort()`
- 2：`Interface{}`排序的 `sort.Stable()`

##### `sort.Sort()`的逻辑

- 他的内部是组合排序，包括：快速排序、堆排序、希尔排序(插入排序)
- 主要区分逻辑是：
- 1：通过数组的长度计算出快排递归的深度`depth=2*ceil(log(n+1))`
- 2：递归调用`quickSort`时，当指定排序数据长度`1<(end - start)<12`时
  - 使用希尔排序，做一次增量为6的排序和做一次增量为1的排序(插入排序)
- 3：然后，判断递归的深度是否等于0，如果等于0则做堆排序
- 4：否则，`depth--`，根据中位数分开两部分进行递归快速排序

##### `Sort.Stable()`的逻辑

- 他的内部是组合排序，包括：插入排序和归并排序
- 主要逻辑是：
- 1：数组以20个元素为一个块进行划分，将每个块进行插入排序
- 2：每次通过`2^i`个块进行归并排序，`i`表示第`i`次归并，直到`20 * 2^i > n `为止