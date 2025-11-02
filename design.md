详细设计说明书 (DDS)

端云协同巡检平台

V1.0

1. 引言

1.1 目的

本文档为“端云协同巡检平台”项目提供详细的软件设计规范。它在《软件需求说明书 (SRS)》(srs_v2.md) 和《概要设计说明书 (SDD)》(sdd_v1.md) 的基础上，对系统内的所有子系统、模块、接口、协议和数据结构进行精确定义。

本文档是编码实现的直接依据。

1.2 范围

本设计覆盖三大子系统及其所有内部模块和外部接口：

边缘端子系统 (Python 3)

云端子系统 (Go)

客户端子系统 (Vue.js 3 + Capacitor)

1.3 参考文件

《软件需求说明书 (SRS) - V2.0》(srs_v2.md)

《概要设计说明书 (SDD) - V1.0》(sdd_v1.md)

《巡检车控制算法开发进展报告》(005.pdf)

docker-compose.yml (草案)

2. 系统体系结构

系统采用“边缘-云-客户”三端架构，通过两种通信协议（MQTT, HTTP/S）和三种核心交互流（A, B, C）解耦。

+--------------------------------+
|      客户端 (Vue.js 3)         |
|  (Web / Android / iOS)         |
|--------------------------------|
|  [交互 A] HTTP POST (指令)    | <-----> |
|  [交互 B] WebSocket (实时状态) | <-----> |
+--------------------------------+       |
         ^                             |
(HTTPS / WSS)                        |
         |                             |
+--------v-----------------------------+
|        云端子系统 (Go)              |
| (Gin, WebSocket, Paho MQTT, ONNX)  |
|------------------------------------|
|   [交互 A] MQTT 发布 (指令)          | <---+
|   [交互 B] MQTT 订阅 (状态)          | <---+-----> [EMQX Broker]
|   [交互 C] HTTP 响应 (AI决策)        | <---+
+------------------------------------+
         ^                             |
         |                             |
 (MQTT / HTTPS)                      |
         |                             |
+--------v-----------------------------+
|      边缘端子系统 (Python 3)        |
| (Coordinator + NetworkService)     |
|------------------------------------|
|  [交互 A] MQTT 订阅 (指令)           |
|  [交互 B] MQTT 发布 (状态)           |
|  [交互 C] HTTP POST (决策请求)       |
+------------------------------------+


3. 接口控制文档 (ICD) - 协议与约定

本节严格定义所有跨子系统的通信接口。

3.1 边缘端 (Python) ⟷ 云端 (Go) (MQTT 协议)

中间件: EMQX (云端部署)

协议: MQTT v5.0

QoS: 所有交互必须使用 QoS 1，以确保在网络不稳定时消息至少送达一次 (NF-2)。

Clean Session: False (边缘端重连时，云端和边缘端应能恢复未确认的消息)。

3.1.1 状态上报 (交互 B: 边缘端 -> 云端)

主题 (Topic): vehicles/{vehicle_id}/status

方向: 边缘端 (Publish) -> 云端 (Subscribe)

Payload (JSON): application/json

{
  "timestamp": 1678886400, // Unix 时间戳 (秒)
  "position": {
    "lat": 39.91234,
    "lng": 116.39745
  },
  "battery": 85.5, // 电池百分比
  "state": "NAVIGATING" // 状态枚举: IDLE, PLANNING, NAVIGATING, OPERATING, AWAITING_CONFIRMATION, ERROR
}


3.1.2 宏观指令 (交互 A: 云端 -> 边缘端)

主题 (Topic): vehicles/{vehicle_id}/command

方向: 云端 (Publish) -> 边缘端 (Subscribe)

Payload (JSON): application/json

{
  "command_id": "uuid-cmd-12345", // 用于边缘端去重
  "command": "START_AUTONOMY", // 枚举: START_AUTONOMY, EMERGENCY_STOP, RESUME_PATH
  "task_id": "task-abc-123" // (可选) 任务标识
}


