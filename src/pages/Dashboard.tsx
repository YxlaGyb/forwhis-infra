import React, { useState, useEffect, useCallback, useRef, useMemo } from 'react';
import { Row, Col, Card, Statistic, Typography, Spin, message, Slider } from 'antd';
import { UserOutlined, FileTextOutlined, MessageOutlined } from '@ant-design/icons';
import { statsApi } from '../api';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';

const { Title, Paragraph, Text } = Typography;

interface UserStats {
    total: number;
    todayNew: number;
    weekNew: number;
}

interface ContentStats {
    posts: number;
    pages: number;
    drafts: number;
    published: number;
    private: number;
}

interface CommentStats {
    total: number;
    pending: number;
    approved: number;
    spam: number;
}

interface DashboardData {
    users?: UserStats;
    contents?: ContentStats;
    comments?: CommentStats;
}

interface TrendData {
    date: string;
    user: number;
    content: number;
    comment: number;
}

interface TooltipPayloadEntry {
    name: string;
    value: number;
    color: string;
}

interface TooltipProps {
    active?: boolean;
    payload?: TooltipPayloadEntry[];
    label?: string;
}

const Dashboard: React.FC = () => {
    const [dashboardData, setDashboardData] = useState<DashboardData>({});
    const [trendData, setTrendData] = useState<TrendData[]>([]);
    const [loading, setLoading] = useState(true);
    const [trendLoading, setTrendLoading] = useState(false);
    const [days, setDays] = useState(7);
    const [hiddenLines, setHiddenLines] = useState<Set<string>>(new Set());
    
    const abortControllerRef = useRef<AbortController | null>(null);
    const trendAbortControllerRef = useRef<AbortController | null>(null);

    useEffect(() => {
        fetchDashboardData();
        return () => {
            if (abortControllerRef.current) {
                abortControllerRef.current.abort();
            }
        };
    }, []);

    useEffect(() => {
        fetchTrendData();
        return () => {
            if (trendAbortControllerRef.current) {
                trendAbortControllerRef.current.abort();
            }
        };
    }, [days]);

    const fetchDashboardData = async () => {
        abortControllerRef.current = new AbortController();
        try {
            setLoading(true);
            const response = await statsApi.getDashboard();
            
            if (abortControllerRef.current?.signal.aborted) return;
            
            if (response.data.code === 200 && response.data.data) {
                setDashboardData(response.data.data);
            } else {
                message.warning('获取数据失败，使用默认展示');
            }
        } catch (error: unknown) {
            if (abortControllerRef.current?.signal.aborted) return;
            console.error('获取仪表盘数据失败:', error);
            const err = error as { response?: { status?: number } };
            if (err.response?.status === 401) {
                message.error('未登录或登录已过期');
            } else {
                message.error('获取仪表盘数据失败');
            }
        } finally {
            if (!abortControllerRef.current?.signal.aborted) {
                setLoading(false);
            }
        }
    };

    const fetchTrendData = async () => {
        trendAbortControllerRef.current = new AbortController();
        try {
            setTrendLoading(true);
            const response = await statsApi.getTrend(days);
            
            if (trendAbortControllerRef.current?.signal.aborted) return;
            
            if (response.data.code === 200 && response.data.data) {
                setTrendData(response.data.data);
            }
        } catch (error: unknown) {
            if (trendAbortControllerRef.current?.signal.aborted) return;
            console.error('获取趋势数据失败:', error);
        } finally {
            if (!trendAbortControllerRef.current?.signal.aborted) {
                setTrendLoading(false);
            }
        }
    };

    const formatDate = useCallback((dateStr: string) => {
        const date = new Date(dateStr);
        const month = date.getMonth() + 1;
        const day = date.getDate();
        return `${month}/${day}`;
    }, []);

    const increments = useMemo(() => {
        const selectedDaysData = trendData.slice(-days);
        let userInc = 0, contentInc = 0, commentInc = 0;
        selectedDaysData.forEach(item => {
            userInc += item.user;
            contentInc += item.content;
            commentInc += item.comment;
        });
        return { userInc, contentInc, commentInc };
    }, [trendData, days]);

    const toggleLine = useCallback((dataKey: string) => {
        const newHidden = new Set(hiddenLines);
        if (newHidden.has(dataKey)) {
            newHidden.delete(dataKey);
        } else {
            newHidden.add(dataKey);
        }
        setHiddenLines(newHidden);
    }, [hiddenLines]);

    const CustomTooltip = ({ active, payload, label }: TooltipProps) => {
        if (active && payload && payload.length) {
            return (
                <div style={{ 
                    background: 'rgba(255, 255, 255, 0.95)', 
                    padding: '12px 16px', 
                    borderRadius: 12, 
                    boxShadow: '0 4px 20px rgba(0,0,0,0.1)',
                    border: '1px solid rgba(0,0,0,0.05)'
                }}>
                    <div style={{ fontWeight: 600, marginBottom: 8, color: '#333' }}>{label}</div>
                    {payload.map((entry, index) => (
                        <div key={index} style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 4 }}>
                            <div style={{ 
                                width: 8, 
                                height: 8, 
                                borderRadius: '50%', 
                                background: entry.color 
                            }} />
                            <span style={{ color: '#666', fontSize: 13 }}>{entry.name}:</span>
                            <span style={{ fontWeight: 600, color: entry.color }}>{entry.value}</span>
                        </div>
                    ))}
                </div>
            );
        }
        return null;
    };

    if (loading) {
        return (
            <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 200 }}>
                <Spin size="large" />
            </div>
        );
    }

    return (
        <div>
            <Title level={2} style={{ marginBottom: 8 }}>数据总览</Title>
            <Paragraph type="secondary" style={{ marginBottom: 24 }}>
                实时展示系统各项核心数据指标及变化趋势
            </Paragraph>

            <div style={{ marginBottom: 16, fontSize: 16, fontWeight: 600, color: '#333' }}>数据总和</div>
            <Row gutter={[24, 24]}>
                <Col xs={24} sm={12} lg={8}>
                    <Card>
                        <Statistic
                            title="总用户数"
                            value={dashboardData.users?.total || 0}
                            prefix={<UserOutlined />}
                            valueStyle={{ color: '#3f8600' }}
                        />
                        <div style={{ marginTop: 8, fontSize: 12, color: '#888' }}>
                            今日新增 {dashboardData.users?.todayNew || 0} | 本周新增 {dashboardData.users?.weekNew || 0}
                        </div>
                    </Card>
                </Col>

                <Col xs={24} sm={12} lg={8}>
                    <Card>
                        <Statistic
                            title="文章总数"
                            value={dashboardData.contents?.posts || 0}
                            prefix={<FileTextOutlined />}
                            valueStyle={{ color: '#1890ff' }}
                        />
                        <div style={{ marginTop: 8, fontSize: 12, color: '#888' }}>
                            已发布 {dashboardData.contents?.published || 0} | 草稿 {dashboardData.contents?.drafts || 0}
                        </div>
                    </Card>
                </Col>

                <Col xs={24} sm={12} lg={8}>
                    <Card>
                        <Statistic
                            title="评论总数"
                            value={dashboardData.comments?.total || 0}
                            prefix={<MessageOutlined />}
                            valueStyle={{ color: '#cf1322' }}
                        />
                        <div style={{ marginTop: 8, fontSize: 12, color: '#888' }}>
                            已审核 {dashboardData.comments?.approved || 0} | 待审核 {dashboardData.comments?.pending || 0}
                        </div>
                    </Card>
                </Col>
            </Row>

            <div style={{ marginTop: 24, marginBottom: 8, fontSize: 16, fontWeight: 600, color: '#333' }}>
                增量数据
                <Text type="secondary" style={{ fontWeight: 400, marginLeft: 8, fontSize: 14 }}>
                    调整时间范围
                </Text>
            </div>
            <Slider
                min={1}
                max={30}
                value={days}
                onChange={setDays}
                style={{ marginBottom: 16 }}
                tooltip={{ formatter: null }}
            />
            <Row gutter={[24, 24]} style={{ marginTop: 8 }}>
                <Col xs={24} sm={8}>
                    <Card style={{ textAlign: 'center' }}>
                        <div style={{ color: '#888', fontSize: 14, marginBottom: 8 }}>新增用户</div>
                        <div style={{ fontSize: 28, fontWeight: 'bold', color: '#3f8600' }}>+{increments.userInc}</div>
                    </Card>
                </Col>
                <Col xs={24} sm={8}>
                    <Card style={{ textAlign: 'center' }}>
                        <div style={{ color: '#888', fontSize: 14, marginBottom: 8 }}>新增文章</div>
                        <div style={{ fontSize: 28, fontWeight: 'bold', color: '#1890ff' }}>+{increments.contentInc}</div>
                    </Card>
                </Col>
                <Col xs={24} sm={8}>
                    <Card style={{ textAlign: 'center' }}>
                        <div style={{ color: '#888', fontSize: 14, marginBottom: 8 }}>新增评论</div>
                        <div style={{ fontSize: 28, fontWeight: 'bold', color: '#cf1322' }}>+{increments.commentInc}</div>
                    </Card>
                </Col>
            </Row>

            <Row gutter={[24, 24]} style={{ marginTop: 24 }}>
                <Col span={24}>
                    <Card 
                        title="趋势图"
                        style={{ minHeight: 400, position: 'relative' }}
                    >
                        <div style={{ position: 'absolute', top: 16, right: 24, display: 'flex', gap: 16, zIndex: 10 }}>
                            {[
                                { key: 'user', name: '用户', color: '#3f8600' },
                                { key: 'content', name: '文章', color: '#1890ff' },
                                { key: 'comment', name: '评论', color: '#cf1322' }
                            ].map(item => (
                                <div 
                                    key={item.key}
                                    onClick={() => toggleLine(item.key)}
                                    style={{ 
                                        display: 'flex', 
                                        alignItems: 'center', 
                                        gap: 6, 
                                        cursor: 'pointer',
                                        opacity: hiddenLines.has(item.key) ? 0.4 : 1,
                                        transition: 'opacity 0.2s'
                                    }}
                                >
                                    <div style={{ 
                                        width: 10, 
                                        height: 10, 
                                        borderRadius: '50%', 
                                        background: item.color 
                                    }} />
                                    <span style={{ fontSize: 13, color: '#666' }}>{item.name}</span>
                                </div>
                            ))}
                        </div>
                        {trendLoading ? (
                            <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 300 }}>
                                <Spin size="large" />
                            </div>
                        ) : trendData.length > 0 ? (
                            <ResponsiveContainer width="100%" height={320}>
                                <AreaChart data={trendData} margin={{ top: 40, right: 30, left: 10, bottom: 10 }}>
                                    <defs>
                                        <linearGradient id="colorUser" x1="0" y1="0" x2="0" y2="1">
                                            <stop offset="5%" stopColor="#3f8600" stopOpacity={0.3}/>
                                            <stop offset="95%" stopColor="#3f8600" stopOpacity={0}/>
                                        </linearGradient>
                                        <linearGradient id="colorContent" x1="0" y1="0" x2="0" y2="1">
                                            <stop offset="5%" stopColor="#1890ff" stopOpacity={0.3}/>
                                            <stop offset="95%" stopColor="#1890ff" stopOpacity={0}/>
                                        </linearGradient>
                                        <linearGradient id="colorComment" x1="0" y1="0" x2="0" y2="1">
                                            <stop offset="5%" stopColor="#cf1322" stopOpacity={0.3}/>
                                            <stop offset="95%" stopColor="#cf1322" stopOpacity={0}/>
                                        </linearGradient>
                                    </defs>
                                    <CartesianGrid strokeDasharray="3 3" stroke="#f5f5f5" vertical={false} />
                                    <XAxis 
                                        dataKey="date" 
                                        tickFormatter={formatDate}
                                        stroke="#d9d9d9"
                                        tick={{ fill: '#999', fontSize: 11 }}
                                        axisLine={{ stroke: '#e8e8e8' }}
                                        tickLine={false}
                                        dy={10}
                                    />
                                    <YAxis 
                                        stroke="#d9d9d9"
                                        tick={{ fill: '#999', fontSize: 11 }}
                                        axisLine={false}
                                        tickLine={false}
                                        dx={-5}
                                    />
                                    <Tooltip content={<CustomTooltip />} />
                                    {!hiddenLines.has('user') && (
                                        <Area 
                                            type="monotone" 
                                            dataKey="user" 
                                            name="用户" 
                                            stroke="#3f8600" 
                                            strokeWidth={2.5}
                                            fillOpacity={1}
                                            fill="url(#colorUser)"
                                            dot={false}
                                            activeDot={{ r: 5, strokeWidth: 0 }}
                                        />
                                    )}
                                    {!hiddenLines.has('content') && (
                                        <Area 
                                            type="monotone" 
                                            dataKey="content" 
                                            name="文章" 
                                            stroke="#1890ff" 
                                            strokeWidth={2.5}
                                            fillOpacity={1}
                                            fill="url(#colorContent)"
                                            dot={false}
                                            activeDot={{ r: 5, strokeWidth: 0 }}
                                        />
                                    )}
                                    {!hiddenLines.has('comment') && (
                                        <Area 
                                            type="monotone" 
                                            dataKey="comment" 
                                            name="评论" 
                                            stroke="#cf1322" 
                                            strokeWidth={2.5}
                                            fillOpacity={1}
                                            fill="url(#colorComment)"
                                            dot={false}
                                            activeDot={{ r: 5, strokeWidth: 0 }}
                                        />
                                    )}
                                </AreaChart>
                            </ResponsiveContainer>
                        ) : (
                            <div style={{ textAlign: 'center', padding: '40px 0' }}>
                                <Text type="secondary">暂无趋势数据</Text>
                            </div>
                        )}
                    </Card>
                </Col>
            </Row>
        </div>
    );
};

export default Dashboard;
