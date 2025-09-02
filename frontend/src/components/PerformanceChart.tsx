'use client';

import React, { useMemo } from 'react';
import { Card, Tabs } from 'antd';
import { LineChartIcon, BarChart3Icon, TrendingDownIcon } from 'lucide-react';
import ReactECharts from 'echarts-for-react';
import { BacktestResult } from '@/types';

interface PerformanceChartProps {
  result: BacktestResult;
}

export function PerformanceChart({ result }: PerformanceChartProps) {
  const { daily_returns, trades } = result;

  // 准备图表数据
  const chartData = useMemo(() => {
    const portfolioData = daily_returns.map(dr => [dr.date, dr.portfolio_value]);
    const returnData = daily_returns.map(dr => [dr.date, dr.cumulative_return * 100]);
    const drawdownData = daily_returns.map(dr => [dr.date, -dr.drawdown * 100]);
    
    // 交易标记点
    const tradeMarks = trades.map(trade => ({
      name: trade.action === 'buy' ? '买入' : '卖出',
      coord: [trade.date, daily_returns.find(dr => dr.date === trade.date)?.portfolio_value || 0],
      symbol: trade.action === 'buy' ? 'triangle' : 'diamond',
      symbolSize: 8,
      itemStyle: {
        color: trade.action === 'buy' ? '#10b981' : '#ef4444'
      }
    }));

    return {
      portfolioData,
      returnData,
      drawdownData,
      tradeMarks
    };
  }, [daily_returns, trades]);

  // 资产净值曲线配置
  const portfolioValueOption = {
    title: {
      text: '资产净值曲线',
      left: 'center',
      textStyle: { fontSize: 16 }
    },
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const [date, value] = params[0].data;
        return `${date}<br/>资产净值: ¥${value.toLocaleString()}`;
      }
    },
    xAxis: {
      type: 'time',
      splitLine: { show: false }
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        formatter: (value: number) => `¥${(value / 10000).toFixed(1)}万`
      }
    },
    series: [
      {
        name: '资产净值',
        type: 'line',
        data: chartData.portfolioData,
        smooth: true,
        lineStyle: { width: 2, color: '#3b82f6' },
        areaStyle: { opacity: 0.1, color: '#3b82f6' },
        markPoint: {
          data: chartData.tradeMarks,
          symbol: 'pin',
          symbolSize: 30
        }
      }
    ],
    grid: { top: 60, bottom: 60, left: 80, right: 40 },
    dataZoom: [
      { type: 'inside', start: 0, end: 100 },
      { type: 'slider', start: 0, end: 100, height: 30 }
    ]
  };

  // 累计收益率配置
  const cumulativeReturnOption = {
    title: {
      text: '累计收益率',
      left: 'center',
      textStyle: { fontSize: 16 }
    },
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const [date, value] = params[0].data;
        return `${date}<br/>累计收益率: ${value.toFixed(2)}%`;
      }
    },
    xAxis: {
      type: 'time',
      splitLine: { show: false }
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        formatter: (value: number) => `${value.toFixed(1)}%`
      }
    },
    series: [
      {
        name: '累计收益率',
        type: 'line',
        data: chartData.returnData,
        smooth: true,
        lineStyle: { width: 2, color: '#10b981' },
        areaStyle: { 
          opacity: 0.1, 
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: '#10b981' },
              { offset: 1, color: 'rgba(16, 185, 129, 0.1)' }
            ]
          }
        }
      }
    ],
    grid: { top: 60, bottom: 60, left: 60, right: 40 },
    dataZoom: [
      { type: 'inside', start: 0, end: 100 },
      { type: 'slider', start: 0, end: 100, height: 30 }
    ]
  };

  // 回撤曲线配置
  const drawdownOption = {
    title: {
      text: '回撤曲线',
      left: 'center',
      textStyle: { fontSize: 16 }
    },
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        const [date, value] = params[0].data;
        return `${date}<br/>回撤: ${Math.abs(value).toFixed(2)}%`;
      }
    },
    xAxis: {
      type: 'time',
      splitLine: { show: false }
    },
    yAxis: {
      type: 'value',
      max: 0,
      axisLabel: {
        formatter: (value: number) => `${value.toFixed(1)}%`
      }
    },
    series: [
      {
        name: '回撤',
        type: 'line',
        data: chartData.drawdownData,
        smooth: true,
        lineStyle: { width: 2, color: '#ef4444' },
        areaStyle: { 
          opacity: 0.2,
          color: {
            type: 'linear',
            x: 0, y: 0, x2: 0, y2: 1,
            colorStops: [
              { offset: 0, color: 'rgba(239, 68, 68, 0.4)' },
              { offset: 1, color: 'rgba(239, 68, 68, 0.1)' }
            ]
          }
        }
      }
    ],
    grid: { top: 60, bottom: 60, left: 60, right: 40 },
    dataZoom: [
      { type: 'inside', start: 0, end: 100 },
      { type: 'slider', start: 0, end: 100, height: 30 }
    ]
  };

  // 交易分布图配置
  const tradeDistributionOption = {
    title: {
      text: '交易收益分布',
      left: 'center',
      textStyle: { fontSize: 16 }
    },
    tooltip: {
      trigger: 'axis',
      axisPointer: { type: 'shadow' }
    },
    xAxis: {
      type: 'category',
      data: trades.filter((_, i) => i % 2 === 1).map((_, i) => `交易 ${i + 1}`)
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        formatter: (value: number) => `${value.toFixed(1)}%`
      }
    },
    series: [
      {
        name: '交易收益',
        type: 'bar',
        data: trades.filter((_, i) => i % 2 === 1).map((trade, i) => {
          const buyTrade = trades[i * 2];
          if (buyTrade) {
            const profit = ((trade.price - buyTrade.price) / buyTrade.price) * 100;
            return {
              value: profit,
              itemStyle: {
                color: profit >= 0 ? '#10b981' : '#ef4444'
              }
            };
          }
          return 0;
        }),
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        }
      }
    ],
    grid: { top: 60, bottom: 60, left: 60, right: 40 }
  };

  const tabItems = [
    {
      key: 'portfolio',
      label: (
        <div className="flex items-center">
          <LineChartIcon className="mr-1 h-4 w-4" />
          资产净值
        </div>
      ),
      children: (
        <div className="h-96">
          <ReactECharts option={portfolioValueOption} style={{ height: '100%' }} />
        </div>
      )
    },
    {
      key: 'returns',
      label: (
        <div className="flex items-center">
          <LineChartIcon className="mr-1 h-4 w-4" />
          累计收益
        </div>
      ),
      children: (
        <div className="h-96">
          <ReactECharts option={cumulativeReturnOption} style={{ height: '100%' }} />
        </div>
      )
    },
    {
      key: 'drawdown',
      label: (
        <div className="flex items-center">
          <TrendingDownIcon className="mr-1 h-4 w-4" />
          回撤分析
        </div>
      ),
      children: (
        <div className="h-96">
          <ReactECharts option={drawdownOption} style={{ height: '100%' }} />
        </div>
      )
    },
    {
      key: 'trades',
      label: (
        <div className="flex items-center">
          <BarChart3Icon className="mr-1 h-4 w-4" />
          交易分析
        </div>
      ),
      children: (
        <div className="h-96">
          <ReactECharts option={tradeDistributionOption} style={{ height: '100%' }} />
        </div>
      )
    }
  ];

  return (
    <Card title="性能图表分析">
      <Tabs items={tabItems} size="large" />
    </Card>
  );
}