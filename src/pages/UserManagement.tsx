import React, { useState, useEffect, useRef } from 'react';
import {
    Table,
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
    Popconfirm
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
import type { User } from '../types';

const { Title } = Typography;
const { Option } = Select;

const formatTimestamp = (timestamp: number): string => {
    if (!timestamp) return '-';
    const date = new Date(timestamp * 1000);
    return date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    });
};

const getRoleColor = (role: string): string => {
    const colorMap: Record<string, string> = {
        'superadmin': 'red',
        'admin': 'orange',
        'editor': 'blue',
        'user': 'green',
        'visitor': 'default'
    };
    return colorMap[role] || 'default';
};

const getRoleName = (role: string): string => {
    const nameMap: Record<string, string> = {
        'superadmin': '超级管理员',
        'admin': '管理员',
        'editor': '编辑',
        'user': '用户',
        'visitor': '访客'
    };
    return nameMap[role] || role;
};

const UserManagement: React.FC = () => {
    const [users, setUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState(false);
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [editingUser, setEditingUser] = useState<User | null>(null);
    const [form] = Form.useForm();
    
    const abortControllerRef = useRef<AbortController | null>(null);

    const { filteredData: filteredUsers, handleSearch } = useSearchFilter<User>({
        data: users,
        searchKeys: ['name', 'email']
    });

    useEffect(() => {
        fetchUsers();
        return () => {
            if (abortControllerRef.current) {
                abortControllerRef.current.abort();
            }
        };
    }, []);

    const fetchUsers = async () => {
        abortControllerRef.current = new AbortController();
        setLoading(true);
        try {
            const response = await userApi.getUsers();
            
            if (abortControllerRef.current?.signal.aborted) return;
            
            if (response.data.code === 200 && response.data.data) {
                setUsers(response.data.data);
            } else {
                message.error(response.data.msg || '获取用户列表失败');
            }
        } catch (error) {
            if (abortControllerRef.current?.signal.aborted) return;
            console.error('获取用户列表失败:', error);
            message.error('获取用户列表失败');
        } finally {
            if (!abortControllerRef.current?.signal.aborted) {
                setLoading(false);
            }
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
            roles: user.roles && user.roles.length > 0 ? user.roles[0] : 'user'
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
                    role: values.roles
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

    const columns = [
        {
            title: 'UID',
            dataIndex: 'uid',
            key: 'uid',
            width: 80,
        },
        {
            title: '姓名',
            dataIndex: 'name',
            key: 'name',
            render: (text: string) => (
                <Space>
                    <UserOutlined />
                    {text}
                </Space>
            ),
        },
        {
            title: '邮箱',
            dataIndex: 'email',
            key: 'email',
        },
        {
            title: '角色',
            dataIndex: 'roles',
            key: 'roles',
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
            render: (timestamp: number) => formatTimestamp(timestamp),
        },
        {
            title: '操作',
            key: 'action',
            render: (_record: unknown, record: User) => (
                <Space size="middle">
                    <Button
                        type="link"
                        icon={<EditOutlined />}
                        onClick={() => handleEditUser(record)}
                    >
                        编辑
                    </Button>
                    <Popconfirm
                        title="确定删除这个用户吗？"
                        onConfirm={() => handleDeleteUser(record.uid)}
                        okText="确定"
                        cancelText="取消"
                    >
                        <Button type="link" danger icon={<DeleteOutlined />}>
                            删除
                        </Button>
                    </Popconfirm>
                </Space>
            ),
        },
    ];

    return (
        <div>
            <Title level={2} style={{ marginBottom: 24 }}>
                <UserOutlined /> 用户管理
            </Title>
            
            <Card>
                <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
                    <Input
                        placeholder="搜索用户姓名或邮箱"
                        prefix={<SearchOutlined />}
                        style={{ width: 300 }}
                        onChange={(e) => handleSearch(e.target.value)}
                        allowClear
                    />
                    <Button type="primary" icon={<PlusOutlined />} onClick={handleAddUser}>
                        添加用户
                    </Button>
                </div>
                
                <Table
                    columns={columns}
                    dataSource={filteredUsers}
                    rowKey="uid"
                    loading={loading}
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
                        label="姓名"
                        name="name"
                        rules={[{ required: true, message: '请输入用户姓名' }]}
                    >
                        <Input placeholder="请输入用户姓名" />
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
                            rules={[{ required: !editingUser, message: '请输入密码' }]}
                        >
                            <Input.Password placeholder="请输入密码" />
                        </Form.Item>
                    )}

                    <Form.Item
                        label="角色"
                        name="roles"
                        rules={[{ required: true, message: '请选择用户角色' }]}
                    >
                        <Select placeholder="请选择用户角色">
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
