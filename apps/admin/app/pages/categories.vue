<script setup lang="ts">
import { computed, h, nextTick, onMounted, reactive, ref } from 'vue'
import type { FormInstance } from 'ant-design-vue'
import {
	DeleteOutlined,
	EditOutlined,
	FolderOpenOutlined,
	FolderOutlined,
	PlusOutlined,
	SearchOutlined,
} from '@ant-design/icons-vue'
import { errorMessage } from '~/api/client'
import { metasApi } from '~/api'
import type { Category } from '~/types/domain'

definePageMeta({ layout: 'admin' })

interface CategoryNode extends Category {
	children?: CategoryNode[]
}

interface TreeItem {
	key: number
	title: string
	item: CategoryNode
	children?: TreeItem[]
}

const { message } = useNotify()

const categories = ref<CategoryNode[]>([])
const expandedKeys = ref<number[]>([])
const modalOpen = ref(false)
const saving = ref(false)
const editingCategory = ref<CategoryNode | null>(null)
const formRef = ref<FormInstance>()

const formState = reactive({
	name: '',
	slug: '',
	parent: undefined as number | undefined,
	description: '',
})

const { searchText, filtered } = useSearchFilter<CategoryNode>(categories, [
	'name',
	'slug',
	'description',
])

function buildTree(items: Category[]): CategoryNode[] {
	const map = new Map<number, CategoryNode>()
	const roots: CategoryNode[] = []

	items.forEach((item) => {
		map.set(item.mid, { ...item, children: [] })
	})

	items.forEach((item) => {
		const node = map.get(item.mid)
		if (!node) return
		if (item.parent && map.has(item.parent)) {
			map.get(item.parent)?.children?.push(node)
		} else {
			roots.push(node)
		}
	})

	return roots
}

async function fetchCategories() {
	try {
		const items = (await metasApi.list('category')) || []
		categories.value = buildTree(items)
		expandedKeys.value = categories.value.map((item) => item.mid)
	} catch (error) {
		message.error(errorMessage(error, '获取分类列表失败'))
	}
}

function openCreate() {
	editingCategory.value = null
	Object.assign(formState, { name: '', slug: '', parent: undefined, description: '' })
	modalOpen.value = true
	nextTick(() => formRef.value?.clearValidate())
}

function openEdit(record: CategoryNode) {
	editingCategory.value = record
	Object.assign(formState, {
		name: record.name,
		slug: record.slug,
		parent: record.parent,
		description: record.description || '',
	})
	modalOpen.value = true
	nextTick(() => formRef.value?.clearValidate())
}

function closeModal() {
	modalOpen.value = false
	editingCategory.value = null
}

async function handleDelete(mid: number) {
	try {
		await metasApi.remove(mid)
		message.success('删除成功')
		await fetchCategories()
	} catch (error) {
		message.error(errorMessage(error, '删除失败'))
	}
}

async function handleOk() {
	try {
		await formRef.value?.validate()
	} catch {
		return
	}

	saving.value = true
	try {
		if (editingCategory.value) {
			await metasApi.update(editingCategory.value.mid, { ...formState })
			message.success('更新成功')
		} else {
			await metasApi.create({ ...formState, type: 'category' })
			message.success('创建成功')
		}
		closeModal()
		await fetchCategories()
	} catch (error) {
		message.error(errorMessage(error, editingCategory.value ? '更新失败' : '创建失败'))
	} finally {
		saving.value = false
	}
}

function toTreeData(items: CategoryNode[]): TreeItem[] {
	return items.map((item) => ({
		key: item.mid,
		title: item.name,
		item,
		children: item.children?.length ? toTreeData(item.children) : undefined,
	}))
}

const treeData = computed(() => toTreeData(filtered.value))

onMounted(fetchCategories)
</script>

