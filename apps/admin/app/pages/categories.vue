<script setup lang="ts">
import { errorMessage } from '~/api/client'
import { metasApi } from '~/api'
import type { Category } from '~/types/domain'

definePageMeta({ layout: 'admin' })

interface CategoryNode extends Category {
	children?: CategoryNode[]
}

interface CategoryRow extends CategoryNode {
	depth: number
}

interface FormHost extends HTMLElement {
	submit: () => void
	clearErrors: () => void
}

const { message, modal } = useNotify()

const categories = ref<CategoryNode[]>([])
const expandedKeys = ref<number[]>([])
const dialogOpen = ref(false)
const saving = ref(false)
const editingCategory = ref<CategoryNode | null>(null)
const form = ref({ name: '', slug: '', parent: '', description: '' })
const formEl = ref<FormHost | null>(null)

const RULES = {
	name: '[{"rule":"required","message":"请输入分类名称"}]',
	slug: '[{"rule":"required","message":"请输入别名"},{"rule":"pattern","value":"^[a-z0-9-]+$","message":"只能使用小写字母、数字和连字符"}]',
	description: '[{"rule":"maxLength","value":200,"message":"描述不能超过 200 个字符"}]',
}

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

function flatten(nodes: CategoryNode[], depth = 0, out: CategoryRow[] = []): CategoryRow[] {
	nodes.forEach((node) => {
		out.push({ ...node, depth })
		if (node.children?.length) flatten(node.children, depth + 1, out)
	})
	return out
}

const flatRows = computed(() => flatten(categories.value))
const { searchText, filtered } = useSearchFilter<CategoryRow>(flatRows, ['name', 'slug', 'description'])

const rows = computed(() => {
	if (searchText.value.trim()) return filtered.value
	return flatRows.value.filter((row) => row.depth === 0 || expandedKeys.value.includes(row.parent ?? -1))
})

function hasChildren(row: CategoryRow): boolean {
	return Boolean(row.children?.length)
}

function isExpanded(row: CategoryRow): boolean {
	return expandedKeys.value.includes(row.mid)
}

function toggleExpand(row: CategoryRow) {
	expandedKeys.value = isExpanded(row)
		? expandedKeys.value.filter((key) => key !== row.mid)
		: [...expandedKeys.value, row.mid]
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
	form.value = { name: '', slug: '', parent: '', description: '' }
	dialogOpen.value = true
	nextTick(() => formEl.value?.clearErrors())
}

function openEdit(row: CategoryRow) {
	editingCategory.value = row
	form.value = {
		name: row.name,
		slug: row.slug,
		parent: row.parent ? String(row.parent) : '',
		description: row.description || '',
	}
	dialogOpen.value = true
	nextTick(() => formEl.value?.clearErrors())
}

