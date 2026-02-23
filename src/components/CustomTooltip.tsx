import React from 'react';
import type { TooltipProps } from '../types';

const CustomTooltip: React.FC<TooltipProps> = ({ active, payload, label }) => {
    if (active && payload && payload.length) {
        return (
            <div style={{ 
                background: 'rgba(255, 255, 255, 0.95)', 
                padding: '12px 16px', 
                borderRadius: 12, 
                boxShadow: '0 4px 20px rgba(0,0,0,0.1)',
                border: '1px solid rgba(0,0,0,0.05)'
            }}>
                <div style={{ fontWeight: 600, marginBottom: 8, color: '#333' }}>{label}</div>
                {payload.map((entry, index) => (
                    <div key={index} style={{ display: 'flex', alignItems: 'center', gap: 8, marginBottom: 4 }}>
                        <div style={{ 
                            width: 8, 
                            height: 8, 
                            borderRadius: '50%', 
                            background: entry.color 
                        }} />
                        <span style={{ color: '#666', fontSize: 13 }}>{entry.name}:</span>
                        <span style={{ fontWeight: 600, color: entry.color }}>{entry.value}</span>
                    </div>
                ))}
            </div>
        );
    }
    return null;
};

export default CustomTooltip;
