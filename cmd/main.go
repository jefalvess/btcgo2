package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"sync"
	"time"

	"runtime"

	"github.com/btcsuite/btcutil/base58"
	"github.com/decred/dcrd/dcrec/secp256k1/v4"
	"golang.org/x/crypto/ripemd160"
)

// Gera o WIF (Wallet Import Format) a partir da chave privada
func GenerateWif(privKeyInt *big.Int) string {
	privKeyHex := fmt.Sprintf("%064x", privKeyInt)
	privKeyBytes, err := hex.DecodeString(privKeyHex)
	if err != nil {
		log.Fatal(err)
	}

	extendedKey := append([]byte{0x80}, privKeyBytes...)
	extendedKey = append(extendedKey, 0x01)

	firstSHA := sha256.Sum256(extendedKey)
	secondSHA := sha256.Sum256(firstSHA[:])
	checksum := secondSHA[:4]

	finalKey := append(extendedKey, checksum...)
	wif := base58.Encode(finalKey)
	return wif
}

func GeneratePublicKey(privKeyInt *big.Int) []byte {
	privKeyBytes := privKeyInt.Bytes()

	// Cria uma nova chave privada usando o pacote secp256k1
	privKey := secp256k1.PrivKeyFromBytes(privKeyBytes)

	// Obtém a chave pública correspondente no formato comprimido
	compressedPubKey := privKey.PubKey().SerializeCompressed()
	return compressedPubKey
}

// Cria o hash160 da chave pública comprimida
func CreatePublicHash160(pubKey []byte) []byte {
	pubKeyHash := hash160(pubKey)
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
	versionedPayload := append([]byte{0x00}, hash160...)
	checksum := checksum(versionedPayload)
	fullPayload := append(versionedPayload, checksum...)
	return base58.Encode(fullPayload)
}

// Realiza o hash160 (SHA256 seguido de RIPEMD160)
func hash160(b []byte) []byte {
	h := sha256.New()
	h.Write(b)
	sha256Hash := h.Sum(nil)

	r := ripemd160.New()
	r.Write(sha256Hash)
	return r.Sum(nil)
}

// Processar a chave pública
func processPublicKey(privKeyInt *big.Int, resultChan chan string, mutex *sync.Mutex) {
	// Gerar a chave pública
	pubKey := GeneratePublicKey(privKeyInt)

	// Criar o hash da chave pública
	pubKeyHash := CreatePublicHash160(pubKey)
	pubKeyHashStr := hex.EncodeToString(pubKeyHash)

	// Comparar com o valor desejado
	if pubKeyHashStr == "739437bb3dd6d1983e66629c5f08c70e52769371" {
		// Imprime no console
		// fmt.Println("Hash Público (RIPEMD-160):", pubKeyHashStr)

		// Abre o arquivo ou cria um novo arquivo txt para escrita
		mutex.Lock() // Protege o acesso ao arquivo
		defer mutex.Unlock()

		file, err := os.OpenFile("hash_publico.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("Erro ao criar o arquivo:", err)
			return
		}
		defer file.Close() // Garante que o arquivo será fechado após a operação
		bigIntStr := privKeyInt.String()

		// Escreve o hash público no arquivo
		_, err = file.WriteString("Hash Público (RIPEMD-160): " + pubKeyHashStr + "   chave:  " + bigIntStr + "\n")
		if err != nil {
			fmt.Println("Erro ao escrever no arquivo:", err)
			return
		}

		fmt.Println("Hash salvo em hash_publico.txt")
	}

	// Envia o hash para o canal
	resultChan <- pubKeyHashStr
}

func main() {
	// Definir o número de núcleos a serem usados
	runtime.GOMAXPROCS(10) // Ajuste o número de núcleos para o seu hardware

	privKeyDecimal := "147573952589614412928" // Exemplo de chave privada decimal

	// Converte para BigInt (chave privada em decimal)
	privKeyInt := new(big.Int)
	privKeyInt.SetString(privKeyDecimal, 10) // Base 10 para BigInt

	// Canal para coletar os resultados
	resultChan := make(chan string)
	var wg sync.WaitGroup

	// Mutex para proteger a escrita no arquivo
	var mutex sync.Mutex

	// Limitar o número de goroutines simultâneas
	numWorkers := 10
	semaphore := make(chan struct{}, numWorkers)

	// Medir o tempo de execução do loop
	startTime := time.Now() // Começa a medir o tempo

	for i := 1000 * 1000 * 1000; i >= 0; i-- {
		wg.Add(1)

		// Cria uma cópia de privKeyInt para passar para a goroutine
		privKeyCopy := new(big.Int).Set(privKeyInt)

		go func() {
			defer wg.Done()

			// Limita a quantidade de goroutines simultâneas
			semaphore <- struct{}{} // Adquire um slot no pool

			processPublicKey(privKeyCopy, resultChan, &mutex)

			// Libera o slot no pool
			<-semaphore
		}()

		if i%1000000 == 0 {
			bigIntStr := privKeyInt.String()
			// Abre o arquivo ou cria um novo arquivo para gravação
			file, err := os.OpenFile("progress.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
			// Escreve a posição atual no arquivo
			_, err = file.WriteString(fmt.Sprintf("Progresso: Iteração->" + bigIntStr + "\n"))
			if err != nil {
				log.Fatal(err)
			}
		}

		// Aumenta o valor da chave privada
		privKeyInt.Sub(privKeyInt, big.NewInt(1))
	}

	// Aguarda até que todas as goroutines terminem
	go func() {
		wg.Wait()
		close(resultChan) // Fecha o canal depois que todas as goroutines terminarem
	}()

	// Medir o tempo de execução total após o loop
	elapsedTime := time.Since(startTime)
	fmt.Printf("Tempo total para processar as goroutines: %s\n", elapsedTime)
}
