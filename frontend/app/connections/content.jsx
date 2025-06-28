"use client";

import {useEffect, useState} from "react";
import RequestHelper from "@/helpers/RequestHelper";
import {ToastHelper} from "@/helpers/ToastHelper";
import {useRouter} from "next/navigation";

// Import components
import PageHeader from "./components/PageHeader";
import LoadingState from "./components/LoadingState";
import EmptyState from "./components/EmptyState";
import ConnectionCard from "./components/ConnectionCard";
import AddConnectionModal from "./components/AddConnectionModal";
import EditConnectionModal from "./components/EditConnectionModal";

export default function Content() {

    const router = useRouter();

    const [mtas, setMtas] = useState([]);
    const [loading, setLoading] = useState(true);
    const [showAddModal, setShowAddModal] = useState(false);
    const [showEditModal, setShowEditModal] = useState(false);
    const [selectedMta, setSelectedMta] = useState(null);
    const [mtaName, setMtaName] = useState("");

    useEffect(() => {
        loadMtas();
    }, []);

    const loadMtas = async () => {

        setLoading(true);

        try {

            const response = await RequestHelper.sendAuthenticatedGetRequest('/mta');

            if (response.data.success) {
                setMtas(response.data.mtas || []);
            } else {
                ToastHelper.errorToast(response.data.message || 'There was an error loading connections');
            }

        } catch (error) {

            ToastHelper.errorToast('There was an error loading connections');

        }

        setLoading(false);

    };

    const handleAddMta = async (name) => {

        try {

            const response = await RequestHelper.sendAuthenticatedPostRequest('/mta', {
                name: name
            });

            if (response.data.success) {

                ToastHelper.successToast('The connection was added successfully');

                setShowAddModal(false);
                setMtaName("");

                loadMtas();

            } else {

                ToastHelper.errorToast(response.data.message || 'We were unable to add the connection');

            }

        } catch (error) {
            ToastHelper.errorToast('We were unable to add the connection');
        }

    };

    const handleEditMta = async (mtaId, name) => {

        try {

            const response = await RequestHelper.sendAuthenticatedPatchRequest(`/mta/${mtaId}`, {
                name: name
            });

            if (response.data.success) {

                ToastHelper.successToast('The connection was updated successfully');

                setShowEditModal(false);
                setMtaName("");
                setSelectedMta(null);

                loadMtas();

            } else {
                ToastHelper.errorToast(response.data.message || 'Failed to update connection');
            }

        } catch (error) {
            ToastHelper.errorToast('Failed to update connection');
        }

    };

    const handleDeleteMta = async (mta) => {

        if (!confirm(`Are you sure you want to delete "${mta.name}"?`)) {
            return;
        }

        try {

            const response = await RequestHelper.sendAuthenticatedDeleteRequest(`/mta/${mta.id}`);

            if (response.data.success) {
                ToastHelper.successToast('Connection deleted successfully');
                loadMtas();
            } else {
                ToastHelper.errorToast(response.data.message || 'Failed to delete connection');
            }

        } catch (error) {
            ToastHelper.errorToast('Failed to delete connection');
        }

    };

    const handleRotateToken = async (mta) => {

        if (!confirm(`Are you sure you want to rotate the secret token for "${mta.name}"? This will invalidate the current webhook URL.`)) {
            return;
        }

        try {

            const response = await RequestHelper.sendAuthenticatedPostRequest(`/mta/${mta.id}/rotate-secret-token`, {});

            if (response.data.success) {
                ToastHelper.successToast('Secret token rotated successfully');
                loadMtas();
            } else {
                ToastHelper.errorToast(response.data.message || 'Failed to rotate secret token');
            }

        } catch (error) {
            ToastHelper.errorToast('Failed to rotate secret token');
        }

    };

    const openEditModal = (mta) => {
        setSelectedMta(mta);
        setMtaName(mta.name);
        setShowEditModal(true);
    };

    const closeAddModal = () => {
        setShowAddModal(false);
    };

    const closeEditModal = () => {
        setShowEditModal(false);
        setSelectedMta(null);
    };

    return (
        <div className="container">

            <PageHeader onAddConnection={() => setShowAddModal(true)}/>

            {loading ? (
                <LoadingState/>
            ) : mtas.length === 0 ? (
                <EmptyState onAddConnection={() => setShowAddModal(true)}/>
            ) : (
                <div className="row">
                    {mtas.map((mta) => (
                        <ConnectionCard
                            key={mta.id}
                            mta={mta}
                            onEdit={openEditModal}
                            onDelete={handleDeleteMta}
                            onRotateToken={handleRotateToken}
                        />
                    ))}
                </div>
            )}

            <AddConnectionModal
                show={showAddModal}
                onClose={closeAddModal}
                onSubmit={handleAddMta}
            />

            <EditConnectionModal
                show={showEditModal}
                mta={selectedMta}
                onClose={closeEditModal}
                onSubmit={handleEditMta}
            />

        </div>
    );

}
