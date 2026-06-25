import axios from 'axios';

// Get token from local storage
const getToken = () => localStorage.getItem('token');

// Utility to create an axios instance for a specific microservice port
const createApiClient = (port) => {
  const instance = axios.create({
    baseURL: `http://localhost:${port}/api/v1`,
  });

  instance.interceptors.request.use((config) => {
    const token = getToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  });

  return instance;
};

export const userApi = createApiClient(8081);
export const catalogApi = createApiClient(8082);
export const orderApi = createApiClient(8083);
export const paymentApi = createApiClient(8084);
