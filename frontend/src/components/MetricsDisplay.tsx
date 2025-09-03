'use client';

import React from 'react';
import { Card, Row, Col, Statistic, Divider } from 'antd';
import {
  TrendingUpIcon,
  TrendingDownIcon,
  BarChart3Icon,
  PieChartIcon,
  ActivityIcon,
  DollarSignIcon,
  CalendarIcon,
  TargetIcon,
} from 'lucide-react';
import { PerformanceMetrics } from '@/types';
import {
  formatPercentage,
  formatCurrency,
  getSharpeRatioColor,
  getDrawdownColor,
  getValueColor,
} from '@/lib/utils';

interface MetricsDisplayProps {
  metrics: PerformanceMetrics;
  initialCash: number;
}

export function MetricsDisplay({ metrics, initialCash }: MetricsDisplayProps) {
  const finalValue = initialCash * (1 + metrics.total_return);

  const metricCards = [
    {
      title: '总收益率',
      value: formatPercentage(metrics.total_return),
      icon: <TrendingUpIcon className="h-6 w-6" />,
      color: getValueColor(metrics.total_return),
      description: '整个回测期间的总收益',
    },
    {
      title: '年化收益率',
      value: formatPercentage(metrics.annualized_return),
      icon: <BarChart3Icon className="h-6 w-6" />,
      color: getValueColor(metrics.annualized_return),
      description: '按年计算的平均收益率',
    },
    {
      title: '最大回撤',
      value: formatPercentage(-metrics.max_drawdown),
      icon: <TrendingDownIcon className="h-6 w-6" />,
      color: getDrawdownColor(metrics.max_drawdown),
      description: '投资组合的最大损失幅度',
    },
    {
      title: '夏普比率',
      value: metrics.sharpe_ratio.toFixed(3),
      icon: <ActivityIcon className="h-6 w-6" />,
      color: getSharpeRatioColor(metrics.sharpe_ratio),
      description: '风险调整后的收益指标',
    },
    {
      title: '索提诺比率',
      value: metrics.sortino_ratio.toFixed(3),
      icon: <TargetIcon className="h-6 w-6" />,
      color: getSharpeRatioColor(metrics.sortino_ratio),
      description: '下行风险调整后的收益指标',
    },
    {
      title: '波动率',
      value: formatPercentage(metrics.volatility),
      icon: <PieChartIcon className="h-6 w-6" />,
      color: 'text-blue-600',
      description: '收益率的标准差',
    },
    {
      title: '卡玛比率',
      value: metrics.calmar_ratio.toFixed(3),
      icon: <DollarSignIcon className="h-6 w-6" />,
      color: getSharpeRatioColor(metrics.calmar_ratio),
      description: '年化收益与最大回撤的比值',
    },
    {
      title: '胜率',
      value: formatPercentage(metrics.win_rate),
      icon: <CalendarIcon className="h-6 w-6" />,
      color: getValueColor(metrics.win_rate - 0.5),
      description: '盈利交易占总交易的比例',
    },
  ];

  const tradeMetrics = [
    {
      label: '总交易次数',
      value: metrics.total_trades,
      suffix: '次',
    },
    {
      label: '盈利交易',
      value: metrics.winning_trades,
      suffix: '次',
      valueStyle: { color: '#52c41a' },
    },
    {
      label: '亏损交易',
      value: metrics.losing_trades,
      suffix: '次',
      valueStyle: { color: '#f5222d' },
    },
    {
      label: '盈亏比',
      value: metrics.profit_factor.toFixed(2),
      suffix: '',
      valueStyle: { color: getValueColor(metrics.profit_factor - 1).replace('text-', '#') },
    },
    {
      label: '平均盈利交易',
      value: formatPercentage(metrics.avg_winning_trade),
      suffix: '',
      valueStyle: { color: '#52c41a' },
    },
    {
      label: '平均亏损交易',
      value: formatPercentage(metrics.avg_losing_trade),
      suffix: '',
      valueStyle: { color: '#f5222d' },
    },
    {
      label: '最大单笔盈利',
      value: formatPercentage(metrics.max_winning_trade),
      suffix: '',
      valueStyle: { color: '#52c41a' },
    },
    {
      label: '最大单笔亏损',
      value: formatPercentage(metrics.max_losing_trade),
      suffix: '',
      valueStyle: { color: '#f5222d' },
    },
  ];

  const portfolioMetrics = [
    {
      label: '初始资金',
      value: formatCurrency(initialCash),
      suffix: '',
    },
    {
      label: '最终资金',
      value: formatCurrency(finalValue),
      suffix: '',
      valueStyle: { color: getValueColor(metrics.total_return).replace('text-', '#') },
    },
    {
      label: '绝对收益',
      value: formatCurrency(finalValue - initialCash),
      suffix: '',
      valueStyle: { color: getValueColor(metrics.total_return).replace('text-', '#') },
    },
    {
      label: '最大回撤期',
      value: metrics.max_drawdown_period,
      suffix: '天',
    },
    {
      label: '恢复期',
      value: metrics.recovery_period,
      suffix: '天',
    },
  ];

  return (
    <div className="space-y-6">
      {/* Main Performance Metrics */}
      <Card title="核心绩效指标" className="w-full">
        <Row gutter={[16, 16]}>
          {metricCards.map((metric, index) => (
            <Col xs={12} sm={8} md={6} key={index}>
              <Card 
                size="small" 
                className="text-center hover:shadow-md transition-shadow"
                styles={{ body: { padding: '16px 8px' } }}
              >
                <div className="flex flex-col items-center space-y-2">
                  <div className={`${metric.color}`}>
                    {metric.icon}
                  </div>
                  <div className="text-xs text-gray-500 text-center leading-tight">
                    {metric.title}
                  </div>
                  <div className={`text-lg font-bold ${metric.color}`}>
                    {metric.value}
                  </div>
                  <div className="text-xs text-gray-400 text-center leading-tight">
                    {metric.description}
                  </div>
                </div>
              </Card>
            </Col>
          ))}
        </Row>
      </Card>

      <Row gutter={16}>
        {/* Portfolio Summary */}
        <Col xs={24} lg={12}>
          <Card title="投资组合概要" size="small">
            <Row gutter={[8, 8]}>
              {portfolioMetrics.map((metric, index) => (
                <Col xs={12} key={index}>
                  <Statistic
                    title={metric.label}
                    value={metric.value}
                    suffix={metric.suffix}
                    valueStyle={metric.valueStyle}
                  />
                </Col>
              ))}
            </Row>
          </Card>
        </Col>

        {/* Trade Analysis */}
        <Col xs={24} lg={12}>
          <Card title="交易分析" size="small">
            <Row gutter={[8, 8]}>
              {tradeMetrics.map((metric, index) => (
                <Col xs={12} key={index}>
                  <Statistic
                    title={metric.label}
                    value={metric.value}
                    suffix={metric.suffix}
                    valueStyle={metric.valueStyle}
                  />
                </Col>
              ))}
            </Row>
          </Card>
        </Col>
      </Row>

      {/* Risk Assessment */}
      <Card title="风险评估" size="small">
        <Row gutter={16}>
          <Col xs={24} sm={8}>
            <div className="text-center p-4">
              <div className="text-2xl font-bold mb-2">
                <span className={getSharpeRatioColor(metrics.sharpe_ratio)}>
                  {metrics.sharpe_ratio >= 2 ? '优秀' : 
                   metrics.sharpe_ratio >= 1 ? '良好' : 
                   metrics.sharpe_ratio >= 0.5 ? '一般' : '较差'}
                </span>
              </div>
              <div className="text-sm text-gray-500">夏普比率评级</div>
              <div className="text-xs text-gray-400 mt-1">
                基于风险调整后收益
              </div>
            </div>
          </Col>
          <Col xs={24} sm={8}>
            <div className="text-center p-4">
              <div className="text-2xl font-bold mb-2">
                <span className={getDrawdownColor(metrics.max_drawdown)}>
                  {Math.abs(metrics.max_drawdown) <= 0.05 ? '低风险' : 
                   Math.abs(metrics.max_drawdown) <= 0.10 ? '中等风险' : 
                   Math.abs(metrics.max_drawdown) <= 0.20 ? '较高风险' : '高风险'}
                </span>
              </div>
              <div className="text-sm text-gray-500">回撤风险评级</div>
              <div className="text-xs text-gray-400 mt-1">
                基于最大回撤幅度
              </div>
            </div>
          </Col>
          <Col xs={24} sm={8}>
            <div className="text-center p-4">
              <div className="text-2xl font-bold mb-2">
                <span className={getValueColor(metrics.win_rate - 0.5)}>
                  {metrics.win_rate >= 0.7 ? '高胜率' : 
                   metrics.win_rate >= 0.5 ? '平衡' : '低胜率'}
                </span>
              </div>
              <div className="text-sm text-gray-500">交易胜率评级</div>
              <div className="text-xs text-gray-400 mt-1">
                基于盈利交易比例
              </div>
            </div>
          </Col>
        </Row>
      </Card>
    </div>
  );
}