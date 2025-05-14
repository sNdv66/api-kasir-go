import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Login from './pages/Login';
import Dashboard from './pages/Dashboard';
import Transaksi from './pages/Transaksi';
import Riwayat from './pages/Riwayat';
import Produk from './pages/Produk';
import Laporan from './pages/Laporan';
import PrivateRoute from './routes/PrivateRoute';
import AppLayout from './components/AppLayout';
import { AuthProvider } from './context/AuthContext';

function App() {
  return (
    <AuthProvider>
      <Router>
        <Routes>
          {/* Login page (tanpa layout) */}
          <Route path="/" element={<Login />} />

          {/* Semua halaman lain dibungkus AppLayout + PrivateRoute */}
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
          </Route>
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;