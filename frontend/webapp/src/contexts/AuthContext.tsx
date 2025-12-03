import React, { createContext, useContext, useState, useEffect } from 'react';

interface User {
  id: string;
  email: string;
  name: string;
  roles: string[];
  team?: string;
}

interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

// Mock user data for demo
const MOCK_USER: User = {
  id: '1',
  email: 'admin@example.com',
  name: 'Admin User',
  roles: ['admin', 'security_engineer'],
  team: 'Security Operations',
};

// Valid demo credentials
const DEMO_CREDENTIALS = [
  { email: 'admin@example.com', password: 'admin' },
  { email: 'admin@example.com', password: '9895' },
];

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [user, setUser] = useState<User | null>(() => {
    // Initialize from localStorage
    const stored = localStorage.getItem('cmp_user');
    return stored ? JSON.parse(stored) : null;
  });
  const [loading, setLoading] = useState(false);

  const login = async (email: string, password: string) => {
    // Simulate API call delay
    await new Promise(resolve => setTimeout(resolve, 500));

    // Check credentials
    const isValid = DEMO_CREDENTIALS.some(
      cred => cred.email.toLowerCase() === email.toLowerCase() && cred.password === password
    );

    if (!isValid) {
      throw new Error('Invalid email or password');
    }

    // Set user and persist
    setUser(MOCK_USER);
    localStorage.setItem('cmp_user', JSON.stringify(MOCK_USER));
  };

  const logout = async () => {
    setUser(null);
    localStorage.removeItem('cmp_user');
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        login,
        logout,
        isAuthenticated: !!user,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
