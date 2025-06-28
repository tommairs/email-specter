import {toast} from 'react-hot-toast';
import React from 'react';

export const ToastHelper = {

    loadingToast: (message) => {
        let loadingColor = 'text-dark';
        let backgroundColor = '#fff';
        let color = '#333';

        return toast(
            (t) => (
                <span>
          <div
              className={'spinner-border spinner-border-sm ' + loadingColor}
              role="status"
          >
            <span className="visually-hidden">Loading...</span>
          </div>
                    &nbsp; &nbsp; &nbsp;
                    {message}
        </span>
            ),
            {
                duration: Infinity,
                style: {
                    borderRadius: '10px',
                    background: backgroundColor,
                    color: color,
                },
            }
        );

    },

    successToast: (message) => {
        let backgroundColor = '#fff';
        let color = '#333';

        return toast.success(message, {
            style: {
                borderRadius: '10px',
                background: backgroundColor,
                color: color,
            },
        });
    },

    errorToast: (message) => {
        let backgroundColor = '#fff';
        let color = '#333';

        return toast.error(message, {
            style: {
                borderRadius: '10px',
                background: backgroundColor,
                color: color,
            },
        });
    },

    dismissToast: (toastId) => {
        toast.dismiss(toastId);
    },

    deletePrompt: (message, cancelLabel, onConfirm) => {
        return toast(
            (t) => (
                <div>

                    <p className="text-natural-dark-3 fs-sm">
                        {message}
                    </p>

                    <div className="d-flex justify-content-end gap-2 mt-4">

                        <button className={"btn btn-sm btn-outline-danger"} onClick={() => toast.dismiss(t.id)}>
                            X
                        </button>

                        <button
                            type="button"
                            className={"btn btn-sm btn-danger"}
                            onClick={() => {
                                toast.dismiss(t.id);
                                onConfirm();
                            }}>
                            {cancelLabel}
                        </button>

                    </div>

                </div>
            ),
            {
                duration: 10000,
                style: {
                    borderRadius: '10px',
                    background: '#fff',
                    color: '#333',
                    padding: '10px',
                    minWidth: '250px',
                },
            }
        );
    },

};