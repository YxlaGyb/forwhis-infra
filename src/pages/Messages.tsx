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
    message,
    Popconfirm,
    Badge
} from 'antd';
import {
    SearchOutlined,
    MessageOutlined,
    EyeOutlined,
    DeleteOutlined
} from '@ant-design/icons';
import { commentApi } from '../api';
import { useSearchFilter } from '../hooks/useSearchFilter';
import type { Message } from '../types';

const { Title } = Typography;
const { TextArea } = Input;

const Messages: React.FC = () => {
    const [messages, setMessages] = useState<Message[]>([]);
    const [loading, setLoading] = useState(false);
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [selectedMessage, setSelectedMessage] = useState<Message | null>(null);
    
    const abortControllerRef = useRef<AbortController | null>(null);

    const { filteredData: filteredMessages, handleSearch } = useSearchFilter<Message>({
        data: messages,
        searchKeys: ['title', 'content', 'sender']
    });

    useEffect(() => {
        fetchMessages();
        return () => {
            if (abortControllerRef.current) {
                abortControllerRef.current.abort();
            }
        };
    }, []);

    const fetchMessages = async () => {
        abortControllerRef.current = new AbortController();
        setLoading(true);
        try {
            const response = await commentApi.getComments();
            
            if (abortControllerRef.current?.signal.aborted) return;
            
            if (response.data.code === 200 && response.data.data) {
                setMessages(response.data.data);
            } else {
                message.error('获取消息列表失败');
            }
        } catch (error) {
            if (abortControllerRef.current?.signal.aborted) return;
            console.error('获取消息列表失败:', error);
            message.error('获取消息列表失败');
        } finally {
            if (!abortControllerRef.current?.signal.aborted) {
                setLoading(false);
            }
        }
    };

    const handleViewMessage = async (messageItem: Message) => {
        setSelectedMessage(messageItem);
        setIsModalVisible(true);
        
        if (messageItem.status === 'unread') {
            try {
                const coid = parseInt(messageItem.id);
                await commentApi.approveComment(coid);
                
                setMessages(messages.map(msg => 
                    msg.id === messageItem.id ? { ...msg, status: 'read' } : msg
                ));
            } catch (error) {
                console.error('标记已读失败:', error);
            }
        }
    };

    const handleDeleteMessage = async (messageId: string) => {
        try {
            const coid = parseInt(messageId);
            await commentApi.deleteComment(coid);
            setMessages(messages.filter(message => message.id !== messageId));
            message.success('消息删除成功');
        } catch (error) {
            console.error('删除消息失败:', error);
            message.error('删除消息失败');
        }
    };

    const handleModalCancel = () => {
        setIsModalVisible(false);
        setSelectedMessage(null);
    };

    const columns = [
        {
            title: 'ID',
            dataIndex: 'id',
            key: 'id',
            width: 80,
        },
        {
            title: '标题',
            dataIndex: 'title',
            key: 'title',
            render: (text: string, record: Message) => (
                <Space>
                    {record.status === 'unread' && <Badge dot />}
                    {text}
                </Space>
            ),
        },
        {
            title: '发送者',
            dataIndex: 'sender',
            key: 'sender',
        },
        {
            title: '接收者',
            dataIndex: 'receiver',
            key: 'receiver',
        },
        {
            title: '状态',
            dataIndex: 'status',
            key: 'status',
            render: (status: string) => (
                <Tag color={status === 'unread' ? 'red' : 'green'}>
                    {status === 'unread' ? '未读' : '已读'}
                </Tag>
            ),
        },
        {
            title: '发送时间',
            dataIndex: 'createTime',
            key: 'createTime',
        },
        {
            title: '操作',
            key: 'action',
            render: (_record: unknown, record: Message) => (
                <Space size="middle">
                    <Button
                        type="link"
                        icon={<EyeOutlined />}
                        onClick={() => handleViewMessage(record)}
                    >
                        查看
                    </Button>
                    <Popconfirm
                        title="确定删除这条消息吗？"
                        onConfirm={() => handleDeleteMessage(record.id)}
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
                <MessageOutlined /> 消息管理
            </Title>
            
            <Card>
                <div style={{ marginBottom: 16, display: 'flex', justifyContent: 'space-between' }}>
                    <Input
                        placeholder="搜索消息标题、内容或发送者"
                        prefix={<SearchOutlined />}
                        style={{ width: 300 }}
                        onChange={(e) => handleSearch(e.target.value)}
                        allowClear
                    />
                </div>
                
                <Table
                    columns={columns}
                    dataSource={filteredMessages}
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
                title="消息详情"
                open={isModalVisible}
                onCancel={handleModalCancel}
                footer={[
                    <Button key="close" onClick={handleModalCancel}>
                        关闭
                    </Button>
                ]}
                width={600}
            >
                {selectedMessage && (
                    <div>
                        <Form layout="vertical">
                            <Form.Item label="标题">
                                <Input value={selectedMessage.title} readOnly />
                            </Form.Item>
                            
                            <Form.Item label="发送者">
                                <Input value={selectedMessage.sender} readOnly />
                            </Form.Item>
                            
                            <Form.Item label="接收者">
                                <Input value={selectedMessage.receiver} readOnly />
                            </Form.Item>
                            
                            <Form.Item label="发送时间">
                                <Input value={selectedMessage.createTime} readOnly />
                            </Form.Item>
                            
                            <Form.Item label="内容">
                                <TextArea 
                                    value={selectedMessage.content} 
                                    readOnly 
                                    rows={6}
                                    style={{ resize: 'none' }}
                                />
                            </Form.Item>
                        </Form>
                    </div>
                )}
            </Modal>
        </div>
    );
};

export default Messages;
