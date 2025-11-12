import axios from 'axios';

// 从环境变量或默认值中获取API基础URL
const API_BASE_URL = process.env.REACT_APP_API_URL || '/api/v1';

// 创建一个Axios实例
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 添加一个请求拦截器，用于在每个请求的Authorization头中附加JWT
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('jwt_token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

/**
 * 登录函数
 * @param username 用户名
 * @param password 密码
 * @returns Promise，包含token
 */
export const login = async (username: string, password: string) => {
  const response = await apiClient.post('/auth/login', { username, password });
  if (response.data && response.data.token) {
    localStorage.setItem('jwt_token', response.data.token);
  }
  return response.data;
};

/**
 * 登出函数
 */
export const logout = () => {
  localStorage.removeItem('jwt_token');
  // 这里可以根据需要添加其他清理逻辑
};

/**
 * 获取当前存储的token
 * @returns token字符串或null
 */
export const getToken = () => {
  return localStorage.getItem('jwt_token');
};

/**
 * 获取车辆列表
 * @returns Promise，包含车辆列表
 */
export const getVehicles = async () => {
  const response = await apiClient.get('/vehicles');
  return response.data;
};

/**
 * 根据ID获取单个车辆信息
 * @param id 车辆ID
 * @returns Promise，包含车辆信息
 */
export const getVehicleById = async (id: string) => {
  const response = await apiClient.get(`/vehicles/${id}`);
  return response.data;
};

/**
 * 根据车辆ID获取决策日志
 * @param vehicleId 车辆ID
 * @returns Promise，包含决策日志列表
 */
export const getDecisionLogsByVehicleId = async (vehicleId: string) => {
  // 注意：API路由是 /vehicles/:id/decision-logs
  const response = await apiClient.get(`/vehicles/${vehicleId}/decision-logs`);
  return response.data;
};

/**
 * 获取所有决策日志
 * @returns Promise，包含所有决策日志列表
 */
export const getAllDecisionLogs = async () => {
  const response = await apiClient.get('/decision-logs');
  return response.data;
};


export default apiClient;
