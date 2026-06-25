import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { catalogApi } from '../api';

export default function Home() {
  const navigate = useNavigate();
  const [restaurants, setRestaurants] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchRestaurants = async () => {
      try {
        const res = await catalogApi.get('/restaurants');
        if (res.data.success) {
          setRestaurants(res.data.data || []);
        }
      } catch (err) {
        console.error("Failed to fetch", err);
      } finally {
        setLoading(false);
      }
    };
    fetchRestaurants();
  }, []);

  return (
    <div className="container animate-fade-in">
      <div className="d-flex justify-between align-center mb-4">
        <h2>Restaurants near you</h2>
      </div>

      {loading ? (
        <p>Loading...</p>
      ) : (
        <div className="grid grid-cols-3">
          {restaurants.map((resto, index) => (
            <div key={index} className="glass-panel">
              <h3 style={{ textTransform: 'capitalize' }}>{resto.name}</h3>
              <p className="text-muted mb-4">{resto.address}</p>
              <button 
                className="btn btn-primary btn-block"
                onClick={() => navigate(`/restaurant/${resto.id}`)}
              >
                View Menu
              </button>
            </div>
          ))}
          {restaurants.length === 0 && (
             <div className="glass-panel" style={{ gridColumn: 'span 3', textAlign: 'center' }}>
               <p>No restaurants found. Please add menus as a merchant first.</p>
             </div>
          )}
        </div>
      )}
    </div>
  );
}
