"use client";

import {useEffect, useState} from "react";
import RequestHelper from "@/helpers/RequestHelper";
import {ToastHelper} from "@/helpers/ToastHelper";
import {useRouter} from "next/navigation";
import {DataTable} from 'primereact/datatable';
import {Column} from 'primereact/column';
import {TabView, TabPanel} from 'primereact/tabview';
import {InputText} from 'primereact/inputtext';
import {Button} from 'primereact/button';
import {ProgressBar} from 'primereact/progressbar';
import {Badge} from 'primereact/badge';
import {Skeleton} from 'primereact/skeleton';
import {Card} from 'primereact/card';
import GlobalHelper from "@/helpers/GlobalHelper";
import Loading from "@/components/Loading";

export default function Content() {

    const [loading, setLoading] = useState(true);

    const [entities, setEntities] = useState({
        destination_domains: [],
        destination_services: [],
        source_domains: [],
        source_ips: []
    });

    const [activeTab, setActiveTab] = useState(0);
    const [globalFilterValue, setGlobalFilterValue] = useState('');

    const fetchEntities = async () => {

        setLoading(true);

        const response = await RequestHelper.sendAuthenticatedGetRequest("/reports/top-entities");
        const data = response.data;

        if (data.success) {

            setEntities({
                destination_domains: (data.data.destination_domains || []).map(name => ({ name })),
                destination_services: (data.data.destination_services || []).map(name => ({ name })),
                source_domains: (data.data.source_domains || []).map(name => ({ name })),
                source_ips: (data.data.source_ips || []).map(name => ({ name })),
            });

        }

        setLoading(false);

    };

    useEffect(() => {
        fetchEntities();
    }, []);

    const copyToClipboard = (text) => {
        GlobalHelper.copyToClipboard(text);
        ToastHelper.successToast('Copied!');
    };

    const onGlobalFilterChange = (e) => {
        setGlobalFilterValue(e.target.value);
    };

    const renderHeader = (title, data) => {

        return (
            <div className="d-flex justify-content-between align-items-center mb-3 flex-wrap gap-2">

                <div className="d-flex align-items-center gap-2">

                    <h5 className="mb-0 me-2">
                        {title}
                    </h5>

                    <span className="badge bg-info text-white">{data.length}</span>

                </div>

                <div className="d-flex align-items-center gap-2">

                    <div className="input-group input-group-sm">

                        <span className="input-group-text bg-light border-end-0">
                            <i className="fa fa-search"></i>
                        </span>

                        <input
                            type="text"
                            className="form-control border-start-0"
                            placeholder="Search..."
                            value={globalFilterValue}
                            onChange={onGlobalFilterChange}
                        />

                    </div>

                </div>

            </div>
        );

    };

    const nameBodyTemplate = (rowData) => {

        return (
            <div className="flex align-items-center gap-2">
                <span className="font-medium">{rowData.name}</span>
            </div>
        );

    };


    const actionBodyTemplate = (rowData) => {

        return (
            <div className="text-center">
                <button className="btn btn-sm btn-dark" onClick={() => copyToClipboard(rowData.name)} title="Copy to clipboard">
                    <i className="fa fa-copy"></i>
                </button>
            </div>
        );

    };

    const getColumnHeader = (tabIndex) => {
        const headers = ['Domain', 'Service', 'Domain', 'IP Address'];
        return headers[tabIndex] || 'Name';
    };

    const tabTitles = [
        'Destination Domains',
        'Destination Services',
        'Source Domains',
        'Source IP Addresses'
    ];

    const entityKeys = [
        'destination_domains',
        'destination_services',
        'source_domains',
        'source_ips'
    ];

    if (loading) {
        return <Loading/>;
    }

    return (
        <div className={"container"}>

            <div className="flex align-items-center gap-3 mb-4">
                <i className="pi pi-chart-bar text-6xl text-primary"></i>
                <div>
                    <h1 className="text-3xl font-bold text-900 m-0">Top Entities</h1>
                    <p className="text-600 m-0 mt-1">
                        This report provides insights into the most frequently occurring entities in your email traffic. Use this data to identify trends and potential issues in your emails.
                    </p>
                </div>
            </div>

            <div className="card">

                <ul className="nav nav-tabs">

                    {tabTitles.map((title, index) => (
                        <li className="nav-item" key={index}>
                            <button className={`nav-link ${activeTab === index ? 'active' : ''}`} onClick={() => {
                                setActiveTab(index);
                                setGlobalFilterValue('');
                            }}>
                                <span className="me-2">{title}</span>
                                <Badge value={entities[entityKeys[index]]?.length || 0} severity="secondary"/>
                            </button>
                        </li>
                    ))}

                </ul>

                <div className="tab-content p-3 border border-top-0">

                    {tabTitles.map((title, index) => (

                        <div key={index} className={`tab-pane fade ${activeTab === index ? 'show active' : ''}`}>

                            <DataTable
                                value={entities[entityKeys[index]] || []}
                                paginator
                                rows={10}
                                rowsPerPageOptions={[5, 10, 25, 50]}
                                globalFilter={globalFilterValue}
                                header={renderHeader(title, entities[entityKeys[index]] || [])}
                                emptyMessage="There's no data to display."
                                className="p-datatable-gridlines"
                                size="small"
                                stripedRows
                                sortMode="multiple"
                                removableSort
                                showGridlines
                            >

                                <Column
                                    field="name"
                                    header={getColumnHeader(index)}
                                    sortable
                                    body={nameBodyTemplate}
                                    style={{minWidth: '12rem'}}
                                />

                                <Column
                                    header="Actions"
                                    body={actionBodyTemplate}
                                    style={{width: '6rem'}}
                                    className="text-center"
                                />

                            </DataTable>

                        </div>

                    ))}

                </div>
            </div>

        </div>
    );

}
