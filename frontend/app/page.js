"use client";

import {useUser} from "@/providers/UserProvider";
import {useEffect} from "react";
import {useRouter} from "next/navigation";

export default function Home() {

    const {isAuth, refreshUser} = useUser();
    const router = useRouter();

    const checkIfUserIsLoggedIn = async () => {

        if (isAuth) {
            router.push('/messages');
        } else {
            router.push('/auth/login');
        }

    };

    useEffect(() => {
        checkIfUserIsLoggedIn();
    }, []);

}
