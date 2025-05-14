import { useState, useEffect } from 'react';
import { 
  TextField, 
  Button, 
  Container, 
  Typography, 
  Box,
  Fade,
  CircularProgress,
  IconButton,
  InputAdornment,
  Alert,
  Collapse,
  Paper
} from '@mui/material';
import { 
  Visibility, 
  VisibilityOff,
  PointOfSale,
  Login as LoginIcon,
  Email
} from '@mui/icons-material';
import axios from 'axios';
import { useAuth } from '../context/AuthContext';
import { useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';

const CashierLogin = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [capsLockWarning, setCapsLockWarning] = useState(false);
  const navigate = useNavigate();
  const { login } = useAuth();

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.getModifierState('CapsLock')) {
      setCapsLockWarning(true);
    } else {
      setCapsLockWarning(false);
    }

    if (e.key === 'Enter') {
      handleLogin();
    }
  };

  const handleLogin = async () => {
    if (isLoading) return;
    
    setIsLoading(true);
    setError('');
    
    try {
      const res = await axios.post('http://127.0.0.1:3000/login', { email, password });
      const { token, user } = res.data;
      
      // Success animation delay
      await new Promise(resolve => setTimeout(resolve, 500));
      
      login(user, token);
      navigate('/dashboard');
    } catch (err: any) {
      setError(err.response?.data?.message || 'Login failed. Please check your credentials.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <Container maxWidth="xs" sx={{ 
      height: '90vh',
      display: 'flex',
      flexDirection: 'column',
      justifyContent: 'center',
      alignItems: 'center',
      background: 'linear-gradient(135deg, #f5f7fa 0%, #e4e8ed 100%)'
    }}>
      <motion.div
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.3 }}
      >
        <Paper elevation={3} sx={{ 
          width: '100%',
          p: 4,
          borderRadius: 6,
          borderTop: '4px solid',
          borderColor: 'primary.main',
          backgroundColor: 'background.paper'
        }}>
          <Box sx={{ 
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'center',
            mb: 3
          }}>
            <PointOfSale sx={{ 
              fontSize: 50,
              color: 'primary.main',
              mb: 1
            }} />
            <Typography component="h1" variant="h5" sx={{ 
              fontWeight: 'bold',
              color: 'primary.main'
            }}>
            </Typography>
            <Typography variant="body2" color="text.secondary">
            </Typography>
          </Box>
          
          <Box component="div" sx={{ width: '100%' }} onKeyDown={handleKeyDown}>
            <TextField
              label="Email / Username"
              variant="outlined"
              fullWidth
              autoFocus
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              error={!!error}
              sx={{ mb: 2 }}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <Email color="action" />
                  </InputAdornment>
                ),
              }}
            />
            
            <TextField
              label="Password"
              type={showPassword ? 'text' : 'password'}
              variant="outlined"
              fullWidth
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              error={!!error}
              sx={{ mb: 1 }}
              InputProps={{
                startAdornment: (
                  <InputAdornment position="start">
                    <LoginIcon color="action" />
                  </InputAdornment>
                ),
                endAdornment: (
                  <InputAdornment position="end">
                    <IconButton
                      aria-label="toggle password visibility"
                      onClick={() => setShowPassword(!showPassword)}
                      edge="end"
                      size="small"
                    >
                      {showPassword ? <VisibilityOff fontSize="small" /> : <Visibility fontSize="small" />}
                    </IconButton>
                  </InputAdornment>
                ),
              }}
            />

            <Collapse in={capsLockWarning}>
              <Alert severity="warning" sx={{ mb: 2, fontSize: '0.8rem' }} icon={false}>
                Caps Lock is on
              </Alert>
            </Collapse>
            
            <Collapse in={!!error}>
              <Alert severity="error" sx={{ mb: 2 }} icon={false}>
                {error}
              </Alert>
            </Collapse>
            
            <Button
              fullWidth
              variant="contained"
              color="primary"
              onClick={handleLogin}
              disabled={isLoading || !email || !password}
              sx={{ 
                py: 1.5,
                borderRadius: 1,
                fontWeight: 'bold',
                textTransform: 'uppercase',
                letterSpacing: '0.5px',
                '&:active': {
                  transform: 'scale(0.98)'
                },
                transition: 'transform 0.2s'
              }}
              startIcon={isLoading ? <CircularProgress size={20} color="inherit" /> : null}
            >
              {isLoading ? 'Processing...' : 'Login'}
            </Button>
          </Box>
        </Paper>
      </motion.div>
    </Container>
  );
};

export default CashierLogin;