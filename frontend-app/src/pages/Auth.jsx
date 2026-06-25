import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { userApi } from '../api';

export default function Auth() {
  const navigate = useNavigate();
  const [isLogin, setIsLogin] = useState(true);
  const [formData, setFormData] = useState({ 
    name: '',
    email: '', 
    password: '', 
    phone: '',
    role: 'CUSTOMER' 
  });
  const [error, setError] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    try {
      if (isLogin) {
        const res = await userApi.post('/auth/login', { 
          email: formData.email, 
          password: formData.password 
        });
        if (res.data.success) {
          localStorage.setItem('token', res.data.data.token);
          if (formData.role === 'MERCHANT') {
             navigate('/dashboard');
          } else {
             navigate('/');
          }
        }
      } else {
        const res = await userApi.post('/auth/register', formData);
        if (res.data.success) {
          setIsLogin(true);
          setError('Registration successful. Please login.');
        }
      }
    } catch (err) {
      if (!err.response) {
        setError('Network Error: Backend API (User Service) is not running.');
      } else {
        setError(err.response?.data?.error || 'Authentication failed');
      }
    }
  };

  return (
    <div className="container" style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh', padding: '3rem 1rem' }}>
      <div className="glass-panel animate-fade-in" style={{ width: '100%', maxWidth: '400px', margin: 'auto' }}>
        <div className="text-center mb-4">
          <h1 className="navbar-brand" style={{ fontSize: '2.2rem', marginBottom: '0.5rem', display: 'inline-block' }}>FoodDelivery</h1>
          <p style={{ color: 'var(--color-text-muted)', fontSize: '1rem', marginTop: '0' }}>{isLogin ? 'Welcome back, please sign in' : 'Create your new account'}</p>
        </div>
        {error && <div style={{ color: error.includes('successful') ? 'var(--color-success)' : 'var(--color-danger)', marginBottom: '1rem', textAlign: 'center' }}>{error}</div>}
        
        <form onSubmit={handleSubmit}>
          {!isLogin && (
            <div className="form-group">
              <label className="form-label">Name</label>
              <input 
                type="text" 
                className="form-control" 
                value={formData.name}
                onChange={e => setFormData({...formData, name: e.target.value})}
                required={!isLogin} 
              />
            </div>
          )}
          
          <div className="form-group">
            <label className="form-label">Email</label>
            <input 
              type="email" 
              className="form-control" 
              value={formData.email}
              onChange={e => setFormData({...formData, email: e.target.value})}
              required 
            />
          </div>

          {!isLogin && (
            <div className="form-group">
              <label className="form-label">Phone</label>
              <input 
                type="text" 
                className="form-control" 
                value={formData.phone}
                onChange={e => setFormData({...formData, phone: e.target.value})}
                required={!isLogin} 
              />
            </div>
          )}

          <div className="form-group">
            <label className="form-label">Password</label>
            <input 
              type="password" 
              className="form-control" 
              value={formData.password}
              onChange={e => setFormData({...formData, password: e.target.value})}
              required 
            />
          </div>

          <div className="form-group">
            <label className="form-label">Role</label>
            <select 
              className="form-control" 
              value={formData.role}
              onChange={e => setFormData({...formData, role: e.target.value})}
            >
              <option value="CUSTOMER">Customer</option>
              <option value="MERCHANT">Merchant (Restaurant)</option>
            </select>
          </div>

          <button type="submit" className="btn btn-primary btn-block">
            {isLogin ? 'Sign In' : 'Sign Up'}
          </button>
        </form>
        
        <p className="text-center mt-4 text-muted">
          {isLogin ? "Don't have an account? " : "Already have an account? "}
          <a href="#" onClick={(e) => { e.preventDefault(); setIsLogin(!isLogin); setError(''); }}>
            {isLogin ? 'Register' : 'Login'}
          </a>
        </p>
      </div>
    </div>
  );
}
