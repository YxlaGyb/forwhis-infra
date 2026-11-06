<script setup lang="ts">
definePageMeta({ layout: 'blank' })

const { logout } = useAuth()

async function handleLogout() {
	await logout()
	await navigateTo('/login')
}
</script>

<template>
	<div class="no-access">
		<p>您没有管理员权限</p>
		<a-button @click="handleLogout">退出登录</a-button>
	</div>
</template>

<style scoped>
.no-access {
	min-height: 100vh;
	display: flex;
	flex-direction: column;
	align-items: center;
	justify-content: center;
	gap: 16px;
}
</style>
