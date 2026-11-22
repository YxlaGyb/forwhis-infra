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
		<mao-card id="no-access-card" padding="lg">
			<mao-empty icon="shield-alert" title="您没有管理员权限" description="请使用管理员账号登录，或联系站点负责人开通权限。" />
			<mao-button color="primary" variant="bordered" @click="handleLogout">
				<mao-icon name="log-out" :size="15" />
				退出登录
			</mao-button>
		</mao-card>
	</div>
</template>

<style scoped>
.no-access {
	min-height: 100vh;
	display: flex;
	align-items: center;
	justify-content: center;
	padding: var(--mt-space-4);
	background: var(--mt-bg);
}

#no-access-card {
	width: min(420px, 100%);
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: var(--mt-space-4);
}
</style>
