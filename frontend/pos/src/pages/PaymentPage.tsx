import React, { useState, useEffect, useRef, useMemo, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import {
  Box, Typography, TextField, Alert, List, ListItem, Grid,
  Accordion, AccordionSummary, AccordionDetails, Fab, Zoom, Divider,
  Paper, Chip, Stack, IconButton, InputAdornment, Container, Modal,
  ThemeProvider, createTheme, alpha
} from "@mui/material";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import PaymentIcon from "@mui/icons-material/Payment";
import CloseIcon from "@mui/icons-material/Close";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import ReceiptLongIcon from "@mui/icons-material/ReceiptLong";
import AttachMoneyIcon from "@mui/icons-material/AttachMoney";
import QrCodeIcon from "@mui/icons-material/QrCode";
import CreditCardIcon from "@mui/icons-material/CreditCard";
import ConfettiExplosion from 'react-confetti-explosion';
import axios from "axios";

const theme = createTheme({
  palette: {
    primary: {
      main: '#4caf50',
      light: '#80e27e',
      dark: '#087f23',
      contrastText: '#fff',
    },
    secondary: {
      main: '#ff9800',
      light: '#ffc947',
      dark: '#c66900',
      contrastText: '#000',
    },
    background: {
      default: '#f8f9fa',
      paper: '#ffffff',
    },
  },
  typography: {
    fontFamily: '"Poppins", "Roboto", "Helvetica", "Arial", sans-serif',
    h5: {
      fontWeight: 600,
    },
    h6: {
      fontWeight: 500,
    },
  },
  shape: {
    borderRadius: 12,
  },
  shadows: [
    'none',
    '0px 2px 8px rgba(0,0,0,0.05)',
    '0px 4px 16px rgba(0,0,0,0.08)',
    // Rest of default
    ...Array(22).fill(''),
  ],
  components: {
    MuiPaper: {
      styleOverrides: {
        root: {
          borderRadius: 12,
        },
      },
    },
    MuiChip: {
      styleOverrides: {
        root: {
          fontWeight: 500,
        },
      },
    },
    MuiDivider: {
      styleOverrides: {
        root: {
          borderColor: 'rgba(0, 0, 0, 0.08)',
        },
      },
    },
    MuiTextField: {
      styleOverrides: {
        root: {
          '& .MuiOutlinedInput-root': {
            borderRadius: 8,
          },
        },
      },
    },
    MuiFab: {
      styleOverrides: {
        root: {
          textTransform: 'none',
          boxShadow: '0 4px 12px rgba(76, 175, 80, 0.2)',
          '&:hover': {
            boxShadow: '0 6px 16px rgba(76, 175, 80, 0.4)',
          },
        },
      },
    },
  },
});

// Animasi ketika pembayaran success
const modalStyle = {
  position: 'absolute',
  top: '50%',
  left: '50%',
  transform: 'translate(-50%, -50%)',
  width: '90%',
  maxWidth: 400,
  bgcolor: 'background.paper',
  boxShadow: 24,
  borderRadius: 4,
  p: 4,
  textAlign: 'center',
  outline: 'none',
};

const PaymentPage = () => {
  const navigate = useNavigate();
  const inputRef = useRef(null);
  
  // Get orders from sessionStorage
  const [orders, setOrders] = useState([]);
  const [total, setTotal] = useState(0);
  const [paymentMethod, setPaymentMethod] = useState("cash");
  const [nominalUang, setNominalUang] = useState("");
  const [kembalian, setKembalian] = useState(null);
  const [error, setError] = useState("");
  const [isProcessing, setIsProcessing] = useState(false);
  const [showSuccessModal, setShowSuccessModal] = useState(false);
  const [showConfetti, setShowConfetti] = useState(false);
  
  // Load stored orders
  useEffect(() => {
    const storedOrders = sessionStorage.getItem("orders");
    if (storedOrders) {
      try {
        const parsedOrders = JSON.parse(storedOrders);
        if (parsedOrders && parsedOrders.length > 0) {
          setOrders(parsedOrders);
          const totalPrice = parsedOrders.reduce(
            (sum, item) => sum + item.price * item.quantity,
            0
          );
          setTotal(totalPrice);
        } else {
          // Redirect if no orders
          navigate("/produk");
        }
      } catch (error) {
        console.error("Error parsing orders:", error);
        navigate("/produk");
      }
    } else {
      // Redirect if no orders
      navigate("/produk");
    }
  }, [navigate]);

  // Auto focus on input field
  useEffect(() => {
    if (inputRef.current) {
      inputRef.current.focus();
    }
  }, []);

  // Quick payment options
  const quickPaymentOptions = useMemo(() => {
    // Calculate payment options based on total
    const baseAmount = Math.ceil(total / 10000) * 10000; // Round up to nearest 10k
    return [
      baseAmount,
      baseAmount + 10000,
      baseAmount + 50000,
      baseAmount + 100000
    ];
  }, [total]);
  
  // Handle quick payment option selection
  const handleQuickPayment = useCallback((amount) => {
    setNominalUang(amount.toString());
    // Automatically calculate change when quick payment option selected
    setTimeout(() => {
      if (amount >= total) {
        setError("");
        setKembalian(amount - total);
      } else {
        setError("Nominal uang kurang dari total pembayaran.");
        setKembalian(null);
      }
    }, 100);
  }, [total]);
  
  // Clear input field
  const clearNominalInput = useCallback(() => {
    setNominalUang("");
    setKembalian(null);
    setError("");
    if (inputRef.current) {
      inputRef.current.focus();
    }
  }, []);

  // Handle payment input change
  const handleNominalChange = useCallback((e) => {
    const value = e.target.value;
    // Allow only numbers
    if (/^\d*$/.test(value)) {
      setNominalUang(value);
      
      const numValue = Number(value);
      if (numValue < total) {
        setError("Nominal uang kurang dari total pembayaran.");
        setKembalian(null);
      } else {
        setError("");
        setKembalian(numValue - total);
      }
    }
  }, [total]);

  // Process payment
  const handlePayment = useCallback(async () => {
    if (Number(nominalUang) < total) return;
    
    setIsProcessing(true);
    
    try {
      const payload = {
        items: orders.map((item) => ({
          product_id: item.id,
          quantity: item.quantity,
        })),
        payment_method: paymentMethod,
        paid_amount: Number(nominalUang),
      };
      const token = localStorage.getItem("token");

      await axios.post("http://localhost:3000/payment", payload, {
        headers: {
          Authorization: `Bearer ${token}`,
        },
      });

      // Show success modal with confetti
      setShowSuccessModal(true);
      setShowConfetti(true);
      
      // Redirect after delay
      setTimeout(() => {
        sessionStorage.removeItem("orders");
        navigate("/produk");
      }, 3000);

    } catch (error) {
      console.error("Payment failed:", error);
      setError(
        error.response?.data?.message || 
        "Terjadi kesalahan dalam pemrosesan pembayaran."
      );
      setIsProcessing(false);
    }
  }, [orders, paymentMethod, nominalUang, total, navigate]);

  const handleGoBack = () => {
    navigate(-1);
  };
  
  // Get payment icon based on method
  const getPaymentIcon = (method) => {
    switch (method) {
      case "cash": return <AttachMoneyIcon />;
      case "qris": return <QrCodeIcon />;
      case "debit": return <CreditCardIcon />;
      default: return <PaymentIcon />;
    }
  };
  
  // Success Modal Component
  const SuccessModal = () => (
    <Modal
      open={showSuccessModal}
      aria-labelledby="payment-success-modal"
      aria-describedby="payment-success-description"
    >
      <Box sx={modalStyle}>
        {showConfetti && (
          <Box sx={{ position: 'absolute', top: 0, left: '50%', transform: 'translateX(-50%)' }}>
            <ConfettiExplosion 
              force={0.8}
              duration={3000}
              particleCount={100}
              width={1600}
            />
          </Box>
        )}
        <CheckCircleIcon 
          color="primary" 
          sx={{ 
            fontSize: 80, 
            mb: 2,
            animation: 'scaleIn 0.5s ease-out',
            '@keyframes scaleIn': {
              '0%': { transform: 'scale(0)' },
              '70%': { transform: 'scale(1.2)' },
              '100%': { transform: 'scale(1)' }
            }
          }} 
        />
        <Typography 
          variant="h5" 
          component="h2" 
          id="payment-success-modal"
          color="primary.dark"
          fontWeight="bold"
          sx={{ 
            mb: 2,
            animation: 'fadeInUp 0.5s ease-out 0.3s both',
            '@keyframes fadeInUp': {
              '0%': { opacity: 0, transform: 'translateY(10px)' },
              '100%': { opacity: 1, transform: 'translateY(0)' }
            }
          }}
        >
          Pembayaran Berhasil!
        </Typography>
        <Typography 
          id="payment-success-description" 
          sx={{ 
            mb: 3,
            color: 'text.secondary',
            animation: 'fadeInUp 0.5s ease-out 0.5s both',
            '@keyframes fadeInUp': {
              '0%': { opacity: 0, transform: 'translateY(10px)' },
              '100%': { opacity: 1, transform: 'translateY(0)' }
            }
          }}
        >
          Terima kasih atas pembelian Anda. Pesanan Anda sedang diproses.
        </Typography>
        {kembalian > 0 && (
          <Paper 
            elevation={0}
            sx={{ 
              p: 2, 
              mb: 3, 
              bgcolor: alpha(theme.palette.secondary.light, 0.2),
              border: '1px dashed',
              borderColor: 'secondary.main',
              animation: 'fadeInUp 0.5s ease-out 0.7s both',
            }}
          >
            <Typography variant="body1" fontWeight="bold" color="secondary.dark">
              Kembali: Rp{kembalian.toLocaleString("id-ID")}
            </Typography>
          </Paper>
        )}
        <Typography 
          variant="body2" 
          color="text.secondary"
          sx={{ 
            animation: 'fadeInUp 0.5s ease-out 0.9s both',
          }}
        >
          Anda akan dialihkan dalam beberapa detik...
        </Typography>
      </Box>
    </Modal>
  );

  return (
    <ThemeProvider theme={theme}>
      <Container maxWidth="sm" sx={{ py: 3 }}>
        <Paper 
          elevation={2} 
          sx={{ 
            p: { xs: 2, sm: 3 }, 
            borderRadius: 3,
            background: 'linear-gradient(to bottom, #ffffff, #f9f9f9)',
            position: 'relative',
            overflow: 'hidden',
            '&::before': {
              content: '""',
              position: 'absolute',
              top: 0,
              left: 0,
              right: 0,
              height: '4px',
              background: 'linear-gradient(90deg, #4caf50, #80e27e)',
            }
          }}
        >
          {/* Header */}
          <Box sx={{ display: "flex", alignItems: "center", mb: 2 }}>
            <IconButton 
              onClick={handleGoBack} 
              sx={{ mr: 1, color: 'text.secondary' }}
            >
              <ArrowBackIcon />
            </IconButton>
            <Typography 
              variant="h5" 
              sx={{ 
                flex: 1,
                textAlign: "center", 
                color: "primary.main",
                letterSpacing: 0.5,
              }}
            >
              Pembayaran
            </Typography>
          </Box>

          {/* Order List */}
          <Accordion 
            defaultExpanded 
            sx={{ 
              mb: 2,
              boxShadow: 'none',
              border: '1px solid',
              borderColor: 'divider',
              '&:before': {
                display: 'none',
              },
            }}
          >
            <AccordionSummary 
              expandIcon={<ExpandMoreIcon />}
              sx={{
                backgroundColor: alpha(theme.palette.primary.main, 0.04),
                '&:hover': {
                  backgroundColor: alpha(theme.palette.primary.main, 0.08),
                }
              }}
            >
              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                <ReceiptLongIcon sx={{ mr: 1, color: 'primary.main' }} />
                <Typography variant="h6">
                  Daftar Pesanan ({orders.length} item)
                </Typography>
              </Box>
            </AccordionSummary>
            <AccordionDetails sx={{ p: 0 }}>
              <List disablePadding>
                {orders.map((order, index) => (
                  <React.Fragment key={order.id || index}>
                    <ListItem 
                      sx={{ 
                        display: "flex", 
                        flexDirection: "column", 
                        alignItems: "start", 
                        py: 1.5,
                        px: 2,
                        transition: 'all 0.2s',
                        '&:hover': {
                          backgroundColor: alpha(theme.palette.primary.main, 0.04),
                        }
                      }}
                    >
                      <Typography fontWeight="medium">{order.name}</Typography>
                      <Grid container spacing={1} sx={{ width: "100%", mt: 0.5 }}>
                        <Grid item xs={6}>
                          <Typography variant="body2" color="text.secondary">
                            {order.quantity} x Rp{(order.price || 0).toLocaleString("id-ID")}
                          </Typography>
                        </Grid>
                        <Grid item xs={6} sx={{ textAlign: "right" }}>
                          <Typography fontWeight="bold">
                            Rp{((order.price || 0) * (order.quantity || 0)).toLocaleString("id-ID")}
                          </Typography>
                        </Grid>
                      </Grid>
                      {order.category && (
                        <Chip
                          label={order.category}
                          size="small"
                          variant="outlined"
                          sx={{ mt: 0.5, height: 20, fontSize: '0.7rem' }}
                        />
                      )}
                    </ListItem>
                    {index < orders.length - 1 && <Divider />}
                  </React.Fragment>
                ))}
              </List>
              <Box sx={{ p: 2, bgcolor: alpha(theme.palette.primary.main, 0.03) }}>
                <Grid container spacing={1}>
                  <Grid item xs={6}>
                    <Typography variant="body1">Total Harga:</Typography>
                  </Grid>
                  <Grid item xs={6}>
                    <Typography variant="body1" align="right" fontWeight="medium">
                      Rp{total.toLocaleString("id-ID")}
                    </Typography>
                  </Grid>
                  
                  <Grid item xs={12}><Divider sx={{ my: 1 }} /></Grid>
                  
                  <Grid item xs={6}>
                    <Typography variant="h6" fontWeight="bold">
                      Total Pembayaran:
                    </Typography>
                  </Grid>
                  <Grid item xs={6}>
                    <Typography variant="h6" align="right" fontWeight="bold" color="primary.dark">
                      Rp{total.toLocaleString("id-ID")}
                    </Typography>
                  </Grid>
                </Grid>
              </Box>
            </AccordionDetails>
          </Accordion>

          {/* Payment Method Section */}
          <Paper 
            elevation={0} 
            sx={{ 
              p: 2, 
              mb: 2, 
              bgcolor: alpha(theme.palette.primary.main, 0.04),
              border: '1px solid',
              borderColor: alpha(theme.palette.primary.main, 0.1),
            }}
          >
            <Typography variant="subtitle1" gutterBottom fontWeight="medium">
              Metode Pembayaran:
            </Typography>
            <Stack 
              direction="row" 
              spacing={1} 
              sx={{ 
                flexWrap: "wrap", 
                gap: 1 
              }}
            >
              <Chip
                icon={<AttachMoneyIcon />}
                label="Tunai"
                color="primary"
                variant={paymentMethod === "cash" ? "filled" : "outlined"}
                onClick={() => setPaymentMethod("cash")}
                sx={{ 
                  fontSize: "0.9rem", 
                  py: 2.5, 
                  borderRadius: 2,
                  fontWeight: 500,
                  transition: "all 0.2s",
                  '&:hover': {
                    transform: 'translateY(-2px)',
                    boxShadow: 1
                  }
                }}
              />
              <Chip
                icon={<QrCodeIcon />}
                label="QRIS"
                color="primary"
                variant={paymentMethod === "qris" ? "filled" : "outlined"}
                onClick={() => setPaymentMethod("qris")}
                sx={{ 
                  fontSize: "0.9rem", 
                  py: 2.5, 
                  borderRadius: 2,
                  fontWeight: 500,
                  transition: "all 0.2s",
                  '&:hover': {
                    transform: 'translateY(-2px)',
                    boxShadow: 1
                  }
                }}
              />
              <Chip
                icon={<CreditCardIcon />}
                label="Kartu Debit"
                color="primary"
                variant={paymentMethod === "debit" ? "filled" : "outlined"}
                onClick={() => setPaymentMethod("debit")}
                sx={{ 
                  fontSize: "0.9rem", 
                  py: 2.5, 
                  borderRadius: 2,
                  fontWeight: 500,
                  transition: "all 0.2s",
                  '&:hover': {
                    transform: 'translateY(-2px)',
                    boxShadow: 1
                  }
                }}
              />
            </Stack>
          </Paper>

          {/* Quick Payment Options */}
          <Paper 
            elevation={0} 
            sx={{ 
              p: 2, 
              mb: 2, 
              bgcolor: alpha(theme.palette.primary.main, 0.04),
              border: '1px solid',
              borderColor: alpha(theme.palette.primary.main, 0.1),
            }}
          >
            <Typography variant="subtitle1" gutterBottom fontWeight="medium">
              Pilih Nominal Cepat:
            </Typography>
            <Stack 
              direction="row" 
              spacing={1} 
              sx={{ 
                flexWrap: "wrap", 
                gap: 1 
              }}
            >
              {quickPaymentOptions.map((amount) => (
                <Chip
                  key={amount}
                  label={`Rp${amount.toLocaleString("id-ID")}`}
                  color="primary"
                  variant={nominalUang === amount.toString() ? "filled" : "outlined"}
                  onClick={() => handleQuickPayment(amount)}
                  sx={{ 
                    fontSize: "0.9rem", 
                    py: 2.5, 
                    borderRadius: 2,
                    fontWeight: 500,
                    transition: "all 0.2s",
                    '&:hover': {
                      transform: 'translateY(-2px)',
                      boxShadow: 1
                    }
                  }}
                />
              ))}
            </Stack>
          </Paper>

          {/* Input Nominal Uang */}
          <TextField
            fullWidth
            label="Masukkan Nominal Uang"
            variant="outlined"
            type="text"
            inputMode="numeric"
            value={nominalUang}
            onChange={handleNominalChange}
            inputRef={inputRef}
            InputProps={{
              startAdornment: 
               <InputAdornment position="start">Rp</InputAdornment>,
              endAdornment: nominalUang && (
                <InputAdornment position="end">
                  <IconButton onClick={clearNominalInput} edge="end">
                    <CloseIcon />
                  </IconButton>
                </InputAdornment>
              )
            }}
            sx={{ 
              mb: 2,
              '& .MuiOutlinedInput-root': {
                borderRadius: 2,
                '&.Mui-focused': {
                  '& fieldset': {
                    borderColor: 'primary.main',
                    borderWidth: 2,
                  },
                },
              },
            }}
          />

          {/* Error message */}
          {error && (
            <Alert 
              severity="error" 
              sx={{ 
                mb: 2,
                borderRadius: 2,
                animation: 'fadeIn 0.3s ease-in-out',
                '@keyframes fadeIn': {
                  '0%': {
                    opacity: 0,
                    transform: 'translateY(-10px)'
                  },
                  '100%': {
                    opacity: 1,
                    transform: 'translateY(0)'
                  }
                }
              }}
            >
              {error}
            </Alert>
          )}
          
          {/* Kembalian alert */}
          {kembalian !== null && kembalian >= 0 && (
            <Alert 
              severity="success" 
              sx={{ 
                mb: 2,
                borderRadius: 2,
                animation: 'fadeIn 0.3s ease-in-out',
                '@keyframes fadeIn': {
                  '0%': {
                    opacity: 0,
                    transform: 'translateY(-10px)'
                  },
                  '100%': {
                    opacity: 1,
                    transform: 'translateY(0)'
                  }
                }
              }}
            >
              <Typography variant="body1" fontWeight="medium">
                Kembali: Rp{kembalian.toLocaleString("id-ID")}
              </Typography>
            </Alert>
          )}
          
          {/* Payment button */}
          <Box 
            sx={{ 
              mt: 2, 
              display: "flex", 
              justifyContent: "center"
            }}
          >
            <Zoom 
              in={Number(nominalUang) >= total} 
              style={{ transitionDelay: Number(nominalUang) >= total ? '100ms' : '0ms' }}
            >
              <Fab 
                color="primary" 
                variant="extended" 
                onClick={handlePayment} 
                disabled={isProcessing || Number(nominalUang) < total}
                size="large"
                sx={{
                  px: 4,
                  py: 3,
                  fontSize: '1rem',
                  position: 'relative',
                  overflow: 'hidden',
                  '&::after': isProcessing ? {
                    content: '""',
                    position: 'absolute',
                    top: 0,
                    left: 0,
                    width: '30%',
                    height: '100%',
                    background: 'rgba(255, 255, 255, 0.3)',
                    animation: 'loading 1s infinite',
                    '@keyframes loading': {
                      '0%': { transform: 'translateX(-100%)' },
                      '100%': { transform: 'translateX(400%)' }
                    }
                  } : {},
                  '&:hover': {
                    transform: 'translateY(-2px) scale(1.05)',
                    transition: 'all 0.2s ease'
                  }
                }}
              >
                <PaymentIcon sx={{ mr: 1, fontSize: '1.4rem' }} />
                {isProcessing ? "Memproses..." : "Bayar Sekarang"}
              </Fab>
            </Zoom>
          </Box>
        </Paper>
        
        {/* Success Modal */}
        <SuccessModal />
      </Container>
    </ThemeProvider>
  );
};

export default PaymentPage;