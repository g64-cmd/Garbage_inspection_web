package services

import (
	"encoding/json"
	"fmt"
	"log"
	"patrol-cloud/internal/models"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

// CommandService 负责将来自 API 的指令发布到 MQTT
type CommandService struct {
	mqttClient mqtt.Client
}

func NewCommandService(client mqtt.Client) *CommandService {
	return &CommandService{mqttClient: client}
}

// SendCommand 遵循 3.1.2 协议发布指令
func (s *CommandService) SendCommand(vehicleID, command, taskID string) (string, error) {
	commandID := uuid.NewString()

	// 1. 构建 Payload
	payload := models.Command{
		CommandID: commandID,
		Command:   command,
		TaskID:    taskID,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal command: %w", err)
	}

	// 2. 定义 Topic
	topic := fmt.Sprintf("vehicles/%s/command", vehicleID)

	// 3. 发布 (QoS 1，如 3.1 所定义)
	token := s.mqttClient.Publish(topic, 1, false, payloadBytes)
	
	// (等待确认不是必须的，但有助于调试)
	if token.Wait() && token.Error() != nil {
		log.Printf("ERROR: Failed to publish command to %s: %v", topic, token.Error())
		return "", token.Error()
	}

	log.Printf("INFO: Command %s published to topic %s", commandID, topic)
	return commandID, nil
}
