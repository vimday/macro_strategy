'use client';

import React from 'react';
import { Card, Alert, Spin } from 'antd';
import { TrophyIcon, AlertCircleIcon } from 'lucide-react';
import { BacktestResult } from '@/types';
import { MetricsDisplay } from './MetricsDisplay';
import { PerformanceChart } from './PerformanceChart';

interface BacktestResultsProps {
  result: BacktestResult;
  loading?: boolean;
  error?: string | null;
}

export function BacktestResults({ result, loading = false, error = null }: BacktestResultsProps) {
  if (loading) {
    return (
      <Card className="mt-8">
        <div className="flex items-center justify-center py-12">
          <Spin size="large" />
          <span className="ml-4 text-lg">正在执行回测策略...</span>
        </div>
      </Card>
    );
  }

  if (error) {
    return (
      <Card className="mt-8">
        <Alert
          message="回测执行失败"
          description={error}
          type="error"
          icon={<AlertCircleIcon className="h-4 w-4" />}
          showIcon
        />
      </Card>
    );
  }

  if (!result) {
    return null;
  }

  return (
    <div className="mt-8 space-y-6">
      {/* Results Header */}
      <Card>
        <div className="flex items-center justify-between">
          <div className="flex items-center">
            <TrophyIcon className="h-6 w-6 text-yellow-500 mr-3" />
            <div>
              <h3 className="text-xl font-bold text-gray-900">
                回测结果 - {result.request.index_id}
              </h3>
              <p className="text-gray-600 mt-1">
                策略类型: {result.request.strategy.type === 'monthly_rotation' ? '月度轮动策略' : result.request.strategy.type}
              </p>
              <p className="text-gray-500 text-sm mt-1">
                回测期间: {result.request.start_date} 至 {result.request.end_date}
              </p>
            </div>
          </div>
          <div className="text-right">
            <div className="text-sm text-gray-500">回测ID</div>
            <div className="text-sm font-mono text-gray-700">{result.id}</div>
            <div className="text-xs text-gray-400 mt-1">
              创建时间: {new Date(result.created_at).toLocaleString('zh-CN')}
            </div>
          </div>
        </div>
      </Card>

      {/* Strategy Configuration Summary */}
      <Card title="策略配置" size="small">
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <div>
            <div className="text-sm text-gray-500">指数</div>
            <div className="font-medium">{result.request.index_id}</div>
          </div>
          <div>
            <div className="text-sm text-gray-500">初始资金</div>
            <div className="font-medium">¥{result.request.initial_cash.toLocaleString()}</div>
          </div>
          <div>
            <div className="text-sm text-gray-500">买入时机</div>
            <div className="font-medium">
              月末前 {(result.request.strategy.parameters as any).buy_days_before_month_end} 个交易日
            </div>
          </div>
          <div>
            <div className="text-sm text-gray-500">卖出时机</div>
            <div className="font-medium">
              月初第 {(result.request.strategy.parameters as any).sell_days_after_month_start} 个交易日
            </div>
          </div>
        </div>
      </Card>

      {/* Performance Metrics */}
      <MetricsDisplay 
        metrics={result.performance_metrics} 
        initialCash={result.request.initial_cash}
      />

      {/* Performance Charts */}
      <PerformanceChart result={result} />

      {/* Quick Summary */}
      <Card title="回测总结" size="small">
        <div className="bg-gray-50 p-4 rounded-lg">
          <div className="prose prose-sm max-w-none">
            <p>
              在 {result.request.start_date} 至 {result.request.end_date} 期间，
              使用月度轮动策略对 <strong>{result.request.index_id}</strong> 进行回测。
            </p>
            <ul className="mt-2">
              <li>
                <strong>总收益率:</strong> {(result.performance_metrics.total_return * 100).toFixed(2)}%
                (年化收益率: {(result.performance_metrics.annualized_return * 100).toFixed(2)}%)
              </li>
              <li>
                <strong>风险指标:</strong> 最大回撤 {(result.performance_metrics.max_drawdown * 100).toFixed(2)}%，
                夏普比率 {result.performance_metrics.sharpe_ratio.toFixed(3)}
              </li>
              <li>
                <strong>交易统计:</strong> 共执行 {result.performance_metrics.total_trades} 笔交易，
                胜率 {(result.performance_metrics.win_rate * 100).toFixed(1)}%
              </li>
            </ul>
            <div className="mt-3 text-xs text-gray-600">
              <strong>风险提示:</strong> 回测结果仅供参考，不构成投资建议。过往业绩不代表未来收益。
            </div>
          </div>
        </div>
      </Card>
    </div>
  );
}