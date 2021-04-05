package main

import (
	"fmt"
	"math/rand"
	"time"
)

// 如何从一个数组输出随机数组（洗牌算法）
func shuffleCards(nums []int) {
	for i := len(nums) - 1; i >= 0; i-- {
		idx := rand.Intn(i + 1)
		if idx == i {
			continue
		}
		nums[idx], nums[i] = nums[i], nums[idx]
	}
}

// 如何随机生成不重复的 10个100 以内的数字
func rand10Nums() []int {
	nums := []int{}
	for i := 1; i <= 100; i++ {
		nums = append(nums, i)
	}
	for i := len(nums) - 1; i >= 100-10; i-- {
		idx := rand.Intn(i + 1)
		if idx == i {
			continue
		}
		nums[idx], nums[i] = nums[i], nums[idx]
	}
	return nums[90:]
}

func main() {
	nums := []int{}
	for i := 1; i <= 54; i++ {
		nums = append(nums, i)
	}
	rand.Seed(time.Now().UnixNano())
	shuffleCards(nums)

	data := make(map[int]int)
	for i := 0; i < len(nums); i++ {
		if _, ok := data[nums[i]]; ok {
			panic("数据错误")
		}
		data[nums[i]]++
	}

	fmt.Println(nums)
	fmt.Println(rand10Nums())
}
