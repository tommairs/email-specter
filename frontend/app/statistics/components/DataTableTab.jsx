"use client";

import {DataTable} from 'primereact/datatable';
import {Column} from 'primereact/column';

export default function DataTableTab({ statistics }) {

    const eventTypes = ["Reception", "Delivery", "Bounce", "TransientFailure"];
    
    return (
        <div className="card border-0 shadow-sm">

            <div className="card-header bg-white">
                <h5 className="card-title mb-0">Daily Overview</h5>
            </div>

            <div className="card-body p-0">

                <DataTable
                    value={statistics.data}
                    paginator
                    rows={10}
                    dataKey="date"
                    stripedRows
                    className="p-datatable-sm"
                    emptyMessage="There are no statistics available for the selected period."
                >

                    <Column field="date" header="Date" sortable />

                    {eventTypes.map(type => (
                        <Column
                            key={type}
                            header={type}
                            body={(rowData) => {
                                const event = rowData.events.find(e => e.event_type === type);
                                return event ? event.count.toLocaleString() : '0';
                            }}
                        />
                    ))}

                    <Column
                        header="Total"
                        body={(rowData) => {
                            const total = rowData.events.reduce((sum, event) => {
                                if (event.event_type !== "Reception") {
                                    return sum + event.count;
                                }
                                return sum;
                            }, 0);
                            return total.toLocaleString();
                        }}
                    />

                </DataTable>

            </div>

        </div>
    );

}
