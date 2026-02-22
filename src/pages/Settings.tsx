import React, { useState, useEffect } from 'react';
import {
    Card,
    Form,
    Input,
    Button,
    Switch,
    Select,
    Typography,
    Divider,
    message,
    Space,
    Row,
    Col
} from 'antd';
import { SettingOutlined, SaveOutlined } from '@ant-design/icons';
import { settingsApi } from '../api';
import type { SettingsData } from '../types';

const { Title, Text } = Typography;
const { Option } = Select;
const { TextArea } = Input;

const Settings: React.FC = () => {
    const [form] = Form.useForm();
    const [loading, setLoading] = useState(false);
    const [saving, setSaving] = useState(false);

    useEffect(() => {
        fetchSettings();
    }, []);

    const fetchSettings = async () => {
        setLoading(true);
        try {
            const response = await settingsApi.getSettings();
            
            if (response.data.code === 200 && response.data.data) {
                form.setFieldsValue({
                    siteName: response.data.data.siteName || '',
                    siteDescription: response.data.data.siteDescription || '',
                    adminEmail: response.data.data.adminEmail || '',
                    language: response.data.data.language || 'zh-CN',
                    timezone: response.data.data.timezone || 'Asia/Shanghai',
                    enableRegistration: response.data.data.enableRegistration ?? true,
                    enableComments: response.data.data.enableComments ?? true,
                    maintenanceMode: response.data.data.maintenanceMode ?? false
                });
            }
        } catch (error) {
            console.error('获取设置失败:', error);
            message.error('获取设置失败');
        } finally {
            setLoading(false);
        }
    };

    const onFinish = async (values: SettingsData) => {
        setSaving(true);
        try {
            const response = await settingsApi.updateSettings(values);
            
            if (response.data.code === 200) {
                message.success('设置保存成功');
            } else {
                message.error(response.data.msg || '设置保存失败');
            }
        } catch (error) {
            console.error('保存设置失败:', error);
            message.error('保存设置失败');
        } finally {
            setSaving(false);
        }
    };

    return (
        <div>
            <Title level={2} style={{ marginBottom: 24 }}>
                <SettingOutlined /> 系统设置
            </Title>

            <Form
                form={form}
                layout="vertical"
                onFinish={onFinish}
                disabled={loading}
            >
                <Row gutter={[24, 16]}>
                    <Col span={24}>
                        <Card title="基本设置" style={{ marginBottom: 24 }}>
                            <Row gutter={[16, 0]}>
                                <Col span={12}>
                                    <Form.Item
                                        label="网站名称"
                                        name="siteName"
                                        rules={[{ required: true, message: '请输入网站名称' }]}
                                    >
                                        <Input placeholder="请输入网站名称" />
                                    </Form.Item>
                                </Col>
                                <Col span={12}>
                                    <Form.Item
                                        label="管理员邮箱"
                                        name="adminEmail"
                                        rules={[
                                            { required: true, message: '请输入管理员邮箱' },
                                            { type: 'email', message: '请输入有效的邮箱地址' }
                                        ]}
                                    >
                                        <Input placeholder="请输入管理员邮箱" />
                                    </Form.Item>
                                </Col>
                            </Row>

                            <Form.Item
                                label="网站描述"
                                name="siteDescription"
                            >
                                <TextArea 
                                    rows={3} 
                                    placeholder="请输入网站描述"
                                    showCount 
                                    maxLength={200}
                                />
                            </Form.Item>

                            <Row gutter={[16, 0]}>
                                <Col span={12}>
                                    <Form.Item
                                        label="语言"
                                        name="language"
                                    >
                                        <Select>
                                            <Option value="zh-CN">简体中文</Option>
                                            <Option value="en-US">English</Option>
                                            <Option value="ja-JP">日本語</Option>
                                        </Select>
                                    </Form.Item>
                                </Col>
                                <Col span={12}>
                                    <Form.Item
                                        label="时区"
                                        name="timezone"
                                    >
                                        <Select>
                                            <Option value="Asia/Shanghai">Asia/Shanghai (UTC+8)</Option>
                                            <Option value="America/New_York">America/New_York (UTC-5)</Option>
                                            <Option value="Europe/London">Europe/London (UTC+0)</Option>
                                        </Select>
                                    </Form.Item>
                                </Col>
                            </Row>
                        </Card>
                    </Col>

                    <Col span={24}>
                        <Card title="功能设置" style={{ marginBottom: 24 }}>
                            <Space direction="vertical" size="middle" style={{ width: '100%' }}>
                                <Form.Item
                                    label="允许用户注册"
                                    name="enableRegistration"
                                    valuePropName="checked"
                                >
                                    <Switch 
                                        checkedChildren="开启" 
                                        unCheckedChildren="关闭" 
                                    />
                                </Form.Item>

                                <Form.Item
                                    label="允许评论"
                                    name="enableComments"
                                    valuePropName="checked"
                                >
                                    <Switch 
                                        checkedChildren="开启" 
                                        unCheckedChildren="关闭" 
                                    />
                                </Form.Item>

                                <Form.Item
                                    label="维护模式"
                                    name="maintenanceMode"
                                    valuePropName="checked"
                                >
                                    <Switch 
                                        checkedChildren="开启" 
                                        unCheckedChildren="关闭" 
                                    />
                                </Form.Item>
                            </Space>
                        </Card>
                    </Col>

                    <Col span={24}>
                        <Card title="高级设置">
                            <Text type="secondary">
                                高级设置功能正在开发中，敬请期待...
                            </Text>
                        </Card>
                    </Col>
                </Row>

                <Divider />

                <Form.Item style={{ textAlign: 'right', marginBottom: 0 }}>
                    <Space>
                        <Button onClick={() => form.resetFields()}>
                            重置
                        </Button>
                        <Button 
                            type="primary" 
                            htmlType="submit" 
                            loading={saving}
                            icon={<SaveOutlined />}
                        >
                            保存设置
                        </Button>
                    </Space>
                </Form.Item>
            </Form>
        </div>
    );
};

export default Settings;
