<script setup lang="ts">
import { computed, h } from 'vue'
import { SafetyOutlined, UserOutlined } from '@ant-design/icons-vue'

const route = useRoute()

const value = computed(() => (route.path === '/roles' ? 'roles' : 'users'))

const options = [
	{ value: 'users', label: h('span', [h(UserOutlined), ' 用户管理']) },
	{ value: 'roles', label: h('span', [h(SafetyOutlined), ' 角色管理 Bata']) },
]

function handleChange(next: string | number) {
	void navigateTo(next === 'roles' ? '/roles' : '/users')
}
</script>

<template>
	<div class="user-role-tabs">
		<a-segmented :value="value" :options="options" @change="handleChange" />
	</div>
</template>

<style scoped>
.user-role-tabs {
	margin-bottom: 16px;
}
</style>
