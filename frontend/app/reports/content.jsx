"use client";

import {useEffect, useState} from "react";
import RequestHelper from "@/helpers/RequestHelper";
import GlobalHelper from "@/helpers/GlobalHelper";

export default function Content() {

    // type ReportRequest struct {
    // 	From                       string `json:"from"`
    // 	To                         string `json:"to"`
    // 	SourceIP                   string `json:"source_ip"`
    // 	SourceDomain               string `json:"source_domain"`
    // 	DestinationDomain          string `json:"destination_domain"`
    // 	DestinationService         string `json:"destination_service"`
    // 	KumoMtaClassification      string `json:"kumo_mta_classification"`
    // 	EmailSpecterClassification string `json:"email_specter_classification"`
    // 	EventType                  string `json:"event_type"`
    // 	GroupBy                    string `json:"group_by"`
    // }

    const [entities, setEntities] = useState({
        destination_domains: [],
        destination_services: [],
        source_domains: [],
        source_ips: []
    });

    const [loading, setLoading] = useState(true);

    const loadTopEntries = async () => {

        const data = await GlobalHelper.fetchTopEntities();

        if (data) {

            setEntities({
                destination_domains: (data.destination_domains || []).map(name => ({ name })),
                destination_services: (data.destination_services || []).map(name => ({ name })),
                source_domains: (data.source_domains || []).map(name => ({ name })),
                source_ips: (data.source_ips || []).map(name => ({ name })),
            });

        }

    };

    useEffect(() => {
        loadTopEntries();
    }, []);

    return (
        <div className="container">

            <div className="d-flex align-items-center mb-4">

                <div className="flex align-items-center gap-3 mb-4">
                    <div>
                        <h1 className="text-3xl font-bold text-900 m-0">
                            Reports
                        </h1>
                        <p className="text-600 m-0 mt-1">
                            Generate and view reports on your email traffic, including delivery status, bounce rates, and more.
                        </p>
                    </div>
                </div>

            </div>

            <div>

                asd

            </div>

        </div>
    );

}