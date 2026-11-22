<script setup lang="ts">
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useMediaQuery } from '@vueuse/core'
import { errorMessage } from '~/api/client'
import { settingsApi } from '~/api'
import type {
	AllSettings,
	FeatureSettings,
	SeoSettings,
	SettingsSection,
	SiteSettings,
} from '~/types/domain'

definePageMeta({ layout: 'admin' })

const { message, modal } = useNotify()
const isMobile = useMediaQuery('(max-width: 768px)')

const sections = [
	{
		key: 'site',
		icon: 'globe',
		title: '站点信息',
		description: '网站名称、Logo、描述等基本信息',
	},
	{
		key: 'feature',
		icon: 'sliders-horizontal',
		title: '功能开关',
		description: '注册、评论、验证码等功能控制',
	},
	{
		key: 'seo',
		icon: 'search',
		title: 'SEO 设置',
		description: '关键词、作者、站点地图',
	},
] as const

const defaultSettings: AllSettings = {
	site: {
		siteName: '',
		siteDescription: '',
		siteLogo: '',
		favicon: '',
		icp: '',
	},
	feature: {
		enableRegistration: true,
		enableComments: true,
		enableEmailVerify: false,
		enableCaptcha: true,
	},
	seo: {
		siteKeywords: '',
		siteAuthor: '',
		robotsTxt: '',
		enableSitemap: true,
	},
}

const activeSection = ref<SettingsSection>('site')
const loading = ref(false)
const saving = ref(false)
const hasChanges = ref(false)
const drawerOpen = ref(false)
const siteNameInvalid = ref(false)
const originalValues = ref<AllSettings>(cloneSettings(defaultSettings))
const formModel = ref<AllSettings>(cloneSettings(defaultSettings))
let applying = false

watch(
	formModel,
	() => {
		if (!applying) hasChanges.value = true
	},
	{ deep: true },
)

const currentSection = computed(
	() => sections.find((item) => item.key === activeSection.value) ?? sections[0]!,
)

function cloneSettings(value: AllSettings): AllSettings {
	return JSON.parse(JSON.stringify(value)) as AllSettings
}

function toBoolean(value: unknown, fallback: boolean): boolean {
	if (typeof value === 'boolean') return value
	if (value === '1' || value === 'true') return true
	if (value === '0' || value === 'false') return false
	return fallback
}

function mergeSettings(data?: Partial<AllSettings>): AllSettings {
	const feature = { ...defaultSettings.feature, ...(data?.feature || {}) }
	return {
		site: { ...defaultSettings.site, ...(data?.site || {}) },
		feature: {
			enableRegistration: toBoolean(feature.enableRegistration, true),
			enableComments: toBoolean(feature.enableComments, true),
			enableEmailVerify: toBoolean(feature.enableEmailVerify, false),
			enableCaptcha: toBoolean(feature.enableCaptcha, true),
		},
		seo: {
			...defaultSettings.seo,
			...(data?.seo || {}),
			enableSitemap: toBoolean(data?.seo?.enableSitemap, true),
		},
	}
}

function applySettings(value: AllSettings) {
	const next = cloneSettings(value)
	applying = true
	formModel.value.site = next.site
	formModel.value.feature = next.feature
	formModel.value.seo = next.seo
	nextTick(() => {
		applying = false
	})
}

function detailValue(event: Event): string {
	return (event as CustomEvent<{ value?: string }>).detail?.value ?? ''
}

function setSite(key: keyof SiteSettings, event: Event) {
	formModel.value.site[key] = detailValue(event)
	if (key === 'siteName') siteNameInvalid.value = false
}

function setSeo(key: Exclude<keyof SeoSettings, 'enableSitemap'>, event: Event) {
	formModel.value.seo[key] = detailValue(event)
}

function setFeature(key: keyof FeatureSettings, event: Event) {
	const detail = (event as CustomEvent<{ checked?: boolean }>).detail
	formModel.value.feature[key] = Boolean(detail?.checked)
}

function setSitemap(event: Event) {
	const detail = (event as CustomEvent<{ checked?: boolean }>).detail
	formModel.value.seo.enableSitemap = Boolean(detail?.checked)
}

