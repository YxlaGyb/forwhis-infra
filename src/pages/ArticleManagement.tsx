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
    Popconfirm,
    Descriptions
} from 'antd';
import {
    SearchOutlined,
    PlusOutlined,
    EditOutlined,
    DeleteOutlined,
    EyeOutlined,
    FileTextOutlined
} from '@ant-design/icons';
import { contentApi } from '../api';
import { useSearchFilter } from '../hooks/useSearchFilter';
import type { Content } from '../types';

const { Title } = Typography;
const { Option } = Select;
const { TextArea } = Input;

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

const getStatusColor = (status: string): string => {
    const colorMap: Record<string, string> = {
        'publish': 'green',
        'draft': 'orange',
        'private': 'blue',
        'trash': 'red'
    };
    return colorMap[status] || 'default';
};

const getStatusName = (status: string): string => {
    const nameMap: Record<string, string> = {
        'publish': '已发布',
        'draft': '草稿',
        'private': '私密',
        'trash': '回收站'
    };
    return nameMap[status] || status;
};

const getTypeName = (type: string): string => {
    const nameMap: Record<string, string> = {
        'post': '文章',
        'page': '页面'
    };
    return nameMap[type] || type;
};

const ArticleManagement: React.FC = () => {
    const [articles, setArticles] = useState<Content[]>([]);
    const [loading, setLoading] = useState(false);
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [isViewModalVisible, setIsViewModalVisible] = useState(false);
    const [editingArticle, setEditingArticle] = useState<Content | null>(null);
    const [viewingArticle, setViewingArticle] = useState<Content | null>(null);
    const [form] = Form.useForm();
    const [statusFilter, setStatusFilter] = useState<string>('');
    const [typeFilter, setTypeFilter] = useState<string>('');

    const abortControllerRef = useRef<AbortController | null>(null);

    const { filteredData: filteredArticles, handleSearch } = useSearchFilter<Content>({
        data: articles,
        searchKeys: ['title', 'slug']
    });

    useEffect(() => {
        fetchArticles();
        return () => {
            if (abortControllerRef.current) {
                abortControllerRef.current.abort();
            }
        };
    }, []);

    const fetchArticles = async () => {
        abortControllerRef.current = new AbortController();
        setLoading(true);
        try {
            const response = await contentApi.getContents({ limit: 100 });

            if (abortControllerRef.current?.signal.aborted) return;

            if (response.data.code === 200 && response.data.data) {
                const items = response.data.data.items || response.data.data;
                setArticles(items);
            } else {
                message.error(response.data.msg || '获取文章列表失败');
            }
        } catch {
            if (abortControllerRef.current?.signal.aborted) return;
            message.error('获取文章列表失败');
        } finally {
            setLoading(false);
        }
    };

    const handleCreate = () => {
        setEditingArticle(null);
        form.resetFields();
        setIsModalVisible(true);
    };

    const handleEdit = (record: Content) => {
        setEditingArticle(record);
        form.setFieldsValue({
            title: record.title,
            slug: record.slug,
            text: record.text,
            status: record.status,
            type: record.type
        });
        setIsModalVisible(true);
    };

    const handleView = (record: Content) => {
        setViewingArticle(record);
        setIsViewModalVisible(true);
    };

    const handleDelete = async (cid: number) => {
        try {
            const response = await contentApi.deleteContent(cid);
            if (response.data.code === 200) {
                message.success('删除成功');
                fetchArticles();
            } else {
                message.error(response.data.msg || '删除失败');
            }
        } catch {
            message.error('删除失败');
        }
    };

    const handlePublish = async (cid: number) => {
        try {
            const response = await contentApi.publishContent(cid);
            if (response.data.code === 200) {
                message.success('发布成功');
                fetchArticles();
            } else {
                message.error(response.data.msg || '发布失败');
            }
        } catch {
            message.error('发布失败');
        }
    };

    const handleModalOk = async () => {
        try {
            const values = await form.validateFields();

            if (editingArticle) {
                const response = await contentApi.updateContent(editingArticle.cid, values);
                if (response.data.code === 200) {
                    message.success('更新成功');
                    setIsModalVisible(false);
                    fetchArticles();
                } else {
                    message.error(response.data.msg || '更新失败');
                }
            } else {
                const response = await contentApi.createContent(values);
                if (response.data.code === 200) {
                    message.success('创建成功');
                    setIsModalVisible(false);
                    fetchArticles();
                } else {
                    message.error(response.data.msg || '创建失败');
                }
            }
        } catch {
            message.error('操作失败');
        }
    };

    const handleModalCancel = () => {
        setIsModalVisible(false);
        form.resetFields();
    };

    const handleViewModalCancel = () => {
        setIsViewModalVisible(false);
        setViewingArticle(null);
    };

    const filteredByStatusAndType = filteredArticles.filter(article => {
        const statusMatch = !statusFilter || article.status === statusFilter;
        const typeMatch = !typeFilter || article.type === typeFilter;
        return statusMatch && typeMatch;
    });

    const columns = [
        {
            title: 'ID',
            dataIndex: 'cid',
            key: 'cid',
            width: 80,
        },
        {
            title: '标题',
            dataIndex: 'title',
            key: 'title',
            width: 80,
            ellipsis: true,
        },
        {
            title: '别名',
            dataIndex: 'slug',
            key: 'slug',
            width: 80,
            ellipsis: true,
        },
        {
            title: '类型',
            dataIndex: 'type',
            key: 'type',
            width: 80,
            render: (type: string) => getTypeName(type),
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            width: 80,
            render: (status: string) => (
                <Tag color={getStatusColor(status)}>
                    {getStatusName(status)}
                </Tag>
            ),
        },
        {
            title: '作者',
            dataIndex: 'author',
            key: 'author',
            width: 80,
        },
        {
            title: '评论数',
            dataIndex: 'commentsNum',
            key: 'commentsNum',
            width: 100,
        },
        {
            title: '创建时间',
            dataIndex: 'created',
            key: 'created',
            width: 120,
            render: (created: number) => formatTimestamp(created),
        },
        {
            title: '操作',
            key: 'action',
            width: 80,
            render: (_: unknown, record: Content) => (
                <Space size="small">
                    <Button
                        type="text"
                        size="small"
                        icon={<EyeOutlined />}
                        onClick={() => handleView(record)}
                    >
                        查看
                    </Button>
                    <Button
                        type="text"
                        size="small"
                        icon={<EditOutlined />}
                        onClick={() => handleEdit(record)}
                    >
                        编辑
                    </Button>
                    {record.status !== 'publish' && (
                        <Button
                            type="text"
                            size="small"
                            onClick={() => handlePublish(record.cid)}
                        >
                            发布
                        </Button>
                    )}
                    <Popconfirm
                        title="确定要删除这篇文章吗？"
                        onConfirm={() => handleDelete(record.cid)}
                        okText="确定"
                        cancelText="取消"
                    >
                        <Button
                            type="text"
                            danger
                            size="small"
                            icon={<DeleteOutlined />}
                        >
                            删除
                        </Button>
                    </Popconfirm>
                </Space>
            ),
        },
    ];

    return (
        <div>
            <Card>
                <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 16 }}>
                    <Title level={4} style={{ margin: 0 }}>
                        <FileTextOutlined style={{ marginRight: 8 }} />
                        文章管理
                    </Title>
                    <Button
                        type="primary"
                        icon={<PlusOutlined />}
                        onClick={handleCreate}
                    >
                        新建文章
                    </Button>
                </div>

                <Space style={{ marginBottom: 16 }} wrap>
                    <Input
                        placeholder="搜索标题或别名"
                        prefix={<SearchOutlined />}
                        onChange={(e) => handleSearch(e.target.value)}
                        style={{ width: 250 }}
                        allowClear
                    />
                    <Select
                        placeholder="筛选状态"
                        value={statusFilter || undefined}
                        onChange={(value) => setStatusFilter(value)}
                        style={{ width: 120 }}
                        allowClear
                    >
                        <Option value="publish">已发布</Option>
                        <Option value="draft">草稿</Option>
                        <Option value="private">私密</Option>
                        <Option value="trash">回收站</Option>
                    </Select>
                    <Select
                        placeholder="筛选类型"
                        value={typeFilter || undefined}
                        onChange={(value) => setTypeFilter(value)}
                        style={{ width: 120 }}
                        allowClear
                    >
                        <Option value="post">文章</Option>
                        <Option value="page">页面</Option>
                    </Select>
                </Space>

                <Table
                    columns={columns}
                    dataSource={filteredByStatusAndType}
                    rowKey="cid"
                    loading={loading}
                    pagination={{
                        pageSize: 10,
                        showSizeChanger: true,
                        showTotal: (total) => `共 ${total} 条`,
                    }}
                />
            </Card>

            <Modal
                title={editingArticle ? '编辑文章' : '新建文章'}
                open={isModalVisible}
                onOk={handleModalOk}
                onCancel={handleModalCancel}
                width={700}
                destroyOnClose
            >
                <Form
                    form={form}
                    layout="vertical"
                    initialValues={{ status: 'draft', type: 'post' }}
                >
                    <Form.Item
                        name="title"
                        label="标题"
                        rules={[{ required: true, message: '请输入文章标题' }]}
                    >
                        <Input placeholder="请输入文章标题" />
                    </Form.Item>

                    <Form.Item
                        name="slug"
                        label="别名"
                        rules={[{ required: true, message: '请输入文章别名' }]}
                    >
                        <Input placeholder="请输入文章别名，用于URL" />
                    </Form.Item>

                    <Form.Item
                        name="text"
                        label="内容"
                        rules={[{ required: true, message: '请输入文章内容' }]}
                    >
                        <TextArea
                            rows={10}
                            placeholder="请输入文章内容，支持Markdown格式"
                        />
                    </Form.Item>

                    <Form.Item
                        name="type"
                        label="类型"
                        rules={[{ required: true, message: '请选择文章类型' }]}
                    >
                        <Select placeholder="请选择文章类型">
                            <Option value="post">文章</Option>
                            <Option value="page">页面</Option>
                        </Select>
                    </Form.Item>

                    <Form.Item
                        name="status"
                        label="状态"
                        rules={[{ required: true, message: '请选择文章状态' }]}
                    >
                        <Select placeholder="请选择文章状态">
                            <Option value="publish">已发布</Option>
                            <Option value="draft">草稿</Option>
                            <Option value="private">私密</Option>
                        </Select>
                    </Form.Item>
                </Form>
            </Modal>

            <Modal
                title="文章详情"
                open={isViewModalVisible}
                onCancel={handleViewModalCancel}
                footer={[
                    <Button key="close" onClick={handleViewModalCancel}>
                        关闭
                    </Button>
                ]}
                width={700}
            >
                {viewingArticle && (
                    <Descriptions bordered column={1}>
                        <Descriptions.Item label="ID">{viewingArticle.cid}</Descriptions.Item>
                        <Descriptions.Item label="标题">{viewingArticle.title}</Descriptions.Item>
                        <Descriptions.Item label="别名">{viewingArticle.slug}</Descriptions.Item>
                        <Descriptions.Item label="类型">{getTypeName(viewingArticle.type)}</Descriptions.Item>
                        <Descriptions.Item label="状态">
                            <Tag color={getStatusColor(viewingArticle.status)}>
                                {getStatusName(viewingArticle.status)}
                            </Tag>
                        </Descriptions.Item>
                        <Descriptions.Item label="作者">{viewingArticle.author}</Descriptions.Item>
                        <Descriptions.Item label="评论数">{viewingArticle.commentsNum}</Descriptions.Item>
                        <Descriptions.Item label="创建时间">{formatTimestamp(viewingArticle.created)}</Descriptions.Item>
                        <Descriptions.Item label="修改时间">{formatTimestamp(viewingArticle.modified)}</Descriptions.Item>
                        <Descriptions.Item label="内容">
                            <div style={{ maxHeight: 300, overflow: 'auto', whiteSpace: 'pre-wrap' }}>
                                {viewingArticle.text}
                            </div>
                        </Descriptions.Item>
                    </Descriptions>
                )}
            </Modal>
        </div>
    );
};

export default ArticleManagement;
