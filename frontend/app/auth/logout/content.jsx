"use client";

import {useEffect} from "react";
import {useRouter} from "next/navigation";
import {useUser} from "@/providers/UserProvider";

export default function View() {

    const router = useRouter();
    const {logout} = useUser();

    useEffect(() => {
        logout();
        router.push('/auth/login');
    }, []);

}

