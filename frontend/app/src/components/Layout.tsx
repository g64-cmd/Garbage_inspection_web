import React, { ReactNode } from 'react';
import {
  AppBar,
  Toolbar,
  Typography,
  Button,
  Box,
  Container,
} from '@mui/material';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';

interface LayoutProps {
  children: ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const { logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <Box sx={ { display: 'flex' } }>
      <AppBar position="fixed">
        <Toolbar>
          <Typography variant="h6" component="div" sx={ { flexGrow: 1 } }>
            Garbage Inspection Platform
          </Typography>
          <Button color="inherit" onClick={handleLogout}>
            Logout
          </Button>
        </Toolbar>
      </AppBar>
      <Container component="main" sx={ { mt: 10, flexGrow: 1 } }>
        {children}
      </Container>
    </Box>
  );
};

export default Layout;
