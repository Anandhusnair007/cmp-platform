import React, { useState } from 'react';
import { Link, useLocation, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import {
  HomeIcon,
  FolderIcon,
  PlusCircleIcon,
  ServerIcon,
  Cog6ToothIcon,
  ArrowRightOnRectangleIcon,
  Bars3Icon,
  XMarkIcon,
  ShieldCheckIcon,
  BellIcon,
  MagnifyingGlassIcon,
} from '@heroicons/react/24/outline';

interface LayoutProps {
  children: React.ReactNode;
}

const navigation = [
  { name: 'Dashboard', href: '/', icon: HomeIcon },
  { name: 'Inventory', href: '/inventory', icon: FolderIcon },
  { name: 'Request Certificate', href: '/request', icon: PlusCircleIcon },
  { name: 'Agents', href: '/agents', icon: ServerIcon },
  { name: 'Admin', href: '/admin', icon: Cog6ToothIcon },
];

export default function Layout({ children }: LayoutProps) {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const { user, logout } = useAuth();
  const location = useLocation();
  const navigate = useNavigate();

  // Use a default user for demo mode
  const displayUser = user || { name: 'Admin User', email: 'admin@example.com', roles: ['admin'], team: 'Security Operations' };

  const handleLogout = async () => {
    await logout();
    navigate('/login');
  };

  return (
    <div className="min-h-screen animated-gradient cyber-grid">
      {/* Mobile sidebar */}
      <div
        className={`fixed inset-0 z-50 lg:hidden ${sidebarOpen ? 'block' : 'hidden'}`}
      >
        <div className="fixed inset-0 bg-slate-950/90 backdrop-blur-sm" onClick={() => setSidebarOpen(false)} />
        <div className="fixed inset-y-0 left-0 w-72 glass border-r border-slate-700/50">
          <SidebarContent
            user={displayUser}
            location={location}
            onLogout={handleLogout}
            onClose={() => setSidebarOpen(false)}
          />
        </div>
      </div>

      {/* Desktop sidebar */}
      <div className="hidden lg:fixed lg:inset-y-0 lg:z-50 lg:flex lg:w-72 lg:flex-col">
        <div className="glass border-r border-slate-700/50 flex grow flex-col gap-y-5 overflow-y-auto px-6 pb-4">
          <SidebarContent
            user={displayUser}
            location={location}
            onLogout={handleLogout}
          />
        </div>
      </div>

      {/* Main content */}
      <div className="lg:pl-72">
        {/* Top bar */}
        <div className="sticky top-0 z-40 glass border-b border-slate-700/50 shadow-cyber">
          <div className="flex h-16 shrink-0 items-center gap-x-4 px-4 sm:gap-x-6 sm:px-6 lg:px-8">
            <button
              type="button"
              className="-m-2.5 p-2.5 text-slate-400 hover:text-slate-300 lg:hidden"
              onClick={() => setSidebarOpen(true)}
            >
              <Bars3Icon className="h-6 w-6" />
            </button>

            {/* Search bar */}
            <div className="flex flex-1 gap-x-4 self-stretch lg:gap-x-6">
              <form className="relative flex flex-1 max-w-md" action="#">
                <label htmlFor="search-field" className="sr-only">
                  Search
                </label>
                <MagnifyingGlassIcon
                  className="pointer-events-none absolute inset-y-0 left-3 h-full w-5 text-slate-500"
                  aria-hidden="true"
                />
                <input
                  id="search-field"
                  className="block h-full w-full border-0 py-0 pl-10 pr-0 bg-transparent text-white placeholder:text-slate-500 focus:ring-0 sm:text-sm"
                  placeholder="Search certificates, agents..."
                  type="search"
                  name="search"
                />
              </form>

              {/* Right side */}
              <div className="flex items-center gap-x-4 lg:gap-x-6 ml-auto">
                {/* Notifications */}
                <button type="button" className="relative p-2 text-slate-400 hover:text-slate-300">
                  <BellIcon className="h-6 w-6" />
                  <span className="absolute top-1 right-1 w-2 h-2 bg-cyan-500 rounded-full pulse-glow"></span>
                </button>

                <div className="hidden lg:block lg:h-6 lg:w-px lg:bg-slate-700" />

                {/* User info */}
                <div className="flex items-center gap-x-3">
                  <div className="text-sm text-right">
                    <div className="font-semibold text-white">{displayUser.name}</div>
                    <div className="text-slate-400 text-xs">{displayUser.email}</div>
                  </div>
                  <div className="relative">
                    <div className="h-9 w-9 rounded-lg bg-gradient-to-br from-cyan-500 to-blue-600 flex items-center justify-center text-white font-bold">
                      {displayUser.name.charAt(0)}
                    </div>
                    <span className="absolute bottom-0 right-0 w-3 h-3 bg-emerald-500 rounded-full border-2 border-slate-900 pulse-glow"></span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Page content */}
        <main className="py-8 px-4 sm:px-6 lg:px-8">
          {children}
        </main>
      </div>
    </div>
  );
}

