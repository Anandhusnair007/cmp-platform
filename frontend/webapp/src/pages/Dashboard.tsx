import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import toast from 'react-hot-toast';
import Modal from '../components/Modal';
import { DashboardSkeleton } from '../components/Skeleton';
import {
  ShieldCheckIcon,
  ServerIcon,
  ClockIcon,
  ExclamationTriangleIcon,
  ArrowPathIcon,
  PlusIcon,
  CheckCircleIcon,
  BoltIcon,
  ChartBarIcon,
} from '@heroicons/react/24/outline';

export default function Dashboard() {
  const [isRequestModalOpen, setIsRequestModalOpen] = useState(false);
  const [isHealthModalOpen, setIsHealthModalOpen] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  // Simulate initial loading
  useEffect(() => {
    const timer = setTimeout(() => setIsLoading(false), 1500);
    return () => clearTimeout(timer);
  }, []);

  const stats = [
    { name: 'Total Certificates', value: '1,284', change: '+12%', icon: ShieldCheckIcon, color: 'text-cyan-400', border: 'border-cyan-600/20' },
    { name: 'Expiring Soon', value: '23', change: '+2', icon: ClockIcon, color: 'text-amber-400', border: 'border-amber-600/20' },
    { name: 'Active Agents', value: '45/48', change: '98%', icon: ServerIcon, color: 'text-emerald-400', border: 'border-emerald-600/20' },
    { name: 'Compliance Score', value: '94%', change: '+1.2%', icon: ChartBarIcon, color: 'text-blue-400', border: 'border-blue-600/20' },
  ];

  const expiringCerts = [
    { id: 1, name: 'api.production.company.com', days: 5, issuer: 'DigiCert' },
    { id: 2, name: 'web.staging.company.com', days: 12, issuer: "Let's Encrypt" },
    { id: 3, name: 'auth.company.com', days: 25, issuer: 'DigiCert' },
    { id: 4, name: 'db.internal.local', days: 28, issuer: 'Internal CA' },
  ];

  const recentActivity = [
    { id: 1, action: 'Certificate Renewed', target: 'mail.company.com', time: '2 mins ago', icon: ArrowPathIcon, color: 'text-emerald-400' },
    { id: 2, action: 'New Agent Connected', target: 'prod-worker-04', time: '15 mins ago', icon: ServerIcon, color: 'text-blue-400' },
    { id: 3, action: 'Certificate Revoked', target: 'dev-test.local', time: '1 hour ago', icon: ExclamationTriangleIcon, color: 'text-red-400' },
  ];

  const handleQuickAction = (action: string) => {
    if (action === 'request') {
      setIsRequestModalOpen(true);
    } else if (action === 'health') {
      setIsHealthModalOpen(true);
    } else {
      toast('Feature coming soon!', { icon: '🚧' });
    }
  };

  if (isLoading) {
    return <DashboardSkeleton />;
  }

  return (
    <div className="space-y-6 fade-in-up">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-3xl font-bold gradient-text neon-glow">Security Overview</h1>
          <p className="text-slate-400 mt-1">Real-time monitoring and certificate management</p>
        </div>
        <div className="flex gap-3">
          <button
            onClick={() => handleQuickAction('health')}
            className="btn-secondary flex items-center gap-2"
          >
            <BoltIcon className="w-5 h-5" />
            System Health
          </button>
          <button
            onClick={() => handleQuickAction('request')}
            className="btn-primary flex items-center gap-2"
          >
            <PlusIcon className="w-5 h-5" />
            Request Certificate
          </button>
        </div>
      </div>

      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        {stats.map((stat) => (
          <div key={stat.name} className={`glass rounded-xl p-6 border ${stat.border} card-hover`}>
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm text-slate-400 font-medium">{stat.name}</p>
                <p className={`text-3xl font-bold mt-1 ${stat.color}`}>{stat.value}</p>
              </div>
              <stat.icon className={`h-12 w-12 opacity-50 ${stat.color}`} />
            </div>
            <div className="mt-4 flex items-center text-sm">
              <span className="text-emerald-400 font-medium">{stat.change}</span>
              <span className="text-slate-500 ml-2">vs last month</span>
            </div>
          </div>
        ))}
      </div>

      {/* Main Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Expiring Certificates */}
        <div className="lg:col-span-2 glass rounded-xl border border-slate-700/50 shadow-cyber overflow-hidden">
          <div className="p-6 border-b border-slate-700/50 bg-slate-900/30 flex justify-between items-center">
            <h2 className="text-lg font-bold text-white flex items-center gap-2">
              <ClockIcon className="h-5 w-5 text-amber-400" />
              Expiring Soon
            </h2>
            <Link to="/inventory" className="text-sm text-cyan-400 hover:text-cyan-300 font-medium">
              View All
            </Link>
          </div>
          <div className="overflow-x-auto">
            <table className="table-cyber">
              <thead>
                <tr>
                  <th>Domain</th>
                  <th>Issuer</th>
                  <th>Expires In</th>
                  <th>Action</th>
                </tr>
              </thead>
              <tbody>
                {expiringCerts.map((cert) => (
                  <tr key={cert.id}>
                    <td className="font-medium text-white">{cert.name}</td>
                    <td className="text-slate-400">{cert.issuer}</td>
                    <td>
                      <span className={`badge ${cert.days <= 7 ? 'badge-danger' : 'badge-warning'}`}>
                        {cert.days} days
                      </span>
                    </td>
                    <td>
                      <Link to={`/certs/${cert.id}`} className="text-cyan-400 hover:text-cyan-300 text-sm font-medium">
                        Renew
                      </Link>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>

        {/* Quick Actions & Activity */}
        <div className="space-y-6">
          {/* Quick Actions Panel */}
          <div className="glass rounded-xl p-6 border border-slate-700/50">
            <h2 className="text-lg font-bold text-white mb-4 flex items-center gap-2">
              <BoltIcon className="h-5 w-5 text-cyan-400" />
              Quick Actions
            </h2>
            <div className="grid grid-cols-2 gap-3">
              <button
                onClick={() => handleQuickAction('request')}
                className="p-3 rounded-lg bg-slate-800 hover:bg-slate-700 border border-slate-600 transition-all text-center group"
              >
                <PlusIcon className="h-6 w-6 text-cyan-400 mx-auto mb-2 group-hover:scale-110 transition-transform" />
                <span className="text-xs font-medium text-slate-300">New Cert</span>
              </button>
              <Link
                to="/agents"
                className="p-3 rounded-lg bg-slate-800 hover:bg-slate-700 border border-slate-600 transition-all text-center group"
              >
                <ServerIcon className="h-6 w-6 text-emerald-400 mx-auto mb-2 group-hover:scale-110 transition-transform" />
                <span className="text-xs font-medium text-slate-300">Add Agent</span>
              </Link>
              <Link
                to="/inventory"
                className="p-3 rounded-lg bg-slate-800 hover:bg-slate-700 border border-slate-600 transition-all text-center group"
              >
                <ArrowPathIcon className="h-6 w-6 text-amber-400 mx-auto mb-2 group-hover:scale-110 transition-transform" />
                <span className="text-xs font-medium text-slate-300">Renewals</span>
              </Link>
              <button
                onClick={() => handleQuickAction('scan')}
                className="p-3 rounded-lg bg-slate-800 hover:bg-slate-700 border border-slate-600 transition-all text-center group"
              >
                <ShieldCheckIcon className="h-6 w-6 text-blue-400 mx-auto mb-2 group-hover:scale-110 transition-transform" />
                <span className="text-xs font-medium text-slate-300">Scan Net</span>
              </button>
            </div>
          </div>

          {/* Recent Activity */}
          <div className="glass rounded-xl p-6 border border-slate-700/50">
            <h2 className="text-lg font-bold text-white mb-4">Recent Activity</h2>
            <div className="space-y-4">
              {recentActivity.map((activity) => (
                <div key={activity.id} className="flex items-start gap-3">
                  <div className={`p-2 rounded-lg bg-slate-900/50 border border-slate-700/50`}>
                    <activity.icon className={`h-4 w-4 ${activity.color}`} />
                  </div>
                  <div>
                    <p className="text-sm font-medium text-white">{activity.action}</p>
                    <p className="text-xs text-slate-400">{activity.target}</p>
                    <p className="text-xs text-slate-500 mt-1">{activity.time}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>

      {/* Request Modal (Simplified) */}
      <Modal
        isOpen={isRequestModalOpen}
        onClose={() => setIsRequestModalOpen(false)}
        title="Request New Certificate"
      >
        <div className="space-y-4">
          <p className="text-slate-300 text-sm">
            Start the wizard to request a new SSL/TLS certificate for your domain.
          </p>
          <div className="space-y-2">
            <label className="text-xs font-semibold text-slate-400 uppercase">Domain Name</label>
            <input type="text" className="input-cyber" placeholder="e.g., app.company.com" />
          </div>
          <div className="flex justify-end gap-3 mt-4">
            <button onClick={() => setIsRequestModalOpen(false)} className="btn-secondary">Cancel</button>
            <Link to="/request" className="btn-primary">Continue to Wizard</Link>
          </div>
        </div>
      </Modal>

      {/* System Health Modal */}
      <Modal
        isOpen={isHealthModalOpen}
        onClose={() => setIsHealthModalOpen(false)}
        title="System Health Status"
      >
        <div className="space-y-4">
          <div className="flex items-center justify-between p-3 bg-emerald-900/20 border border-emerald-600/30 rounded-lg">
            <div className="flex items-center gap-3">
              <CheckCircleIcon className="h-6 w-6 text-emerald-400" />
              <div>
                <p className="font-bold text-white">All Systems Operational</p>
                <p className="text-xs text-emerald-400">99.99% Uptime</p>
              </div>
            </div>
            <div className="h-2 w-2 bg-emerald-400 rounded-full animate-pulse"></div>
          </div>

          <div className="space-y-2">
            <div className="flex justify-between text-sm">
              <span className="text-slate-400">API Latency</span>
              <span className="text-emerald-400 font-mono">24ms</span>
            </div>
            <div className="w-full bg-slate-800 rounded-full h-1.5">
              <div className="bg-emerald-400 h-1.5 rounded-full" style={{ width: '15%' }}></div>
            </div>
          </div>

          <div className="space-y-2">
            <div className="flex justify-between text-sm">
              <span className="text-slate-400">Database Load</span>
              <span className="text-blue-400 font-mono">12%</span>
            </div>
            <div className="w-full bg-slate-800 rounded-full h-1.5">
              <div className="bg-blue-400 h-1.5 rounded-full" style={{ width: '12%' }}></div>
            </div>
          </div>

          <div className="flex justify-end mt-4">
            <button onClick={() => setIsHealthModalOpen(false)} className="btn-secondary">Close</button>
          </div>
        </div>
      </Modal>
    </div>
  );
}