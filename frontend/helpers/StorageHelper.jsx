export default class StorageHelper {

    static isClient() {
        return typeof window !== 'undefined';
    }

    static get(key) {

        if (!this.isClient()) return null;

        const data = localStorage.getItem(key);

        return data || null;

    }

    static set(key, data) {

        if (this.isClient()) {
            localStorage.setItem(key, data);
        }

    }

    static remove(key) {

        if (this.isClient()) {
            localStorage.removeItem(key);
        }

    }

}