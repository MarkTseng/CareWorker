package main

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
)

func genRandString() string {
	buf := new(bytes.Buffer)
	io.CopyN(buf, rand.Reader, 32)
	return hex.EncodeToString(buf.Bytes())
}

func DoHash(pass, salt string) string {
	h := sha256.New()
	h.Write([]byte(pass))
	h.Write([]byte(salt))
	return hex.EncodeToString(h.Sum(nil))
}

func main() {
	salt := string("fc8e4231f10c2990ff39b670e56335fdb3cf42048eb9f2759bd8b9305f0391f3")
	password := string("123456")
	salt_pass := DoHash(password, salt)
	fmt.Println(salt_pass)
}
