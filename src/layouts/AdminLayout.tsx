import React, { useState, useEffect } from 'react';
import { Layout, Menu, Avatar, Dropdown, Space, Typography, theme, message, Drawer } from 'antd';
import { useNavigate, useLocation, Outlet } from 'react-router-dom';
import {
    DashboardOutlined,
    UserOutlined,
    BellOutlined,
    SettingOutlined,
    LogoutOutlined,
    MenuFoldOutlined,
    MenuUnfoldOutlined,
    MenuOutlined,
    FileTextOutlined,
    CommentOutlined
} from '@ant-design/icons';
import axios from 'axios';
import { API_BASE_URL } from '../api';
import type { AdminLayoutProps } from '../types';

const { Header, Sider, Content } = Layout;
const { Text } = Typography;

const AdminLayout: React.FC<AdminLayoutProps> = ({ onLogout }) => {
    const [collapsed, setCollapsed] = useState(true);
    const [adminInfo, setAdminInfo] = useState({ name: '管理员', identity: '管理员' });
    const [mobile, setMobile] = useState(false);
    const [drawerVisible, setDrawerVisible] = useState(false);
    const navigate = useNavigate();
    const location = useLocation();
    const {
        token: { colorBgContainer, borderRadiusLG },
    } = theme.useToken();

    useEffect(() => {
        fetchAdminInfo();
        
        const handleResize = () => {
            const isMobile = window.innerWidth < 768;
            setMobile(isMobile);
            setCollapsed(isMobile);
        };
        
        window.addEventListener('resize', handleResize);
        handleResize();
        
        return () => window.removeEventListener('resize', handleResize);
    }, []);

    const fetchAdminInfo = async () => {
    try {
      // 从localStorage获取用户信息
      const userInfo = localStorage.getItem('userInfo');
      if (userInfo) {
        const parsedInfo = JSON.parse(userInfo);
        setAdminInfo({
          name: parsedInfo.name || '管理员',
          identity: parsedInfo.group || '管理员'
        });
      }

      // 从服务器获取最新用户信息
      const response = await axios.get(`${API_BASE_URL}/api/v1/user/me`, {
        withCredentials: true,
        headers: {
          'Content-Type': 'application/json'
        }
      });

      if (response.data.code === 200) {
        const userData = response.data.data;
        
        const userRoles = userData.roles || [];
        const isAdmin = userRoles.includes('admin') || userRoles.includes('superadmin');
        const isSuperadmin = userRoles.includes('superadmin');
        
        if (!isAdmin) {
          message.error('权限不足，仅管理员可访问此页面');
          localStorage.removeItem('userInfo');
          window.location.reload();
          return;
        }
        
        setAdminInfo({
          name: userData.name || '管理员',
          identity: isSuperadmin ? '超级管理员' : (userRoles[0] || '管理员')
        });
        
        // 更新localStorage中的用户信息
        localStorage.setItem('userInfo', JSON.stringify({
          uid: userData.uid,
          name: userData.name,
          email: userData.email,
          group: userData.roles?.[0] || '用户',
          isAdmin: isAdmin
        }));
      }
    } catch (error: unknown) {
      console.error('获取管理员信息失败:', error);
      const err = error as { response?: { status?: number } };
      if (err.response?.status === 401) {
        // 未登录，清除localStorage
        localStorage.removeItem('userInfo');
      }
    }
  };

    const menuItems = [
        {
        key: '/',
        icon: <DashboardOutlined />,
        label: '仪表盘',
        },
        {
        key: '/users',
        icon: <UserOutlined />,
        label: '用户管理',
        },
        {
        key: '/articles',
        icon: <FileTextOutlined />,
        label: '文章管理',
        },
        {
        key: '/comments',
        icon: <CommentOutlined />,
        label: '评论管理',
        },
        {
        key: '/settings',
        icon: <SettingOutlined />,
        label: '系统设置',
        },
    ];

    const userMenuItems = [
        {
        key: 'profile',
        icon: <UserOutlined />,
        label: '个人信息',
        },
        {
        key: 'logout',
        icon: <LogoutOutlined />,        label: '退出登录',
        onClick: onLogout,
        },
    ];

    const handleMenuClick = (key: string) => {
        navigate(key);
        if (mobile) {
            setDrawerVisible(false);
        }
    };

    const siderWidth = collapsed ? 80 : 200;

    return (
        <Layout style={{ minHeight: '100vh' }}>
        {mobile ? (
            <>
                <Drawer
                    title={
                        <Text strong style={{ color: '#1890ff' }}>
                            Lmoadll
                        </Text>
                    }
                    placement="left"
                    onClose={() => setDrawerVisible(false)}
                    open={drawerVisible}
                    width={200}
                >
                    <Menu
                        mode="inline"
                        selectedKeys={[location.pathname]}
                        items={menuItems}
                        onClick={({ key }) => handleMenuClick(key)}
                        style={{ borderRight: 0 }}
                    />
                </Drawer>
                <Sider 
                    trigger={null} 
                    collapsed={true}
                    width={0}
                    style={{ display: 'none' }}
                />
            </>
        ) : (
            <Sider 
            trigger={null} 
            collapsible 
            collapsed={collapsed}
            style={{
            background: colorBgContainer,
            boxShadow: '2px 0 8px rgba(0, 0, 0, 0.1)',
            position: 'fixed',
            left: 0,
            top: 0,
            bottom: 0,
            zIndex: 100,
            overflow: 'auto'
            }}
        >
            <div style={{
            height: 64,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            borderBottom: '1px solid #f0f0f0'
            }}>
            <Text strong style={{ fontSize: collapsed ? 14 : 16, color: '#1890ff' }}>
                {collapsed ? 'L' : 'Lmoadll'}
            </Text>
            </div>
            
            <Menu
            mode="inline"
            selectedKeys={[location.pathname]}
            items={menuItems}
            onClick={({ key }) => handleMenuClick(key)}
            style={{ borderRight: 0, marginTop: 16 }}
            />
        </Sider>
        )}
        
        <Layout style={{ marginLeft: mobile ? 0 : siderWidth }}>
            <Header style={{
            padding: '0 24px',
            background: colorBgContainer,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
            boxShadow: '0 2px 8px rgba(0, 0, 0, 0.1)',
            position: 'fixed',
            top: 0,
            right: 0,
            left: mobile ? 0 : siderWidth,
            zIndex: 99
            }}>
            <div>
                {mobile ? (
                    <MenuOutlined 
                        onClick={() => setDrawerVisible(true)}
                        style={{ fontSize: 18, cursor: 'pointer' }}
                    />
                ) : (
                    React.createElement(collapsed ? MenuUnfoldOutlined : MenuFoldOutlined, {
                    className: 'trigger',
                    onClick: () => setCollapsed(!collapsed),
                    style: { fontSize: 18, cursor: 'pointer' }
                    })
                )}
            </div>
            
            <Space size="middle">
                <BellOutlined style={{ fontSize: 18, cursor: 'pointer' }} />
                <Dropdown 
                menu={{ items: userMenuItems }}
                placement="bottomRight"
                arrow
                >
                <Space style={{ cursor: 'pointer' }}>
                    <Avatar size="small" icon={<UserOutlined />} />
                    <Text>{adminInfo.name}</Text>
                </Space>
                </Dropdown>
            </Space>
            </Header>
            
            <Content style={{
            margin: '24px',
            marginTop: 88,
            padding: 24,
            background: colorBgContainer,
            borderRadius: borderRadiusLG,
            minHeight: 280
            }}>
            <Outlet />
            </Content>
        </Layout>
        </Layout>
    );
};

export default AdminLayout;