3.2 边缘端 (Python) ⟷ 云端 (Go) (HTTP 协议)

3.2.1 同步决策 (交互 C: 边缘端 -> 云端)

接口: POST /api/v1/decisions/recognize

方向: 边缘端 (Request) -> 云端 (Response)

约束: 边缘端发起请求时必须设置客户端总超时 timeout=5.0 秒 (REQ-E-3)。

Request:

Content-Type: multipart/form-data

Body:

image (File): 裁切后的图片二进制数据 (e.g., image.jpg)。

metadata (Text): 描述元数据的 JSON 字符串。

{
  "vehicle_id": "v-001",
  "timestamp": 1678886401
}


Response (Success): 200 OK

Content-Type: application/json

Body (JSON):

{
  "image_id": "uuid-img-12345", // 云端生成的此事件ID
  "action": "pickup", // 枚举: "pickup", "abandon"
  "confidence": 0.95, // 本地 ONNX 模型对该决策的置信度
  "reason": "is_trash_type_A" // (可选) 决策原因
}


Response (Failure):

400 Bad Request: (e.g., image 或 metadata 缺失)。

500 Internal Server Error: (e.g., AI 模型加载失败)。

3.3 客户端 (Vue) ⟷ 云端 (Go) (HTTP/S & WSS)

3.3.1 HTTP/S RESTful API

认证: 除 /login 外，所有接口必须在 Header 中携带 Authorization: Bearer {JWT} (REQ-S-2)。

用户认证

接口: POST /api/v1/auth/login

Request (JSON): {"username": "admin", "password": "secure_password_123"}

Response (JSON, 200 OK): {"token": "eyJhbGciOiJILET..."}

宏观指令 (交互 A)

接口: POST /api/v1/commands/send (REQ-S-3)

Request (JSON):

{
  "vehicle_id": "v-001",
  "command": "START_AUTONOMY" // 同 3.1.2 中 `command` 枚举
}


Response (JSON, 202 Accepted):

{"status": "queued", "command_id": "uuid-cmd-12345"}


LLM 智能交互

接口: POST /api/v1/llm/plan

Request (JSON): {"prompt": "为A区规划一条巡检路线"}

Response (JSON, 200 OK): {"plan_id": "uuid-plan-789", "steps": ["...", "..."]}

3.3.2 WebSocket (WSS) 实时遥测

接口: GET /ws/telemetry (REQ-S-6)

协议: HTTP/1.1 升级到 WebSocket。

消息 (Server -> Client):

事件名称: (如果使用 Socket.IO) telemetry_update / (如果原生 WebSocket) 消息本身。

方向: 云端 (Push) -> 客户端 (Listen)。

Payload (JSON): 必须与 3.1.1 状态上报的 Payload 格式完全一致。

{
  "vehicle_id": "v-001", // (云端应附加此字段，以便客户端区分)
  "timestamp": 1678886400,
  "position": { ... },
  "battery": 85.5,
  "state": "NAVIGATING"
}


4. 子系统详细设计

4.1 边缘端子系统 (Python 3)

系统层次结构: 严格遵循 [cite: 8-92] 的星状架构。main.py 启动 Coordinator 和所有 Services (作为 threading.Thread)，并传入 Coordinator 的 queue.Queue 实例。

4.1.1 Coordinator (Class)

文件: main_coordinator.py

职责: 核心状态机，系统的“大脑”。

属性:

state: (Enum) IDLE, NAVIGATING, AWAITING_CONFIRMATION, OPERATING, ERROR

main_queue: (queue.Queue) 接收所有服务发来的消息。

service_queues: (dict) 存储指向每个服务队列的引用 (e.g., {"motion": motion_q, "network": network_q})。

核心方法:

run(self): 主循环。msg = self.main_queue.get()，switch msg.type。

handle_navigating(self, msg): 处理 NAVIGATING 状态下的消息。

if msg.type == 'TARGET_FOUND':

if msg.confidence < LOW_THRESHOLD:

self._send_to_service("motion", Message("STOP"))

