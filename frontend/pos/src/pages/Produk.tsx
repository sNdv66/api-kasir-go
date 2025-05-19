import React, { useEffect, useState, useCallback } from "react";
import axios from "axios";
import { useAuth } from "../context/AuthContext";
import { 
  Container, 
  Box, 
  Grid, 
  Card, 
  CardMedia, 
  CardContent, 
  InputBase, 
  Paper, 
  IconButton, 
  Typography,
  useMediaQuery,
  useTheme,
  Modal,
  Button,
  Divider,
  Badge,
  Skeleton,
  Drawer,
  Fade,
  Chip,
  Snackbar,
  Alert,
  Backdrop,
  CircularProgress,
  Fab
} from "@mui/material";
import { styled, alpha } from '@mui/material/styles';
import SearchIcon from "@mui/icons-material/Search";
import NoPhotographyIcon from "@mui/icons-material/NoPhotography";
import ShoppingCartIcon from "@mui/icons-material/ShoppingCart";
import CloseIcon from "@mui/icons-material/Close";
import AddIcon from "@mui/icons-material/Add";
import RemoveIcon from "@mui/icons-material/Remove";
import SaveIcon from "@mui/icons-material/Save";
import DeleteOutlineIcon from "@mui/icons-material/DeleteOutline";
import ShoppingBagIcon from "@mui/icons-material/ShoppingBag";
import ArrowBackIcon from "@mui/icons-material/ArrowBack";
import { motion } from "framer-motion";
import { createContext, useContext , ReactNode } from "react";
import { useNavigate } from "react-router-dom";


// Define interfaces
interface Product {
  id: string;
  name: string;
  price: number;
  stock: number;
  image_url?: string;
}

interface OrderItem extends Product {
  quantity: number;
}

// Custom styled components
const SearchBar = styled('div')(({ theme }) => ({
  position: 'relative',
  borderRadius: theme.shape.borderRadius,
  backgroundColor: alpha(theme.palette.common.white, 0.15),
  '&:hover': {
    backgroundColor: alpha(theme.palette.common.white, 0.25),
  },
  marginRight: theme.spacing(2),
  marginLeft: 0,
  width: '100%',
  [theme.breakpoints.up('sm')]: {
    width: '100%',
  },
  boxShadow: '0 2px 8px rgba(0,0,0,0.08)',
  border: '1px solid',
  borderColor: alpha(theme.palette.common.black, 0.05),
}));

const SearchIconWrapper = styled('div')(({ theme }) => ({
  padding: theme.spacing(0, 1),
  height: '100%',
  position: 'absolute',
  pointerEvents: 'none',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
  color: alpha(theme.palette.common.black, 0.54),
}));

const StyledInputBase = styled(InputBase)(({ theme }) => ({
  color: 'inherit',
  width: '100%',
  '& .MuiInputBase-input': {
    padding: theme.spacing(1.5, 1, 1.5, 0),
    // vertical padding + font size from searchIcon
    paddingLeft: `calc(1em + ${theme.spacing(4)})`,
    transition: theme.transitions.create('width'),
    width: '100%',
  },
}));

const ProductCard = styled(motion.div)(({ theme }) => ({
  height: '100%',
  cursor: 'pointer',
  borderRadius: theme.shape.borderRadius,
  overflow: 'hidden',
  boxShadow: '0 4px 12px rgba(0,0,0,0.05)',
  transition: 'all 0.3s ease',
  '&:hover': {
    transform: 'translateY(-4px)',
    boxShadow: '0 10px 20px rgba(0,0,0,0.08)',
  },
  '&:active': {
    transform: 'scale(0.98)',
  },
  background: theme.palette.background.paper,
}));

const CartFab = styled(Fab)(({ theme }) => ({
  position: 'fixed',
  bottom: theme.spacing(4),
  right: theme.spacing(4),
  zIndex: 999,
  boxShadow: '0 4px 14px rgba(0,0,0,0.15)',
}));

