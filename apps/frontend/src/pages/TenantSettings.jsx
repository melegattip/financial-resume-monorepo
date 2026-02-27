import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { FaBuilding, FaUsers, FaEnvelope, FaExclamationTriangle, FaCopy, FaCheck, FaTrash, FaUserCog, FaPlus } from 'react-icons/fa';
import { useTenant } from '../contexts/TenantContext';
import RoleGuard from '../components/RoleGuard';
import tenantService from '../services/tenantService';
import toast from '../utils/notifications';

const ROLES = ['owner', 'admin', 'member', 'viewer'];

const roleBadgeClass = (role) => {
  const map = {
    owner: 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400',
    admin: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400',
    member: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
    viewer: 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300',
  };
  return map[role] || map.viewer;
};

// ─── General Tab ──────────────────────────────────────────────────────────────

function GeneralTab({ tenant, onRefresh, isAdmin }) {
  const [name, setName] = useState(tenant?.name || '');
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    setName(tenant?.name || '');
  }, [tenant]);

  const handleSave = async (e) => {
    e.preventDefault();
    if (!name.trim()) return;
    try {
      setSaving(true);
      await tenantService.updateMyTenant({ name: name.trim() });
      toast.success('Nombre actualizado');
      onRefresh();
    } catch (err) {
      toast.error(err.response?.data?.error || 'Error actualizando tenant');
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="space-y-6">
      <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6">
        <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">Información del espacio</h3>
        <dl className="space-y-3 text-sm">
          <div className="flex justify-between">
            <dt className="text-gray-500 dark:text-gray-400">ID</dt>
            <dd className="font-mono text-gray-700 dark:text-gray-300">{tenant?.id}</dd>
          </div>
          <div className="flex justify-between">
            <dt className="text-gray-500 dark:text-gray-400">Slug</dt>
            <dd className="font-mono text-gray-700 dark:text-gray-300">{tenant?.slug}</dd>
          </div>
          <div className="flex justify-between">
            <dt className="text-gray-500 dark:text-gray-400">Plan</dt>
            <dd className="capitalize text-gray-700 dark:text-gray-300">{tenant?.plan || 'free'}</dd>
          </div>
        </dl>
      </div>

      {isAdmin && (
        <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-gray-100 mb-4">Editar nombre</h3>
          <form onSubmit={handleSave} className="flex gap-3">
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Nombre del espacio"
              required
            />
            <button
              type="submit"
              disabled={saving || !name.trim()}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {saving ? 'Guardando…' : 'Guardar'}
            </button>
          </form>
        </div>
      )}
    </div>
  );
}

// ─── Members Tab ──────────────────────────────────────────────────────────────

