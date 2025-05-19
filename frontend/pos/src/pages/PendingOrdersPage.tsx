import React, { useEffect, useState } from 'react';
import axios from 'axios';
import { useAuth } from '../context/AuthContext';

function PendingOrders() {
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const { token } = useAuth();

  useEffect(() => {
    axios.get('http://127.0.0.1:3000/pending-orders', {
      headers: {
        Authorization: `Bearer ${token}`
      }
    })
    .then(res => {
      setOrders(res.data || []);
      console.log(res.data);
      setLoading(false);
    })
    .catch(err => {
      console.error('Gagal mengambil data:', err.response?.data || err.message);
      setError(err.message);
      setLoading(false);
    });
  }, [token]);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!orders.length) return <div>No pending orders found</div>;

  return (
    <div>
      <h2>Pending Orders</h2>
      {orders.map((order) => (
        <div key={order.id} style={{ border: '1px solid #ccc', padding: '10px', marginBottom: '20px' }}>
          <p><strong>Order ID:</strong> {order.id}</p>
          <p><strong>Tanggal:</strong> {order.created_at ? new Date(order.created_at).toLocaleString() : 'N/A'}</p>
          <h4>Produk:</h4>
          <p><strong> brach_id : </strong>
              {order.branch_id} </p>
          <ul>
            {order.orders?.map((item, idx) => (
              <li key={idx}>
                <p>Nama: {item.name || 'N/A'}</p>
                <p>Jumlah: {item.quantity || 'N/A'}</p>
                <p>Harga: Rp{item.subtotal ? item.subtotal.toLocaleString() : 'N/A'}</p>
              </li>
            ))}
          </ul>
          <p><strong>Total Harga:</strong> Rp{order.total ? order.total.toLocaleString() : 'N/A'}</p>
        </div>
      ))}
    </div>
  );
}

export default PendingOrders;