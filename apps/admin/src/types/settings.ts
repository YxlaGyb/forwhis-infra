export interface SiteSettings {
    siteName: string;
    siteDescription: string;
    siteLogo: string;
    favicon: string;
    icp: string;
}

export interface ContactSettings {
    adminEmail: string;
    contactPhone: string;
    contactAddress: string;
}

export interface SystemSettings {
    language: string;
    timezone: string;
    dateFormat: string;
    timeFormat: string;
}

export interface FeatureSettings {
    enableRegistration: boolean;
    enableComments: boolean;
    enableEmailVerify: boolean;
    enableCaptcha: boolean;
}

export interface SecuritySettings {
    maintenanceMode: boolean;
    loginFailLimit: number;
    sessionTimeout: number;
    passwordMinLength: number;
}

export interface SeoSettings {
    siteKeywords: string;
    siteAuthor: string;
    robotsTxt: string;
    enableSitemap: boolean;
}

export interface AllSettings {
    site: SiteSettings;
    feature: FeatureSettings;
    seo: SeoSettings;
}

export type SettingsSection = 'site' | 'feature' | 'seo';

export interface SettingsSectionConfig {
    key: SettingsSection;
    title: string;
    icon: string;
    description: string;
}
