import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { AuthProvider, useAuth } from './context/AuthContext';
import LoginPage from './pages/LoginPage';
import DashboardPage from './pages/DashboardPage';
import VehicleDetailPage from './pages/VehicleDetailPage'; // Import the new page
import PrivateRoute from './components/PrivateRoute';
import Layout from './components/Layout';

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          {/* Public Route: Login Page */}
          <Route path="/login" element={<LoginPage />} />

          {/* Private Routes: All routes under this element are protected and use the main layout */}
          <Route element={<Layout><PrivateRoute /></Layout>}>
            <Route path="/dashboard" element={<DashboardPage />} />
            <Route path="/vehicles/:id" element={<VehicleDetailPage />} /> {/* Add the new route */}
            {/* Add other private routes here as the application grows */}
          </Route>

          {/* Root path redirect logic */}
          {/* If you land on "/", it will redirect to /dashboard if logged in, or /login if not */}
          <Route path="*" element={<RootRedirect />} />
        </Routes>
      </Router>
    </AuthProvider>
  );
}

// This component handles the redirection logic from the root path or any unmatched path.
const RootRedirect: React.FC = () => {
  const { isAuthenticated } = useAuth();
  return <Navigate to={isAuthenticated ? "/dashboard" : "/login"} replace />;
};

export default App;
