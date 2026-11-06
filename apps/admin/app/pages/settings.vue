<script setup lang="ts">
import { computed, h, nextTick, onMounted, ref, watch } from 'vue'
import { useMediaQuery } from '@vueuse/core'
import type { FormInstance } from 'ant-design-vue'
import {
	ExclamationCircleOutlined,
	FileSearchOutlined,
	GlobalOutlined,
	MenuOutlined,
	ReloadOutlined,
	SaveOutlined,
	SettingOutlined,
} from '@ant-design/icons-vue'
import { errorMessage } from '~/api/client'
import { settingsApi } from '~/api'
import type { AllSettings, SettingsSection } from '~/types/domain'

definePageMeta({ layout: 'admin' })

const { message } = useNotify()
const isMobile = useMediaQuery('(max-width: 768px)')

const sectionConfigs = [
	{ key: 'site', title: '站点信息', description: '网站名称、Logo、描述等基本信息' },
	{ key: 'feature', title: '功能开关', description: '注册、评论、验证码等功能控制' },
	{ key: 'seo', title: 'SEO设置', description: '关键词、作者、站点地图' },
] as const

const iconMap = {
	site: GlobalOutlined,
	feature: SettingOutlined,
	seo: FileSearchOutlined,
}

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
const drawerVisible = ref(false)
const originalValues = ref<AllSettings>(cloneSettings(defaultSettings))
const formRef = ref<FormInstance>()
const formModel = ref<AllSettings>(cloneSettings(defaultSettings))
let applying = false

watch(
	formModel,
	() => {
		if (!applying) hasChanges.value = true
	},
	{ deep: true },
)

const currentConfig = computed(
	() => sectionConfigs.find((item) => item.key === activeSection.value) || sectionConfigs[0],
)

const menuItems = sectionConfigs.map((config) => ({
	key: config.key,
	icon: () => h(iconMap[config.key]),
	label: h('div', { class: 'section-label' }, [
		h('div', { class: 'section-label-title' }, config.title),
		h('div', { class: 'section-label-desc' }, config.description),
	]),
}))

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

async function fetchSettings() {
	loading.value = true
	try {
		const data = await settingsApi.get()
		const merged = mergeSettings(data)
		originalValues.value = cloneSettings(merged)
		applySettings(merged)
		hasChanges.value = false
	} catch {
		message.warning('无法连接到服务器, 使用本地默认配置')
		originalValues.value = cloneSettings(defaultSettings)
		applySettings(defaultSettings)
		hasChanges.value = false
	} finally {
		loading.value = false
	}
}

function handleSectionChange(key: string | number) {
	activeSection.value = String(key) as SettingsSection
	if (isMobile.value) drawerVisible.value = false
}

function handleReset() {
	applySettings(originalValues.value)
	hasChanges.value = false
	message.info('已重置为上次保存的设置')
}

function handleResetToDefault() {
	applySettings(defaultSettings)
	hasChanges.value = true
	message.info('已重置为默认设置')
}

