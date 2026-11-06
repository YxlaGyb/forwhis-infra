<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { theme } from 'ant-design-vue'
import { useMediaQuery } from '@vueuse/core'
import { LogoutOutlined, UserOutlined } from '@ant-design/icons-vue'

const route = useRoute()
const { displayName, identity, logout } = useAuth()
const { token } = theme.useToken()

const isMobile = useMediaQuery('(max-width: 767px)')
const collapsed = ref(false)
const drawerVisible = ref(false)

watch(
	isMobile,
	(mobile) => {
		collapsed.value = mobile
		if (!mobile) drawerVisible.value = false
	},
	{ immediate: true },
)

const isEditorPage = computed(() => route.path.startsWith('/editor'))
const siderWidth = computed(() => (collapsed.value ? 80 : 200))

const colorBgContainer = computed(() => token.value.colorBgContainer)
const borderRadiusLG = computed(() => token.value.borderRadiusLG)

const siderStyle = computed(() => ({
	background: colorBgContainer.value,
	boxShadow: '2px 0 8px rgba(0, 0, 0, 0.1)',
	position: 'fixed' as const,
	left: 0,
	top: 0,
	bottom: 0,
	zIndex: 100,
	overflow: 'auto' as const,
}))

const headerStyle = computed(() => ({
	padding: '0 24px',
	background: colorBgContainer.value,
	display: 'flex',
	alignItems: 'center',
	justifyContent: 'space-between',
	boxShadow: '0 2px 8px rgba(0, 0, 0, 0.1)',
	position: 'fixed' as const,
	top: 0,
	right: 0,
	left: `${isMobile.value ? 0 : siderWidth.value}px`,
	zIndex: 99,
}))

const contentStyle = computed(() => {
	if (isEditorPage.value) {
		return {
			margin: 0,
			padding: 0,
			background: 'transparent',
			borderRadius: 0,
			minHeight: '280px',
			maxWidth: '100%',
			overflow: 'hidden' as const,
		}
	}
	return {
		margin: isMobile.value ? '76px 0 0' : '24px',
		marginTop: isMobile.value ? '76px' : '88px',
		padding: isMobile.value ? '12px' : '24px',
		background: colorBgContainer.value,
		borderRadius: isMobile.value ? 0 : `${borderRadiusLG.value}px`,
		minHeight: '280px',
		maxWidth: '100%',
		overflow: 'hidden' as const,
	}
})

function handleMenuClick(key: string) {
	void navigateTo(key)
	if (isMobile.value) drawerVisible.value = false
}

async function handleLogout() {
	await logout()
	await navigateTo('/login')
}

function handleUserMenuClick(info: { key: string | number }) {
	if (String(info.key) === 'logout') void handleLogout()
}
</script>

<template>
	<a-layout :style="{ minHeight: '100vh' }">
		<a-drawer
			v-if="isMobile"
			placement="left"
			:open="drawerVisible"
			:width="200"
			@close="drawerVisible = false"
		>
			<template #title>
				<span style="color: #1890ff; font-weight: 600">Lmoadll</span>
			</template>
			<SideMenu @select="handleMenuClick" />
		</a-drawer>

		<a-layout-sider
			v-else
			:collapsed="collapsed"
			collapsible
			:trigger="null"
			:width="200"
			:collapsed-width="80"
			:style="siderStyle"
		>
			<div class="layout-brand">
				<span :style="{ fontSize: collapsed ? '14px' : '16px' }">
					{{ collapsed ? 'L' : 'Lmoadll' }}
				</span>
			</div>
			<SideMenu style="margin-top: 16px" @select="handleMenuClick" />
		</a-layout-sider>

		<a-layout :style="{ marginLeft: isMobile ? '0' : `${siderWidth}px` }">
			<a-layout-header v-if="!isEditorPage" :style="headerStyle">
				<div class="layout-header-left">
					<MenuOutlined v-if="isMobile" class="layout-trigger" @click="drawerVisible = true" />
					<MenuUnfoldOutlined
						v-else-if="collapsed"
						class="layout-trigger"
						@click="collapsed = false"
					/>
					<MenuFoldOutlined v-else class="layout-trigger" @click="collapsed = true" />
				</div>

				<a-space :size="16">
					<BellOutlined class="layout-bell" />
					<a-dropdown placement="bottomRight" arrow>
						<a-space class="layout-user">
							<a-avatar size="small">
								<template #icon><UserOutlined /></template>
							</a-avatar>
							<span>{{ displayName }}</span>
							<span class="layout-identity">{{ identity }}</span>
						</a-space>
						<template #overlay>
							<a-menu @click="handleUserMenuClick">
								<a-menu-item key="profile">
									<UserOutlined />
									个人信息
								</a-menu-item>
								<a-menu-item key="logout">
									<LogoutOutlined />
									退出登录
								</a-menu-item>
							</a-menu>
						</template>
					</a-dropdown>
				</a-space>
			</a-layout-header>

			<a-layout-content :style="contentStyle">
				<slot />
			</a-layout-content>
		</a-layout>
	</a-layout>
</template>

<style scoped>
.layout-brand {
	height: 64px;
	display: flex;
	align-items: center;
	justify-content: center;
	border-bottom: 1px solid #f0f0f0;
}

.layout-brand span {
	font-weight: 600;
	color: #1890ff;
}

.layout-header-left {
	display: flex;
	align-items: center;
}

.layout-trigger {
	font-size: 18px;
	cursor: pointer;
}

.layout-bell {
	font-size: 18px;
	cursor: pointer;
}

.layout-user {
	cursor: pointer;
}

.layout-identity {
	color: #94a3b8;
	font-size: 12px;
}
</style>
