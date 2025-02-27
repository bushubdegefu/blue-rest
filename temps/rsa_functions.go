package temps

import (
	"fmt"
	"os"
	"text/template"
)

func RSAHelper() {
	// ####################################################
	//  rabbit template
	helper_tmpl, err := template.New("RenderData").Parse(rsaUtilsTemplate)
	if err != nil {
		panic(err)
	}

	// creating config folder for .env files
	err = os.MkdirAll("helper", os.ModePerm)
	if err != nil {
		fmt.Printf("Frame - 10: %v\n", err)
	}

	helper_file, err := os.Create("helper/rsa_helper.go")
	if err != nil {
		panic(err)
	}
	defer helper_file.Close()

	err = helper_tmpl.Execute(helper_file, nil)
	if err != nil {
		panic(err)
	}
}

var rsaUtilsTemplate = `
package helper

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

// Generate RSA Keys
func generateRSAKeys(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	// Generate private and public keys
	priv, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return priv, &priv.PublicKey, nil
}

// RSA Encrypt
func rsaEncrypt(publicKey *rsa.PublicKey, message []byte) ([]byte, error) {
	// Encrypt the message using RSA and SHA256 as hashing
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, message, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// RSA Decrypt
func rsaDecrypt(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	// Decrypt the message using RSA and SHA256
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

`

var rsaExampleMainFunction = `
package main


func main() {
	// Step 1: Generate RSA Keys
	privateKey, publicKey, err := generateRSAKeys(2048)
	if err != nil {
		log.Fatal("Error generating RSA keys:", err)
	}

	// Step 2: Encrypt a message using the RSA public key
	message := []byte("This is a secret message.")
	fmt.Printf("Original Message: %s\n", message)

	// Encrypt message using RSA Public Key
	ciphertext, err := rsaEncrypt(publicKey, message)
	if err != nil {
		log.Fatal("Error encrypting message:", err)
	}

	// Print the encrypted ciphertext (as a hexadecimal string)
	fmt.Printf("Encrypted Message: %x\n", ciphertext)

	// Step 3: Decrypt the ciphertext using the RSA private key
	plaintext, err := rsaDecrypt(privateKey, ciphertext)
	if err != nil {
		log.Fatal("Error decrypting message:", err)
	}

	// Print the decrypted message
	fmt.Printf("Decrypted Message: %s\n", plaintext)
}`
