import { useState, useEffect } from 'react';
import { orderApi, paymentApi } from '../api';
import { useNavigate } from 'react-router-dom';

export default function Orders() {
  const navigate = useNavigate();
  const [orders, setOrders] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchOrders();
  }, []);

  const fetchOrders = async () => {
    try {
      const res = await orderApi.get('/orders/customer');
      if (res.data.success) {
        setOrders(res.data.data || []);
      }
    } catch (err) {
      console.error("Failed to fetch orders", err);
    } finally {
      setLoading(false);
    }
  };

  const handlePay = async (orderId) => {
    try {
      const res = await paymentApi.get(`/payments/order/${orderId}`);
      if (res.data.success && res.data.data?.snap_url) {
        window.location.href = res.data.data.snap_url;
      } else {
        alert("Payment URL not found or still generating.");
      }
    } catch (err) {
      alert("Failed to fetch payment info: " + (err.response?.data?.error || err.message));
    }
  };

  if (loading) return <div className="container"><p>Loading your orders...</p></div>;

  return (
    <div className="container animate-fade-in">
      <div className="d-flex justify-between align-center mb-4">
        <h2>Your Orders</h2>
        <button className="btn" onClick={() => navigate('/')}>Back to Home</button>
      </div>

      {orders.length === 0 ? (
        <div className="glass-panel text-center">
          <p className="text-muted">You haven't placed any orders yet.</p>
          <button className="btn btn-primary mt-4" onClick={() => navigate('/')}>Find Restaurants</button>
        </div>
      ) : (
        <div className="grid">
          {orders.map((order) => (
            <div key={order.id} className="glass-panel">
              <div className="d-flex justify-between align-center mb-2">
                <span className="text-muted" style={{ fontSize: '0.9rem' }}>Order #{order.id.split('-')[0]}</span>
                <span style={{ 
                  padding: '4px 8px', 
                  borderRadius: '12px', 
                  fontSize: '0.8rem',
                  backgroundColor: order.status === 'PAID' ? 'rgba(46, 213, 115, 0.2)' : 'rgba(255, 165, 2, 0.2)',
                  color: order.status === 'PAID' ? 'var(--color-success)' : 'var(--color-warning)'
                }}>
                  {order.status}
                </span>
              </div>
              <p><strong>Address:</strong> {order.delivery_address}</p>
              <h4 className="mt-2" style={{ color: 'var(--color-success)' }}>Rp {order.total_amount.toLocaleString()}</h4>
              
              {order.status === 'PENDING' && (
                <button 
                  className="btn btn-primary mt-4" 
                  style={{ width: '100%' }}
                  onClick={() => handlePay(order.id)}
                >
                  Pay Now
                </button>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
