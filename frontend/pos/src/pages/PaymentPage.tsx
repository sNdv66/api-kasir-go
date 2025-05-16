
import {
  Box,
  Typography,
  Card,
  CardContent,
  Grid,
  TextField,
  MenuItem,
  Button,
  Divider
} from "@mui/material";
import { useEffect, useState } from "react";
import axios from "axios";

interface OrderItem {
  id: string;
  name: string;
  price: number;
  quantity: number;
  stock: number;
  image_url?: string;
  category?: string;
}

const PaymentPage = () => {
  const [orders, setOrders] = useState<OrderItem[]>([]);
  const [total, setTotal] = useState(0);
  const [paymentMethod, setPaymentMethod] = useState("cash");
  const [paidAmount, setPaidAmount] = useState<number>(0);
  const [change, setChange] = useState(0);

  useEffect(() => {
    const storedOrders = sessionStorage.getItem("orders");
    if (storedOrders) {
      const parsedOrders = JSON.parse(storedOrders);
      setOrders(parsedOrders);
      const totalHarga = parsedOrders.reduce(
        (sum: number, item: OrderItem) => sum + item.price * item.quantity,
        0
      );
      setTotal(totalHarga);
    }
  }, []);

  useEffect(() => {
    setChange(paidAmount - total);
  }, [paidAmount, total]);
  
  const handleSubmit = async () => {
    try {
      const payload = {
        items: orders.map((item) => ({
          product_id: item.id,
          quantity: item.quantity,
        })),
        payment_method: paymentMethod,
        paid_amount: paidAmount,
      };
      const token = localStorage.getItem("token"); // atau sesuaikan sumber token

      const response = await axios.post("http://localhost:3000/payment", payload, {
      headers: {
       Authorization: `Bearer ${token}`,
      },
      });

      alert("Pembayaran berhasil!");
      sessionStorage.removeItem("orders");
      window.location.href = "/"; // redirect ke halaman awal

    } catch (error: any) {
      console.error("Gagal bayar:", error);
      alert(
  error.response?.data?.message ||
  JSON.stringify(error.response?.data) ||
  "Terjadi kesalahan."
);
    }
  };
  
  

  return (
    <Box p={2}>
      <Typography variant="h5" gutterBottom>
        Pembayaran
      </Typography>

      <Grid container spacing={2}>
        <Grid item xs={12} md={6}>
          {orders.map((item) => (
            <Card key={item.id} sx={{ mb: 2 }}>
              <CardContent>
                <Typography variant="subtitle1">{item.name}</Typography>
                <Typography>
                  {item.quantity} x Rp{item.price.toLocaleString()}
                </Typography>
                <Typography variant="body2">
                  Subtotal: Rp{(item.price * item.quantity).toLocaleString()}
                </Typography>
              </CardContent>
            </Card>
          ))}
          <Divider />
          <Typography variant="h6" mt={2}>
            Total: Rp{total.toLocaleString()}
          </Typography>
        </Grid>

        <Grid item xs={12} md={6}>
          <Typography variant="h6">Form Pembayaran</Typography>
          <TextField
            fullWidth
            select
            label="Metode Pembayaran"
            value={paymentMethod}
            onChange={(e) => setPaymentMethod(e.target.value)}
            sx={{ mt: 2 }}
          >
            <MenuItem value="cash">Cash</MenuItem>
            <MenuItem value="qris">QRIS</MenuItem>
            <MenuItem value="debit">Debit</MenuItem>
          </TextField>

          <TextField
            fullWidth
            label="Jumlah Bayar"
            type="number"
            value={paidAmount}
            onChange={(e) => setPaidAmount(Number(e.target.value))}
            sx={{ mt: 2 }}
          />

          <Typography sx={{ mt: 2 }}>
            Kembalian: Rp{change > 0 ? change.toLocaleString() : 0}
          </Typography>

                <Button
        variant="contained"
        color="primary"
        fullWidth
        sx={{ mt: 3 }}
        disabled={paidAmount < total}
        onClick={handleSubmit}
      >
        Bayar Sekarang
      </Button>
      
        </Grid>
      </Grid>
    </Box>
  );
};

export default PaymentPage;