async function fetchSettings() {
	loading.value = true
	try {
		const data = await settingsApi.get()
		const merged = mergeSettings(data)
		originalValues.value = cloneSettings(merged)
		applySettings(merged)
		hasChanges.value = false
	} catch {
		message.warning('无法连接到服务器，使用本地默认配置')
		originalValues.value = cloneSettings(defaultSettings)
		applySettings(defaultSettings)
		hasChanges.value = false
	} finally {
		loading.value = false
	}
}

function handleSectionChange(event: Event) {
	const key = (event as CustomEvent<{ value?: string }>).detail?.value ?? ''
	if (!key) return
	activeSection.value = key as SettingsSection
	drawerOpen.value = false
}

function handleReset() {
	applySettings(originalValues.value)
	hasChanges.value = false
	message.info('已重置为上次保存的设置')
}

function handleResetToDefault() {
	modal.confirm({
		title: '重置为默认设置',
		content: '所有设置将恢复为系统默认值，保存后生效。确定继续吗？',
		okText: '重置',
		cancelText: '取消',
		onOk: () => {
			applySettings(defaultSettings)
			hasChanges.value = true
			message.info('已重置为默认设置')
		},
	})
}

async function handleSave() {
	if (!formModel.value.site.siteName.trim()) {
		siteNameInvalid.value = true
		activeSection.value = 'site'
		message.error('请输入网站名称')
		return
	}

	const values = cloneSettings(formModel.value)

	saving.value = true
	try {
		await settingsApi.update(values)
		message.success('设置保存成功')
		originalValues.value = cloneSettings(values)
		hasChanges.value = false
	} catch (error) {
		message.error(errorMessage(error, '设置保存失败'))
	} finally {
		saving.value = false
	}
}

onMounted(fetchSettings)
</script>

