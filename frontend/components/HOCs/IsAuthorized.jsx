'use client';
import {useRouter} from 'next/navigation';
import {useEffect} from 'react';
import {useUser} from "@/providers/UserProvider";

function isAuthorized(WrappedComponent) {

    return function WithAuthorization(props) {

        const router = useRouter();
        const {isAuth} = useUser();

        useEffect(() => {

            if (!isAuth) {
                router.push("/auth/login");
            }

        }, [isAuth]);

        return <WrappedComponent {...props} />;

    };

}

export default isAuthorized;