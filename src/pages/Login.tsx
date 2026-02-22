import React, { useState } from 'react';
import { Form, Input, Button, Card, Typography, message } from 'antd';
import { UserOutlined, LockOutlined, EyeInvisibleOutlined, EyeTwoTone } from '@ant-design/icons';
import { authApi } from '../api';
import { storage } from '../utils/storage';
import { config } from '../config';
import type { LoginProps } from '../types';

const { Title, Text } = Typography;

const Login: React.FC<LoginProps> = ({ onLogin }) => {
    const [loading, setLoading] = useState(false);
    const [form] = Form.useForm();

    const onFinish = async (values: { email: string; password: string }) => {
        setLoading(true);
        
        try {
            const response = await authApi.login(values.email, values.password);

            if (response.data.code === 200) {
                storage.setUserInfo({
                    uid: response.data.data.uid,
                    name: response.data.data.name,
                    email: response.data.data.email
                });
                
                message.success('登录成功');
                onLogin();
            } else {
                message.error(response.data.msg || '登录失败');
            }
        } catch (error: unknown) {
            console.error('登录错误:', error);
            const err = error as { response?: { data?: { message?: string } }; code?: string };
            if (err.response?.data?.message) {
                message.error(err.response.data.message);
            } else if (err.code === 'NETWORK_ERROR') {
                message.error('网络错误，请检查网络连接');
            } else {
                message.error('登录失败，请检查网络连接');
            }
        } finally {
            setLoading(false);
        }
    };

    return (
        <div style={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            minHeight: '100vh',
            backgroundImage: `url("${config.loginBackgroundUrl}")`,
            backgroundSize: 'cover',
            backgroundPosition: 'center',
            backgroundRepeat: 'no-repeat',
        }}>
            <Card 
                style={{ 
                    width: 400, 
                    boxShadow: '0 8px 32px rgba(0, 0, 0, 0.3)',
                    borderRadius: '12px',
                    backgroundColor: 'rgba(255, 255, 255, 0.5)',
                    backdropFilter: 'blur(20px)',
                    border: '1px solid rgba(255, 255, 255, 0.3)'
                }}
            >
                <div style={{ textAlign: 'center', marginBottom: 32 }}>
                    <Title level={2} style={{ marginBottom: 8, color: '#2996fbff' }}>
                        林影心语 <Text style={{ color: 'red', fontSize: '16px' }}>admin</Text>
                    </Title>
                </div>

                <Form
                    form={form}
                    name="login"
                    onFinish={onFinish}
                    autoComplete="off"
                    layout="vertical"
                >
                    <Form.Item
                        label="邮箱"
                        name="email"
                        rules={[{ required: true, message: '请输入邮箱地址' }]}
                    >
                        <Input 
                            prefix={<UserOutlined />} 
                            placeholder="请输入邮箱地址"
                            size="large"
                            style={{ 
                                backgroundColor: 'rgba(255, 255, 255, 0.85)', 
                                backdropFilter: 'blur(0px)',
                                border: '1px solid rgba(255, 255, 255, 0.5)'
                            }}
                        />
                    </Form.Item>

                    <Form.Item
                        label="密码"
                        name="password"
                        rules={[{ required: true, message: '请输入密码' }]}
                    >
                        <Input.Password
                            prefix={<LockOutlined />}
                            placeholder="请输入密码"
                            size="large"
                            iconRender={(visible) => (visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />)}
                            style={{ 
                                backgroundColor: 'rgba(255, 255, 255, 0.85)', 
                                backdropFilter: 'blur(0px)',
                                border: '1px solid rgba(255, 255, 255, 0.5)'
                            }}
                        />
                    </Form.Item>

                    <Form.Item>
                        <Button 
                            type="primary"
                            htmlType="submit" 
                            loading={loading}
                            size="large"
                            style={{ 
                                width: '100%',
                                backgroundColor: 'rgba(135, 133, 162, 0.6)',
                                borderColor: 'rgba(135, 133, 162, 1)',
                                backdropFilter: 'blur(10px)',
                                transition: 'all 0.3s ease'
                            }}
                            onMouseEnter={(e) => e.currentTarget.style.backgroundColor = 'rgba(135, 133, 162, 1)'}
                            onMouseLeave={(e) => e.currentTarget.style.backgroundColor = 'rgba(135, 133, 162, 0.9)'}
                        >
                            登录
                        </Button>
                    </Form.Item>

                    <Text 
                        style={{ color: '#f6f6f6', cursor: 'pointer' }} 
                        onClick={() => window.open(config.siteUrl, '_blank')}
                    >
                        返回首页
                    </Text>
                </Form>
            </Card>
        </div>
    );
};

export default Login;
