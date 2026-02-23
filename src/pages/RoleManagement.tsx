import React, { useState } from 'react';
import {
    Card,
    Button,
    Space,
    Tag,
    Modal,
    Form,
    Input,
    Checkbox,
    message,
    Segmented,
    Avatar
} from 'antd';
import {
    PlusOutlined,
    EditOutlined,
    DeleteOutlined,
    SafetyOutlined,
    UserOutlined
} from '@ant-design/icons';
import { useNavigate, useLocation } from 'react-router-dom';
import { ResponsiveTable } from '../components/ResponsiveTable';
import type { ResponsiveColumn, ActionItem } from '../components/ResponsiveTable';

const { TextArea } = Input;

interface Role {
    id: number;
    name: string;
    description: string;
    permissions: string[];
    userCount: number;
}

const PERMISSION_CATEGORIES = [
    {
        category: '用户管理',
        permissions: [
            { key: 'user:read', name: '查看用户' },
            { key: 'user:create', name: '创建用户' },
            { key: 'user:edit', name: '编辑用户' },
            { key: 'user:delete', name: '删除用户' },
        ]
    },
    {
        category: '内容管理',
        permissions: [
            { key: 'post:read', name: '查看文章' },
            { key: 'post:create', name: '创建文章' },
            { key: 'post:edit', name: '编辑文章' },
            { key: 'post:delete', name: '删除文章' },
            { key: 'post:publish', name: '发布文章' },
            { key: 'comment:read', name: '查看评论' },
            { key: 'comment:create', name: '创建评论' },
            { key: 'comment:delete', name: '删除评论' },
            { key: 'comment:moderate', name: '审核评论' },
        ]
    },
    {
        category: '系统管理',
        permissions: [
            { key: 'system:config', name: '系统配置' },
            { key: 'system:logs', name: '查看日志' },
            { key: 'system:backup', name: '数据备份' },
            { key: 'system:plugin', name: '插件管理' },
            { key: 'settings:edit', name: '编辑系统设置' },
        ]
    }
];

const ROLE_PERMISSIONS: Record<string, string[]> = {
    superadmin: ['user:read', 'user:create', 'user:edit', 'user:delete', 'post:read', 'post:create', 'post:edit', 'post:delete', 'post:publish', 'comment:read', 'comment:create', 'comment:delete', 'comment:moderate', 'system:config', 'system:logs', 'system:backup', 'system:plugin', 'settings:edit'],
    admin: ['user:read', 'user:create', 'user:edit', 'user:delete', 'post:read', 'post:create', 'post:edit', 'post:delete', 'post:publish', 'comment:read', 'comment:create', 'comment:delete', 'comment:moderate', 'system:logs', 'settings:edit'],
    editor: ['post:read', 'post:create', 'post:edit', 'post:publish', 'comment:read', 'comment:create'],
    user: ['post:read', 'comment:read', 'comment:create'],
};