<template>
	<div>
		<a-card>
			<div class="page-header">
				<a-typography-title :level="4" :style="{ margin: 0 }">
					<FolderOutlined :style="{ marginRight: '8px' }" />
					分类管理
				</a-typography-title>
				<a-button type="primary" :icon="h(PlusOutlined)" @click="openCreate"> 新增分类 </a-button>
			</div>

			<a-input
				v-model:value="searchText"
				placeholder="搜索分类名称"
				allow-clear
				:style="{ width: '300px', maxWidth: '100%', marginBottom: '16px' }"
			>
				<template #prefix><SearchOutlined /></template>
			</a-input>

			<a-tree
				:tree-data="treeData"
				:expanded-keys="expandedKeys"
				:motion="false"
				show-icon
				:style="{ background: 'transparent' }"
				@update:expanded-keys="(keys: number[]) => (expandedKeys = keys)"
			>
				<template #icon="{ expanded }">
					<FolderOpenOutlined v-if="expanded" />
					<FolderOutlined v-else />
				</template>
				<template #title="{ dataRef }">
					<div class="category-node">
						<div class="category-node-main">
							<span class="category-name">{{ dataRef.item.name }}</span>
							<span class="category-slug">/{{ dataRef.item.slug }}</span>
							<span v-if="dataRef.item.count > 0" class="category-count">
								({{ dataRef.item.count }} 篇文章)
							</span>
						</div>
						<div class="category-actions" @click.stop>
							<a-tooltip title="编辑">
								<a-button
									type="text"
									size="small"
									:icon="h(EditOutlined)"
									@click="openEdit(dataRef.item)"
								/>
							</a-tooltip>
							<a-popconfirm
								title="确定要删除这个分类吗？"
								description="删除后无法恢复, 相关文章将失去分类"
								ok-text="确定"
								cancel-text="取消"
								@confirm="handleDelete(dataRef.item.mid)"
							>
								<a-button type="text" size="small" danger :icon="h(DeleteOutlined)" />
							</a-popconfirm>
						</div>
					</div>
				</template>
			</a-tree>
		</a-card>

		<a-modal
			:title="editingCategory ? '编辑分类' : '新增分类'"
			:open="modalOpen"
			:width="500"
			:confirm-loading="saving"
			ok-text="确定"
			cancel-text="取消"
			@ok="handleOk"
			@cancel="closeModal"
		>
			<a-form ref="formRef" :model="formState" layout="vertical">
				<a-form-item
					label="分类名称"
					name="name"
					:rules="[{ required: true, message: '请输入分类名称' }]"
				>
					<a-input v-model:value="formState.name" placeholder="例如：技术文章" />
				</a-form-item>

				<a-form-item
					label="别名 (Slug)"
					name="slug"
					:rules="[
						{ required: true, message: '请输入别名' },
						{ pattern: /^[a-z0-9-]+$/, message: '只能使用小写字母、数字和连字符' },
					]"
				>
					<a-input v-model:value="formState.slug" placeholder="例如：tech-articles" />
				</a-form-item>

				<a-form-item label="父分类" name="parent">
					<a-input-search
						:value="formState.parent"
						placeholder="输入分类ID, 留空为顶级分类"
						allow-clear
						disabled
					/>
				</a-form-item>

				<a-form-item
					label="描述"
					name="description"
					:rules="[{ max: 200, message: '描述不能超过 200 个字符' }]"
				>
					<a-textarea
						v-model:value="formState.description"
						:rows="3"
						placeholder="分类的描述信息"
						:maxlength="200"
						show-count
					/>
				</a-form-item>
			</a-form>
		</a-modal>
	</div>
</template>

<style scoped>
.page-header {
	display: flex;
	justify-content: space-between;
	align-items: center;
	margin-bottom: 16px;
	flex-wrap: wrap;
	gap: 12px;
}

:deep(.category-node) {
	display: flex;
	align-items: center;
	justify-content: space-between;
	padding: 8px 0;
	gap: 16px;
}

:deep(.category-node-main) {
	display: flex;
	align-items: center;
	gap: 8px;
	min-width: 0;
}

:deep(.category-name) {
	font-weight: 600;
	color: #1e293b;
}

:deep(.category-slug),
:deep(.category-count) {
	font-size: 12px;
	color: #94a3b8;
}

:deep(.category-actions) {
	display: flex;
	align-items: center;
	gap: 4px;
}
</style>