function MembersTab({ isAdmin, currentUserID }) {
  const [members, setMembers] = useState([]);
  const [loading, setLoading] = useState(true);

  const load = async () => {
    try {
      setLoading(true);
      const data = await tenantService.listMembers();
      setMembers(data || []);
    } catch (err) {
      toast.error('Error cargando miembros');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, []);

  const handleRoleChange = async (userID, newRole) => {
    try {
      await tenantService.updateMemberRole(userID, newRole);
      toast.success('Rol actualizado');
      load();
    } catch (err) {
      toast.error(err.response?.data?.error || 'Error actualizando rol');
    }
  };

  const handleRemove = async (userID) => {
    if (!window.confirm('¿Eliminar este miembro del espacio?')) return;
    try {
      await tenantService.removeMember(userID);
      toast.success('Miembro eliminado');
      load();
    } catch (err) {
      toast.error(err.response?.data?.error || 'Error eliminando miembro');
    }
  };

  if (loading) {
    return <div className="text-center py-8 text-gray-500 dark:text-gray-400">Cargando miembros…</div>;
  }

  return (
    <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden">
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-700/50">
            <th className="text-left px-4 py-3 text-gray-600 dark:text-gray-400 font-medium">Usuario</th>
            <th className="text-left px-4 py-3 text-gray-600 dark:text-gray-400 font-medium">Rol</th>
            <th className="text-left px-4 py-3 text-gray-600 dark:text-gray-400 font-medium">Desde</th>
            {isAdmin && <th className="px-4 py-3" />}
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
          {members.map((m) => (
            <tr key={m.user_id || m.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/30">
              <td className="px-4 py-3">
                <div className="font-medium text-gray-900 dark:text-gray-100">{m.user_email || m.user_id}</div>
              </td>
              <td className="px-4 py-3">
                {isAdmin && m.role !== 'owner' && m.user_id !== currentUserID ? (
                  <select
                    value={m.role}
                    onChange={(e) => handleRoleChange(m.user_id, e.target.value)}
                    className="text-xs px-2 py-1 border border-gray-300 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-300"
                  >
                    {ROLES.filter(r => r !== 'owner').map(r => (
                      <option key={r} value={r}>{r}</option>
                    ))}
                  </select>
                ) : (
                  <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${roleBadgeClass(m.role)}`}>
                    {m.role}
                  </span>
                )}
              </td>
              <td className="px-4 py-3 text-gray-500 dark:text-gray-400">
                {m.joined_at ? new Date(m.joined_at).toLocaleDateString() : '—'}
              </td>
              {isAdmin && (
                <td className="px-4 py-3 text-right">
                  {m.role !== 'owner' && m.user_id !== currentUserID && (
                    <button
                      onClick={() => handleRemove(m.user_id)}
                      className="p-1.5 text-red-400 hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition-colors"
                      title="Eliminar miembro"
                    >
                      <FaTrash className="w-3.5 h-3.5" />
                    </button>
                  )}
                </td>
              )}
            </tr>
          ))}
          {members.length === 0 && (
            <tr>
              <td colSpan={isAdmin ? 4 : 3} className="px-4 py-8 text-center text-gray-500 dark:text-gray-400">
                No hay miembros.
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}

// ─── Invitations Tab ──────────────────────────────────────────────────────────

function InvitationsTab() {
  const [invitations, setInvitations] = useState([]);
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);
  const [newRole, setNewRole] = useState('member');
  const [copiedCode, setCopiedCode] = useState(null);

  const load = async () => {
    try {
      setLoading(true);
      const data = await tenantService.listInvitations();
      setInvitations(data || []);
    } catch (err) {
      toast.error('Error cargando invitaciones');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, []);

  const handleCreate = async () => {
    try {
      setCreating(true);
      await tenantService.createInvitation({ role: newRole });
      toast.success('Invitación creada');
      load();
    } catch (err) {
      toast.error(err.response?.data?.error || 'Error creando invitación');
    } finally {
      setCreating(false);
    }
  };

  const handleRevoke = async (code) => {
    if (!window.confirm('¿Revocar esta invitación?')) return;
    try {
      await tenantService.revokeInvitation(code);
      toast.success('Invitación revocada');
      load();
    } catch (err) {
      toast.error(err.response?.data?.error || 'Error revocando invitación');
    }
  };

  const handleCopy = (code) => {
    navigator.clipboard.writeText(code);
    setCopiedCode(code);
    setTimeout(() => setCopiedCode(null), 2000);
  };

  return (
    <div className="space-y-4">
      {/* Create invitation */}
      <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-4">
        <h3 className="text-sm font-semibold text-gray-900 dark:text-gray-100 mb-3">Nueva invitación</h3>
        <div className="flex gap-3 items-center">
          <select
            value={newRole}
            onChange={(e) => setNewRole(e.target.value)}
            className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 text-sm"
          >
            {ROLES.filter(r => r !== 'owner').map(r => (
              <option key={r} value={r}>{r}</option>
            ))}
          </select>
          <button
            onClick={handleCreate}
            disabled={creating}
            className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white text-sm rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors"
          >
            <FaPlus className="w-3.5 h-3.5" />
            {creating ? 'Creando…' : 'Crear código'}
          </button>
        </div>
      </div>

      {/* List */}
      <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden">
        {loading ? (
          <div className="text-center py-8 text-gray-500 dark:text-gray-400">Cargando…</div>
        ) : invitations.length === 0 ? (
          <div className="text-center py-8 text-gray-500 dark:text-gray-400">No hay invitaciones activas.</div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-700/50">
                <th className="text-left px-4 py-3 text-gray-600 dark:text-gray-400 font-medium">Código</th>
                <th className="text-left px-4 py-3 text-gray-600 dark:text-gray-400 font-medium">Rol</th>
                <th className="text-left px-4 py-3 text-gray-600 dark:text-gray-400 font-medium">Expira</th>
                <th className="px-4 py-3" />
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
              {invitations.map((inv) => (
                <tr key={inv.code} className="hover:bg-gray-50 dark:hover:bg-gray-700/30">
                  <td className="px-4 py-3">
                    <div className="flex items-center gap-2">
                      <span className="font-mono text-gray-700 dark:text-gray-300">{inv.code}</span>
                      <button
                        onClick={() => handleCopy(inv.code)}
                        className="p-1 text-gray-400 hover:text-blue-500 transition-colors"
                        title="Copiar código"
                      >
                        {copiedCode === inv.code
                          ? <FaCheck className="w-3.5 h-3.5 text-green-500" />
                          : <FaCopy className="w-3.5 h-3.5" />
                        }
                      </button>
                    </div>
                  </td>
                  <td className="px-4 py-3">
                    <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${roleBadgeClass(inv.role)}`}>
                      {inv.role}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-gray-500 dark:text-gray-400">
                    {inv.expires_at ? new Date(inv.expires_at).toLocaleDateString() : '—'}
                  </td>
                  <td className="px-4 py-3 text-right">
                    <button
                      onClick={() => handleRevoke(inv.code)}
                      className="p-1.5 text-red-400 hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition-colors"
                      title="Revocar invitación"
                    >
                      <FaTrash className="w-3.5 h-3.5" />
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
}

// ─── Danger Zone Tab ──────────────────────────────────────────────────────────

function DangerZoneTab({ tenantName }) {
  const navigate = useNavigate();
  const { refreshTenant } = useTenant();
  const [confirmName, setConfirmName] = useState('');
  const [deleting, setDeleting] = useState(false);

  const handleDelete = async () => {
    if (confirmName !== tenantName) return;
    if (!window.confirm('Esta acción es irreversible. ¿Confirmar eliminación?')) return;
    try {
      setDeleting(true);
      await tenantService.deleteMyTenant();
      toast.success('Espacio eliminado. Redirigiendo…');
      refreshTenant();
      navigate('/dashboard');
    } catch (err) {
      toast.error(err.response?.data?.error || 'Error eliminando espacio');
    } finally {
      setDeleting(false);
    }
  };

  return (
    <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-xl p-6">
      <div className="flex items-start gap-3 mb-4">
        <FaExclamationTriangle className="w-5 h-5 text-red-500 mt-0.5 flex-shrink-0" />
        <div>
          <h3 className="text-base font-semibold text-red-700 dark:text-red-400">Eliminar espacio</h3>
          <p className="text-sm text-red-600 dark:text-red-400 mt-1">
            Esta acción elimina el espacio y todos sus datos. Es irreversible.
          </p>
        </div>
      </div>
      <div className="space-y-3">
        <label className="block text-sm text-red-700 dark:text-red-400">
          Escribe <span className="font-mono font-semibold">{tenantName}</span> para confirmar:
        </label>
        <input
          type="text"
          value={confirmName}
          onChange={(e) => setConfirmName(e.target.value)}
          className="w-full px-3 py-2 border border-red-300 dark:border-red-700 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-red-500"
          placeholder={tenantName}
        />
        <button
          onClick={handleDelete}
          disabled={confirmName !== tenantName || deleting}
          className="px-4 py-2 bg-red-600 text-white text-sm rounded-lg hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          {deleting ? 'Eliminando…' : 'Eliminar espacio definitivamente'}
        </button>
      </div>
    </div>
  );
}

// ─── Main Page ────────────────────────────────────────────────────────────────

const TABS = [
  { id: 'general', label: 'General', icon: FaBuilding },
  { id: 'members', label: 'Miembros', icon: FaUsers },
  { id: 'invitations', label: 'Invitaciones', icon: FaEnvelope, permission: 'invite_members' },
  { id: 'danger', label: 'Peligro', icon: FaExclamationTriangle, role: 'owner' },
];

const TenantSettings = () => {
  const { currentTenant, myRole, isAdmin, hasPermission, loading, refreshTenant } = useTenant();
  const [activeTab, setActiveTab] = useState('general');
  const currentUserID = JSON.parse(localStorage.getItem('auth_user') || '{}')?.id;

  const visibleTabs = TABS.filter(tab => {
    if (tab.role && myRole !== tab.role) return false;
    if (tab.permission && !hasPermission(tab.permission)) return false;
    return true;
  });

  if (loading) {
    return (
      <div className="flex items-center justify-center py-20 text-gray-500 dark:text-gray-400">
        Cargando…
      </div>
    );
  }

  if (!currentTenant) {
    return (
      <div className="text-center py-20 text-gray-500 dark:text-gray-400">
        No se pudo cargar el espacio.
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto px-4 py-6 space-y-6">
      {/* Page header */}
      <div className="flex items-center gap-3">
        <div className="p-2 bg-blue-50 dark:bg-blue-900/30 rounded-lg">
          <FaUserCog className="w-5 h-5 text-blue-600 dark:text-blue-400" />
        </div>
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Configuración del espacio</h1>
          <p className="text-sm text-gray-500 dark:text-gray-400">{currentTenant.name}</p>
        </div>
      </div>

      {/* Tabs */}
      <div className="border-b border-gray-200 dark:border-gray-700">
        <nav className="flex space-x-1 -mb-px">
          {visibleTabs.map(tab => {
            const Icon = tab.icon;
            return (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`flex items-center gap-2 px-4 py-2.5 text-sm font-medium border-b-2 transition-colors ${
                  activeTab === tab.id
                    ? 'border-blue-500 text-blue-600 dark:text-blue-400'
                    : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:border-gray-300 dark:hover:border-gray-600'
                } ${tab.id === 'danger' ? 'text-red-500 dark:text-red-400 hover:text-red-600 hover:border-red-300' : ''}`}
              >
                <Icon className="w-4 h-4" />
                {tab.label}
              </button>
            );
          })}
        </nav>
      </div>

      {/* Tab content */}
      <div>
        {activeTab === 'general' && (
          <GeneralTab tenant={currentTenant} onRefresh={refreshTenant} isAdmin={isAdmin} />
        )}
        {activeTab === 'members' && (
          <MembersTab isAdmin={isAdmin} currentUserID={currentUserID} />
        )}
        {activeTab === 'invitations' && (
          <RoleGuard permission="invite_members">
            <InvitationsTab />
          </RoleGuard>
        )}
        {activeTab === 'danger' && (
          <RoleGuard role="owner">
            <DangerZoneTab tenantName={currentTenant.name} />
          </RoleGuard>
        )}
      </div>
    </div>
  );
};

export default TenantSettings;
