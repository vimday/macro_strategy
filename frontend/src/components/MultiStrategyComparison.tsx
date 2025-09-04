'use client';

import React from 'react';
import { Card, Table, Alert, Spin } from 'antd';
import { 
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer,
  LineChart, Line
} from 'recharts';
import { 
  TrendingUpIcon, 
  TrendingDownIcon, 
  BarChart3Icon,
  ActivityIcon,
  TargetIcon
} from 'lucide-react';
import { MultiStrategyBacktestResult } from '@/types';
import { 
  formatPercentage, 
  formatCurrency,
  getSharpeRatioColor,
  getValueColor
} from '@/lib/utils';

interface MultiStrategyComparisonProps {
  result: MultiStrategyBacktestResult;
  loading?: boolean;
  error?: string | null;
}

export function MultiStrategyComparison({ result, loading = false, error = null }: MultiStrategyComparisonProps) {
  if (loading) {
    return (
      <Card className="mt-8">
        <div className="flex items-center justify-center py-12">
          <Spin size="large" />
          <span className="ml-4 text-lg">正在执行多策略回测对比...</span>
        </div>
      </Card>
    );
  }

  if (error) {
    return (
      <Card className="mt-8">
        <Alert
          message="多策略回测对比失败"
          description={error}
          type="error"
          showIcon
        />
      </Card>
    );
  }

  if (!result || !result.results || result.results.length === 0) {
    return (
      <Card className="mt-8">
        <Alert
          message="无回测结果"
          description="未找到多策略回测结果数据"
          type="info"
          showIcon
        />
      </Card>
    );
  }

  // Prepare data for charts
  const strategyNames = result.results.map(r => r.request.strategy.description || r.request.strategy.type);
  
  // Prepare metrics comparison data
  const metricsData = [
    {
      metric: '总收益率',
      ...(Object.fromEntries(result.results.map((r, i) => [strategyNames[i], r.performance_metrics.total_return])) as Record<string, number>)
    },
    {
      metric: '年化收益率',
      ...(Object.fromEntries(result.results.map((r, i) => [strategyNames[i], r.performance_metrics.annualized_return])) as Record<string, number>)
    },
    {
      metric: '最大回撤',
      ...(Object.fromEntries(result.results.map((r, i) => [strategyNames[i], -r.performance_metrics.max_drawdown])) as Record<string, number>)
    },
    {
      metric: '夏普比率',
      ...(Object.fromEntries(result.results.map((r, i) => [strategyNames[i], r.performance_metrics.sharpe_ratio])) as Record<string, number>)
    },
    {
      metric: '胜率',
      ...(Object.fromEntries(result.results.map((r, i) => [strategyNames[i], r.performance_metrics.win_rate])) as Record<string, number>)
    }
  ];

  // Prepare detailed metrics table
  const columns = [
    {
      title: '指标',
      dataIndex: 'metric',
      key: 'metric',
    },
    ...strategyNames.map((name, index) => ({
      title: name,
      dataIndex: `strategy_${index}`,
      key: `strategy_${index}`,
      render: (value: number) => {
        if (typeof value === 'number') {
          if (value >= 1 || value <= -1) {
            return formatCurrency(value);
          } else {
            return formatPercentage(value);
          }
        }
        return value;
      }
    }))
  ];

  const metricsTableData = [
    {
      key: 'total_return',
      metric: '总收益率',
      ...(Object.fromEntries(result.results.map((r, i) => [`strategy_${i}`, r.performance_metrics.total_return])) as Record<string, number>)
    },
    {
      key: 'annualized_return',
      metric: '年化收益率',
      ...(Object.fromEntries(result.results.map((r, i) => [`strategy_${i}`, r.performance_metrics.annualized_return])) as Record<string, number>)
    },
    {
      key: 'max_drawdown',
      metric: '最大回撤',
      ...(Object.fromEntries(result.results.map((r, i) => [`strategy_${i}`, -r.performance_metrics.max_drawdown])) as Record<string, number>)
    },
    {
      key: 'sharpe_ratio',
      metric: '夏普比率',
      ...(Object.fromEntries(result.results.map((r, i) => [`strategy_${i}`, r.performance_metrics.sharpe_ratio])) as Record<string, number>)
    },
    {
      key: 'sortino_ratio',
      metric: '索提诺比率',
      ...(Object.fromEntries(result.results.map((r, i) => [`strategy_${i}`, r.performance_metrics.sortino_ratio])) as Record<string, number>)
    },
    {
      key: 'win_rate',
      metric: '胜率',
      ...(Object.fromEntries(result.results.map((r, i) => [`strategy_${i}`, r.performance_metrics.win_rate])) as Record<string, number>)
    },
    {
      key: 'profit_factor',
      metric: '盈亏比',
      ...(Object.fromEntries(result.results.map((r, i) => [`strategy_${i}`, r.performance_metrics.profit_factor])) as Record<string, number>)
    },
    {
      key: 'total_trades',
      metric: '总交易次数',
      ...(Object.fromEntries(result.results.map((r, i) => [`strategy_${i}`, r.performance_metrics.total_trades])) as Record<string, number>)
    }
  ];

  // Prepare cumulative returns data for line chart
  const maxLength = Math.max(...result.results.map(r => r.daily_returns.length));
  const cumulativeReturnsData: Record<string, string | number>[] = [];
  
  for (let i = 0; i < maxLength; i++) {
    const dataPoint: Record<string, string | number> = { date: '' };
    
    result.results.forEach((resultItem, index) => {
      if (i < resultItem.daily_returns.length) {
        const dailyReturn = resultItem.daily_returns[i];
        dataPoint.date = dailyReturn.date;
        dataPoint[strategyNames[index]] = dailyReturn.cumulative_return * 100;
      }
    });
    
    if (dataPoint.date) {
      cumulativeReturnsData.push(dataPoint);
    }
  }

  // Prepare drawdown data
  const drawdownData: Record<string, string | number>[] = [];
  
  for (let i = 0; i < maxLength; i++) {
    const dataPoint: Record<string, string | number> = { date: '' };
    
    result.results.forEach((resultItem, index) => {
      if (i < resultItem.daily_returns.length) {
        const dailyReturn = resultItem.daily_returns[i];
        dataPoint.date = dailyReturn.date;
        dataPoint[strategyNames[index]] = -dailyReturn.drawdown * 100;
      }
    });
    
    if (dataPoint.date) {
      drawdownData.push(dataPoint);
    }
  }

  return (
    <div className="mt-8 space-y-6">
      {/* Header */}
      <Card>
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-2xl font-bold text-gray-900">多策略回测对比</h2>
            <p className="text-gray-600 mt-1">
              资产: {result.request.asset_id} | 期间: {result.request.start_date} 至 {result.request.end_date}
            </p>
          </div>
          <div className="text-right">
            <div className="text-sm text-gray-500">对比ID</div>
            <div className="text-sm font-mono text-gray-700">{result.id}</div>
          </div>
        </div>
      </Card>

      {/* Performance Comparison Chart */}
      <Card title="收益对比" extra={<BarChart3Icon className="h-5 w-5" />}>
        <div className="h-80">
          <ResponsiveContainer width="100%" height="100%">
            <BarChart
              data={metricsData}
              margin={{ top: 20, right: 30, left: 20, bottom: 60 }}
            >
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="metric" angle={-45} textAnchor="end" height={60} />
              <YAxis 
                tickFormatter={(value) => `${(value * 100).toFixed(1)}%`}
              />
              <Tooltip 
                formatter={(value) => [`${(Number(value) * 100).toFixed(2)}%`, '']}
                labelFormatter={(value) => `指标: ${value}`}
              />
              <Legend />
              {strategyNames.map((name, index) => (
                <Bar 
                  key={name} 
                  dataKey={name} 
                  fill={`hsl(${index * 60}, 70%, 50%)`} 
                  name={name}
                />
              ))}
            </BarChart>
          </ResponsiveContainer>
        </div>
      </Card>

      {/* Cumulative Returns Chart */}
      <Card title="累计收益走势" extra={<TrendingUpIcon className="h-5 w-5" />}>
        <div className="h-80">
          <ResponsiveContainer width="100%" height="100%">
            <LineChart
              data={cumulativeReturnsData}
              margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
            >
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="date" />
              <YAxis 
                tickFormatter={(value) => `${value.toFixed(1)}%`}
              />
              <Tooltip 
                formatter={(value) => [`${Number(value).toFixed(2)}%`, '']}
                labelFormatter={(value) => `日期: ${value}`}
              />
              <Legend />
              {strategyNames.map((name, index) => (
                <Line 
                  key={name} 
                  type="monotone" 
                  dataKey={name} 
                  stroke={`hsl(${index * 60}, 70%, 50%)`} 
                  activeDot={{ r: 8 }} 
                  name={name}
                />
              ))}
            </LineChart>
          </ResponsiveContainer>
        </div>
      </Card>

      {/* Drawdown Chart */}
      <Card title="回撤走势" extra={<TrendingDownIcon className="h-5 w-5" />}>
        <div className="h-80">
          <ResponsiveContainer width="100%" height="100%">
            <LineChart
              data={drawdownData}
              margin={{ top: 5, right: 30, left: 20, bottom: 5 }}
            >
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="date" />
              <YAxis 
                tickFormatter={(value) => `${value.toFixed(1)}%`}
              />
              <Tooltip 
                formatter={(value) => [`${Number(value).toFixed(2)}%`, '']}
                labelFormatter={(value) => `日期: ${value}`}
              />
              <Legend />
              {strategyNames.map((name, index) => (
                <Line 
                  key={name} 
                  type="monotone" 
                  dataKey={name} 
                  stroke={`hsl(${index * 60}, 70%, 50%)`} 
                  activeDot={{ r: 8 }} 
                  name={name}
                />
              ))}
            </LineChart>
          </ResponsiveContainer>
        </div>
      </Card>

      {/* Detailed Metrics Comparison */}
      <Card title="详细指标对比" extra={<ActivityIcon className="h-5 w-5" />}>
        <Table 
          dataSource={metricsTableData} 
          columns={columns} 
          pagination={false} 
          size="small"
          scroll={{ x: true }}
        />
      </Card>

      {/* Strategy Rankings */}
      {result.comparison && (
        <Card title="策略排名" extra={<TargetIcon className="h-5 w-5" />}>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            <Card size="small" className="text-center">
              <div className="text-lg font-bold text-green-600">🏆</div>
              <div className="text-sm text-gray-500 mt-1">最佳策略</div>
              <div className="font-medium">{result.comparison.best_strategy}</div>
            </Card>
            <Card size="small" className="text-center">
              <div className="text-lg font-bold text-blue-600">📊</div>
              <div className="text-sm text-gray-500 mt-1">对比总结</div>
              <div className="font-medium text-sm">{result.comparison.summary}</div>
            </Card>
            <Card size="small" className="text-center">
              <div className="text-lg font-bold text-red-600">⚠️</div>
              <div className="text-sm text-gray-500 mt-1">最差策略</div>
              <div className="font-medium">{result.comparison.worst_strategy}</div>
            </Card>
          </div>
        </Card>
      )}

      {/* Individual Strategy Results */}
      <Card title="各策略详细结果">
        <div className="space-y-4">
          {result.results.map((strategyResult, index) => (
            <Card size="small" key={strategyResult.id}>
              <div className="flex justify-between items-center">
                <h3 className="font-bold">{strategyNames[index]}</h3>
                <div className="text-sm text-gray-500">
                  执行时间: {strategyResult.duration.toFixed(2)}ms
                </div>
              </div>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mt-4">
                <div>
                  <div className="text-sm text-gray-500">总收益率</div>
                  <div className="text-lg font-bold" style={{ color: getValueColor(strategyResult.performance_metrics.total_return) }}>
                    {(strategyResult.performance_metrics.total_return * 100).toFixed(2)}%
                  </div>
                </div>
                <div>
                  <div className="text-sm text-gray-500">年化收益率</div>
                  <div className="text-lg font-bold" style={{ color: getValueColor(strategyResult.performance_metrics.annualized_return) }}>
                    {(strategyResult.performance_metrics.annualized_return * 100).toFixed(2)}%
                  </div>
                </div>
                <div>
                  <div className="text-sm text-gray-500">夏普比率</div>
                  <div className="text-lg font-bold" style={{ color: getSharpeRatioColor(strategyResult.performance_metrics.sharpe_ratio) }}>
                    {strategyResult.performance_metrics.sharpe_ratio.toFixed(3)}
                  </div>
                </div>
                <div>
                  <div className="text-sm text-gray-500">胜率</div>
                  <div className="text-lg font-bold" style={{ color: getValueColor(strategyResult.performance_metrics.win_rate - 0.5) }}>
                    {(strategyResult.performance_metrics.win_rate * 100).toFixed(1)}%
                  </div>
                </div>
              </div>
            </Card>
          ))}
        </div>
      </Card>
    </div>
  );
}