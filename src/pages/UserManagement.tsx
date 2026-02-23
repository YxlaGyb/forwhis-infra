import React, { useState, useEffect } from 'react';
import {
    Card,
    Input,
    Button,
    Space,
    Typography,
    Tag,
    Modal,
    Form,
    Select,
    message,
    Avatar
} from 'antd';
import {
    SearchOutlined,
    PlusOutlined,
    EditOutlined,
    DeleteOutlined,
    UserOutlined
} from '@ant-design/icons';
import { userApi } from '../api';
import { useSearchFilter } from '../hooks/useSearchFilter';
import { formatTimestamp } from '../utils/format';
import { getRoleColor, getRoleName } from '../utils/mapping';
import { ResponsiveTable } from '../components/ResponsiveTable';
import type { ResponsiveColumn, ActionItem } from '../components/ResponsiveTable';
import type { User } from '../types';

const { Title } = Typography;
const { Option } = Select;

const UserManagement: React.FC = () => {
    const [users, setUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState(false);
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [editingUser, setEditingUser] = useState<User | null>(null);
    const [form] = Form.useForm();

    const { filteredData: filteredUsers, handleSearch } = useSearchFilter<User>({
        data: users,
        searchKeys: ['name', 'email']
    });

    useEffect(() => {
        fetchUsers();
    }, []);

    const fetchUsers = async () => {
        setLoading(true);
        try {
            const response = await userApi.getUsers();
            
            if (response.data.code === 200 && response.data.data) {
                setUsers(response.data.data);
            } else {
                message.error(response.data.msg || '获取用户列表失败');
            }
        } catch (error) {
            console.error('获取用户列表失败:', error);
            message.error('获取用户列表失败');
        } finally {
            setLoading(false);
        }
    };

    const handleAddUser = () => {
        setEditingUser(null);
        setIsModalVisible(true);
        form.resetFields();
    };

    const handleEditUser = (user: User) => {
        setEditingUser(user);
        setIsModalVisible(true);
        form.setFieldsValue({
            name: user.name,
            email: user.email,
            roles: user.roles && user.roles.length > 0 ? user.roles : ['user']
        });
    };

    const handleDeleteUser = async (uid: number) => {
        try {
            const response = await userApi.deleteUser(uid);
            if (response.data.code === 200) {
                setUsers(users.filter(user => user.uid !== uid));
                message.success('用户删除成功');
            } else {
                message.error(response.data.msg || '用户删除失败');
            }
        } catch (error) {
            console.error('删除用户失败:', error);
            message.error('删除用户失败');
        }
    };

    const handleModalOk = async () => {
        try {
            const values = await form.validateFields();
            
            if (editingUser) {
                const response = await userApi.updateUser(editingUser.uid, {
                    name: values.name,
                    email: values.email,
                    roles: values.roles
                });
                
                if (response.data.code === 200) {
                    message.success('用户更新成功');
                    fetchUsers();
                } else {
                    message.error(response.data.msg || '用户更新失败');
                }
            } else {
                if (!values.password) {
                    message.error('请输入密码');
                    return;
                }
                
                const response = await userApi.createUser({
                    name: values.name,
                    email: values.email,
                    password: values.password,
                    role: values.roles
                });
                
                if (response.data.code === 200) {
                    message.success('用户添加成功');
                    fetchUsers();
                } else {
                    message.error(response.data.msg || '用户添加失败');
                }
            }
            
            setIsModalVisible(false);
            form.resetFields();
        } catch (error) {
            console.error('表单验证失败:', error);
            message.error('操作失败');
        }
    };

    const handleModalCancel = () => {
        setIsModalVisible(false);
        form.resetFields();
        setEditingUser(null);
    };

    // 表格列配置
    const columns: ResponsiveColumn<User>[] = [
        {
            title: 'UID',
            dataIndex: 'uid',
            key: 'uid',
            width: 80,
            priority: 'low',
        },
        {
            title: '昵称',
            dataIndex: 'name',
            key: 'name',
            width: 150,
            priority: 'low',
            render: (text: string) => (
                <Space>
                    <Avatar size="small" icon={<UserOutlined />} />
                    {text}
                </Space>
            ),
        },
        {
            title: '邮箱',
            dataIndex: 'email',
            key: 'email',
            width: 200,
            priority: 'low',
        },
        {
            title: '角色',
            dataIndex: 'roles',
            key: 'roles',
            width: 150,
            priority: 'high',
            render: (roles: string[]) => (
                <Space>
                    {roles && roles.map((role, index) => (
                        <Tag key={index} color={getRoleColor(role)}>
                            {getRoleName(role)}
                        </Tag>
                    ))}
                </Space>
            ),
        },
        {
            title: '创建时间',
            dataIndex: 'createdAt',
            key: 'createdAt',
            width: 160,
            priority: 'medium',
            render: (timestamp: number) => formatTimestamp(timestamp),
        },
        {
            title: '最后登录',
            dataIndex: 'lastLogin',
            key: 'lastLogin',
            width: 160,
            priority: 'medium',
            render: (timestamp: number) => timestamp ? formatTimestamp(timestamp) : '-',
        },
        {
            title: '登录IP',
            dataIndex: 'lastip',
            key: 'lastip',
            width: 140,
            priority: 'medium',
        },
    ];

    // 操作按钮配置
    const getActions = (record: User): ActionItem[] => [
        {
            key: 'edit',
            text: '编辑',
            icon: <EditOutlined />,
            onClick: () => handleEditUser(record),
        },
        {
            key: 'delete',
            text: '删除',
            icon: <DeleteOutlined />,
            danger: true,
            onClick: () => {
                Modal.confirm({
                    title: '确定删除这个用户吗？',
                    content: '删除后无法恢复',
                    okText: '确定',
                    cancelText: '取消',
                    onOk: () => handleDeleteUser(record.uid),
                });
            },
        },
    ];

    // 卡片标题
    const getCardTitle = (record: User) => (
        <Space>
            <Avatar icon={<UserOutlined />} />
            {record.name}
        </Space>
    );

    // 卡片副标题
    const getCardSubtitle = (record: User) => record.email;

    return (
        <div>
            <Title level={4} style={{ marginBottom: 16 }}>
                <UserOutlined /> 用户管理
            </Title>
            
            <Card>
                <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
                    <Input
                        placeholder="搜索用户昵称或邮箱"
                        prefix={<SearchOutlined />}
                        style={{ width: 300 }}
                        onChange={(e) => handleSearch(e.target.value)}
                        allowClear
                    />
                    <Button type="primary" icon={<PlusOutlined />} onClick={handleAddUser}>
                        添加用户
                    </Button>
                </div>
                
                <ResponsiveTable
                    columns={columns}
                    dataSource={filteredUsers}
                    rowKey="uid"
                    loading={loading}
                    actions={getActions}
                    cardTitle={getCardTitle}
                    cardSubtitle={getCardSubtitle}
                    pagination={{
                        showSizeChanger: true,
                        showQuickJumper: true,
                        showTotal: (total) => `共 ${total} 条记录`,
                    }}
                />
            </Card>

            <Modal
                title={editingUser ? '编辑用户' : '添加用户'}
                open={isModalVisible}
                onOk={handleModalOk}
                onCancel={handleModalCancel}
                width={600}
            >
                <Form
                    form={form}
                    layout="vertical"
                    name="userForm"
                >
                    <Form.Item
                        label="昵称"
                        name="name"
                        rules={[{ required: true, message: '请输入用户昵称' }]}
                    >
                        <Input placeholder="请输入用户昵称" />
                    </Form.Item>

                    <Form.Item
                        label="邮箱"
                        name="email"
                        rules={[
                            { required: true, message: '请输入邮箱地址' },
                            { type: 'email', message: '请输入有效的邮箱地址' }
                        ]}
                    >
                        <Input placeholder="请输入邮箱地址" />
                    </Form.Item>

                    {!editingUser && (
                        <Form.Item
                            label="密码"
                            name="password"
                            rules={[
                                { required: true, message: '请输入密码' },
                                { min: 8, message: '密码至少8个字符' },
                                { pattern: /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/, message: '密码需包含大小写字母和数字' }
                            ]}
                        >
                            <Input.Password placeholder="请输入密码（至少8位，包含大小写字母和数字）" />
                        </Form.Item>
                    )}

                    <Form.Item
                        label="角色"
                        name="roles"
                        rules={[{ required: true, message: '请选择用户角色' }]}
                    >
                        <Select
                            mode="multiple"
                            placeholder="请选择用户角色"
                            allowClear
                        >
                            <Option value="user">普通用户</Option>
                            <Option value="editor">编辑</Option>
                            <Option value="admin">管理员</Option>
                            <Option value="superadmin">超级管理员</Option>
                        </Select>
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    );
};

export default UserManagement;
