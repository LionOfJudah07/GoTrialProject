import { loadDashboard } from './pages/dashboard.js';
import { loadTransactions } from './pages/transactions.js';
import { loadBudgets } from './pages/budgets.js';
import { loadProfile } from './pages/profile.js';
import { loadSettings } from './pages/settings.js';

const routes = {
    '#dashboard': loadDashboard,
    '#transactions': loadTransactions,
    '#budgets': loadBudgets,
    '#profile': loadProfile,
    '#settings': loadSettings,
};

export function initRouter() {
    window.addEventListener('hashchange', handleRoute);
    handleRoute();
}

async function handleRoute() {
    const hash = window.location.hash || '#dashboard';
    const loader = routes[hash];
    const container = document.getElementById('page-container');
    if (loader) {
        container.innerHTML = '<div class="text-center py-10">Loading...</div>';
        await loader(container);
        document.querySelectorAll('.nav-link').forEach(link => {
            link.classList.remove('active');
            if (link.getAttribute('href') === hash) {
                link.classList.add('active');
            }
        });
    }
}