const RoleManagement: React.FC = () => {
    const navigate = useNavigate();
    const location = useLocation();
    const [roles, setRoles] = useState<Role[]>([
        { id: 1, name: 'superadmin', description: '超级管理员，拥有所有权限', permissions: ROLE_PERMISSIONS.superadmin, userCount: 1 },
        { id: 2, name: 'admin', description: '管理员，可以管理用户和内容', permissions: ROLE_PERMISSIONS.admin, userCount: 0 },
        { id: 3, name: 'editor', description: '编辑，可以管理文章', permissions: ROLE_PERMISSIONS.editor, userCount: 0 },
        { id: 4, name: 'user', description: '普通用户', permissions: ROLE_PERMISSIONS.user, userCount: 1 },
    ]);
    const [loading] = useState(false);
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [editingRole, setEditingRole] = useState<Role | null>(null);
    const [form] = Form.useForm();
    const [selectedPermissions, setSelectedPermissions] = useState<string[]>([]);

    const columns: ResponsiveColumn<Role>[] = [
        {
            title: '角色名称',
            dataIndex: 'name',
            key: 'name',
            width: 120,
            render: (name: string) => {
                const colorMap: Record<string, string> = {
                    superadmin: 'red',
                    admin: 'orange',
                    editor: 'blue',
                    user: 'green',
                };
                const nameMap: Record<string, string> = {
                    superadmin: '超级管理员',
                    admin: '管理员',
                    editor: '编辑',
                    user: '普通用户',
                };
                return <Tag color={colorMap[name] || 'default'}>{nameMap[name] || name}</Tag>;
            },
        },
        {
            title: '描述',
            dataIndex: 'description',
            key: 'description',
            width: 200,
            priority: 'medium',
        },
        {
            title: '权限数量',
            dataIndex: 'permissions',
            key: 'permissions',
            width: 100,
            priority: 'low',
            render: (permissions: string[]) => `${permissions.length} 个权限`,
        },
        {
            title: '用户数量',
            dataIndex: 'userCount',
            key: 'userCount',
            width: 100,
            priority: 'low',
        },
    ];

    const getActions = (record: Role): ActionItem[] => [
        {
            key: 'edit',
            text: '编辑',
            icon: <EditOutlined />,
            onClick: () => handleEdit(record),
            disabled: record.name === 'superadmin',
        },
        {
            key: 'delete',
            text: '删除',
            icon: <DeleteOutlined />,
            danger: true,
            onClick: () => {
                if (record.name !== 'superadmin' && record.userCount === 0) {
                    handleDelete(record.id);
                }
            },
            disabled: record.name === 'superadmin' || record.userCount > 0,
        },
    ];

    const getCardTitle = (record: Role) => {
        const colorMap: Record<string, string> = {
            superadmin: 'red',
            admin: 'orange',
            editor: 'blue',
            user: 'green',
        };
        const nameMap: Record<string, string> = {
            superadmin: '超级管理员',
            admin: '管理员',
            editor: '编辑',
            user: '普通用户',
        };
        return (
            <Space>
                <Avatar size="small" icon={<SafetyOutlined />} />
                <Tag color={colorMap[record.name] || 'default'}>{nameMap[record.name] || record.name}</Tag>
            </Space>
        );
    };

    const getCardSubtitle = (record: Role) => record.description;

    const handleAdd = () => {
        setEditingRole(null);
        setSelectedPermissions([]);
        form.resetFields();
        setIsModalVisible(true);
    };

    const handleEdit = (role: Role) => {
        setEditingRole(role);
        setSelectedPermissions(role.permissions);
        form.setFieldsValue({
            name: role.name,
            description: role.description,
        });
        setIsModalVisible(true);
    };

    const handleDelete = (id: number) => {
        setRoles(roles.filter(r => r.id !== id));
        message.success('角色删除成功');
    };

    const handleModalOk = async () => {
        try {
            const values = await form.validateFields();
            
            if (editingRole) {
                setRoles(roles.map(r => 
                    r.id === editingRole.id 
                        ? { ...r, ...values, permissions: selectedPermissions }
                        : r
                ));
                message.success('角色更新成功');
            } else {
                const newRole: Role = {
                    id: Date.now(),
                    name: values.name,
                    description: values.description,
                    permissions: selectedPermissions,
                    userCount: 0,
                };
                setRoles([...roles, newRole]);
                message.success('角色创建成功');
            }
            
            setIsModalVisible(false);
            form.resetFields();
            setSelectedPermissions([]);
        } catch (error) {
            console.error('表单验证失败:', error);
        }
    };

    const handleModalCancel = () => {
        setIsModalVisible(false);
        form.resetFields();
        setSelectedPermissions([]);
    };

    const handlePermissionChange = (checkedValues: string[]) => {
        setSelectedPermissions(checkedValues);
    };

    const isSystemRole = editingRole?.name === 'superadmin';

    return (
        <div>
            <div style={{ marginBottom: 16 }}>
                <Segmented
                    value={location.pathname === '/roles' ? 'roles' : 'users'}
                    onChange={(value) => navigate(value === 'roles' ? '/roles' : '/users')}
                    options={[
                        { label: <><UserOutlined /> 用户管理</>, value: 'users' },
                        { label: <><SafetyOutlined /> 角色管理 Bata</>, value: 'roles' },
                    ]}
                />
            </div>

            <Card>
                <div style={{ marginBottom: 16 }}>
                    <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
                        添加角色
                    </Button>
                </div>
                
                <ResponsiveTable
                    columns={columns}
                    dataSource={roles}
                    rowKey="id"
                    loading={loading}
                    actions={getActions}
                    cardTitle={getCardTitle}
                    cardSubtitle={getCardSubtitle}
                    pagination={false}
                />
            </Card>

            <Modal
                title={editingRole ? '编辑角色' : '添加角色'}
                open={isModalVisible}
                onOk={handleModalOk}
                onCancel={handleModalCancel}
                width={700}
            >
                <Form
                    form={form}
                    layout="vertical"
                >
                    <Form.Item
                        label="角色名称"
                        name="name"
                        rules={[
                            { required: true, message: '请输入角色名称' },
                            { pattern: /^[a-z0-9_]+$/, message: '只能使用小写字母、数字和下划线' }
                        ]}
                    >
                        <Input 
                            placeholder="请输入角色名称，如：editor" 
                            disabled={isSystemRole || !!editingRole}
                        />
                    </Form.Item>

                    <Form.Item
                        label="描述"
                        name="description"
                        rules={[{ required: true, message: '请输入角色描述' }]}
                    >
                        <TextArea 
                            rows={2} 
                            placeholder="请输入角色描述"
                            disabled={isSystemRole}
                        />
                    </Form.Item>

                    <Form.Item
                        label="权限配置"
                        required
                    >
                        {isSystemRole ? (
                            <Tag color="red">超级管理员拥有所有权限，不可修改</Tag>
                        ) : (
                            <Space direction="vertical" style={{ width: '100%' }}>
                                {PERMISSION_CATEGORIES.map(cat => (
                                    <Card 
                                        key={cat.category} 
                                        size="small" 
                                        title={cat.category}
                                        style={{ marginBottom: 8 }}
                                    >
                                        <Checkbox.Group
                                            value={selectedPermissions}
                                            onChange={handlePermissionChange}
                                        >
                                            <Space wrap>
                                                {cat.permissions.map(perm => (
                                                    <Checkbox key={perm.key} value={perm.key}>
                                                        {perm.name}
                                                    </Checkbox>
                                                ))}
                                            </Space>
                                        </Checkbox.Group>
                                    </Card>
                                ))}
                            </Space>
                        )}
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    );
};

export default RoleManagement;