async function handleSave() {
	try {
		await formRef.value?.validate()
	} catch {
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
	<div class="settings-page" :class="{ mobile: isMobile }">
		<template v-if="isMobile">
			<a-card size="small">
				<div class="mobile-header">
					<a-typography-title :level="5" :style="{ margin: 0 }">
						<SettingOutlined :style="{ marginRight: '8px' }" />
						系统设置
					</a-typography-title>
					<a-button :icon="h(MenuOutlined)" @click="drawerVisible = true">
						{{ currentConfig.title }}
					</a-button>
				</div>
			</a-card>

			<a-drawer
				:open="drawerVisible"
				placement="left"
				:width="280"
				:body-style="{ padding: 0 }"
				@close="drawerVisible = false"
			>
				<template #title>
					<SettingOutlined :style="{ marginRight: '8px' }" />
					系统设置
				</template>
				<a-menu
					mode="inline"
					:items="menuItems"
					:selected-keys="[activeSection]"
					class="settings-menu"
					@click="({ key }: { key: string | number }) => handleSectionChange(key)"
				/>
			</a-drawer>
		</template>

		<a-card v-else class="settings-nav" :body-style="{ padding: 0 }">
			<div class="settings-nav-header">
				<a-typography-title :level="5" :style="{ margin: 0 }">
					<SettingOutlined :style="{ marginRight: '8px' }" />
					系统设置
				</a-typography-title>
			</div>
			<a-menu
				mode="inline"
				:items="menuItems"
				:selected-keys="[activeSection]"
				class="settings-menu"
				@click="({ key }: { key: string | number }) => handleSectionChange(key)"
			/>
		</a-card>

		<div class="settings-content">
			<a-form
				ref="formRef"
				:model="formModel"
				layout="vertical"
				:disabled="loading"
				:style="{ maxWidth: isMobile ? '100%' : '800px' }"
			>
				<a-card :loading="loading">
					<template #title>
						<a-space>
							<component :is="iconMap[activeSection]" />
							<span>{{ currentConfig.title }}</span>
							<a-badge v-if="hasChanges" dot color="red">
								<span class="dirty-text">有未保存的更改</span>
							</a-badge>
						</a-space>
					</template>
					<template #extra>
						<a-space :size="8" wrap>
							<a-tooltip title="重置为上次保存的设置">
								<a-button
									:icon="h(ReloadOutlined)"
									:disabled="!hasChanges || saving"
									:size="isMobile ? 'small' : 'middle'"
									@click="handleReset"
								>
									<template v-if="!isMobile">重置</template>
								</a-button>
							</a-tooltip>
							<a-button
								type="primary"
								:icon="h(SaveOutlined)"
								:loading="saving"
								:disabled="!hasChanges"
								:size="isMobile ? 'small' : 'middle'"
								@click="handleSave"
							>
								<template v-if="!isMobile">保存</template>
							</a-button>
						</a-space>
					</template>

					<div v-show="activeSection === 'site'">
						<a-space direction="vertical" :style="{ width: '100%' }" :size="16">
							<a-alert
								message="站点信息设置"
								description="配置网站的基本信息, 这些信息将显示在网站标题、搜索引擎结果和社交媒体分享中。"
								type="info"
								show-icon
							/>
							<a-row :gutter="[24, 16]">
								<a-col :xs="24" :md="12">
									<a-form-item
										label="网站名称"
										:name="['site', 'siteName']"
										:rules="[{ required: true, message: '请输入网站名称' }]"
										tooltip="显示在浏览器标签页和网站头部"
									>
										<a-input
											v-model:value="formModel.site.siteName"
											placeholder="例如：我的CMS系统"
											:maxlength="50"
											show-count
										/>
									</a-form-item>
								</a-col>
								<a-col :xs="24" :md="12">
									<a-form-item
										label="网站Logo"
										:name="['site', 'siteLogo']"
										tooltip="建议尺寸 200x60 像素"
									>
										<a-input
											v-model:value="formModel.site.siteLogo"
											placeholder="https://example.com/logo.png"
										/>
									</a-form-item>
								</a-col>
							</a-row>
							<a-form-item
								label="网站描述"
								:name="['site', 'siteDescription']"
								tooltip="用于搜索引擎优化（SEO）"
							>
								<a-textarea
									v-model:value="formModel.site.siteDescription"
									:rows="3"
									placeholder="简要描述您的网站, 不超过200字"
									:maxlength="200"
									show-count
								/>
							</a-form-item>
							<a-row :gutter="[24, 16]">
								<a-col :xs="24" :md="12">
									<a-form-item
										label="网站图标 (Favicon)"
										:name="['site', 'favicon']"
										tooltip="建议尺寸 32x32 或 64x64 像素"
									>
										<a-input
											v-model:value="formModel.site.favicon"
											placeholder="https://example.com/favicon.ico"
										/>
									</a-form-item>
								</a-col>
								<a-col :xs="24" :md="12">
									<a-form-item label="ICP备案号" :name="['site', 'icp']" tooltip="中国大陆网站必填">
										<a-input
											v-model:value="formModel.site.icp"
											placeholder="例如：京ICP备12345678号"
										/>
									</a-form-item>
								</a-col>
							</a-row>
						</a-space>
					</div>

					<div v-show="activeSection === 'feature'">
						<a-space direction="vertical" :style="{ width: '100%' }" :size="16">
							<a-alert
								message="功能开关"
								description="控制网站的各项功能模块的开启和关闭。"
								type="info"
								show-icon
							/>
							<a-card size="small" title="用户功能">
								<a-space direction="vertical" :style="{ width: '100%' }">
									<a-form-item
										label="允许用户注册"
										:name="['feature', 'enableRegistration']"
										value-prop-name="checked"
										style="margin-bottom: 8px"
									>
										<a-switch
											v-model:checked="formModel.feature.enableRegistration"
											checked-children="开启"
											un-checked-children="关闭"
										/>
									</a-form-item>
									<a-typography-text type="secondary" :style="{ fontSize: '12px' }">
										关闭后新用户将无法注册, 仅允许现有用户登录
									</a-typography-text>
								</a-space>
							</a-card>
							<a-card size="small" title="内容功能">
								<a-space direction="vertical" :style="{ width: '100%' }">
									<a-form-item
										label="允许评论"
										:name="['feature', 'enableComments']"
										value-prop-name="checked"
										style="margin-bottom: 8px"
									>
										<a-switch
											v-model:checked="formModel.feature.enableComments"
											checked-children="开启"
											un-checked-children="关闭"
										/>
									</a-form-item>
									<a-typography-text type="secondary" :style="{ fontSize: '12px' }">
										关闭后全站将禁用评论功能
									</a-typography-text>
								</a-space>
							</a-card>
							<a-card size="small" title="安全功能">
								<a-space direction="vertical" :style="{ width: '100%' }">
									<a-form-item
										label="邮箱验证"
										:name="['feature', 'enableEmailVerify']"
										value-prop-name="checked"
										style="margin-bottom: 8px"
									>
										<a-switch
											v-model:checked="formModel.feature.enableEmailVerify"
											checked-children="开启"
											un-checked-children="关闭"
										/>
									</a-form-item>
									<a-typography-text type="secondary" :style="{ fontSize: '12px' }">
										开启后新用户注册需要验证邮箱
									</a-typography-text>
									<a-form-item
										label="验证码"
										:name="['feature', 'enableCaptcha']"
										value-prop-name="checked"
										style="margin: 16px 0 8px"
									>
										<a-switch
											v-model:checked="formModel.feature.enableCaptcha"
											checked-children="开启"
											un-checked-children="关闭"
										/>
									</a-form-item>
									<a-typography-text type="secondary" :style="{ fontSize: '12px' }">
										开启后在登录和注册时需要输入验证码
									</a-typography-text>
								</a-space>
							</a-card>
						</a-space>
					</div>

					<div v-show="activeSection === 'seo'">
						<a-space direction="vertical" :style="{ width: '100%' }" :size="16">
							<a-alert
								message="SEO设置"
								description="优化搜索引擎收录和展示效果。"
								type="info"
								show-icon
							/>
							<a-row :gutter="[24, 16]">
								<a-col :xs="24" :md="12">
									<a-form-item
										label="网站关键词"
										:name="['seo', 'siteKeywords']"
										tooltip="多个关键词用逗号分隔"
									>
										<a-input
											v-model:value="formModel.seo.siteKeywords"
											placeholder="CMS, 博客, 内容管理"
										/>
									</a-form-item>
								</a-col>
								<a-col :xs="24" :md="12">
									<a-form-item label="网站作者" :name="['seo', 'siteAuthor']">
										<a-input
											v-model:value="formModel.seo.siteAuthor"
											placeholder="网站所有者或组织名称"
										/>
									</a-form-item>
								</a-col>
							</a-row>
							<a-form-item
								label="Robots.txt"
								:name="['seo', 'robotsTxt']"
								tooltip="控制搜索引擎爬虫的访问规则"
							>
								<a-textarea
									v-model:value="formModel.seo.robotsTxt"
									:rows="6"
									:style="{ fontFamily: 'monospace' }"
									placeholder="User-agent: *&#10;Allow: /&#10;Disallow: /admin/&#10;Sitemap: https://example.com/sitemap.xml"
								/>
							</a-form-item>
							<a-form-item
								label="启用站点地图"
								:name="['seo', 'enableSitemap']"
								value-prop-name="checked"
							>
								<a-switch
									v-model:checked="formModel.seo.enableSitemap"
									checked-children="开启"
									un-checked-children="关闭"
								/>
							</a-form-item>
						</a-space>
					</div>
				</a-card>

				<a-divider />

				<a-card size="small" title="危险操作" :style="{ borderColor: '#ff4d4f' }">
					<a-space wrap>
						<ExclamationCircleOutlined :style="{ color: '#ff4d4f' }" />
						<a-typography-text type="secondary">重置所有设置为默认值</a-typography-text>
						<a-button danger size="small" @click="handleResetToDefault">重置为默认</a-button>
					</a-space>
				</a-card>
			</a-form>
		</div>
	</div>
</template>

<style scoped>
.settings-page {
	height: calc(100vh - 112px);
	display: flex;
	gap: 24px;
	overflow: hidden;
}

.settings-page.mobile {
	height: calc(100vh - 140px);
	flex-direction: column;
	gap: 12px;
}

.settings-nav {
	width: 260px;
	flex-shrink: 0;
}

.settings-nav-header {
	padding: 16px 20px;
	border-bottom: 1px solid #f0f0f0;
}

.settings-content {
	flex: 1;
	overflow: auto;
}

.mobile-header {
	display: flex;
	align-items: center;
	justify-content: space-between;
}

.dirty-text {
	font-size: 12px;
	color: #999;
}

:deep(.settings-menu) {
	border: none;
	padding: 8px 0;
}

:deep(.settings-menu .ant-menu-item) {
	height: auto;
	padding: 12px 16px;
	margin: 4px 0;
	border-radius: 0;
}

:deep(.section-label-title) {
	font-weight: 500;
	font-size: 14px;
	margin-bottom: 4px;
	color: #262626;
}

:deep(.section-label-desc) {
	font-size: 12px;
	color: #8c8c8c;
	line-height: 1.3;
}

:deep(.ant-menu-item-selected .section-label-title) {
	color: #1890ff;
	font-weight: 600;
}

:deep(.ant-menu-item-selected .section-label-desc) {
	color: #69c0ff;
}
</style>
