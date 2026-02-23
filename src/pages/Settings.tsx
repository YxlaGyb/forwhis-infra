import React, { useState, useEffect, useCallback } from 'react';
import {
    Card,
    Form,
    Input,
    Button,
    Switch,
    Typography,
    message,
    Space,
    Row,
    Col,
    Menu,
    Badge,
    Tooltip,
    Divider,
    Alert
} from 'antd';
import {
    SettingOutlined,
    SaveOutlined,
    GlobalOutlined,
    FileSearchOutlined,
    ReloadOutlined,
    ExclamationCircleOutlined
} from '@ant-design/icons';
import { settingsApi } from '../api';
import type {
    AllSettings,
    SettingsSection,
    SettingsSectionConfig
} from '../types';

const { Title, Text } = Typography;
const { TextArea } = Input;

const sectionConfigs: SettingsSectionConfig[] = [
    {
        key: 'site',
        title: '站点信息',
        icon: 'GlobalOutlined',
        description: '网站名称、Logo、描述等基本信息'
    },
    {
        key: 'feature',
        title: '功能开关',
        icon: 'SettingOutlined',
        description: '注册、评论、验证码等功能控制'
    },
    {
        key: 'seo',
        title: 'SEO设置',
        icon: 'FileSearchOutlined',
        description: '关键词、作者、站点地图'
    }
];

const defaultSettings: AllSettings = {
    site: {
        siteName: '',
        siteDescription: '',
        siteLogo: '',
        favicon: '',
        icp: ''
    },
    feature: {
        enableRegistration: true,
        enableComments: true,
        enableEmailVerify: false,
        enableCaptcha: true
    },
    seo: {
        siteKeywords: '',
        siteAuthor: '',
        robotsTxt: '',
        enableSitemap: true
    }
};

const iconMap: Record<string, React.ReactNode> = {
    GlobalOutlined: <GlobalOutlined />,
    SettingOutlined: <SettingOutlined />,
    FileSearchOutlined: <FileSearchOutlined />
};

