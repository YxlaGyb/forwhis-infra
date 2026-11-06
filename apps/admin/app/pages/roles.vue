<script setup lang="ts">
import { computed, h, nextTick, reactive, ref } from 'vue'
import type { FormInstance } from 'ant-design-vue'
import { DeleteOutlined, EditOutlined, PlusOutlined, SafetyOutlined } from '@ant-design/icons-vue'
import type { ActionItem, ResponsiveColumn } from '~/types/domain'
import { getRoleColor, getRoleName } from '~/utils/mapping'

definePageMeta({ layout: 'admin' })

interface Role {
	id: number
	name: string
	description: string
	permissions: string[]
	userCount: number
}

const PERMISSION_CATEGORIES = [
	{
		category: '用户管理',
		permissions: [
			{ key: 'user:read', name: '查看用户' },
			{ key: 'user:create', name: '创建用户' },
			{ key: 'user:edit', name: '编辑用户' },
			{ key: 'user:delete', name: '删除用户' },
		],
	},
	{
		category: '内容管理',
		permissions: [
			{ key: 'post:read', name: '查看文章' },
			{ key: 'post:create', name: '创建文章' },
			{ key: 'post:edit', name: '编辑文章' },
			{ key: 'post:delete', name: '删除文章' },
			{ key: 'post:publish', name: '发布文章' },
			{ key: 'comment:read', name: '查看评论' },
			{ key: 'comment:create', name: '创建评论' },
			{ key: 'comment:delete', name: '删除评论' },
			{ key: 'comment:moderate', name: '审核评论' },
		],
	},
	{
		category: '系统管理',
		permissions: [
			{ key: 'system:config', name: '系统配置' },
			{ key: 'system:logs', name: '查看日志' },
			{ key: 'system:backup', name: '数据备份' },
			{ key: 'system:plugin', name: '插件管理' },
			{ key: 'settings:edit', name: '编辑系统设置' },
		],
	},
]

const ROLE_PERMISSIONS: Record<string, string[]> = {
	superadmin: [
		'user:read',
		'user:create',
		'user:edit',
		'user:delete',
		'post:read',
		'post:create',
		'post:edit',
		'post:delete',
		'post:publish',
		'comment:read',
		'comment:create',
		'comment:delete',
		'comment:moderate',
		'system:config',
		'system:logs',
		'system:backup',
		'system:plugin',
		'settings:edit',
	],
	admin: [
		'user:read',
		'user:create',
		'user:edit',
		'user:delete',
		'post:read',
		'post:create',
		'post:edit',
		'post:delete',
		'post:publish',
		'comment:read',
		'comment:create',
		'comment:delete',
		'comment:moderate',
		'system:logs',
		'settings:edit',
	],
	editor: [
		'post:read',
		'post:create',
		'post:edit',
		'post:publish',
		'comment:read',
		'comment:create',
	],
	user: ['post:read', 'comment:read', 'comment:create'],
}

const { message } = useNotify()

const roles = ref<Role[]>([
	{
		id: 1,
		name: 'superadmin',
		description: '超级管理员, 拥有所有权限',
		permissions: ROLE_PERMISSIONS.superadmin ?? [],
		userCount: 1,
	},
	{
		id: 2,
		name: 'admin',
		description: '管理员, 可以管理用户和内容',
		permissions: ROLE_PERMISSIONS.admin ?? [],
		userCount: 0,
	},
	{
		id: 3,
		name: 'editor',
		description: '编辑, 可以管理文章',
		permissions: ROLE_PERMISSIONS.editor ?? [],
		userCount: 0,
	},
	{
		id: 4,
		name: 'user',
		description: '普通用户',
		permissions: ROLE_PERMISSIONS.user ?? [],
		userCount: 1,
	},
])

const loading = ref(false)
const modalOpen = ref(false)
const editingRole = ref<Role | null>(null)
const selectedPermissions = ref<string[]>([])
const formRef = ref<FormInstance>()

const formState = reactive({ name: '', description: '' })

const isSystemRole = computed(() => editingRole.value?.name === 'superadmin')

const columns: ResponsiveColumn[] = [
	{ key: 'name', title: '角色名称', dataIndex: 'name' },
	{ key: 'description', title: '描述', dataIndex: 'description', priority: 'medium' },
	{ key: 'permissions', title: '权限数量', dataIndex: 'permissions', priority: 'low' },
	{ key: 'userCount', title: '用户数量', dataIndex: 'userCount', priority: 'low' },
]

function openCreate() {
	editingRole.value = null
	selectedPermissions.value = []
	Object.assign(formState, { name: '', description: '' })
	modalOpen.value = true
	nextTick(() => formRef.value?.clearValidate())
}

