import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import Transaksi from './pages/Transaksi';
import Riwayat from './pages/Riwayat';
import Produk from './pages/Produk';
import Laporan from './pages/Laporan';
import PaymentPage from './pages/PaymentPage';
import PrivateRoute from './routes/PrivateRoute';
import AppLayout from './components/AppLayout';
import { AuthProvider } from './context/AuthContext';


function App() {
  return (
     <AuthProvider>
    <Router>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route
          path="/"
          element={
            <PrivateRoute>
              <AppLayout />
            </PrivateRoute>
          }
        >
          <Route path="dashboard" element={<Dashboard />} />
          <Route path="transaksi" element={<Transaksi />} />
          <Route path="riwayat" element={<Riwayat />} />
          <Route path="produk" element={<Produk />} />
          <Route path="laporan" element={<Laporan />} />
          <Route path="payment" element={<PaymentPage />} />
        </Route>
      </Routes>
    </Router>
</AuthProvider>
  );
}

export default App;