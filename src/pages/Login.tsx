import React, { useState } from 'react';
import { Form, Input, Button, Card, Typography, message } from 'antd';
import { UserOutlined, LockOutlined, EyeInvisibleOutlined, EyeTwoTone } from '@ant-design/icons';
import { authApi } from '../api';

const { Title, Text } = Typography;

interface LoginProps {
    onLogin: () => void;
}

const Login: React.FC<LoginProps> = ({ onLogin }) => {
    const [loading, setLoading] = useState(false);
    const [form] = Form.useForm();

    const onFinish = async (values: { email: string; password: string }) => {
        setLoading(true);
        
        try {
            const response = await authApi.login(values.email, values.password);

            if (response.data.code === 200) {
                localStorage.setItem('userInfo', JSON.stringify({
                    uid: response.data.data.uid,
                    name: response.data.data.name,
                    email: response.data.data.email
                }));
                
                message.success('登录成功');
                onLogin();
            } else {
                message.error(response.data.msg || '登录失败');
            }
        } catch (error: any) {
            console.error('登录错误:', error);
            if (error.response?.data?.message) {
                message.error(error.response.data.message);
            } else if (error.code === 'NETWORK_ERROR') {
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
            backgroundImage: 'url("https://api.klbbx.cc/random_image/Genshin_Impact")',
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
                    border: '1px solid rgba(255, 255, 255, -10)'
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
                        label="Email"
                        name="email"
                        rules={[{ required: true, message: '请输入邮箱地址' }]}
                    >
                        <Input 
                            prefix={<UserOutlined />} 
                            placeholder="请输入邮箱地址"
                            size="large"
                            style={{ 
                                backgroundColor: 'rgba(255, 255, 255, -1)', 
                                backdropFilter: 'blur(0px)',
                                border: '1px solid rgba(255, 255, 255, 0.5)'
                            }}
                        />
                    </Form.Item>

                    <Form.Item
                        label="Password"
                        name="password"
                        rules={[{ required: true, message: '请输入密码' }]}
                    >
                        <Input.Password
                            prefix={<LockOutlined />}
                            placeholder="请输入密码"
                            size="large"
                            iconRender={(visible) => (visible ? <EyeTwoTone /> : <EyeInvisibleOutlined />)}
                            style={{ 
                                backgroundColor: 'rgba(255, 255, 255, -1)', 
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
                            Login
                        </Button>
                    </Form.Item>

                    <Text 
                        style={{ color: '#f6f6f6', cursor: 'pointer' }} 
                        onClick={() => window.open('https://banana404.top', '_blank')}
                    >
                        Go to home
                    </Text>
                </Form>
            </Card>
        </div>
    );
};

export default Login;
