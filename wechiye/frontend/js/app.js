import { initRouter } from './router.js';
import { initI18n, t } from './i18n.js';
import { showToast } from './components/toast.js';
import { setupPasswordModal } from './pages/setup.js';

window.wails = window.wails || {};

async function init() {
    lucide.createIcons();
    await initI18n();
    
    const themeToggle = document.getElementById('theme-toggle');
    themeToggle.addEventListener('click', () => {
        document.documentElement.classList.toggle('dark');
        localStorage.setItem('theme', document.documentElement.classList.contains('dark') ? 'dark' : 'light');
        lucide.createIcons();
    });
    if (localStorage.getItem('theme') === 'dark') {
        document.documentElement.classList.add('dark');
    }
    
    setupPasswordModal();
    initRouter();
    
    if (window.wails?.Events) {
        window.wails.Events.On('needs-password', () => {
            document.getElementById('password-modal').classList.remove('hidden');
        });
        window.wails.Events.On('db-ready', () => {
            document.getElementById('password-modal').classList.add('hidden');
            showToast(t('database_unlocked'), 'success');
            window.dispatchEvent(new HashChangeEvent('hashchange'));
        });
    }
}

init();