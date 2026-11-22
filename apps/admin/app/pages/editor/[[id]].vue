<script setup lang="ts">
import { computed, nextTick, onMounted, reactive, ref, watch } from 'vue'
import { useEventListener, useMediaQuery } from '@vueuse/core'
import { errorMessage } from '~/api/client'
import { contentsApi, metasApi } from '~/api'
import type { Category, ContentType } from '~/types/domain'
import { CONTENT_LIMITS, EDITOR_TITLE_MAX_LENGTH } from '~/utils/constants'

definePageMeta({ layout: 'admin' })

interface ToolbarItem {
	icon: string
	title: string
	action: () => void
}

const route = useRoute()
const config = useRuntimeConfig()
const { message, modal } = useNotify()
const isMobile = useMediaQuery('(max-width: 768px)')

const routeId = computed(() => {
	const value = route.params.id
	const raw = Array.isArray(value) ? value[0] : value
	const parsed = Number(raw)
	return Number.isFinite(parsed) && parsed > 0 ? parsed : null
})

const isEditMode = computed(() => routeId.value !== null)
const currentCid = ref<number | null>(routeId.value)

const contentField = ref<HTMLElement | null>(null)
const saving = ref(false)
const publishing = ref(false)
const lastSaved = ref<Date | null>(null)
const isDirty = ref(false)
const settingsOpen = ref(false)
const categories = ref<Category[]>([])

const formData = reactive({
	title: '',
	text: '',
	status: 'draft' as 'draft' | 'publish' | 'private',
	type: 'post' as ContentType,
	slug: '',
	categoryMid: 0,
})

const isPost = computed(() => formData.type === 'post')
const contentLimit = computed(() => CONTENT_LIMITS[formData.type])
const typeLabel = computed(() =>
	formData.type === 'post' ? '帖子' : formData.type === 'article' ? '文章' : '页面',
)
const contentCount = computed(() => formData.text.replace(/\s/g, '').length)

const toolbarItems: ToolbarItem[] = [
	{ icon: 'image', title: '图片', action: () => insertMarkdown('![图片](', ')') },
	{ icon: 'smile', title: '表情', action: () => message.info('表情功能开发中') },
	{ icon: 'hash', title: '话题', action: () => insertMarkdown('#话题#') },
	{ icon: 'bold', title: '粗体', action: () => insertMarkdown('**', '**') },
	{ icon: 'italic', title: '斜体', action: () => insertMarkdown('*', '*') },
	{ icon: 'list', title: '列表', action: () => insertMarkdown('\n- ') },
	{ icon: 'list-ordered', title: '有序列表', action: () => insertMarkdown('\n1. ') },
	{ icon: 'code', title: '代码', action: () => insertMarkdown('```\n', '\n```') },
	{ icon: 'link', title: '链接', action: () => insertMarkdown('[', '](url)') },
]

function getTextarea(): HTMLTextAreaElement | null {
	return contentField.value?.shadowRoot?.querySelector('textarea') ?? null
}

function insertMarkdown(before: string, after = '') {
	const textarea = getTextarea()
	if (!textarea) return

	const start = textarea.selectionStart
	const end = textarea.selectionEnd
	const selected = formData.text.substring(start, end)
	formData.text =
		formData.text.substring(0, start) + before + selected + after + formData.text.substring(end)
	isDirty.value = true

	nextTick(() => {
		textarea.focus()
		const position = start + before.length + selected.length
		textarea.setSelectionRange(position, position)
	})
}

function readValue(event: Event): string {
	return (event as CustomEvent<{ value?: string }>).detail?.value ?? ''
}

function handleTitleInput(event: Event) {
	formData.title = readValue(event)
	isDirty.value = true
}

function handleContentInput(event: Event) {
	const value = readValue(event)
	if (value.replace(/\s/g, '').length > contentLimit.value) {
		const textarea = getTextarea()
		if (textarea) textarea.value = formData.text
		return
	}
	formData.text = value
	isDirty.value = true
}

function selectType(type: ContentType) {
	formData.type = type
	if (type === 'post') formData.title = ''
	isDirty.value = true
}

function selectStatus(status: 'draft' | 'publish' | 'private') {
	formData.status = status
	isDirty.value = true
}

function selectCategory(mid: number) {
	formData.categoryMid = mid
	isDirty.value = true
}

async function fetchCategories() {
	try {
		categories.value = (await metasApi.list('category')) || []
	} catch {
		categories.value = []
	}
}

async function fetchArticle(cid: number) {
	try {
		const data = await contentsApi.get(cid)
		formData.title = data.title || ''
		formData.text = data.text || ''
		formData.status = (data.status as 'draft' | 'publish' | 'private') || 'draft'
		formData.type = (data.type as ContentType) || 'post'
		formData.slug = data.slug || ''
		formData.categoryMid = data.categoryMid || data.category?.mid || 0
		isDirty.value = false
	} catch (error) {
		message.error(errorMessage(error, '获取内容失败'))
		await navigateTo('/articles')
	}
}

