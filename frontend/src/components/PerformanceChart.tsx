'use client';

import React from 'react';
import { Card } from 'antd';
import {
  LineChart,
  Line,
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';
import { TrendingUpIcon, TrendingDownIcon, BarChart3Icon } from 'lucide-react';
import { BacktestResult } from '@/types';
import { formatCurrency } from '@/lib/utils';

interface PerformanceChartProps {
  result: BacktestResult;
}

// Format data for charts
const formatCumulativeData = (result: BacktestResult) => {
  return result.daily_returns.map(dailyReturn => ({
    date: dailyReturn.date,
    '累计收益 (%)': (dailyReturn.cumulative_return * 100).toFixed(2),
    '策略价值': dailyReturn.portfolio_value,
  }));
};

const formatDrawdownData = (result: BacktestResult) => {
  return result.daily_returns.map(dailyReturn => ({
    date: dailyReturn.date,
    '回撤 (%)': (-dailyReturn.drawdown * 100).toFixed(2),
  }));
};

const formatDailyReturnData = (result: BacktestResult) => {
  return result.daily_returns.map(dailyReturn => ({
    date: dailyReturn.date,
    '日收益率 (%)': (dailyReturn.daily_return * 100).toFixed(2),
  }));
};

export function PerformanceChart({ result }: PerformanceChartProps) {
  const cumulativeData = formatCumulativeData(result);
  const drawdownData = formatDrawdownData(result);
  const dailyReturnData = formatDailyReturnData(result);

  return (
    <div className="space-y-6">
      {/* Cumulative Returns Chart */}
      <Card 
        title="累计收益走势" 
        extra={<TrendingUpIcon className="h-5 w-5" />}
        className="w-full"
      >
        <div className="h-80">
          <ResponsiveContainer width="100%" height="100%">
            <LineChart
              data={cumulativeData}
              margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
            >
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="date" />
              <YAxis 
                yAxisId="left"
                tickFormatter={(value) => `${value}%`}
              />
              <YAxis 
                yAxisId="right" 
                orientation="right"
                tickFormatter={(value) => formatCurrency(Number(value))}
              />
              <Tooltip 
                formatter={(value, name) => {
                  if (name === '累计收益 (%)') {
                    return [`${value}%`, name];
                  }
                  return [formatCurrency(Number(value)), name];
                }}
              />
              <Legend />
              <Line
                yAxisId="left"
                type="monotone"
                dataKey="累计收益 (%)"
                stroke="#8884d8"
                activeDot={{ r: 8 }}
                name="累计收益 (%)"
              />
              <Line
                yAxisId="right"
                type="monotone"
                dataKey="策略价值"
                stroke="#82ca9d"
                name="策略价值"
              />
            </LineChart>
          </ResponsiveContainer>
        </div>
      </Card>

      {/* Drawdown Chart */}
      <Card 
        title="回撤分析" 
        extra={<TrendingDownIcon className="h-5 w-5" />}
        className="w-full"
      >
        <div className="h-80">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart
              data={drawdownData}
              margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
            >
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="date" />
              <YAxis 
                tickFormatter={(value) => `${value}%`}
              />
              <Tooltip 
                formatter={(value) => [`${value}%`, '回撤']}
              />
              <Bar 
                dataKey="回撤 (%)" 
                fill="#ff6b6b"
                name="回撤"
              />
            </BarChart>
          </ResponsiveContainer>
        </div>
      </Card>

      {/* Daily Returns Distribution */}
      <Card 
        title="日收益率分布" 
        extra={<BarChart3Icon className="h-5 w-5" />}
        className="w-full"
      >
        <div className="h-80">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart
              data={dailyReturnData}
              margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
            >
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="date" />
              <YAxis 
                tickFormatter={(value) => `${value}%`}
              />
              <Tooltip 
                formatter={(value) => [`${value}%`, '日收益率']}
              />
              <Bar 
                dataKey="日收益率 (%)" 
                name="日收益率"
              >
                {dailyReturnData.map((entry, index) => (
                  <CustomBar 
                    key={`cell-${index}`} 
                    fill={Number(entry['日收益率 (%)']) >= 0 ? '#00C49F' : '#FF8042'} 
                  />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </div>
      </Card>
    </div>
  );
}

// Custom component for conditional coloring
const CustomBar = ({ fill }: { fill: string }) => (
  <rect fill={fill} />
);