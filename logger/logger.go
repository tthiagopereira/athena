package logger

import (
	"bytes"
	"fmt"
	"net/http"
)

type Logger struct {
	ElasticsearchURL string
}

func NewLogger(elasticsearchURL string) *Logger {
	return &Logger{
		ElasticsearchURL: elasticsearchURL,
	}
}

func (l *Logger) SendLog(jsonData []byte) error {
	resp, err := http.Post(l.ElasticsearchURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("Erro ao enviar log para o Elasticsearch. Status code: %d", resp.StatusCode)
	}

	return nil
}