self._send_to_service("network", Message("REQUEST_DECISION", payload=msg.payload))

self.state = AWAITING_CONFIRMATION

else: (处理高置信度目标) ...

handle_awaiting_confirmation(self, msg): (REQ-E-4, REQ-E-5)

if msg.type == 'DECISION_COMPLETE' and msg.payload['action'] == 'pickup':

self._send_to_service("manipulation", Message("EXECUTE_PICKUP", ...))

self.state = OPERATING

if msg.type == 'DECISION_FAILED' or (msg.type == 'DECISION_COMPLETE' and msg.payload['action'] == 'abandon'):

self._log_failure(msg.payload)

self._send_to_service("motion", Message("RESUME_PATH"))

self.state = NAVIGATING

_send_to_service(self, service_name, msg): 辅助函数，self.service_queues[service_name].put(msg)。

4.1.2 NetworkService (Class)

文件: network_service.py (继承 BaseService)

职责: 作为唯一的网络网关，处理 交互 A、B、C。

属性:

coordinator_queue: (queue.Queue)

self_queue: (queue.Queue) 自身的任务队列。

mqtt_client: (paho.mqtt.client)

http_session: (requests.Session)

核心方法:

run(self): 循环 msg = self.self_queue.get()。

if msg.type == 'REQUEST_DECISION': self._handle_decision_request(msg.payload)

if msg.type == 'PUBLISH_STATUS': self.publish_status(msg.payload)

setup(self): (在线程启动时调用)

self.http_session = requests.Session()

self._init_mqtt()

_init_mqtt(self):

self.mqtt_client = paho.mqtt.client.Client(...)

self.mqtt_client.on_message = self._on_mqtt_message

self.mqtt_client.connect(...)

self.mqtt_client.subscribe(f"vehicles/{self.vehicle_id}/command", qos=1)

self.mqtt_client.loop_start()

_on_mqtt_message(self, client, userdata, msg): [交互 A]

data = json.loads(msg.payload)

self.coordinator_queue.put(Message("COMMAND_RECEIVED", payload=data))

publish_status(self, status_payload): [交互 B]

self.mqtt_client.publish(f"vehicles/{self.vehicle_id}/status", payload=json.dumps(status_payload), qos=1)

_handle_decision_request(self, payload): [交互 C]

image = payload['image']

metadata = json.dumps({"vehicle_id": self.vehicle_id, ...})

files = {'image': image, 'metadata': metadata}

try:

response = self.http_session.post(DECISION_URL, files=files, timeout=5.0)

response.raise_for_status()

data = response.json()

self.coordinator_queue.put(Message("DECISION_COMPLETE", payload=data))

except (requests.exceptions.Timeout, requests.exceptions.RequestException) as e:

self.coordinator_queue.put(Message("DECISION_FAILED", payload={"reason": str(e)}))

4.2 云端子系统 (Go)

系统层次结构: main.go 启动 api (Gin), background (Goroutines), services (业务逻辑), db (数据访问)。

4.2.1 模块: main.go

职责: 应用入口，组装所有依赖。

实现:

db_pool = db.InitPostgreSQL(PG_DSN)

repo = db.NewRepository(db_pool)

storage = storage.InitMinIO(...)

mqtt_client = background.InitMQTTClient(EMQX_HOST, ...)

telemetry_hub = services.NewTelemetryHub()

go telemetry_hub.Run()

ai_service = services.NewAIService(ONNX_MODEL_PATH)

llm_service = services.NewLLMService(LLM_API_KEY)

... (初始化其他 services, 注入依赖)

listener = background.NewMQTTListener(mqtt_client, telemetry_hub.BroadcastChannel)

go listener.StartListening()

router = api.SetupRouter(auth_service, command_service, ...)

router.Run(":8080")

