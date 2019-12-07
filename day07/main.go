package main

import "fmt"

func HeapPermutation(a []int, size int, ch chan []int) {
	if size == 1 {
		ch <- a
		//fmt.Println(a)
	}

	for i := 0; i < size; i++ {
		HeapPermutation(a, size-1, ch)

		if size%2 == 1 {
			a[0], a[size-1] = a[size-1], a[0]
		} else {
			a[i], a[size-1] = a[size-1], a[i]
		}
	}
}

func permutations(arr []int, ch chan []int) {
	var helper func([]int, int)

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			ch <- tmp
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
}

func main() {
	ch := make(chan []int)

	a := []int{0, 1, 2, 3, 4}
	go func() {
		permutations(a, ch)
		ch <- []int{0}
	}()

	for {
		b := <-ch
		if len(b) == 1 {
			break
		}
		fmt.Println(b)
	}
}
