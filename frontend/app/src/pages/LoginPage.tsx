import React, { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';
import {
  Container,
  TextField,
  Button,
  Typography,
  Box,
  Alert,
} from '@mui/material';

const LoginPage: React.FC = () => {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const { login } = useAuth();
  const navigate = useNavigate();

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    setError(null); // 重置错误信息

    if (!username || !password) {
      setError('Username and password are required.');
      return;
    }

    try {
      await login(username, password);
      navigate('/'); // 登录成功后跳转到主页
    } catch (err: any) {
      // Enhanced error logging
      console.error('Login attempt failed:', err);
      if (err.response) {
        // The request was made and the server responded with a status code
        console.error('Server Response Data:', err.response.data);
        setError(err.response.data.error || 'Login failed. Please check your username and password.');
      } else if (err.request) {
        // The request was made but no response was received
        console.error('No response from server:', err.request);
        setError('Login failed: No response from server.');
      } else {
        // Something happened in setting up the request that triggered an Error
        console.error('Request setup error:', err.message);
        setError('Login failed: An unexpected error occurred during request setup.');
      }
    }
  };

  return (
    <Container component="main" maxWidth="xs">
      <Box
        sx={ {
          marginTop: 8,
          display: 'flex',
          flexDirection: 'column',
          alignItems: 'center',
        } }
      >
        <Typography component="h1" variant="h5">
          Sign In
        </Typography>
        <Box component="form" onSubmit={handleSubmit} noValidate sx={ { mt: 1 } }>
          <TextField
            margin="normal"
            required
            fullWidth
            id="username"
            label="Username"
            name="username"
            autoComplete="username"
            autoFocus
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <TextField
            margin="normal"
            required
            fullWidth
            name="password"
            label="Password"
            type="password"
            id="password"
            autoComplete="current-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          {error && (
            <Alert severity="error" sx={ { width: '100%', mt: 2 } }>
              {error}
            </Alert>
          )}
          <Button
            type="submit"
            fullWidth
            variant="contained"
            sx={ { mt: 3, mb: 2 } }
          >
            Sign In
          </Button>
        </Box>
      </Box>
    </Container>
  );
};

export default LoginPage;
