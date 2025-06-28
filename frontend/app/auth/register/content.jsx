"use client";

import FormHelper from "@/helpers/FormHelper";
import {useEffect} from "react";
import RequestHelper from "@/helpers/RequestHelper";
import {ToastHelper} from "@/helpers/ToastHelper";
import {useRouter} from "next/navigation";

export default function View() {

    const router = useRouter();

    const submit = async (e) => {

        const formData = FormHelper.getFormData(e);

        const response = await RequestHelper.sendAuthenticatedPostRequest('/register', {
            full_name: formData.full_name,
            email_address: formData.email_address,
            password: formData.password
        });

        const data = response.data;

        if (data.success) {

            ToastHelper.successToast(data.message);

            router.push('/auth/login');

        } else {

            ToastHelper.errorToast(data.message);

        }

    };

    return (
        <div className="container">
            <div className="row">

                <div className="col-md-6 offset-md-3 mt-5">

                    <h1 className="text-center">
                        Register
                    </h1>

                    <p className="mt-4 mb-4">
                        If this installation doesn't have an admin user, you can create one here. Once created, you will be able to log in and manage the application.
                    </p>

                    <form onSubmit={(e) => submit(e)}>

                        <div className="mb-3">
                            <label htmlFor="full_name" className="form-label">Name</label>
                            <input type="text" className="form-control" id="full_name" name="full_name"/>
                        </div>

                        <div className="mb-3">
                            <label htmlFor="email_address" className="form-label">Email Address</label>
                            <input type="email" className="form-control" id="email_address" name="email_address"/>
                        </div>

                        <div className="mb-3">
                            <label htmlFor="password" className="form-label">Password</label>
                            <input type="password" className="form-control" id="password" name="password"/>
                        </div>

                        <button type="submit" className="btn btn-dark btn-lg w-100">
                            Register
                        </button>

                    </form>

                </div>

            </div>
        </div>
    );

}

