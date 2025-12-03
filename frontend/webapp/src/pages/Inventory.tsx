import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { format, differenceInDays } from 'date-fns';
import toast from 'react-hot-toast';
import Modal from '../components/Modal';
import {
  MagnifyingGlassIcon,
  FunnelIcon,
  DocumentCheckIcon,
  ClockIcon,
  ShieldCheckIcon,
  ExclamationTriangleIcon,
  ArrowPathIcon,
  CheckCircleIcon,
} from '@heroicons/react/24/outline';

export default function Inventory() {
  const [searchTerm, setSearchTerm] = useState('');
  const [filterStatus, setFilterStatus] = useState('all');
  const [selectedCert, setSelectedCert] = useState<any>(null);
  const [isRenewModalOpen, setIsRenewModalOpen] = useState(false);
  const [isViewModalOpen, setIsViewModalOpen] = useState(false);
  const [isRenewing, setIsRenewing] = useState(false);

  // Mock certificate data
  const [certs, setCerts] = useState([
    {
      id: '1',
      subject: 'api.production.company.com',
      issuer: 'DigiCert',
      not_before: '2024-01-15',
      not_after: new Date(Date.now() + 5 * 24 * 60 * 60 * 1000).toISOString(),
      status: 'active',
      type: 'SSL/TLS',
      algorithm: 'RSA 2048'
    },
    {
      id: '2',
      subject: 'web.staging.company.com',
      issuer: 'Let\'s Encrypt',
      not_before: '2024-02-01',
      not_after: new Date(Date.now() + 12 * 24 * 60 * 60 * 1000).toISOString(),
      status: 'active',
      type: 'SSL/TLS',
      algorithm: 'RSA 2048'
    },
    {
      id: '3',
      subject: 'auth.company.com',
      issuer: 'DigiCert',
      not_before: '2024-01-20',
      not_after: new Date(Date.now() + 25 * 24 * 60 * 60 * 1000).toISOString(),
      status: 'active',
      type: 'SSL/TLS',
      algorithm: 'RSA 4096'
    },
    {
      id: '4',
      subject: 'cdn.company.com',
      issuer: 'Cloudflare',
      not_before: '2024-03-01',
      not_after: new Date(Date.now() + 45 * 24 * 60 * 60 * 1000).toISOString(),
      status: 'active',
      type: 'SSL/TLS',
      algorithm: 'ECDSA P-256'
    },
    {
      id: '5',
      subject: 'mail.company.com',
      issuer: 'DigiCert',
      not_before: '2023-12-01',
      not_after: new Date(Date.now() + 90 * 24 * 60 * 60 * 1000).toISOString(),
      status: 'active',
      type: 'SSL/TLS',
      algorithm: 'RSA 2048'
    },
    {
      id: '6',
      subject: 'vpn.company.com',
      issuer: 'Internal CA',
      not_before: '2024-01-10',
      not_after: new Date(Date.now() + 180 * 24 * 60 * 60 * 1000).toISOString(),
      status: 'active',
      type: 'SSL/TLS',
      algorithm: 'RSA 4096'
    },
  ]);

  const handleRenewClick = (cert: any) => {
    setSelectedCert(cert);
    setIsRenewModalOpen(true);
  };

  const handleViewClick = (cert: any) => {
    setSelectedCert(cert);
    setIsViewModalOpen(true);
  };

  const confirmRenew = async () => {
    setIsRenewing(true);
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 2000));

    // Update cert in list
    setCerts(prev => prev.map(c => {
      if (c.id === selectedCert.id) {
        return {
          ...c,
          not_after: new Date(Date.now() + 365 * 24 * 60 * 60 * 1000).toISOString(), // Extend by 1 year
          status: 'active'
        };
      }
      return c;
    }));

    setIsRenewing(false);
    setIsRenewModalOpen(false);
    toast.success(`Successfully renewed certificate for ${selectedCert.subject}`, {
      icon: '✅',
      duration: 4000,
    });
  };

  const filteredCerts = certs.filter(cert => {
    const matchesSearch = cert.subject.toLowerCase().includes(searchTerm.toLowerCase()) ||
      cert.issuer.toLowerCase().includes(searchTerm.toLowerCase());
    const daysLeft = differenceInDays(new Date(cert.not_after), new Date());

    if (filterStatus === 'expiring') return matchesSearch && daysLeft <= 30;
    if (filterStatus === 'critical') return matchesSearch && daysLeft <= 7;
    return matchesSearch;
  });

  const stats = {
    total: certs.length,
    expiring: certs.filter(c => differenceInDays(new Date(c.not_after), new Date()) <= 30).length,
    critical: certs.filter(c => differenceInDays(new Date(c.not_after), new Date()) <= 7).length,
  };

  return (
    <div className="space-y-6 fade-in-up">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <div>
          <h1 className="text-3xl font-bold gradient-text neon-glow">Certificate Inventory</h1>
          <p className="text-slate-400 mt-1">Manage and monitor all SSL/TLS certificates</p>
        </div>
        <Link
          to="/request"
          className="btn-primary"
        >
          <span className="flex items-center gap-2">
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
            Request Certificate
          </span>
        </Link>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div className="glass rounded-xl p-6 border border-cyan-600/20">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-slate-400 font-medium">Total Certificates</p>
              <p className="text-3xl font-bold text-cyan-400 mt-1">{stats.total}</p>
            </div>
            <DocumentCheckIcon className="h-12 w-12 text-cyan-400 opacity-50" />
          </div>
        </div>
        <div className="glass rounded-xl p-6 border border-amber-600/20">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-slate-400 font-medium">Expiring Soon (30d)</p>
              <p className="text-3xl font-bold text-amber-400 mt-1">{stats.expiring}</p>
            </div>
            <ClockIcon className="h-12 w-12 text-amber-400 opacity-50" />
          </div>
        </div>
        <div className="glass rounded-xl p-6 border border-red-600/20">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm text-slate-400 font-medium">Critical (7d)</p>
              <p className="text-3xl font-bold text-red-400 mt-1">{stats.critical}</p>
            </div>
            <ExclamationTriangleIcon className="h-12 w-12 text-red-400 opacity-50" />
          </div>
        </div>
      </div>

      {/* Search and Filter */}
      <div className="glass rounded-xl p-6 border border-slate-700/50">
        <div className="flex flex-col sm:flex-row gap-4">
          <div className="flex-1 relative">
            <MagnifyingGlassIcon className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-slate-500" />
            <input
              type="text"
              placeholder="Search certificates..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="input-cyber pl-10"
            />
          </div>
          <div className="flex gap-2">
            <button
              onClick={() => setFilterStatus('all')}
              className={`px-4 py-2 rounded-lg font-medium transition-all ${filterStatus === 'all'
                  ? 'bg-cyan-600 text-white'
                  : 'bg-slate-800 text-slate-300 hover:bg-slate-700'
                }`}
            >
              All
            </button>
            <button
              onClick={() => setFilterStatus('expiring')}
              className={`px-4 py-2 rounded-lg font-medium transition-all ${filterStatus === 'expiring'
                  ? 'bg-amber-600 text-white'
                  : 'bg-slate-800 text-slate-300 hover:bg-slate-700'
                }`}
            >
              Expiring
            </button>
            <button
              onClick={() => setFilterStatus('critical')}
              className={`px-4 py-2 rounded-lg font-medium transition-all ${filterStatus === 'critical'
                  ? 'bg-red-600 text-white'
                  : 'bg-slate-800 text-slate-300 hover:bg-slate-700'
                }`}
            >
              Critical
            </button>
          </div>
        </div>
      </div>

      {/* Certificates Table */}
      <div className="glass rounded-xl shadow-cyber-lg overflow-hidden">
        <div className="p-6 border-b border-slate-700/50 bg-slate-900/30">
          <div className="flex items-center justify-between">
            <h2 className="text-xl font-bold text-white flex items-center gap-2">
              <FunnelIcon className="h-5 w-5 text-cyan-400" />
              Certificates ({filteredCerts.length})
            </h2>
            <button className="text-slate-400 hover:text-cyan-400 transition-colors">
              <ArrowPathIcon className="h-5 w-5" />
            </button>
          </div>
        </div>

        <div className="overflow-x-auto">
          <table className="table-cyber">
            <thead>
              <tr>
                <th>Domain</th>
                <th>Issuer</th>
                <th>Type</th>
                <th>Algorithm</th>
                <th>Issued</th>
                <th>Expires</th>
                <th>Days Left</th>
                <th>Status</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {filteredCerts.map((cert) => {
                const daysLeft = differenceInDays(new Date(cert.not_after), new Date());
                const isExpiringSoon = daysLeft <= 30;
                const isCritical = daysLeft <= 7;

                return (
                  <tr key={cert.id}>
                    <td>
                      <Link
                        to={`/certs/${cert.id}`}
                        className="font-medium text-white hover:text-cyan-400 transition-colors flex items-center gap-2"
                      >
                        <ShieldCheckIcon className="h-4 w-4 text-cyan-400" />
                        {cert.subject}
                      </Link>
                    </td>
                    <td className="text-slate-300">{cert.issuer}</td>
                    <td className="text-slate-300">{cert.type}</td>
                    <td className="text-slate-300 font-mono text-xs">{cert.algorithm}</td>
                    <td className="text-slate-300">{format(new Date(cert.not_before), 'MMM dd, yyyy')}</td>
                    <td className="text-slate-300">{format(new Date(cert.not_after), 'MMM dd, yyyy')}</td>
                    <td className={`font-semibold ${isCritical ? 'text-red-400' : isExpiringSoon ? 'text-amber-400' : 'text-emerald-400'
                      }`}>
                      {daysLeft} days
                    </td>
                    <td>
                      <span className={`badge ${isCritical ? 'badge-danger' : isExpiringSoon ? 'badge-warning' : 'badge-success'
                        }`}>
                        {isCritical ? 'Critical' : isExpiringSoon ? 'Expiring' : 'Active'}
                      </span>
                    </td>
                    <td>
                      <div className="flex gap-2">
                        <button
                          onClick={() => handleViewClick(cert)}
                          className="text-cyan-400 hover:text-cyan-300 text-sm font-medium"
                        >
                          View
                        </button>
                        <button
                          onClick={() => handleRenewClick(cert)}
                          className="text-amber-400 hover:text-amber-300 text-sm font-medium"
                        >
                          Renew
                        </button>
                      </div>
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>

          {filteredCerts.length === 0 && (
            <div className="p-12 text-center text-slate-500">
              <DocumentCheckIcon className="h-16 w-16 mx-auto mb-4 opacity-30" />
              <p className="text-lg font-medium">No certificates found</p>
              <p className="text-sm mt-1">Try adjusting your search or filters</p>
            </div>
          )}
        </div>
      </div>

      {/* Renew Modal */}
      <Modal
        isOpen={isRenewModalOpen}
        onClose={() => setIsRenewModalOpen(false)}
        title="Renew Certificate"
      >
        <div className="space-y-4">
          <div className="flex items-center gap-3 p-4 bg-amber-900/20 border border-amber-600/30 rounded-lg">
            <ExclamationTriangleIcon className="h-6 w-6 text-amber-400" />
            <div>
              <p className="font-semibold text-amber-400">Confirm Renewal</p>
              <p className="text-sm text-slate-300">
                Are you sure you want to renew the certificate for <span className="font-bold text-white">{selectedCert?.subject}</span>?
              </p>
            </div>
          </div>

          <div className="bg-slate-900/50 p-4 rounded-lg border border-slate-700/50 space-y-2">
            <div className="flex justify-between text-sm">
              <span className="text-slate-400">Current Expiration:</span>
              <span className="text-white font-mono">{selectedCert && format(new Date(selectedCert.not_after), 'MMM dd, yyyy')}</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-slate-400">New Expiration:</span>
              <span className="text-emerald-400 font-mono font-bold">{format(new Date(Date.now() + 365 * 24 * 60 * 60 * 1000), 'MMM dd, yyyy')}</span>
            </div>
          </div>

          <div className="flex justify-end gap-3 mt-6">
            <button
              onClick={() => setIsRenewModalOpen(false)}
              className="px-4 py-2 text-slate-300 hover:text-white font-medium"
            >
              Cancel
            </button>
            <button
              onClick={confirmRenew}
              disabled={isRenewing}
              className="btn-primary flex items-center gap-2"
            >
              {isRenewing ? (
                <>
                  <div className="spinner w-4 h-4"></div>
                  Renewing...
                </>
              ) : (
                <>
                  <ArrowPathIcon className="h-5 w-5" />
                  Confirm Renewal
                </>
              )}
            </button>
          </div>
        </div>
      </Modal>

      {/* View Details Modal */}
      <Modal
        isOpen={isViewModalOpen}
        onClose={() => setIsViewModalOpen(false)}
        title="Certificate Details"
        size="lg"
      >
        {selectedCert && (
          <div className="space-y-6">
            <div className="flex items-center gap-4 p-4 bg-slate-900/50 rounded-lg border border-slate-700/50">
              <div className="p-3 bg-cyan-900/30 rounded-lg">
                <ShieldCheckIcon className="h-8 w-8 text-cyan-400" />
              </div>
              <div>
                <h3 className="text-xl font-bold text-white">{selectedCert.subject}</h3>
                <p className="text-sm text-slate-400">Serial: 7A:4F:2C:1D:9E:5B:3A:8F</p>
              </div>
            </div>

            <div className="grid grid-cols-2 gap-6">
              <div>
                <h4 className="text-sm font-semibold text-slate-400 uppercase tracking-wider mb-3">Issued To</h4>
                <div className="space-y-2">
                  <p className="text-sm"><span className="text-slate-500 w-24 inline-block">Common Name:</span> <span className="text-white">{selectedCert.subject}</span></p>
                  <p className="text-sm"><span className="text-slate-500 w-24 inline-block">Organization:</span> <span className="text-white">Acme Corp</span></p>
                  <p className="text-sm"><span className="text-slate-500 w-24 inline-block">Unit:</span> <span className="text-white">IT Security</span></p>
                </div>
              </div>
              <div>
                <h4 className="text-sm font-semibold text-slate-400 uppercase tracking-wider mb-3">Issued By</h4>
                <div className="space-y-2">
                  <p className="text-sm"><span className="text-slate-500 w-24 inline-block">Common Name:</span> <span className="text-white">{selectedCert.issuer}</span></p>
                  <p className="text-sm"><span className="text-slate-500 w-24 inline-block">Organization:</span> <span className="text-white">{selectedCert.issuer} Inc.</span></p>
                </div>
              </div>
            </div>

            <div className="border-t border-slate-700/50 pt-4">
              <h4 className="text-sm font-semibold text-slate-400 uppercase tracking-wider mb-3">Validity Period</h4>
              <div className="grid grid-cols-2 gap-4">
                <div className="p-3 bg-slate-900/30 rounded-lg border border-slate-700/30">
                  <p className="text-xs text-slate-500">Not Before</p>
                  <p className="text-sm font-mono text-white">{format(new Date(selectedCert.not_before), 'MMM dd, yyyy HH:mm:ss')}</p>
                </div>
                <div className="p-3 bg-slate-900/30 rounded-lg border border-slate-700/30">
                  <p className="text-xs text-slate-500">Not After</p>
                  <p className="text-sm font-mono text-white">{format(new Date(selectedCert.not_after), 'MMM dd, yyyy HH:mm:ss')}</p>
                </div>
              </div>
            </div>

            <div className="flex justify-end pt-4">
              <button
                onClick={() => setIsViewModalOpen(false)}
                className="btn-secondary"
              >
                Close
              </button>
            </div>
          </div>
        )}
      </Modal>
    </div>
  );
}