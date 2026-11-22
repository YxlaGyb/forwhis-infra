<script setup lang="ts">
import { errorMessage } from '~/api/client'

definePageMeta({ layout: 'blank' })

const config = useRuntimeConfig()
const { message } = useNotify()
const { login, canAccessAdmin } = useAuth()

const loading = ref(false)

const backgroundUrl = computed(() => config.public.loginBgUrl || '/110546943_p0.webp')

const RULES = {
	email: '[{"rule":"required","message":"请输入邮箱地址"},{"rule":"email","message":"邮箱格式不正确"}]',
	password: '[{"rule":"required","message":"请输入密码"}]',
}

function readForm(event: Event) {
	const values = (event as CustomEvent<{ values?: Record<string, string> }>).detail?.values ?? {}
	return {
		email: (values.email ?? '').trim(),
		password: values.password ?? '',
	}
}

async function handleSubmit(event: Event) {
	const { email, password } = readForm(event)
	loading.value = true
	try {
		const user = await login(email, password)
		message.success('登录成功')
		await navigateTo(user && canAccessAdmin.value ? '/' : '/no-access')
	} catch (error) {
		message.error(errorMessage(error, '登录失败, 请检查网络连接'))
	} finally {
		loading.value = false
	}
}

function openSite() {
	window.open(config.public.siteUrl, '_blank')
}
</script>

<template>
	<div class="login-shell" :style="{ backgroundImage: `url('${backgroundUrl}')` }">
		<mao-card id="login-card" padding="lg">
			<div class="mt-card-body login-body">
				<div class="login-brand">
					<mao-icon name="feather" :size="22" />
					<span class="login-brand-name">林影心语</span>
					<mao-badge tone="danger">admin</mao-badge>
				</div>

				<mao-form label="管理员登录" @mt-submit="handleSubmit">
					<mao-form-item name="email" label="邮箱" :rules="RULES.email">
						<mao-text-field
							name="email"
							icon="mail"
							placeholder="请输入邮箱地址"
							autocomplete="username"
						/>
					</mao-form-item>

					<mao-form-item name="password" label="密码" :rules="RULES.password">
						<mao-text-field
							name="password"
							type="password"
							icon="lock"
							placeholder="请输入密码"
							autocomplete="current-password"
						/>
					</mao-form-item>

					<mao-button
						color="primary"
						type="submit"
						size="lg"
						block
						:loading="loading || undefined"
					>
						登录
					</mao-button>
				</mao-form>

				<button type="button" class="login-back" @click="openSite">
					<mao-icon name="arrow-left" :size="14" />
					返回首页
				</button>
			</div>
		</mao-card>
	</div>
</template>

<style scoped>
.login-shell {
	min-height: 100vh;
	display: flex;
	align-items: center;
	justify-content: center;
	padding: var(--mt-space-4);
	background-color: var(--mt-bg);
	background-size: cover;
	background-position: center;
	background-repeat: no-repeat;
}

#login-card {
	width: 400px;
	max-width: 100%;
	backdrop-filter: blur(18px);
	background: color-mix(in srgb, var(--mt-surface) 86%, transparent);
	box-shadow: var(--mt-shadow-overlay);
}

.login-body {
	display: flex;
	flex-direction: column;
	gap: var(--mt-space-4);
}

.login-brand {
	display: flex;
	align-items: center;
	justify-content: center;
	gap: var(--mt-space-2);
	color: var(--mt-text);
}

.login-brand-name {
	font-family: var(--mt-font-display);
	font-size: var(--mt-fs-xl);
	font-weight: 600;
	letter-spacing: 0.02em;
}

.login-back {
	display: inline-flex;
	align-items: center;
	justify-content: center;
	gap: 6px;
	margin: 0 auto;
	padding: 0;
	border: 0;
	background: none;
	color: var(--mt-muted);
	font: inherit;
	font-size: var(--mt-fs-sm);
	cursor: pointer;
}

.login-back:hover {
	color: var(--mt-text);
}

@media (max-width: 480px) {
	#login-card {
		width: 100%;
	}
}
</style>
