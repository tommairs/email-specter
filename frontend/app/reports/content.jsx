"use client";

import {useEffect, useState} from "react";
import RequestHelper from "@/helpers/RequestHelper";
import GlobalHelper from "@/helpers/GlobalHelper";
import {Calendar} from "primereact/calendar";
import {Dropdown} from "primereact/dropdown";
import {InputText} from "primereact/inputtext";
import {ToastHelper} from "@/helpers/ToastHelper";
import ReportTable from "@/app/reports/components/ReportTable";

const BOUNCE_CLASSIFICATIONS = ["recipient", "block", "spam", "blacklist", "dmarc", "auth", "envelope", "message", "rate", "capacity", "greylist", "network", "protocol", "config", "other"];
const STATUS_OPTIONS = ["Reception", "TransientFailure", "Bounce", "Delivery"];
const GROUP_BY_OPTIONS = ["source_ip", "source_domain", "destination_domain", "destination_service", "kumo_mta_classification", "email_specter_classification"];

const getDateOffset = (days) => {

    const date = new Date();
    date.setDate(date.getDate() - days);

    return date;

};

export default function Content() {

    const [filters, setFilters] = useState({
        from: getDateOffset(30),
        to: new Date(),
        source_ip: '',
        source_domain: '',
        destination_domain: '',
        destination_service: '',
        kumo_mta_classification: '',
        email_specter_classification: '',
        event_type: '',
        group_by: 'destination_domain',
        max_results: 1000,
    });

    const [entities, setEntities] = useState({
        destination_domains: [],
        destination_services: [],
        source_domains: [],
        source_ips: [],
        kumo_mta_classifications: []
    });

    const [loading, setLoading] = useState(false);
    const [report, setReport] = useState(null);

    useEffect(() => {

        const loadEntities = async () => {
            const data = await GlobalHelper.fetchTopEntities();
            if (data) setEntities(data);
        };

        loadEntities();

    }, []);

    const updateFilter = (field, value) => {
        setFilters(prev => ({...prev, [field]: value}));
    };

    const generateReport = async () => {

        setLoading(true);

        const response = await RequestHelper.sendAuthenticatedPostRequest('/reports/generate', {
            from: filters.from?.toISOString().split('T')[0] || "",
            to: filters.to?.toISOString().split('T')[0] || "",
            source_ip: filters.source_ip || "",
            source_domain: filters.source_domain || "",
            destination_service: filters.destination_service || "",
            destination_domain: filters.destination_domain || "",
            event_type: filters.event_type || "",
            email_specter_classification: filters.email_specter_classification || "",
            kumo_mta_classifications: filters.kumo_mta_classification || "",
            group_by: filters.group_by || "",
            max_results: filters.max_results || 1000,
        });

        if (response.data.success) {
            setReport(response.data.data);
        } else {
            ToastHelper.errorToast(response.data.message);
        }

        setLoading(false);

    };

    const FilterField = ({label, children, colClass = "col-md-3"}) => (
        <div className={colClass}>
            <label className="form-label">{label}</label>
            {children}
        </div>
    );

    return (
        <div className="container">

            <div className="flex align-items-center gap-3 mb-4">
                <div>
                    <h1 className="text-3xl font-bold text-900 m-0">Reports</h1>
                    <p className="text-600 m-0 mt-1">View delivery counts by domain, IP, classification, and more.</p>
                </div>
            </div>

            <div className="card mb-4">
                <div className="card-body">

                    <div className="row g-3">

                        <FilterField label="From Date">
                            <Calendar
                                value={filters.from}
                                onChange={(e) => updateFilter("from", e.value)}
                                dateFormat="yy-mm-dd"
                                className="w-100"
                                showIcon={false}
                            />
                        </FilterField>

                        <FilterField label="To Date">
                            <Calendar
                                value={filters.to}
                                onChange={(e) => updateFilter("to", e.value)}
                                dateFormat="yy-mm-dd"
                                className="w-100"
                                showIcon={false}
                            />
                        </FilterField>

                        <FilterField label="Source IP">
                            <Dropdown
                                value={filters.source_ip}
                                options={entities.source_ips}
                                onChange={(e) => updateFilter("source_ip", e.value)}
                                placeholder="All Source IPs"
                                className="w-100"
                                filter
                                showClear
                            />
                        </FilterField>

                        <FilterField label="Source Domain">
                            <InputText
                                value={filters.source_domain}
                                onChange={(e) => updateFilter("source_domain", e.target.value)}
                                className="w-100"
                            />
                        </FilterField>

                    </div>

                    <div className="row mt-4">

                        <FilterField label="Destination Domain">
                            <InputText
                                value={filters.destination_domain}
                                onChange={(e) => updateFilter("destination_domain", e.target.value)}
                                className="w-100"
                            />
                        </FilterField>

                        <FilterField label="Destination Service">
                            <Dropdown
                                value={filters.destination_service}
                                options={entities.destination_services}
                                onChange={(e) => updateFilter("destination_service", e.value)}
                                placeholder="All Services"
                                className="w-100"
                                filter
                                showClear
                            />
                        </FilterField>

                        <FilterField label="Email Specter Classification">
                            <Dropdown
                                value={filters.email_specter_classification}
                                options={BOUNCE_CLASSIFICATIONS}
                                onChange={(e) => updateFilter("email_specter_classification", e.value)}
                                placeholder="All Classifications"
                                className="w-100"
                                filter
                                showClear
                            />
                        </FilterField>

                        <FilterField label="KumoMTA Classification">
                            <Dropdown
                                value={filters.kumo_mta_classification}
                                options={entities.kumo_mta_classifications}
                                onChange={(e) => updateFilter("kumo_mta_classification", e.value)}
                                placeholder="All Classifications"
                                className="w-100"
                                filter
                                showClear
                            />
                        </FilterField>

                    </div>

                    <div className="row mt-4">

                        <FilterField label="Status">
                            <Dropdown
                                value={filters.event_type}
                                options={STATUS_OPTIONS}
                                onChange={(e) => updateFilter("event_type", e.value)}
                                placeholder="All Statuses"
                                className="w-100"
                                filter
                                showClear
                            />
                        </FilterField>

                        <FilterField label="Group By">
                            <Dropdown
                                value={filters.group_by}
                                options={GROUP_BY_OPTIONS}
                                onChange={(e) => updateFilter("group_by", e.value)}
                                placeholder="All Groups"
                                className="w-100"
                                filter
                                showClear
                            />
                        </FilterField>

                        <FilterField label="Max Results">
                            <InputText
                                type="number"
                                value={filters.max_results}
                                onChange={(e) => updateFilter("max_results", parseInt(e.target.value) || 1000)}
                                className="w-100"
                            />
                        </FilterField>

                    </div>

                    <button type="button" className="btn btn-lg btn-primary mt-4" onClick={generateReport} disabled={loading}>
                        {loading ? (
                            <span className="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                        ) : (
                            <i className="fa fa-file-alt me-2"></i>
                        )}
                        Generate Report
                    </button>

                </div>
            </div>

            {report && <ReportTable reports={report}/>}

        </div>
    );

}