<template>
	<div class="page">
		<header class="page-head">
			<div>
				<h1 class="page-title">系统设置</h1>
				<p class="page-desc">站点信息、功能开关与 SEO 配置</p>
			</div>
			<div class="toolbar">
				<mao-badge v-if="hasChanges" tone="accent">未保存</mao-badge>
				<mao-button
					variant="bordered"
					:disabled="(!hasChanges || saving) || undefined"
					@click="handleReset"
				>
					<mao-icon name="rotate-ccw" :size="15" />
					重置
				</mao-button>
				<mao-button
					color="primary"
					:loading="saving || undefined"
					:disabled="(!hasChanges || saving) || undefined"
					@click="handleSave"
				>
					<mao-icon name="save" :size="15" />
					保存
				</mao-button>
			</div>
		</header>

		<div class="settings-layout">
			<aside v-if="!isMobile" class="settings-nav">
				<div class="nav-head">
					<mao-icon name="settings" :size="15" />
					<span class="strong">设置分组</span>
				</div>
				<mao-menu
					:value="activeSection"
					class="settings-menu"
					label="设置分组"
					@mt-change="handleSectionChange"
				>
					<mao-menu-item
						v-for="section in sections"
						:key="section.key"
						:value="section.key"
						:icon="section.icon"
					>
						<span class="nav-label">
							<span class="nav-title">{{ section.title }}</span>
							<span class="nav-desc">{{ section.description }}</span>
						</span>
					</mao-menu-item>
				</mao-menu>
			</aside>

			<mao-card v-else padding="sm">
				<button type="button" class="nav-trigger" @click="drawerOpen = true">
					<mao-icon :name="currentSection.icon" :size="16" />
					<span class="strong">{{ currentSection.title }}</span>
					<mao-icon name="chevrons-up-down" :size="16" class="faint" />
				</button>
			</mao-card>

			<div class="settings-main">
				<mao-card padding="lg">
					<div class="panel-head">
						<div class="row">
							<mao-icon :name="currentSection.icon" :size="18" />
							<span class="section-title">{{ currentSection.title }}</span>
						</div>
						<span v-if="hasChanges" class="tiny faint">有未保存的更改</span>
					</div>
					<p class="section-desc">{{ currentSection.description }}</p>

					<mao-skeleton v-if="loading" lines="6" animated />
					<template v-else>
						<section v-show="activeSection === 'site'" class="section-body">
							<mao-notice tone="blue" title="站点信息">
								这些信息会显示在网站标题、搜索引擎结果和社交媒体分享中。
							</mao-notice>
							<div class="form-grid">
								<mao-text-field
									label="网站名称"
									hint="显示在浏览器标签页和网站头部"
									placeholder="例如：林影心语"
									maxlength="50"
									show-count
									required
									:value="formModel.site.siteName"
									:invalid="siteNameInvalid || undefined"
									@mt-input="setSite('siteName', $event)"
								/>
								<mao-text-field
									label="网站 Logo"
									hint="建议尺寸 200×60 像素"
									placeholder="https://example.com/logo.png"
									:value="formModel.site.siteLogo"
									@mt-input="setSite('siteLogo', $event)"
								/>
								<mao-text-field
									class="span-2"
									label="网站描述"
									hint="用于搜索引擎优化"
									placeholder="简要描述你的网站，不超过 200 字"
									maxlength="200"
									show-count
									:value="formModel.site.siteDescription"
									@mt-input="setSite('siteDescription', $event)"
								/>
								<mao-text-field
									label="网站图标 Favicon"
									hint="建议尺寸 32×32 或 64×64 像素"
									placeholder="https://example.com/favicon.ico"
									:value="formModel.site.favicon"
									@mt-input="setSite('favicon', $event)"
								/>
								<mao-text-field
									label="ICP 备案号"
									hint="中国大陆网站必填"
									placeholder="例如：京 ICP 备 12345678 号"
									:value="formModel.site.icp"
									@mt-input="setSite('icp', $event)"
								/>
							</div>
						</section>

						<section v-show="activeSection === 'feature'" class="section-body">
							<mao-notice tone="blue" title="功能开关">
								控制网站各项功能模块的开启与关闭，保存后立即生效。
							</mao-notice>
							<div class="switch-group">
								<div class="switch-row">
									<div>
										<div class="strong">允许用户注册</div>
										<div class="tiny muted">关闭后新用户将无法注册，仅允许现有用户登录</div>
									</div>
									<mao-switch
										:checked="formModel.feature.enableRegistration || undefined"
										aria-label="允许用户注册"
										@mt-change="setFeature('enableRegistration', $event)"
									/>
								</div>
								<div class="switch-row">
									<div>
										<div class="strong">允许评论</div>
										<div class="tiny muted">关闭后全站将禁用评论功能</div>
									</div>
									<mao-switch
										:checked="formModel.feature.enableComments || undefined"
										aria-label="允许评论"
										@mt-change="setFeature('enableComments', $event)"
									/>
								</div>
								<div class="switch-row">
									<div>
										<div class="strong">邮箱验证</div>
										<div class="tiny muted">开启后新用户注册需要验证邮箱</div>
									</div>
									<mao-switch
										:checked="formModel.feature.enableEmailVerify || undefined"
										aria-label="邮箱验证"
										@mt-change="setFeature('enableEmailVerify', $event)"
									/>
								</div>
								<div class="switch-row">
									<div>
										<div class="strong">验证码</div>
										<div class="tiny muted">开启后在登录和注册时需要输入验证码</div>
									</div>
									<mao-switch
										:checked="formModel.feature.enableCaptcha || undefined"
										aria-label="验证码"
										@mt-change="setFeature('enableCaptcha', $event)"
									/>
								</div>
							</div>
						</section>

						<section v-show="activeSection === 'seo'" class="section-body">
							<mao-notice tone="blue" title="SEO 设置">优化搜索引擎收录与展示效果。</mao-notice>
							<div class="form-grid">
								<mao-text-field
									label="网站关键词"
									hint="多个关键词用逗号分隔"
									placeholder="博客, 内容管理, 写作"
									:value="formModel.seo.siteKeywords"
									@mt-input="setSeo('siteKeywords', $event)"
								/>
								<mao-text-field
									label="网站作者"
									placeholder="网站所有者或组织名称"
									:value="formModel.seo.siteAuthor"
									@mt-input="setSeo('siteAuthor', $event)"
								/>
								<mao-text-area
									class="span-2 mono"
									label="Robots.txt"
									hint="控制搜索引擎爬虫的访问规则"
									rows="6"
									placeholder="User-agent: *&#10;Allow: /&#10;Disallow: /admin/&#10;Sitemap: https://example.com/sitemap.xml"
									:value="formModel.seo.robotsTxt"
									@mt-input="setSeo('robotsTxt', $event)"
								/>
							</div>
							<div class="switch-group">
								<div class="switch-row">
									<div>
										<div class="strong">启用站点地图</div>
										<div class="tiny muted">自动生成 sitemap.xml 供搜索引擎抓取</div>
									</div>
									<mao-switch
										:checked="formModel.seo.enableSitemap || undefined"
										aria-label="启用站点地图"
										@mt-change="setSitemap($event)"
									/>
								</div>
							</div>
						</section>
					</template>
				</mao-card>

				<mao-card id="danger-card" padding="md">
					<div class="row-between">
						<div class="row">
							<mao-icon name="triangle-alert" :size="17" class="danger-icon" />
							<div>
								<div class="strong">危险操作</div>
								<div class="tiny muted">将所有设置恢复为系统默认值</div>
							</div>
						</div>
						<mao-button
							color="danger"
							variant="bordered"
							size="sm"
							@click="handleResetToDefault"
						>
							重置为默认
						</mao-button>
					</div>
				</mao-card>
			</div>
		</div>

		<mao-drawer
			:open="drawerOpen || undefined"
			heading="系统设置"
			placement="start"
			size="sm"
			closable
			@mt-close="drawerOpen = false"
		>
			<mao-menu
				:value="activeSection"
				class="settings-menu"
				label="设置分组"
				@mt-change="handleSectionChange"
			>
				<mao-menu-item
					v-for="section in sections"
					:key="section.key"
					:value="section.key"
					:icon="section.icon"
				>
					<span class="nav-label">
						<span class="nav-title">{{ section.title }}</span>
						<span class="nav-desc">{{ section.description }}</span>
					</span>
				</mao-menu-item>
			</mao-menu>
		</mao-drawer>
	</div>
