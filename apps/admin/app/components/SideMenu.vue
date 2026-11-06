<script setup lang="ts">
import { h } from 'vue'
import {
	CommentOutlined,
	DashboardOutlined,
	FileTextOutlined,
	FolderOutlined,
	SafetyOutlined,
	SettingOutlined,
	UserOutlined,
} from '@ant-design/icons-vue'

const emit = defineEmits<{ select: [key: string] }>()

const route = useRoute()

const items = [
	{ key: '/', icon: () => h(DashboardOutlined), label: '仪表盘' },
	{ key: '/users', icon: () => h(UserOutlined), label: '用户管理' },
	{ key: '/articles', icon: () => h(FileTextOutlined), label: '内容管理' },
	{ key: '/categories', icon: () => h(FolderOutlined), label: '分类管理' },
	{ key: '/comments', icon: () => h(CommentOutlined), label: '评论管理' },
	{ key: '/reviews', icon: () => h(SafetyOutlined), label: '审核管理' },
	{ key: '/settings', icon: () => h(SettingOutlined), label: '系统设置' },
]

const selectedKeys = computed(() => [route.path])
</script>

<template>
	<a-menu
		mode="inline"
		:items="items"
		:selected-keys="selectedKeys"
		:style="{ borderRight: 0 }"
		@click="(info: { key: string | number }) => emit('select', String(info.key))"
	/>
</template>
