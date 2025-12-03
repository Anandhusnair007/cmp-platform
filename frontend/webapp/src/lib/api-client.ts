import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8082/api/v1';

// Create axios instance
export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add auth token interceptor
apiClient.interceptors.request.use(
  (config) => {
    const token = getAuthToken();
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);

// Handle 401 responses (unauthorized)
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Clear token and redirect to login
      clearAuthToken();
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);

// Token management (stored in memory only, not localStorage)
let authToken: string | null = null;

export const setAuthToken = (token: string) => {
  authToken = token;
};

export const getAuthToken = (): string | null => {
  return authToken;
};

export const clearAuthToken = () => {
  authToken = null;
};

// API methods (will be replaced by generated SDK)
export const authAPI = {
  login: async (email: string, password: string) => {
    const response = await apiClient.post('/auth/login', { email, password });
    return response.data;
  },
  getCurrentUser: async () => {
    const response = await apiClient.get('/auth/me');
    return response.data;
  },
  logout: async () => {
    const response = await apiClient.post('/auth/logout');
    return response.data;
  },
};

export const certsAPI = {
  list: async (params?: any) => {
    const response = await apiClient.get('/certs', { params });
    return response.data;
  },
  get: async (id: string) => {
    const response = await apiClient.get(`/certs/${id}`);
    return response.data;
  },
  request: async (data: any) => {
    const response = await apiClient.post('/certs/request', data);
    return response.data;
  },
  revoke: async (id: string, reason?: string) => {
    const response = await apiClient.post(`/certs/${id}/revoke`, { reason });
    return response.data;
  },
};

export const inventoryAPI = {
  get: async (params?: any) => {
    const response = await apiClient.get('/inventory', { params });
    return response.data;
  },
  getExpiring: async (days: number = 30) => {
    const response = await apiClient.get('/inventory/expiring', { params: { days } });
    return response.data;
  },
};

export const agentsAPI = {
  list: async () => {
    const response = await apiClient.get('/agents');
    return response.data;
  },
  install: async (agentId: string, certId: string, path: string, reloadCmd?: string) => {
    const response = await apiClient.post(`/agents/${agentId}/install`, {
      cert_id: certId,
      path,
      reload_cmd: reloadCmd,
    });
    return response.data;
  },
};
