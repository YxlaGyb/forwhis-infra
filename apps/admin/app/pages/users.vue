<script setup lang="ts">
import { errorMessage } from '~/api/client'
import { usersApi } from '~/api'
import type { ActionItem, AdminUser, ResponsiveColumn } from '~/types/domain'
import { ROLE_OPTIONS, USER_LIMITS } from '~/utils/constants'
import { formatTimestamp } from '~/utils/format'
import { getRoleName, getRoleTone } from '~/utils/mapping'

definePageMeta({ layout: 'admin' })

interface FormHost extends HTMLElement {
	submit: () => void
	clearErrors: () => void
}

const { message, modal } = useNotify()

const users = ref<AdminUser[]>([])
const loading = ref(false)
const dialogOpen = ref(false)
const saving = ref(false)
const editingUser = ref<AdminUser | null>(null)
const form = ref({ name: '', email: '', password: '', roles: 'user' })
const formEl = ref<FormHost | null>(null)

const { searchText, filtered } = useSearchFilter<AdminUser>(users, ['name', 'email'])

const columns: ResponsiveColumn[] = [
	{ key: 'uid', title: 'UID', dataIndex: 'uid', priority: 'low' },
	{ key: 'name', title: '昵称', dataIndex: 'name', priority: 'low' },
	{ key: 'email', title: '邮箱', dataIndex: 'email', priority: 'low' },
	{ key: 'roles', title: '角色', dataIndex: 'roles', priority: 'high' },
	{ key: 'createdAt', title: '创建时间', dataIndex: 'createdAt', priority: 'medium' },
	{ key: 'lastLogin', title: '最后登录', dataIndex: 'lastLogin', priority: 'medium' },
	{ key: 'lastip', title: '登录IP', dataIndex: 'lastip', priority: 'medium' },
]

const RULES = {
	name: `[{"rule":"required","message":"请输入用户昵称"},{"rule":"maxLength","value":${USER_LIMITS.nameMaxLength},"message":"昵称最多 ${USER_LIMITS.nameMaxLength} 个字符"}]`,
	email: '[{"rule":"required","message":"请输入邮箱地址"},{"rule":"email","message":"请输入有效的邮箱地址"}]',
	password: `[{"rule":"required","message":"请输入密码"},{"rule":"minLength","value":${USER_LIMITS.passwordMinLength},"message":"密码至少 ${USER_LIMITS.passwordMinLength} 个字符"},{"rule":"maxLength","value":${USER_LIMITS.passwordMaxLength},"message":"密码最多 ${USER_LIMITS.passwordMaxLength} 个字符"},{"rule":"pattern","value":"${USER_LIMITS.passwordPattern.source}","message":"密码需包含字母和数字"}]`,
	roles: '[{"rule":"required","message":"请选择用户角色"}]',
}

const dialogHeading = computed(() => (editingUser.value ? '编辑用户' : '添加用户'))

async function fetchUsers() {
	loading.value = true
	try {
		users.value = (await usersApi.list()) || []
	} catch (error) {
		message.error(errorMessage(error, '获取用户列表失败'))
	} finally {
		loading.value = false
	}
}

function openCreate() {
	editingUser.value = null
	form.value = { name: '', email: '', password: '', roles: 'user' }
	dialogOpen.value = true
	nextTick(() => formEl.value?.clearErrors())
}

function openEdit(user: AdminUser) {
	editingUser.value = user
	form.value = {
		name: user.name,
		email: user.email,
		password: '',
		roles: (user.roles?.length ? user.roles : ['user']).join(','),
	}
	dialogOpen.value = true
	nextTick(() => formEl.value?.clearErrors())
}

function closeDialog() {
	dialogOpen.value = false
	editingUser.value = null
}

function splitRoles(value: string): string[] {
	return value
		.split(',')
		.map((item) => item.trim())
		.filter(Boolean)
}

async function handleSubmit(event: Event) {
	const values = (event as CustomEvent<{ values?: Record<string, string> }>).detail?.values ?? {}
	const roles = splitRoles(values.roles ?? form.value.roles)
	saving.value = true
	try {
		const payload = {
			name: values.name ?? '',
			email: values.email ?? '',
			roles,
		}
		if (editingUser.value) {
			await usersApi.update(editingUser.value.uid, payload)
			message.success('用户更新成功')
		} else {
			await usersApi.create({ ...payload, password: values.password ?? '' })
			message.success('用户添加成功')
		}
		closeDialog()
		await fetchUsers()
	} catch (error) {
		message.error(errorMessage(error, editingUser.value ? '用户更新失败' : '用户添加失败'))
	} finally {
		saving.value = false
	}
}

