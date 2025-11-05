# 项目状态与路线图 (截至 2025年11月6日)

本文档旨在总结当前项目的完成状态，并为下一阶段的开发工作提供清晰的路线图和约定。

## 1. 今日完成的工作

今天，我们主要围绕**提升系统健壮性**和**扩展核心功能**两个方面展开工作，取得了以下关键进展：

### 1.1 健壮性修复

- **配置管理**: 彻底移除了所有硬编码配置。现在所有配置（如数据库连接、JWT密钥、MinIO凭证、LLM密钥等）都通过 `internal/config` 包从环境变量加载，使应用配置更安全、更灵活。
- **WebSocket 安全**: 修复了 WebSocket 的 `CheckOrigin` 漏洞，现在只有在环境变量 `WEBSOCKET_ALLOWED_ORIGINS` 中指定的域名才能建立连接，防止了跨站 WebSocket 劫持攻击。
- **开发环境依赖隔离**: 使用 Go 的构建标签（Build Tags）将测试用户创建逻辑 (`createTempUser`) 从生产代码中分离，确保生产构建的纯净性。现在只有在执行 `go build -tags dev` 时，该功能才会被编译进去。
- **后台任务持久性**: 为 `DecisionService` 中的异步任务增加了文件队列作为“安全网”。如果数据库日志记录失败，相关数据将被保存到 `failed_tasks.log` 文件中，防止数据丢失。
- **自动化测试**: 为项目引入了 `testify` 测试框架，并为核心的 `AuthService` 编写了单元测试，修复了所有相关的依赖和构建问题，确保测试流程可以顺利运行。

### 1.2 功能实现

- **Qwen LLM 集成**: 重构了 `LLMService`，使其能够与通义千问（Qwen）的 OpenAI 兼容 API 对接。相关的 API 端点 `/api/v1/llm/plan` 也已实现并启用。
- **数据库持久层扩展**:
    - 使用 `golang-migrate` 的思想，在 `migrations` 目录下创建了管理 `users`, `decision_logs`, `vehicles`, `vehicle_telemetry` 表结构的 SQL 脚本。
    - 扩展了数据模型 (`models.go`) 和数据访问层 (`repository.go`)，增加了对新表的增删改查支持。
    - 更新了 `MQTTListener`，现在它可以在接收到设备遥测数据的同时，将其异步存入数据库。
- **后端 API 扩展**:
    - 创建了 `VehicleHandler`, `TelemetryHandler`, `LogHandler`，以提供前端所需的车辆列表、历史轨迹和决策日志等数据。
    - 在 `router.go` 中注册了所有新的 API 路由。
- **前端项目初始化**:
    - 在根目录创建了 `frontend` 文件夹，并使用 `create-react-app` 初始化了名为 `app` 的 React (TypeScript) 项目。
    - 安装了前端开发所需的全部核心依赖，包括 `axios`, `react-router-dom`, `leaflet`, `chart.js`, 以及 `Material-UI`。
    - 对默认模板进行了清理。

## 2. 后续工作规划 (下一步)

**目标**: 继续完成前端应用的开发，实现一个功能完善、界面美观、可打包成原生 App 的原型。

**核心步骤**:

1.  **搭建前端认证流程**:
    -   **API 服务层**: 创建一个 `axios` 实例，封装与后端 API 的所有通信。
    -   **认证服务**: 实现一个 `authService`，负责调用登录接口、在 `localStorage` 中存取 JWT。
    -   **全局状态管理**: 使用 React Context (`AuthContext`) 在整个应用中管理用户的登录状态。
    -   **私有路由**: 创建一个 `PrivateRoute` 组件，用于保护需要登录才能访问的页面。

2.  **实现核心页面组件**:
    -   **登录页 (`LoginPage.tsx`)**: 使用 Material-UI 构建一个包含用户名和密码输入框的登录表单。
    -   **主控台 (`DashboardPage.tsx`)**: 作为登录后的主界面，将包含车辆列表等核心视图。
    -   **车辆详情页 (`VehicleDetailPage.tsx`)**: 这是最复杂的页面，将包含：
        -   **实时地图**: 使用 Leaflet 和 WebSocket，实时展示车辆位置。
        -   **历史轨迹**: 调用后端 API，在地图上绘制历史轨迹。
        -   **数据图表**: 使用 Chart.js 展示电量等历史数据。
        -   **日志列表**: 分页展示决策日志和相关图片。

3.  **完善整体布局**: 创建通用的布局组件，包括导航栏、侧边栏等。

4.  **打包成原生应用**: 在 Web 功能完成后，使用 **Capacitor** 将其打包成 Android 和 iOS 应用。

## 3. 前后端接口约定

以下是当前后端已准备好，供前端调用的 API 接口：

- **认证**: `POST /api/v1/auth/login`
- **车辆**: 
    - `GET /api/v1/vehicles`
    - `GET /api/v1/vehicles/:id`
- **历史轨迹**: `GET /api/v1/vehicles/:id/telemetry?start_time=<RFC3339>&end_time=<RFC3339>`
- **决策日志**: `GET /api/v1/vehicles/:id/decision-logs?page=<number>&pageSize=<number>`
- **LLM 规划**: `POST /api/v1/llm/plan`
- **发送指令**: `POST /api/v1/commands/send`
- **实时遥测 (WebSocket)**: `GET /ws/telemetry?token=<jwt_token>`
