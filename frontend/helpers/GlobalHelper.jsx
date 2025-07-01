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

    static addCommasToNumber = (number) => {

        if (typeof number !== 'number') {
            return number; // Return as is if not a number
        }

        return number.toLocaleString('en-US'); // Format number with commas

    };

    static formatDate = (date, format) => {

        if (!(date instanceof Date)) {
            date = new Date(date);
        }

        const options = {
            year: 'numeric',
            month: '2-digit',
            day: '2-digit',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit',
            hour12: false,
        };

        return date.toLocaleDateString('en-US', options).replace(/(\d+)\/(\d+)\/(\d+), (\d+):(\d+):(\d+)/, (match, month, day, year, hour, minute, second) => {
            return format.replace('yyyy', year).replace('MM', month).replace('dd', day).replace('HH', hour).replace('mm', minute).replace('ss', second);
        });

    }


}