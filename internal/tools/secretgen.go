package main

import (
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"os"
)


func main() {
	n := flag.Int("n", 1, "number of secrets to generate")
	bytes := flag.Int("b", 32, "bytes to entropy per secret")
	flag.Parse()

	for i := 0; i < *n; i++ {
		secret, err := generateSecret(*bytes)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		fmt.Println(secret)
	}
}

func generateSecret(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
			return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
