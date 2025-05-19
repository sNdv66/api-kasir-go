import {
  AppBar,
  Box,
  Drawer,
  IconButton,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Toolbar,
  Typography,
  Badge,
  Avatar,
  Menu,
  MenuItem,
  Divider,
  useMediaQuery,
  useTheme,
} from "@mui/material";
import {
  Dashboard,
  PointOfSale,
  History,
  Inventory,
  BarChart,
  Menu as MenuIcon,
  Notifications as NotificationsIcon,
} from "@mui/icons-material";
import { useNavigate, Outlet, useLocation } from "react-router-dom";
import { useState ,useEffect} from "react";
import { useAuth } from "../context/AuthContext";
import axios from "axios";



const drawerWidth = 240;

const menuItems = [
  { text: "Dashboard", icon: <Dashboard />, path: "/dashboard" },
  { text: "Transaksi", icon: <PointOfSale />, path: "/transaksi" },
  { text: "Riwayat", icon: <History />, path: "/pending-orders"},
  { text: "Produk", icon: <Inventory />, path: "/produk" },
  { text: "Laporan", icon: <BarChart />, path: "/laporan" },
  
];

const AppLayout = () => {
  const { user, token ,logout} = useAuth();
  const [branchName, setBranchName] = useState("Memuat...");

  useEffect(() => {
    if (user?.branch_id && token) {
      axios
        .get(`http://127.0.0.1:3000/branches/${user.branch_id}`, {
          headers: { Authorization: `Bearer ${token}` },
        })
        .then((res) => setBranchName(res.data.name))
        .catch(() => setBranchName("N/A"));
    }
  }, [user, token]);
  
  const navigate = useNavigate();
  const location = useLocation();
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down("md"));

  const [mobileOpen, setMobileOpen] = useState(false);
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [notifAnchorEl, setNotifAnchorEl] = useState<null | HTMLElement>(null);
  const [notifications] = useState([
    { id: 1, message: "Transaksi berhasil #TRX123", type: "success" },
    { id: 2, message: "Transaksi gagal #TRX124", type: "error" },
  ]);

  const handleDrawerToggle = () => {
    setMobileOpen(!mobileOpen);
  };

  const drawer = (
    <Box sx={{ width: drawerWidth }}>
      <Toolbar />
      <List>
        {menuItems.map((item) => (
          <ListItem disablePadding key={item.text}>
            <ListItemButton
              selected={location.pathname === item.path}
              onClick={() => {
                navigate(item.path);
                if (isMobile) setMobileOpen(false);
              }}
            >
              <ListItemIcon>{item.icon}</ListItemIcon>
              <ListItemText primary={item.text} />
            </ListItemButton>
          </ListItem>
        ))}
      </List>
    </Box>
  );

  return (
    <Box sx={{ display: "flex" }}>
      {/* AppBar/Header */}
      <AppBar position="fixed" sx={{ zIndex: theme.zIndex.drawer + 1 }}>
        <Toolbar sx={{ display: "flex", justifyContent: "space-between" }}>
          <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
            {isMobile && (
              <IconButton color="inherit" onClick={handleDrawerToggle}>
                <MenuIcon />
              </IconButton>
            )}
          </Box>

          <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
            <IconButton color="inherit" onClick={(e) => setNotifAnchorEl(e.currentTarget)}>
              <Badge badgeContent={notifications.length} color="error">
                <NotificationsIcon />
              </Badge>
            </IconButton>

            <Menu
              anchorEl={notifAnchorEl}
              open={Boolean(notifAnchorEl)}
              onClose={() => setNotifAnchorEl(null)}
              PaperProps={{ style: { width: 200 } }}
            >
              {notifications.length === 0 ? (
                <MenuItem disabled>Tidak ada notifikasi</MenuItem>
              ) : (
                notifications.map((notif) => (
                  <MenuItem key={notif.id}>
                    <Typography
                      variant="body2"
                      color={notif.type === "error" ? "error" : "success.main"}
                    >
                      {notif.message}
                    </Typography>
                  </MenuItem>
                ))
              )}
            </Menu>

            <IconButton onClick={(e) => setAnchorEl(e.currentTarget)}>
              <Avatar>{user?.name?.charAt(0).toUpperCase()}</Avatar>
            </IconButton>
          </Box>

          <Menu anchorEl={anchorEl} open={Boolean(anchorEl)} onClose={() => setAnchorEl(null)}>
            <MenuItem disabled>Cabang: {branchName}</MenuItem>
            <MenuItem disabled>Role: {user?.role}</MenuItem>
            <Divider />
            <MenuItem onClick={() => { logout(); navigate("/"); }}>Logout</MenuItem>
          </Menu>
        </Toolbar>
      </AppBar>

      {/* Drawer */}
      <Box component="nav" sx={{ width: { md: drawerWidth }, flexShrink: { md: 0 } }}>
        <Drawer
          variant={isMobile ? "temporary" : "permanent"}
          open={isMobile ? mobileOpen : true}
          onClose={handleDrawerToggle}
          ModalProps={{ keepMounted: true }}
          sx={{
            "& .MuiDrawer-paper": {
              width: drawerWidth,
              boxSizing: "border-box",
            },
          }}
        >
          {drawer}
        </Drawer>
      </Box>

      {/* Konten */}
      <Box component="main" sx={{ flexGrow: 1, p: 2, mt: 5 }}>
        <Outlet />
      </Box>
    </Box>
  );
};

export default AppLayout;