#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
AKShare 数据获取脚本
用于从 AKShare 获取 A 股历史数据
"""

import sys
import json
import akshare as ak
import pandas as pd
from datetime import datetime
import warnings

warnings.filterwarnings('ignore')

def get_stock_zh_a_hist(symbol, start_date, end_date):
    """
    获取 A 股历史数据
    
    Args:
        symbol: 股票代码（如：sh000300）
        start_date: 开始日期（格式：20200101）
        end_date: 结束日期（格式：20231231）
    
    Returns:
        JSON 格式的历史数据
    """
    try:
        # 调用 AKShare 获取数据 - 使用更新的API
        df = ak.stock_zh_index_daily(symbol=symbol)
        
        # 过滤日期范围
        if not df.empty:
            df['date'] = pd.to_datetime(df['date'])
            start_dt = datetime.strptime(start_date, '%Y%m%d')
            end_dt = datetime.strptime(end_date, '%Y%m%d')
            df = df[(df['date'] >= start_dt) & (df['date'] <= end_dt)]
        
        # 重命名列以便 Go 代码解析
        df = df.rename(columns={
            'date': '日期',
            'open': '开盘',
            'high': '最高', 
            'low': '最低',
            'close': '收盘',
            'volume': '成交量'
        })
        
        # 确保日期格式正确
        if '日期' in df.columns:
            df['日期'] = df['日期'].dt.strftime('%Y-%m-%d')
        
        # 转换为字典列表
        data = df.to_dict('records')
        
        # 输出 JSON
        print(json.dumps(data, ensure_ascii=False, default=str))
        
    except Exception as e:
        # 输出错误信息
        error_data = {"error": str(e)}
        print(json.dumps(error_data, ensure_ascii=False))
        sys.exit(1)

def get_stock_info(symbol):
    """
    获取股票基本信息
    
    Args:
        symbol: 股票代码
        
    Returns:
        JSON 格式的股票信息
    """
    try:
        # 获取股票信息
        df = ak.stock_individual_info_em(symbol=symbol)
        data = df.to_dict('records')
        print(json.dumps(data, ensure_ascii=False, default=str))
        
    except Exception as e:
        error_data = {"error": str(e)}
        print(json.dumps(error_data, ensure_ascii=False))
        sys.exit(1)

def get_index_list():
    """
    获取指数列表
    
    Returns:
        JSON 格式的指数列表
    """
    try:
        # 获取主要指数列表
        df = ak.index_stock_info()
        data = df.to_dict('records')
        print(json.dumps(data, ensure_ascii=False, default=str))
        
    except Exception as e:
        error_data = {"error": str(e)}
        print(json.dumps(error_data, ensure_ascii=False))
        sys.exit(1)

def main():
    """
    主函数：解析命令行参数并调用相应函数
    """
    if len(sys.argv) < 2:
        print(json.dumps({"error": "缺少参数"}, ensure_ascii=False))
        sys.exit(1)
    
    command = sys.argv[1]
    
    if command == "get_stock_zh_a_hist":
        if len(sys.argv) != 5:
            print(json.dumps({"error": "参数不正确：需要 symbol, start_date, end_date"}, ensure_ascii=False))
            sys.exit(1)
        
        symbol = sys.argv[2]
        start_date = sys.argv[3]
        end_date = sys.argv[4]
        get_stock_zh_a_hist(symbol, start_date, end_date)
        
    elif command == "get_stock_info":
        if len(sys.argv) != 3:
            print(json.dumps({"error": "参数不正确：需要 symbol"}, ensure_ascii=False))
            sys.exit(1)
        
        symbol = sys.argv[2]
        get_stock_info(symbol)
        
    elif command == "get_index_list":
        get_index_list()
        
    else:
        print(json.dumps({"error": f"未知命令：{command}"}, ensure_ascii=False))
        sys.exit(1)

if __name__ == "__main__":
    main()