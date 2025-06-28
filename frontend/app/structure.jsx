"use client";

import {useEffect, useState} from "react";
import Link from "next/link";
import {useUser} from "@/providers/UserProvider";

export default function Structure({children}) {

    const[isClient, setIsClient] = useState(false);

    const {isAuth} = useUser();

    useEffect(() => {
        require('bootstrap/dist/js/bootstrap.bundle.js');
    }, []);

    useEffect(() => {
        setIsClient(true);
    }, []);

    if (!isClient) {
        return null;
    }

    return (
        <>
            <nav className="navbar navbar-expand-lg bg-dark" data-bs-theme="dark">
                <div className="container">

                    <Link className="navbar-brand" href="/">
                        Email Specter
                    </Link>

                    <button className="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarColor02" aria-controls="navbarColor02" aria-expanded="false" aria-label="Toggle navigation">
                        <span className="navbar-toggler-icon"></span>
                    </button>

                    <div className="collapse navbar-collapse" id="navbarColor02">

                        <ul className="navbar-nav me-auto">

                            <li className="nav-item">
                                <Link className="nav-link" href="/">
                                    Home
                                </Link>
                            </li>

                            <li className="nav-item">
                                <a className="nav-link" href="#">Features</a>
                            </li>

                            <li className="nav-item">
                                <a className="nav-link" href="#">Pricing</a>
                            </li>

                            <li className="nav-item">
                                <a className="nav-link" href="#">About</a>
                            </li>

                        </ul>

                        <ul className="navbar-nav ms-auto">

                            {

                                isAuth ? (

                                    <li className="nav-item">
                                        <Link className="nav-link" href="/auth/logout">
                                            Logout
                                        </Link>
                                    </li>

                                ) : (

                                    <li className="nav-item">
                                        <Link className="nav-link" href="/auth/login">
                                            Login
                                        </Link>
                                    </li>

                                )

                            }

                        </ul>

                    </div>

                </div>
            </nav>

            <main>
                <div className="container mt-4">
                    {children}
                </div>
            </main>

            <footer className="bg-dark text-white mt-4">
                <div className="container">
                    <div className="d-flex justify-content-lg-between p-4">
                        <p className="mb-0">
                            Released under the MIT License
                        </p>
                        <p className="mb-0">
                            <Link href="https://github.com/maileroo/email-specter" className="text-white text-decoration-none me-3">
                                View on GitHub
                            </Link>
                        </p>
                    </div>
                </div>
            </footer>

        </>
    );

}