function closeDialog() {
	dialogOpen.value = false
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

function requestDelete(row: CategoryRow) {
	void modal.confirm({
		title: '确定要删除这个分类吗？',
		content: '删除后无法恢复, 相关文章将失去分类',
		okText: '确定',
		cancelText: '取消',
		okType: 'danger',
		onOk: () => handleDelete(row.mid),
	})
}

async function handleSubmit(event: Event) {
	const values = (event as CustomEvent<{ values?: Record<string, string> }>).detail?.values ?? {}
	const parent = values.parent?.trim() ? Number(values.parent) : undefined
	if (parent !== undefined && !Number.isFinite(parent)) {
		message.warning('父分类 ID 必须是数字')
		return
	}

	saving.value = true
	try {
		const payload = {
			name: values.name ?? '',
			slug: values.slug ?? '',
			parent,
			description: values.description ?? '',
		}
		if (editingCategory.value) {
			await metasApi.update(editingCategory.value.mid, payload)
			message.success('更新成功')
		} else {
			await metasApi.create({ ...payload, type: 'category' })
			message.success('创建成功')
		}
		closeDialog()
		await fetchCategories()
	} catch (error) {
		message.error(errorMessage(error, editingCategory.value ? '更新失败' : '创建失败'))
	} finally {
		saving.value = false
	}
}

onMounted(fetchCategories)
</script>

<template>
	<div class="page">
		<header class="page-head">
			<div>
				<h1 class="page-title">分类管理</h1>
				<p class="page-desc">维护内容分类的层级、别名与描述</p>
			</div>
			<mao-button color="primary" @click="openCreate">
				<mao-icon name="plus" :size="15" />
				新增分类
			</mao-button>
		</header>

		<mao-card padding="md">
			<div class="mt-card-body">
				<mao-text-field
					:value="searchText"
					icon="search"
					clearable
					placeholder="搜索分类名称"
					aria-label="搜索分类"
					class="search"
					@mt-input="searchText = ($event as CustomEvent<{ value: string }>).detail.value"
				/>

				<mao-empty v-if="!rows.length" icon="folder-open" title="暂无分类" description="创建第一个分类后就会出现在这里" />

				<ul v-else class="cat-list">
					<li v-for="row in rows" :key="row.mid" class="cat-item" :class="{ 'is-child': row.depth > 0 }">
						<button
							v-if="hasChildren(row)"
							type="button"
							class="cat-toggle"
							:aria-expanded="isExpanded(row)"
							:aria-label="isExpanded(row) ? '收起子分类' : '展开子分类'"
							@click="toggleExpand(row)"
						>
							<mao-icon :name="isExpanded(row) ? 'chevron-down' : 'chevron-right'" :size="15" />
						</button>
						<span v-else class="cat-toggle is-placeholder" aria-hidden="true"/>

						<mao-icon :name="row.depth > 0 ? 'folder' : 'folder-open'" :size="16" class="cat-icon" />

						<span class="cat-name">{{ row.name }}</span>
						<span class="cat-slug">/{{ row.slug }}</span>
						<span v-if="(row.count ?? 0) > 0" class="cat-count">{{ row.count ?? 0 }} 篇</span>

						<span class="cat-actions">
							<mao-button size="sm" variant="flat" aria-label="编辑分类" @click="openEdit(row)">
								<mao-icon name="pencil" :size="15" />
							</mao-button>
							<mao-button size="sm" variant="flat" color="danger" aria-label="删除分类" @click="requestDelete(row)">
								<mao-icon name="trash-2" :size="15" />
							</mao-button>
						</span>
					</li>
				</ul>
			</div>
		</mao-card>

		<mao-dialog
			:open="dialogOpen || undefined"
			:heading="editingCategory ? '编辑分类' : '新增分类'"
			size="sm"
			closable
			@mt-close="closeDialog"
		>
			<mao-form ref="formEl" label="分类信息" @mt-submit="handleSubmit">
				<mao-form-item name="name" label="分类名称" :rules="RULES.name">
					<mao-text-field name="name" :value="form.name" placeholder="例如：技术文章" />
				</mao-form-item>

				<mao-form-item name="slug" label="别名 (Slug)" :rules="RULES.slug">
					<mao-text-field name="slug" :value="form.slug" placeholder="例如：tech-articles" />
				</mao-form-item>

				<mao-form-item name="parent" label="父分类 ID" hint="留空表示顶级分类">
					<mao-text-field name="parent" :value="form.parent" placeholder="输入分类 ID" />
				</mao-form-item>

				<mao-form-item name="description" label="描述" :rules="RULES.description">
					<mao-text-area
						name="description"
						:value="form.description"
						:rows="3"
						:maxlength="200"
						show-count
						placeholder="分类的描述信息"
					/>
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
.search {
	width: 300px;
	max-width: 100%;
	margin-bottom: var(--mt-space-3);
}

.cat-list {
	margin: 0;
	padding: 0;
	list-style: none;
	display: flex;
	flex-direction: column;
	gap: 2px;
}

.cat-item {
	display: flex;
	align-items: center;
	gap: var(--mt-space-2);
	padding: 8px 10px;
	border-radius: var(--mt-radius);
	color: var(--mt-text);
}

.cat-item:hover {
	background: var(--mt-surface-2);
}

.cat-item.is-child {
	margin-left: 22px;
}

.cat-toggle {
	display: inline-flex;
	align-items: center;
	justify-content: center;
	width: 20px;
	height: 20px;
	flex: none;
	padding: 0;
	border: 0;
	border-radius: var(--mt-radius-sm);
	background: transparent;
	color: var(--mt-muted);
	cursor: pointer;
}

.cat-toggle:hover {
	background: var(--mt-border);
	color: var(--mt-text);
}

.cat-toggle.is-placeholder {
	cursor: default;
	background: transparent;
}

.cat-icon {
	color: var(--mt-muted);
}

.cat-name {
	font-weight: 500;
}

.cat-slug,
.cat-count {
	font-size: var(--mt-fs-xs);
	color: var(--mt-faint);
}

.cat-actions {
	display: flex;
	align-items: center;
	gap: 2px;
	margin-left: auto;
}

@media (max-width: 480px) {
	.cat-slug {
		display: none;
	}
}
</style>
