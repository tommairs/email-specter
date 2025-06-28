"use client";

import {useEffect, useState} from "react";
import Link from "next/link";
import {useUser} from "@/providers/UserProvider";
import {usePathname} from "next/navigation";

export default function Structure({children}) {

    const [isClient, setIsClient] = useState(false);

    const {isAuth} = useUser();
    const pathname = usePathname();

    useEffect(() => {
        require('bootstrap/dist/js/bootstrap.bundle.js');
    }, []);

    useEffect(() => {
        setIsClient(true);
    }, []);

    if (!isClient) {
        return null;
    }

    const matchRoute = (route, currentPath) => {
        return currentPath.startsWith(route) ? ' active' : '';
    };

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
                                <Link className={"nav-link " + matchRoute('/connections', pathname)} href="/connections">
                                    Connections
                                </Link>
                            </li>

                            <li className="nav-item">
                                <Link className={"nav-link " + matchRoute('/statistics', pathname)} href="/statistics">
                                    Statistics
                                </Link>
                            </li>

                            <li className="nav-item">
                                <Link className={"nav-link " + matchRoute('/reports', pathname)} href="/reports">
                                    Reports
                                </Link>
                            </li>

                            <li className="nav-item dropdown">

                                <a className={"nav-link dropdown-toggle" + matchRoute('/providers', pathname)} href="#" role="button" data-bs-toggle="dropdown" aria-expanded="false">
                                    Provider Insights
                                </a>

                                <ul className="dropdown-menu">

                                    <li>
                                        <Link className="dropdown-item" href="/providers/event-statistics">
                                            Event Statistics
                                        </Link>
                                    </li>

                                    <li>
                                        <Link className="dropdown-item" href="/providers/bounce-types">
                                            Bounce Types
                                        </Link>
                                    </li>

                                </ul>

                            </li>

                            <li className="nav-item">
                                <Link className={"nav-link " + matchRoute('/top-entities', pathname)} href="/top-entities">
                                    Top Entities
                                </Link>
                            </li>

                            <li className="nav-item">
                                <Link className={"nav-link " + matchRoute('/messages', pathname)} href="/messages">
                                    Messages
                                </Link>
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