watch(routeId, (next) => {
	if (next && next !== currentCid.value) {
		currentCid.value = next
		void fetchArticle(next)
	}
})

function validateForm(): boolean {
	if (!formData.text.trim()) {
		message.error('请输入内容')
		return false
	}
	if (contentCount.value > contentLimit.value) {
		message.error(`${typeLabel.value}内容不能超过${contentLimit.value}字`)
		return false
	}
	if (!isPost.value && !formData.title.trim()) {
		message.error(`请输入${typeLabel.value}标题`)
		return false
	}
	return true
}

async function handleSave(publish = false) {
	if (saving.value || publishing.value) return
	if (!validateForm()) return

	const flag = publish ? publishing : saving
	flag.value = true

	try {
		const payload = {
			title: formData.title,
			text: formData.text,
			status: publish ? 'publish' : formData.status,
			type: formData.type,
			categoryMid: formData.categoryMid,
		}

		if (currentCid.value) {
			await contentsApi.update(currentCid.value, payload)
		} else {
			const created = await contentsApi.create(payload)
			if (created?.cid) {
				currentCid.value = created.cid
				await navigateTo(`/editor/${created.cid}`, { replace: true })
			}
		}

		message.success(publish ? '发布成功' : '保存成功')
		lastSaved.value = new Date()
		isDirty.value = false
	} catch (error) {
		message.error(errorMessage(error, '保存失败'))
	} finally {
		flag.value = false
	}
}

async function handleDelete() {
	if (!currentCid.value) return
	const confirmed = await modal.confirm({
		title: `确定删除这个${typeLabel.value}吗？`,
		content: '删除后无法恢复',
		okText: '确定',
		cancelText: '取消',
		okType: 'danger',
	})
	if (!confirmed) return
	try {
		await contentsApi.remove(currentCid.value)
		message.success('删除成功')
		await navigateTo('/articles')
	} catch (error) {
		message.error(errorMessage(error, '删除失败'))
	}
}

function handlePreview() {
	if (!currentCid.value) {
		message.info('请先保存内容')
		return
	}
	window.open(`${config.public.siteUrl}/post/${formData.slug || currentCid.value}`, '_blank')
}

function formatTime(date: Date) {
	return date.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' })
}

useEventListener(window, 'beforeunload', (event: BeforeUnloadEvent) => {
	if (isDirty.value) {
		event.preventDefault()
		event.returnValue = ''
	}
})

useEventListener(window, 'keydown', (event: KeyboardEvent) => {
	if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === 's') {
		event.preventDefault()
		void handleSave(false)
	}
})

onBeforeRouteLeave(() => {
	if (isDirty.value && !saving.value && !publishing.value) {
		return window.confirm('有未保存的更改，确定离开吗？')
	}
	return true
})

onMounted(async () => {
	await fetchCategories()
	if (currentCid.value) await fetchArticle(currentCid.value)
})
</script>

<template>
	<div class="editor">
		<header class="editor-bar">
			<div class="bar-group">
				<mao-button variant="flat" size="sm" aria-label="返回列表" @click="navigateTo('/articles')">
					<mao-icon name="x" :size="18" />
				</mao-button>
				<span class="editor-heading">
					{{ isEditMode ? `编辑${typeLabel}` : `新建${typeLabel}` }}
				</span>
				<mao-badge v-if="isDirty" tone="accent">未保存</mao-badge>
			</div>

			<div class="bar-group bar-actions">
				<span v-if="lastSaved" class="meta hide-sm">上次保存 {{ formatTime(lastSaved) }}</span>
				<span class="meta nowrap">{{ contentCount }}/{{ contentLimit }} 字</span>
				<mao-button
					v-if="isMobile"
					variant="flat"
					size="sm"
					aria-label="内容设置"
					@click="settingsOpen = true"
				>
					<mao-icon name="settings-2" :size="18" />
				</mao-button>
				<mao-button
					variant="bordered"
					size="sm"
					:disabled="(!isEditMode || saving) || undefined"
					@click="handlePreview"
				>
					<mao-icon name="eye" :size="15" />
					<span class="hide-sm">预览</span>
				</mao-button>
				<mao-button
					variant="bordered"
					size="sm"
					:loading="saving || undefined"
					@click="handleSave(false)"
				>
					<mao-icon name="save" :size="15" />
					<span class="hide-sm">保存草稿</span>
				</mao-button>
				<mao-button
					color="primary"
					size="sm"
					:loading="publishing || undefined"
					@click="handleSave(true)"
				>
					<mao-icon name="send" :size="15" />
					{{ formData.status === 'publish' ? '更新发布' : '立即发布' }}
				</mao-button>
				<mao-button
					v-if="isEditMode"
					color="danger"
					variant="bordered"
					size="sm"
					@click="handleDelete"
				>
					<mao-icon name="trash-2" :size="15" />
					<span class="hide-sm">删除</span>
				</mao-button>
			</div>
		</header>

		<div class="editor-toolbar" role="toolbar" aria-label="Markdown 格式">
			<mao-tooltip v-for="item in toolbarItems" :key="item.title" :content="item.title">
				<mao-button variant="flat" size="sm" :aria-label="item.title" @click="item.action">
					<mao-icon :name="item.icon" :size="16" />
				</mao-button>
			</mao-tooltip>
		</div>

		<div class="editor-body">
			<div class="editor-main">
				<div class="editor-sheet">
					<mao-text-field
						v-if="!isPost"
						class="editor-title"
						:value="formData.title"
						placeholder="请输入标题..."
						:maxlength="EDITOR_TITLE_MAX_LENGTH"
						show-count
						@mt-input="handleTitleInput"
					/>
					<mao-text-area
						ref="contentField"
						class="editor-content"
						:value="formData.text"
						:placeholder="isPost ? '分享你的想法...' : '开始写作...'"
						rows="16"
						@mt-input="handleContentInput"
					/>
				</div>
				<div class="editor-foot">
					<span class="meta">支持 Markdown 语法</span>
					<span class="meta">Ctrl+S 保存 · {{ contentCount }}/{{ contentLimit }} 字</span>
				</div>
			</div>

			<aside v-if="!isMobile" class="editor-aside">
				<EditorSettings
					:categories="categories"
					:category-mid="formData.categoryMid"
					:type="formData.type"
					:status="formData.status"
					:type-label="typeLabel"
					:can-delete="false"
					@select-category="selectCategory"
					@select-type="selectType"
					@select-status="selectStatus"
				/>
			</aside>
		</div>

		<mao-drawer
			v-if="isMobile"
			:open="settingsOpen || undefined"
			heading="内容设置"
			placement="end"
			size="sm"
			closable
			@mt-close="settingsOpen = false"
		>
			<EditorSettings
				:categories="categories"
				:category-mid="formData.categoryMid"
				:type="formData.type"
				:status="formData.status"
				:type-label="typeLabel"
				:can-delete="!!currentCid"
				@select-category="selectCategory"
				@select-type="selectType"
				@select-status="selectStatus"
				@delete="handleDelete"
			/>
		</mao-drawer>
	</div>
