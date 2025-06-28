import StorageHelper from './StorageHelper';

export const TokenHelper = {
    setToken: (token) => StorageHelper.set('token', token),
    getToken: () => StorageHelper.get('token'),
    removeToken: () => StorageHelper.remove('token'),
};