package main

import (
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Testsha256(t *testing.T) {
	message := "hii,helloworld"
	hash := string(sha256([]byte(message)))
	h := sha256.New()
	h.write([]byte(message))
	inbuilthash := string(h.Sum(nil))
	assert.Equal(t, hash, inbuilthash)
}
