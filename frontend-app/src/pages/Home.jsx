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
            <div key={index} className="glass-panel" style={{ display: 'flex', flexDirection: 'column', gap: '15px', padding: '24px' }}>
              <div>
                <h3 style={{ textTransform: 'capitalize', fontSize: '1.4rem', marginBottom: '8px', display: 'flex', alignItems: 'center', gap: '8px' }}>
                  <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="var(--color-primary)" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M3 2v7c0 1.1.9 2 2 2h4a2 2 0 0 0 2-2V2"></path><path d="M7 2v20"></path><path d="M21 15V2v0a5 5 0 0 0-5 5v6c0 1.1.9 2 2 2h3Zm0 0v7"></path></svg>
                  {resto.name}
                </h3>
                <p className="text-muted" style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path><circle cx="12" cy="10" r="3"></circle></svg>
                  {resto.address}
                </p>
              </div>
              <button 
                className="btn btn-primary btn-block"
                style={{ marginTop: 'auto', padding: '12px', fontSize: '1.05rem', fontWeight: '500' }}
                onClick={() => navigate(`/restaurant/${resto.id}`, { state: { restaurantName: resto.name } })}
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
