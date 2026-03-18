import React, { useState, useEffect, useRef } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import {
    Input,
    Button,
    message,
    Typography,
    Divider,
    Tag,
    Tooltip,
    Popconfirm,
    Drawer
} from 'antd';
import {
    SaveOutlined,
    SendOutlined,
    EyeOutlined,
    CloseOutlined,
    PictureOutlined,
    LinkOutlined,
    BoldOutlined,
    ItalicOutlined,
    UnorderedListOutlined,
    OrderedListOutlined,
    CodeOutlined,
    DeleteOutlined,
    SmileOutlined,
    PlusOutlined,
    SettingOutlined,
    TagOutlined,
    FolderOutlined
} from '@ant-design/icons';
import { contentApi, metaApi } from '../api';

const { Text } = Typography;
const { TextArea } = Input;

interface Category {
    mid: number;
    name: string;
    slug: string;
    parent?: number;
}

interface EditorState {
    title: string;
    text: string;
    status: 'draft' | 'publish' | 'private';
    type: 'post' | 'page';
    slug: string;
    categoryMid: number;
}

const ArticleEditor: React.FC = () => {
    const navigate = useNavigate();
    const { id } = useParams<{ id?: string }>();
    const isEditMode = Boolean(id);
    const textAreaRef = useRef<HTMLTextAreaElement>(null);
    const [isMobile, setIsMobile] = useState(false);
    const [saving, setSaving] = useState(false);
    const [publishing, setPublishing] = useState(false);
    const [titleCount, setTitleCount] = useState(0);
    const [contentCount, setContentCount] = useState(0);
    const [lastSaved, setLastSaved] = useState<Date | null>(null);
    const [isDirty, setIsDirty] = useState(false);
    const [showToolbar, setShowToolbar] = useState(true);
    const [categoryDrawerOpen, setCategoryDrawerOpen] = useState(false);
    const [categories, setCategories] = useState<Category[]>([]);

    const [formData, setFormData] = useState<EditorState>({
        title: '',
        text: '',
        status: 'draft',
        type: 'post',
        slug: '',
        categoryMid: 0
    });

    useEffect(() => {
        const checkMobile = () => {
            setIsMobile(window.innerWidth <= 768);
        };
        checkMobile();
        window.addEventListener('resize', checkMobile);
        return () => window.removeEventListener('resize', checkMobile);
    }, []);

    useEffect(() => {
        fetchCategories();
    }, []);

    const fetchCategories = async () => {
        try {
            const response = await metaApi.getMetas({ type: 'category' });
            if (response.data.code === 200 && response.data.data) {
                setCategories(response.data.data || []);
            }
        } catch (error) {
            console.error('获取分类失败:', error);
        }
    };

    useEffect(() => {
        if (isEditMode && id) {
            fetchArticle(parseInt(id));
        }
    }, [isEditMode, id]);

    useEffect(() => {
        setTitleCount(formData.title.length);
        setContentCount(formData.text.replace(/\s/g, '').length);
    }, [formData.title, formData.text]);

    useEffect(() => {
        const handleBeforeUnload = (e: BeforeUnloadEvent) => {
            if (isDirty) {
                e.preventDefault();
                e.returnValue = '';
            }
        };
        window.addEventListener('beforeunload', handleBeforeUnload);
        return () => window.removeEventListener('beforeunload', handleBeforeUnload);
    }, [isDirty]);

    const fetchArticle = async (cid: number) => {
        try {
            const response = await contentApi.getContent(cid);
            if (response.data.code === 200 && response.data.data) {
                const data = response.data.data;
                setFormData({
                    title: data.title || '',
                    text: data.text || '',
                    status: data.status || 'draft',
                    type: data.type || 'post',
                    slug: data.slug || '',
                    categoryMid: data.categoryMid || data.category?.mid || 0
                });
            } else {
                message.error('获取文章失败');
                navigate('/articles');
            }
        } catch (error) {
            console.error('获取文章失败:', error);
            message.error('获取文章失败');
            navigate('/articles');
        }
    };

    const handleTitleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const title = e.target.value;
        if (title.length <= 30) {
            setFormData(prev => ({ ...prev, title }));
            setIsDirty(true);
        }
    };

    const handleContentChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        const text = e.target.value;
        if (text.replace(/\s/g, '').length <= 30000) {
            setFormData(prev => ({ ...prev, text }));
            setIsDirty(true);
        }
    };

    const insertMarkdown = (before: string, after: string = '') => {
        const textarea = textAreaRef.current;
        if (!textarea) return;

        const start = textarea.selectionStart;
        const end = textarea.selectionEnd;
        const selectedText = formData.text.substring(start, end);
        const newText = formData.text.substring(0, start) + before + selectedText + after + formData.text.substring(end);

        setFormData(prev => ({ ...prev, text: newText }));
        setIsDirty(true);

        setTimeout(() => {
            textarea.focus();
            const newCursorPos = start + before.length + selectedText.length;
            textarea.setSelectionRange(newCursorPos, newCursorPos);
        }, 0);
    };

    const toolbarItems = [
        { icon: <PictureOutlined />, title: '图片', action: () => insertMarkdown('![图片](', ')') },
        { icon: <SmileOutlined />, title: '表情', action: () => message.info('表情功能开发中') },
        { icon: <TagOutlined />, title: '话题', action: () => insertMarkdown('#话题#') },
        { icon: <BoldOutlined />, title: '粗体', action: () => insertMarkdown('**', '**') },
        { icon: <ItalicOutlined />, title: '斜体', action: () => insertMarkdown('*', '*') },
        { icon: <UnorderedListOutlined />, title: '列表', action: () => insertMarkdown('\n- ') },
        { icon: <OrderedListOutlined />, title: '有序列表', action: () => insertMarkdown('\n1. ') },
        { icon: <CodeOutlined />, title: '代码', action: () => insertMarkdown('```\n', '\n```') },
        { icon: <LinkOutlined />, title: '链接', action: () => insertMarkdown('[', '](url)') },
    ];

    const validateForm = (): boolean => {
        if (!formData.title.trim()) {
            message.error('请输入文章标题');
            return false;
        }
        if (!formData.text.trim()) {
            message.error('请输入文章内容');
            return false;
        }
        if (!formData.categoryMid || formData.categoryMid === 0) {
            message.error('请选择文章分类');
            return false;
        }
        return true;
    };

    const handleSave = async (publish: boolean = false) => {
        if (!validateForm()) return;

        const saveAction = publish ? setPublishing : setSaving;
        saveAction(true);

        try {
            const payload = {
                title: formData.title,
                text: formData.text,
                content: formData.text,
                status: publish ? 'publish' : formData.status,
                type: formData.type,
                categoryMid: formData.categoryMid
            };

            let response;
            if (isEditMode && id) {
                response = await contentApi.updateContent(parseInt(id), payload);
            } else {
                response = await contentApi.createContent(payload);
            }

            if (response.data.code === 200) {
                message.success(publish ? '发布成功' : '保存成功');
                setLastSaved(new Date());
                setIsDirty(false);

                if (!isEditMode && response.data.data?.cid) {
                    navigate(`/editor/${response.data.data.cid}`, { replace: true });
                }
            } else {
                message.error(response.data.msg || '操作失败');
            }
        } catch (error) {
            console.error('保存失败:', error);
            message.error('保存失败');
        } finally {
            saveAction(false);
        }
    };

    const handleDelete = async () => {
        if (!isEditMode || !id) return;

        try {
            const response = await contentApi.deleteContent(parseInt(id));
            if (response.data.code === 200) {
                message.success('删除成功');
                navigate('/articles');
            } else {
                message.error(response.data.msg || '删除失败');
            }
        } catch (error) {
            console.error('删除失败:', error);
            message.error('删除失败');
        }
    };

    const handlePreview = () => {
        if (isEditMode && id) {
            window.open(`/post/${formData.slug || id}`, '_blank');
        } else {
            message.info('请先保存文章');
        }
    };

    const formatTime = (date: Date) => {
        return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' });
    };

    // 移动端布局
    if (isMobile) {
        return (
            <div style={{
                display: 'flex',
                flexDirection: 'column',
                height: '100vh',
                background: '#fff',
                overflow: 'hidden'
            }}>
                {/* 顶部导航栏 */}
                <div style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    padding: '12px 16px',
                    borderBottom: '1px solid #f0f0f0',
                    background: '#fff',
                    flexShrink: 0
                }}>
                    <Button
                        type="text"
                        icon={<CloseOutlined style={{ fontSize: 20 }} />}
                        onClick={() => navigate('/articles')}
                    />
                    <div style={{ display: 'flex', gap: 12 }}>
                        <Button
                            type="text"
                            size="small"
                            onClick={() => handleSave(false)}
                            loading={saving}
                            style={{ color: '#666' }}
                        >
                            草稿箱
                        </Button>
                        <Button
                            type="primary"
                            size="small"
                            onClick={() => handleSave(true)}
                            loading={publishing}
                            style={{
                                background: '#1890ff',
                                borderRadius: 4,
                                padding: '2px 12px'
                            }}
                        >
                            发布
                        </Button>
                    </div>
                </div>

                {/* 编辑区域 */}
                <div style={{
                    flex: 1,
                    display: 'flex',
                    flexDirection: 'column',
                    padding: '16px',
                    paddingBottom: showToolbar ? '140px' : '90px',
                    overflow: 'auto'
                }}>
                    {/* 标题输入 */}
                    <div style={{ marginBottom: 12, flexShrink: 0 }}>
                        <Input
                            placeholder="标题（必填）"
                            value={formData.title}
                            onChange={handleTitleChange}
                            style={{
                                fontSize: 20,
                                fontWeight: 500,
                                border: 'none',
                                padding: 0,
                                boxShadow: 'none'
                            }}
                        />
                        <div style={{
                            textAlign: 'right',
                            fontSize: 12,
                            color: titleCount >= 30 ? '#ff4d4f' : '#999',
                            marginTop: 4
                        }}>
                            {titleCount}/30
                        </div>
                    </div>

                    <Divider style={{ margin: '8px 0', flexShrink: 0 }} />

                    {/* 正文输入区域 */}
                    <div style={{ flex: 1, display: 'flex', flexDirection: 'column', overflow: 'hidden' }}>
                        <TextArea
                            ref={textAreaRef}
                            placeholder="请尽情发挥吧"
                            value={formData.text}
                            onChange={handleContentChange}
                            style={{
                                flex: 1,
                                border: 'none',
                                resize: 'none',
                                fontSize: 16,
                                lineHeight: 1.8,
                                padding: 0,
                                minHeight: 0
                            }}
                        />

                        {/* 字数统计 */}
                        <div style={{
                            textAlign: 'right',
                            fontSize: 12,
                            color: contentCount >= 30000 ? '#ff4d4f' : '#999',
                            marginTop: 8,
                            flexShrink: 0
                        }}>
                            {contentCount}/30000
                        </div>
                    </div>
                </div>

                {/* 底部工具栏 - 固定在屏幕底部 */}
                <div style={{
                    position: 'fixed',
                    bottom: 0,
                    left: 0,
                    right: 0,
                    background: '#fafafa',
                    zIndex: 100,
                    borderTop: '1px solid #f0f0f0'
                }}>
                    <div style={{
                        display: 'flex',
                        justifyContent: 'space-between',
                        alignItems: 'center',
                        padding: '8px 16px'
                    }}>
                        <div style={{ display: 'flex', gap: 20 }}>
                            {toolbarItems.slice(0, 5).map((item, index) => (
                                <div
                                    key={index}
                                    onClick={item.action}
                                    style={{
                                        cursor: 'pointer',
                                        color: '#666',
                                        fontSize: 20
                                    }}
                                >
                                    {item.icon}
                                </div>
                            ))}
                            <div
                                onClick={() => setShowToolbar(!showToolbar)}
                                style={{
                                    cursor: 'pointer',
                                    color: '#666',
                                    fontSize: 20
                                }}
                            >
                                <PlusOutlined />
                            </div>
                        </div>
                        <div
                            onClick={() => setCategoryDrawerOpen(true)}
                            style={{
                                cursor: 'pointer',
                                color: '#666',
                                fontSize: 20
                            }}
                        >
                            <SettingOutlined />
                        </div>
                    </div>

                    {/* 扩展工具栏 */}
                    {showToolbar && (
                        <div style={{
                            display: 'flex',
                            gap: 20,
                            padding: '8px 16px',
                            borderTop: '1px solid #f5f5f5',
                            overflowX: 'auto'
                        }}>
                            {toolbarItems.slice(5).map((item, index) => (
                                <Tooltip key={index} title={item.title}>
                                    <div
                                        onClick={item.action}
                                        style={{
                                            cursor: 'pointer',
                                            color: '#666',
                                            fontSize: 18
                                        }}
                                    >
                                        {item.icon}
                                    </div>
                                </Tooltip>
                            ))}
                        </div>
                    )}
                </div>

                {/* 设置抽屉 */}
                <Drawer
                    title="文章设置"
                    placement="bottom"
                    onClose={() => setCategoryDrawerOpen(false)}
                    open={categoryDrawerOpen}
                    height="50%"
                >
                    <div style={{ marginBottom: 24 }}>
                        <Text type="secondary">文章分类</Text>
                        <div style={{ marginTop: 8, display: 'flex', gap: 12, flexWrap: 'wrap' }}>
                            {categories.length === 0 ? (
                                <Text type="secondary">暂无分类，请先创建分类</Text>
                            ) : (
                                categories.map(cat => (
                                    <Tag
                                        key={cat.mid}
                                        color={formData.categoryMid === cat.mid ? 'blue' : undefined}
                                        onClick={() => {
                                            setFormData(prev => ({ ...prev, categoryMid: cat.mid }));
                                            setIsDirty(true);
                                        }}
                                        style={{ cursor: 'pointer', padding: '4px 12px' }}
                                    >
                                        <FolderOutlined style={{ marginRight: 4 }} />
                                        {cat.name}
                                    </Tag>
                                ))
                            )}
                        </div>
                    </div>
                    <div style={{ marginBottom: 24 }}>
                        <Text type="secondary">文章类型</Text>
                        <div style={{ marginTop: 8, display: 'flex', gap: 12 }}>
                            <Tag
                                color={formData.type === 'post' ? 'blue' : undefined}
                                onClick={() => setFormData(prev => ({ ...prev, type: 'post' }))}
                                style={{ cursor: 'pointer', padding: '4px 12px' }}
                            >
                                文章
                            </Tag>
                            <Tag
                                color={formData.type === 'page' ? 'blue' : undefined}
                                onClick={() => setFormData(prev => ({ ...prev, type: 'page' }))}
                                style={{ cursor: 'pointer', padding: '4px 12px' }}
                            >
                                页面
                            </Tag>
                        </div>
                    </div>
                    <div style={{ marginBottom: 24 }}>
                        <Text type="secondary">发布状态</Text>
                        <div style={{ marginTop: 8, display: 'flex', gap: 12 }}>
                            <Tag
                                color={formData.status === 'draft' ? 'orange' : undefined}
                                onClick={() => setFormData(prev => ({ ...prev, status: 'draft' }))}
                                style={{ cursor: 'pointer', padding: '4px 12px' }}
                            >
                                草稿
                            </Tag>
                            <Tag
                                color={formData.status === 'publish' ? 'green' : undefined}
                                onClick={() => setFormData(prev => ({ ...prev, status: 'publish' }))}
                                style={{ cursor: 'pointer', padding: '4px 12px' }}
                            >
                                发布
                            </Tag>
                            <Tag
                                color={formData.status === 'private' ? 'purple' : undefined}
                                onClick={() => setFormData(prev => ({ ...prev, status: 'private' }))}
                                style={{ cursor: 'pointer', padding: '4px 12px' }}
                            >
                                私密
                            </Tag>
                        </div>
                    </div>
                    {isEditMode && (
                        <div>
                            <Popconfirm
                                title="确定删除这篇文章吗？"
                                description="删除后无法恢复"
                                onConfirm={handleDelete}
                                okText="确定"
                                cancelText="取消"
                                okButtonProps={{ danger: true }}
                            >
                                <Button danger block icon={<DeleteOutlined />}>
                                    删除文章
                                </Button>
                            </Popconfirm>
                        </div>
                    )}
                </Drawer>
            </div>
        );
    }

    // 桌面端布局
    return (
        <div style={{
            display: 'flex',
            flexDirection: 'column',
            height: '100vh',
            background: '#f5f5f5'
        }}>
            {/* 顶部工具栏 */}
            <div style={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                padding: '12px 24px',
                background: '#fff',
                borderBottom: '1px solid #e8e8e8',
                flexShrink: 0
            }}>
                <div style={{ display: 'flex', alignItems: 'center', gap: 16 }}>
                    <Button
                        icon={<CloseOutlined />}
                        onClick={() => navigate('/articles')}
                    />
                    <Text strong style={{ fontSize: 16 }}>
                        {isEditMode ? '编辑文章' : '新建文章'}
                    </Text>
                    {isDirty && <Tag color="orange">未保存</Tag>}
                </div>

                <div style={{ display: 'flex', alignItems: 'center', gap: 12 }}>
                    {lastSaved && (
                        <Text type="secondary" style={{ fontSize: 12 }}>
                            上次保存: {formatTime(lastSaved)}
                        </Text>
                    )}
                    <Text type="secondary" style={{ fontSize: 12 }}>
                        {contentCount} 字
                    </Text>
                    <Button
                        icon={<EyeOutlined />}
                        onClick={handlePreview}
                        disabled={!isEditMode}
                    >
                        预览
                    </Button>
                    <Button
                        icon={<SaveOutlined />}
                        loading={saving}
                        onClick={() => handleSave(false)}
                    >
                        保存草稿
                    </Button>
                    <Button
                        type="primary"
                        icon={<SendOutlined />}
                        loading={publishing}
                        onClick={() => handleSave(true)}
                    >
                        {formData.status === 'publish' ? '更新发布' : '立即发布'}
                    </Button>
                    {isEditMode && (
                        <Popconfirm
                            title="确定删除这篇文章吗？"
                            description="删除后无法恢复"
                            onConfirm={handleDelete}
                            okText="确定"
                            cancelText="取消"
                            okButtonProps={{ danger: true }}
                        >
                            <Button danger icon={<DeleteOutlined />}>
                                删除
                            </Button>
                        </Popconfirm>
                    )}
                </div>
            </div>

            {/* Markdown 工具栏 */}
            <div style={{
                display: 'flex',
                padding: '8px 24px',
                background: '#fff',
                borderBottom: '1px solid #e8e8e8',
                gap: 4,
                flexShrink: 0
            }}>
                {toolbarItems.map((item, index) => (
                    <Tooltip key={index} title={item.title}>
                        <Button
                            type="text"
                            size="small"
                            icon={item.icon}
                            onClick={item.action}
                        />
                    </Tooltip>
                ))}
            </div>

            {/* 编辑区域 */}
            <div style={{
                flex: 1,
                display: 'flex',
                padding: '24px',
                overflow: 'auto',
                gap: 24
            }}>
                <div style={{
                    flex: 1,
                    maxWidth: 900,
                    background: '#fff',
                    borderRadius: 8,
                    boxShadow: '0 2px 8px rgba(0, 0, 0, 0.06)',
                    display: 'flex',
                    flexDirection: 'column',
                    minHeight: 0
                }}>
                    {/* 标题输入 */}
                    <div style={{ padding: '32px 40px 16px' }}>
                        <Input
                            placeholder="请输入文章标题..."
                            value={formData.title}
                            onChange={handleTitleChange}
                            style={{
                                fontSize: 32,
                                fontWeight: 600,
                                border: 'none',
                                padding: 0,
                                boxShadow: 'none'
                            }}
                        />
                    </div>

                    <Divider style={{ margin: 0 }} />

                    {/* 正文输入 */}
                    <div style={{ flex: 1, display: 'flex', flexDirection: 'column', minHeight: 0 }}>
                        <TextArea
                            ref={textAreaRef}
                            value={formData.text}
                            onChange={handleContentChange}
                            placeholder="开始写作..."
                            style={{
                                flex: 1,
                                border: 'none',
                                resize: 'none',
                                fontSize: 16,
                                lineHeight: 1.8,
                                padding: '24px 40px',
                                minHeight: 400
                            }}
                        />
                    </div>

                    {/* 底部信息栏 */}
                    <div style={{
                        display: 'flex',
                        justifyContent: 'space-between',
                        alignItems: 'center',
                        padding: '12px 40px',
                        borderTop: '1px solid #f0f0f0',
                        background: '#fafafa',
                        borderRadius: '0 0 8px 8px'
                    }}>
                        <Text type="secondary" style={{ fontSize: 12 }}>
                            支持 Markdown 语法
                        </Text>
                        <Text type="secondary" style={{ fontSize: 12 }}>
                            快捷键: Ctrl+S 保存 | {contentCount} 字
                        </Text>
                    </div>
                </div>

                {/* 右侧设置面板 */}
                <div style={{
                    width: 280,
                    flexShrink: 0,
                    display: 'flex',
                    flexDirection: 'column',
                    gap: 16
                }}>
                    {/* 分类选择 */}
                    <div style={{
                        background: '#fff',
                        borderRadius: 8,
                        padding: 16,
                        boxShadow: '0 2px 8px rgba(0, 0, 0, 0.06)'
                    }}>
                        <Text strong style={{ display: 'block', marginBottom: 12 }}>
                            <FolderOutlined style={{ marginRight: 8 }} />
                            文章分类
                        </Text>
                        <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8 }}>
                            {categories.length === 0 ? (
                                <Text type="secondary" style={{ fontSize: 12 }}>
                                    暂无分类，请先创建
                                </Text>
                            ) : (
                                categories.map(cat => (
                                    <Tag
                                        key={cat.mid}
                                        color={formData.categoryMid === cat.mid ? 'blue' : undefined}
                                        onClick={() => {
                                            setFormData(prev => ({ ...prev, categoryMid: cat.mid }));
                                            setIsDirty(true);
                                        }}
                                        style={{ cursor: 'pointer', padding: '4px 12px' }}
                                    >
                                        {cat.name}
                                    </Tag>
                                ))
                            )}
                        </div>
                    </div>

                    {/* 文章类型 */}
                    <div style={{
                        background: '#fff',
                        borderRadius: 8,
                        padding: 16,
                        boxShadow: '0 2px 8px rgba(0, 0, 0, 0.06)'
                    }}>
                        <Text strong style={{ display: 'block', marginBottom: 12 }}>
                            文章类型
                        </Text>
                        <div style={{ display: 'flex', gap: 8 }}>
                            <Tag
                                color={formData.type === 'post' ? 'blue' : undefined}
                                onClick={() => setFormData(prev => ({ ...prev, type: 'post' }))}
                                style={{ cursor: 'pointer', padding: '4px 12px' }}
                            >
                                文章
                            </Tag>
                            <Tag
                                color={formData.type === 'page' ? 'blue' : undefined}
                                onClick={() => setFormData(prev => ({ ...prev, type: 'page' }))}
                                style={{ cursor: 'pointer', padding: '4px 12px' }}
                            >
                                页面
                            </Tag>
                        </div>
                    </div>

                    {/* 发布状态 */}
                    <div style={{
                        background: '#fff',
                        borderRadius: 8,
                        padding: 16,
                        boxShadow: '0 2px 8px rgba(0, 0, 0, 0.06)'
                    }}>
                        <Text strong style={{ display: 'block', marginBottom: 12 }}>
                            发布状态
                        </Text>
                        <div style={{ display: 'flex', flexWrap: 'wrap', gap: 8 }}>
                            <Tag
                                color={formData.status === 'draft' ? 'orange' : undefined}
                                onClick={() => setFormData(prev => ({ ...prev, status: 'draft' }))}
                                style={{ cursor: 'pointer', padding: '4px 12px' }}
                            >
                                草稿
                            </Tag>
                            <Tag
                                color={formData.status === 'publish' ? 'green' : undefined}
                                onClick={() => setFormData(prev => ({ ...prev, status: 'publish' }))}
                                style={{ cursor: 'pointer', padding: '4px 12px' }}
                            >
                                发布
                            </Tag>
                            <Tag
                                color={formData.status === 'private' ? 'purple' : undefined}
                                onClick={() => setFormData(prev => ({ ...prev, status: 'private' }))}
                                style={{ cursor: 'pointer', padding: '4px 12px' }}
                            >
                                私密
                            </Tag>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default ArticleEditor;