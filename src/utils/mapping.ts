export const getStatusColor = (status: string): string => {
    const colorMap: Record<string, string> = {
        'publish': 'green',
        'draft': 'orange',
        'private': 'blue',
        'trash': 'red',
        'approved': 'green',
        'pending': 'orange',
        'spam': 'red'
    };
    return colorMap[status] || 'default';
};

export const getStatusName = (status: string): string => {
    const nameMap: Record<string, string> = {
        'publish': '已发布',
        'draft': '草稿',
        'private': '私密',
        'trash': '回收站',
        'approved': '已批准',
        'pending': '待审核',
        'spam': '垃圾'
    };
    return nameMap[status] || status;
};

export const getTypeName = (type: string): string => {
    const nameMap: Record<string, string> = {
        'post': '文章',
        'page': '页面'
    };
    return nameMap[type] || type;
};

export const getRoleColor = (role: string): string => {
    const colorMap: Record<string, string> = {
        'superadmin': 'red',
        'admin': 'orange',
        'editor': 'blue',
        'user': 'green',
        'visitor': 'default'
    };
    return colorMap[role] || 'default';
};

export const getRoleName = (role: string): string => {
    const nameMap: Record<string, string> = {
        'superadmin': '超级管理员',
        'admin': '管理员',
        'editor': '编辑',
        'user': '用户',
        'visitor': '访客'
    };
    return nameMap[role] || role;
};
