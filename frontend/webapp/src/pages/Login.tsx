import React, { useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import toast from 'react-hot-toast';
import { ShieldCheckIcon, LockClosedIcon, ArrowRightIcon } from '@heroicons/react/24/outline';

export default function Login() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      await login(email, password);
      toast.success('🔐 Login successful! Welcome to CMP', {
        duration: 2000,
        icon: '✅',
      });
      // Use window.location to force reload and pick up auth state
      setTimeout(() => {
        window.location.href = '/';
      }, 500);
    } catch (error: any) {
      const errorMessage = error.message || error.response?.data?.error || 'Authentication failed. Please check your credentials.';
      toast.error(`❌ ${errorMessage}`, {
        duration: 4000,
      });
      console.error('Login error:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center cyber-grid animated-gradient px-4 relative overflow-hidden">
      {/* Floating Background Elements */}
      <div className="absolute inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-20 left-10 w-72 h-72 bg-cyan-500 rounded-full mix-blend-multiply filter blur-3xl opacity-10 animate-pulse"></div>
        <div className="absolute bottom-20 right-10 w-72 h-72 bg-blue-500
 rounded-full mix-blend-multiply filter blur-3xl opacity-10 animate-pulse" style={{ animationDelay: '2s' }}></div>
        <div className="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-indigo-500 rounded-full mix-blend-multiply filter blur-3xl opacity-5 animate-pulse" style={{ animationDelay: '4s' }}></div>
      </div>

      <div className="max-w-md w-full space-y-8 relative z-10">
        {/* Logo and Header */}
        <div className="text-center fade-in-up">
          <div className="flex justify-center mb-6">
            <div className="relative">
              <div className="absolute inset-0 bg-cyan-500 blur-2xl opacity-30 shield-pulse"></div>
              <ShieldCheckIcon className="h-20 w-20 text-cyan-400 relative z-10" />
            </div>
          </div>
          <h1 className="text-4xl font-bold gradient-text neon-glow mb-2">
            Certificate Management
          </h1>
          <p className="text-xl font-semibold text-cyan-400 mb-2">
            Platform
          </p>
          <p className="text-sm text-slate-400 font-medium">
            Enterprise-Grade Security Infrastructure
          </p>
          <div className="flex items-center justify-center gap-2 mt-3 text-xs text-slate-500">
            <span className="flex items-center gap-1">
              <div className="w-2 h-2 bg-emerald-500 rounded-full pulse-glow"></div>
              Secure
            </span>
            <span>•</span>
            <span className="flex items-center gap-1">
              <div className="w-2 h-2 bg-cyan-500 rounded-full pulse-glow"></div>
              Encrypted
            </span>
            <span>•</span>
            <span className="flex items-center gap-1">
              <div className="w-2 h-2 bg-blue-500 rounded-full pulse-glow"></div>
              Compliant
            </span>
          </div>
        </div>

        {/* Login Form */}
        <div className="glass rounded-2xl shadow-cyber-lg p-8 fade-in-up" style={{ animationDelay: '0.2s' }}>
          <div className="mb-6">
            <h2 className="text-2xl font-bold text-white flex items-center gap-2">
              <LockClosedIcon className="h-6 w-6 text-cyan-400" />
              Secure Access
            </h2>
            <p className="text-sm text-slate-400 mt-1">
              Sign in to access your certificate dashboard
            </p>
          </div>

          <form onSubmit={handleSubmit} className="space-y-5">
            <div>
              <label htmlFor="email" className="block text-sm font-semibold text-slate-300 mb-2">
                Email Address
              </label>
              <input
                id="email"
                name="email"
                type="email"
                autoComplete="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="input-cyber"
                placeholder="admin@example.com"
              />
            </div>

            <div>
              <label htmlFor="password" className="block text-sm font-semibold text-slate-300 mb-2">
                Password
              </label>
              <input
                id="password"
                name="password"
                type="password"
                autoComplete="current-password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="input-cyber"
                placeholder="Enter your password"
              />
            </div>

            <div className="flex items-center justify-between text-sm">
              <label className="flex items-center text-slate-400 hover:text-slate-300 cursor-pointer">
                <input
                  type="checkbox"
                  className="mr-2 rounded border-slate-600 bg-slate-800 text-cyan-500 focus:ring-cyan-500 focus:ring-offset-slate-900"
                />
                Remember me
              </label>
              <a href="#" className="text-cyan-400 hover:text-cyan-300 font-medium">
                Forgot password?
              </a>
            </div>

            <button
              type="submit"
              disabled={loading}
              className="w-full btn-primary flex items-center justify-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed group"
            >
              {loading ? (
                <>
                  <div className="spinner w-5 h-5"></div>
                  <span>Authenticating...</span>
                </>
              ) : (
                <>
                  <span>Sign In</span>
                  <ArrowRightIcon className="h-5 w-5 group-hover:translate-x-1 transition-transform" />
                </>
              )}
            </button>
          </form>

          {/* Demo Credentials */}
          <div className="mt-6 p-4 bg-slate-900/50 border border-slate-700 rounded-lg">
            <p className="text-xs font-semibold text-cyan-400 mb-2 flex items-center gap-2">
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Demo Credentials
            </p>
            <div className="space-y-1 text-xs text-slate-400 font-mono">
              <p><span className="text-slate-500">Email:</span> admin@example.com</p>
              <p><span className="text-slate-500">Password:</span> admin or 9895</p>
            </div>
          </div>

          {/* Security Info */}
          <div className="mt-6 pt-6 border-t border-slate-700">
            <div className="flex items-center justify-center gap-6 text-xs text-slate-500">
              <div className="flex items-center gap-1">
                <svg className="w-4 h-4 text-emerald-500" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M2.166 4.999A11.954 11.954 0 0010 1.944 11.954 11.954 0 0017.834 5c.11.65.166 1.32.166 2.001 0 5.225-3.34 9.67-8 11.317C5.34 16.67 2 12.225 2 7c0-.682.057-1.35.166-2.001zm11.541 3.708a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                </svg>
                <span>SSL Secured</span>
              </div>
              <div className="flex items-center gap-1">
                <svg className="w-4 h-4 text-cyan-500" fill="currentColor" viewBox="0 0 20 20">
                  <path fillRule="evenodd" d="M18 8a6 6 0 01-7.743 5.743L10 14l-1 1-1 1H6v2H2v-4l4.257-4.257A6 6 0 1118 8zm-6-4a1 1 0 100 2 2 2 0 012 2 1 1 0 102 0 4 4 0 00-4-4z" clipRule="evenodd" />
                </svg>
                <span>Encrypted</span>
              </div>
              <div className="flex items-center gap-1">
                <svg className="w-4 h-4 text-blue-500" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M9 2a1 1 0 000 2h2a1 1 0 100-2H9z" />
                  <path fillRule="evenodd" d="M4 5a2 2 0 012-2 3 3 0 003 3h2a3 3 0 003-3 2 2 0 012 2v11a2 2 0 01-2 2H6a2 2 0 01-2-2V5zm9.707 5.707a1 1 0 00-1.414-1.414L9 12.586l-1.293-1.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                </svg>
                <span>Compliant</span>
              </div>
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="text-center text-xs text-slate-500 fade-in-up" style={{ animationDelay: '0.4s' }}>
          <p>© 2024 Certificate Management Platform. All rights reserved.</p>
          <p className="mt-1">Enterprise Security Solutions</p>
        </div>
      </div>
    </div>
  );
}
