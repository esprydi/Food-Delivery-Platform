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
                  <div key={m.id} className="glass-panel d-flex flex-col justify-between">
                    <div>
                      <h4>{m.name}</h4>
                      <p className="text-muted">{m.description}</p>
                      <strong style={{ color: 'var(--color-success)' }}>Rp {m.price.toLocaleString()}</strong>
                    </div>
                    <button className="btn mt-2" onClick={() => handleEditClick(m)}>Edit</button>
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
