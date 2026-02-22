import axios from 'axios';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://127.0.0.1:2324';

export { API_BASE_URL };

const api = axios.create({
    baseURL: API_BASE_URL,
    withCredentials: true,
    headers: {
        'Content-Type': 'application/json'
    }
});

api.interceptors.response.use(
    response => response,
    error => {
        if (error.response?.status === 401) {
            localStorage.removeItem('userInfo');
            window.location.href = '/login';
        } else if (error.response?.status === 403) {
            localStorage.removeItem('userInfo');
            window.location.href = '/login';
        }
        return Promise.reject(error);
    }
);

export const authApi = {
    login: (email: string, passwd: string) =>
        api.post('/api/v1/user/login', { email, passwd }),
    
    register: (email: string, passwd: string, name?: string, code?: string) =>
        api.post('/api/v1/user/register', { email, passwd, name, code }),
    
    logout: () =>
        api.post('/api/v1/user/logout'),
    
    getCurrentUser: () =>
        api.get('/api/v1/user/me'),
    
    updateProfile: (data: { name?: string; mail?: string; gender?: string; birthday?: string; bio?: string }) =>
        api.post('/api/v1/user/me/update', data),
    
    changePassword: (oldPassword: string, newPassword: string) =>
        api.post('/api/v1/user/passwd', { oldPassword, newPassword })
};

export const statsApi = {
    getDashboard: () =>
        api.get('/api/v1/admin/stats/dashboard'),
    
    getUserStats: () =>
        api.get('/api/v1/admin/stats/users'),
    
    getContentStats: () =>
        api.get('/api/v1/admin/stats/contents'),
    
    getCommentStats: () =>
        api.get('/api/v1/admin/stats/comments'),
    
    getTrend: (days: number = 30) =>
        api.get('/api/v1/admin/stats/trend', { params: { days } })
};

export const adminApi = {
    getPendingReviews: () =>
        api.get('/api/v1/admin/reviews'),
    
    approveReview: (reviewId: number) =>
        api.post(`/api/v1/admin/reviews/${reviewId}/approve`),
    
    rejectReview: (reviewId: number, reason: string) =>
        api.post(`/api/v1/admin/reviews/${reviewId}/reject`, { reason })
};

export const userApi = {
    getUsers: () =>
        api.get('/api/v1/user/users'),
    
    getUserByUsername: (username: string) =>
        api.get('/api/v1/user/users', { params: { username } }),
    
    createUser: (data: { name: string; email: string; password: string; role?: string; status?: string }) =>
        api.post('/api/v1/user/users', data),
    
    updateUser: (uid: number, data: { name?: string; email?: string; role?: string; status?: string }) =>
        api.post(`/api/v1/user/users/${uid}/update`, data),
    
    deleteUser: (uid: number) =>
        api.post(`/api/v1/user/users/del`, { uid })
};

export const contentApi = {
    getContents: (params?: { page?: number; limit?: number; status?: string; type?: string }) =>
        api.get('/api/v1/content', { params }),
    
    getContent: (cid: number) =>
        api.get(`/api/v1/content/${cid}`),
    
    createContent: (data: { title: string; content: string; text?: string; status?: string; type?: string; categories?: number[]; tags?: number[] }) =>
        api.post('/api/v1/content', data),
    
    updateContent: (cid: number, data: { title?: string; content?: string; text?: string; status?: string; categories?: number[], tags?: number[] }) =>
        api.post(`/api/v1/content/${cid}`, data),
    
    deleteContent: (cid: number) =>
        api.post(`/api/v1/content/${cid}/del`),
    
    publishContent: (cid: number) =>
        api.post(`/api/v1/content/${cid}/publish`),
    
    getMyContents: (params?: { page?: number; limit?: number }) =>
        api.get('/api/v1/content/my', { params })
};

export const commentApi = {
    getComments: (params?: { page?: number; limit?: number; status?: string; cid?: number }) =>
        api.get('/api/v1/comment', { params }),
    
    getMyComments: () =>
        api.get('/api/v1/comment/my'),
    
    getComment: (coid: number) =>
        api.get(`/api/v1/comment/${coid}`),
    
    createComment: (data: { cid: number; text: string; author?: string; mail?: string; url?: string; parent?: number }) =>
        api.post('/api/v1/comment', data),
    
    updateComment: (coid: number, data: { status?: string }) =>
        api.post(`/api/v1/comment/${coid}`, data),
    
    approveComment: (coid: number) =>
        api.post(`/api/v1/comment/${coid}/approve`),
    
    deleteComment: (coid: number) =>
        api.post(`/api/v1/comment/${coid}/del`),
    
    getContentComments: (cid: number) =>
        api.get(`/api/v1/comment/content/${cid}`)
};

export const metaApi = {
    getMetas: (params?: { type?: string }) =>
        api.get('/api/v1/meta', { params }),
    
    getMeta: (mid: number) =>
        api.get(`/api/v1/meta/${mid}`),
    
    createMeta: (data: { name: string; type?: string; slug?: string; description?: string; parent?: number }) =>
        api.post('/api/v1/meta', data),
    
    updateMeta: (mid: number, data: { name?: string; slug?: string; description?: string; parent?: number }) =>
        api.post(`/api/v1/meta/${mid}`, data),
    
    deleteMeta: (mid: number) =>
        api.post(`/api/v1/meta/${mid}/del`)
};

export const likeApi = {
    toggleLike: (targetType: string, targetId: number) =>
        api.post('/api/v1/like/toggle', { target_type: targetType, target_id: targetId }),
    
    getLikeCount: (targetType: string, targetId: number) =>
        api.get('/api/v1/like/count', { params: { target_type: targetType, target_id: targetId } }),
    
    getLikeStatus: (targetType: string, targetId: number) =>
        api.get('/api/v1/like/status', { params: { target_type: targetType, target_id: targetId } }),
    
    getMyLikesCount: () =>
        api.get('/api/v1/like/my/count')
};

export const settingsApi = {
    getSettings: () =>
        api.get('/api/v1/admin/settings'),
    
    updateSettings: (data: {
        siteName?: string;
        siteDescription?: string;
        adminEmail?: string;
        language?: string;
        timezone?: string;
        enableRegistration?: boolean;
        enableComments?: boolean;
        maintenanceMode?: boolean;
    }) =>
        api.post('/api/v1/admin/settings', data)
};

export default api;
