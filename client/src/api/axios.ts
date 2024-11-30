import axios from 'axios';

const API_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:5000/api';

export const api = axios.create({
    baseURL: API_URL,
    headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
    },
    withCredentials: true // Important for CORS with credentials
});

// Request interceptor to add auth token
api.interceptors.request.use((config) => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
}, (error) => {
    return Promise.reject(error);
});

// Response interceptor for error handling
api.interceptors.response.use(
    (response) => response,
    (error) => {
        if (error.response?.status === 401) {
            // Handle unauthorized access
            localStorage.removeItem('token');
            window.location.href = '/login';
        }
        return Promise.reject(error);
    }
);

// API functions
export const todoApi = {
    getTodos: () => api.get('/todos'),
    createTodo: (body: string) => api.post('/todos', { body }),
    updateTodo: (id: number) => api.patch(`/todos/${id}`),
    deleteTodo: (id: number) => api.delete(`/todos/${id}`),

};

export const userApi = {
    fetchProfile: () => api.get('/users/profile'),
}
