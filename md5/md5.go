package main

const BITS_8 = 0xff   // bytes - 512 bits
const CHUNK_SIZE = 64 // bytes - 512 bits

// States
const A = 0x67452301
const B = 0xefcdab89
const C = 0x98badcfe
const D = 0x10325476

/*
 To pad the original data as needed by the algorithm.

 1. Add a 1 bit after the data.
 2. Pad with 0 until getting to a multiple of 512 minus 64 (here will be set the original data size % 2 ^ 64)
*/
func padData(data []byte) []byte {
	originalSize := len(data)

	paddedData := data

	// Set padding
	paddedData = append(paddedData, 0x80) // 10000000
	for len(paddedData)%64 != 56 {
		paddedData = append(paddedData, 0x0)
	}

	// Append the original data size (64 bits)
	tempSize := int64(originalSize)
	bytes := int64(8)

	for i := int64(bytes - 1); int64(i) >= 0; i-- {
		paddedData = append(paddedData, byte((tempSize>>(bytes*i))&BITS_8))
	}

	return paddedData
}

// Sines of integers (radians)
func getK() [64]int {
	return [64]int{
		0xd76aa478, 0xe8c7b756, 0x242070db, 0xc1bdceee, 0xf57c0faf, 0x4787c62a,
		0xa8304613, 0xfd469501, 0x698098d8, 0x8b44f7af, 0xffff5bb1, 0x895cd7be,
		0x6b901122, 0xfd987193, 0xa679438e, 0x49b40821, 0xf61e2562, 0xc040b340,
		0x265e5a51, 0xe9b6c7aa, 0xd62f105d, 0x02441453, 0xd8a1e681, 0xe7d3fbc8,
		0x21e1cde6, 0xc33707d6, 0xf4d50d87, 0x455a14ed, 0xa9e3e905, 0xfcefa3f8,
		0x676f02d9, 0x8d2a4c8a, 0xfffa3942, 0x8771f681, 0x6d9d6122, 0xfde5380c,
		0xa4beea44, 0x4bdecfa9, 0xf6bb4b60, 0xbebfbc70, 0x289b7ec6, 0xeaa127fa,
		0xd4ef3085, 0x04881d05, 0xd9d4d039, 0xe6db99e5, 0x1fa27cf8, 0xc4ac5665,
		0xf4292244, 0x432aff97, 0xab9423a7, 0xfc93a039, 0x655b59c3, 0x8f0ccc92,
		0xffeff47d, 0x85845dd1, 0x6fa87e4f, 0xfe2ce6e0, 0xa3014314, 0x4e0811a1,
		0xf7537e82, 0xbd3af235, 0x2ad7d2bb, 0xeb86d391,
	}
}

// Per-round shift amounts
func getS() [64]int {
	return [64]int{
		7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22, 5, 9, 14,
		20, 5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14, 20, 4, 11, 16, 23, 4, 11,
		16, 23, 4, 11, 16, 23, 4, 11, 16, 23, 6, 10, 15, 21, 6, 10, 15, 21, 6,
		10, 15, 21, 6, 10, 15, 21,
	}
}

// func updateHash(data []byte) {

// }
