import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { login as apiLogin, logout as apiLogout, getToken } from '../services/api';

// 定义Context中值的类型
interface AuthContextType {
  isAuthenticated: boolean;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
}

// 创建AuthContext
const AuthContext = createContext<AuthContextType | undefined>(undefined);

// 创建一个自定义Hook，方便组件使用AuthContext
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

// 创建AuthProvider组件
interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(false);

  // 在组件加载时检查本地存储中是否已有token
  useEffect(() => {
    const token = getToken();
    if (token) {
      setIsAuthenticated(true);
    }
  }, []);

  const login = async (username: string, password: string) => {
    try {
      await apiLogin(username, password);
      setIsAuthenticated(true);
    } catch (error) {
      console.error('Login failed:', error);
      // 重新抛出错误，以便登录页面可以捕获并显示错误信息
      throw error;
    }
  };

  const logout = () => {
    apiLogout();
    setIsAuthenticated(false);
  };

  const value = {
    isAuthenticated,
    login,
    logout,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
