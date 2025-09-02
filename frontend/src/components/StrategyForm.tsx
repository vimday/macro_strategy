'use client';

import React from 'react';
import {
  Card,
  Form,
  Select,
  InputNumber,
  Button,
  Row,
  Col,
  DatePicker,
} from 'antd';
import { PlayIcon, BarChart3Icon } from 'lucide-react';
import dayjs from 'dayjs';
import { useIndexes } from '@/hooks/useBacktest';
import { BacktestFormData, StrategyType } from '@/types';

const { Option } = Select;

interface StrategyFormProps {
  onSubmit: (values: BacktestFormData) => void;
  loading: boolean;
}

export function StrategyForm({ onSubmit, loading }: StrategyFormProps) {
  const [form] = Form.useForm();
  const { data: indexes, isLoading: indexesLoading } = useIndexes();

  const handleSubmit = (values: any) => {
    const formData: BacktestFormData = {
      indexId: values.indexId,
      strategyType: values.strategyType,
      buyDaysBeforeMonthEnd: values.buyDaysBeforeMonthEnd,
      sellDaysAfterMonthStart: values.sellDaysAfterMonthStart,
      startDate: values.dateRange[0].format('YYYY-MM-DD'),
      endDate: values.dateRange[1].format('YYYY-MM-DD'),
      initialCash: values.initialCash,
    };
    onSubmit(formData);
  };

  return (
    <Card 
      title={
        <div className="flex items-center">
          <PlayIcon className="mr-2 h-5 w-5" />
          策略配置
        </div>
      }
      className="mb-8"
    >
      <Form
        form={form}
        layout="vertical"
        onFinish={handleSubmit}
        initialValues={{
          strategyType: 'monthly_rotation' as StrategyType,
          buyDaysBeforeMonthEnd: 1,
          sellDaysAfterMonthStart: 1,
          initialCash: 1000000,
          dateRange: [dayjs('2020-01-01'), dayjs('2023-12-31')],
        }}
      >
        <Row gutter={16}>
          <Col xs={24} md={12}>
            <Form.Item
              name="indexId"
              label="选择指数"
              rules={[{ required: true, message: '请选择指数' }]}
            >
              <Select
                placeholder="请选择要回测的指数"
                loading={indexesLoading}
                showSearch
                filterOption={(input, option) =>
                  option?.children?.toString().toLowerCase().includes(input.toLowerCase())
                }
              >
                {indexes?.map(index => (
                  <Option key={index.id} value={index.id}>
                    {index.name} ({index.symbol})
                  </Option>
                ))}
              </Select>
            </Form.Item>
          </Col>
          <Col xs={24} md={12}>
            <Form.Item
              name="strategyType"
              label="策略类型"
              rules={[{ required: true, message: '请选择策略类型' }]}
            >
              <Select placeholder="选择策略类型">
                <Option value="monthly_rotation">月度轮动策略</Option>
              </Select>
            </Form.Item>
          </Col>
        </Row>

        <Row gutter={16}>
          <Col xs={24} md={12}>
            <Form.Item
              name="buyDaysBeforeMonthEnd"
              label="月末前几个交易日买入"
              rules={[
                { required: true, message: '请输入买入时机' },
                { type: 'number', min: 1, max: 20, message: '请输入1-20之间的数字' }
              ]}
            >
              <InputNumber 
                min={1} 
                max={20} 
                className="w-full"
                placeholder="1"
              />
            </Form.Item>
          </Col>
          <Col xs={24} md={12}>
            <Form.Item
              name="sellDaysAfterMonthStart"
              label="月初第几个交易日卖出"
              rules={[
                { required: true, message: '请输入卖出时机' },
                { type: 'number', min: 1, max: 20, message: '请输入1-20之间的数字' }
              ]}
            >
              <InputNumber 
                min={1} 
                max={20} 
                className="w-full"
                placeholder="1"
              />
            </Form.Item>
          </Col>
        </Row>

        <Row gutter={16}>
          <Col xs={24} md={12}>
            <Form.Item
              name="dateRange"
              label="回测时间范围"
              rules={[{ required: true, message: '请选择回测时间范围' }]}
            >
              <DatePicker.RangePicker
                className="w-full"
                format="YYYY-MM-DD"
                disabledDate={(current) => current && current > dayjs().endOf('day')}
              />
            </Form.Item>
          </Col>
          <Col xs={24} md={12}>
            <Form.Item
              name="initialCash"
              label="初始资金 (元)"
              rules={[
                { required: true, message: '请输入初始资金' },
                { type: 'number', min: 10000, message: '初始资金不能少于10,000元' }
              ]}
            >
              <InputNumber
                min={10000}
                max={100000000}
                className="w-full"
                formatter={value => `${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')}
                parser={value => value?.replace(/\$\s?|(,*)/g, '')}
                placeholder="1,000,000"
              />
            </Form.Item>
          </Col>
        </Row>

        <Form.Item>
          <Button
            type="primary"
            htmlType="submit"
            loading={loading}
            size="large"
            icon={<BarChart3Icon className="h-4 w-4" />}
            className="w-full"
          >
            {loading ? '回测运行中...' : '开始回测'}
          </Button>
        </Form.Item>
      </Form>
    </Card>
  );
}