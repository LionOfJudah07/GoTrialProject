import { call } from '../api.js';
import { showToast } from '../components/toast.js';
import { t } from '../i18n.js';

export function setupPasswordModal() {
    const modal = document.getElementById('password-modal');
    const input = document.getElementById('master-password');
    const remember = document.getElementById('remember-me');
    const btn = document.getElementById('unlock-btn');
    
    btn.addEventListener('click', async () => {
        const pwd = input.value;
        if (!pwd) return;
        try {
            await call('SetupMasterPassword', pwd, remember.checked);
        } catch (err) {
            showToast(t('invalid_password'), 'error');
        }
    });
    
    input.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') btn.click();
    });
}