import React, { useState } from 'react';
import { useLocation, useNavigate } from 'react-router-dom';
import axios from 'axios';
import { useAuth } from '../context/AuthContext';

function EditOrder() {
  const { token } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();
  const { order } = location.state;

  const [items, setItems] = useState(order.orders);

  const handleChangeQuantity = (index, value) => {
    const updated = [...items];
    updated[index].quantity = parseInt(value) || 0;
    updated[index].subtotal = updated[index].quantity * (updated[index].subtotal / order.orders[index].quantity);
    setItems(updated);
  };

  const handleSave = async () => {
    try {
      await axios.patch(`http://127.0.0.1:3000/pending-orders/${order.id}`, {
        orders: items,
      }, {
        headers: { Authorization: `Bearer ${token}` }
      });
      alert('Pesanan berhasil diperbarui!');
      navigate('/pending-orders');
    } catch (err) {
      console.error('Gagal update:', err.response?.data || err.message);
    }
  };

  return (
    <div>
      <h2>Edit Pesanan</h2>
      <ul style={{ listStyle: 'none', padding: 0 }}>
        {items.map((item, idx) => (
          <li key={idx} style={{ marginBottom: '10px' }}>
            <p>{item.name}</p>
            <input
              type="number"
              value={item.quantity}
              min={1}
              onChange={(e) => handleChangeQuantity(idx, e.target.value)}
            />
            <p>Subtotal: Rp{item.subtotal.toLocaleString()}</p>
          </li>
        ))}
      </ul>
      <button onClick={handleSave}>Simpan Perubahan</button>
    </div>
  );
}

export default EditOrder;
