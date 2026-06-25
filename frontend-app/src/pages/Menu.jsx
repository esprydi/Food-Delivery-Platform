import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { catalogApi, orderApi } from '../api';

export default function Menu() {
  const { restaurantId } = useParams();
  const navigate = useNavigate();
  const [menus, setMenus] = useState([]);
  const [cart, setCart] = useState([]);
  const [loading, setLoading] = useState(true);
  const [checkoutLoading, setCheckoutLoading] = useState(false);

  useEffect(() => {
    const fetchMenus = async () => {
      try {
        const res = await catalogApi.get(`/restaurants/${restaurantId}/menus`);
        if (res.data.success) {
          setMenus(res.data.data || []);
        }
      } catch (err) {
        console.error("Failed to fetch menus", err);
      } finally {
        setLoading(false);
      }
    };
    fetchMenus();
  }, [restaurantId]);

  const addToCart = (menu) => {
    setCart(prev => {
      const existing = prev.find(item => item.menu_id === menu.id);
      if (existing) {
        return prev.map(item => item.menu_id === menu.id ? { ...item, quantity: item.quantity + 1 } : item);
      }
      return [...prev, { menu_id: menu.id, name: menu.name, quantity: 1, price: menu.price }];
    });
  };

  const totalAmount = cart.reduce((sum, item) => sum + (item.price * item.quantity), 0);

  const handleCheckout = async () => {
    if (cart.length === 0) return;
    setCheckoutLoading(true);
    try {
      const payload = {
        restaurant_id: restaurantId,
        delivery_address: "123 React Street",
        items: cart.map(c => ({ 
          menu_item_id: c.menu_id, 
          menu_item_name: c.name,
          quantity: c.quantity,
          unit_price: c.price
        }))
      };
      const res = await orderApi.post('/orders', payload);
      if (res.data.success) {
        alert("Order created! Proceeding to Payment...");
        navigate('/orders');
      }
    } catch (err) {
      alert("Checkout failed: " + (err.response?.data?.error || err.message));
    } finally {
      setCheckoutLoading(false);
    }
  };

  return (
    <div className="container animate-fade-in">
      <div className="d-flex justify-between align-center mb-4">
        <h2 style={{ textTransform: 'capitalize' }}>{restaurantId.replace('_', ' ')} Menu</h2>
        <button className="btn" onClick={() => navigate('/')}>Back</button>
      </div>

      <div className="grid" style={{ gridTemplateColumns: '2fr 1fr' }}>
        {/* Menu List */}
        <div className="grid grid-cols-2">
          {loading ? <p>Loading menu...</p> : menus.map((menu) => (
            <div key={menu.id} className="glass-panel d-flex flex-col justify-between">
              <div>
                <h3>{menu.name}</h3>
                <p className="text-muted mb-2">{menu.description}</p>
                <h4 style={{ color: 'var(--color-success)' }}>Rp {menu.price.toLocaleString()}</h4>
              </div>
              <button className="btn btn-primary mt-2" onClick={() => addToCart(menu)}>
                Add to Cart
              </button>
            </div>
          ))}
        </div>

        {/* Cart */}
        <div className="glass-panel" style={{ height: 'fit-content' }}>
          <h3>Your Cart</h3>
          {cart.length === 0 ? (
            <p className="text-muted mt-2">Cart is empty</p>
          ) : (
            <div className="mt-4">
              {cart.map((item, idx) => (
                <div key={idx} className="d-flex justify-between mb-2 pb-2" style={{ borderBottom: '1px solid rgba(255,255,255,0.1)' }}>
                  <div>
                    <span>{item.quantity}x {item.name}</span>
                  </div>
                  <span>Rp {(item.price * item.quantity).toLocaleString()}</span>
                </div>
              ))}
              <div className="d-flex justify-between mt-4 mb-4">
                <strong>Total</strong>
                <strong style={{ color: 'var(--color-success)' }}>Rp {totalAmount.toLocaleString()}</strong>
              </div>
              <button 
                className="btn btn-primary btn-block" 
                onClick={handleCheckout}
                disabled={checkoutLoading}
              >
                {checkoutLoading ? 'Processing...' : 'Checkout & Pay'}
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
