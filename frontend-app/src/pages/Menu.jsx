import { useState, useEffect } from 'react';
import { useParams, useNavigate, useLocation } from 'react-router-dom';
import { catalogApi, orderApi } from '../api';

export default function Menu() {
  const { restaurantId } = useParams();
  const location = useLocation();
  const navigate = useNavigate();
  const restaurantName = location.state?.restaurantName || restaurantId;
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

  const updateQuantity = (menuId, delta) => {
    setCart(prev => prev.map(item => {
      if (item.menu_id === menuId) {
        return { ...item, quantity: item.quantity + delta };
      }
      return item;
    }).filter(item => item.quantity > 0));
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
        <h2 style={{ textTransform: 'capitalize' }}>{restaurantName} Menu</h2>
        <button className="btn" onClick={() => navigate('/')}>Back</button>
      </div>

      <div className="grid" style={{ gridTemplateColumns: '2fr 1fr', alignItems: 'start' }}>
        {/* Menu List */}
        <div className="grid grid-cols-2">
          {loading ? <p>Loading menu...</p> : menus.map((menu) => (
            <div key={menu.id} className="glass-panel" style={{ position: 'relative', padding: '20px', display: 'flex', flexDirection: 'column', height: '100%' }}>
              <div className="d-flex justify-between align-start" style={{ gap: '15px' }}>
                <div style={{ flex: 1 }}>
                  <h3 style={{ margin: '0 0 8px 0', fontSize: '1.2rem' }}>{menu.name}</h3>
                  <p className="text-muted" style={{ margin: '0 0 15px 0', fontSize: '0.9rem', lineHeight: '1.4' }}>{menu.description}</p>
                </div>
                <strong style={{ color: 'var(--color-success)', fontSize: '1.2rem', whiteSpace: 'nowrap' }}>Rp {menu.price.toLocaleString()}</strong>
              </div>
              <button 
                className="btn btn-primary" 
                style={{ width: '100%', marginTop: 'auto', padding: '8px 12px', fontSize: '0.95rem', display: 'flex', alignItems: 'center', justifyContent: 'center', gap: '6px' }} 
                onClick={() => addToCart(menu)}
              >
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="9" cy="21" r="1"></circle><circle cx="20" cy="21" r="1"></circle><path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"></path></svg>
                Add to Cart
              </button>
            </div>
          ))}
        </div>

        {/* Cart */}
        <div className="cart-panel" style={{ height: 'fit-content' }}>
          <h3 style={{ borderBottom: '1px solid rgba(255,255,255,0.1)', paddingBottom: '15px', marginBottom: '20px' }}>Your Cart</h3>
          {cart.length === 0 ? (
            <p className="text-muted mt-2">Cart is empty</p>
          ) : (
            <div className="mt-4">
              {cart.map((item, idx) => (
                <div key={idx} className="d-flex justify-between align-center mb-4 pb-3" style={{ borderBottom: '1px solid rgba(255,255,255,0.05)' }}>
                  <div className="d-flex flex-col" style={{ gap: '8px', marginRight: '15px' }}>
                    <div>
                      <span style={{ fontWeight: '600', fontSize: '1.1rem', display: 'block', marginBottom: '2px' }}>{item.name}</span>
                      <span style={{ color: 'var(--color-text-muted)', fontSize: '0.85rem' }}>Rp {item.price.toLocaleString()}</span>
                    </div>
                    {/* Horizontal Quantity Stepper */}
                    <div className="d-flex align-center" style={{ background: 'rgba(15, 23, 42, 0.8)', borderRadius: '12px', border: '1px solid rgba(255,255,255,0.1)', width: 'fit-content' }}>
                      <button style={{ background: 'transparent', border: 'none', color: 'var(--color-text-muted)', padding: '4px 8px', cursor: 'pointer', display: 'flex', alignItems: 'center' }} onClick={() => updateQuantity(item.menu_id, -1)}>
                        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><line x1="5" y1="12" x2="19" y2="12"></line></svg>
                      </button>
                      <span style={{ fontSize: '0.85rem', fontWeight: 'bold', minWidth: '20px', textAlign: 'center', color: 'white' }}>{item.quantity}</span>
                      <button style={{ background: 'transparent', border: 'none', color: 'var(--color-primary)', padding: '4px 8px', cursor: 'pointer', display: 'flex', alignItems: 'center' }} onClick={() => updateQuantity(item.menu_id, 1)}>
                        <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line></svg>
                      </button>
                    </div>
                  </div>
                  {/* Total Price */}
                  <strong style={{ color: 'var(--color-success)', fontSize: '1.15rem' }}>Rp {(item.price * item.quantity).toLocaleString()}</strong>
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
