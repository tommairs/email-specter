"use client";

import {useEffect, useState} from "react";
import RequestHelper from "@/helpers/RequestHelper";
import {ToastHelper} from "@/helpers/ToastHelper";
import {useRouter} from "next/navigation";
import Loading from "@/components/Loading";
import SummaryCards from "./components/SummaryCards";
import ChartsTab from "./components/ChartsTab";
import DataTableTab from "./components/DataTableTab";
import AnalyticsTab from "./components/AnalyticsTab";

export default function Content() {

    const [loading, setLoading] = useState(true);

    const [statistics, setStatistics] = useState({});
    const [chartData, setChartData] = useState({});
    const [chartOptions, setChartOptions] = useState({});
    const [activeTab, setActiveTab] = useState('charts');

    const processDataForCharts = (data) => {

        if (!data.data || data.data.length === 0) return;

        const dates = data.data.map(item => item.date);
        const eventTypes = ['Delivery', 'Reception', 'Bounce', 'TransientFailure'];

        const datasets = eventTypes.map((eventType, index) => {

            const colors = ['#42A5F5', '#66BB6A', '#EF5350', '#FFA726'];

            return {
                label: eventType,
                data: data.data.map(dayData => {
                    const event = dayData.events.find(e => e.event_type === eventType);
                    return event ? event.count : 0;
                }),
                borderColor: colors[index],
                backgroundColor: colors[index] + '20',
                tension: 0.4
            };

        });

        const lineChartData = {
            labels: dates,
            datasets: datasets
        };

        const totals = eventTypes.map(eventType => {

            return data.data.reduce((sum, dayData) => {
                const event = dayData.events.find(e => e.event_type === eventType);
                return sum + (event ? event.count : 0);
            }, 0);

        });

        const pieChartData = {
            labels: eventTypes,
            datasets: [{
                data: totals,
                backgroundColor: ['#42A5F5', '#66BB6A', '#EF5350', '#FFA726'],
                hoverBackgroundColor: ['#64B5F6', '#81C784', '#E57373', '#FFB74D']
            }]
        };

        const barChartData = {
            labels: dates,
            datasets: [{
                label: 'Total Events',
                data: data.data.map(dayData =>
                    dayData.events.reduce((sum, event) => sum + event.count, 0)
                ),
                backgroundColor: '#42A5F5',
                borderColor: '#1E88E5',
                borderWidth: 1
            }]
        };

        setChartData({
            line: lineChartData,
            pie: pieChartData,
            bar: barChartData
        });
    };

    const loadStatistics = async () => {

        const response = await RequestHelper.sendAuthenticatedGetRequest('/reports/aggregated-data');
        const data = response.data;

        if (data.success) {
            setStatistics(data);
            processDataForCharts(data);
        } else {
            ToastHelper.errorToast(response.message);
        }

        setLoading(false);

    };

    const calculateSummaryStats = () => {

        if (!statistics.data) return {};

        const totals = {
            delivery: 0,
            reception: 0,
            bounce: 0,
            transientFailure: 0,
            total: 0
        };

        statistics.data.forEach(dayData => {

            dayData.events.forEach(event => {

                totals.total += event.count;

                switch (event.event_type) {
                    case 'Delivery':
                        totals.delivery += event.count;
                        break;
                    case 'Reception':
                        totals.reception += event.count;
                        break;
                    case 'Bounce':
                        totals.bounce += event.count;
                        break;
                    case 'TransientFailure':
                        totals.transientFailure += event.count;
                        break;
                }

            });

        });

        return totals;

    };

    useEffect(() => {

        const options = {
            responsive: true,
            maintainAspectRatio: false,
            plugins: {
                legend: {
                    position: 'top',
                },
                title: {
                    display: false
                }
            },
            scales: {
                y: {
                    beginAtZero: true
                }
            }
        };

        setChartOptions({
            line: options,
            bar: options,
            pie: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: {
                        position: 'right',
                    }
                }
            }
        });

        loadStatistics();

    }, []);

    const summaryStats = calculateSummaryStats();

    const tabs = [
        {id: 'charts', name: 'Charts', icon: 'fas fa-chart-line'},
        {id: 'analytics', name: 'Analytics', icon: 'fa fa-list-alt'},
        {id: 'data', name: 'Daily Overview', icon: 'fas fa-table'}
    ];

    return (
        <div className="container py-4">
            <div className="d-flex justify-content-between align-items-center mb-4">

                <div>
                    <h2 className="mb-1 fw-bold">Statistics</h2>
                    <p className="text-muted mb-0">Monitor total events, delivery success, bounce rates, and transient failures at a glance.</p>
                </div>

                <button className="btn btn-outline-dark btn-sm" onClick={loadStatistics} disabled={loading}>
                    <i className="fas fa-sync me-2"></i>
                    Refresh
                </button>

            </div>

            {
                loading && <Loading />
            }

            {!loading && statistics.data && (

                <>

                    <SummaryCards summaryStats={summaryStats}/>

                    <div className="card border-0 shadow-sm">

                        <div className="card-header bg-white border-bottom-0">
                            <ul className="nav nav-tabs card-header-tabs" role="tablist">

                                {tabs.map((tab) => (
                                    <li key={tab.id} className="nav-item" role="presentation">
                                        <button className={`nav-link ${activeTab === tab.id ? 'active' : ''}`} type="button" role="tab" onClick={() => setActiveTab(tab.id)}>
                                            <i className={`${tab.icon} me-2`}></i>
                                            {tab.name}
                                        </button>
                                    </li>
                                ))}

                            </ul>

                        </div>

                        <div className="card-body">

                            <div className="tab-content">

                                <div className={`tab-pane fade ${activeTab === 'charts' ? 'show active' : ''}`}>
                                    <ChartsTab chartData={chartData} chartOptions={chartOptions}/>
                                </div>

                                <div className={`tab-pane fade ${activeTab === 'analytics' ? 'show active' : ''}`}>
                                    <AnalyticsTab summaryStats={summaryStats} statistics={statistics}/>
                                </div>

                                <div className={`tab-pane fade ${activeTab === 'data' ? 'show active' : ''}`}>
                                    <DataTableTab statistics={statistics}/>
                                </div>

                            </div>

                        </div>

                    </div>

                </>
            )}

        </div>
    );

}