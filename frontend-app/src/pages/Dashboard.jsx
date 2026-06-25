import { useState, useEffect } from 'react';
import { catalogApi, orderApi } from '../api';
import { useNavigate } from 'react-router-dom';

export default function Dashboard() {
  const navigate = useNavigate();
  const [restaurant, setRestaurant] = useState(null);
  const [loading, setLoading] = useState(true);
  const [restaurantForm, setRestaurantForm] = useState({ name: '', address: '' });
  
  const [menuForm, setMenuForm] = useState({ name: '', description: '', price: '' });
  const [menus, setMenus] = useState([]);
  
  const [editingMenuId, setEditingMenuId] = useState(null);

  useEffect(() => {
    fetchRestaurant();
  }, []);

  const fetchRestaurant = async () => {
    try {
      const res = await catalogApi.get('/merchant/restaurants/me');
      if (res.data.success) {
        setRestaurant(res.data.data);
        fetchMenus(res.data.data.id);
      }
    } catch (err) {
      if (err.response?.status !== 404) {
        console.error("Failed to fetch restaurant", err);
      }
    } finally {
      setLoading(false);
    }
  };

  const fetchMenus = async (restId) => {
    try {
      const res = await catalogApi.get(`/restaurants/${restId}/menus`);
      if (res.data.success) {
        setMenus(res.data.data);
      }
    } catch (err) {
      console.error("Failed to fetch menus", err);
    }
  };

  const handleCreateRestaurant = async (e) => {
    e.preventDefault();
    try {
      const res = await catalogApi.post('/merchant/restaurants', restaurantForm);
      if (res.data.success) {
        setRestaurant(res.data.data);
      }
    } catch (err) {
      alert("Failed: " + (err.response?.data?.error || err.message));
    }
  };

  const handleSubmitMenu = async (e) => {
    e.preventDefault();
    try {
      const rawPrice = menuForm.price.toString().replace(/\D/g, '');
      const payload = {
        restaurant_id: restaurant.id,
        name: menuForm.name,
        description: menuForm.description,
        price: parseFloat(rawPrice) || 0
      };

      if (editingMenuId) {
        const res = await catalogApi.put(`/merchant/menus/${editingMenuId}`, payload);
        if (res.data.success) {
          setEditingMenuId(null);
          setMenuForm({ name: '', description: '', price: '' });
          fetchMenus(restaurant.id);
        }
      } else {
        const res = await catalogApi.post('/merchant/menus', payload);
        if (res.data.success) {
          setMenuForm({ name: '', description: '', price: '' });
          fetchMenus(restaurant.id);
        }
      }
    } catch (err) {
      alert("Failed: " + (err.response?.data?.error || err.message));
    }
  };

  const handleEditClick = (menu) => {
    setEditingMenuId(menu.id);
    setMenuForm({
      name: menu.name,
      description: menu.description,
      price: menu.price.toString()
    });
  };

  const handleDeleteMenu = async (menuId) => {
    if (window.confirm("Are you sure you want to delete this menu?")) {
      try {
        const res = await catalogApi.delete(`/merchant/menus/${menuId}`);
        if (res.data.success) {
          fetchMenus(restaurant.id);
        }
      } catch (err) {
        alert("Failed to delete menu: " + (err.response?.data?.error || err.message));
      }
    }
  };

  const cancelEdit = () => {
    setEditingMenuId(null);
    setMenuForm({ name: '', description: '', price: '' });
  };

  if (loading) return <div className="container"><p>Loading dashboard...</p></div>;

  return (
    <div className="container animate-fade-in">
      <div className="d-flex justify-between align-center mb-4">
        <h2>Merchant Dashboard</h2>
        <button className="btn" onClick={() => {
          localStorage.removeItem('token');
          navigate('/auth');
        }}>Logout</button>
      </div>

      {!restaurant ? (
        <div className="glass-panel" style={{ maxWidth: '500px', margin: '0 auto' }}>
          <h3>Create Your Restaurant</h3>
          <p className="text-muted mb-4">You need to set up your restaurant before you can add menus.</p>
          <form onSubmit={handleCreateRestaurant}>
            <div className="form-group">
              <label className="form-label">Restaurant Name</label>
              <input type="text" className="form-control" value={restaurantForm.name} onChange={e => setRestaurantForm({...restaurantForm, name: e.target.value})} required />
            </div>
            <div className="form-group">
              <label className="form-label">Address</label>
              <input type="text" className="form-control" value={restaurantForm.address} onChange={e => setRestaurantForm({...restaurantForm, address: e.target.value})} required />
            </div>
            <button type="submit" className="btn btn-primary btn-block">Create Restaurant</button>
          </form>
        </div>
      ) : (
        <div className="grid" style={{ gridTemplateColumns: '1fr 2fr' }}>
          <div className="glass-panel" style={{ height: 'fit-content' }}>
            <h3>Restaurant Details</h3>
            <p><strong>Name:</strong> {restaurant.name}</p>
            <p><strong>Address:</strong> {restaurant.address}</p>

            <h3 className="mt-4">{editingMenuId ? "Edit Menu Item" : "Add Menu Item"}</h3>
            <form onSubmit={handleSubmitMenu} className="mt-2">
              <div className="form-group">
                <label className="form-label">Menu Name</label>
                <input type="text" className="form-control" value={menuForm.name} onChange={e => setMenuForm({...menuForm, name: e.target.value})} required />
              </div>
              <div className="form-group">
                <label className="form-label">Description</label>
                <textarea className="form-control" value={menuForm.description} onChange={e => setMenuForm({...menuForm, description: e.target.value})} required />
              </div>
              <div className="form-group">
                <label className="form-label">Price (Rp)</label>
                <input type="text" className="form-control" value={menuForm.price} onChange={e => setMenuForm({...menuForm, price: e.target.value})} required />
              </div>
              <div className="d-flex" style={{ gap: '10px' }}>
                <button type="submit" className="btn btn-primary" style={{ flex: 1 }}>{editingMenuId ? "Save Changes" : "Add Menu"}</button>
                {editingMenuId && (
                  <button type="button" className="btn" style={{ flex: 1 }} onClick={cancelEdit}>Cancel</button>
                )}
              </div>
            </form>
          </div>

          <div className="glass-panel">
            <h3>Your Menu Items</h3>
            {menus.length === 0 ? <p className="text-muted mt-2">No menus added yet.</p> : (
              <div className="grid grid-cols-2 mt-4">
                {menus.map(m => (
                  <div key={m.id} className="glass-panel" style={{ position: 'relative', padding: '20px', display: 'flex', flexDirection: 'column', height: '100%' }}>
                    <div className="d-flex justify-between align-start" style={{ gap: '15px' }}>
                      <div style={{ flex: 1 }}>
                        <h4 style={{ margin: '0 0 8px 0', fontSize: '1.2rem' }}>{m.name}</h4>
                        <p className="text-muted" style={{ margin: '0 0 15px 0', fontSize: '0.9rem', lineHeight: '1.4' }}>{m.description}</p>
                      </div>
                      <div className="d-flex" style={{ gap: '8px' }}>
                        <button 
                          className="btn" 
                          style={{ padding: '0', borderRadius: '50%', display: 'flex', alignItems: 'center', justifyContent: 'center', width: '36px', height: '36px' }} 
                          onClick={() => handleEditClick(m)} 
                          title="Edit"
                        >
                          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M12 20h9"></path><path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4L16.5 3.5z"></path></svg>
                        </button>
                        <button 
                          className="btn" 
                          style={{ padding: '0', borderRadius: '50%', display: 'flex', alignItems: 'center', justifyContent: 'center', width: '36px', height: '36px', backgroundColor: 'rgba(220, 53, 69, 0.9)', color: 'white', border: 'none' }} 
                          onClick={() => handleDeleteMenu(m.id)} 
                          title="Delete"
                        >
                          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="3 6 5 6 21 6"></polyline><path d="M19 6v14a2 2 0 0 1-2 2H7a2 2 0 0 1-2-2V6m3 0V4a2 2 0 0 1 2-2h4a2 2 0 0 1 2 2v2"></path></svg>
                        </button>
                      </div>
                    </div>
                    <div style={{ marginTop: 'auto' }}>
                      <strong style={{ color: 'var(--color-success)', fontSize: '1.1rem' }}>Rp {m.price.toLocaleString()}</strong>
                    </div>
                  </div>
                ))}
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
