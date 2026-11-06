<script setup lang="ts">
import { h, nextTick, onMounted, reactive, ref } from 'vue'
import type { FormInstance } from 'ant-design-vue'
import {
	DeleteOutlined,
	EditOutlined,
	PlusOutlined,
	SearchOutlined,
	UserOutlined,
} from '@ant-design/icons-vue'
import { errorMessage } from '~/api/client'
import { usersApi } from '~/api'
import type { ActionItem, AdminUser, ResponsiveColumn } from '~/types/domain'
import { ROLE_OPTIONS, USER_LIMITS } from '~/utils/constants'
import { formatTimestamp } from '~/utils/format'
import { getRoleColor, getRoleName } from '~/utils/mapping'

definePageMeta({ layout: 'admin' })

const { message, modal } = useNotify()

const users = ref<AdminUser[]>([])
const loading = ref(false)
const modalOpen = ref(false)
const saving = ref(false)
const editingUser = ref<AdminUser | null>(null)
const formRef = ref<FormInstance>()

const formState = reactive({
	name: '',
	email: '',
	password: '',
	roles: ['user'] as string[],
})

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
	Object.assign(formState, { name: '', email: '', password: '', roles: ['user'] })
	modalOpen.value = true
	nextTick(() => formRef.value?.clearValidate())
}

function openEdit(user: AdminUser) {
	editingUser.value = user
	Object.assign(formState, {
		name: user.name,
		email: user.email,
		password: '',
		roles: user.roles?.length ? [...user.roles] : ['user'],
	})
	modalOpen.value = true
	nextTick(() => formRef.value?.clearValidate())
}

function closeModal() {
	modalOpen.value = false
	editingUser.value = null
}

async function handleOk() {
	try {
		await formRef.value?.validate()
	} catch {
		return
	}

	saving.value = true
	try {
		if (editingUser.value) {
			await usersApi.update(editingUser.value.uid, {
				name: formState.name,
				email: formState.email,
				roles: formState.roles,
			})
			message.success('用户更新成功')
		} else {
			await usersApi.create({
				name: formState.name,
				email: formState.email,
				password: formState.password,
				roles: formState.roles,
			})
			message.success('用户添加成功')
		}
		closeModal()
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
		{
			key: 'edit',
			text: '编辑',
			icon: EditOutlined,
			onClick: () => openEdit(record),
		},
		{
			key: 'delete',
			text: '删除',
			icon: DeleteOutlined,
			danger: true,
			onClick: () => {
				modal.confirm({
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
	<div>
		<UserRoleTabs />

		<a-card>
			<div class="toolbar">
				<a-input
					v-model:value="searchText"
					placeholder="搜索用户昵称或邮箱"
					allow-clear
					class="toolbar-search"
				>
					<template #prefix><SearchOutlined /></template>
				</a-input>
				<a-button type="primary" :icon="h(PlusOutlined)" @click="openCreate"> 添加用户 </a-button>
			</div>

			<ResponsiveTable
				:columns="columns"
				:data-source="filtered"
				row-key="uid"
				:loading="loading"
				:actions="getActions"
				:pagination="{
					showSizeChanger: true,
					showQuickJumper: true,
					showTotal: (total: number) => `共 ${total} 条记录`,
				}"
			>
				<template #card-title="{ record }">
					<a-space>
						<a-avatar>
							<template #icon><UserOutlined /></template>
						</a-avatar>
						{{ record.name }}
					</a-space>
				</template>
				<template #card-subtitle="{ record }">{{ record.email }}</template>
				<template #cell-roles="{ record }">
					<a-space wrap>
						<a-tag v-for="role in record.roles" :key="role" :color="getRoleColor(role)">
							{{ getRoleName(role) }}
						</a-tag>
					</a-space>
				</template>
				<template #cell-createdAt="{ record }">
					{{ formatTimestamp(record.createdAt) }}
				</template>
				<template #cell-lastLogin="{ record }">
					{{ record.lastLogin ? formatTimestamp(record.lastLogin) : '-' }}
				</template>
			</ResponsiveTable>
		</a-card>

		<a-modal
			:title="editingUser ? '编辑用户' : '添加用户'"
			:open="modalOpen"
			:width="600"
			:confirm-loading="saving"
			ok-text="确定"
			cancel-text="取消"
			@ok="handleOk"
			@cancel="closeModal"
		>
			<a-form ref="formRef" :model="formState" layout="vertical" name="userForm">
				<a-form-item
					label="昵称"
					name="name"
					:rules="[{ required: true, message: '请输入用户昵称' }]"
				>
					<a-input v-model:value="formState.name" placeholder="请输入用户昵称" />
				</a-form-item>

				<a-form-item
					label="邮箱"
					name="email"
					:rules="[
						{ required: true, message: '请输入邮箱地址' },
						{ type: 'email', message: '请输入有效的邮箱地址' },
					]"
				>
					<a-input v-model:value="formState.email" placeholder="请输入邮箱地址" />
				</a-form-item>

				<a-form-item
					v-if="!editingUser"
					label="密码"
					name="password"
					:rules="[
						{ required: true, message: '请输入密码' },
						{ min: USER_LIMITS.passwordMinLength, message: '密码至少8个字符' },
						{ max: USER_LIMITS.passwordMaxLength, message: '密码最多16个字符' },
						{ pattern: USER_LIMITS.passwordPattern, message: '密码需包含字母和数字' },
					]"
				>
					<a-input-password
						v-model:value="formState.password"
						placeholder="请输入密码（8-16位，包含字母和数字）"
					/>
				</a-form-item>

				<a-form-item
					label="角色"
					name="roles"
					:rules="[{ required: true, type: 'array', message: '请选择用户角色' }]"
				>
					<a-select
						v-model:value="formState.roles"
						mode="multiple"
						placeholder="请选择用户角色"
						allow-clear
						:options="ROLE_OPTIONS"
					/>
				</a-form-item>
			</a-form>
		</a-modal>
	</div>
</template>

<style scoped>
.toolbar {
	margin-bottom: 16px;
	display: flex;
	justify-content: space-between;
	align-items: center;
	gap: 12px;
}

.toolbar-search {
	width: 300px;
	max-width: 100%;
}

@media (max-width: 480px) {
	.toolbar {
		flex-direction: column;
		align-items: stretch;
	}

	.toolbar-search {
		width: 100%;
	}
}
</style>
