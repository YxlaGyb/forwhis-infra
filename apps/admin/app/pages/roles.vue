<script setup lang="ts">
import type { ActionItem, ResponsiveColumn } from '~/types/domain'
import { getRoleName, getRoleTone } from '~/utils/mapping'

definePageMeta({ layout: 'admin' })

interface Role {
	id: number
	name: string
	description: string
	permissions: string[]
	userCount: number
}

interface PermissionCategory {
	category: string
	permissions: { key: string; name: string }[]
}

const PERMISSION_CATEGORIES: PermissionCategory[] = [
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
	superadmin: PERMISSION_CATEGORIES.flatMap((group) => group.permissions.map((item) => item.key)),
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
	editor: ['post:read', 'post:create', 'post:edit', 'post:publish', 'comment:read', 'comment:create'],
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
const formEl = ref<{ submit: () => void; clearErrors: () => void } | null>(null)
const formState = ref({ name: '', description: '' })

const isSystemRole = computed(() => editingRole.value?.name === 'superadmin')
const permissionValue = computed(() => selectedPermissions.value.join(','))

const RULES = {
	name: '[{"rule":"required","message":"请输入角色名称"},{"rule":"pattern","value":"^[a-z0-9_]+$","message":"只能使用小写字母、数字和下划线"}]',
	description: '[{"rule":"required","message":"请输入角色描述"}]',
}

const columns: ResponsiveColumn[] = [
	{ key: 'name', title: '角色名称', dataIndex: 'name' },
	{ key: 'description', title: '描述', dataIndex: 'description', priority: 'medium' },
	{ key: 'permissions', title: '权限数量', dataIndex: 'permissions', priority: 'low' },
	{ key: 'userCount', title: '用户数量', dataIndex: 'userCount', priority: 'low' },
]

function openCreate() {
	editingRole.value = null
	selectedPermissions.value = []
	formState.value = { name: '', description: '' }
	modalOpen.value = true
	nextTick(() => formEl.value?.clearErrors())
}

function openEdit(role: Role) {
	editingRole.value = role
	selectedPermissions.value = [...role.permissions]
	formState.value = { name: role.name, description: role.description }
	modalOpen.value = true
	nextTick(() => formEl.value?.clearErrors())
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

function handlePermissionChange(event: Event) {
	const values = (event as CustomEvent<{ values?: string[] }>).detail?.values
	selectedPermissions.value = values ?? []
}

function categorySelected(category: PermissionCategory): boolean {
	return category.permissions.every((item) => selectedPermissions.value.includes(item.key))
}

function toggleCategory(category: PermissionCategory) {
	const keys = category.permissions.map((item) => item.key)
	selectedPermissions.value = categorySelected(category)
		? selectedPermissions.value.filter((key) => !keys.includes(key))
		: [...new Set([...selectedPermissions.value, ...keys])]
}

async function handleSubmit(event: Event) {
	const values = (event as CustomEvent<{ values?: Record<string, string> }>).detail?.values ?? {}

	if (editingRole.value) {
		roles.value = roles.value.map((role) =>
			role.id === editingRole.value?.id
				? {
						...role,
						name: values.name ?? role.name,
						description: values.description ?? role.description,
						permissions: [...selectedPermissions.value],
					}
				: role,
		)
		message.success('角色更新成功')
	} else {
		roles.value = [
			...roles.value,
			{
				id: Date.now(),
				name: values.name ?? '',
				description: values.description ?? '',
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
			icon: 'pencil',
			disabled: record.name === 'superadmin',
			onClick: () => openEdit(record),
		},
		{
			key: 'delete',
			text: '删除',
			icon: 'trash-2',
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
	<div class="page">
		<UserRoleTabs />

		<header class="page-head">
			<div>
				<h1 class="page-title">角色权限</h1>
				<p class="page-desc">配置角色的描述与可用权限</p>
			</div>
			<mao-button color="primary" @click="openCreate">
				<mao-icon name="plus" :size="15" />
				添加角色
			</mao-button>
		</header>

		<mao-card padding="md">
			<div class="mt-card-body">
				<ResponsiveTable
					:columns="columns"
					:data-source="roles"
					row-key="id"
					:loading="loading"
					:actions="getActions"
					:pagination="false"
				>
					<template #cell-name="{ record }">
						<span class="row">
							<mao-avatar :label="getRoleName((record as Role).name)" size="sm" />
							<mao-badge :tone="getRoleTone((record as Role).name)">
								{{ getRoleName((record as Role).name) }}
							</mao-badge>
						</span>
					</template>
					<template #cell-permissions="{ record }">
						{{ (record as Role).permissions.length }} 个权限
					</template>
				</ResponsiveTable>
			</div>
		</mao-card>

		<mao-dialog
			:open="modalOpen || undefined"
			:heading="editingRole ? '编辑角色' : '添加角色'"
			size="lg"
			closable
			@mt-close="closeModal"
		>
			<mao-form ref="formEl" label="角色信息" @mt-submit="handleSubmit">
				<mao-form-item name="name" label="角色名称" :rules="RULES.name">
					<mao-text-field
						name="name"
						:value="formState.name"
						:readonly="!!editingRole || undefined"
						:disabled="isSystemRole || undefined"
						placeholder="请输入角色名称, 如：editor"
					/>
				</mao-form-item>

				<mao-form-item name="description" label="描述" :rules="RULES.description">
					<mao-text-area
						name="description"
						:value="formState.description"
						:rows="2"
						:disabled="isSystemRole || undefined"
						placeholder="请输入角色描述"
					/>
				</mao-form-item>
			</mao-form>

			<div class="stack permissions">
				<div class="row-between">
					<h4 class="section-title">权限配置</h4>
					<mao-badge tone="neutral">{{ selectedPermissions.length }} 项已选</mao-badge>
				</div>

				<mao-notice v-if="isSystemRole" tone="info">
					超级管理员拥有所有权限, 不可修改
				</mao-notice>

				<mao-card v-for="group in PERMISSION_CATEGORIES" :key="group.category" padding="sm" bordered>
					<div class="mt-card-body">
						<div class="row-between group-head">
							<span class="strong">{{ group.category }}</span>
							<mao-button
								size="sm"
								variant="flat"
								:disabled="isSystemRole || undefined"
								@click="toggleCategory(group)"
							>
								{{ categorySelected(group) ? '清空' : '全选' }}
							</mao-button>
						</div>

						<mao-checkbox-group
							:value="permissionValue"
							:disabled="isSystemRole || undefined"
							:label="group.category"
							@mt-change="handlePermissionChange"
						>
							<mao-checkbox v-for="item in group.permissions" :key="item.key" :value="item.key">
								{{ item.name }}
							</mao-checkbox>
						</mao-checkbox-group>
					</div>
				</mao-card>
			</div>

			<mao-button slot="actions" variant="flat" @click="closeModal">取消</mao-button>
			<mao-button slot="actions" color="primary" @click="formEl?.submit()">确定</mao-button>
		</mao-dialog>
	</div>
</template>

<style scoped>
.permissions {
	margin-top: var(--mt-space-4);
}

.group-head {
	margin-bottom: var(--mt-space-2);
}

mao-checkbox-group {
	display: grid;
	grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
	gap: var(--mt-space-2) var(--mt-space-3);
}
</style>
