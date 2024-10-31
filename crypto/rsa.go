package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"log"
)

const bitSize = 2048

func main() {
	privKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Private key: %+v\n", privKey)

	publicKey := x509.MarshalPKCS1PublicKey(&privKey.PublicKey)
	fmt.Printf("Public key: %+v\n", hex.EncodeToString(publicKey))

	message := "Hello, World"

	encryptedBytes, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, &privKey.PublicKey, []byte(message), nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("encrypted bytes: ", encryptedBytes)

	decryptedBytes, err := privKey.Decrypt(rand.Reader, encryptedBytes, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		panic(err)
	}
	fmt.Println("decrypted message: ", string(decryptedBytes))

	msgHash := sha256.New()
	_, err = msgHash.Write([]byte(message))
	if err != nil {
		panic(err)
	}
	msgHashSum := msgHash.Sum(nil)

	signature, err := rsa.SignPSS(rand.Reader, privKey, crypto.SHA256, msgHashSum, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println("Signature: ", signature)

	err = rsa.VerifyPSS(&privKey.PublicKey, crypto.SHA256, msgHashSum, signature, nil)
	if err != nil {
		fmt.Println("could not verify signature: ", err)
		return
	} else {
		fmt.Println("signature verified")
	}
}
