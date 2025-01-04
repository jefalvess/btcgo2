package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/btcsuite/btcutil/base58"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"golang.org/x/crypto/ripemd160"
)

// Gera o WIF (Wallet Import Format) a partir da chave privada
func GenerateWif(privKeyInt *big.Int) string {
	// Converte a chave privada para hexadecimal
	privKeyHex := fmt.Sprintf("%064x", privKeyInt)

	// Decodifica a chave privada em bytes
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	// Adiciona o prefixo (0x80 para Bitcoin) e sufixo (0x01 para chave privada WIF)
	extendedKey := append([]byte{0x80}, privKeyBytes...)
	extendedKey = append(extendedKey, 0x01)

	// Calcula o checksum (SHA256 duplo)
	firstSHA := sha256.Sum256(extendedKey)
	secondSHA := sha256.Sum256(firstSHA[:])
	checksum := secondSHA[:4]

	// Adiciona o checksum ao final
	finalKey := append(extendedKey, checksum...)

	// Codifica a chave final em Base58 e retorna o WIF
	wif := base58.Encode(finalKey)
	return wif
}

// Cria o hash160 da chave pública comprimida a partir da chave privada
func CreatePublicHash160(privKeyInt *big.Int) []byte {
	privKeyBytes := privKeyInt.Bytes()

	// Cria a chave privada usando a biblioteca secp256k1
	privKey := secp256k1.PrivKeyFromBytes(privKeyBytes)

	// Gera a chave pública comprimida
	compressedPubKey := privKey.PubKey().SerializeCompressed()

	// Gera o hash160 da chave pública comprimida
	pubKeyHash := hash160(compressedPubKey)

	return pubKeyHash
}

// Calcula o checksum necessário para o endereço Bitcoin
func checksum(payload []byte) []byte {
	hash1 := sha256.Sum256(payload)
	hash2 := sha256.Sum256(hash1[:])
	return hash2[:4]
}

// Gera o endereço Bitcoin a partir do hash160
func Hash160ToAddress(hash160 []byte) string {
	// Adiciona o prefixo para a rede Bitcoin (0x00 para Mainnet)
	versionedPayload := append([]byte{0x00}, hash160...)
	// Calcula o checksum
	checksum := checksum(versionedPayload)
	// Adiciona o checksum ao final
	fullPayload := append(versionedPayload, checksum...)
	// Codifica o endereço final em Base58
	return base58.Encode(fullPayload)
}

// Realiza o hash160 (SHA256 seguido de RIPEMD160)
func hash160(b []byte) []byte {
	// Primeiro aplica o SHA-256
	h := sha256.New()
	h.Write(b)
	sha256Hash := h.Sum(nil)

	// Em seguida aplica o RIPEMD-160
	r := ripemd160.New()
	r.Write(sha256Hash)
	return r.Sum(nil)
}
