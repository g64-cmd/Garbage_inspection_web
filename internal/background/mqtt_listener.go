package background

import (
	"encoding/json"
	"log"
	"patrol-cloud/internal/models"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// MQTTListener 遵循 4.2.6 设计
type MQTTListener struct {
	Client     mqtt.Client
	HubChannel chan<- []byte // (只写通道，推向 TelemetryHub)
}

func NewMQTTListener(client mqtt.Client, hubChannel chan<- []byte) *MQTTListener {
	return &MQTTListener{
		Client:     client,
		HubChannel: hubChannel,
	}
}

// StartListening 订阅主题并启动监听
func (l *MQTTListener) StartListening() {
	// (订阅通配符主题，如 4.2.6 所述)
	const topic = "vehicles/+/status"
	
	if token := l.Client.Subscribe(topic, 1, l.onStatusMessage); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to subscribe to MQTT topic %s: %v", topic, token.Error())
	}
	log.Printf("INFO: MQTTListener subscribed to topic: %s", topic)
}

// onStatusMessage 是 4.2.6 的 _onStatusMessage 实现
func (l *MQTTListener) onStatusMessage(client mqtt.Client, msg mqtt.Message) {
	log.Printf("DEBUG: Received MQTT message on topic: %s", msg.Topic())

	// 1. 解析 vehicle_id (如 3.3.2 所需)
	topicParts := strings.Split(msg.Topic(), "/")
	if len(topicParts) < 3 {
		log.Printf("WARN: Received message on unexpected topic: %s", msg.Topic())
		return
	}
	vehicleID := topicParts[1]

	// 2. 反序列化 Payload (来自 3.1.1)
	var status models.VehicleStatus
	if err := json.Unmarshal(msg.Payload(), &status); err != nil {
		log.Printf("WARN: Failed to unmarshal status from %s: %v", vehicleID, err)
		return
	}

	// 3. 构建遥测消息 (3.3.2 格式)
	updateMsg := models.TelemetryUpdate{
		VehicleID:     vehicleID,
		VehicleStatus: status,
	}

	// 4. 序列化并推送到 Hub Channel
	updateBytes, err := json.Marshal(updateMsg)
	if err != nil {
		log.Printf("ERROR: Failed to marshal telemetry update: %v", err)
		return
	}

	l.HubChannel <- updateBytes
}
