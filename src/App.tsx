import { useState, useEffect } from 'react';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import axios from 'axios';
import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import Login from './pages/Login';
import AdminLayout from './layouts/AdminLayout';
import Dashboard from './pages/Dashboard';
import UserManagement from './pages/UserManagement';
import ArticleManagement from './pages/ArticleManagement';
import CommentManagement from './pages/CommentManagement';
import Settings from './pages/Settings';
import { API_BASE_URL } from './api';
import './App.css';

const App = () => {
    const [isLoggedIn, setIsLoggedIn] = useState(false);
    const [canAccessAdmin, setCanAccessAdmin] = useState(false);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        checkLoginStatus();
    }, []);

    const checkLoginStatus = async () => {
        try {
            const userInfo = localStorage.getItem('userInfo');
            if (userInfo) {
                const response = await axios.get(`${API_BASE_URL}/api/v1/user/me`, {
                    withCredentials: true,
                    headers: { 'Content-Type': 'application/json' }
                });
                
                if (response.data.code === 200) {
                    const data = response.data.data;
                    localStorage.setItem('userInfo', JSON.stringify({
                        uid: data.uid,
                        name: data.name,
                        email: data.email,
                        permissions: data.permissions,
                        roles: data.roles,
                        canAccessAdmin: data.canAccessAdmin
                    }));
                    
                    if (data.canAccessAdmin) {
                        setCanAccessAdmin(true);
                        setIsLoggedIn(true);
                    } else {
                        handleLogout();
                    }
                } else {
                    localStorage.removeItem('userInfo');
                }
            }
        } catch (error) {
            console.error('检查登录状态失败:', error);
            localStorage.removeItem('userInfo');
        } finally {
            setLoading(false);
        }
    };

    const handleLogin = () => {
        checkLoginStatus();
    };

    const handleLogout = async () => {
        try {
            await axios.post(`${API_BASE_URL}/api/v1/user/logout`, {}, {
                withCredentials: true,
                headers: { 'Content-Type': 'application/json' }
            });
        } catch (error) {
            console.error('登出API调用失败:', error);
        } finally {
            localStorage.removeItem('userInfo');
            document.cookie = 'forestwhisper=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
            setIsLoggedIn(false);
            setCanAccessAdmin(false);
        }
    };

    if (loading) {
        return (
            <ConfigProvider locale={zhCN}>
                <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh' }}>
                    <div>加载中...</div>
                </div>
            </ConfigProvider>
        );
    }

    if (!canAccessAdmin && isLoggedIn) {
        return (
            <ConfigProvider locale={zhCN}>
                <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: '100vh', flexDirection: 'column', gap: 16 }}>
                    <div>您没有管理员权限</div>
                    <button onClick={handleLogout} style={{ padding: '8px 16px', cursor: 'pointer' }}>退出登录</button>
                </div>
            </ConfigProvider>
        );
    }

    return (
        <ConfigProvider locale={zhCN}>
            <BrowserRouter>
                <Routes>
                    <Route path="/login" element={
                        isLoggedIn ? <Navigate to="/" replace /> : <Login onLogin={handleLogin} />
                    } />
                    <Route path="/" element={
                        isLoggedIn ? <AdminLayout onLogout={handleLogout} /> : <Navigate to="/login" replace />
                    }>
                        <Route index element={<Dashboard />} />
                        <Route path="users" element={<UserManagement />} />
                        <Route path="articles" element={<ArticleManagement />} />
                        <Route path="comments" element={<CommentManagement />} />
                        <Route path="settings" element={<Settings />} />
                    </Route>
                </Routes>
            </BrowserRouter>
        </ConfigProvider>
    );
};

export default App;
