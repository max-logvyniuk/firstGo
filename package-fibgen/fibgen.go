package fibgen

// func Fibgen() {
// 	jobs := make(chan int, 100)
// 	results := make(chan int, 100)

// 	go Worker(jobs, results)
// 	go Worker(jobs, results)

// 	for i := 0; i < 100; i++ {
// 		jobs <- i
// 	}
// 	close(jobs)

// 	for j := 0; j < 100; j++ {
// 		fmt.Println(<-results)
// 	}
// }

func Worker(jobs <-chan int, results chan<- int) {
	for n := range jobs {
		results <- fib(n)
	}
}

func fib(n int) int {
	if n <= 1 {
		return n
	}

	return fib(n-1) + fib(n-2)
}
