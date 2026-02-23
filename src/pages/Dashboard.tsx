import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { Row, Col, Card, Statistic, Typography, Spin, message, Slider } from 'antd';
import { UserOutlined, FileTextOutlined, MessageOutlined, EyeOutlined, EditOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { statsApi } from '../api';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import type { DashboardData, TrendData, TooltipProps } from '../types';

const { Title, Paragraph, Text } = Typography;

const Dashboard: React.FC = () => {
    const navigate = useNavigate();
    const [dashboardData, setDashboardData] = useState<DashboardData>({});
    const [trendData, setTrendData] = useState<TrendData[]>([]);
    const [loading, setLoading] = useState(true);
    const [trendLoading, setTrendLoading] = useState(false);
    const [days, setDays] = useState(7);
    const [sliderValue, setSliderValue] = useState(7);
    const [hiddenLines, setHiddenLines] = useState<Set<string>>(new Set());
    const [isMobile, setIsMobile] = useState(false);

    useEffect(() => {
        const checkMobile = () => {
            setIsMobile(window.innerWidth <= 768);
        };
        checkMobile();
        window.addEventListener('resize', checkMobile);
        return () => window.removeEventListener('resize', checkMobile);
    }, []);

    useEffect(() => {
        fetchDashboardData();
    }, []);

    useEffect(() => {
        fetchTrendData();
    }, [days]);

    const fetchDashboardData = async () => {
        try {
            setLoading(true);
            const response = await statsApi.getDashboard();
            
            if (response.data.code === 200 && response.data.data) {
                setDashboardData(response.data.data);
            } else {
                message.warning('获取数据失败，使用默认展示');
            }
        } catch (error: unknown) {
            console.error('获取仪表盘数据失败:', error);
            const err = error as { response?: { status?: number } };
            if (err.response?.status === 401) {
                message.error('未登录或登录已过期');
            } else {
                message.error('获取仪表盘数据失败');
            }
        } finally {
            setLoading(false);
        }
    };

    const fetchTrendData = async () => {
        try {
            setTrendLoading(true);
            const response = await statsApi.getTrend(days);
            
            if (response.data.code === 200 && response.data.data) {
                setTrendData(response.data.data);
            }
        } catch (error: unknown) {
            console.error('获取趋势数据失败:', error);
        } finally {
            setTrendLoading(false);
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
        let userInc = 0, contentInc = 0, commentInc = 0, viewsInc = 0;
        selectedDaysData.forEach(item => {
            userInc += item.user;
            contentInc += item.content;
            commentInc += item.comment;
            viewsInc += item.views;
        });
        return { userInc, contentInc, commentInc, viewsInc };
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

            <div style={{ marginBottom: 8, fontSize: 16, fontWeight: 600, color: '#333' }}>数据总和</div>
            <Paragraph type="secondary" style={{ marginBottom: 16 }}>
                网站在建立以来, 各项指标的总和数据
            </Paragraph>
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
                    <Card
                        hoverable
                        onClick={() => navigate('/articles')}
                        bodyStyle={{ padding: 24, cursor: 'pointer' }}
                    >
                        <Statistic
                            title={
                                <span style={{ display: 'flex', alignItems: 'center', gap: 8 }}>
                                    文章总数
                                    <EditOutlined style={{ fontSize: 12, opacity: 0.5 }} />
                                </span>
                            }
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

                <Col xs={24} sm={12} lg={8}>
                    <Card>
                        <Statistic
                            title="文章访问量"
                            value={dashboardData.contents?.views || 0}
                            prefix={<EyeOutlined />}
                            valueStyle={{ color: '#722ed1' }}
                        />
                        <div style={{ marginTop: 8, fontSize: 12, color: '#888' }}>
                            总阅读量统计
                        </div>
                    </Card>
                </Col>
            </Row>

            <div style={{ marginTop: 24, marginBottom: 8, fontSize: 16, fontWeight: 600, color: '#333' }}>
                增量数据
                <Text type="secondary" style={{ fontWeight: 400, marginLeft: 8, fontSize: 14 }}>
                    最近 {sliderValue} 天
                </Text>
            </div>
            <Paragraph type="secondary" style={{ marginBottom: 16 }}>
                网站在指定时间段内, 各项指标的增量数据
            </Paragraph>
            <Slider
                min={1}
                max={30}
                value={sliderValue}
                onChange={setSliderValue}
                onAfterChange={setDays}
                style={{ marginBottom: 16 }}
                tooltip={{ formatter: null }}
            />
            <Row gutter={[24, 24]} style={{ marginTop: 8 }}>
                <Col xs={24} sm={12} lg={6}>
                    <Card style={{ textAlign: 'center' }}>
                        <div style={{ color: '#888', fontSize: 14, marginBottom: 8 }}>新增用户</div>
                        <div style={{ fontSize: 28, fontWeight: 'bold', color: '#3f8600' }}>+{increments.userInc}</div>
                    </Card>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <Card style={{ textAlign: 'center' }}>
                        <div style={{ color: '#888', fontSize: 14, marginBottom: 8 }}>新增文章</div>
                        <div style={{ fontSize: 28, fontWeight: 'bold', color: '#1890ff' }}>+{increments.contentInc}</div>
                    </Card>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <Card style={{ textAlign: 'center' }}>
                        <div style={{ color: '#888', fontSize: 14, marginBottom: 8 }}>新增评论</div>
                        <div style={{ fontSize: 28, fontWeight: 'bold', color: '#cf1322' }}>+{increments.commentInc}</div>
                    </Card>
                </Col>
                <Col xs={24} sm={12} lg={6}>
                    <Card style={{ textAlign: 'center' }}>
                        <div style={{ color: '#888', fontSize: 14, marginBottom: 8 }}>新增访问</div>
                        <div style={{ fontSize: 28, fontWeight: 'bold', color: '#722ed1' }}>+{increments.viewsInc}</div>
                    </Card>
                </Col>
            </Row>

            <Row gutter={[24, 24]} style={{ marginTop: 24 }}>
                <Col span={24}>
                    <div style={{ marginBottom: 8, fontSize: 16, fontWeight: 600, color: '#333' }}>趋势图</div>
                    <Paragraph type="secondary" style={{ marginBottom: 16 }}>
                        如果要用一张可视化的图表来表示网站的增量数据状态, 那就是下面这张图
                    </Paragraph>
                    <Card
                        style={{ minHeight: isMobile ? 320 : 360, position: 'relative' }}
                    >
                        <div style={{
                            position: 'absolute',
                            top: isMobile ? 8 : 16,
                            right: isMobile ? 8 : 24,
                            left: isMobile ? 8 : 'auto',
                            display: 'flex',
                            gap: isMobile ? 8 : 16,
                            zIndex: 10,
                            flexWrap: 'wrap',
                            justifyContent: isMobile ? 'flex-start' : 'flex-end'
                        }}>
                            {[
                                { key: 'user', name: '用户', color: '#3f8600' },
                                { key: 'content', name: '文章', color: '#1890ff' },
                                { key: 'comment', name: '评论', color: '#cf1322' },
                                { key: 'views', name: '访问', color: '#722ed1' }
                            ].map(item => (
                                <div
                                    key={item.key}
                                    onClick={() => toggleLine(item.key)}
                                    style={{
                                        display: 'flex',
                                        alignItems: 'center',
                                        gap: isMobile ? 4 : 6,
                                        cursor: 'pointer',
                                        opacity: hiddenLines.has(item.key) ? 0.4 : 1,
                                        transition: 'opacity 0.2s'
                                    }}
                                >
                                    <div style={{
                                        width: isMobile ? 8 : 10,
                                        height: isMobile ? 8 : 10,
                                        borderRadius: '50%',
                                        background: item.color
                                    }} />
                                    <span style={{ fontSize: isMobile ? 11 : 13, color: '#666' }}>{item.name}</span>
                                </div>
                            ))}
                        </div>
                        {trendLoading ? (
                            <div style={{ display: 'flex', justifyContent: 'center', alignItems: 'center', height: 300 }}>
                                <Spin size="large" />
                            </div>
                        ) : trendData.length > 0 ? (
                            <ResponsiveContainer width="100%" height={isMobile ? 280 : 320}>
                                <AreaChart data={trendData} margin={{ top: isMobile ? 50 : 40, right: isMobile ? 10 : 30, left: isMobile ? -20 : 10, bottom: isMobile ? 5 : 10 }}>
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
                                        <linearGradient id="colorViews" x1="0" y1="0" x2="0" y2="1">
                                            <stop offset="5%" stopColor="#722ed1" stopOpacity={0.3}/>
                                            <stop offset="95%" stopColor="#722ed1" stopOpacity={0}/>
                                        </linearGradient>
                                    </defs>
                                    <CartesianGrid strokeDasharray="3 3" stroke="#f5f5f5" vertical={false} />
                                    <XAxis
                                        dataKey="date"
                                        tickFormatter={formatDate}
                                        stroke="#d9d9d9"
                                        tick={{ fill: '#999', fontSize: isMobile ? 9 : 11 }}
                                        axisLine={{ stroke: '#e8e8e8' }}
                                        tickLine={false}
                                        dy={isMobile ? 5 : 10}
                                        interval={isMobile ? 'preserveStartEnd' : 0}
                                    />
                                    <YAxis
                                        stroke="#d9d9d9"
                                        tick={{ fill: '#999', fontSize: isMobile ? 9 : 11 }}
                                        axisLine={false}
                                        tickLine={false}
                                        dx={isMobile ? -15 : -5}
                                        width={isMobile ? 30 : 60}
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
                                    {!hiddenLines.has('views') && (
                                        <Area 
                                            type="monotone" 
                                            dataKey="views" 
                                            name="访问" 
                                            stroke="#722ed1" 
                                            strokeWidth={2.5}
                                            fillOpacity={1}
                                            fill="url(#colorViews)"
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
