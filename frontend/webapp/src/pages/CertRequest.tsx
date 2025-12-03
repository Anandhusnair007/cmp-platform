import React, { useState } from 'react';
import { useMutation, useQueryClient } from 'react-query';
import { useNavigate } from 'react-router-dom';
import { certsAPI } from '../lib/api-client';
import toast from 'react-hot-toast';

export default function CertRequest() {
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [formData, setFormData] = useState({
    owner_id: 'default-owner',
    common_name: '',
    sans: '',
    key_algorithm: 'rsa' as 'rsa' | 'ecdsa',
    key_size: 2048,
    adapter_id: 'vault-staging',
    agent_id: '',
    path: '',
    reload_cmd: '',
  });

  const requestMutation = useMutation(
    (data: any) => certsAPI.request(data),
    {
      onSuccess: (response) => {
        toast.success(`Certificate request submitted! Request ID: ${response.request_id}`);
        queryClient.invalidateQueries(['certs']);
        queryClient.invalidateQueries(['inventory']);
        
        // Reset form
        setFormData({
          ...formData,
          common_name: '',
          sans: '',
        });
        
        // Navigate to inventory after a delay
        setTimeout(() => {
          navigate('/inventory');
        }, 2000);
      },
      onError: (error: any) => {
        toast.error(error.response?.data?.error || 'Failed to submit certificate request');
      },
    }
  );

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    const sans = formData.sans
      .split(',')
      .map((s) => s.trim())
      .filter((s) => s.length > 0);

    const requestData = {
      owner_id: formData.owner_id,
      common_name: formData.common_name,
      sans: sans.length > 0 ? sans : undefined,
      key_algorithm: formData.key_algorithm,
      key_size: formData.key_size,
      adapter_id: formData.adapter_id,
      install_targets:
        formData.agent_id && formData.path
          ? [
              {
                agent_id: formData.agent_id,
                path: formData.path,
                reload_cmd: formData.reload_cmd || undefined,
              },
            ]
          : undefined,
    };

    requestMutation.mutate(requestData);
  };

  return (
    <div>
      <h1 className="text-3xl font-bold text-white mb-6">Request Certificate</h1>

      <div className="bg-gray-800 rounded-lg shadow-lg">
        <div className="px-6 py-6">
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label htmlFor="common_name" className="block text-sm font-medium text-gray-300 mb-2">
                Common Name *
              </label>
              <input
                type="text"
                id="common_name"
                required
                value={formData.common_name}
                onChange={(e) => setFormData({ ...formData, common_name: e.target.value })}
                className="mt-1 block w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="app.example.com"
              />
            </div>

            <div>
              <label htmlFor="sans" className="block text-sm font-medium text-gray-300 mb-2">
                Subject Alternative Names (comma-separated)
              </label>
              <input
                type="text"
                id="sans"
                value={formData.sans}
                onChange={(e) => setFormData({ ...formData, sans: e.target.value })}
                className="mt-1 block w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="www.example.com, api.example.com"
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label htmlFor="key_algorithm" className="block text-sm font-medium text-gray-300 mb-2">
                  Key Algorithm
                </label>
                <select
                  id="key_algorithm"
                  value={formData.key_algorithm}
                  onChange={(e) => setFormData({ ...formData, key_algorithm: e.target.value as 'rsa' | 'ecdsa' })}
                  className="mt-1 block w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-md text-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                >
                  <option value="rsa">RSA</option>
                  <option value="ecdsa">ECDSA</option>
                </select>
              </div>

              <div>
                <label htmlFor="key_size" className="block text-sm font-medium text-gray-300 mb-2">
                  Key Size (bits)
                </label>
                <input
                  type="number"
                  id="key_size"
                  value={formData.key_size}
                  onChange={(e) => setFormData({ ...formData, key_size: parseInt(e.target.value) })}
                  className="mt-1 block w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                />
              </div>
            </div>

            <div>
              <label htmlFor="adapter_id" className="block text-sm font-medium text-gray-300 mb-2">
                CA Adapter
              </label>
              <input
                type="text"
                id="adapter_id"
                value={formData.adapter_id}
                onChange={(e) => setFormData({ ...formData, adapter_id: e.target.value })}
                className="mt-1 block w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                placeholder="vault-staging"
              />
            </div>

            <div className="border-t border-gray-700 pt-6">
              <h3 className="text-lg font-medium text-white mb-4">Installation (Optional)</h3>
              
              <div className="space-y-4">
                <div>
                  <label htmlFor="agent_id" className="block text-sm font-medium text-gray-300 mb-2">
                    Agent ID
                  </label>
                  <input
                    type="text"
                    id="agent_id"
                    value={formData.agent_id}
                    onChange={(e) => setFormData({ ...formData, agent_id: e.target.value })}
                    className="mt-1 block w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    placeholder="agent-1"
                  />
                </div>

                <div>
                  <label htmlFor="path" className="block text-sm font-medium text-gray-300 mb-2">
                    Certificate Path
                  </label>
                  <input
                    type="text"
                    id="path"
                    value={formData.path}
                    onChange={(e) => setFormData({ ...formData, path: e.target.value })}
                    className="mt-1 block w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    placeholder="/etc/ssl/certs/app.pem"
                  />
                </div>

                <div>
                  <label htmlFor="reload_cmd" className="block text-sm font-medium text-gray-300 mb-2">
                    Reload Command
                  </label>
                  <input
                    type="text"
                    id="reload_cmd"
                    value={formData.reload_cmd}
                    onChange={(e) => setFormData({ ...formData, reload_cmd: e.target.value })}
                    className="mt-1 block w-full px-4 py-2 bg-gray-700 border border-gray-600 rounded-md text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    placeholder="systemctl reload nginx"
                  />
                </div>
              </div>
            </div>

            <div>
              <button
                type="submit"
                disabled={requestMutation.isLoading}
                className="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {requestMutation.isLoading ? 'Submitting...' : 'Submit Request'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}