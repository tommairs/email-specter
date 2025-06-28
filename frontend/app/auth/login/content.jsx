"use client";

import FormHelper from "@/helpers/FormHelper";
import {useEffect} from "react";
import RequestHelper from "@/helpers/RequestHelper";
import {ToastHelper} from "@/helpers/ToastHelper";
import {useRouter} from "next/navigation";
import {TokenHelper} from "@/helpers/TokenHelper";
import {useUser} from "@/providers/UserProvider";

export default function View() {

    const router = useRouter();
    const {isAuth, refreshUser} = useUser();

    const checkIfInstallationIsSetup = async () => {

        const response = await RequestHelper.sendAuthenticatedGetRequest('/auth/can-register');

        const data = response.data;

        if (data['success'] && data['can_register']) {

            router.push('/auth/register');

            ToastHelper.errorToast("You must create an admin user before you can login.");

        }

    };

    const submit = async (e) => {

        const formData = FormHelper.getFormData(e);

        const response = await RequestHelper.sendAuthenticatedPostRequest('/login', {
            email_address: formData.email_address,
            password: formData.password
        });

        const data = response.data;

        if (data.success) {

            ToastHelper.successToast(data.message);
            TokenHelper.setToken(data.data.token);

            refreshUser();

            router.push('/dashboard');

        } else {

            ToastHelper.errorToast(data.message);

        }

    };

    const checkIfUserIsLoggedIn = async () => {

        if (isAuth) {
            router.push('/dashboard');
        }

    };

    useEffect(() => {
        checkIfInstallationIsSetup();
        checkIfUserIsLoggedIn();
    }, []);

    return (
        <div className="container">
            <div className="row">

                <div className="col-md-6 offset-md-3 mt-5">

                    <h1 className="text-center">
                        Login
                    </h1>

                    <form onSubmit={(e) => submit(e)}>

                        <div className="mb-3">
                            <label htmlFor="email_address" className="form-label">Email Address</label>
                            <input type="email" className="form-control" id="email_address" name="email_address"/>
                        </div>

                        <div className="mb-3">
                            <label htmlFor="password" className="form-label">Password</label>
                            <input type="password" className="form-control" id="password" name="password"/>
                        </div>

                        <button type="submit" className="btn btn-dark btn-lg w-100">
                            Login
                        </button>

                    </form>

                </div>

            </div>
        </div>
    );

}

