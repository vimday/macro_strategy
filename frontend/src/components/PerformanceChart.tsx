'use client';

import React, { useMemo } from 'react';
import { Card, Tabs, Table } from 'antd';
import type { TabsProps, ColumnsType } from 'antd/es/tabs';
import ReactECharts from 'echarts-for-react';
import { BacktestResult, Trade, DailyReturn } from '@/types';
import { formatPercentage, formatCurrency, formatDate } from '@/lib/utils';

interface PerformanceChartProps {
  result: BacktestResult;
}

export function PerformanceChart({ result }: PerformanceChartProps) {
  // Portfolio value chart data
  const portfolioValueOption = useMemo(() => {
    const dates = result.daily_returns.map(dr => formatDate(dr.date));
    const values = result.daily_returns.map(dr => dr.portfolio_value);

    return {
      title: {
        text: '投资组合价值变化',
        left: 'center',
        textStyle: { fontSize: 16 },
      },
      tooltip: {
        trigger: 'axis',
        formatter: (params: any) => {
          const data = params[0];
          return `
            日期: ${data.name}<br/>
            投资组合价值: ${formatCurrency(data.value)}
          `;
        },
      },
      xAxis: {
        type: 'category',
        data: dates,
        axisLabel: {
          formatter: (value: string) => {
            const date = new Date(value);
            return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}`;
          },
        },
      },
      yAxis: {
        type: 'value',
        axisLabel: {
          formatter: (value: number) => formatCurrency(value),
        },
      },
      series: [
        {
          name: '投资组合价值',
          type: 'line',
          data: values,
          smooth: true,
          lineStyle: { color: '#1890ff', width: 2 },
          areaStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                { offset: 0, color: 'rgba(24, 144, 255, 0.3)' },
                { offset: 1, color: 'rgba(24, 144, 255, 0.1)' },
              ],
            },
          },
        },
      ],
      grid: {
        left: '10%',
        right: '10%',
        bottom: '15%',
        top: '15%',
      },
    };
  }, [result.daily_returns]);

  // Cumulative returns chart
  const cumulativeReturnsOption = useMemo(() => {
    const dates = result.daily_returns.map(dr => formatDate(dr.date));
    const returns = result.daily_returns.map(dr => dr.cumulative_return * 100);

    return {
      title: {
        text: '累计收益率',
        left: 'center',
        textStyle: { fontSize: 16 },
      },
      tooltip: {
        trigger: 'axis',
        formatter: (params: any) => {
          const data = params[0];
          return `
            日期: ${data.name}<br/>
            累计收益率: ${data.value.toFixed(2)}%
          `;
        },
      },
      xAxis: {
        type: 'category',
        data: dates,
        axisLabel: {
          formatter: (value: string) => {
            const date = new Date(value);
            return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}`;
          },
        },
      },
      yAxis: {
        type: 'value',
        axisLabel: {
          formatter: (value: number) => `${value}%`,
        },
      },
      series: [
        {
          name: '累计收益率',
          type: 'line',
          data: returns,
          smooth: true,
          lineStyle: { color: '#52c41a', width: 2 },
          areaStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                { offset: 0, color: 'rgba(82, 196, 26, 0.3)' },
                { offset: 1, color: 'rgba(82, 196, 26, 0.1)' },
              ],
            },
          },
        },
      ],
      grid: {
        left: '10%',
        right: '10%',
        bottom: '15%',
        top: '15%',
      },
    };
  }, [result.daily_returns]);

  // Drawdown chart
  const drawdownOption = useMemo(() => {
    const dates = result.daily_returns.map(dr => formatDate(dr.date));
    const drawdowns = result.daily_returns.map(dr => -dr.drawdown * 100);

    return {
      title: {
        text: '回撤分析',
        left: 'center',
        textStyle: { fontSize: 16 },
      },
      tooltip: {
        trigger: 'axis',
        formatter: (params: any) => {
          const data = params[0];
          return `
            日期: ${data.name}<br/>
            回撤: ${data.value.toFixed(2)}%
          `;
        },
      },
      xAxis: {
        type: 'category',
        data: dates,
        axisLabel: {
          formatter: (value: string) => {
            const date = new Date(value);
            return `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}`;
          },
        },
      },
      yAxis: {
        type: 'value',
        max: 0,
        axisLabel: {
          formatter: (value: number) => `${value}%`,
        },
      },
      series: [
        {
          name: '回撤',
          type: 'line',
          data: drawdowns,
          smooth: true,
          lineStyle: { color: '#f5222d', width: 2 },
          areaStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                { offset: 0, color: 'rgba(245, 34, 45, 0.3)' },
                { offset: 1, color: 'rgba(245, 34, 45, 0.1)' },
              ],
            },
          },
        },
      ],
      grid: {
        left: '10%',
        right: '10%',
        bottom: '15%',
        top: '15%',
      },
    };
  }, [result.daily_returns]);

  // Trade analysis chart
  const tradeAnalysisOption = useMemo(() => {
    if (result.trades.length === 0) return null;

    const monthlyData: { [key: string]: { buy: number; sell: number } } = {};
    
    result.trades.forEach(trade => {
      const date = new Date(trade.date);
      const monthKey = `${date.getFullYear()}-${(date.getMonth() + 1).toString().padStart(2, '0')}`;
      
      if (!monthlyData[monthKey]) {
        monthlyData[monthKey] = { buy: 0, sell: 0 };
      }
      
      if (trade.action === 'buy') {
        monthlyData[monthKey].buy += 1;
      } else {
        monthlyData[monthKey].sell += 1;
      }
    });

    const months = Object.keys(monthlyData).sort();
    const buyData = months.map(month => monthlyData[month].buy);
    const sellData = months.map(month => monthlyData[month].sell);

    return {
      title: {
        text: '交易分布',
        left: 'center',
        textStyle: { fontSize: 16 },
      },
      tooltip: {
        trigger: 'axis',
        axisPointer: { type: 'shadow' },
      },
      legend: {
        data: ['买入', '卖出'],
        top: '10%',
      },
      xAxis: {
        type: 'category',
        data: months,
      },
      yAxis: {
        type: 'value',
      },
      series: [
        {
          name: '买入',
          type: 'bar',
          data: buyData,
          itemStyle: { color: '#52c41a' },
        },
        {
          name: '卖出',
          type: 'bar',
          data: sellData,
          itemStyle: { color: '#f5222d' },
        },
      ],
      grid: {
        left: '10%',
        right: '10%',
        bottom: '15%',
        top: '25%',
      },
    };
  }, [result.trades]);

  // Trade table columns
  const tradeColumns: ColumnsType<Trade> = [
    {
      title: '日期',
      dataIndex: 'date',
      key: 'date',
      render: (date: string) => formatDate(date),
      sorter: (a, b) => new Date(a.date).getTime() - new Date(b.date).getTime(),
    },
    {
      title: '操作',
      dataIndex: 'action',
      key: 'action',
      render: (action: string) => (
        <span className={action === 'buy' ? 'text-green-600' : 'text-red-600'}>
          {action === 'buy' ? '买入' : '卖出'}
        </span>
      ),
      filters: [
        { text: '买入', value: 'buy' },
        { text: '卖出', value: 'sell' },
      ],
      onFilter: (value, record) => record.action === value,
    },
    {
      title: '价格',
      dataIndex: 'price',
      key: 'price',
      render: (price: number) => price.toFixed(2),
      sorter: (a, b) => a.price - b.price,
    },
    {
      title: '数量',
      dataIndex: 'quantity',
      key: 'quantity',
      render: (quantity: number) => quantity.toFixed(0),
      sorter: (a, b) => a.quantity - b.quantity,
    },
    {
      title: '金额',
      dataIndex: 'amount',
      key: 'amount',
      render: (amount: number) => formatCurrency(amount),
      sorter: (a, b) => a.amount - b.amount,
    },
    {
      title: '手续费',
      dataIndex: 'commission',
      key: 'commission',
      render: (commission: number) => formatCurrency(commission),
      sorter: (a, b) => a.commission - b.commission,
    },
  ];

  const tabItems: TabsProps['items'] = [
    {
      key: 'portfolio',
      label: '投资组合价值',
      children: (
        <ReactECharts 
          option={portfolioValueOption} 
          style={{ height: '400px', width: '100%' }}
          opts={{ renderer: 'canvas' }}
        />
      ),
    },
    {
      key: 'returns',
      label: '累计收益',
      children: (
        <ReactECharts 
          option={cumulativeReturnsOption} 
          style={{ height: '400px', width: '100%' }}
          opts={{ renderer: 'canvas' }}
        />
      ),
    },
    {
      key: 'drawdown',
      label: '回撤分析',
      children: (
        <ReactECharts 
          option={drawdownOption} 
          style={{ height: '400px', width: '100%' }}
          opts={{ renderer: 'canvas' }}
        />
      ),
    },
    {
      key: 'trades',
      label: '交易分布',
      children: tradeAnalysisOption ? (
        <ReactECharts 
          option={tradeAnalysisOption} 
          style={{ height: '400px', width: '100%' }}
          opts={{ renderer: 'canvas' }}
        />
      ) : (
        <div className="flex items-center justify-center h-96 text-gray-500">
          暂无交易数据
        </div>
      ),
    },
    {
      key: 'tradeTable',
      label: '交易明细',
      children: (
        <Table
          columns={tradeColumns}
          dataSource={result.trades.map((trade, index) => ({ ...trade, key: index }))}
          pagination={{
            pageSize: 10,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => `共 ${total} 条记录`,
          }}
          scroll={{ x: 800 }}
          size="small"
        />
      ),
    },
  ];

  return (
    <Card title="绩效分析图表" className="w-full">
      <Tabs items={tabItems} type="card" />
    </Card>
  );
}