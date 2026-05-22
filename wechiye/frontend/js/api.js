export async function call(method, ...args) {
    try {
        return await window.wails.Call(method, ...args);
    } catch (err) {
        console.error(`API call ${method} failed:`, err);
        throw err;
    }
}

export const getUser = () => call('GetUser');
export const updateUser = (user) => call('UpdateUser', user);
export const updateAvatar = (avatar) => call('UpdateAvatar', avatar);

export const getAccounts = () => call('GetAccounts');
export const addAccount = (name, initBalance) => call('AddAccount', name, initBalance);
export const updateAccount = (id, name, initBalance) => call('UpdateAccount', id, name, initBalance);
export const deleteAccount = (id) => call('DeleteAccount', id);

export const getTransactions = (filter) => call('GetTransactions', filter);
export const addTransaction = (txn) => call('AddTransaction', txn);
export const updateTransaction = (txn) => call('UpdateTransaction', txn);
export const deleteTransaction = (id) => call('DeleteTransaction', id);

export const getBudgets = (monthYear) => call('GetBudgets', monthYear);
export const setBudget = (budget) => call('SetBudget', budget);
export const deleteBudget = (id) => call('DeleteBudget', id);

export const generatePairingQR = () => call('GeneratePairingQR');
export const processScannedQR = (data) => call('ProcessScannedQR', data);
export const getCoupleStatus = () => call('GetCoupleStatus');
export const exportSharedData = () => call('ExportEncryptedSharedData');
export const importSharedData = (b64) => call('ImportEncryptedSharedData', b64);

export const exportBackup = () => call('ExportBackup');
export const restoreBackup = (path) => call('RestoreBackup', path);