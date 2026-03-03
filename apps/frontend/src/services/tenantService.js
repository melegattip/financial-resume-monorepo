import apiClient from './apiClient';

const tenantService = {
  // ─── Tenant ───────────────────────────────────────────────────────────────

  async getMyTenant() {
    const response = await apiClient.get('/tenants/me');
    return response.data.tenant;
  },

  async listMyTenants() {
    const response = await apiClient.get('/tenants/list');
    return response.data.tenants;
  },

  async switchTenant(tenantId) {
    const response = await apiClient.post('/users/switch-tenant', { tenant_id: tenantId });
    return response.data; // { access_token, refresh_token, expires_at }
  },

  async updateMyTenant(data) {
    const response = await apiClient.put('/tenants/me', data);
    return response.data;
  },

  async deleteMyTenant() {
    const response = await apiClient.delete('/tenants/me');
    return response.data;
  },

  // ─── Permissions ──────────────────────────────────────────────────────────

  async getMyPermissions() {
    const response = await apiClient.get('/tenants/me/permissions');
    return response.data.permissions;
  },

  // ─── Members ──────────────────────────────────────────────────────────────

  async listMembers() {
    const response = await apiClient.get('/tenants/me/members');
    return response.data.members;
  },

  async updateMemberRole(userID, role) {
    const response = await apiClient.put(`/tenants/me/members/${userID}/role`, { role });
    return response.data;
  },

  async removeMember(userID) {
    const response = await apiClient.delete(`/tenants/me/members/${userID}`);
    return response.data;
  },

  // ─── Invitations ──────────────────────────────────────────────────────────

  async listInvitations() {
    const response = await apiClient.get('/tenants/me/invitations');
    return response.data.invitations;
  },

  async createInvitation(data) {
    const response = await apiClient.post('/tenants/me/invitations', data);
    return response.data.invitation;
  },

  async revokeInvitation(code) {
    const response = await apiClient.delete(`/tenants/me/invitations/${code}`);
    return response.data;
  },

  async joinTenant(code) {
    const response = await apiClient.post('/tenants/join', { code });
    return response.data;
  },

  // ─── Audit Logs ───────────────────────────────────────────────────────────

  async getAuditLogs(limit = 50, offset = 0) {
    const response = await apiClient.get('/tenants/me/audit', {
      params: { limit, offset },
    });
    return response.data;
  },
};

export default tenantService;