4.2.2 模块: api/*_handler.go (Gin Handlers)

职责: HTTP/WebSocket 接口，负责数据校验和调用 Service。

DecisionHandler(decision_service *services.DecisionService):

image, _ := c.FormFile("image")

metadata_str := c.Form("metadata")

json.Unmarshal([]byte(metadata_str), &metadata)

image_bytes = ... (从 image 读取)

result, err := decision_service.ProcessDecision(image_bytes, metadata)

if err != nil { c.JSON(500, ...); return }

c.JSON(200, result) (响应体必须符合 3.2.1)

WebSocketHandler(telemetry_hub *services.TelemetryHub):

conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

telemetry_hub.RegisterClient(conn)

(go conn.ReadLoop(...) -> 触发 hub.UnregisterClient)

4.2.3 模块: services/decision_service.go

职责: 编排 交互 C，解耦 API 和 AI。

Struct: DecisionService { AIService *AIService; Repo *db.Repository; Storage *storage.MinIOClient }

Method: ProcessDecision(image []byte, metadata models.Metadata) (*models.DecisionResult, error)

result, err := s.AIService.Recognize(image)

if err != nil { return nil, err }

go func() {
  // 异步存日志和图片，并处理潜在错误
  imageURL, err := s.Storage.Upload(image)
  if err != nil { log.Printf("ERROR: failed to upload image: %v", err) }

  err = s.Repo.LogDecision(result, imageURL, ...)
  if err != nil { log.Printf("ERROR: failed to log decision: %v", err) }
}()

return result, nil

4.2.4 模块: services/ai_service.go

职责: 封装 CGo + ONNX Runtime 调用 (REQ-S-9)。

Struct: AIService { onnx_session ... }

Method: Recognize(image []byte) (*models.DecisionResult, error)

tensor_input = PreProcess(image)

tensor_output = C.run_onnx_inference(self.onnx_session, tensor_input)

result = PostProcess(tensor_output) (将 tensor 转换为 {"action": "pickup", ...})

return result, nil

4.2.5 模块: services/llm_service.go

职责: 封装外部 LLM API 调用。

Struct: LLMService { http_client *http.Client; api_key string }

Method: PlanMission(prompt string) (*models.LLMPlan, error)

req_body = json.Build(...)

http_req, _ = http.NewRequest("POST", LLM_API_URL, req_body)

http_req.Header.Set("Authorization", "Bearer "+s.api_key)

http_resp, _ := s.http_client.Do(http_req)

... (解析 http_resp)

return plan, nil

4.2.6 模块: background/mqtt_listener.go

职责: 交互 B 的数据源，常驻 Goroutine。

Struct: MQTTListener { Client paho.Client; HubChannel chan<- []byte }

Method: StartListening()

self.Client.Subscribe("vehicles/+/status", 1, self._onStatusMessage)

Method: _onStatusMessage(client paho.Client, msg paho.Message)

// TODO: 附加 vehicle_id (从 msg.Topic() 中解析)

// e.g., msg_with_id = ...

self.HubChannel <- msg.Payload() (将消息推入 Go Channel)

4.2.7 模块: services/telemetry_service.go

职责: 交互 B 的广播中枢。

Struct: TelemetryHub

Clients: map[*websocket.Conn]bool (并发安全)

BroadcastChannel: chan []byte (从 MQTTListener 接收)

Register: chan *websocket.Conn

Unregister: chan *websocket.Conn

Method: Run(): 核心 Goroutine，必须使用 select 语句

select { ... }

case conn := <- h.Register: h.Clients[conn] = true

case conn := <- h.Unregister: ... delete(h.Clients, conn)

case message := <- h.BroadcastChannel:

for conn := range h.Clients { conn.WriteMessage(websocket.TextMessage, message) } (并发写入)

4.3 客户端子系统 (Vue 3 + Capacitor)

系统层次结构: main.js (初始化 Vue, Pinia, Router) -> App.vue -> router.js -> views/* -> components/*。services/ 和 stores/ 作为全局单例。

4.3.1 模块: services/http.js (Axios)

职责: 封装 axios 实例，必须实现 JWT 认证拦截器 (REQ-C-6)。

实现:

const apiClient = axios.create({ baseURL: "/api/v1" })

apiClient.interceptors.request.use(config => { ... } )

const authStore = useAuthStore()

if (authStore.token) { config.headers.Authorization = \Bearer ${authStore.token}` }`

return config

export default apiClient

4.3.2 模块: services/websocket.service.js

职责: 封装 socket.io-client，处理 交互 B。

实现: (单例模式)

let socket;

connect():

const authStore = useAuthStore()

socket = io(WS_URL, { auth: { token: authStore.token } })

const vehicleStore = useVehicleStore()

socket.on("telemetry_update", (data) => { vehicleStore.updateVehicleStatus(data) }) (REQ-C-4)

disconnect(): if (socket) socket.disconnect()

4.3.3 模块: stores/auth.store.js (Pinia)

职责: 存储用户状态和 JWT。

State: token: string | null (从 localStorage 初始化), user: object | null

Actions:

login(username, password):

const { token } = await authService.login(username, password)

this.token = token

localStorage.setItem("token", token)

router.push("/dashboard")

logout(): this.token = null; localStorage.removeItem("token"); ...

4.3.4 模块: stores/vehicle.store.js (Pinia)

职责: 存储所有车辆的实时状态，驱动 UI 响应式更新。

State: vehicles: Record<string, VehicleState> (e.g., {"v-001": { position: ..., battery: ... }})

Actions:

updateVehicleStatus(data): [核心] (REQ-C-4)

const id = data.vehicle_id

if (!this.vehicles[id]) { this.vehicles[id] = {} }

this.vehicles[id].position = data.position

this.vehicles[id].state = data.state

...

4.3.5 模块: views/VehicleDetailView.vue (Vue Component)

职责: 展示单个车辆的实时信息。

实现:

<script setup>
import { onMounted } from 'vue'
import { useVehicleStore } from '@/stores/vehicle.store'
import { websocketService } from '@/services/websocket.service'
import VehicleMap from '@/components/VehicleMap.vue'

const vehicleStore = useVehicleStore()
const vehicleId = 'v-001' // (应从路由参数获取)

// 订阅 WebSocket
onMounted(() => {
  websocketService.connect()
})
</script>

<template>
  <!-- Pinia 状态 store.vehicles[vehicleId] 改变，UI 自动更新 -->
  <VehicleMap :position="vehicleStore.vehicles[vehicleId]?.position" />
  <p>状态: {{ vehicleStore.vehicles[vehicleId]?.state }}</p>
</template>


5. 数据库详细设计 (PostgreSQL)

(见 sdd_v1.md 5. 数据库概要设计)

users

id (uuid, PK), username (varchar, unique, not null), hashed_password (varchar, not null)

vehicles

id (varchar, PK), name (varchar, not null), model (varchar), current_status (jsonb)

decision_logs

id (uuid, PK, default gen_random_uuid()), vehicle_id (varchar, not null, FK -> vehicles.id), timestamp (timestamptz, not null, default now()), image_url (varchar), server_decision (varchar), vehicle_action (varchar), request_metadata (jsonb)

patrol_logs

...

6. 部署与环境

基础设施: docker-compose.yml 用于启动 PostgreSQL, EMQX, MinIO。

云端部署:

Dockerfile: Go 应用的 Dockerfile 必须基于 nvidia/cuda:*-cudnn*-devel-* 镜像，以包含 GPU 驱动和 CGo 编译链。

Dockerfile: 必须 COPY onnxruntime.so 等库文件到镜像中。

Docker Compose: 生产 docker-compose.yml 中，Go 服务必须配置 deploy.resources.reservations.devices 以便访问 NVIDIA GPU。

边缘端部署: Python 3 环境（带 requests, paho-mqtt）必须预装在巡检车的操作系统上。

客户端部署:

Web: npm run build 生成的 dist/ 目录，由 Nginx 或 Gin 提供静态文件服务。

App: npx cap sync 同步 dist/ 目录到 android/ 和 ios/，然后使用 Android Studio / Xcode 编译打包。