async function handleDelete(uid: number) {
	try {
		await usersApi.remove(uid)
		users.value = users.value.filter((user) => user.uid !== uid)
		message.success('用户删除成功')
	} catch (error) {
		message.error(errorMessage(error, '用户删除失败'))
	}
}

function getActions(record: AdminUser): ActionItem[] {
	return [
		{ key: 'edit', text: '编辑', icon: 'pencil', onClick: () => openEdit(record) },
		{
			key: 'delete',
			text: '删除',
			icon: 'trash-2',
			danger: true,
			onClick: () => {
				void modal.confirm({
					title: '确定删除这个用户吗？',
					content: '删除后无法恢复',
					okText: '确定',
					cancelText: '取消',
					okType: 'danger',
					onOk: () => handleDelete(record.uid),
				})
			},
		},
	]
}

onMounted(fetchUsers)
</script>

<template>
	<div class="page">
		<UserRoleTabs />

		<mao-card padding="md">
			<div class="mt-card-body">
				<div class="toolbar">
					<mao-text-field
						:value="searchText"
						icon="search"
						clearable
						placeholder="搜索用户昵称或邮箱"
						aria-label="搜索用户"
						class="toolbar-search"
						@mt-input="searchText = ($event as CustomEvent<{ value: string }>).detail.value"
					/>
					<mao-button color="primary" icon="plus" @click="openCreate">添加用户</mao-button>
				</div>

				<ResponsiveTable
					:columns="columns"
					:data-source="filtered"
					row-key="uid"
					:loading="loading"
					:actions="getActions"
					:pagination="{ showSizeChanger: true, showTotal: (total) => `共 ${total} 条记录` }"
				>
					<template #card-title="{ record }">
						<span class="row">
							<mao-avatar :label="record.name" size="sm" />
							{{ record.name }}
						</span>
					</template>
					<template #card-subtitle="{ record }">{{ record.email }}</template>
					<template #cell-roles="{ record }">
						<span class="tag-row">
							<mao-badge v-for="role in record.roles" :key="role" :tone="getRoleTone(role)">
								{{ getRoleName(role) }}
							</mao-badge>
						</span>
					</template>
					<template #cell-createdAt="{ record }">
						{{ formatTimestamp(record.createdAt) }}
					</template>
					<template #cell-lastLogin="{ record }">
						{{ record.lastLogin ? formatTimestamp(record.lastLogin) : '-' }}
					</template>
				</ResponsiveTable>
			</div>
		</mao-card>

		<mao-dialog
			:open="dialogOpen || undefined"
			:heading="dialogHeading"
			size="lg"
			closable
			@mt-close="closeDialog"
		>
			<mao-form ref="formEl" label="用户信息" @mt-submit="handleSubmit">
				<mao-form-item name="name" label="昵称" :rules="RULES.name">
					<mao-text-field name="name" :value="form.name" placeholder="请输入用户昵称" />
				</mao-form-item>

				<mao-form-item name="email" label="邮箱" :rules="RULES.email">
					<mao-text-field name="email" :value="form.email" placeholder="请输入邮箱地址" />
				</mao-form-item>

				<mao-form-item v-if="!editingUser" name="password" label="密码" :rules="RULES.password">
					<mao-text-field
						name="password"
						type="password"
						:value="form.password"
						placeholder="请输入密码（8-16位，包含字母和数字）"
					/>
				</mao-form-item>

				<mao-form-item name="roles" label="角色" hint="可多选，至少选择一个角色" :rules="RULES.roles">
					<mao-checkbox-group name="roles" :value="form.roles" label="用户角色">
						<mao-checkbox v-for="option in ROLE_OPTIONS" :key="option.value" :value="option.value">
							{{ option.label }}
						</mao-checkbox>
					</mao-checkbox-group>
				</mao-form-item>
			</mao-form>

			<mao-button slot="actions" variant="flat" @click="closeDialog">取消</mao-button>
			<mao-button slot="actions" color="primary" :loading="saving || undefined" @click="formEl?.submit()">
				确定
			</mao-button>
		</mao-dialog>
	</div>
</template>

<style scoped>
.toolbar {
	margin-bottom: var(--mt-space-4);
}

.toolbar-search {
	width: 300px;
	max-width: 100%;
}

.tag-row {
	display: flex;
	flex-wrap: wrap;
	gap: 6px;
}

@media (max-width: 768px) {
	.toolbar {
		flex-direction: column;
		align-items: stretch;
	}

	.toolbar-search {
		width: 100%;
	}
}
</style>
