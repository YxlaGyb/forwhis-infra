<script setup lang="ts">
import { computed, h, nextTick, onMounted, reactive, ref, watch } from 'vue'
import { useEventListener, useMediaQuery } from '@vueuse/core'
import {
	BoldOutlined,
	CloseOutlined,
	CodeOutlined,
	DeleteOutlined,
	EyeOutlined,
	ItalicOutlined,
	LinkOutlined,
	OrderedListOutlined,
	PictureOutlined,
	PlusOutlined,
	SaveOutlined,
	SendOutlined,
	SettingOutlined,
	SmileOutlined,
	TagOutlined,
	UnorderedListOutlined,
} from '@ant-design/icons-vue'
import type { Component } from 'vue'
import { errorMessage } from '~/api/client'
import { contentsApi, metasApi } from '~/api'
import type { Category, ContentType } from '~/types/domain'
import { CONTENT_LIMITS, EDITOR_TITLE_MAX_LENGTH } from '~/utils/constants'

definePageMeta({ layout: 'admin' })

interface ToolbarItem {
	icon: Component
	title: string
	action: () => void
}

const route = useRoute()
const config = useRuntimeConfig()
const { message } = useNotify()
const isMobile = useMediaQuery('(max-width: 768px)')

const routeId = computed(() => {
	const value = route.params.id
	const raw = Array.isArray(value) ? value[0] : value
	const parsed = Number(raw)
	return Number.isFinite(parsed) && parsed > 0 ? parsed : null
})

const isEditMode = computed(() => routeId.value !== null)
const currentCid = ref<number | null>(routeId.value)

const editorRoot = ref<HTMLElement | null>(null)
const saving = ref(false)
const publishing = ref(false)
const lastSaved = ref<Date | null>(null)
const isDirty = ref(false)
const showToolbar = ref(true)
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
const titleCount = computed(() => formData.title.length)

const toolbarItems: ToolbarItem[] = [
	{ icon: PictureOutlined, title: '图片', action: () => insertMarkdown('![图片](', ')') },
	{ icon: SmileOutlined, title: '表情', action: () => message.info('表情功能开发中') },
	{ icon: TagOutlined, title: '话题', action: () => insertMarkdown('#话题#') },
	{ icon: BoldOutlined, title: '粗体', action: () => insertMarkdown('**', '**') },
	{ icon: ItalicOutlined, title: '斜体', action: () => insertMarkdown('*', '*') },
	{ icon: UnorderedListOutlined, title: '列表', action: () => insertMarkdown('\n- ') },
	{ icon: OrderedListOutlined, title: '有序列表', action: () => insertMarkdown('\n1. ') },
	{ icon: CodeOutlined, title: '代码', action: () => insertMarkdown('```\n', '\n```') },
	{ icon: LinkOutlined, title: '链接', action: () => insertMarkdown('[', '](url)') },
]

