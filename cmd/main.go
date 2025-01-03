package main

import (
	"btcgo/cmd/utils"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"runtime"
	"sync"
	"time"
)

type Wallets struct {
	Wallets          []string `json:"Wallets"`
	FileName         string
	SearchingWallets string
	DataWallet       map[string]bool
	DataWalletID     map[int]string
}

func SalvarBufferEmArquivo(buffer []string, filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Erro ao abrir o arquivo: %v", err)
	}
	defer file.Close()

	for _, line := range buffer {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			log.Fatalf("Erro ao escrever no arquivo: %v", err)
		}
	}
}

func processarIntervalo(inicio *big.Int, fim *big.Int, Wallets *Wallets, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()
	const checkpointInterval = 1000000 * 1000 // 1 bilhao de registros

	privKeyHex := ""
	privKeyInt := new(big.Int)

	for i := new(big.Int).Set(inicio); i.Cmp(fim) <= 0; i.Add(i, big.NewInt(1)) {
		privKeyHex = fmt.Sprintf("%064x", i) // Preenche com zeros à esquerda para garantir 64 caracteres
		privKeyInt.SetString(privKeyHex, 16)
		address := utils.CreatePublicHash160(privKeyInt)

		if _, exists := Wallets.DataWallet[string(address)]; exists {
			wallet := utils.Hash160ToAddress(address)
			wif := utils.GenerateWif(privKeyInt)

			line := fmt.Sprintf("%s -> %s -> %s -> %s", privKeyHex, wallet, wif, time.Now().Format("2006-01-02 15:04:05"))
			SalvarBufferEmArquivo([]string{line}, "wallets.txt")

		}

		// Verifica se atingiu o checkpoint
		if workerID == 1 && i.Int64()%checkpointInterval == 0 {
			checkpoint := fmt.Sprintf("Checkpoint: Worker %d processou até %s em %s", workerID, i.String(), time.Now().Format("2006-01-02 15:04:05"))
			SalvarBufferEmArquivo([]string{checkpoint}, "checkpoint-1.txt")
		}
	}
}

func main() {
	// Configura o Go para usar todos os núcleos disponíveis
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Carregar os dados de wallets
	bytes, err := os.ReadFile("data/Wallets.json")
	if err != nil {
		log.Fatal(err)
	}

	var Wallets Wallets
	if err := json.Unmarshal(bytes, &Wallets); err != nil {
		log.Fatal(err)
	}

	// Inicializando os mapas de DataWallet e DataWalletID corretamente
	Wallets.DataWallet = make(map[string]bool)
	Wallets.DataWalletID = make(map[int]string)

	for i, address := range Wallets.Wallets {
		Wallets.DataWallet[string(utils.Decode(address)[1:21])] = true
		Wallets.DataWalletID[i] = address
	}

	// Valores grandes (total e numWorkers)
	total := new(big.Int)
	total.SetString("146346217550346335726", 10) // Total
	start := new(big.Int)
	start.SetString("46346217550346335726", 10) // Início

	numWorkers := runtime.NumCPU() // Usando o número de CPUs disponíveis
	fmt.Print(numWorkers)
	interval := new(big.Int).Div(new(big.Int).Sub(total, start), big.NewInt(int64(numWorkers)))

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		inicio := new(big.Int).Add(start, new(big.Int).Mul(big.NewInt(int64(i)), interval))
		fim := new(big.Int).Add(inicio, interval)

		if i == numWorkers-1 { // Ajusta o intervalo do último worker
			fim = total
		}

		wg.Add(1)
		go processarIntervalo(inicio, fim, &Wallets, &wg, i+1)
	}

	wg.Wait()

	fmt.Println("Processamento concluído! Todos os workers finalizaram suas tarefas.")

}
