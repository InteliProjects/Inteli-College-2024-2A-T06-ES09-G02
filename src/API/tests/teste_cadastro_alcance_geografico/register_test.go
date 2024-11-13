package teste_cadastro_alcance_geografico

import (
	"2024-2A-T06-ES09-G02/src/API/handlers"
	"2024-2A-T06-ES09-G02/src/API/logging"
	"2024-2A-T06-ES09-G02/src/API/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap"
)

// Estrutura para o log em JSON
type LogData struct {
	State          string  `json:"state"`
	LatencyIn      float64 `json:"latency_in"`
	ProcessingTime float64 `json:"processing_time"`
	LatencyOut     float64 `json:"latency_out"`
	DownloadTime   float64 `json:"download_time"`
	StatusCode     int     `json:"status_code"`
	ResponseText   string  `json:"response_text"`
	ResponseSize   int     `json:"response_size"`
	TotalTime      float64 `json:"total_time"`
}

// Simula latência normal baseada em uma média e desvio padrão
func simulateLatency(mean, stdDev float64) float64 {
	latency := rand.NormFloat64()*stdDev + mean
	if latency < 0 {
		return 0
	}
	return latency
}

// Simula o tempo de download baseado no tamanho da resposta e na velocidade de download
func calculateDownloadTime(responseSize int, downloadSpeedKbps float64) float64 {
	downloadSpeedBps := downloadSpeedKbps * 1000 / 8
	downloadTime := float64(responseSize) / downloadSpeedBps
	return downloadTime
}

// Função que escreve os logs no arquivo JSON adicionando ao array
func logToFile(logData LogData) {
	// Tenta abrir o arquivo de logs existente
	file, err := os.OpenFile("logs.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Erro ao abrir arquivo de logs:", err)
		return
	}
	defer file.Close()

	// Carrega os logs existentes no array, ou inicializa um novo array vazio se o arquivo estiver vazio
	var logs []LogData

	// Lê o conteúdo do arquivo
	fileContent, err := ioutil.ReadAll(file)
	if err == nil && len(fileContent) > 0 {
		// Deserializa o conteúdo do arquivo JSON para o slice de logs, se o conteúdo não estiver vazio
		err = json.Unmarshal(fileContent, &logs)
		if err != nil {
			fmt.Println("Erro ao parsear arquivo JSON:", err)
			return
		}
	}

	// Adiciona o novo log ao array de logs
	logs = append(logs, logData)

	// Move o ponteiro do arquivo para o início e trunca o arquivo para reescrevê-lo
	file.Seek(0, 0)
	file.Truncate(0)

	// Serializa o array de logs atualizado de volta para o arquivo
	updatedLogs, err := json.MarshalIndent(logs, "", "  ") // Formatação bonita com indentação
	if err != nil {
		fmt.Println("Erro ao converter logs para JSON:", err)
		return
	}

	// Escreve os logs atualizados no arquivo
	_, err = file.Write(updatedLogs)
	if err != nil {
		fmt.Println("Erro ao escrever logs atualizados no arquivo:", err)
	}
}

// Função que gera o log de uma requisição em formato JSON
func logRequest(latencyIn, processingTime, latencyOut, downloadTime float64, statusCode int, state, responseText string, responseSize int) {
	totalTime := latencyIn + processingTime + latencyOut + downloadTime

	// Cria a estrutura de log em JSON
	logData := LogData{
		State:          state,
		LatencyIn:      latencyIn,
		ProcessingTime: processingTime,
		LatencyOut:     latencyOut,
		DownloadTime:   downloadTime,
		StatusCode:     statusCode,
		ResponseText:   responseText,
		ResponseSize:   responseSize,
		TotalTime:      totalTime,
	}

	// Escreve o log no arquivo logs.json
	logToFile(logData)
}

