"use client";

export default function PageHeader({onAddConnection}) {

    return (
        <>

            <div className="d-flex justify-content-between align-items-center mb-4">

                <h1 className="text-2xl font-bold">Connections</h1>

                <button className="btn btn-primary" onClick={onAddConnection}>
                    <i className="fas fa-plus me-2"></i>
                    Add Connection
                </button>

            </div>

            <p className="mb-4 text-muted">
                Create webhook connections to receive real-time updates from your KumoMTA instances.
                <br/>
                Each connection provides a unique webhook URL that you can configure in your KumoMTA instance to send events directly to this service.
            </p>

        </>
    );

}
