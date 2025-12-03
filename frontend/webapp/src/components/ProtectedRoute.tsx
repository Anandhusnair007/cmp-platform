import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

interface ProtectedRouteProps {
  children: React.ReactNode;
}

export default function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { isAuthenticated, loading } = useAuth();

  // DEMO MODE: Allow access for showcase purposes
  const DEMO_MODE = false; // Set to false for production auth

  if (loading && !DEMO_MODE) {
    return (
      <div className="min-h-screen flex items-center justify-center animated-gradient">
        <div className="spinner w-12 h-12"></div>
      </div>
    );
  }

  if (!isAuthenticated && !DEMO_MODE) {
    return <Navigate to="/login" replace />;
  }

  return <>{children}</>;
}
