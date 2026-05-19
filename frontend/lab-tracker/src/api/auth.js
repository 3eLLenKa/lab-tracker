import axios from 'axios';

const api = axios.create({
  baseURL: '/api',
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const register = (data) => api.post('/v1/auth/register', data);
export const login    = (data) => api.post('/v1/auth/login', data);

export const listGroups = () => api.get('/v1/groups/list');

export const listLabWorks = (params) => api.get('/v1/labworks', { params });
export const getLabWork = (id) => api.get(`/v1/labworks/${id}`);
export const createLabWork = (data) => api.post('/v1/labworks', data);
export const updateLabWork = (id, data) => api.put(`/v1/labworks/${id}`, data);
export const deleteLabWork = (id) => api.delete(`/v1/labworks/${id}`);
export const listStudentAssignments = () => api.get('/v1/student/assignments');
export const getStudentProgress = () => api.get('/v1/student/progress');
export const listTeacherSubmissions = () => api.get('/v1/teacher/submissions');
export const submitAssignment = (data) => api.post('/v1/submissions/create', data);
export const setGrade = (data) => api.post('/v1/grades/set', data);
export const getAdminStats = () => api.get('/v1/admin/stats');
export const exportAdminCsv = () => api.get('/v1/admin/export.csv', { responseType: 'blob' });

export default api;