function getTextarea(): HTMLTextAreaElement | null {
	return editorRoot.value?.querySelector('textarea') ?? null
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

function handleTitleInput(event: Event) {
	const input = event.target as HTMLInputElement
	if (input.value.length <= EDITOR_TITLE_MAX_LENGTH) {
		formData.title = input.value
		isDirty.value = true
	} else {
		input.value = formData.title
	}
}

function handleContentInput(event: Event) {
	const textarea = event.target as HTMLTextAreaElement
	const value = textarea.value
	if (value.replace(/\s/g, '').length <= contentLimit.value) {
		formData.text = value
		isDirty.value = true
	} else {
		textarea.value = formData.text
	}
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
	<div ref="editorRoot">
		<div v-if="isMobile" class="editor-mobile">
			<div class="mobile-top">
				<a-button type="text" :icon="h(CloseOutlined)" @click="navigateTo('/articles')" />
				<div class="mobile-actions">
					<a-button type="text" size="small" :loading="saving" @click="handleSave(false)">
						草稿箱
					</a-button>
					<a-button type="primary" size="small" :loading="publishing" @click="handleSave(true)">
						发布
					</a-button>
				</div>
			</div>

			<div class="mobile-body" :style="{ paddingBottom: showToolbar ? '140px' : '90px' }">
				<template v-if="!isPost">
					<div class="mobile-title-wrap">
						<a-input
							class="mobile-title"
							placeholder="标题（必填）"
							:value="formData.title"
							@input="handleTitleInput"
						/>
						<div class="counter" :class="{ limit: titleCount >= EDITOR_TITLE_MAX_LENGTH }">
							{{ titleCount }}/{{ EDITOR_TITLE_MAX_LENGTH }}
						</div>
					</div>
					<a-divider :style="{ margin: '8px 0' }" />
				</template>

				<div class="mobile-content">
					<a-textarea
						class="mobile-textarea"
						:placeholder="isPost ? '分享你的想法...' : '请尽情发挥吧'"
						:value="formData.text"
						@input="handleContentInput"
					/>
					<div class="counter" :class="{ limit: contentCount >= contentLimit }">
						{{ contentCount }}/{{ contentLimit }}
					</div>
				</div>
			</div>

			<div class="mobile-toolbar">
				<div class="mobile-toolbar-row">
					<div class="mobile-toolbar-icons">
						<div
							v-for="item in toolbarItems.slice(0, 5)"
							:key="item.title"
							class="toolbar-icon"
							@click="item.action"
						>
							<component :is="item.icon" />
						</div>
						<div class="toolbar-icon" @click="showToolbar = !showToolbar">
							<PlusOutlined />
						</div>
					</div>
					<div class="toolbar-icon" @click="settingsOpen = true">
						<SettingOutlined />
					</div>
				</div>

				<div v-if="showToolbar" class="mobile-toolbar-extra">
					<a-tooltip v-for="item in toolbarItems.slice(5)" :key="item.title" :title="item.title">
						<div class="toolbar-icon small" @click="item.action">
							<component :is="item.icon" />
						</div>
					</a-tooltip>
				</div>
			</div>

			<a-drawer
				title="内容设置"
				placement="bottom"
				:open="settingsOpen"
				height="50%"
				@close="settingsOpen = false"
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
			</a-drawer>
		</div>

		<div v-else class="editor-desktop">
			<div class="desktop-top">
				<div class="desktop-top-left">
					<a-button :icon="h(CloseOutlined)" @click="navigateTo('/articles')" />
					<span class="editor-heading">{{
						isEditMode ? `编辑${typeLabel}` : `新建${typeLabel}`
					}}</span>
					<a-tag v-if="isDirty" color="orange">未保存</a-tag>
				</div>

				<div class="desktop-top-right">
					<a-typography-text v-if="lastSaved" type="secondary" class="tiny-text">
						上次保存: {{ formatTime(lastSaved) }}
					</a-typography-text>
					<a-typography-text type="secondary" class="tiny-text">
						{{ contentCount }}/{{ contentLimit }} 字
					</a-typography-text>
					<a-button :icon="h(EyeOutlined)" :disabled="!isEditMode" @click="handlePreview">
						预览
					</a-button>
					<a-button :icon="h(SaveOutlined)" :loading="saving" @click="handleSave(false)">
						保存草稿
					</a-button>
					<a-button
						type="primary"
						:icon="h(SendOutlined)"
						:loading="publishing"
						@click="handleSave(true)"
					>
						{{ formData.status === 'publish' ? '更新发布' : '立即发布' }}
					</a-button>
					<a-popconfirm
						v-if="isEditMode"
						:title="`确定删除这个${typeLabel}吗？`"
						description="删除后无法恢复"
						ok-text="确定"
						cancel-text="取消"
						@confirm="handleDelete"
					>
						<a-button danger :icon="h(DeleteOutlined)">删除</a-button>
					</a-popconfirm>
				</div>
			</div>

			<div class="desktop-toolbar">
				<a-tooltip v-for="item in toolbarItems" :key="item.title" :title="item.title">
					<a-button type="text" size="small" @click="item.action">
						<template #icon><component :is="item.icon" /></template>
					</a-button>
				</a-tooltip>
			</div>

			<div class="desktop-body">
				<div class="editor-card">
					<template v-if="!isPost">
						<div class="editor-title-wrap">
							<a-input
								class="editor-title"
								placeholder="请输入标题..."
								:value="formData.title"
								@input="handleTitleInput"
							/>
						</div>
						<a-divider :style="{ margin: 0 }" />
					</template>

					<div class="editor-body">
						<a-textarea
							class="editor-textarea"
							:placeholder="isPost ? '分享你的想法...' : '开始写作...'"
							:value="formData.text"
							@input="handleContentInput"
						/>
					</div>

					<div class="editor-footer">
						<a-typography-text type="secondary" class="tiny-text"
							>支持 Markdown 语法</a-typography-text
						>
						<a-typography-text type="secondary" class="tiny-text">
							快捷键: Ctrl+S 保存 | {{ contentCount }}/{{ contentLimit }} 字
						</a-typography-text>
					</div>
				</div>

				<div class="settings-panel">
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
				</div>
			</div>
		</div>
	</div>
</template>

<style scoped>
.editor-mobile {
	display: flex;
	flex-direction: column;
	height: 100vh;
	background: #fff;
	overflow: hidden;
}

.mobile-top {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 12px 16px;
	border-bottom: 1px solid #f0f0f0;
	flex-shrink: 0;
}

.mobile-actions {
	display: flex;
	gap: 12px;
}

.mobile-body {
	flex: 1;
	display: flex;
	flex-direction: column;
	padding: 16px;
	overflow: auto;
}

.mobile-title-wrap {
	margin-bottom: 12px;
	flex-shrink: 0;
}

.mobile-title {
	font-size: 20px;
	font-weight: 500;
}

.mobile-content {
	flex: 1;
	display: flex;
	flex-direction: column;
	overflow: hidden;
}

.mobile-textarea {
	flex: 1;
	border: none;
	resize: none;
	font-size: 16px;
	line-height: 1.8;
	padding: 0;
	min-height: 200px;
}

.counter {
	text-align: right;
	font-size: 12px;
	color: #999;
	margin-top: 4px;
}

.counter.limit {
	color: #ff4d4f;
}

.mobile-toolbar {
	position: fixed;
	bottom: 0;
	left: 0;
	right: 0;
	background: #fafafa;
	z-index: 100;
	border-top: 1px solid #f0f0f0;
}

.mobile-toolbar-row {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 8px 16px;
}

.mobile-toolbar-icons {
	display: flex;
	gap: 20px;
}

.toolbar-icon {
	cursor: pointer;
	color: #666;
	font-size: 20px;
}

.toolbar-icon.small {
	font-size: 18px;
}

.mobile-toolbar-extra {
	display: flex;
	gap: 20px;
	padding: 8px 16px;
	border-top: 1px solid #f5f5f5;
	overflow-x: auto;
}

.editor-desktop {
	display: flex;
	flex-direction: column;
	height: 100vh;
	background: #f5f5f5;
}

.desktop-top {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 12px 24px;
	background: #fff;
	border-bottom: 1px solid #e8e8e8;
	flex-shrink: 0;
	gap: 12px;
	flex-wrap: wrap;
}

.desktop-top-left,
.desktop-top-right {
	display: flex;
	align-items: center;
	gap: 12px;
	flex-wrap: wrap;
}

.editor-heading {
	font-weight: 600;
	font-size: 16px;
}

.tiny-text {
	font-size: 12px;
}

.desktop-toolbar {
	display: flex;
	padding: 8px 24px;
	background: #fff;
	border-bottom: 1px solid #e8e8e8;
	gap: 4px;
	flex-shrink: 0;
}

.desktop-body {
	flex: 1;
	display: flex;
	padding: 24px;
	overflow: auto;
	gap: 24px;
}

.editor-card {
	flex: 1;
	max-width: 900px;
	background: #fff;
	border-radius: 8px;
	box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
	display: flex;
	flex-direction: column;
	min-height: 0;
}

.editor-title-wrap {
	padding: 32px 40px 16px;
}

.editor-title {
	font-size: 32px;
	font-weight: 600;
}

.editor-body {
	flex: 1;
	display: flex;
	flex-direction: column;
	min-height: 0;
}

.editor-textarea {
	flex: 1;
	border: none;
	resize: none;
	font-size: 16px;
	line-height: 1.8;
	padding: 24px 40px;
	min-height: 400px;
}

.editor-footer {
	display: flex;
	justify-content: space-between;
	align-items: center;
	padding: 12px 40px;
	border-top: 1px solid #f0f0f0;
	background: #fafafa;
	border-radius: 0 0 8px 8px;
}

.settings-panel {
	width: 280px;
	flex-shrink: 0;
	display: flex;
	flex-direction: column;
	gap: 16px;
}

@media (max-width: 1024px) {
	.desktop-body {
		flex-direction: column;
	}

	.settings-panel {
		width: 100%;
		flex-direction: row;
		flex-wrap: wrap;
	}

	.settings-panel > * {
		flex: 1;
		min-width: 240px;
	}
}
</style>
