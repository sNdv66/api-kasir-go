import {
  AppBar,
  Toolbar,
  Typography,
  IconButton,
  Avatar,
  Menu,
  MenuItem,
  Divider,
  Box,
  Badge,
} from "@mui/material";
import {
  Notifications as NotificationsIcon,
} from "@mui/icons-material";
import { useAuth } from "../context/AuthContext";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

const Header = ({ branchName }: { branchName: string }) => {
  const { user, logout } = useAuth();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [notifAnchorEl, setNotifAnchorEl] = useState<null | HTMLElement>(null);
  const [notifications, setNotifications] = useState([
    { id: 1, message: "Transaksi berhasil #TRX123", type: "success" },
    { id: 2, message: "Transaksi gagal #TRX124", type: "error" },
  ]);

  const open = Boolean(anchorEl);
  const openNotif = Boolean(notifAnchorEl);
  const navigate = useNavigate();

  const handleClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleNotifClick = (event: React.MouseEvent<HTMLElement>) => {
    setNotifAnchorEl(event.currentTarget);
  };

  const handleNotifClose = () => {
    setNotifAnchorEl(null);
  };

  const handleLogout = () => {
    logout();
    navigate("/");
  };

  return (
    <AppBar position="static">
      <Toolbar sx={{ display: "flex", justifyContent: "space-between" }}>
        <Typography variant="h6">Dashboard</Typography>

        <Box sx={{ display: "flex", alignItems: "center", gap: 1 }}>
          {/* Notifikasi */}
          <IconButton color="inherit" onClick={handleNotifClick}>
            <Badge badgeContent={notifications.length} color="error">
              <NotificationsIcon />
            </Badge>
          </IconButton>

          <Menu
            anchorEl={notifAnchorEl}
            open={openNotif}
            onClose={handleNotifClose}
            PaperProps={{ style: { width: 300 } }}
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

          {/* Avatar dan Menu Akun */}
          <IconButton onClick={handleClick}>
            <Avatar>{user?.name?.charAt(0).toUpperCase()}</Avatar>
          </IconButton>
        </Box>

        <Menu anchorEl={anchorEl} open={open} onClose={handleClose}>
          <MenuItem disabled>Cabang: {branchName}</MenuItem>
          <MenuItem disabled>Role: {user?.role}</MenuItem>
          <Divider />
          <MenuItem onClick={handleLogout}>Logout</MenuItem>
        </Menu>
      </Toolbar>
    </AppBar>
  );
};

export default Header;