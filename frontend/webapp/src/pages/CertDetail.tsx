import { useParams, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from 'react-query';
import { certsAPI } from '../lib/api-client';
import toast from 'react-hot-toast';
import { format } from 'date-fns';
import { ArrowLeftIcon, ArrowDownTrayIcon, TrashIcon, ArrowPathIcon } from '@heroicons/react/24/outline';

export default function CertDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const queryClient = useQueryClient();

  const { data: cert, isLoading, error } = useQuery(
    ['cert', id],
    () => certsAPI.get(id!),
    { enabled: !!id }
  );

  const revokeMutation = useMutation(
    (reason: string) => certsAPI.revoke(id!, reason),
    {
      onSuccess: () => {
        toast.success('Certificate revoked successfully');
        queryClient.invalidateQueries(['cert', id]);
        queryClient.invalidateQueries(['certs']);
      },
      onError: (error: any) => {
        toast.error(error.response?.data?.error || 'Failed to revoke certificate');
      },
    }
  );

  const handleDownload = () => {
    if (cert?.cert_pem) {
      const blob = new Blob([cert.cert_pem], { type: 'application/x-pem-file' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `${cert.subject.replace(/\s+/g, '_')}.pem`;
      a.click();
      URL.revokeObjectURL(url);
      toast.success('Certificate downloaded');
    }
  };

  if (isLoading) {
    return <div className="text-center py-12"><div className="spinner w-8 h-8 mx-auto"></div></div>;
  }

  if (error || !cert) {
    return (
      <div className="text-center py-12">
        <p className="text-red-400">Failed to load certificate</p>
        <button onClick={() => navigate('/inventory')} className="mt-4 text-blue-400 hover:text-blue-300">
          Back to Inventory
        </button>
      </div>
    );
  }

  const daysUntilExpiry = cert.days_until_expiry || Math.ceil(
    (new Date(cert.not_after).getTime() - new Date().getTime()) / (1000 * 60 * 60 * 24)
  );

  return (
    <div>
      <button
        onClick={() => navigate('/inventory')}
        className="mb-6 flex items-center text-gray-400 hover:text-white"
      >
        <ArrowLeftIcon className="h-5 w-5 mr-2" />
        Back to Inventory
      </button>

      <div className="bg-gray-800 rounded-lg shadow-lg p-6 mb-6">
        <div className="flex justify-between items-start mb-6">
          <div>
            <h1 className="text-2xl font-bold text-white mb-2">{cert.subject}</h1>
            <p className="text-gray-400">Certificate Details</p>
          </div>
          <div className="flex gap-2">
            <button
              onClick={handleDownload}
              className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md flex items-center"
            >
              <ArrowDownTrayIcon className="h-5 w-5 mr-2" />
              Download PEM
            </button>
            <button
              onClick={() => navigate(`/request?renew=${id}`)}
              className="px-4 py-2 bg-gray-700 hover:bg-gray-600 text-white rounded-md flex items-center"
            >
              <ArrowPathIcon className="h-5 w-5 mr-2" />
              Rotate
            </button>
            <button
              onClick={() => {
                if (confirm('Are you sure you want to revoke this certificate?')) {
                  revokeMutation.mutate('unspecified');
                }
              }}
              disabled={revokeMutation.isLoading || cert.status === 'revoked'}
              className="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-md flex items-center disabled:opacity-50"
            >
              <TrashIcon className="h-5 w-5 mr-2" />
              Revoke
            </button>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label className="text-sm text-gray-400">Status</label>
            <div className="mt-1">
              <span className={`px-3 py-1 rounded-full text-sm font-medium ${
                cert.status === 'active' ? 'bg-green-900 text-green-300' :
                cert.status === 'expired' ? 'bg-red-900 text-red-300' :
                cert.status === 'revoked' ? 'bg-gray-700 text-gray-300' :
                'bg-yellow-900 text-yellow-300'
              }`}>
                {cert.status}
              </span>
            </div>
          </div>

          <div>
            <label className="text-sm text-gray-400">Days Until Expiry</label>
            <p className="mt-1 text-white font-medium">
              {daysUntilExpiry > 0 ? `${daysUntilExpiry} days` : 'Expired'}
            </p>
          </div>

          <div>
            <label className="text-sm text-gray-400">Issuer</label>
            <p className="mt-1 text-white">{cert.issuer}</p>
          </div>

          <div>
            <label className="text-sm text-gray-400">Fingerprint</label>
            <p className="mt-1 text-white font-mono text-sm">{cert.fingerprint}</p>
          </div>

          <div>
            <label className="text-sm text-gray-400">Valid From</label>
            <p className="mt-1 text-white">{format(new Date(cert.not_before), 'PPpp')}</p>
          </div>

          <div>
            <label className="text-sm text-gray-400">Valid To</label>
            <p className="mt-1 text-white">{format(new Date(cert.not_after), 'PPpp')}</p>
          </div>

          <div>
            <label className="text-sm text-gray-400">Key Algorithm</label>
            <p className="mt-1 text-white">{cert.key_algo} ({cert.key_size} bits)</p>
          </div>

          <div>
            <label className="text-sm text-gray-400">Owner</label>
            <p className="mt-1 text-white">{cert.owner_id}</p>
          </div>
        </div>

        {cert.sans && cert.sans.length > 0 && (
          <div className="mt-6">
            <label className="text-sm text-gray-400">Subject Alternative Names</label>
            <div className="mt-1 flex flex-wrap gap-2">
              {cert.sans.map((san: string, idx: number) => (
                <span key={idx} className="px-3 py-1 bg-gray-700 rounded-md text-sm text-white">
                  {san}
                </span>
              ))}
            </div>
          </div>
        )}

        {cert.audit_logs && cert.audit_logs.length > 0 && (
          <div className="mt-8">
            <h2 className="text-lg font-semibold text-white mb-4">Audit Log</h2>
            <div className="space-y-2">
              {cert.audit_logs.map((log: any) => (
                <div key={log.id} className="bg-gray-700 rounded p-3 text-sm">
                  <div className="flex justify-between">
                    <span className="text-white">{log.action}</span>
                    <span className="text-gray-400">{format(new Date(log.timestamp), 'PPpp')}</span>
                  </div>
                  <div className="text-gray-400 mt-1">by {log.performed_by}</div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