// Simula o processamento completo da requisição
func processRequest(state string, meanLatency, stdDevLatency, downloadSpeedKbps float64, wg *sync.WaitGroup, results *[]float64, logger *zap.Logger) {
	defer wg.Done()

	latencyIn := simulateLatency(meanLatency, stdDevLatency)
	time.Sleep(time.Duration(latencyIn * float64(time.Second)))

	startTime := time.Now()
	responseText, statusCode := testRegisterHandler(state, logger) // Substitui CadastroUsuario
	processingTime := time.Since(startTime).Seconds()

	latencyOut := simulateLatency(meanLatency, stdDevLatency)
	time.Sleep(time.Duration(latencyOut * float64(time.Second)))

	responseSize := len(responseText)
	downloadTime := calculateDownloadTime(responseSize, downloadSpeedKbps)
	time.Sleep(time.Duration(downloadTime * float64(time.Second)))

	logRequest(latencyIn, processingTime, latencyOut, downloadTime, statusCode, state, responseText, responseSize)

	totalTime := latencyIn + processingTime + latencyOut + downloadTime
	*results = append(*results, totalTime)
}

// Função que simula a chamada da função RegisterHandler
func testRegisterHandler(state string, logger *zap.Logger) (string, int) {
	user := models.User{ // Cria um mock do usuário
		Username: "Usuario_" + state,
		Email:    "usuario_" + state + "@exemplo.com",
	}

	// Converte o usuário para JSON
	userJSON, _ := json.Marshal(user)

	// Cria uma requisição HTTP simulada
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(userJSON))
	rec := httptest.NewRecorder()

	// Chama a função RegisterHandler com logger de teste
	handlers.RegisterHandler(rec, req, logger)

	time.Sleep(400 *time.Millisecond)

	// Captura o status e a resposta
	res := rec.Result()
	defer res.Body.Close()

	responseBytes, _ := json.Marshal(res)
	return string(responseBytes), res.StatusCode
}

// Função para executar os testes em todos os estados
func runTestsForState(state string, meanLatency, stdDevLatency, downloadSpeedKbps float64, results *[]float64, logger *zap.Logger) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go processRequest(state, meanLatency, stdDevLatency, downloadSpeedKbps, &wg, results, logger)
		time.Sleep(1 * time.Second)
	}
	wg.Wait()
}

// Função que verifica se o teste passou ou falhou
func evaluateResults(results []float64, threshold float64) bool {
	passed := 0
	for _, result := range results {
		if result <= threshold {
			passed++
		}
	}
	successRate := (float64(passed) / float64(len(results))) * 100
	fmt.Printf("Porcentagem de requisições dentro do tempo: %.2f%%\n", successRate)
	return successRate >= 80
}

