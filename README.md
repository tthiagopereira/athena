
# Athena

## Descrição
Nomeado após Atena, a deusa da sabedoria, simbolizando a vigilância inteligente do sistema.
Este sistema é responsável por monitorar o desempenho do sistema e enviar alertas quando ocorrem eventos importantes, como alta utilização de recursos ou horários específicos do dia.

## Recursos
- Monitoramento de uso de disco e memória.
- Envio de alertas por SMS e e-mail.
- Funciona em horários específicos (7h da manhã e 17h da tarde).
- Alertas adicionais quando a porcentagem de espaço em disco usado atinge 80%.

## Requisitos de Configuração
- Go

## Como Usar
1. Clone este repositório.
2. Configure as variáveis de ambiente para SMS e e-mail.
3. Execute o sistema usando `go run main.go`.

## Configuração
### Variáveis de Ambiente
- `SMS_API_KEY`: Chave de API para enviar SMS.
- `EMAIL_SMTP_HOST`: Host SMTP para envio de e-mails.
- `EMAIL_SMTP_PORT`: Porta SMTP.

## Autor
[THIAGO PEREIRA DOS SANTOS]
