import StorageHelper from '../helpers/StorageHelper';
import axios from 'axios';

export default class RequestHelper {

    static async sendAuthenticatedDeleteRequest(endpoint) {

        try {

            let response = await axios.delete(process.env.NEXT_PUBLIC_API_URL + endpoint, {
                headers: {
                    Authorization: 'Bearer ' + StorageHelper.get('token'),
                },
            });

            return {
                code: response.status,
                data: response.data,
            };

        } catch (e) {

            return {
                code: e.response.status,
                data: e.response.data,
            };

        }

    }

    static async sendAuthenticatedPostRequest(endpoint, data) {

        try {

            let response = await axios.post(process.env.NEXT_PUBLIC_API_URL + endpoint, data, {
                headers: {
                    Authorization: 'Bearer ' + StorageHelper.get('token'),
                },
            });

            return {
                code: response.status,
                data: response.data,
            };

        } catch (e) {

            return {
                code: e.response.status,
                data: e.response.data,
            };

        }

    }

    static async sendAuthenticatedPostUploadRequest(endpoint, data, onUploadProgress) {

        try {

            let response = await axios.post(process.env.NEXT_PUBLIC_API_URL + endpoint, data, {
                headers: {
                    Authorization: 'Bearer ' + StorageHelper.get('token'),
                },
                onUploadProgress: (progressEvent) => {
                    if (onUploadProgress) {
                        const percentComplete = Math.round(
                            (progressEvent.loaded * 100) / progressEvent.total
                        );
                        onUploadProgress(percentComplete);
                    }
                },
            });

            return {
                code: response.status,
                data: response.data,
            };

        } catch (e) {

            return {
                code: e.response.status,
                data: e.response.data,
            };

        }

    }

    static async sendAuthenticatedPutRequest(endpoint, data) {

        try {

            let response = await axios.put(process.env.NEXT_PUBLIC_API_URL + endpoint, data, {
                headers: {
                    Authorization: 'Bearer ' + StorageHelper.get('token'),
                },
            });

            return {
                code: response.status,
                data: response.data,
            };

        } catch (e) {

            return {
                code: e.response.status,
                data: e.response.data,
            };

        }

    }

    static async sendAuthenticatedPatchRequest(endpoint, data) {

        try {

            let response = await axios.patch(process.env.NEXT_PUBLIC_API_URL + endpoint, data, {
                headers: {
                    Authorization: 'Bearer ' + StorageHelper.get('token'),
                },
            });

            return {
                code: response.status,
                data: response.data,
            };

        } catch (e) {

            return {
                code: e.response.status,
                data: e.response.data,
            };

        }

    }

    static async sendAuthenticatedGetRequest(endpoint , token) {

        try {

            let response = await axios.get(process.env.NEXT_PUBLIC_API_URL + endpoint, {
                headers: {
                    Authorization: 'Bearer ' + (token || StorageHelper.get('token')),
                },
            });

            return {
                code: response.status,
                data: response.data,
            };

        } catch (e) {

            return {
                code: e.response?.status,
                data: e.response.data,
            };

        }

    }

    static async sendUnauthenticatedPostRequest(endpoint, data) {

        try {

            let response = await axios.post(process.env.NEXT_PUBLIC_API_URL + endpoint, data);

            return {
                code: response.status,
                data: response.data,
            };

        } catch (e) {

            return {
                code: e.response.status,
                data: e.response.data,
            };

        }

    }

    static async sendUnauthenticatedGetRequest(endpoint) {

        try {

            let response = await axios.get(process.env.NEXT_PUBLIC_API_URL + endpoint);

            return {
                code: response.status,
                data: response.data,
            };

        } catch (e) {

            return {
                code: e.response.status,
                data: e.response.data,
            };

        }

    }

}