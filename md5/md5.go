// Adapted from https://en.wikipedia.org/wiki/MD5
// Support for little endian encoding only for now
package main

const MASK_8_BITS = 0xff
const MASK_32_BITS = 0xffffffff
const CHUNK_SIZE = 64 // bytes - 512 bits
const BYTES_PER_WORD = 4
const BITS_PER_BYTE = 8

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
	for len(paddedData)%CHUNK_SIZE != 56 {
		paddedData = append(paddedData, 0x0)
	}

	// Append the original data size (64 bits). Just little endian by now
	tempSize := int64(originalSize) * BITS_PER_BYTE
	bytes := int64(8)

	for i := int64(0); int64(i) < bytes; i++ {
		paddedData = append(paddedData, byte((tempSize>>(bytes*i))&MASK_8_BITS))
	}

	return convertToLittleEndian(paddedData)
	// return paddedData
}

// Sines of integers (radians)
func getK() [64]uint32 {
	return [64]uint32{
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
func getS() [64]uint32 {
	return [64]uint32{
		7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22, 7, 12, 17, 22,
		5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14, 20, 5, 9, 14, 20,
		4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23, 4, 11, 16, 23,
		6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21, 6, 10, 15, 21,
	}
}

func hash(data []byte) [16]byte {
	a0 := uint32(0x67452301)
	b0 := uint32(0xefcdab89)
	c0 := uint32(0x98badcfe)
	d0 := uint32(0x10325476)
	A := a0
	B := b0
	C := c0
	D := d0
	K := getK()
	S := getS()
	var g uint32
	var F uint32
	digest := [16]byte{}

	// Set padding
	data = padData(data)

	// Process the hash itself
	for i := 0; i < len(data)/CHUNK_SIZE; i++ {
		M := [16]uint32{}
		for j := 0; j < 16; j++ {
			offset := i*16 + j*4
			M[j] = uint32(data[offset])<<24 | uint32(data[offset+1])<<16 | uint32(data[offset+2])<<8 | uint32(data[offset+3])
			// M[j] = uint32(data[offset+3])<<24 | uint32(data[offset+3])<<16 | uint32(data[offset+1])<<8 | uint32(data[offset])
			// M[j] = uint32(data[offset]) | uint32(data[offset+1])<<8 | uint32(data[offset+2])<<16 | uint32(data[offset+3])<<24
		}

		for k := uint32(0); k < 64; k++ {
			if k <= 15 {
				F = (B & C) | (^B & D)
				g = k
			} else if k >= 16 && k <= 31 {
				F = (D & B) | (^D & C)
				g = (5*k + 1) % 16
			} else if k >= 32 && k <= 47 {
				F = B ^ C ^ D
				g = (3*k + 5) % 16
			} else if k >= 48 && k <= 63 {
				F = C ^ (B | ^D)
				g = (7 * k) % 16
			}

			// F = add32Bit(add32Bit(F, A), add32Bit(K[k], M[g]))
			F = F + A + K[k] + M[g]
			A = D
			D = C
			C = B
			B = B + leftRotate(F, S[k])
		}
		a0 += A
		b0 += B
		c0 += C
		d0 += D
	}

	// Form the resulting digest.
	// I didn't want to write the all the digest byte assignations manually, that's why I used this loop
	processedData := [4]uint32{a0, b0, c0, d0}
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			digest[i*4+j] = byte((processedData[i] >> (j * 8)))
		}
	}

	return digest
}

func leftRotate(val uint32, shift uint32) uint32 {
	var mod_shift = shift & 31
	return ((val << mod_shift) & MASK_32_BITS) | ((val & MASK_32_BITS) >> (32 - mod_shift))
}

func add32Bit(a uint32, b uint32) uint32 {
	return uint32(a+b) & MASK_32_BITS
}

func convertToLittleEndian(data []byte) []byte {
	var temp0 byte
	var temp1 byte
	var temp2 byte

	for i := 0; i < len(data); i += BYTES_PER_WORD {
		temp0 = data[i+3]
		temp1 = data[i+2]
		temp2 = data[i+1]
		data[i+3] = data[i]
		data[i+2] = temp2
		data[i+1] = temp1
		data[i] = temp0
	}

	return data
}