// Teste automatizado para RegisterHandler
func TestRegisterHandler(t *testing.T) {
	// Inicia o logger para testes
	err := logging.InitTestLogger()
	if err != nil {
		t.Fatalf("Erro ao inicializar logger de teste: %v", err)
	}
	logger := logging.GetTestLogger()

	states := map[string]struct {
		meanLatency       float64
		stdDevLatency     float64
		downloadSpeedKbps float64
	}{
		"SP": {0.01650, 0.01, 95315.0},
		"RJ": {0.01781, 0.01, 91055.3},
		"MG": {0.02654, 0.01, 60995.7},
		"RS": {0.03050, 0.01, 60018.1},
		"PA": {0.05606, 0.01, 41395.6},
		"SC": {0.02775, 0.01, 73512.2},
		"PR": {0.02367, 0.01, 67796.6},
		"MS": {0.03567, 0.01, 51616.1},
		"ES": {0.03298, 0.01, 56730.6},
		"GO": {0.02859, 0.01, 65946.2},
		"MT": {0.03194, 0.01, 89959.2},
		"BA": {0.04459, 0.01, 36541.9},
		"TO": {0.04470, 0.01, 55911.7},
		"RO": {0.04615, 0.01, 73750.8},
		"AC": {0.07499, 0.01, 60657.2},
		"SE": {0.04564, 0.01, 65269.6},
		"PI": {0.05731, 0.01, 52569.7},
		"AL": {0.05850, 0.01, 43606.9},
		"PE": {0.05594, 0.01, 48559.2},
		"PB": {0.04667, 0.01, 74056.9},
		"RN": {0.05510, 0.01, 74966.8},
		"CE": {0.05620, 0.01, 52344.4},
		"MA": {0.06423, 0.01, 35426.4},
		"AM": {0.07000, 0.01, 21773.8},
		"RR": {0.10162, 0.01, 18686.0},
		"AP": {0.09391, 0.01, 14718.1},
	}

	var allResults []float64

	for state, data := range states {
		fmt.Printf("\nExecutando teste para o estado: %s\n", state)
		var stateResults []float64
		runTestsForState(state, data.meanLatency, data.stdDevLatency, data.downloadSpeedKbps, &stateResults, logger)
		allResults = append(allResults, stateResults...)
	}

	if !evaluateResults(allResults, 1.0) {
		t.Errorf("Teste falhou. Algumas requisições excederam o limite de tempo.")
	}

	// Define o tempo máximo de resposta em milissegundos (1 segundo)
	tempoDeRespostaMaximo := 1.0 // Segundos

	// Teste com filtro de estado (SP)
	percentual, err := calcularPercentualDentroDoSLA("SP", tempoDeRespostaMaximo)
	if err != nil {
		t.Fatalf("Erro ao calcular percentual dentro do SLA: %v", err)
	}

	// Verifica se o percentual está acima de 80% (limite de aprovação)
	if percentual < 0.95 {
		t.Errorf("Percentual de respostas dentro do SLA é menor que 80%% para SP: %.2f%%", percentual*100)
	}

	// Teste sem filtro de estado (considera todos os estados)
	percentualTodosEstados, err := calcularPercentualDentroDoSLA("", tempoDeRespostaMaximo)
	if err != nil {
		t.Fatalf("Erro ao calcular percentual dentro do SLA para todos os estados: %v", err)
	}

	// Verifica se o percentual está acima de 80% para todos os estados
	if percentualTodosEstados < 0.95 {
		t.Errorf("Percentual de respostas dentro do SLA é menor que 80%% para todos os estados: %.2f%%", percentualTodosEstados*100)
	}

	print(percentual*100, "%")
	print(percentualTodosEstados*100, "%")
}

// Função que lê os logs e calcula a porcentagem de respostas dentro do tempo de SLA
func calcularPercentualDentroDoSLA(estado string, tempoDeRespostaMaximo float64) (float64, error) {
	// Abre o arquivo de logs
	file, err := os.Open("logs.json")
	if err != nil {
		return 0, fmt.Errorf("erro ao abrir o arquivo de logs: %v", err)
	}
	defer file.Close()

	// Lê o conteúdo do arquivo
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return 0, fmt.Errorf("erro ao ler o arquivo de logs: %v", err)
	}

	// Carrega os logs no array de LogData
	var logs []LogData
	err = json.Unmarshal(fileContent, &logs)
	if err != nil {
		return 0, fmt.Errorf("erro ao parsear o JSON de logs: %v", err)
	}

	// Variáveis para calcular o percentual
	var totalLogs, logsDentroDoSLA int

	// Itera sobre os logs
	for _, log := range logs {
		// Se o estado não for vazio, filtra os logs pelo estado
		if estado == "" || log.State == estado {
			totalLogs++
			if log.TotalTime <= tempoDeRespostaMaximo {
				logsDentroDoSLA++
			}
		}
	}

	// Verifica se há logs suficientes para calcular
	if totalLogs == 0 {
		return 0, fmt.Errorf("nenhum log encontrado para o estado %s", estado)
	}

	// Calcula a porcentagem
	percentual := float64(logsDentroDoSLA) / float64(totalLogs)

	return percentual, nil
}
