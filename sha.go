package main

import (
	"bufio"
	"fmt"
	"os"
)

func Ch(x, y, z uint32) uint32 {
	return (x & y) ^ (^x & z)
}

func Maj(x, y, z uint32) uint32 {
	return (x & y) ^ (x & z) ^ (y & z)
}

func SmallSigma0(x uint32) uint32 {
	return Rotr(x, 7) ^ Rotr(x, 18) ^ Shr(x, 3)
}
func SmallSigma1(x uint32) uint32 {
	return Rotr(x, 17) ^ Rotr(x, 19) ^ Shr(x, 10)
}
func LargeSigma0(x uint32) uint32 {
	return Rotr(x, 2) ^ Rotr(x, 13) ^ Rotr(x, 22)
}
func LargeSigma1(x uint32) uint32 {
	return Rotr(x, 6) ^ Rotr(x, 11) ^ Rotr(x, 25)
}

// Right rotation
func Rotr(x, n uint32) uint32 {
	return x<<(32-n) | x>>n
}

// Right shift
func Shr(x, n uint32) uint32 {
	return x >> n
}

func main() {

	// if len(os.Args) != 2 {
	// 	os.Exit(1)
	// }
	// initializing hash to the values
	H := []uint32{
		0x6a09e667,
		0xbb67ae85,
		0x3c6ef372,
		0xa54ff53a,
		0x510e527f,
		0x9b05688c,
		0x1f83d9ab,
		0x5be0cd19,
	}
	k := []uint32{
		0x428a2f98,
		0x71374491,
		0xb5c0fbcf,
		0xe9b5dba5,
		0x3956c25b,
		0x59f111f1,
		0x923f82a4,
		0xab1c5ed5,
		0xd807aa98,
		0x12835b01,
		0x243185be,
		0x550c7dc3,
		0x72be5d74,
		0x80deb1fe,
		0x9bdc06a7,
		0xc19bf174,
		0xe49b69c1,
		0xefbe4786,
		0x0fc19dc6,
		0x240ca1cc,
		0x2de92c6f,
		0x4a7484aa,
		0x5cb0a9dc,
		0x76f988da,
		0x983e5152,
		0xa831c66d,
		0xb00327c8,
		0xbf597fc7,
		0xc6e00bf3,
		0xd5a79147,
		0x06ca6351,
		0x14292967,
		0x27b70a85,
		0x2e1b2138,
		0x4d2c6dfc,
		0x53380d13,
		0x650a7354,
		0x766a0abb,
		0x81c2c92e,
		0x92722c85,
		0xa2bfe8a1,
		0xa81a664b,
		0xc24b8b70,
		0xc76c51a3,
		0xd192e819,
		0xd6990624,
		0xf40e3585,
		0x106aa070,
		0x19a4c116,
		0x1e376c08,
		0x2748774c,
		0x34b0bcb5,
		0x391c0cb3,
		0x4ed8aa4a,
		0x5b9cca4f,
		0x682e6ff3,
		0x748f82ee,
		0x78a5636f,
		0x84c87814,
		0x8cc70208,
		0x90befffa,
		0xa4506ceb,
		0xbef9a3f7,
		0xc67178f2,
	}
	// Receiving input
	fmt.Println("please give the input")
	read := bufio.NewReader(os.Stdin)
	input, _ := read.ReadString('\n')
	s := " "
	for i := 0; i < len(input)-1; i++ {
		s = s + string(input[i])
	}
	// Padding process
	pad := Padding([]byte(input), 64)
	// Create a message block by dividing it into 64 bytes
	msgblocks := Split(pad, 64)
	for _, bl := range msgblocks {
		// Create 64 word (uint32) array from message block
		words := uint32Array(bl)
		for i := 16; i < 64; i++ {
			w := SmallSigma1(words[i-2]) + words[i-7] + SmallSigma0(words[i-15]) + words[i-16]
			words = append(words, w)
		}
		// Initial value before rotation processing
		a := H[0]
		b := H[1]
		c := H[2]
		d := H[3]
		e := H[4]
		f := H[5]
		g := H[6]
		h := H[7]
		// Rotation processing
		for t, w := range words {
			T1 := h + LargeSigma1(e) + Ch(e, f, g) + k[t] + w
			T2 := LargeSigma0(a) + Maj(a, b, c)
			h = g
			g = f
			f = e
			e = d + T1
			d = c
			c = b
			b = a
			a = T1 + T2
		}

		// Updata hash value
		H[0] = a + H[0]
		H[1] = b + H[1]
		H[2] = c + H[2]
		H[3] = d + H[3]
		H[4] = e + H[4]
		H[5] = f + H[5]
		H[6] = g + H[6]
		H[7] = h + H[7]
	}
	// Display hash value
	for _, h := range H {
		fmt.Printf("%x", h)
	}

}

// Make the input byte array a multiple of length
func Padding(input []byte, length int) []byte {
	l := len(input)
	bits := l * 8
	mod := l % length
	padcount := length - mod
	if mod > length-8 {
		padcount += 64
	}
	for i := 0; i < padcount; i++ {
		if i == 0 {
			// Put 0x80 as an input delimiter
			input = append(input, 0x80)
		} else {
			// Others are filled with 0x00
			input = append(input, 0x00)
		}
	}
	// The last 8 bytes (uint64) are the number of bits in the input
	for i := 1; i <= 8; i++ {
		input[len(input)-i] = byte(bits & 0xff)
		bits = bits >> 8
	}
	return input
}

// Divide the input byte array into length pieces
func Split(input []byte, length int) [][]byte {
	var barr [][]byte
	n := len(input) / length
	for i := 0; i < n; i++ {
		barr = append(barr, input[i*length:(i+1)*length])
	}
	return barr
}

// Convert byte array to uint32 array
func uint32Array(b []byte) []uint32 {
	var arr []uint32
	for i := 0; i < len(b)/4; i++ {
		var v uint32
		v += uint32(b[i*4]) << 24
		v += uint32(b[i*4+1]) << 16
		v += uint32(b[i*4+2]) << 8
		v += uint32(b[i*4+3])
		arr = append(arr, v)
	}
	return arr
}
