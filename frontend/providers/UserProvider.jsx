'use client';

import {TokenHelper} from '@/helpers/TokenHelper';
import React, {createContext, useState, useContext, useEffect} from 'react';
import RequestHelper from '@/helpers/RequestHelper';

const UserContext = createContext();

export const UserProvider = ({children}) => {

    const [user, setUser] = useState(null);
    const [isAuth, setIsAuth] = useState(!!TokenHelper.getToken());

    const getUser = () => {
        return user;
    };

    const logout = () => {
        setIsAuth(false);
        setUser(null);
        TokenHelper.removeToken();
    };

    const refreshUser = async () => {

        const response = await RequestHelper.sendAuthenticatedGetRequest('/account');
        const data = response.data;

        if (response.data['success']) {
            setUser(data['user']);
            setIsAuth(true);
        } else {
            setIsAuth(false);
            setUser(null);
        }

    }

    useEffect(() => {

        if (TokenHelper.getToken()) {

            if (getUser()) {
                setIsAuth(true);
            } else {
                refreshUser();
            }

        }

    }, [user]);

    return (
        <UserContext.Provider value={{user, isAuth, refreshUser, getUser, logout}}>
            {children}
        </UserContext.Provider>
    );

};

export const useUser = () => {
    return useContext(UserContext);
};