function openEdit(role: Role) {
	editingRole.value = role
	selectedPermissions.value = [...role.permissions]
	Object.assign(formState, { name: role.name, description: role.description })
	modalOpen.value = true
	nextTick(() => formRef.value?.clearValidate())
}

function closeModal() {
	modalOpen.value = false
	editingRole.value = null
	selectedPermissions.value = []
}

function handleDelete(id: number) {
	roles.value = roles.value.filter((role) => role.id !== id)
	message.success('角色删除成功')
}

async function handleOk() {
	try {
		await formRef.value?.validate()
	} catch {
		return
	}

	if (editingRole.value) {
		roles.value = roles.value.map((role) =>
			role.id === editingRole.value?.id
				? { ...role, ...formState, permissions: [...selectedPermissions.value] }
				: role,
		)
		message.success('角色更新成功')
	} else {
		roles.value = [
			...roles.value,
			{
				id: Date.now(),
				name: formState.name,
				description: formState.description,
				permissions: [...selectedPermissions.value],
				userCount: 0,
			},
		]
		message.success('角色创建成功')
	}

	closeModal()
}

function getActions(record: Role): ActionItem[] {
	return [
		{
			key: 'edit',
			text: '编辑',
			icon: EditOutlined,
			disabled: record.name === 'superadmin',
			onClick: () => openEdit(record),
		},
		{
			key: 'delete',
			text: '删除',
			icon: DeleteOutlined,
			danger: true,
			disabled: record.name === 'superadmin' || record.userCount > 0,
			onClick: () => {
				if (record.name !== 'superadmin' && record.userCount === 0) handleDelete(record.id)
			},
		},
	]
}
</script>

<template>
	<div>
		<UserRoleTabs />

		<a-card>
			<div class="toolbar">
				<a-button type="primary" :icon="h(PlusOutlined)" @click="openCreate"> 添加角色 </a-button>
			</div>

			<ResponsiveTable
				:columns="columns"
				:data-source="roles"
				row-key="id"
				:loading="loading"
				:actions="getActions"
				:pagination="false"
			>
				<template #card-title="{ record }">
					<a-space>
						<a-avatar size="small">
							<template #icon><SafetyOutlined /></template>
						</a-avatar>
						<a-tag :color="getRoleColor(record.name)">{{ getRoleName(record.name) }}</a-tag>
					</a-space>
				</template>
				<template #card-subtitle="{ record }">{{ record.description }}</template>
				<template #cell-permissions="{ record }"> {{ record.permissions.length }} 个权限 </template>
			</ResponsiveTable>
		</a-card>

		<a-modal
			:title="editingRole ? '编辑角色' : '添加角色'"
			:open="modalOpen"
			:width="700"
			ok-text="确定"
			cancel-text="取消"
			@ok="handleOk"
			@cancel="closeModal"
		>
			<a-form ref="formRef" :model="formState" layout="vertical">
				<a-form-item
					label="角色名称"
					name="name"
					:rules="[
						{ required: true, message: '请输入角色名称' },
						{ pattern: /^[a-z0-9_]+$/, message: '只能使用小写字母、数字和下划线' },
					]"
				>
					<a-input
						v-model:value="formState.name"
						placeholder="请输入角色名称, 如：editor"
						:disabled="isSystemRole || !!editingRole"
					/>
				</a-form-item>

				<a-form-item
					label="描述"
					name="description"
					:rules="[{ required: true, message: '请输入角色描述' }]"
				>
					<a-textarea
						v-model:value="formState.description"
						:rows="2"
						placeholder="请输入角色描述"
						:disabled="isSystemRole"
					/>
				</a-form-item>

				<a-form-item label="权限配置" required>
					<a-tag v-if="isSystemRole" color="red">超级管理员拥有所有权限, 不可修改</a-tag>
					<a-space v-else direction="vertical" style="width: 100%">
						<a-card
							v-for="category in PERMISSION_CATEGORIES"
							:key="category.category"
							size="small"
							:title="category.category"
							style="margin-bottom: 8px"
						>
							<a-checkbox-group v-model:value="selectedPermissions">
								<a-space wrap>
									<a-checkbox
										v-for="permission in category.permissions"
										:key="permission.key"
										:value="permission.key"
									>
										{{ permission.name }}
									</a-checkbox>
								</a-space>
							</a-checkbox-group>
						</a-card>
					</a-space>
				</a-form-item>
			</a-form>
		</a-modal>
	</div>
</template>

<style scoped>
.toolbar {
	margin-bottom: 16px;
}
</style>
