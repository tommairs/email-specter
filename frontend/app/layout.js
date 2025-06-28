import 'bootstrap/dist/css/bootstrap.min.css';
import './bootstrap.min.css';
import "./globals.css";
import Structure from "@/app/structure";
import {Toaster} from "react-hot-toast";
import {UserProvider} from "@/providers/UserProvider";

export const metadata = {
    title: "Email Specter",
    description: "A tool to analyze your KumoMTA email logs and identify any issues or anomalies.",
};

export default function RootLayout({children}) {

    return (
        <html lang="en">
        <head>
            <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/6.7.2/css/all.min.css"/>
        </head>
        <body className={`antialiased`}>
        <Toaster position="bottom-center" containerStyle={{zIndex: 10000}}/>
        <UserProvider>
            <Structure>
                {children}
            </Structure>
        </UserProvider>
        </body>
        </html>
    );

}
