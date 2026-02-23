import React, { useState } from 'react';
import { Card, Button, Dropdown, Empty, Pagination } from 'antd';
import { MoreOutlined } from '@ant-design/icons';
import type { TableProps } from 'antd/es/table';
import './styles.css';

export interface ResponsiveColumn<T = any> {
    title: string;
    dataIndex?: string;
    key: string;
    width?: number | string;
    render?: (value: any, record: T, index: number) => React.ReactNode;
    responsive?: 'desktop' | 'mobile' | 'all';
    priority?: 'high' | 'medium' | 'low';
    ellipsis?: boolean;
}

export interface ResponsiveTableProps<T> extends Omit<TableProps<T>, 'columns'> {
    columns: ResponsiveColumn<T>[];
    dataSource: T[];
    rowKey: string | ((record: T) => string);
    loading?: boolean;
    actions?: (record: T) => ActionItem[];
    cardTitle?: (record: T) => React.ReactNode;
    cardSubtitle?: (record: T) => React.ReactNode;
    emptyText?: string;
}

export interface ActionItem {
    key: string;
    text: string;
    icon?: React.ReactNode;
    danger?: boolean;
    onClick: () => void;
    hidden?: boolean;
}

export function ResponsiveTable<T extends Record<string, any>>({
    columns,
    dataSource,
    rowKey,
    loading: _loading,
    actions,
    cardTitle,
    cardSubtitle,
    emptyText = '暂无数据',
    pagination,
}: ResponsiveTableProps<T>) {
    const [currentPage, setCurrentPage] = useState(1);
    const [pageSize, setPageSize] = useState(10);

    // 获取当前页的数据
    const getCurrentPageData = () => {
        if (!pagination) return dataSource;
        const start = (currentPage - 1) * pageSize;
        const end = start + pageSize;
        return dataSource.slice(start, end);
    };

    // 处理分页变化
    const handlePageChange = (page: number, size?: number) => {
        setCurrentPage(page);
        if (size && size !== pageSize) {
            setPageSize(size);
        }
    };

    const getRowKey = (record: T): string => {
        if (typeof rowKey === 'function') {
            return rowKey(record);
        }
        return record[rowKey] as string;
    };

    // 渲染卡片
    const renderCard = (record: T, index: number) => {
        const key = getRowKey(record);
        const actionItems = actions ? actions(record).filter(item => !item.hidden) : [];

        // 获取高优先级字段显示在卡片上
        const priorityFields = columns.filter(col =>
            col.priority === 'high' || col.priority === 'medium'
        );

        return (
            <Card
                key={key}
                className="responsive-table-card"
                bodyStyle={{ padding: '12px 16px' }}
            >
                <div className="card-header">
                    <div className="card-title-section">
                        {cardTitle && (
                            <div className="card-title">{cardTitle(record)}</div>
                        )}
                        {cardSubtitle && (
                            <div className="card-subtitle">{cardSubtitle(record)}</div>
                        )}
                    </div>
                </div>

                <div className="card-content">
                    {priorityFields.map(col => {
                        const value = col.dataIndex ? record[col.dataIndex] : null;
                        const displayValue = col.render
                            ? col.render(value, record, index)
                            : value;

                        if (displayValue === null || displayValue === undefined || displayValue === '') return null;

                        return (
                            <div key={col.key} className="card-field">
                                <span className="card-field-label" data-field={col.title}>{col.title}:</span>
                                <span className="card-field-value">{displayValue}</span>
                            </div>
                        );
                    })}
                </div>

                {actionItems.length > 0 && (
                    <div className="card-actions">
                        {actionItems.slice(0, 3).map(item => (
                            <Button
                                key={item.key}
                                type={item.danger ? 'default' : 'primary'}
                                size="small"
                                danger={item.danger}
                                icon={item.icon}
                                onClick={item.onClick}
                                className="card-action-btn"
                            >
                                {item.text}
                            </Button>
                        ))}
                        {actionItems.length > 3 && (
                            <Dropdown
                                menu={{
                                    items: actionItems.slice(3).map(item => ({
                                        key: item.key,
                                        label: (
                                            <span style={{ color: item.danger ? '#ff4d4f' : undefined }}>
                                                {item.icon && <span style={{ marginRight: 8 }}>{item.icon}</span>}
                                                {item.text}
                                            </span>
                                        ),
                                        onClick: item.onClick,
                                    }))
                                }}
                                placement="topRight"
                            >
                                <Button
                                    type="default"
                                    size="small"
                                    icon={<MoreOutlined />}
                                    className="card-action-btn"
                                >
                                    更多
                                </Button>
                            </Dropdown>
                        )}
                    </div>
                )}
            </Card>
        );
    };

    const currentData = getCurrentPageData();
    const total = dataSource.length;

    return (
        <div className="responsive-table-container">
            {dataSource.length === 0 ? (
                <Empty description={emptyText} />
            ) : (
                <>
                    <div className="card-list">
                        {currentData.map((record, index) => renderCard(record, index))}
                    </div>
                    {pagination && total > 0 && (
                        <div className="table-pagination">
                            <Pagination
                                current={currentPage}
                                pageSize={pageSize}
                                total={total}
                                onChange={handlePageChange}
                                showSizeChanger={pagination.showSizeChanger}
                                showQuickJumper={pagination.showQuickJumper}
                                showTotal={pagination.showTotal}
                                pageSizeOptions={['10', '20', '50', '100']}
                            />
                        </div>
                    )}
                </>
            )}
        </div>
    );
}

export default ResponsiveTable;
