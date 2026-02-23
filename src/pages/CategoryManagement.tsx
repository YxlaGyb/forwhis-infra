import React, { useState, useEffect } from 'react';
import {
    Card,
    Input,
    Button,
    Space,
    Typography,
    message,
    Modal,
    Form,
    Tree,
    Popconfirm,
    Tooltip
} from 'antd';
import {
    SearchOutlined,
    PlusOutlined,
    EditOutlined,
    DeleteOutlined,
    FolderOutlined,
    FolderOpenOutlined
} from '@ant-design/icons';
import { metaApi } from '../api';
import { useSearchFilter } from '../hooks/useSearchFilter';
import type { TreeDataNode } from 'antd';

const { Title, Text } = Typography;

interface Category {
    mid: number;
    name: string;
    slug: string;
    type: string;
    description?: string;
    parent?: number;
    count?: number;
    children?: Category[];
}

const CategoryManagement: React.FC = () => {
    const [categories, setCategories] = useState<Category[]>([]);
    const [loading, setLoading] = useState(false);
    const [isModalVisible, setIsModalVisible] = useState(false);
    const [editingCategory, setEditingCategory] = useState<Category | null>(null);
    const [form] = Form.useForm();
    const [expandedKeys, setExpandedKeys] = useState<React.Key[]>([]);

    const { filteredData: filteredCategories, handleSearch } = useSearchFilter<Category>({
        data: categories,
        searchKeys: ['name', 'slug', 'description']
    });

    useEffect(() => {
        fetchCategories();
    }, []);

    const fetchCategories = async () => {
        setLoading(true);
        try {
            const response = await metaApi.getMetas({ type: 'category' });
            if (response.data.code === 200 && response.data.data) {
                const items = response.data.data || [];
                const treeData = buildTree(items);
                setCategories(treeData);
                setExpandedKeys(treeData.map(c => c.mid));
            } else {
                message.error(response.data.msg || '获取分类列表失败');
            }
        } catch (error) {
            console.error('获取分类列表失败:', error);
            message.error('获取分类列表失败');
        } finally {
            setLoading(false);
        }
    };

    const buildTree = (items: Category[]): Category[] => {
        const map = new Map<number, Category>();
        const roots: Category[] = [];

        items.forEach(item => {
            map.set(item.mid, { ...item, children: [] });
        });

        items.forEach(item => {
            const node = map.get(item.mid)!;
            if (item.parent && map.has(item.parent)) {
                const parent = map.get(item.parent)!;
                if (!parent.children) parent.children = [];
                parent.children.push(node);
            } else {
                roots.push(node);
            }
        });

        return roots;
    };

    const flattenTree = (items: Category[]): Category[] => {
        const result: Category[] = [];
        const traverse = (nodes: Category[]) => {
            nodes.forEach(node => {
                result.push(node);
                if (node.children && node.children.length > 0) {
                    traverse(node.children);
                }
            });
        };
        traverse(items);
        return result;
    };

    const handleAdd = () => {
        setEditingCategory(null);
        form.resetFields();
        setIsModalVisible(true);
    };

    const handleEdit = (record: Category) => {
        setEditingCategory(record);
        form.setFieldsValue({
            name: record.name,
            slug: record.slug,
            description: record.description,
            parent: record.parent
        });
        setIsModalVisible(true);
    };

    const handleDelete = async (mid: number) => {
        try {
            const response = await metaApi.deleteMeta(mid);
            if (response.data.code === 200) {
                message.success('删除成功');
                fetchCategories();
            } else {
                message.error(response.data.msg || '删除失败');
            }
        } catch (error) {
            console.error('删除失败:', error);
            message.error('删除失败');
        }
    };

    const handleModalOk = async () => {
        try {
            const values = await form.validateFields();

            if (editingCategory) {
                const response = await metaApi.updateMeta(editingCategory.mid, values);
                if (response.data.code === 200) {
                    message.success('更新成功');
                    setIsModalVisible(false);
                    fetchCategories();
                } else {
                    message.error(response.data.msg || '更新失败');
                }
            } else {
                const response = await metaApi.createMeta({ ...values, type: 'category' });
                if (response.data.code === 200) {
                    message.success('创建成功');
                    setIsModalVisible(false);
                    fetchCategories();
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
        setEditingCategory(null);
    };

    const convertToTreeData = (items: Category[]): TreeDataNode[] => {
        return items.map(item => ({
            key: item.mid,
            title: (
                <div style={{
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    padding: '8px 0'
                }}>
                    <Space>
                        <Text strong>{item.name}</Text>
                        <Text type="secondary" style={{ fontSize: 12 }}>
                            /{item.slug}
                        </Text>
                        {item.count !== undefined && item.count > 0 && (
                            <Text type="secondary" style={{ fontSize: 12 }}>
                                ({item.count} 篇文章)
                            </Text>
                        )}
                    </Space>
                    <Space size="small">
                        <Tooltip title="编辑">
                            <Button
                                type="text"
                                size="small"
                                icon={<EditOutlined />}
                                onClick={(e) => {
                                    e.stopPropagation();
                                    handleEdit(item);
                                }}
                            />
                        </Tooltip>
                        <Popconfirm
                            title="确定要删除这个分类吗？"
                            description="删除后无法恢复，相关文章将失去分类"
                            onConfirm={(e) => {
                                e?.stopPropagation();
                                handleDelete(item.mid);
                            }}
                            okText="确定"
                            cancelText="取消"
                        >
                            <Button
                                type="text"
                                size="small"
                                danger
                                icon={<DeleteOutlined />}
                                onClick={(e) => e.stopPropagation()}
                            />
                        </Popconfirm>
                    </Space>
                </div>
            ),
            children: item.children ? convertToTreeData(item.children) : undefined
        }));
    };

    const allCategories = flattenTree(categories);

    return (
        <div>
            <Card>
                <div style={{
                    display: 'flex',
                    justifyContent: 'space-between',
                    alignItems: 'center',
                    marginBottom: 16,
                    flexWrap: 'wrap',
                    gap: 12
                }}>
                    <Title level={4} style={{ margin: 0 }}>
                        <FolderOutlined style={{ marginRight: 8 }} />
                        分类管理
                    </Title>
                    <Button
                        type="primary"
                        icon={<PlusOutlined />}
                        onClick={handleAdd}
                    >
                        新增分类
                    </Button>
                </div>

                <Space style={{ marginBottom: 16 }} wrap>
                    <Input
                        placeholder="搜索分类名称"
                        prefix={<SearchOutlined />}
                        onChange={(e) => handleSearch(e.target.value)}
                        style={{ width: 300 }}
                        allowClear
                    />
                </Space>

                <Tree
                    treeData={convertToTreeData(filteredCategories)}
                    expandedKeys={expandedKeys}
                    onExpand={(keys) => setExpandedKeys(keys)}
                    showIcon
                    icon={({ expanded }) =>
                        expanded ? <FolderOpenOutlined /> : <FolderOutlined />
                    }
                    style={{ background: 'transparent' }}
                />
            </Card>

            <Modal
                title={editingCategory ? '编辑分类' : '新增分类'}
                open={isModalVisible}
                onOk={handleModalOk}
                onCancel={handleModalCancel}
                width={500}
            >
                <Form
                    form={form}
                    layout="vertical"
                >
                    <Form.Item
                        label="分类名称"
                        name="name"
                        rules={[{ required: true, message: '请输入分类名称' }]}
                    >
                        <Input placeholder="例如：技术文章" />
                    </Form.Item>

                    <Form.Item
                        label="别名 (Slug)"
                        name="slug"
                        rules={[
                            { required: true, message: '请输入别名' },
                            { pattern: /^[a-z0-9-]+$/, message: '只能使用小写字母、数字和连字符' }
                        ]}
                    >
                        <Input placeholder="例如：tech-articles" />
                    </Form.Item>

                    <Form.Item
                        label="父分类"
                        name="parent"
                    >
                        <Input.Search
                            placeholder="输入分类ID，留空为顶级分类"
                            allowClear
                            disabled
                        />
                    </Form.Item>

                    <Form.Item
                        label="描述"
                        name="description"
                    >
                        <Input.TextArea
                            rows={3}
                            placeholder="分类的描述信息"
                            maxLength={200}
                            showCount
                        />
                    </Form.Item>
                </Form>
            </Modal>
        </div>
    );
};

export default CategoryManagement;
