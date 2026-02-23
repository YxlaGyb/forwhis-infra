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
    Descriptions,
    Avatar
} from 'antd';
import {
    SearchOutlined,
    DeleteOutlined,
    EyeOutlined,
    CommentOutlined,
    CheckOutlined,
    CloseOutlined,
    EditOutlined,
    UserOutlined
} from '@ant-design/icons';
import { commentApi } from '../api';
import { useSearchFilter } from '../hooks/useSearchFilter';
import { formatTimestamp } from '../utils/format';
import { getStatusColor, getStatusName } from '../utils/mapping';
import { ResponsiveTable } from '../components/ResponsiveTable';
import type { ResponsiveColumn, ActionItem } from '../components/ResponsiveTable';
import type { Comment } from '../types';

const { Title } = Typography;
const { Option } = Select;
const { TextArea } = Input;

const CommentManagement: React.FC = () => {
    const [comments, setComments] = useState<Comment[]>([]);
    const [loading, setLoading] = useState(false);
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [isViewModalVisible, setIsViewModalVisible] = useState(false);
    const [editingComment, setEditingComment] = useState<Comment | null>(null);
    const [viewingComment, setViewingComment] = useState<Comment | null>(null);
    const [form] = Form.useForm();
    const [statusFilter, setStatusFilter] = useState<string>('');

    const { filteredData: filteredComments, handleSearch } = useSearchFilter<Comment>({
        data: comments,
        searchKeys: ['text', 'author', 'mail']
    });

    useEffect(() => {
        fetchComments();
    }, []);

    const fetchComments = async () => {
        setLoading(true);
        try {
            const response = await commentApi.getComments({ limit: 100 });

            if (response.data.code === 200 && response.data.data) {
                const items = response.data.data.items || response.data.data;
                setComments(items);
            } else {
                message.error(response.data.msg || '获取评论列表失败');
            }
        } catch (error) {
            console.error('获取评论列表失败:', error);
            message.error('获取评论列表失败');
        } finally {
            setLoading(false);
        }
    };

    const handleEdit = (record: Comment) => {
        setEditingComment(record);
        form.setFieldsValue({
            text: record.text,
            status: record.status
        });
        setIsModalVisible(true);
    };

    const handleView = (record: Comment) => {
        setViewingComment(record);
        setIsViewModalVisible(true);
    };

    const handleDelete = async (coid: number) => {
        try {
            const response = await commentApi.deleteComment(coid);
            if (response.data.code === 200) {
                message.success('删除成功');
                fetchComments();
            } else {
                message.error(response.data.msg || '删除失败');
            }
        } catch (error) {
            console.error('删除失败:', error);
            message.error('删除失败');
        }
    };

    const handleApprove = async (coid: number) => {
        try {
            const response = await commentApi.approveComment(coid);
            if (response.data.code === 200) {
                message.success('审核通过');
                fetchComments();
            } else {
                message.error(response.data.msg || '审核失败');
            }
        } catch (error) {
            console.error('审核失败:', error);
            message.error('审核失败');
        }
    };

    const handleUpdateStatus = async (coid: number, status: string) => {
        try {
            const response = await commentApi.updateComment(coid, { status });
            if (response.data.code === 200) {
                message.success('状态更新成功');
                fetchComments();
            } else {
                message.error(response.data.msg || '状态更新失败');
            }
        } catch (error) {
            console.error('状态更新失败:', error);
            message.error('状态更新失败');
        }
    };

    const handleModalOk = async () => {
        try {
            const values = await form.validateFields();

            if (editingComment) {
                const response = await commentApi.updateComment(editingComment.coid, values);
                if (response.data.code === 200) {
                    message.success('更新成功');
                    setIsModalVisible(false);
                    fetchComments();
                } else {
                    message.error(response.data.msg || '更新失败');
                }
            }
        } catch (error) {
            console.error('操作失败:', error);
            message.error('操作失败');
        }
    };

    const handleModalCancel = () => {
        setIsModalVisible(false);
        form.resetFields();
        setEditingComment(null);
    };

    const handleViewModalCancel = () => {
        setIsViewModalVisible(false);
        setViewingComment(null);
    };

    const filteredByStatus = filteredComments.filter(comment => {
        return !statusFilter || comment.status === statusFilter;
    });

    // 表格列配置
    const columns: ResponsiveColumn<Comment>[] = [
        {
            title: 'ID',
            dataIndex: 'coid',
            key: 'coid',
            width: 80,
            priority: 'low',
        },
        {
            title: '评论内容',
            dataIndex: 'text',
            key: 'text',
            width: 250,
            ellipsis: true,
            priority: 'high',
            render: (text: string) => (
                <span style={{ maxWidth: 300, display: 'inline-block', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                    {text}
                </span>
            ),
        },
        {
            title: '评论者',
            dataIndex: 'author',
            key: 'author',
            width: 120,
            priority: 'low',
        },
        {
            title: '邮箱',
            dataIndex: 'mail',
            key: 'mail',
            ellipsis: true,
            width: 180,
            priority: 'low',
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 100,
            priority: 'high',
            render: (status: string) => (
                <Tag color={getStatusColor(status)}>
                    {getStatusName(status)}
                </Tag>
            ),
        },
        {
            title: '文章ID',
            dataIndex: 'cid',
            key: 'cid',
            width: 80,
            priority: 'low',
        },
        {
            title: '发布时间',
            dataIndex: 'created',
            key: 'created',
            width: 160,
            priority: 'medium',
            render: (created: number) => formatTimestamp(created),
        },
    ];

    // 操作按钮配置
    const getActions = (record: Comment): ActionItem[] => [
        {
            key: 'view',
            text: '查看',
            icon: <EyeOutlined />,
            onClick: () => handleView(record),
        },
        {
            key: 'approve',
            text: '通过',
            icon: <CheckOutlined />,
            hidden: record.status === 'approved',
            onClick: () => handleApprove(record.coid),
        },
        {
            key: 'spam',
            text: '标为垃圾',
            icon: <CloseOutlined />,
            hidden: record.status === 'spam',
            onClick: () => handleUpdateStatus(record.coid, 'spam'),
        },
        {
            key: 'edit',
            text: '编辑',
            icon: <EditOutlined />,
            onClick: () => handleEdit(record),
        },
        {
            key: 'delete',
            text: '删除',
            icon: <DeleteOutlined />,
            danger: true,
            onClick: () => {
                Modal.confirm({
                    title: '确定要删除这条评论吗？',
                    content: '删除后无法恢复',
                    okText: '确定',
                    cancelText: '取消',
                    onOk: () => handleDelete(record.coid),
                });
            },
        },
    ];

    // 卡片标题
    const getCardTitle = (record: Comment) => (
        <Space>
            <Avatar size="small" icon={<UserOutlined />} />
            {record.author}
        </Space>
    );

    // 卡片副标题
    const getCardSubtitle = (record: Comment) => record.mail;

    return (
        <div>
            <Card>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
                    <Title level={4} style={{ margin: 0 }}>
                        <CommentOutlined style={{ marginRight: 8 }} />
                        评论管理
                    </Title>
                </div>

                <Space style={{ marginBottom: 16 }} wrap>
                    <Input
                        placeholder="搜索评论内容、评论者或邮箱"
                        prefix={<SearchOutlined />}
                        onChange={(e) => handleSearch(e.target.value)}
                        style={{ width: 300 }}
                        allowClear
                    />
                    <Select
                        placeholder="筛选状态"
                        value={statusFilter || undefined}
                        onChange={(value) => setStatusFilter(value)}
                        style={{ width: 120 }}
                        allowClear
                    >
                        <Option value="approved">已通过</Option>
                        <Option value="pending">待审核</Option>
                        <Option value="spam">垃圾评论</Option>
                        <Option value="trash">回收站</Option>
                    </Select>
                </Space>

                <ResponsiveTable
                    columns={columns}
                    dataSource={filteredByStatus}
                    rowKey="coid"
                    loading={loading}
                    actions={getActions}
                    cardTitle={getCardTitle}
                    cardSubtitle={getCardSubtitle}
                    pagination={{
                        pageSize: 10,
                        showSizeChanger: true,
                        showTotal: (total) => `共 ${total} 条`,
                    }}
                />
            </Card>

            <Modal
                title="编辑评论"
                open={isModalVisible}
                onOk={handleModalOk}
                onCancel={handleModalCancel}
                width={600}
                destroyOnClose
            >
                <Form
                    form={form}
                    layout="vertical"
                >
                    <Form.Item
                        name="text"
                        label="评论内容"
                        rules={[{ required: true, message: '请输入评论内容' }]}
                    >
                        <TextArea
                            rows={6}
                            placeholder="请输入评论内容"
                        />
                    </Form.Item>

                    <Form.Item
                        name="status"
                        label="状态"
                        rules={[{ required: true, message: '请选择评论状态' }]}
                    >
                        <Select placeholder="请选择评论状态">
                            <Option value="approved">已通过</Option>
                            <Option value="pending">待审核</Option>
                            <Option value="spam">垃圾评论</Option>
                            <Option value="trash">回收站</Option>
                        </Select>
                    </Form.Item>
                </Form>
            </Modal>

            <Modal
                title="评论详情"
                open={isViewModalVisible}
                onCancel={handleViewModalCancel}
                footer={[
                    <Button key="close" onClick={handleViewModalCancel}>
                        关闭
                    </Button>
                ]}
                width={600}
            >
                {viewingComment && (
                    <Descriptions bordered column={1}>
                        <Descriptions.Item label="ID">{viewingComment.coid}</Descriptions.Item>
                        <Descriptions.Item label="评论内容">
                            <div style={{ maxHeight: 200, overflow: 'auto', whiteSpace: 'pre-wrap' }}>
                                {viewingComment.text}
                            </div>
                        </Descriptions.Item>
                        <Descriptions.Item label="评论者">{viewingComment.author}</Descriptions.Item>
                        <Descriptions.Item label="邮箱">{viewingComment.mail}</Descriptions.Item>
                        <Descriptions.Item label="网站">{viewingComment.url || '-'}</Descriptions.Item>
                        <Descriptions.Item label="IP地址">{viewingComment.ip || '-'}</Descriptions.Item>
                        <Descriptions.Item label="状态">
                            <Tag color={getStatusColor(viewingComment.status)}>
                                {getStatusName(viewingComment.status)}
                            </Tag>
                        </Descriptions.Item>
                        <Descriptions.Item label="文章ID">{viewingComment.cid}</Descriptions.Item>
                        <Descriptions.Item label="父评论ID">{viewingComment.parent || '-'}</Descriptions.Item>
                        <Descriptions.Item label="发布时间">{formatTimestamp(viewingComment.created)}</Descriptions.Item>
                    </Descriptions>
                )}
            </Modal>
        </div>
    );
};

export default CommentManagement;
