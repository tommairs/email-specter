"use client";

export default function SummaryCards({summaryStats}) {

    const stats = [
        {
            label: "Accepted",
            value: summaryStats.reception,
            icon: "fa fa-globe",
            bg: "bg-primary",
            iconColor: "text-white"
        },
        {
            label: "Deliveries",
            value: summaryStats.delivery,
            icon: "fa fa-check-circle",
            bg: "bg-success",
            iconColor: "text-white"
        },
        {
            label: "Bounces",
            value: summaryStats.bounce,
            icon: "fa fa-times-circle",
            bg: "bg-danger",
            iconColor: "text-white"
        },
        {
            label: "Transient Failures",
            value: summaryStats.transientFailure,
            icon: "fa fa-exclamation-triangle",
            bg: "bg-warning",
            iconColor: "text-dark"
        }
    ];

    return (
        <div className="row g-4 mb-4">
            {stats.map((stat, index) => (
                <div key={index} className="col-12 col-sm-6 col-lg-3">
                    <div className="card h-100 border-0 shadow-sm">
                        <div className="card-body p-3">
                            <div className="d-flex align-items-center">
                                <div className={`rounded-circle p-3 me-3 ${stat.bg}`}>
                                    <i className={`${stat.icon} fs-4 ${stat.iconColor}`}></i>
                                </div>
                                <div>
                                    <h6 className="card-title text-muted mb-1 small">{stat.label}</h6>
                                    <h4 className="mb-0 fw-bold">
                                        {stat.value ? stat.value.toLocaleString() : '0'}
                                    </h4>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            ))}
        </div>
    );

}
