package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	KeySize   = 32
	NonceSize = 12
)

// encrypt encrypts plaintext using a key and returns the ciphertext
func encrypt(k, v []byte) ([]byte, error) {
	block, err := aes.NewCipher(k)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	ciphertext := aesgcm.Seal(nonce, nonce, v, nil)
	return []byte(base64.URLEncoding.EncodeToString(ciphertext)), nil
}

// decrypt decrypts ciphertext using a key and returns the plaintext
func decrypt(k, v []byte) ([]byte, error) {
	data, err := base64.URLEncoding.DecodeString(string(v))
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(k))
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesgcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func makeKey(in string) []byte {
	l := len(in)
	if l > KeySize {
		in = in[:KeySize]
	} else if l < KeySize {
		for i := 0; i < KeySize-l; i++ {
			in += string(in[i%l])
		}
	}
	return []byte(in)
}

var (
	inFile      string
	outFile     string
	key         string
	decryptMode bool
	encryptMode bool
	mode        string
)

func main() {
	flag.StringVar(&inFile, "in", "-", "The target input file. By default input is read from stdin.")
	flag.StringVar(&inFile, "i", "-", "Alias for in.")
	flag.StringVar(&outFile, "out", "-", "The target output file. By default output goes to stdout.")
	flag.StringVar(&outFile, "o", "-", "Alias for out.")
	flag.StringVar(&key, "key", "", "The encryption/decryption key")
	flag.StringVar(&key, "k", "", "Alias for key")
	flag.BoolVar(&decryptMode, "d", false, "Decryption mode. Input data will be decrypted")
	flag.BoolVar(&decryptMode, "decrypt", false, "Decryption mode. Input data will be decrypted")
	flag.BoolVar(&encryptMode, "e", false, "Encryption mode. Input data will be encrypted")
	flag.BoolVar(&encryptMode, "encrypt", false, "Encryption mode. Input data will be encrypted")
	flag.Parse()

	var fn func([]byte, []byte) ([]byte, error)

	if decryptMode {
		fn = decrypt
		mode = "decrypt"
	} else if encryptMode {
		fn = encrypt
		mode = "encrypt"
	} else {
		log.Fatalf("missing the mode flag: expected exactly one of '-encrypt' or '-decrypt'")
	}

	for key == "" {
		fmt.Fprintf(os.Stderr, "Enter pass key to %s: ", mode)
		if b, err := terminal.ReadPassword(int(syscall.Stdin)); err == nil {
			key = string(b)
			fmt.Fprintln(os.Stderr, "")
		} else {
			log.Fatalln("failed to read pass key:", err)
		}
	}
	inData := readData(inFile)
	outData, err := fn(makeKey(key), inData)
	if err != nil {
		log.Fatalf("failed to %s data: %v", mode, err)
	}
	writeData(outFile, outData)
}

func readData(f string) []byte {
	var r io.Reader
	switch f {
	case "-":
		r = os.Stdin
	default:
		file, err := os.Open(f)
		if err != nil {
			log.Fatalf("failed to read %s: %v", f, err)
		}
		defer file.Close()
		r = file
	}

	inData, err := io.ReadAll(r)
	if err != nil {
		log.Fatalf("failed to read data from %s: %v", f, err)
	}
	return inData
}

func writeData(f string, d []byte) {
	switch outFile {
	case "-":
		fmt.Print(string(d))
	default:
		if err := os.WriteFile(f, d, 0600); err != nil {
			log.Fatalf("failed to write %sed data to %s: %v", mode, f, err)
		}
		fmt.Printf("wrote %sed data to %s\n", mode, f)
	}
}
