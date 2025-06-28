export default class FormHelper {

    static getFormData(e) {

        e.preventDefault();

        const formData = new FormData(e.target);
        const data = {};

        formData.forEach((value, key) => {
            data[key] = value;
        });

        return data;

    }

}