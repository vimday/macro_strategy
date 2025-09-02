'use client';

import React from 'react';
import { Card, Row, Col, Statistic, Descriptions, Tag } from 'antd';
import { 
  TrendingUpIcon, 
  TrendingDownIcon, 
  TrophyIcon,
  AlertTriangleIcon,
  BarChart3Icon,
  PercentIcon 
} from 'lucide-react';
import { BacktestResult } from '@/types';
import { formatPercentage, formatNumber, getValueColor } from '@/lib/utils';

interface MetricsDisplayProps {
  result: BacktestResult;
}

export function MetricsDisplay({ result }: MetricsDisplayProps) {
  const { performance_metrics: metrics } = result;

  return (
    <div className="mb-6">
      {/* 核心指标 */}
      <Card 
        title={
          <div className="flex items-center">
            <TrophyIcon className="mr-2 h-5 w-5" />
            核心业绩指标
          </div>
        } 
        className="mb-4"
      >
        <Row gutter={[16, 16]}>
          <Col xs={24} sm={12} md={6}>
            <Statistic
              title="总收益率"
              value={formatPercentage(metrics.total_return)}
              valueStyle={{ color: metrics.total_return >= 0 ? '#10b981' : '#ef4444' }}
              prefix={metrics.total_return >= 0 ? <TrendingUpIcon className="h-4 w-4" /> : <TrendingDownIcon className="h-4 w-4" />}
            />
          </Col>
          <Col xs={24} sm={12} md={6}>
            <Statistic
              title="年化收益率"
              value={formatPercentage(metrics.annualized_return)}
              valueStyle={{ color: metrics.annualized_return >= 0 ? '#10b981' : '#ef4444' }}
              prefix={<BarChart3Icon className="h-4 w-4" />}
            />
          </Col>
          <Col xs={24} sm={12} md={6}>
            <Statistic
              title="最大回撤"
              value={formatPercentage(metrics.max_drawdown)}
              valueStyle={{ color: '#ef4444' }}
              prefix={<TrendingDownIcon className="h-4 w-4" />}
            />
          </Col>
          <Col xs={24} sm={12} md={6}>
            <Statistic
              title="夏普比率"
              value={formatNumber(metrics.sharpe_ratio, 3)}
              valueStyle={{ color: metrics.sharpe_ratio >= 0 ? '#10b981' : '#ef4444' }}
              prefix={<PercentIcon className="h-4 w-4" />}
            />
          </Col>
        </Row>
      </Card>

      {/* 风险指标 */}
      <Card title="风险指标" className="mb-4">
        <Row gutter={[16, 16]}>
          <Col xs={24} sm={12} md={6}>
            <Statistic
              title="波动率"
              value={formatPercentage(metrics.volatility)}
              valueStyle={{ color: '#6b7280' }}
            />
          </Col>
          <Col xs={24} sm={12} md={6}>
            <Statistic
              title="索提诺比率"
              value={formatNumber(metrics.sortino_ratio, 3)}
              valueStyle={{ color: metrics.sortino_ratio >= 0 ? '#10b981' : '#ef4444' }}
            />
          </Col>
          <Col xs={24} sm={12} md={6}>
            <Statistic
              title="卡玛比率"
              value={formatNumber(metrics.calmar_ratio, 3)}
              valueStyle={{ color: metrics.calmar_ratio >= 0 ? '#10b981' : '#ef4444' }}
            />
          </Col>
          <Col xs={24} sm={12} md={6}>
            <Statistic
              title="最大回撤期"
              value={`${metrics.max_drawdown_period} 天`}
              valueStyle={{ color: '#6b7280' }}
            />
          </Col>
        </Row>
      </Card>

      {/* 交易指标 */}
      <Card title="交易统计" className="mb-4">
        <Row gutter={[16, 16]}>
          <Col xs={24} sm={12} md={6}>
            <Statistic
              title="总交易次数"
              value={metrics.total_trades}
              valueStyle={{ color: '#6b7280' }}
            />
          </Col>
          <Col xs={24} sm={12} md={6}>
            <Statistic
              title="胜率"
              value={formatPercentage(metrics.win_rate)}
              valueStyle={{ 
                color: metrics.win_rate > 0.5 ? '#10b981' : 
                       metrics.win_rate > 0.4 ? '#f59e0b' : '#ef4444' 
              }}
            />
          </Col>
          <Col xs={24} sm={12} md={6}>
            <Statistic
              title="盈亏比"
              value={formatNumber(metrics.profit_factor, 2)}
              valueStyle={{ color: metrics.profit_factor >= 1 ? '#10b981' : '#ef4444' }}
            />
          </Col>
          <Col xs={24} sm={12} md={6}>
            <Statistic
              title="恢复期"
              value={`${metrics.recovery_period} 天`}
              valueStyle={{ color: '#6b7280' }}
            />
          </Col>
        </Row>
      </Card>

      {/* 详细信息 */}
      <Card title="详细信息">
        <Descriptions bordered column={{ xs: 1, sm: 1, md: 2 }} size="small">
          <Descriptions.Item label="盈利交易数">
            <Tag color="green">{metrics.winning_trades}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="亏损交易数">
            <Tag color="red">{metrics.losing_trades}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="平均盈利交易">
            <span className="text-green-600">
              {formatPercentage(metrics.avg_winning_trade)}
            </span>
          </Descriptions.Item>
          <Descriptions.Item label="平均亏损交易">
            <span className="text-red-600">
              {formatPercentage(metrics.avg_losing_trade)}
            </span>
          </Descriptions.Item>
          <Descriptions.Item label="最大单笔盈利">
            <span className="text-green-600">
              {formatPercentage(metrics.max_winning_trade)}
            </span>
          </Descriptions.Item>
          <Descriptions.Item label="最大单笔亏损">
            <span className="text-red-600">
              {formatPercentage(metrics.max_losing_trade)}
            </span>
          </Descriptions.Item>
          <Descriptions.Item label="策略类型">
            <Tag color="blue">{result.request.strategy.type}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="回测时长">
            <Tag>{`${Math.round(result.duration / 1000000)} ms`}</Tag>
          </Descriptions.Item>
        </Descriptions>
      </Card>
    </div>
  );
}