#!/bin/bash
# Setup AKShare environment for Macro Strategy Platform

echo "Setting up AKShare environment..."

# Create virtual environment if it doesn't exist
if [ ! -d "venv" ]; then
    echo "Creating Python virtual environment..."
    python3 -m venv venv
fi

# Activate virtual environment and install dependencies
echo "Installing Python dependencies..."
source venv/bin/activate
pip install --upgrade pip
pip install akshare pandas

# Test AKShare installation
echo "Testing AKShare installation..."
python3 -c "import akshare as ak; print('AKShare version:', ak.__version__)"

# Test the AKShare client script
echo "Testing AKShare client script..."
python3 backend/scripts/akshare_client.py get_stock_zh_a_hist sh000300 20240101 20240105 | head -n 1

echo "AKShare setup completed successfully!"
echo "The backend is now configured to use real AKShare data for A-share indexes."
echo ""
echo "To start the application:"
echo "1. Backend: cd backend && go run cmd/main.go"
echo "2. Frontend: cd frontend && npm run dev"