import GlobalHelper from "@/helpers/GlobalHelper";

export default function ReportTable({ reports }) {


    return (
        <div className="table-responsive">
            <table className="table table-striped table-bordered">
                <thead>
                    <tr>
                        <th>Type</th>
                        <th>Count</th>
                    </tr>
                </thead>
                <tbody>
                    {reports.map((report, index) => (
                        <tr key={index}>
                            <td>{report.item}</td>
                            <td>{GlobalHelper.addCommasToNumber(report.count)}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );

}