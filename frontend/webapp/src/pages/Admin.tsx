import toast from 'react-hot-toast';

export default function Admin() {
  return (
    <div>
      <h1 className="text-3xl font-bold text-white mb-6">Admin</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="bg-gray-800 rounded-lg shadow-lg p-6">
          <h2 className="text-xl font-semibold text-white mb-4">Adapter Configuration</h2>
          <p className="text-gray-400 mb-4">Configure Certificate Authority adapters</p>
          <button
            onClick={() => toast('Adapter configuration coming soon', { icon: 'ℹ️' })}
            className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md"
          >
            Manage Adapters
          </button>
        </div>

        <div className="bg-gray-800 rounded-lg shadow-lg p-6">
          <h2 className="text-xl font-semibold text-white mb-4">RBAC Management</h2>
          <p className="text-gray-400 mb-4">Configure roles and permissions</p>
          <button
            onClick={() => toast('RBAC management coming soon', { icon: 'ℹ️' })}
            className="px-4 py-2 bg-blue-600 hover:bg-blue-700 text-white rounded-md"
          >
            Manage Roles
          </button>
        </div>
      </div>
    </div>
  );
}
