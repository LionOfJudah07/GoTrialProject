export function formatDate(dateStr) {
    return new Date(dateStr).toLocaleDateString();
}

export function formatCurrency(amount) {
    return new Intl.NumberFormat(undefined, { style: 'currency', currency: 'USD' }).format(amount);
}