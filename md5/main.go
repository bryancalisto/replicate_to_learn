package main

func main() {
	arr := []byte{}
	for i := 0; i < 55; i++ {
		arr = append(arr, byte(i))
	}

	padData(arr)
}
