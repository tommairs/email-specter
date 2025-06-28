export default class GlobalHelper {

    static copyToClipboard(text) {
        navigator.clipboard.writeText(text);
    }

}