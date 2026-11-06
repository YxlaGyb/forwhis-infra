<script setup lang="ts">
import { reactive, ref, computed } from 'vue'
import { LockOutlined, UserOutlined } from '@ant-design/icons-vue'
import { errorMessage } from '~/api/client'

definePageMeta({ layout: 'blank' })

const config = useRuntimeConfig()
const { message } = useNotify()
const { login, canAccessAdmin } = useAuth()

const form = reactive({ email: '', password: '' })
const loading = ref(false)

const backgroundUrl = computed(() => config.public.loginBgUrl || '/110546943_p0.webp')

async function handleFinish() {
	loading.value = true
	try {
		const user = await login(form.email, form.password)
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
	<div class="login-container" :style="{ backgroundImage: `url('${backgroundUrl}')` }">
		<a-card class="login-card" :bordered="false">
			<div class="login-title">
				<a-typography-title :level="2" style="margin-bottom: 8px; color: #2996ff">
					林影心语
					<span style="color: red; font-size: 16px">admin</span>
				</a-typography-title>
			</div>

			<a-form :model="form" layout="vertical" autocomplete="off" @finish="handleFinish">
				<a-form-item
					label="邮箱"
					name="email"
					:rules="[{ required: true, message: '请输入邮箱地址' }]"
				>
					<a-input
						v-model:value="form.email"
						size="large"
						placeholder="请输入邮箱地址"
						class="login-input"
					>
						<template #prefix><UserOutlined /></template>
					</a-input>
				</a-form-item>

				<a-form-item
					label="密码"
					name="password"
					:rules="[{ required: true, message: '请输入密码' }]"
				>
					<a-input-password
						v-model:value="form.password"
						size="large"
						placeholder="请输入密码"
						class="login-input"
					>
						<template #prefix><LockOutlined /></template>
					</a-input-password>
				</a-form-item>

				<a-form-item>
					<a-button
						type="primary"
						html-type="submit"
						size="large"
						:loading="loading"
						block
						class="login-button"
					>
						登录
					</a-button>
				</a-form-item>

				<a-typography-text class="login-back" @click="openSite"> 返回首页 </a-typography-text>
			</a-form>
		</a-card>
	</div>
</template>

<style scoped>
.login-container {
	min-height: 100vh;
	display: flex;
	align-items: center;
	justify-content: center;
	background-size: cover;
	background-position: center;
	background-repeat: no-repeat;
	background-color: #667eea;
}

.login-card {
	width: 400px;
	box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
	border-radius: 12px;
	background-color: rgba(255, 255, 255, 0.5);
	backdrop-filter: blur(20px);
	border: 1px solid rgba(255, 255, 255, 0.3);
}

.login-title {
	text-align: center;
	margin-bottom: 32px;
}

.login-input {
	background-color: rgba(255, 255, 255, 0.85);
	border: 1px solid rgba(255, 255, 255, 0.5);
}

.login-button {
	background-color: rgba(135, 133, 162, 0.6);
	border-color: rgba(135, 133, 162, 1);
	backdrop-filter: blur(10px);
	transition: all 0.3s ease;
}

.login-button:hover {
	background-color: rgba(135, 133, 162, 1) !important;
}

.login-back {
	color: #f6f6f6;
	cursor: pointer;
}

@media (max-width: 480px) {
	.login-card {
		width: calc(100% - 32px);
	}
}
</style>