interface SidebarContentProps {
  user: { name: string; email: string; roles: string[]; team?: string };
  location: { pathname: string };
  onLogout: () => void;
  onClose?: () => void;
}

function SidebarContent({ user, location, onLogout, onClose }: SidebarContentProps) {
  return (
    <>
      {/* Logo/Header */}
      <div className="flex h-16 shrink-0 items-center border-b border-slate-700/50">
        <div className="flex items-center gap-3">
          <div className="relative">
            <div className="absolute inset-0 bg-cyan-500 blur-xl opacity-50 shield-pulse"></div>
            <ShieldCheckIcon className="h-10 w-10 text-cyan-400 relative z-10" />
          </div>
          <div>
            <h1 className="text-xl font-bold gradient-text">CMP</h1>
            <p className="text-[10px] text-slate-500 font-medium">Certificate Platform</p>
          </div>
        </div>
        {onClose && (
          <button
            type="button"
            className="ml-auto text-slate-400 hover:text-slate-300"
            onClick={onClose}
          >
            <XMarkIcon className="h-6 w-6" />
          </button>
        )}
      </div>

      {/* User Card */}
      <div className="glass rounded-xl p-4 border border-cyan-600/20 shadow-cyber">
        <div className="flex items-center gap-3 mb-3">
          <div className="h-12 w-12 rounded-xl bg-gradient-to-br from-cyan-500 to-blue-600 flex items-center justify-center text-white font-bold text-lg shadow-lg">
            {user.name.charAt(0)}
          </div>
          <div className="flex-1 min-w-0">
            <p className="text-sm font-semibold text-white truncate">{user.name}</p>
            <p className="text-xs text-slate-400 truncate">{user.team || 'Security Team'}</p>
          </div>
        </div>
        <div className="flex gap-2">
          {user.roles.map((role) => (
            <span key={role} className="badge badge-info text-[10px]">
              {role}
            </span>
          ))}
        </div>
      </div>

      {/* Navigation */}
      <nav className="flex flex-1 flex-col mt-2">
        <ul role="list" className="flex flex-1 flex-col gap-y-1">
          {navigation.map((item) => {
            const isActive = location.pathname === item.href;
            return (
              <li key={item.name}>
                <Link
                  to={item.href}
                  onClick={onClose}
                  className={`group flex gap-x-3 rounded-lg p-3 text-sm font-semibold transition-all duration-200 ${isActive
                      ? 'bg-gradient-to-r from-cyan-600/20 to-blue-600/20 text-white border border-cyan-600/30 shadow-cyan-500/20'
                      : 'text-slate-300 hover:text-white hover:bg-slate-800/50 border border-transparent'
                    }`}
                >
                  <item.icon className={`h-5 w-5 shrink-0 ${isActive ? 'text-cyan-400' : 'text-slate-500'}`} />
                  {item.name}
                  {isActive && (
                    <div className="ml-auto">
                      <div className="w-1.5 h-1.5 bg-cyan-400 rounded-full pulse-glow"></div>
                    </div>
                  )}
                </Link>
              </li>
            );
          })}
        </ul>

        {/* System Status */}
        <div className="mt-auto mb-4 p-4 bg-emerald-900/20 border border-emerald-600/30 rounded-lg">
          <div className="flex items-center gap-2 mb-2">
            <div className="w-2 h-2 bg-emerald-500 rounded-full pulse-glow"></div>
            <span className="text-xs font-semibold text-emerald-400">System Status</span>
          </div>
          <p className="text-xs text-slate-400">All services operational</p>
        </div>

        {/* Logout Button */}
        <div className="border-t border-slate-700/50 pt-4">
          <button
            onClick={onLogout}
            className="group flex gap-x-3 rounded-lg p-3 text-sm font-semibold text-slate-300 hover:text-white hover:bg-red-900/20 hover:border-red-600/30 w-full border border-transparent transition-all duration-200"
          >
            <ArrowRightOnRectangleIcon className="h-5 w-5 shrink-0 text-slate-500 group-hover:text-red-400" />
            Logout
          </button>
        </div>
      </nav>
    </>
  );
}
