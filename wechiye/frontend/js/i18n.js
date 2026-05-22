let translations = {};
let currentLang = localStorage.getItem('lang') || 'en';

export async function initI18n() {
    await loadLanguage(currentLang);
}

async function loadLanguage(lang) {
    try {
        const res = await fetch(`./i18n/${lang}.json`);
        translations = await res.json();
        currentLang = lang;
        localStorage.setItem('lang', lang);
        applyTranslations();
    } catch (e) {
        console.error('Failed to load language', e);
    }
}

function applyTranslations() {
    document.querySelectorAll('[data-i18n]').forEach(el => {
        const key = el.getAttribute('data-i18n');
        el.textContent = translations[key] || key;
    });
}

export function t(key) {
    return translations[key] || key;
}

export function setLanguage(lang) {
    loadLanguage(lang);
}