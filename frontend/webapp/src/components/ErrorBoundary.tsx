import React, { Component, ErrorInfo, ReactNode } from 'react';
import { ExclamationTriangleIcon, ArrowPathIcon } from '@heroicons/react/24/outline';

interface Props {
    children: ReactNode;
}

interface State {
    hasError: boolean;
    error: Error | null;
}

export default class ErrorBoundary extends Component<Props, State> {
    public state: State = {
        hasError: false,
        error: null,
    };

    public static getDerivedStateFromError(error: Error): State {
        return { hasError: true, error };
    }

    public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
        console.error('Uncaught error:', error, errorInfo);
    }

    public render() {
        if (this.state.hasError) {
            return (
                <div className="min-h-screen flex items-center justify-center bg-slate-950 p-4">
                    <div className="max-w-md w-full glass p-8 rounded-xl border border-red-900/50 text-center">
                        <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-red-900/20 mb-6">
                            <ExclamationTriangleIcon className="h-10 w-10 text-red-500" />
                        </div>
                        <h1 className="text-2xl font-bold text-white mb-2">Something went wrong</h1>
                        <p className="text-slate-400 mb-6">
                            The application encountered an unexpected error. Our team has been notified.
                        </p>

                        {this.state.error && (
                            <div className="bg-slate-900/50 p-4 rounded-lg border border-slate-800 mb-6 text-left overflow-auto max-h-32">
                                <p className="font-mono text-xs text-red-400">
                                    {this.state.error.toString()}
                                </p>
                            </div>
                        )}

                        <button
                            onClick={() => window.location.reload()}
                            className="btn-primary w-full flex items-center justify-center gap-2"
                        >
                            <ArrowPathIcon className="h-5 w-5" />
                            Reload Application
                        </button>
                    </div>
                </div>
            );
        }

        return this.props.children;
    }
}
