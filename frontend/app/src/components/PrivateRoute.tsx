import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const PrivateRoute: React.FC = () => {
  const { isAuthenticated } = useAuth();

  // 如果已认证，则渲染子路由（通过Outlet）
  // 否则，重定向到登录页面
  return isAuthenticated ? <Outlet /> : <Navigate to="/login" replace />;
};

export default PrivateRoute;
