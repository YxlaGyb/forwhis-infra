<script setup lang="ts">
import { useMediaQuery } from '@vueuse/core'

const route = useRoute()
const { displayName, identity, logout } = useAuth()
const { theme, toggle: toggleTheme } = useTheme()

const isMobile = useMediaQuery('(max-width: 767px)')
const collapsed = ref(false)
const drawerOpen = ref(false)

watch(
	isMobile,
	(mobile) => {
		collapsed.value = mobile
		if (!mobile) drawerOpen.value = false
	},
	{ immediate: true },
)

const isEditorPage = computed(() => route.path.startsWith('/editor'))

const avatarText = computed(() => displayName.value.slice(0, 1).toUpperCase())

function handleMenuSelect(key: string) {
	void navigateTo(key)
	drawerOpen.value = false
}

async function handleLogout() {
	await logout()
	await navigateTo('/login')
}
</script>

<template>
	<div
		class="shell"
		:class="{
			'is-compact': collapsed && !isMobile,
			'is-mobile': isMobile,
			'is-editor': isEditorPage,
		}"
	>
		<aside v-if="!isMobile" class="sider">
			<div class="brand">
				<mao-icon name="feather" :size="20" />
				<span v-show="!collapsed" class="brand-name">forwhis</span>
			</div>
			<div class="nav-wrap">
				<SideMenu :compact="collapsed" @select="handleMenuSelect" />
			</div>
		</aside>

		<mao-drawer
			v-if="isMobile"
			:open="drawerOpen || undefined"
			placement="start"
			size="sm"
			heading="forwhis"
			@mt-close="drawerOpen = false"
		>
			<SideMenu @select="handleMenuSelect" />
		</mao-drawer>

		<div class="main">
			<header v-if="!isEditorPage" class="topbar">
				<mao-button
					variant="flat"
					size="sm"
					aria-label="切换导航"
					@click="isMobile ? (drawerOpen = true) : (collapsed = !collapsed)"
				>
					<mao-icon
						:name="isMobile ? 'menu' : collapsed ? 'panel-left-open' : 'panel-left-close'"
						:size="18"
					/>
				</mao-button>

				<div class="topbar-end">
					<mao-button
						variant="flat"
						size="sm"
						:aria-label="theme === 'dark' ? '切换到浅色主题' : '切换到深色主题'"
						@click="toggleTheme"
					>
						<mao-icon :name="theme === 'dark' ? 'sun' : 'moon'" :size="18" />
					</mao-button>

					<mao-dropdown placement="bottom-end">
						<mao-button slot="trigger" variant="flat" size="sm">
							<span class="user-trigger">
								<mao-avatar :label="avatarText" size="sm" tone="accent" />
								<span class="user-name">{{ displayName }}</span>
								<span class="user-identity">{{ identity }}</span>
								<mao-icon name="chevron-down" :size="14" />
							</span>
						</mao-button>
						<mao-menu-item icon="log-out" @click="handleLogout">退出登录</mao-menu-item>
					</mao-dropdown>
				</div>
			</header>

			<main class="content">
				<slot />
			</main>
		</div>
	</div>
</template>

<style scoped>
.shell {
	--sider-w: 200px;
	display: flex;
	min-height: 100vh;
	background: var(--mt-bg);
	color: var(--mt-text);
}

.shell.is-compact {
	--sider-w: 72px;
}

.shell.is-mobile {
	--sider-w: 0px;
}

.sider {
	position: fixed;
	top: 0;
	bottom: 0;
	left: 0;
	z-index: 100;
	display: flex;
	flex-direction: column;
	width: var(--sider-w);
	background: var(--mt-surface);
	border-right: 1px solid var(--mt-border);
	transition: width var(--mt-dur) var(--mt-ease);
}

.brand {
	display: flex;
	align-items: center;
	gap: 10px;
	height: 56px;
	padding: 0 18px;
	border-bottom: 1px solid var(--mt-border);
	color: var(--mt-accent);
	overflow: hidden;
}

.brand-name {
	font-weight: 600;
	font-size: var(--mt-fs-md);
	letter-spacing: 0.02em;
	white-space: nowrap;
}

.nav-wrap {
	flex: 1;
	overflow-y: auto;
	overflow-x: hidden;
}

.main {
	flex: 1;
	min-width: 0;
	display: flex;
	flex-direction: column;
	margin-left: var(--sider-w);
	transition: margin-left var(--mt-dur) var(--mt-ease);
}

.topbar {
	position: sticky;
	top: 0;
	z-index: 90;
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: var(--mt-space-2);
	height: 56px;
	padding: 0 16px;
	background: var(--mt-surface);
	border-bottom: 1px solid var(--mt-border);
}

.topbar-end {
	display: flex;
	align-items: center;
	gap: var(--mt-space-2);
}

.user-trigger {
	display: inline-flex;
	align-items: center;
	gap: var(--mt-space-2);
}

.user-name {
	color: var(--mt-text);
	font-size: var(--mt-fs-sm);
	font-weight: 500;
}

.user-identity {
	color: var(--mt-muted);
	font-size: var(--mt-fs-xs);
}

.content {
	flex: 1;
	min-width: 0;
	padding: 20px 24px 32px;
}

.shell.is-mobile .content {
	padding: 12px 12px 24px;
}

.shell.is-editor .content {
	padding: 0;
}

@media (max-width: 767px) {
	.user-name,
	.user-identity {
		display: none;
	}
}
</style>
