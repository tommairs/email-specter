"use client";

import {Chart} from 'primereact/chart';
import {Card} from 'primereact/card';

export default function ChartsTab({ chartData, chartOptions }) {

    return (
        <div className="row g-4">

            <div className="col-12">
                <div className="card border-0 shadow-sm">

                    <div className="card-header bg-white">
                        <h5 className="card-title mb-0">
                            Trends
                        </h5>
                    </div>

                    <div className="card-body" style={{ height: '400px' }}>
                        {chartData.line && (
                            <Chart
                                type="line"
                                data={chartData.line}
                                options={chartOptions.line}
                                style={{ height: '100%' }}
                            />
                        )}
                    </div>

                </div>
            </div>
            
            <div className="col-12 col-md-6">
                <div className="card border-0 shadow-sm h-100">
                    <div className="card-header bg-white">
                        <h5 className="card-title mb-0">Distribution</h5>
                    </div>
                    <div className="card-body" style={{ height: '300px' }}>
                        {chartData.pie && (
                            <Chart
                                type="pie"
                                data={chartData.pie}
                                options={chartOptions.pie}
                                style={{ height: '100%' }}
                            />
                        )}
                    </div>
                </div>
            </div>
            
            <div className="col-12 col-md-6">
                <div className="card border-0 shadow-sm h-100">
                    <div className="card-header bg-white">
                        <h5 className="card-title mb-0">Daily Breakdown</h5>
                    </div>
                    <div className="card-body" style={{ height: '300px' }}>
                        {chartData.bar && (
                            <Chart
                                type="bar"
                                data={chartData.bar}
                                options={chartOptions.bar}
                                style={{ height: '100%' }}
                            />
                        )}
                    </div>
                </div>
            </div>

        </div>
    );
}
