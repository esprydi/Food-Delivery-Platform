import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import Auth from './pages/Auth';
import Home from './pages/Home';
import Menu from './pages/Menu';
import Dashboard from './pages/Dashboard';
import Orders from './pages/Orders';
import './index.css';

// Simple protected route wrapper
const ProtectedRoute = ({ children }) => {
  const token = localStorage.getItem('token');
  if (!token) {
    return <Navigate to="/auth" replace />;
  }
  return children;
};

function App() {
  return (
    <Router>
      <div className="navbar container">
        <div className="navbar-brand">FoodDelivery</div>
        <div className="navbar-nav">
          <a href="#" onClick={(e) => {
            e.preventDefault();
            localStorage.removeItem('token');
            window.location.href = '/auth';
          }}>Logout</a>
        </div>
      </div>
      
      <Routes>
        <Route path="/auth" element={<Auth />} />
        
        <Route path="/" element={
          <ProtectedRoute>
            <Home />
          </ProtectedRoute>
        } />
        
        <Route path="/restaurant/:restaurantId" element={
          <ProtectedRoute>
            <Menu />
          </ProtectedRoute>
        } />

        <Route path="/dashboard" element={
          <ProtectedRoute>
            <Dashboard />
          </ProtectedRoute>
        } />

        <Route path="/orders" element={
          <ProtectedRoute>
            <Orders />
          </ProtectedRoute>
        } />
      </Routes>
    </Router>
  );
}

export default App;
