package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Shopify/sarama"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	kafkaBrokers := []string{"3.136.87.228:9094"}
	topic := "alert"

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionNone
	config.Producer.Flush.Frequency = 500 * time.Millisecond

	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(kafkaBrokers, config)
	if err != nil {
		fmt.Printf("Erro ao criar o produtor Kafka: %s\n", err)
		return
	}
	defer producer.Close()

	elasticsearchURL := "http://3.139.239.25:9200/log_logshare_kafka/_doc/"

	for {
		now := time.Now()
		hour := now.Hour()
		min := now.Minute()

		diskInfo, _ := disk.Usage("/")
		usedPercent := diskInfo.UsedPercent

		data := collectData()
		if err := sendLog(data, elasticsearchURL); err != nil {
			fmt.Printf("Erro ao enviar dados para o Elasticsearch: %s\n", err)
		}
		if hour == 12 && min == 58 || hour == 17 && min == 55 || (usedPercent > 80) {
			if err := sendToKafka(data, producer, topic); err != nil {
				fmt.Printf("Erro ao enviar dados para o Kafka: %s\n", err)
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

func collectData() map[string]interface{} {
	memInfo, _ := mem.VirtualMemory()
	diskInfo, _ := disk.Usage("/")
	loadInfo, _ := load.Avg()

	data := map[string]interface{}{
		"espaco_total_disco_gb":          fmt.Sprintf("%.2f", float64(diskInfo.Total)/1024/1024/1024),
		"espaco_usado_disco_gb":          fmt.Sprintf("%.2f", float64(diskInfo.Used)/1024/1024/1024),
		"memoria_total_gb":               fmt.Sprintf("%.2f", float64(memInfo.Total)/1024/1024/1024),
		"memoria_usada_gb":               fmt.Sprintf("%.2f", float64(memInfo.Used)/1024/1024/1024),
		"porcentagem_espaco_usado_disco": fmt.Sprintf("%.2f", diskInfo.UsedPercent),
		"porcentagem_memoria_usada":      fmt.Sprintf("%.2f", memInfo.UsedPercent),
		"utilizacao_cpu_1_min":           fmt.Sprintf("%.2f", loadInfo.Load1),
		"timestamp":                      time.Now().UTC().Format(time.RFC3339),
		"type":                           fmt.Sprintf("Apache Kafka"),
	}

	return data
}

func sendLog(data map[string]interface{}, elasticsearchURL string) error {

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Erro ao converter para JSON: %s", err)
	}

	resp, err := http.Post(elasticsearchURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("Erro ao enviar JSON para o Elasticsearch: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Erro ao enviar JSON para o Elasticsearch. Status code: %d", resp.StatusCode)
	}

	fmt.Println("JSON enviado com sucesso para o Elasticsearch!")
	return nil
}

func sendToKafka(data map[string]interface{}, producer sarama.SyncProducer, topic string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Erro ao converter para JSON: %s", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(jsonData),
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("Erro ao enviar mensagem para o Kafka: %s", err)
	}

	fmt.Println("Mensagem enviada com sucesso para o Kafka!")
	return nil
}