const Settings: React.FC = () => {
    const [form] = Form.useForm<AllSettings>();
    const [activeSection, setActiveSection] = useState<SettingsSection>('site');
    const [loading, setLoading] = useState(false);
    const [saving, setSaving] = useState(false);
    const [hasChanges, setHasChanges] = useState(false);
    const [originalValues, setOriginalValues] = useState<AllSettings>(defaultSettings);

    const fetchSettings = useCallback(async () => {
        setLoading(true);
        try {
            const response = await settingsApi.getSettings();
            if (response?.data?.code === 200 && response?.data?.data) {
                const data = response.data.data;
                const mergedSettings = { ...defaultSettings, ...data };
                form.setFieldsValue(mergedSettings);
                setOriginalValues(mergedSettings);
            } else {
                throw new Error('Invalid response');
            }
        } catch (error) {
            console.error('获取设置失败:', error);
            message.warning('无法连接到服务器，使用本地默认配置');
            form.setFieldsValue(defaultSettings);
            setOriginalValues(defaultSettings);
        } finally {
            setLoading(false);
        }
    }, [form]);

    useEffect(() => {
        fetchSettings();
    }, [fetchSettings]);

    const handleValuesChange = () => {
        setHasChanges(true);
    };

    const handleSave = async () => {
        try {
            const values = await form.validateFields();
            setSaving(true);
            const response = await settingsApi.updateSettings(values);

            if (response.data.code === 200) {
                message.success('设置保存成功');
                setOriginalValues(values as AllSettings);
                setHasChanges(false);
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

    const handleReset = () => {
        form.setFieldsValue(originalValues);
        setHasChanges(false);
        message.info('已重置为上次保存的设置');
    };

    const handleResetToDefault = () => {
        form.setFieldsValue(defaultSettings);
        setHasChanges(true);
        message.info('已重置为默认设置');
    };

    const renderSiteSection = () => (
        <Space direction="vertical" style={{ width: '100%' }} size="large">
            <Alert
                message="站点信息设置"
                description="配置网站的基本信息，这些信息将显示在网站标题、搜索引擎结果和社交媒体分享中。"
                type="info"
                showIcon
            />
            <Row gutter={[24, 16]}>
                <Col xs={24} md={12}>
                    <Form.Item
                        label="网站名称"
                        name={['site', 'siteName']}
                        rules={[{ required: true, message: '请输入网站名称' }]}
                        tooltip="显示在浏览器标签页和网站头部"
                    >
                        <Input placeholder="例如：我的CMS系统" maxLength={50} showCount />
                    </Form.Item>
                </Col>
                <Col xs={24} md={12}>
                    <Form.Item
                        label="网站Logo"
                        name={['site', 'siteLogo']}
                        tooltip="建议尺寸 200x60 像素"
                    >
                        <Input placeholder="https://example.com/logo.png" />
                    </Form.Item>
                </Col>
            </Row>
            <Form.Item
                label="网站描述"
                name={['site', 'siteDescription']}
                tooltip="用于搜索引擎优化（SEO）"
            >
                <TextArea
                    rows={3}
                    placeholder="简要描述您的网站，不超过200字"
                    maxLength={200}
                    showCount
                />
            </Form.Item>
            <Row gutter={[24, 16]}>
                <Col xs={24} md={12}>
                    <Form.Item
                        label="网站图标 (Favicon)"
                        name={['site', 'favicon']}
                        tooltip="建议尺寸 32x32 或 64x64 像素"
                    >
                        <Input placeholder="https://example.com/favicon.ico" />
                    </Form.Item>
                </Col>
                <Col xs={24} md={12}>
                    <Form.Item
                        label="ICP备案号"
                        name={['site', 'icp']}
                        tooltip="中国大陆网站必填"
                    >
                        <Input placeholder="例如：京ICP备12345678号" />
                    </Form.Item>
                </Col>
            </Row>
        </Space>
    );

    const renderFeatureSection = () => (
        <Space direction="vertical" style={{ width: '100%' }} size="large">
            <Alert
                message="功能开关"
                description="控制网站的各项功能模块的开启和关闭。"
                type="info"
                showIcon
            />
            <Card size="small" title="用户功能">
                <Space direction="vertical" style={{ width: '100%' }}>
                    <Form.Item
                        label="允许用户注册"
                        name={['feature', 'enableRegistration']}
                        valuePropName="checked"
                        style={{ marginBottom: 8 }}
                    >
                        <Switch checkedChildren="开启" unCheckedChildren="关闭" />
                    </Form.Item>
                    <Text type="secondary" style={{ fontSize: 12 }}>
                        关闭后新用户将无法注册，仅允许现有用户登录
                    </Text>
                </Space>
            </Card>
            <Card size="small" title="内容功能">
                <Space direction="vertical" style={{ width: '100%' }}>
                    <Form.Item
                        label="允许评论"
                        name={['feature', 'enableComments']}
                        valuePropName="checked"
                        style={{ marginBottom: 8 }}
                    >
                        <Switch checkedChildren="开启" unCheckedChildren="关闭" />
                    </Form.Item>
                    <Text type="secondary" style={{ fontSize: 12 }}>
                        关闭后全站将禁用评论功能
                    </Text>
                </Space>
            </Card>
            <Card size="small" title="安全功能">
                <Space direction="vertical" style={{ width: '100%' }}>
                    <Form.Item
                        label="邮箱验证"
                        name={['feature', 'enableEmailVerify']}
                        valuePropName="checked"
                        style={{ marginBottom: 8 }}
                    >
                        <Switch checkedChildren="开启" unCheckedChildren="关闭" />
                    </Form.Item>
                    <Text type="secondary" style={{ fontSize: 12, marginBottom: 16 }}>
                        开启后新用户注册需要验证邮箱
                    </Text>
                    <Form.Item
                        label="验证码"
                        name={['feature', 'enableCaptcha']}
                        valuePropName="checked"
                        style={{ marginBottom: 8 }}
                    >
                        <Switch checkedChildren="开启" unCheckedChildren="关闭" />
                    </Form.Item>
                    <Text type="secondary" style={{ fontSize: 12 }}>
                        开启后在登录和注册时需要输入验证码
                    </Text>
                </Space>
            </Card>
        </Space>
    );

    const renderSeoSection = () => (
        <Space direction="vertical" style={{ width: '100%' }} size="large">
            <Alert
                message="SEO设置"
                description="优化搜索引擎收录和展示效果。"
                type="info"
                showIcon
            />
            <Row gutter={[24, 16]}>
                <Col xs={24} md={12}>
                    <Form.Item
                        label="网站关键词"
                        name={['seo', 'siteKeywords']}
                        tooltip="多个关键词用逗号分隔"
                    >
                        <Input placeholder="CMS, 博客, 内容管理" />
                    </Form.Item>
                </Col>
                <Col xs={24} md={12}>
                    <Form.Item
                        label="网站作者"
                        name={['seo', 'siteAuthor']}
                    >
                        <Input placeholder="网站所有者或组织名称" />
                    </Form.Item>
                </Col>
            </Row>
            <Form.Item
                label="Robots.txt"
                name={['seo', 'robotsTxt']}
                tooltip="控制搜索引擎爬虫的访问规则"
            >
                <TextArea
                    rows={6}
                    placeholder={`User-agent: *
Allow: /
Disallow: /admin/
Sitemap: https://example.com/sitemap.xml`}
                    style={{ fontFamily: 'monospace' }}
                />
            </Form.Item>
            <Form.Item
                label="启用站点地图"
                name={['seo', 'enableSitemap']}
                valuePropName="checked"
            >
                <Switch checkedChildren="开启" unCheckedChildren="关闭" />
            </Form.Item>
        </Space>
    );

    const sectionRenderers: Record<SettingsSection, () => React.ReactNode> = {
        site: renderSiteSection,
        feature: renderFeatureSection,
        seo: renderSeoSection
    };

    const currentConfig = sectionConfigs.find(c => c.key === activeSection)!;

    // 渲染内容区域，带错误处理
    const renderContent = () => {
        try {
            const renderer = sectionRenderers[activeSection];
            if (!renderer) {
                return (
                    <Alert
                        message="页面加载失败"
                        description="无法加载当前设置页面，请尝试刷新或切换其他选项"
                        type="error"
                        showIcon
                        action={
                            <Button onClick={() => window.location.reload()}>
                                刷新页面
                            </Button>
                        }
                    />
                );
            }
            return renderer();
        } catch (error) {
            console.error('渲染设置页面失败:', error);
            return (
                <Alert
                    message="页面渲染出错"
                    description="设置页面加载时发生错误，已自动切换到默认视图"
                    type="warning"
                    showIcon
                    action={
                        <Button onClick={() => setActiveSection('site')}>
                            返回站点信息
                        </Button>
                    }
                />
            );
        }
    };

    return (
        <div style={{ height: 'calc(100vh - 112px)', display: 'flex', gap: 24 }}>
            {/* 左侧导航 */}
            <Card
                style={{
                    width: 260,
                    flexShrink: 0
                }}
                bodyStyle={{ padding: 0 }}
            >
                <div style={{ padding: '16px 20px', borderBottom: '1px solid #f0f0f0' }}>
                    <Title level={5} style={{ margin: 0 }}>
                        <SettingOutlined style={{ marginRight: 8 }} />
                        系统设置
                    </Title>
                </div>
                <Menu
                    mode="inline"
                    selectedKeys={[activeSection]}
                    onClick={({ key }) => setActiveSection(key as SettingsSection)}
                    style={{ border: 'none', padding: '8px 0' }}
                    theme="light"
                    items={sectionConfigs.map(config => ({
                        key: config.key,
                        icon: (
                            <div style={{ 
                                display: 'flex', 
                                alignItems: 'center', 
                                justifyContent: 'center',
                                fontSize: 18,
                                color: activeSection === config.key ? '#1890ff' : '#595959'
                            }}>
                                {iconMap[config.icon]}
                            </div>
                        ),
                        label: (
                            <div style={{ 
                                padding: '4px 0',
                                lineHeight: '1.4'
                            }}>
                                <div style={{ 
                                    fontWeight: activeSection === config.key ? 600 : 500, 
                                    fontSize: 14, 
                                    marginBottom: 4,
                                    color: activeSection === config.key ? '#1890ff' : '#262626'
                                }}>
                                    {config.title}
                                </div>
                                <div style={{ 
                                    fontSize: 12, 
                                    color: activeSection === config.key ? '#69c0ff' : '#8c8c8c', 
                                    lineHeight: '1.3'
                                }}>
                                    {config.description}
                                </div>
                            </div>
                        ),
                        style: {
                            height: 'auto',
                            padding: '12px 16px',
                            margin: '4px 0',
                            borderRadius: 0
                        }
                    }))}
                />
            </Card>

            {/* 右侧内容 */}
            <div style={{ flex: 1, overflow: 'auto' }}>
                <Form
                    form={form}
                    layout="vertical"
                    onValuesChange={handleValuesChange}
                    disabled={loading}
                    style={{ maxWidth: 800 }}
                >
                    <Card
                        title={
                            <Space>
                                {iconMap[currentConfig.icon]}
                                <span>{currentConfig.title}</span>
                                {hasChanges && (
                                    <Badge dot color="red">
                                        <span style={{ fontSize: 12, color: '#999' }}>有未保存的更改</span>
                                    </Badge>
                                )}
                            </Space>
                        }
                        extra={
                            <Space>
                                <Tooltip title="重置为上次保存的设置">
                                    <Button
                                        icon={<ReloadOutlined />}
                                        onClick={handleReset}
                                        disabled={!hasChanges || saving}
                                    >
                                        重置
                                    </Button>
                                </Tooltip>
                                <Button
                                    type="primary"
                                    icon={<SaveOutlined />}
                                    onClick={handleSave}
                                    loading={saving}
                                    disabled={!hasChanges}
                                >
                                    保存
                                </Button>
                            </Space>
                        }
                        loading={loading}
                    >
                        {renderContent()}
                    </Card>

                    <Divider />

                    <Card size="small" title="危险操作" style={{ borderColor: '#ff4d4f' }}>
                        <Space>
                            <ExclamationCircleOutlined style={{ color: '#ff4d4f' }} />
                            <Text type="secondary">重置所有设置为默认值</Text>
                            <Button
                                danger
                                size="small"
                                onClick={handleResetToDefault}
                            >
                                重置为默认
                            </Button>
                        </Space>
                    </Card>
                </Form>
            </div>
        </div>
    );
};

export default Settings;
