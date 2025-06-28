"use client";

export default function LoadingState() {

    return (
        <div className="text-center py-4">
            <div className="spinner-border" role="status">
                <span className="visually-hidden">Loading...</span>
            </div>
        </div>
    );

}
