"use client";

export default function AnalyticsTab({ summaryStats, statistics }) {

    const calculateMetrics = () => {

        if (!summaryStats || summaryStats.total === 0) {

            return {
                deliveryRate: 0,
                bounceRate: 0,
            };

        }

        const deliveryRate = ((summaryStats.delivery / summaryStats.reception) * 100).toFixed(2);
        const bounceRate = ((summaryStats.bounce / summaryStats.reception) * 100).toFixed(2);

        return { deliveryRate, bounceRate };

    };

    const metrics = calculateMetrics();

    const getProgressBarClass = (value) => {

        if (value >= 90) return 'bg-success';
        if (value >= 70) return 'bg-warning';

        return 'bg-danger';

    };

    const getTrendData = () => {

        if (!statistics.data || statistics.data.length < 2) return null;

        const recent = statistics.data.slice(-2);

        let previousTotal;
        let currentTotal;

        previousTotal = recent[0].events.reduce((sum, event) => {
            return sum + (event.event_type === 'Delivery' ? event.count : 0);
        }, 0);

        currentTotal = recent[1].events.reduce((sum, event) => {
            return sum + (event.event_type === 'Delivery' ? event.count : 0);
        }, 0);

        const change = ((currentTotal - previousTotal) / previousTotal * 100).toFixed(1);

        return {
            previous: previousTotal,
            current: currentTotal,
            change: change,
            isPositive: change > 0
        };

    };

    const trendData = getTrendData();

    return (
        <div className="row g-4">

            <div className="col-12">

                <div className="card border-0 shadow-sm">

                    <div className="card-header bg-white">
                        <h5 className="card-title mb-0">
                            Performance Metrics
                        </h5>
                    </div>

                    <div className="card-body">

                        <div className="row g-4">

                            <div className="col-12 col-md-6">
                                <div className="text-center">
                                    <h6 className="text-muted mb-2">Delivery Rate</h6>
                                    <div className="progress mb-2" style={{ height: '8px' }}>
                                        <div
                                            className={`progress-bar ${getProgressBarClass(metrics.deliveryRate)}`}
                                            style={{ width: `${metrics.deliveryRate}%` }}
                                        ></div>
                                    </div>
                                    <h4 className="mb-0 text-success">{metrics.deliveryRate}%</h4>
                                </div>
                            </div>

                            <div className="col-12 col-md-6">
                                <div className="text-center">
                                <h6 className="text-muted mb-2">Bounce Rate</h6>
                                    <div className="progress mb-2" style={{ height: '8px' }}>
                                        <div
                                            className="progress-bar bg-danger"
                                            style={{ width: `${metrics.bounceRate}%` }}
                                        ></div>
                                    </div>
                                    <h4 className="mb-0 text-danger">{metrics.bounceRate}%</h4>
                                </div>
                            </div>

                        </div>

                    </div>

                </div>

            </div>

            {trendData && (
                <div className="col-12 col-md-6">
                    <div className="card border-0 shadow-sm">

                        <div className="card-header bg-white">
                            <h5 className="card-title mb-0">Trend</h5>
                        </div>

                        <div className="card-body">

                            <div className="d-flex justify-content-between align-items-center mb-3">
                                <span className="text-muted">Previous Day</span>
                                <span className="fw-bold">{trendData.previous.toLocaleString()}</span>
                            </div>

                            <div className="d-flex justify-content-between align-items-center mb-3">
                                <span className="text-muted">Latest Day</span>
                                <span className="fw-bold">{trendData.current.toLocaleString()}</span>
                            </div>

                            <div className="d-flex justify-content-between align-items-center">
                                <span className="text-muted">Change</span>
                                <span className={`fw-bold ${trendData.isPositive ? 'text-success' : 'text-danger'}`}>
                                    <i className={`bi ${trendData.isPositive ? 'bi-arrow-up' : 'bi-arrow-down'} me-1`}></i>
                                    {Math.abs(trendData.change)}%
                                </span>
                            </div>

                        </div>

                    </div>
                </div>
            )}

            <div className="col-12 col-md-6">
                <div className="card border-0 shadow-sm">

                    <div className="card-header bg-white">
                        <h5 className="card-title mb-0">Event Breakdown</h5>
                    </div>

                    <div className="card-body">

                        <div className="list-group list-group-flush">

                            <div className="list-group-item d-flex justify-content-between align-items-center">
                                <span><i className="bi bi-check-circle text-success me-2"></i>Delivered</span>
                                <span className="badge bg-success rounded-pill">{summaryStats.delivery?.toLocaleString() || '0'}</span>
                            </div>

                            <div className="list-group-item d-flex justify-content-between align-items-center">
                                <span><i className="bi bi-inbox text-primary me-2"></i>Accepted</span>
                                <span className="badge bg-primary rounded-pill">{summaryStats.reception?.toLocaleString() || '0'}</span>
                            </div>

                            <div className="list-group-item d-flex justify-content-between align-items-center">
                                <span><i className="bi bi-x-circle text-danger me-2"></i>Bounces</span>
                                <span className="badge bg-danger rounded-pill">{summaryStats.bounce?.toLocaleString() || '0'}</span>
                            </div>

                            <div className="list-group-item d-flex justify-content-between align-items-center">
                                <span><i className="bi bi-exclamation-triangle text-warning me-2"></i>Transient Failures</span>
                                <span className="badge bg-warning rounded-pill">{summaryStats.transientFailure?.toLocaleString() || '0'}</span>
                            </div>

                        </div>

                    </div>

                </div>
            </div>
        </div>
    );

}
