import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { FaBuilding, FaUsers, FaEnvelope, FaExclamationTriangle, FaCopy, FaCheck, FaTrash, FaUserCog, FaPlus, FaKey, FaSignInAlt, FaExchangeAlt, FaCheckCircle } from 'react-icons/fa';
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

// ─── My Spaces Section ────────────────────────────────────────────────────────

function MySpacesSection() {
  const { currentTenant, availableTenants, switching, switchTenant } = useTenant();

  if (availableTenants.length <= 1) return null;

  const roleBadge = (role) => {
    const map = {
      owner: 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400',
      admin: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400',
      member: 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400',
      viewer: 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300',
    };
    return map[role] || map.viewer;
  };

  const handleSwitch = async (tenantId) => {
    try {
      await switchTenant(tenantId);
      toast.success('Espacio activado');
      window.location.href = '/dashboard';
    } catch (err) {
      toast.error(err?.response?.data?.error || 'Error al cambiar de espacio');
    }
  };

  return (
    <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden">
      <div className="px-6 py-4 border-b border-gray-100 dark:border-gray-700 flex items-center gap-3">
        <div className="p-2 bg-blue-50 dark:bg-blue-900/30 rounded-lg flex-shrink-0">
          <FaExchangeAlt className="w-4 h-4 text-blue-600 dark:text-blue-400" />
        </div>
        <div>
          <h3 className="text-base font-semibold text-gray-900 dark:text-gray-100">Mis espacios</h3>
          <p className="text-sm text-gray-500 dark:text-gray-400">Cambiá de espacio sin cerrar sesión.</p>
        </div>
      </div>
      <ul className="divide-y divide-gray-100 dark:divide-gray-700">
        {availableTenants.map(t => {
          const isActive = t.id === currentTenant?.id;
          return (
            <li key={t.id} className={`flex items-center justify-between px-6 py-4 ${isActive ? 'bg-blue-50/60 dark:bg-blue-900/10' : 'hover:bg-gray-50 dark:hover:bg-gray-700/30'}`}>
              <div className="flex items-center gap-3 min-w-0">
                {isActive
                  ? <FaCheckCircle className="w-4 h-4 text-blue-500 flex-shrink-0" />
                  : <div className="w-4 h-4 rounded-full border-2 border-gray-300 dark:border-gray-600 flex-shrink-0" />
                }
                <div className="min-w-0">
                  <p className="text-sm font-medium text-gray-900 dark:text-gray-100 truncate">{t.name}</p>
                  <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium mt-0.5 ${roleBadge(t.role)}`}>
                    {t.role}
                  </span>
                </div>
              </div>
              {isActive ? (
                <span className="text-xs text-blue-600 dark:text-blue-400 font-medium flex-shrink-0">Activo</span>
              ) : (
                <button
                  onClick={() => handleSwitch(t.id)}
                  disabled={switching}
                  className="flex-shrink-0 px-3 py-1.5 text-xs font-medium bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                >
                  {switching ? 'Cambiando…' : 'Activar'}
                </button>
              )}
            </li>
          );
        })}
      </ul>
    </div>
  );
}

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
      <MySpacesSection />
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
      <JoinTenantSection />
    </div>
  );
}

// Helper: relative expiry label
function expiryLabel(expiresAt) {
  if (!expiresAt) return <span className="text-amber-500 dark:text-amber-400">Sin expiración</span>;
  const diff = new Date(expiresAt) - Date.now();
  if (diff <= 0) return <span className="text-red-500 dark:text-red-400">Expirada</span>;
  const hours = Math.floor(diff / 3600000);
  if (hours < 1) return <span className="text-orange-500 dark:text-orange-400">&lt;1h</span>;
  if (hours < 24) return <span className="text-amber-600 dark:text-amber-400">Expira en {hours}h</span>;
  const days = Math.floor(hours / 24);
  return <span className="text-gray-500 dark:text-gray-400">Expira en {days}d</span>;
}

// ─── Members + Invitations (unified card) ─────────────────────────────────────

function MembersTab({ isAdmin, currentUserID, showInvitations }) {
  const [members, setMembers] = useState([]);
  const [invitations, setInvitations] = useState([]);
  const [loading, setLoading] = useState(true);
  const [creating, setCreating] = useState(false);
  const [newRole, setNewRole] = useState('member');
  const [copiedCode, setCopiedCode] = useState(null);

  const load = async () => {
    try {
      setLoading(true);
      const [membersData, invData] = await Promise.all([
        tenantService.listMembers(),
        showInvitations ? tenantService.listInvitations() : Promise.resolve([]),
      ]);
      setMembers(membersData || []);
      setInvitations(invData || []);
    } catch (err) {
      toast.error('Error cargando datos');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => { load(); }, []); // eslint-disable-line react-hooks/exhaustive-deps

  const adminCount = members.filter(m => m.role === 'admin').length;

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

  const handleCreate = async () => {
    try {
      setCreating(true);
      await tenantService.createInvitation({ role: newRole });
      toast.success('Invitación creada (válida 24h)');
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

  if (loading) {
    return <div className="text-center py-8 text-gray-500 dark:text-gray-400">Cargando…</div>;
  }

  return (
    <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden">
      {/* ── Miembros ── */}
      <div className="px-4 py-3 border-b border-gray-100 dark:border-gray-700 bg-gray-50 dark:bg-gray-700/50">
        <div className="flex items-center gap-2">
          <FaUsers className="w-4 h-4 text-gray-500 dark:text-gray-400" />
          <span className="text-sm font-semibold text-gray-900 dark:text-gray-100">Miembros</span>
        </div>
      </div>
      <table className="w-full text-sm">
        <thead>
          <tr className="border-b border-gray-200 dark:border-gray-700 bg-gray-50/50 dark:bg-gray-700/30">
            <th className="text-left px-4 py-2.5 text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide">Usuario</th>
            <th className="text-left px-4 py-2.5 text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide">Rol</th>
            <th className="text-left px-4 py-2.5 text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide">Desde</th>
            {isAdmin && <th className="px-4 py-2.5" />}
          </tr>
        </thead>
        <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
          {members.map((m) => {
            const isLastAdmin = m.role === 'admin' && adminCount <= 1;
            const canEditRole = isAdmin && m.role !== 'owner' && m.user_id !== currentUserID;
            return (
              <tr key={m.user_id || m.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/30">
                <td className="px-4 py-3">
                  <div className="font-medium text-gray-900 dark:text-gray-100">
                    {m.user_name || m.user_email || m.user_id}
                  </div>
                  {m.user_name && m.user_email && (
                    <div className="text-xs text-gray-500 dark:text-gray-400">{m.user_email}</div>
                  )}
                </td>
                <td className="px-4 py-3">
                  {canEditRole && !isLastAdmin ? (
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
                    <span
                      className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${roleBadgeClass(m.role)}`}
                      title={isLastAdmin ? 'Único administrador — asigna otro antes de cambiar este rol' : undefined}
                    >
                      {m.role}
                    </span>
                  )}
                </td>
                <td className="px-4 py-3 text-gray-500 dark:text-gray-400 text-xs">
                  {m.joined_at ? new Date(m.joined_at).toLocaleDateString('es-AR') : '—'}
                </td>
                {isAdmin && (
                  <td className="px-4 py-3 text-right">
                    {m.role !== 'owner' && m.user_id !== currentUserID && !isLastAdmin && (
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
            );
          })}
          {members.length === 0 && (
            <tr>
              <td colSpan={isAdmin ? 4 : 3} className="px-4 py-6 text-center text-gray-500 dark:text-gray-400 text-sm">
                No hay miembros.
              </td>
            </tr>
          )}
        </tbody>
      </table>

      {/* ── Invitaciones ── */}
      {showInvitations && (
        <>
          <div className="flex items-center justify-between px-4 py-3 border-t border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-700/50">
            <div className="flex items-center gap-2">
              <FaEnvelope className="w-4 h-4 text-gray-500 dark:text-gray-400" />
              <span className="text-sm font-semibold text-gray-900 dark:text-gray-100">Invitaciones activas</span>
              <span className="text-xs text-gray-400 dark:text-gray-500">(válidas 24h)</span>
            </div>
            <div className="flex items-center gap-2">
              <select
                value={newRole}
                onChange={(e) => setNewRole(e.target.value)}
                className="px-2 py-1.5 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 text-xs"
              >
                {ROLES.filter(r => r !== 'owner').map(r => (
                  <option key={r} value={r}>{r}</option>
                ))}
              </select>
              <button
                onClick={handleCreate}
                disabled={creating}
                className="flex items-center gap-1.5 px-3 py-1.5 bg-blue-600 text-white text-xs rounded-lg hover:bg-blue-700 disabled:opacity-50 transition-colors"
              >
                <FaPlus className="w-3 h-3" />
                {creating ? 'Creando…' : 'Nuevo código'}
              </button>
            </div>
          </div>
          {invitations.length === 0 ? (
            <div className="px-4 py-5 text-center text-gray-500 dark:text-gray-400 text-sm border-t border-gray-100 dark:border-gray-700">
              No hay invitaciones activas.
            </div>
          ) : (
            <table className="w-full text-sm border-t border-gray-100 dark:border-gray-700">
              <thead>
                <tr className="border-b border-gray-200 dark:border-gray-700 bg-gray-50/50 dark:bg-gray-700/30">
                  <th className="text-left px-4 py-2.5 text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide">Código</th>
                  <th className="text-left px-4 py-2.5 text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide">Rol</th>
                  <th className="text-left px-4 py-2.5 text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wide">Vence</th>
                  <th className="px-4 py-2.5" />
                </tr>
              </thead>
              <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                {invitations.map((inv) => (
                  <tr key={inv.code} className="hover:bg-gray-50 dark:hover:bg-gray-700/30">
                    <td className="px-4 py-3">
                      <div className="flex items-center gap-2">
                        <span className="font-mono text-sm text-gray-700 dark:text-gray-300">{inv.code}</span>
                        <button
                          onClick={() => handleCopy(inv.code)}
                          className="p-1 text-gray-400 hover:text-blue-500 transition-colors"
                          title="Copiar código"
                        >
                          {copiedCode === inv.code
                            ? <FaCheck className="w-3 h-3 text-green-500" />
                            : <FaCopy className="w-3 h-3" />
                          }
                        </button>
                      </div>
                    </td>
                    <td className="px-4 py-3">
                      <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${roleBadgeClass(inv.role)}`}>
                        {inv.role}
                      </span>
                    </td>
                    <td className="px-4 py-3 text-xs">
                      {expiryLabel(inv.expires_at)}
                    </td>
                    <td className="px-4 py-3 text-right">
                      <button
                        onClick={() => handleRevoke(inv.code)}
                        className="p-1.5 text-red-400 hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-900/20 rounded transition-colors"
                        title="Revocar invitación"
                      >
                        <FaTrash className="w-3 h-3" />
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </>
      )}
    </div>
  );
}

// ─── InvitationsTab (standalone, used in standalone mode) ─────────────────────

function InvitationsTab() {
  return <MembersTab isAdmin showInvitations />;
}

// ─── Danger Zone Tab ──────────────────────────────────────────────────────────

export function DangerZoneTab({ tenantName }) {
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

// ─── Join Tenant Section ──────────────────────────────────────────────────────

function JoinTenantSection() {
  const { switchTenant, refreshTenant } = useTenant();
  const [code, setCode] = useState('');
  const [joining, setJoining] = useState(false);

  const handleJoin = async (e) => {
    e.preventDefault();
    if (!code.trim()) return;
    try {
      setJoining(true);
      const result = await tenantService.joinTenant(code.trim());
      const tenantName = result.tenant?.name || 'el nuevo espacio';
      const tenantId = result.tenant?.id;
      setCode('');
      // Auto-switch to the newly joined tenant
      if (tenantId) {
        await switchTenant(tenantId);
        toast.success(`¡Te uniste a "${tenantName}"! Espacio activado.`);
        window.location.href = '/dashboard';
      } else {
        await refreshTenant();
        toast.success(`¡Te uniste a "${tenantName}"!`);
      }
    } catch (err) {
      toast.error(err.response?.data?.error || 'Código inválido o expirado');
    } finally {
      setJoining(false);
    }
  };

  return (
    <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-6">
      <div className="flex items-start gap-3 mb-4">
        <div className="p-2 bg-green-50 dark:bg-green-900/30 rounded-lg flex-shrink-0">
          <FaKey className="w-4 h-4 text-green-600 dark:text-green-400" />
        </div>
        <div>
          <h3 className="text-base font-semibold text-gray-900 dark:text-gray-100">Unirse a otro espacio</h3>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
            Ingresá un código de invitación para unirte y activar el espacio automáticamente.
          </p>
        </div>
      </div>
      <form onSubmit={handleJoin} className="flex gap-3">
        <input
          type="text"
          value={code}
          onChange={(e) => setCode(e.target.value)}
          placeholder="Código de invitación"
          className="flex-1 px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-gray-100 font-mono text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <button
          type="submit"
          disabled={joining || !code.trim()}
          className="flex items-center gap-2 px-4 py-2 bg-green-600 text-white text-sm rounded-lg hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          <FaSignInAlt className="w-3.5 h-3.5" />
          {joining ? 'Uniéndose…' : 'Unirse'}
        </button>
      </form>
    </div>
  );
}

// ─── Main Page ────────────────────────────────────────────────────────────────

const TABS = [
  { id: 'general', label: 'General', icon: FaBuilding },
  { id: 'members', label: 'Miembros', icon: FaUsers },
];

const TenantSettings = ({ widgetMode = false }) => {
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

  // Widget mode: todas las secciones en la misma pantalla sin sub-tabs
  if (widgetMode) {
    return (
      <div className="space-y-4">
        <GeneralTab tenant={currentTenant} onRefresh={refreshTenant} isAdmin={isAdmin} />
        <MembersTab isAdmin={isAdmin} currentUserID={currentUserID} showInvitations={hasPermission('invite_members')} />
      </div>
    );
  }

  // Standalone mode: mantiene sub-tabs propios
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
          <MembersTab isAdmin={isAdmin} currentUserID={currentUserID} showInvitations={hasPermission('invite_members')} />
        )}
      </div>
    </div>
  );
};

export default TenantSettings;
