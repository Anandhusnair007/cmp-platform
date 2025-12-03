import React, { useState, useEffect, useRef } from 'react';
import { format } from 'date-fns';
import toast from 'react-hot-toast';
import Modal from '../components/Modal';
import {
  ServerIcon,
  SignalIcon,
  CpuChipIcon,
  ClockIcon,
  CheckCircleIcon,
  XCircleIcon,
  ArrowPathIcon,
  PlusIcon,
  CommandLineIcon,
  CloudArrowUpIcon,
} from '@heroicons/react/24/outline';

export default function Agents() {
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedAgent, setSelectedAgent] = useState<any>(null);
  const [isLogModalOpen, setIsLogModalOpen] = useState(false);
  const [isDeployModalOpen, setIsDeployModalOpen] = useState(false);
  const [logs, setLogs] = useState<string[]>([]);
  const [isDeploying, setIsDeploying] = useState(false);
  const logsEndRef = useRef<HTMLDivElement>(null);

  // Mock agent data
  const mockAgents = [
    {
      id: '1',
      hostname: 'prod-web-01.company.com',
      ip: '10.0.1.101',
      status: 'online',
      last_checkin: new Date(),
      os: 'Ubuntu 22.04 LTS',
      version: '2.5.1',
      certs_managed: 12,
      cpu: '24%',
      memory: '2.1 GB / 8 GB',
      uptime: '45 days',
    },
    {
      id: '2',
      hostname: 'prod-api-01.company.com',
      ip: '10.0.1.102',
      status: 'online',
      last_checkin: new Date(Date.now() - 2 * 60 * 1000),
      os: 'Ubuntu 22.04 LTS',
      version: '2.5.1',
      certs_managed: 8,
      cpu: '18%',
      memory: '1.8 GB / 8 GB',
      uptime: '45 days',
    },
    {
      id: '3',
      hostname: 'staging-web-01.company.com',
      ip: '10.0.2.101',
      status: 'offline',
      last_checkin: new Date(Date.now() - 3600000),
      os: 'Ubuntu 20.04 LTS',
      version: '2.4.8',
      certs_managed: 4,
      cpu: 'N/A',
      memory: 'N/A',
      uptime: 'Offline',
    },
    {
      id: '4',
      hostname: 'prod-db-01.company.com',
      ip: '10.0.1.103',
      status: 'online',
      last_checkin: new Date(Date.now() - 1 * 60 * 1000),
      os: 'Red Hat Enterprise Linux 8',
      version: '2.5.1',
      certs_managed: 3,
      cpu: '32%',
      memory: '4.2 GB / 16 GB',
      uptime: '120 days',
    },
    {
      id: '5',
      hostname: 'prod-lb-01.company.com',
      ip: '10.0.1.104',
      status: 'online',
      last_checkin: new Date(),
      os: 'Ubuntu 22.04 LTS',
      version: '2.5.1',
      certs_managed: 15,
      cpu: '42%',
      memory: '3.1 GB / 8 GB',
      uptime: '30 days',
    },
  ];

  const handleViewLogs = (agent: any) => {
    setSelectedAgent(agent);
    // Generate mock logs
    const mockLogs = [
      `[${format(new Date(Date.now() - 10000), 'HH:mm:ss')}] INFO: Agent service started v2.5.1`,
      `[${format(new Date(Date.now() - 9000), 'HH:mm:ss')}] INFO: Connected to control plane at 10.0.0.5:8443`,
      `[${format(new Date(Date.now() - 8000), 'HH:mm:ss')}] INFO: Syncing certificate inventory...`,
      `[${format(new Date(Date.now() - 7000), 'HH:mm:ss')}] INFO: Found 12 managed certificates`,
      `[${format(new Date(Date.now() - 6000), 'HH:mm:ss')}] INFO: Checking for renewals...`,
      `[${format(new Date(Date.now() - 5000), 'HH:mm:ss')}] INFO: No renewals required`,
      `[${format(new Date(Date.now() - 4000), 'HH:mm:ss')}] INFO: Heartbeat sent (CPU: ${agent.cpu})`,
      `[${format(new Date(Date.now() - 1000), 'HH:mm:ss')}] INFO: Waiting for commands...`,
    ];
    setLogs(mockLogs);
    setIsLogModalOpen(true);
  };

  const handleDeploy = (agent: any) => {
    setSelectedAgent(agent);
    setIsDeployModalOpen(true);
  };

  const confirmDeploy = async () => {
    setIsDeploying(true);
    await new Promise(resolve => setTimeout(resolve, 2000));
    setIsDeploying(false);
    setIsDeployModalOpen(false);
    toast.success(`Certificate deployed successfully to ${selectedAgent.hostname}`, {
      icon: '🚀',
      duration: 4000,
    });
  };

  const filteredAgents = mockAgents.filter(agent =>
    agent.hostname.toLowerCase().includes(searchTerm.toLowerCase()) ||
    agent.ip.includes(searchTerm)
  );

  const stats = {
    total: mockAgents.length,
    online: mockAgents.filter(a => a.status === 'online').length,
    offline: mockAgents.filter(a => a.status === 'offline').length,
    certs: mockAgents.reduce((sum, a) => sum + a.certs_managed, 0),
  };

  return (
    <div className="space-y-6 fade-in-up">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-3xl font-bold gradient-text neon-glow">Deployment Agents</h1>
          <p className="text-slate-400 mt-1">Monitor and manage certificate deployment agents</p>
        </div>
        <button className="btn-primary">
          <span className="flex items-center gap-2">
            <PlusIcon className="w-5 h-5" />
            Add Agent
          </span>
        </button>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <div className="glass rounded-xl p-6 border border-cyan-600/20">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-slate-400 font-medium">Total Agents</p>
              <p className="text-3xl font-bold text-cyan-400 mt-1">{stats.total}</p>
            </div>
            <ServerIcon className="h-12 w-12 text-cyan-400 opacity-50" />
          </div>
        </div>
        <div className="glass rounded-xl p-6 border border-emerald-600/20">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-slate-400 font-medium">Online</p>
              <p className="text-3xl font-bold text-emerald-400 mt-1">{stats.online}</p>
            </div>
            <CheckCircleIcon className="h-12 w-12 text-emerald-400 opacity-50" />
          </div>
        </div>
        <div className="glass rounded-xl p-6 border border-red-600/20">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-slate-400 font-medium">Offline</p>
              <p className="text-3xl font-bold text-red-400 mt-1">{stats.offline}</p>
            </div>
            <XCircleIcon className="h-12 w-12 text-red-400 opacity-50" />
          </div>
        </div>
        <div className="glass rounded-xl p-6 border border-blue-600/20">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-slate-400 font-medium">Certs Managed</p>
              <p className="text-3xl font-bold text-blue-400 mt-1">{stats.certs}</p>
            </div>
            <CpuChipIcon className="h-12 w-12 text-blue-400 opacity-50" />
          </div>
        </div>
      </div>

      {/* Search Bar */}
      <div className="glass rounded-xl p-4 border border-slate-700/50">
        <div className="relative">
          <svg className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-slate-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            type="text"
            placeholder="Search agents by hostname or IP..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="input-cyber pl-10"
          />
        </div>
      </div>

      {/* Agents Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {filteredAgents.map((agent) => (
          <div key={agent.id} className="glass rounded-xl shadow-cyber overflow-hidden card-hover border border-slate-700/50">
            {/* Agent Header */}
            <div className={`p-4 border-b border-slate-700/50 ${agent.status === 'online' ? 'bg-emerald-900/10' : 'bg-red-900/10'
              }`}>
              <div className="flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className={`p-2 rounded-lg ${agent.status === 'online' ? 'bg-emerald-900/30' : 'bg-red-900/30'
                    }`}>
                    <ServerIcon className={`h-6 w-6 ${agent.status === 'online' ? 'text-emerald-400' : 'text-red-400'
                      }`} />
                  </div>
                  <div>
                    <h3 className="text-lg font-bold text-white">{agent.hostname}</h3>
                    <p className="text-sm text-slate-400 font-mono">{agent.ip}</p>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <span className={`badge ${agent.status === 'online' ? 'badge-success' : 'badge-danger'
                    }`}>
                    {agent.status}
                  </span>
                  {agent.status === 'online' && (
                    <div className="w-3 h-3 bg-emerald-500 rounded-full pulse-glow"></div>
                  )}
                </div>
              </div>
            </div>

            {/* Agent Details */}
            <div className="p-6 space-y-4">
              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-xs text-slate-500 font-medium mb-1">Operating System</p>
                  <p className="text-sm text-white font-medium">{agent.os}</p>
                </div>
                <div>
                  <p className="text-xs text-slate-500 font-medium mb-1">Agent Version</p>
                  <p className="text-sm text-white font-medium font-mono">{agent.version}</p>
                </div>
                <div>
                  <p className="text-xs text-slate-500 font-medium mb-1">Certificates</p>
                  <p className="text-sm text-cyan-400 font-bold">{agent.certs_managed} managed</p>
                </div>
                <div>
                  <p className="text-xs text-slate-500 font-medium mb-1">Uptime</p>
                  <p className="text-sm text-white font-medium">{agent.uptime}</p>
                </div>
              </div>

              {agent.status === 'online' && (
                <div className="grid grid-cols-2 gap-4 pt-4 border-t border-slate-700/50">
                  <div>
                    <p className="text-xs text-slate-500 font-medium mb-1 flex items-center gap-1">
                      <CpuChipIcon className="h-3 w-3" />
                      CPU Usage
                    </p>
                    <p className="text-sm text-white font-medium">{agent.cpu}</p>
                  </div>
                  <div>
                    <p className="text-xs text-slate-500 font-medium mb-1 flex items-center gap-1">
                      <SignalIcon className="h-3 w-3" />
                      Memory
                    </p>
                    <p className="text-sm text-white font-medium">{agent.memory}</p>
                  </div>
                </div>
              )}

              <div className="flex items-center justify-between pt-4 border-t border-slate-700/50">
                <div className="flex items-center gap-2 text-xs text-slate-400">
                  <ClockIcon className="h-4 w-4" />
                  <span>Last check-in: {format(new Date(agent.last_checkin), 'MMM dd, HH:mm')}</span>
                </div>
                <button className="text-cyan-400 hover:text-cyan-300 text-sm font-medium flex items-center gap-1">
                  <ArrowPathIcon className="h-4 w-4" />
                  Refresh
                </button>
              </div>

              {/* Actions */}
              <div className="flex gap-2 pt-4">
                <button
                  onClick={() => handleViewLogs(agent)}
                  className="flex-1 px-4 py-2 bg-slate-800 hover:bg-slate-700 border border-slate-600 rounded-lg text-sm font-medium text-white transition-all flex items-center justify-center gap-2"
                >
                  <CommandLineIcon className="h-4 w-4" />
                  View Logs
                </button>
                <button
                  onClick={() => handleDeploy(agent)}
                  disabled={agent.status === 'offline'}
                  className="flex-1 px-4 py-2 bg-cyan-600 hover:bg-cyan-500 disabled:bg-slate-700 disabled:text-slate-500 rounded-lg text-sm font-medium text-white transition-all flex items-center justify-center gap-2"
                >
                  <CloudArrowUpIcon className="h-4 w-4" />
                  Deploy Cert
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Logs Modal */}
      <Modal
        isOpen={isLogModalOpen}
        onClose={() => setIsLogModalOpen(false)}
        title={`Agent Logs: ${selectedAgent?.hostname}`}
        size="lg"
      >
        <div className="bg-slate-950 rounded-lg border border-slate-800 p-4 font-mono text-xs h-96 overflow-y-auto">
          {logs.map((log, index) => (
            <div key={index} className="text-emerald-400 mb-1">
              <span className="text-slate-500 mr-2">$</span>
              {log}
            </div>
          ))}
          <div ref={logsEndRef} />
        </div>
        <div className="flex justify-end mt-4">
          <button
            onClick={() => setIsLogModalOpen(false)}
            className="btn-secondary"
          >
            Close
          </button>
        </div>
      </Modal>

      {/* Deploy Modal */}
      <Modal
        isOpen={isDeployModalOpen}
        onClose={() => setIsDeployModalOpen(false)}
        title="Deploy Certificate"
      >
        <div className="space-y-4">
          <p className="text-sm text-slate-300">
            Select a certificate to deploy to <span className="font-bold text-white">{selectedAgent?.hostname}</span>.
          </p>

          <div className="space-y-2">
            <label className="text-xs font-semibold text-slate-400 uppercase">Certificate</label>
            <select className="input-cyber">
              <option>api.production.company.com (Expires in 5 days)</option>
              <option>web.staging.company.com (Expires in 12 days)</option>
              <option>auth.company.com (Expires in 25 days)</option>
            </select>
          </div>

          <div className="space-y-2">
            <label className="text-xs font-semibold text-slate-400 uppercase">Target Path</label>
            <input type="text" className="input-cyber" defaultValue="/etc/ssl/certs/" />
          </div>

          <div className="space-y-2">
            <label className="text-xs font-semibold text-slate-400 uppercase">Post-Deploy Command</label>
            <input type="text" className="input-cyber" defaultValue="systemctl reload nginx" />
          </div>

          <div className="flex justify-end gap-3 mt-6">
            <button
              onClick={() => setIsDeployModalOpen(false)}
              className="px-4 py-2 text-slate-300 hover:text-white font-medium"
            >
              Cancel
            </button>
            <button
              onClick={confirmDeploy}
              disabled={isDeploying}
              className="btn-primary flex items-center gap-2"
            >
              {isDeploying ? (
                <>
                  <div className="spinner w-4 h-4"></div>
                  Deploying...
                </>
              ) : (
                <>
                  <CloudArrowUpIcon className="h-5 w-5" />
                  Deploy
                </>
              )}
            </button>
          </div>
        </div>
      </Modal>

      {filteredAgents.length === 0 && (
        <div className="glass rounded-xl p-12 text-center border border-slate-700/50">
          <ServerIcon className="h-16 w-16 mx-auto mb-4 text-slate-600" />
          <p className="text-lg font-medium text-slate-400">No agents found</p>
          <p className="text-sm text-slate-500 mt-1">Try adjusting your search</p>
        </div>
      )}
    </div>
  );
}