</template>

<style scoped>
.editor {
	display: flex;
	flex-direction: column;
	height: 100vh;
	height: 100dvh;
	min-height: 0;
	background: var(--mt-bg);
	color: var(--mt-text);
}

.editor-bar {
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: var(--mt-space-3);
	padding: 10px 20px;
	background: var(--mt-surface);
	border-bottom: 1px solid var(--mt-border);
	flex: none;
	flex-wrap: wrap;
}

.bar-group {
	display: flex;
	align-items: center;
	gap: var(--mt-space-2);
	min-width: 0;
}

.bar-actions {
	flex-wrap: wrap;
	justify-content: flex-end;
}

.editor-heading {
	font-size: var(--mt-fs-md);
	font-weight: 600;
	white-space: nowrap;
}

.meta {
	font-size: var(--mt-fs-xs);
	color: var(--mt-muted);
}

.editor-toolbar {
	display: flex;
	align-items: center;
	gap: 2px;
	padding: 6px 16px;
	background: var(--mt-surface);
	border-bottom: 1px solid var(--mt-border);
	overflow-x: auto;
	flex: none;
}

.editor-body {
	flex: 1;
	min-height: 0;
	display: flex;
	gap: var(--mt-space-4);
	padding: var(--mt-space-4);
	overflow: auto;
}

.editor-main {
	flex: 1;
	min-width: 0;
	display: flex;
	flex-direction: column;
	gap: var(--mt-space-2);
}

.editor-sheet {
	flex: 1;
	min-height: 0;
	display: flex;
	flex-direction: column;
	gap: var(--mt-space-3);
	padding: var(--mt-space-4) var(--mt-space-5);
	background: var(--mt-surface);
	border: 1px solid var(--mt-border);
	border-radius: var(--mt-radius-lg);
}

.editor-title::part(input) {
	font-size: var(--mt-fs-xl);
	font-weight: 600;
	padding-top: 4px;
	padding-bottom: 4px;
}

.editor-content::part(textarea) {
	min-height: 46vh;
	font-size: var(--mt-fs-md);
	line-height: 1.8;
	resize: vertical;
}

.editor-foot {
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: var(--mt-space-3);
	flex-wrap: wrap;
	padding: 0 2px;
}

.editor-aside {
	width: 288px;
	flex: none;
	display: flex;
	flex-direction: column;
	gap: var(--mt-space-3);
}

@media (max-width: 900px) {
	.editor-body {
		flex-direction: column;
	}

	.editor-aside {
		width: 100%;
	}
}

@media (max-width: 768px) {
	.editor-bar {
		padding: 8px 12px;
		gap: var(--mt-space-2);
	}

	.editor-toolbar {
		padding: 4px 10px;
	}

	.editor-body {
		padding: var(--mt-space-3);
		gap: var(--mt-space-3);
	}

	.editor-sheet {
		padding: var(--mt-space-3);
	}

	.editor-content::part(textarea) {
		min-height: 40vh;
	}

	.meta.hide-sm,
	.hide-sm {
		display: none;
	}
}
</style>
