#!/usr/bin/env python3
"""
Quick test script to verify win rate calculation fix
"""

import requests
import json

def test_win_rate():
    print("🧪 Testing win rate calculation...")
    
    # Run a backtest
    request_data = {
        "index_id": "csi300",
        "strategy": {
            "type": "monthly_rotation",
            "parameters": {
                "buy_days_before_month_end": 1,
                "sell_days_after_month_start": 1
            }
        },
        "start_date": "2024-01-01",
        "end_date": "2024-03-31",
        "initial_cash": 1000000
    }
    
    try:
        response = requests.post("http://localhost:8080/api/v1/backtest", json=request_data, timeout=30)
        
        if response.status_code != 200:
            print(f"❌ HTTP Error: {response.status_code}")
            return
            
        result = response.json()
        
        if not result.get("success"):
            print(f"❌ API Error: {result.get('error')}")
            return
            
        data = result["data"]
        trades = data["trades"]
        metrics = data["performance_metrics"]
        
        print(f"✅ Backtest completed successfully!")
        print(f"Total trades: {len(trades)}")
        
        # Manual calculation
        buy_trades = [t for t in trades if t["action"] == "buy"]
        sell_trades = [t for t in trades if t["action"] == "sell"]
        
        print(f"Buy trades: {len(buy_trades)}")
        print(f"Sell trades: {len(sell_trades)}")
        
        # Calculate round trips manually
        winning_trades = 0
        total_round_trips = 0
        
        for i in range(min(len(buy_trades), len(sell_trades))):
            buy = buy_trades[i]
            sell = sell_trades[i]
            
            pnl = (sell["price"] - buy["price"]) * buy["quantity"] - buy["commission"] - sell["commission"]
            total_round_trips += 1
            
            if pnl > 0:
                winning_trades += 1
                
            print(f"Round trip {i+1}: Buy ${buy['price']:.2f} -> Sell ${sell['price']:.2f} = PnL ${pnl:.2f} ({'WIN' if pnl > 0 else 'LOSS'})")
        
        manual_win_rate = winning_trades / total_round_trips if total_round_trips > 0 else 0
        reported_win_rate = metrics["win_rate"]
        
        print(f"\n📊 WIN RATE ANALYSIS:")
        print(f"Manual calculation: {manual_win_rate:.2%} ({winning_trades}/{total_round_trips})")
        print(f"Reported win rate: {reported_win_rate:.2%}")
        print(f"Total trades (API): {metrics['total_trades']}")
        
        if abs(manual_win_rate - reported_win_rate) < 0.001:
            print("✅ Win rate calculation is CORRECT!")
        else:
            print("❌ Win rate calculation is INCORRECT!")
            
        print(f"\n📈 Other metrics:")
        print(f"Total return: {metrics['total_return']:.4%}")
        print(f"Sharpe ratio: {metrics['sharpe_ratio']:.3f}")
        
    except Exception as e:
        print(f"❌ Test failed: {e}")

if __name__ == "__main__":
    test_win_rate()