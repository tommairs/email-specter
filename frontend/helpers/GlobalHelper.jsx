import RequestHelper from "@/helpers/RequestHelper";

export default class GlobalHelper {

    static copyToClipboard(text) {
        navigator.clipboard.writeText(text);
    }

    static fetchTopEntities = async () => {

        const response = await RequestHelper.sendAuthenticatedGetRequest("/reports/top-entities");
        const data = response.data;

        if (data.success) {
            return data.data;
        }

        return null;

    };


}