// Main component
const Products = () => {
  const { user, token } = useAuth();
  const userId = user?.id;
  const branchId = user?.branch_id || "";
  const [products, setProducts] = useState<Product[]>([]);
  const [filtered, setFiltered] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);
  const [search, setSearch] = useState("");
  const [orders, setOrders] = useState<OrderItem[]>([]);
  const [openOrderDrawer, setOpenOrderDrawer] = useState(false);
  const [notification, setNotification] = useState({ open: false, message: '', type: 'success' as 'success' | 'error' });
  const [processingOrder, setProcessingOrder] = useState(false);
  const [savingOrder, setSavingOrder] = useState(false);
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('sm'));
  const isMedium = useMediaQuery(theme.breakpoints.down('md'));
  const navigate = useNavigate();
  
  // Calculate grid columns based on viewport size
  const gridCols = isMobile ? 5 : isMedium ? 4 : 3;
  
  // Memoized fetch products function
  const fetchProducts = useCallback(async () => {
    if (!token || !branchId) return;

    try {
      setLoading(true);
      const res = await axios.get(
        `http://127.0.0.1:3000/branches/${branchId}/products`,
        { headers: { Authorization: `Bearer ${token}` } }
      );
      setProducts(res.data);
      setFiltered(res.data);
    } catch (err) {
      console.error("Failed to fetch products:", err);
      setNotification({
        open: true,
        message: 'Gagal memuat produk. Silakan coba lagi.',
        type: 'error'
      });
    } finally {
      setLoading(false);
    }
  }, [branchId, token]);

  // Handle search
  const handleSearch = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value.toLowerCase();
    setSearch(value);
    setFiltered(products.filter(p => p.name.toLowerCase().includes(value)));
  };

  // Cart management functions
  const addToOrder = (product: Product) => {
    if (product.stock <= 0) {
      setNotification({
        open: true,
        message: 'Stok produk habis!',
        type: 'error'
      });
      return;
    }
    
    setOrders(prevOrders => {
      const existingItem = prevOrders.find(item => item.id === product.id);
      
      if (existingItem) {
        if (existingItem.quantity >= product.stock) {
          setNotification({
            open: true,
            message: 'Jumlah melebihi stok yang tersedia!',
            type: 'error'
          });
          return prevOrders;
        }
        
        const updatedOrders = prevOrders.map(item =>
          item.id === product.id 
            ? { ...item, quantity: item.quantity + 1 } 
            : item
        );
        
        setNotification({
          open: true,
          message: `${product.name} ditambahkan ke keranjang`,
          type: 'success'
        });
        
        return updatedOrders;
      } else {
        setNotification({
          open: true,
          message: `${product.name} ditambahkan ke keranjang`,
          type: 'success'
        });
        
        return [...prevOrders, { ...product, quantity: 1 }];
      }
    });
  };
  
  const removeFromOrder = (productId: string) => {
    setOrders(prevOrders => prevOrders.filter(item => item.id !== productId));
  };
  
  const updateQuantity = (productId: string, newQuantity: number) => {
    if (newQuantity < 1) {
      removeFromOrder(productId);
      return;
    }
    
    const product = products.find(p => p.id === productId);
    
    if (product && newQuantity > product.stock) {
      setNotification({
        open: true,
        message: 'Jumlah melebihi stok yang tersedia!',
        type: 'error'
      });
      return;
    }
    
    setOrders(prevOrders =>
      prevOrders.map(item =>
        item.id === productId ? { ...item, quantity: newQuantity } : item
      )
    );
 };
  

  const clearOrders = () => {
    setOrders([]);
    setNotification({
      open: true,
      message: 'Keranjang telah dikosongkan',
      type: 'success'
    });
  };

  // Price formatter
  const formatPrice = (price: number) => {
    return new Intl.NumberFormat("id-ID", {
      style: "currency",
      currency: "IDR",
      minimumFractionDigits: 0,
    }).format(price);
  };
  
  // Calculate total price
  const calculateTotal = () => {
    return orders.reduce((sum, item) => sum + (item.price * item.quantity), 0);
  };

  // Calculate total items
  const totalItems = orders.reduce((sum, item) => sum + item.quantity, 0);
  // handleProcessOrdercons
  const handleProcessOrder = () => {
  if (orders.length === 0) {
    setNotification({
      open: true,
      message: 'Keranjang masih kosong!',
      type: 'warning'
    });
    return;
  }

  setProcessingOrder(true);

  // Simpan orders ke sessionStorage
  sessionStorage.setItem("orders", JSON.stringify(orders));

  // Pindah ke halaman payment
  navigate("/payment");
};
  // Handle save order
  const handleSaveOrder = async () => {
  if (!user) {
    setNotification({
      open: true,
      message: 'User belum login!',
      type: 'error'
    });
    return;
  }

  const payload = {
    branch_id: user.branch_id,
    user_id: user.id,
    note: '',
    orders: orders.map(item => ({
      product_id: item.id,
      quantity: item.quantity,
      subtotal: item.quantity * item.price
    }))
  };

  console.log('Payload yang dikirim ke backend:', payload);

  try {
    await axios.post('http://127.0.0.1:3000/pending-orders', payload, {
      headers: {
        Authorization: `Bearer ${token}`
      }
    });
    navigate('/pending-orders');
  } catch (err) {
    console.error('Gagal menyimpan pesanan:', err.response?.data || err.message);
    setNotification({
      open: true,
      message: 'Gagal menyimpan pesanan',
      type: 'error'
    });
  } finally {
    setSavingOrder(false);
  }
};


  // Handle close notification
  const handleCloseNotification = () => {
    setNotification({ ...notification, open: false });
  };

  // Fetch products on mount
  useEffect(() => {
    fetchProducts();
  }, [fetchProducts]);

  return (
    <Container 
      maxWidth="lg" 
      sx={{ 
        py: { xs: 1, sm: 2 }, 
        px: { xs: 1, sm: 2 },
        height: '110vh', 
        display: 'flex', 
        flexDirection: 'column' 
      }}
    >
      {/* Header with search */}
      <Box 
        sx={{ 
          display: 'flex', 
          flexDirection: { xs: 'column', sm: 'row' },
          alignItems: 'center', 
          mb: 1, 
          gap: 1 
        }}
      >
        <Box sx={{ width: '100%', display: 'flex', alignItems: 'center' }}>
          <Typography 
            variant="h5" 
            component="h1" 
            sx={{ 
              fontWeight: 'bold', 
              mr: 2,
              display: { xs: 'none', md: 'block' }
            }}
          >
            Produk
          </Typography>
          
          <SearchBar sx={{ flexGrow: 1 }}>
            <SearchIconWrapper>
              <SearchIcon />
            </SearchIconWrapper>
            <StyledInputBase
              placeholder="Cari produk..."
              value={search}
              onChange={handleSearch}
              fullWidth
            />
          </SearchBar>
        </Box>
        
        {!isMobile && (
          <Badge 
            badgeContent={totalItems} 
            color="secondary"
            sx={{ 
              '& .MuiBadge-badge': { 
                fontSize: 10, 
                height: 18, 
                minWidth: 18 
              }
            }}
          >
            <Button
              variant="contained"
              color="secondary"
              startIcon={<ShoppingCartIcon />}
              onClick={() => setOpenOrderDrawer(true)}
              disabled={totalItems === 0}
              sx={{ 
                borderRadius: 2,
                px: 2,
                boxShadow: 2,
                whiteSpace: 'nowrap'
              }}
            >
              {totalItems > 0 ? `${totalItems} item` : "Keranjang"}
            </Button>
          </Badge>
        )}
      </Box>

      {/* Products section */}
      <Box 
        sx={{ 
          flexGrow: 1, 
          overflowY: 'auto',
          px: { xs: 0, sm: 0 },
          pb: { xs: 8, sm: 2 } // Extra padding at bottom for mobile (for the FAB)
        }}
      >
        {loading ? (
          <Grid container spacing={1.5}>
            {[...Array(12)].map((_, i) => (
              <Grid item xs={gridCols} key={i}>
                <Card sx={{ height: '100%', boxShadow: '0 2px 8px rgba(0,0,0,0.05)' }}>
                  <Skeleton variant="rectangular" height={120} />
                  <CardContent sx={{ p: 1.5 }}>
                    <Skeleton width="80%" height={20} />
                    <Box sx={{ display: 'flex', justifyContent: 'space-between', mt: 1 }}>
                      <Skeleton width="40%" height={20} />
                      <Skeleton width="20%" height={20} />
                    </Box>
                  </CardContent>
                </Card>
              </Grid>
            ))}
          </Grid>
        ) : filtered.length === 0 ? (
          <Box 
            sx={{ 
              display: 'flex', 
              flexDirection: 'column',
              justifyContent: 'center', 
              alignItems: 'center', 
              height: '50vh',
              textAlign: 'center',
              color: 'text.secondary',
              p: 3
            }}
          >
            <SearchIcon sx={{ fontSize: 48, mb: 2, opacity: 0.5 }} />
            <Typography variant="h6">Tidak ada produk ditemukan</Typography>
            <Typography variant="body2" sx={{ mt: 1, maxWidth: 300 }}>
              Coba ubah kata kunci pencarian atau periksa kembali koneksi internet Anda
            </Typography>
            <Button 
              variant="outlined" 
              sx={{ mt: 3 }}
              onClick={() => {
                setSearch('');
                setFiltered(products);
              }}
            >
              Tampilkan semua produk
            </Button>
          </Box>
        ) : (
          <Grid container spacing={1.5}>
            {filtered.map((product) => (
              <Grid item xs={gridCols} key={product.id}>
                <ProductCard
                  whileHover={{ y: -4 }}
                  whileTap={{ scale: 0.98 }}
                  onClick={() => addToOrder(product)}
                >
                  <Card sx={{ height: '100%', boxShadow: 'none' }}>
                    <CardMedia
                      component="div"
                      sx={{ 
                        height: 120,
                        width: 140,
                        bgcolor: 'grey.50', 
                        display: 'flex', 
                        alignItems: 'center', 
                        justifyContent: 'center',
                        position: 'relative'
                      }}
                    >
                      {product.image_url ? (
                        <img
                          src={product.image_url}
                          alt={product.name}
                          style={{ 
                            width: '100%', 
                            height: '100%', 
                            objectFit: 'cover',
                            transition: 'transform 0.3s ease'
                          }}
                        />
                      ) : (
                        <NoPhotographyIcon sx={{ fontSize: 40, color: 'grey.300' }} />
                      )}
                      
                      {product.stock <= 0 && (
                        <Box sx={{ 
                          position: 'absolute', 
                          top: 0, 
                          left: 0, 
                          right: 0, 
                          bottom: 0,
                          bgcolor: 'rgba(0,0,0,0.5)',
                          display: 'flex',
                          alignItems: 'center',
                          justifyContent: 'center'
                        }}>
                          <Typography 
                            sx={{ 
                              color: 'white', 
                              fontWeight: 'bold',
                              textTransform: 'uppercase',
                              fontSize: '0.75rem',
                              letterSpacing: 1
                            }}
                          >
                            Stok Habis
                          </Typography>
                        </Box>
                      )}
                      
                      {product.stock > 0 && product.stock <= 5 && (
                        <Chip 
                          label={`Stok: ${product.stock}`}
                          size="small"
                          color="warning"
                          sx={{ 
                            position: 'absolute',
                            top: 8,
                            right: 8,
                            fontSize: '0.625rem',
                            height: 20
                          }}
                        />
                      )}
                    </CardMedia>
                    <CardContent sx={{ p: 1.5 }}>
                      <Typography 
                        variant="subtitle2" 
                        component="div" 
                        noWrap
                        sx={{ fontWeight: 'medium' }}
                      >
                        {product.name}
                      </Typography>
                      <Box sx={{ 
                        display: 'flex', 
                        justifyContent: 'space-between',
                        alignItems: 'center',
                        mt: 0.5
                      }}>
                        <Typography 
                          variant="body2" 
                          color="primary" 
                          fontWeight="bold"
                        >
                          {formatPrice(product.price)}
                        </Typography>
                        
                        {orders.some(item => item.id === product.id) && (
                          <Chip 
                            size="small"
                            label={orders.find(item => item.id === product.id)?.quantity || 0}
                            color="primary"
                            sx={{ 
                              height: 20, 
                              minWidth: 20,
                              fontSize: '0.625rem'
                            }}
                          />
                        )}
                      </Box>
                    </CardContent>
                  </Card>
                </ProductCard>
              </Grid>
            ))}
          </Grid>
        )}
      </Box>

      {/* Mobile floating cart button */}
      {isMobile && (
        <CartFab 
          color="primary" 
          aria-label="cart"
          onClick={() => setOpenOrderDrawer(true)}
          disabled={totalItems === 0}
        >
          <Badge 
            badgeContent={totalItems} 
            color="error"
            sx={{ '& .MuiBadge-badge': { fontSize: 10 } }}
          >
            <ShoppingCartIcon />
          </Badge>
        </CartFab>
      )}

      {/* Cart Drawer */}
      <Drawer
        anchor={isMobile ? 'bottom' : 'right'}
        open={openOrderDrawer}
        onClose={() => setOpenOrderDrawer(false)}
        PaperProps={{
          sx: {
            width: isMobile ? '99%' : 400,
            maxWidth: '99%',
            borderRadius: isMobile ? '16px 16px 0 0' : 0,
            maxHeight: isMobile ? '85vh' : '100vh',
          }
        }}
      >
        <Box sx={{ 
          display: 'flex', 
          flexDirection: 'column',
          height: '100%'
        }}>
          {/* Cart Header */}
          <Box sx={{ 
            p: 2, 
            display: 'flex', 
            justifyContent: 'space-between', 
            alignItems: 'center',
            borderBottom: '1px solid',
            borderColor: 'divider'
          }}>
            <Box sx={{ display: 'flex', alignItems: 'center' }}>
              {isMobile && (
                <IconButton edge="start" onClick={() => setOpenOrderDrawer(false)} sx={{ mr: 1 }}>
                  <ArrowBackIcon />
                </IconButton>
              )}
              <Typography variant="h6" sx={{ display: 'flex', alignItems: 'center' }}>
                <ShoppingBagIcon sx={{ mr: 1 }} /> 
                Pesanan {totalItems > 0 && `(${totalItems})`}
              </Typography>
            </Box>
            
            {!isMobile && (
              <IconButton onClick={() => setOpenOrderDrawer(false)}>
                <CloseIcon />
              </IconButton>
            )}
          </Box>
          
          {/* Cart Items */}
          <Box sx={{ flexGrow: 1, overflowY: 'auto', p: 0 }}>
            {orders.length === 0 ? (
              <Box sx={{ 
                display: 'flex', 
                flexDirection: 'column', 
                alignItems: 'center', 
                justifyContent: 'center', 
                height: '100%',
                textAlign: 'center',
                color: 'text.secondary',
                p: 4
              }}>
                <ShoppingCartIcon sx={{ fontSize: 64, mb: 2, opacity: 0.3 }} />
                <Typography variant="h6">Keranjang Kosong</Typography>
                <Typography variant="body2" sx={{ mt: 1, maxWidth: 300 }}>
                  Tambahkan produk dengan mengklik produk di halaman utama
                </Typography>
                <Button 
                  variant="outlined" 
                  sx={{ mt: 3 }}
                  onClick={() => setOpenOrderDrawer(false)}
                >
                  Belanja Sekarang
                </Button>
              </Box>
            ) : (
              <Box sx={{ p: 0 }}>
                {orders.map((item, index) => (
                  <Fade key={item.id} in={true} timeout={300} style={{ transitionDelay: `${index * 50}ms` }}>
                    <Box>
                      <Box sx={{ p: 2 }}>
                        <Box sx={{ display: 'flex' }}>
                          <Box 
                            sx={{ 
                              width: 60,
                              height: 60,
                              borderRadius: 1,
                              bgcolor: 'grey.100',
                              mr: 2,
                              overflow: 'hidden',
                              flexShrink: 0,
                              display: 'flex',
                              alignItems: 'center',
                              justifyContent: 'center'
                            }}
                          >
                            {item.image_url ? (
                              <img
                                src={item.image_url}
                                alt={item.name}
                                style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                              />
                            ) : (
                              <NoPhotographyIcon sx={{ fontSize: 24, color: 'grey.400' }} />
                            )}
                          </Box>
                          
                          <Box sx={{ flexGrow: 1 }}>
                            <Box sx={{ display: 'flex', justifyContent: 'space-between' }}>
                              <Typography variant="subtitle2" fontWeight="medium" noWrap sx={{ maxWidth: '70%' }}>
                                {item.name}
                              </Typography>
                              <Typography variant="subtitle2" fontWeight="bold">
                                {formatPrice(item.price * item.quantity)}
                              </Typography>
                            </Box>
                            
                            <Typography variant="caption" color="text.secondary">
                              {formatPrice(item.price)} per item
                            </Typography>
                            
                            <Box sx={{ display: 'flex', alignItems: 'center', mt: 1, justifyContent: 'space-between' }}>
                              <Box sx={{ 
                                display: 'flex',
                                alignItems: 'center',
                                border: '1px solid',
                                borderColor: 'divider',
                                borderRadius: 1,
                                overflow: 'hidden'
                              }}>
                                <IconButton 
                                  size="small" 
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    updateQuantity(item.id, item.quantity - 1);
                                  }}
                                  sx={{ p: 0.5 }}
                                >
                                  <RemoveIcon fontSize="small" />
                                </IconButton>
                                
                                <Typography sx={{ px: 2, userSelect: 'none' }}>
                                  {item.quantity}
                                </Typography>
                                
                                <IconButton 
                                  size="small"
                                  onClick={(e) => {
                                    e.stopPropagation();
                                    updateQuantity(item.id, item.quantity + 1);
                                  }}
                                  sx={{ p: 0.5 }}
                                >
                                  <AddIcon fontSize="small" />
                                </IconButton>
                              </Box>
                              
                              <IconButton
                                size="small"
                                color="error"
                                onClick={(e) => {
                                  e.stopPropagation();
                                  removeFromOrder(item.id);
                                }}
                              >
                                <DeleteOutlineIcon fontSize="small" />
                              </IconButton>
                            </Box>
                          </Box>
                        </Box>
                      </Box>
                      <Divider />
                    </Box>
                  </Fade>
                ))}
                
                {orders.length > 0 && (
                  <Box sx={{ p: 2, textAlign: 'right' }}>
                    <Button
                      variant="text"
                      color="error"
                      size="small"
                      startIcon={<DeleteOutlineIcon />}
                      onClick={clearOrders}
                    >
                      Kosongkan Keranjang
                    </Button>
                  </Box>
                )}
              </Box>
            )}
          </Box>
          
          {/* Cart Footer with Total & Actions */}
          {orders.length > 0 && (
            <Box sx={{ 
              p: 2, 
              borderTop: '1px solid', 
              borderColor: 'divider',
              backgroundColor: 'background.paper'
            }}>
              <Box sx={{ 
                display: 'flex', 
                justifyContent: 'space-between', 
                mb: 2,
                p: 1.5,
                bgcolor: alpha(theme.palette.primary.main, 0.05),
                borderRadius: 1
              }}>
                <Typography variant="subtitle1" fontWeight="bold">Total Pembayaran</Typography>
                <Typography variant="subtitle1" fontWeight="bold" color="primary">
                  {formatPrice(calculateTotal())}
                </Typography>
              </Box>
              
              <Grid container spacing={1}>
                <Grid item xs={6}>
                <Button 
                fullWidth 
                variant="outlined"
                color="secondary"
                size="large"
                onClick={handleSaveOrder}
                disabled={processingOrder || savingOrder}
                startIcon={<SaveIcon />}
                sx={{ borderRadius: 2 }}
                 >
                 {savingOrder ? <CircularProgress size={24} /> : "Simpan"}
                </Button>
                </Grid>
                <Grid item xs={6}>
                  <Button 
                    fullWidth 
                    variant="contained" 
                    color="primary"
                    size="large"
                    onClick={handleProcessOrder}
                    disabled={processingOrder}
                    sx={{ borderRadius: 2 }}
                  >
                    {processingOrder ? <CircularProgress size={24} color="inherit" /> : "Proses Pembayaran"}
                  </Button>
                </Grid>
              </Grid>
            </Box>
          )}
        </Box>
      </Drawer>
      
      {/* Notifications */}
      <Snackbar
        open={notification.open}
        autoHideDuration={300}
        onClose={handleCloseNotification}
        anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
      >
        <Alert 
          onClose={handleCloseNotification} 
          severity={notification.type}
          variant="filled"
          sx={{ width: '99%' }}
        >
          {notification.message}
        </Alert>
      </Snackbar>
      
      {/* Loading backdrop for processing actions */}
      <Backdrop
        sx={{ zIndex: theme.zIndex.drawer + 1 }}
        open={processingOrder || savingOrder}
      >
        <CircularProgress color="inherit" />
      </Backdrop>
    </Container>
  );
};

export default Products;