<script setup lang="ts">
const props = defineProps<{ compact?: boolean }>()
const emit = defineEmits<{ select: [key: string] }>()

const route = useRoute()

const items = [
	{ value: '/', icon: 'layout-dashboard', label: '仪表盘' },
	{ value: '/users', icon: 'users', label: '用户管理' },
	{ value: '/articles', icon: 'file-text', label: '内容管理' },
	{ value: '/categories', icon: 'folder-tree', label: '分类管理' },
	{ value: '/comments', icon: 'message-square', label: '评论管理' },
	{ value: '/reviews', icon: 'shield-check', label: '审核管理' },
	{ value: '/settings', icon: 'settings-2', label: '系统设置' },
]

function handleChange(event: Event) {
	const value = (event as CustomEvent<{ value?: string }>).detail?.value
	if (value) emit('select', value)
}
</script>

<template>
	<nav class="nav" :class="{ 'is-compact': props.compact }">
		<mao-menu :value="route.path" orientation="vertical" label="主导航" @mt-change="handleChange">
			<mao-menu-item
				v-for="item in items"
				:key="item.value"
				:value="item.value"
				:icon="item.icon"
				:title="props.compact ? item.label : undefined"
			>
				<span class="nav-label">{{ item.label }}</span>
			</mao-menu-item>
		</mao-menu>
	</nav>
</template>

<style scoped>
.nav {
	padding: 12px 8px;
}

.nav.is-compact {
	padding: 12px 6px;
}

.nav.is-compact .nav-label {
	display: none;
}

.nav.is-compact mao-menu-item {
	justify-content: center;
	padding: 9px 0;
}
</style>
