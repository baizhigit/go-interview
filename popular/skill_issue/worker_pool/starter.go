package main

// func say(id int, phrase string) {
// 	time.Sleep(20 * time.Millisecond)
// 	fmt.Printf("Worker %d says: %s\n", id, phrase)
// }

// func makePool(poolSize int, handler func(int, string)) (func(string), func()) {
// 	handle := func(s string) {

// 	}

// 	wait := func() {

// 	}

// 	return handle, wait
// }

// func main() {
// 	fmt.Println("main start")

// 	phrases := []string{}

// 	for i := range 100 {
// 		phrases = append(phrases, fmt.Sprintf("phrase %d", i))
// 	}

// 	handle, wait := makePool(5, say)

// 	for _, phrase := range phrases {
// 		handle(phrase)
// 	}

// 	wait()

// 	fmt.Println("main end")
// }
