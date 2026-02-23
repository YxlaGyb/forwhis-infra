import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
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
    Descriptions
} from 'antd';
import {
    SearchOutlined,
    PlusOutlined,
    EditOutlined,
    DeleteOutlined,
    EyeOutlined,
    FileTextOutlined,
    SendOutlined
} from '@ant-design/icons';
import { contentApi } from '../api';
import { useSearchFilter } from '../hooks/useSearchFilter';
import { formatTimestamp } from '../utils/format';
import { getStatusColor, getStatusName } from '../utils/mapping';
import { ResponsiveTable } from '../components/ResponsiveTable';
import type { ResponsiveColumn, ActionItem } from '../components/ResponsiveTable';
import type { Content } from '../types';

const { Title } = Typography;
const { Option } = Select;
const { TextArea } = Input;

const ArticleManagement: React.FC = () => {
    const navigate = useNavigate();
    const [articles, setArticles] = useState<Content[]>([]);
    const [loading, setLoading] = useState(false);
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [isViewModalVisible, setIsViewModalVisible] = useState(false);
    const [editingArticle] = useState<Content | null>(null);
    const [viewingArticle, setViewingArticle] = useState<Content | null>(null);
    const [form] = Form.useForm();
    const [statusFilter, setStatusFilter] = useState<string>('');

    const { filteredData: filteredArticles, handleSearch } = useSearchFilter<Content>({
        data: articles,
        searchKeys: ['title']
    });

    useEffect(() => {
        fetchArticles();
    }, []);

    const fetchArticles = async () => {
        setLoading(true);
        try {
            const response = await contentApi.getContents({ limit: 100 });

            if (response.data.code === 200 && response.data.data) {
                const rawItems = response.data.data.items || response.data.data;
                const items = rawItems.map((item: Record<string, unknown>) => ({
                    ...item,
                    commentsNum: item.commentsNum ?? item.comments_num ?? 0,
                    author: typeof item.author === 'object' && item.author !== null 
                        ? (item.author as { name: string }).name 
                        : item.author
                }));
                setArticles(items);
            } else {
                message.error(response.data.msg || '获取文章列表失败');
            }
        } catch (error) {
            console.error('获取文章列表失败:', error);
            message.error('获取文章列表失败');
        } finally {
            setLoading(false);
        }
    };

    const handleCreate = () => {
        navigate('/editor');
    };

    const handleEdit = (record: Content) => {
        navigate(`/editor/${record.cid}`);
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
        } catch (error) {
            console.error('删除失败:', error);
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
        } catch (error) {
            console.error('发布失败:', error);
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
        } catch (error) {
            console.error('操作失败:', error);
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

    const filteredByStatus = filteredArticles.filter(article => {
        return !statusFilter || article.status === statusFilter;
    });

    // 表格列配置
    const columns: ResponsiveColumn<Content>[] = [
        {
            title: 'ID',
            dataIndex: 'cid',
            key: 'cid',
            width: 80,
            priority: 'low',
        },
        {
            title: '标题',
            dataIndex: 'title',
            key: 'title',
            width: 200,
            ellipsis: true,
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
            title: '创建时间',
            dataIndex: 'created',
            key: 'created',
            width: 160,
            priority: 'high',
            render: (created: number) => formatTimestamp(created),
        },
        {
            title: '作者',
            dataIndex: 'author',
            key: 'author',
            width: 120,
            priority: 'low',
            render: (author: string | { uid: number; name: string }) => {
                if (typeof author === 'string') return author;
                return author?.name || '-';
            },
        },
    ];

    // 操作按钮配置
    const getActions = (record: Content): ActionItem[] => [
        {
            key: 'view',
            text: '查看',
            icon: <EyeOutlined />,
            onClick: () => handleView(record),
        },
        {
            key: 'edit',
            text: '编辑',
            icon: <EditOutlined />,
            onClick: () => handleEdit(record),
        },
        {
            key: 'publish',
            text: '发布',
            icon: <SendOutlined />,
            hidden: record.status === 'publish',
            onClick: () => handlePublish(record.cid),
        },
        {
            key: 'delete',
            text: '删除',
            icon: <DeleteOutlined />,
            danger: true,
            onClick: () => {
                Modal.confirm({
                    title: '确定要删除这篇文章吗？',
                    content: '删除后无法恢复',
                    okText: '确定',
                    cancelText: '取消',
                    onOk: () => handleDelete(record.cid),
                });
            },
        },
    ];

    // 卡片标题
    const getCardTitle = (record: Content) => record.title;

    // 卡片副标题
    const getCardSubtitle = (record: Content) => {
        const authorName = typeof record.author === 'string' 
            ? record.author 
            : record.author?.name || '-';
        return `${authorName} · ${formatTimestamp(record.created)}`;
    };

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
                        placeholder="搜索标题"
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
                </Space>

                <ResponsiveTable
                    columns={columns}
                    dataSource={filteredByStatus}
                    rowKey="cid"
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
                        <Descriptions.Item label="状态">
                            <Tag color={getStatusColor(viewingArticle.status)}>
                                {getStatusName(viewingArticle.status)}
                            </Tag>
                        </Descriptions.Item>
                        <Descriptions.Item label="作者">
                            {typeof viewingArticle.author === 'string' 
                                ? viewingArticle.author 
                                : viewingArticle.author?.name || '-'}
                        </Descriptions.Item>
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
