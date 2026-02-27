import React, { useState, useEffect, useCallback } from 'react';
import { FaHistory, FaChevronLeft, FaChevronRight } from 'react-icons/fa';
import tenantService from '../services/tenantService';
import toast from '../utils/notifications';

const PAGE_SIZE = 50;

const actionBadgeClass = (action) => {
  if (action?.includes('deleted')) return 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-400';
  if (action?.includes('created')) return 'bg-green-100 text-green-700 dark:bg-green-900/30 dark:text-green-400';
  if (action?.includes('updated')) return 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400';
  return 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300';
};

const AuditLogs = () => {
  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(true);
  const [offset, setOffset] = useState(0);
  const [hasMore, setHasMore] = useState(false);

  const load = useCallback(async (currentOffset) => {
    try {
      setLoading(true);
      const data = await tenantService.getAuditLogs(PAGE_SIZE, currentOffset);
      const fetchedLogs = data.logs || [];
      setLogs(fetchedLogs);
      setHasMore(fetchedLogs.length === PAGE_SIZE);
    } catch (err) {
      toast.error('Error cargando audit logs');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load(offset);
  }, [offset, load]);

  const handlePrev = () => setOffset(Math.max(0, offset - PAGE_SIZE));
  const handleNext = () => setOffset(offset + PAGE_SIZE);

  const currentPage = Math.floor(offset / PAGE_SIZE) + 1;

  return (
    <div className="max-w-6xl mx-auto px-4 py-6 space-y-6">
      {/* Header */}
      <div className="flex items-center gap-3">
        <div className="p-2 bg-blue-50 dark:bg-blue-900/30 rounded-lg">
          <FaHistory className="w-5 h-5 text-blue-600 dark:text-blue-400" />
        </div>
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-gray-100">Registro de actividad</h1>
          <p className="text-sm text-gray-500 dark:text-gray-400">Historial de eventos del espacio</p>
        </div>
      </div>

      {/* Table */}
      <div className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden">
        {loading ? (
          <div className="text-center py-12 text-gray-500 dark:text-gray-400">Cargando registros…</div>
        ) : logs.length === 0 ? (
          <div className="text-center py-12 text-gray-500 dark:text-gray-400">
            No hay registros de actividad.
          </div>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-700/50">
                <th className="text-left px-4 py-3 text-gray-600 dark:text-gray-400 font-medium">Fecha</th>
                <th className="text-left px-4 py-3 text-gray-600 dark:text-gray-400 font-medium">Acción</th>
                <th className="text-left px-4 py-3 text-gray-600 dark:text-gray-400 font-medium">Usuario</th>
                <th className="text-left px-4 py-3 text-gray-600 dark:text-gray-400 font-medium">Entidad</th>
                <th className="text-left px-4 py-3 text-gray-600 dark:text-gray-400 font-medium">ID</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
              {logs.map((log) => (
                <tr key={log.id} className="hover:bg-gray-50 dark:hover:bg-gray-700/30">
                  <td className="px-4 py-3 text-gray-500 dark:text-gray-400 whitespace-nowrap">
                    {log.created_at
                      ? new Date(log.created_at).toLocaleString('es-AR', {
                          dateStyle: 'short',
                          timeStyle: 'short',
                        })
                      : '—'}
                  </td>
                  <td className="px-4 py-3">
                    <span className={`inline-flex items-center px-2 py-0.5 rounded text-xs font-medium ${actionBadgeClass(log.action)}`}>
                      {log.action}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-gray-700 dark:text-gray-300 font-mono text-xs">
                    {log.user_id ? log.user_id.slice(0, 8) + '…' : '—'}
                  </td>
                  <td className="px-4 py-3 text-gray-700 dark:text-gray-300 capitalize">
                    {log.entity_type || '—'}
                  </td>
                  <td className="px-4 py-3 text-gray-500 dark:text-gray-400 font-mono text-xs">
                    {log.entity_id ? log.entity_id.slice(0, 12) + '…' : '—'}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>

      {/* Pagination */}
      {!loading && (logs.length > 0 || offset > 0) && (
        <div className="flex items-center justify-between">
          <span className="text-sm text-gray-500 dark:text-gray-400">
            Página {currentPage}
          </span>
          <div className="flex gap-2">
            <button
              onClick={handlePrev}
              disabled={offset === 0}
              className="flex items-center gap-1.5 px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded-lg text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
            >
              <FaChevronLeft className="w-3 h-3" />
              Anterior
            </button>
            <button
              onClick={handleNext}
              disabled={!hasMore}
              className="flex items-center gap-1.5 px-3 py-1.5 text-sm border border-gray-300 dark:border-gray-600 rounded-lg text-gray-700 dark:text-gray-300 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
            >
              Siguiente
              <FaChevronRight className="w-3 h-3" />
            </button>
          </div>
        </div>
      )}
    </div>
  );
};

export default AuditLogs;
