"use client";

export default function EmptyState({onAddConnection}) {

    return (
        <div className="text-center py-5">

            <div className="mb-3">
                <i className="fas fa-plug fa-3x text-muted"></i>
            </div>

            <h5>No connections yet</h5>
            <p className="text-muted">Create your first connection to start receiving webhook data.</p>

            <button className="btn btn-primary" onClick={onAddConnection}>
                Add Your First Connection
            </button>

        </div>
    );

}
