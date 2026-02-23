import React, { useState, useEffect } from 'react';
import {
    Card,
    Button,
    Space,
    Typography,
    Tag,
    message,
    Modal,
    Avatar,
    Empty,
    Spin,
    Input,
    Segmented
} from 'antd';
import {
    CheckOutlined,
    CloseOutlined,
    EyeOutlined,
    FileTextOutlined,
    CommentOutlined,
    ClockCircleOutlined,
    UserOutlined,
    CloseCircleFilled,
    InboxOutlined
} from '@ant-design/icons';
import { contentApi, commentApi } from '../api';
import { formatTimestamp } from '../utils/format';

const { Title, Text, Paragraph } = Typography;
const { TextArea } = Input;

interface ReviewItem {
    id: number;
    type: 'content' | 'comment';
    title?: string;
    content?: string;
    author: string;
    authorId?: number;
    created: number;
    status: 'pending' | 'approved' | 'rejected';
}

const ReviewManagement: React.FC = () => {
    const [activeType, setActiveType] = useState<'content' | 'comment'>('content');
    const [loading, setLoading] = useState(false);
    const [reviews, setReviews] = useState<ReviewItem[]>([]);
    const [isRejectModalVisible, setIsRejectModalVisible] = useState(false);
    const [isDetailModalVisible, setIsDetailModalVisible] = useState(false);
    const [selectedItem, setSelectedItem] = useState<ReviewItem | null>(null);
    const [rejectReason, setRejectReason] = useState('');
    const [processingId, setProcessingId] = useState<number | null>(null);

    useEffect(() => {
        fetchReviews();
    }, [activeType]);

    const fetchReviews = async () => {
        setLoading(true);
        try {
            if (activeType === 'content') {
                const response = await contentApi.getContents({ status: 'pending', limit: 100 });
                if (response.data.code === 200) {
                    const items = (response.data.data?.items || response.data.data || []).map((item: any) => ({
                        id: item.cid,
                        type: 'content' as const,
                        title: item.title,
                        content: item.text || item.content,
                        author: typeof item.author === 'object' ? item.author.name : (item.author || '未知'),
                        authorId: typeof item.author === 'object' ? item.author.uid : undefined,
                        created: item.created,
                        status: 'pending' as const
                    }));
                    setReviews(items);
                }
            } else {
                const response = await commentApi.getComments({ status: 'pending', limit: 100 });
                if (response.data.code === 200) {
                    const items = (response.data.data?.items || response.data.data || []).map((item: any) => ({
                        id: item.coid,
                        type: 'comment' as const,
                        content: item.text,
                        author: item.author,
                        created: item.created,
                        status: 'pending' as const
                    }));
                    setReviews(items);
                }
            }
        } catch (error) {
            console.error('获取审核列表失败:', error);
            message.error('获取审核列表失败');
        } finally {
            setLoading(false);
        }
    };

    const handleApprove = async (item: ReviewItem) => {
        setProcessingId(item.id);
        try {
            if (item.type === 'content') {
                const response = await contentApi.publishContent(item.id);
                if (response.data.code === 200) {
                    message.success('文章已通过审核');
                    fetchReviews();
                }
            } else {
                const response = await commentApi.approveComment(item.id);
                if (response.data.code === 200) {
                    message.success('评论已通过审核');
                    fetchReviews();
                }
            }
        } catch (error) {
            console.error('审核失败:', error);
            message.error('审核失败');
        } finally {
            setProcessingId(null);
        }
    };

    const handleReject = (item: ReviewItem) => {
        setSelectedItem(item);
        setRejectReason('');
        setIsRejectModalVisible(true);
    };

    const handleRejectConfirm = async () => {
        if (!selectedItem) return;

        setProcessingId(selectedItem.id);
        try {
            if (selectedItem.type === 'content') {
                const response = await contentApi.updateContent(selectedItem.id, { status: 'rejected' });
                if (response.data.code === 200) {
                    message.success('文章已拒绝');
                    setIsRejectModalVisible(false);
                    fetchReviews();
                }
            } else {
                const response = await commentApi.updateComment(selectedItem.id, { status: 'spam' });
                if (response.data.code === 200) {
                    message.success('评论已拒绝');
                    setIsRejectModalVisible(false);
                    fetchReviews();
                }
            }
        } catch (error) {
            console.error('拒绝失败:', error);
            message.error('拒绝失败');
        } finally {
            setProcessingId(null);
        }
    };

    const handleViewDetail = (item: ReviewItem) => {
        setSelectedItem(item);
        setIsDetailModalVisible(true);
    };

    return (
        <div style={{ padding: '0 0 24px' }}>
            {/* Type Selector */}
            <Card style={{ marginBottom: 24, borderRadius: 12 }} bodyStyle={{ padding: '16px 24px' }}>
                <div style={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', flexWrap: 'wrap', gap: 16 }}>
                    <Segmented
                        value={activeType}
                        onChange={(value) => setActiveType(value as 'content' | 'comment')}
                        options={[
                            {
                                label: (
                                    <Space>
                                        <FileTextOutlined />
                                        文章
                                    </Space>
                                ),
                                value: 'content'
                            },
                            {
                                label: (
                                    <Space>
                                        <CommentOutlined />
                                        评论
                                    </Space>
                                ),
                                value: 'comment'
                            }
                        ]}
                        size="large"
                        style={{ background: '#f5f5f5' }}
                    />
                    <Button
                        icon={<ClockCircleOutlined />}
                        onClick={fetchReviews}
                        loading={loading}
                    >
                        刷新
                    </Button>
                </div>
            </Card>

            {/* Review List */}
            <Spin spinning={loading}>
                {reviews.length === 0 ? (
                    <Card style={{ borderRadius: 12 }}>
                        <Empty
                            image={<InboxOutlined style={{ fontSize: 64, color: '#d9d9d9' }} />}
                            description={
                                <div>
                                    <Text strong style={{ fontSize: 16 }}>暂无待审核内容</Text>
                                    <br />
                                    <Text type="secondary">当前没有需要审核的{activeType === 'content' ? '文章' : '评论'}</Text>
                                </div>
                            }
                        />
                    </Card>
                ) : (
                    <Space direction="vertical" style={{ width: '100%' }} size="middle">
                        {reviews.map((item) => (
                            <Card
                                key={item.id}
                                style={{
                                    borderRadius: 12,
                                    transition: 'all 0.3s ease',
                                    border: '1px solid #f0f0f0',
                                    cursor: 'pointer'
                                }}
                                bodyStyle={{ padding: 20 }}
                                hoverable
                                onClick={() => handleViewDetail(item)}
                            >
                                <div style={{ display: 'flex', gap: 16 }}>
                                    {/* Icon */}
                                    <div style={{
                                        width: 48,
                                        height: 48,
                                        borderRadius: 12,
                                        background: item.type === 'content' ? '#e6f7ff' : '#f6ffed',
                                        display: 'flex',
                                        alignItems: 'center',
                                        justifyContent: 'center',
                                        flexShrink: 0
                                    }}>
                                        {item.type === 'content' ? (
                                            <FileTextOutlined style={{ fontSize: 24, color: '#1890ff' }} />
                                        ) : (
                                            <CommentOutlined style={{ fontSize: 24, color: '#52c41a' }} />
                                        )}
                                    </div>

                                    {/* Content */}
                                    <div style={{ flex: 1, minWidth: 0 }}>
                                        <div style={{ marginBottom: 8 }}>
                                            <Tag color={item.type === 'content' ? 'blue' : 'green'}>
                                                {item.type === 'content' ? '文章' : '评论'}
                                            </Tag>
                                            <Tag color="warning">待审核</Tag>
                                        </div>

                                        {item.title && (
                                            <Title level={5} style={{ margin: '0 0 8px', fontWeight: 600 }}>
                                                {item.title}
                                            </Title>
                                        )}

                                        <Paragraph
                                            ellipsis={{ rows: 2 }}
                                            style={{ marginBottom: 12, color: '#666' }}
                                        >
                                            {item.content}
                                        </Paragraph>

                                        <div style={{ display: 'flex', alignItems: 'center', gap: 16, flexWrap: 'wrap' }}>
                                            <Space>
                                                <Avatar size="small" icon={<UserOutlined />} />
                                                <Text type="secondary">{item.author}</Text>
                                            </Space>
                                            <Text type="secondary" style={{ fontSize: 12 }}>
                                                <ClockCircleOutlined style={{ marginRight: 4 }} />
                                                {formatTimestamp(item.created)}
                                            </Text>
                                        </div>
                                    </div>

                                    {/* Actions */}
                                    <div style={{
                                        display: 'flex',
                                        flexDirection: 'column',
                                        gap: 8,
                                        justifyContent: 'center',
                                        flexShrink: 0
                                    }}>
                                        <Button
                                            type="primary"
                                            icon={<CheckOutlined />}
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                handleApprove(item);
                                            }}
                                            loading={processingId === item.id}
                                            style={{
                                                background: '#52c41a',
                                                borderColor: '#52c41a',
                                                borderRadius: 8
                                            }}
                                        >
                                            通过
                                        </Button>
                                        <Button
                                            danger
                                            icon={<CloseOutlined />}
                                            onClick={(e) => {
                                                e.stopPropagation();
                                                handleReject(item);
                                            }}
                                            loading={processingId === item.id}
                                            style={{ borderRadius: 8 }}
                                        >
                                            拒绝
                                        </Button>
                                    </div>
                                </div>
                            </Card>
                        ))}
                    </Space>
                )}
            </Spin>

            {/* Reject Modal */}
            <Modal
                title={
                    <Space>
                        <CloseCircleFilled style={{ color: '#ff4d4f' }} />
                        <span>拒绝审核</span>
                    </Space>
                }
                open={isRejectModalVisible}
                onOk={handleRejectConfirm}
                onCancel={() => setIsRejectModalVisible(false)}
                okText="确认拒绝"
                cancelText="取消"
                okButtonProps={{ danger: true, size: 'large' }}
                cancelButtonProps={{ size: 'large' }}
                width={500}
            >
                <div style={{ padding: '16px 0' }}>
                    <Text>请输入拒绝原因（可选）：</Text>
                    <TextArea
                        value={rejectReason}
                        onChange={(e) => setRejectReason(e.target.value)}
                        placeholder="例如：内容违规、广告信息、与主题无关等"
                        rows={4}
                        style={{ marginTop: 12, borderRadius: 8 }}
                        maxLength={200}
                        showCount
                    />
                </div>
            </Modal>

            {/* Detail Modal */}
            <Modal
                title={
                    <Space>
                        <EyeOutlined />
                        <span>内容详情</span>
                    </Space>
                }
                open={isDetailModalVisible}
                onCancel={() => setIsDetailModalVisible(false)}
                footer={[
                    <Button key="close" onClick={() => setIsDetailModalVisible(false)}>
                        关闭
                    </Button>
                ]}
                width={700}
            >
                {selectedItem && (
                    <div style={{ padding: '16px 0' }}>
                        <div style={{ marginBottom: 16 }}>
                            <Tag color={selectedItem.type === 'content' ? 'blue' : 'green'}>
                                {selectedItem.type === 'content' ? '文章' : '评论'}
                            </Tag>
                            <Tag color="warning">待审核</Tag>
                        </div>

                        {selectedItem.title && (
                            <Title level={4} style={{ marginBottom: 16 }}>
                                {selectedItem.title}
                            </Title>
                        )}

                        <div style={{
                            background: '#f5f5f5',
                            padding: 16,
                            borderRadius: 8,
                            marginBottom: 16,
                            maxHeight: 300,
                            overflow: 'auto'
                        }}>
                            <Text style={{ whiteSpace: 'pre-wrap' }}>
                                {selectedItem.content}
                            </Text>
                        </div>

                        <div style={{ marginBottom: 16 }}>
                            <Space size="large">
                                <Space>
                                    <Avatar size="small" icon={<UserOutlined />} />
                                    <Text type="secondary">作者:</Text>
                                    <Text strong>{selectedItem.author}</Text>
                                </Space>
                                <Space>
                                    <ClockCircleOutlined style={{ color: '#999' }} />
                                    <Text type="secondary">提交时间:</Text>
                                    <Text strong>{formatTimestamp(selectedItem.created)}</Text>
                                </Space>
                            </Space>
                        </div>

                        <div style={{ marginTop: 24, display: 'flex', gap: 12 }}>
                            <Button
                                type="primary"
                                icon={<CheckOutlined />}
                                onClick={() => {
                                    setIsDetailModalVisible(false);
                                    handleApprove(selectedItem);
                                }}
                                style={{ background: '#52c41a', borderColor: '#52c41a' }}
                            >
                                通过
                            </Button>
                            <Button
                                danger
                                icon={<CloseOutlined />}
                                onClick={() => {
                                    setIsDetailModalVisible(false);
                                    handleReject(selectedItem);
                                }}
                            >
                                拒绝
                            </Button>
                        </div>
                    </div>
                )}
            </Modal>
        </div>
    );
};

export default ReviewManagement;
