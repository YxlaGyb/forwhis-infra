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
                message.error('获取用户列表失败');
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
        form.setFieldsValue(user);
    };

    const handleDeleteUser = async (userId: string) => {
        try {
            await userApi.deleteUser(parseInt(userId));
            setUsers(users.filter(user => user.id !== userId));
            message.success('用户删除成功');
        } catch (error) {
            console.error('删除用户失败:', error);
            message.error('删除用户失败');
        }
    };

    const handleModalOk = async () => {
        try {
            const values = await form.validateFields();
            
            if (editingUser) {
                const response = await userApi.updateUser(parseInt(editingUser.id), {
                    name: values.name,
                    email: values.email,
                    role: values.role,
                    status: values.status
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
                    role: values.role,
                    status: values.status
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
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
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
            dataIndex: 'role',
            key: 'role',
            render: (role: string) => (
                <Tag color={role === 'admin' ? 'red' : 'blue'}>
                    {role === 'admin' ? '管理员' : '普通用户'}
                </Tag>
            ),
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            render: (status: string) => (
                <Tag color={status === 'active' ? 'green' : 'red'}>
                    {status === 'active' ? '活跃' : '禁用'}
                </Tag>
            ),
        },
        {
            title: '创建时间',
            dataIndex: 'createTime',
            key: 'createTime',
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
                        onConfirm={() => handleDeleteUser(record.id)}
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
                    rowKey="id"
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
                        name="role"
                        rules={[{ required: true, message: '请选择用户角色' }]}
                    >
                        <Select placeholder="请选择用户角色">
                            <Option value="user">普通用户</Option>
                            <Option value="admin">管理员</Option>
                        </Select>
                    </Form.Item>

                    <Form.Item
                        label="状态"
                        name="status"
                        rules={[{ required: true, message: '请选择用户状态' }]}
                    >
                        <Select placeholder="请选择用户状态">
                            <Option value="active">活跃</Option>
                            <Option value="inactive">禁用</Option>
                        </Select>
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    );
};

export default UserManagement;