</template>

<style scoped>
.settings-layout {
	display: flex;
	gap: var(--mt-space-4);
	align-items: flex-start;
	min-width: 0;
}

.settings-nav {
	width: 232px;
	flex: none;
	position: sticky;
	top: 84px;
	padding: var(--mt-space-2);
	border: 1px solid var(--mt-border);
	border-radius: var(--mt-radius-lg);
	background: var(--mt-surface);
}

.nav-head {
	display: flex;
	align-items: center;
	gap: 8px;
	padding: 6px 10px 10px;
	color: var(--mt-muted);
}

.settings-menu {
	border: 0;
	background: transparent;
}

.settings-menu :deep(.mt-menu-item) {
	align-items: flex-start;
	padding: 9px 10px;
}

.nav-label {
	display: flex;
	flex-direction: column;
	gap: 2px;
	min-width: 0;
}

.nav-title {
	font-size: var(--mt-fs-md);
	font-weight: 500;
}

.nav-desc {
	font-size: var(--mt-fs-xs);
	color: var(--mt-muted);
	line-height: 1.35;
}

.nav-trigger {
	display: flex;
	align-items: center;
	gap: 10px;
	width: 100%;
	padding: 4px 2px;
	border: 0;
	background: transparent;
	color: var(--mt-text);
	font: inherit;
	cursor: pointer;
}

.nav-trigger .faint {
	margin-left: auto;
}

.settings-main {
	flex: 1;
	min-width: 0;
	display: flex;
	flex-direction: column;
	gap: var(--mt-space-4);
}

.panel-head {
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: var(--mt-space-3);
}

.panel-head .row {
	align-items: center;
	gap: 8px;
}

.section-body {
	display: flex;
	flex-direction: column;
	gap: var(--mt-space-4);
	margin-top: var(--mt-space-3);
}

.span-2 {
	grid-column: 1 / -1;
}

.mono {
	font-family: var(--mt-mono);
}

.switch-group {
	display: flex;
	flex-direction: column;
	border: 1px solid var(--mt-border);
	border-radius: var(--mt-radius);
	overflow: hidden;
}

.switch-row {
	display: flex;
	align-items: center;
	justify-content: space-between;
	gap: var(--mt-space-4);
	padding: 12px 14px;
}

.switch-row + .switch-row {
	border-top: 1px solid var(--mt-border);
}

.danger-icon {
	color: var(--mt-danger);
}

#danger-card {
	border-color: color-mix(in srgb, var(--mt-danger) 42%, var(--mt-border));
}

@media (max-width: 768px) {
	.settings-layout {
		flex-direction: column;
		gap: var(--mt-space-3);
	}

	.settings-main {
		width: 100%;
	}
}
</style>
