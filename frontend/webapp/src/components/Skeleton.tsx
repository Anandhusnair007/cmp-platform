import React from 'react';

interface SkeletonProps {
    className?: string;
    variant?: 'text' | 'circular' | 'rectangular';
    width?: string | number;
    height?: string | number;
}

export default function Skeleton({
    className = '',
    variant = 'rectangular',
    width,
    height
}: SkeletonProps) {
    const baseClasses = "animate-pulse bg-slate-800/50 rounded";
    const variantClasses = {
        text: "h-4 w-full rounded",
        circular: "rounded-full",
        rectangular: "rounded-lg",
    };

    const style = {
        width,
        height,
    };

    return (
        <div
            className={`${baseClasses} ${variantClasses[variant]} ${className}`}
            style={style}
        />
    );
}

export function DashboardSkeleton() {
    return (
        <div className="space-y-6">
            {/* Header Skeleton */}
            <div className="flex justify-between items-center">
                <div className="space-y-2">
                    <Skeleton width={200} height={32} />
                    <Skeleton width={300} height={20} />
                </div>
                <Skeleton width={120} height={40} />
            </div>

            {/* Stats Grid Skeleton */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                {[1, 2, 3, 4].map((i) => (
                    <div key={i} className="glass p-6 rounded-xl border border-slate-800">
                        <div className="flex justify-between items-start">
                            <div className="space-y-2">
                                <Skeleton width={80} height={16} />
                                <Skeleton width={40} height={32} />
                            </div>
                            <Skeleton variant="circular" width={48} height={48} />
                        </div>
                    </div>
                ))}
            </div>

            {/* Main Content Skeleton */}
            <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                <div className="lg:col-span-2 glass p-6 rounded-xl border border-slate-800 space-y-4">
                    <div className="flex justify-between">
                        <Skeleton width={150} height={24} />
                        <Skeleton width={80} height={24} />
                    </div>
                    {[1, 2, 3, 4, 5].map((i) => (
                        <div key={i} className="flex justify-between items-center py-2">
                            <div className="flex items-center gap-3">
                                <Skeleton variant="circular" width={32} height={32} />
                                <div className="space-y-1">
                                    <Skeleton width={120} height={16} />
                                    <Skeleton width={80} height={12} />
                                </div>
                            </div>
                            <Skeleton width={60} height={24} />
                        </div>
                    ))}
                </div>
                <div className="glass p-6 rounded-xl border border-slate-800 space-y-4">
                    <Skeleton width={120} height={24} />
                    {[1, 2, 3].map((i) => (
                        <Skeleton key={i} height={60} />
                    ))}
                </div>
            </div>
        </div>
    );
}
