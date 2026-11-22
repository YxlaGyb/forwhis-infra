<script setup lang="ts">
const route = useRoute()

const value = computed(() => (route.path === '/roles' ? 'roles' : 'users'))

function handleChange(event: Event) {
	const next = (event as CustomEvent<{ value?: string }>).detail?.value
	if (!next) return
	void navigateTo(next === 'roles' ? '/roles' : '/users')
}
</script>

<template>
	<div class="user-role-tabs">
		<mao-segmented :value="value" label="用户与角色" @mt-change="handleChange">
			<mao-segment value="users">
				<mao-icon name="users" :size="15" />
				用户管理
			</mao-segment>
			<mao-segment value="roles">
				<mao-icon name="shield-check" :size="15" />
				角色管理 Beta
			</mao-segment>
		</mao-segmented>
	</div>
</template>

<style scoped>
.user-role-tabs {
	margin-bottom: 16px;
}
</style>
