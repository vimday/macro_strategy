#!/bin/bash

# Setup script for AKShare Python environment

echo "🚀 Setting up AKShare environment..."

# Check if Python3 is available
if ! command -v python3 &> /dev/null; then
    echo "❌ Python3 not found. Please install Python3 first."
    exit 1
fi

echo "✅ Python3 found"

# Check if pip is available
if ! command -v pip3 &> /dev/null; then
    echo "❌ pip3 not found. Please install pip3 first."
    exit 1
fi

echo "✅ pip3 found"

# Create virtual environment
echo "🔧 Creating virtual environment..."
python3 -m venv akshare_env

# Activate virtual environment
echo "🔧 Activating virtual environment..."
source akshare_env/bin/activate

# Upgrade pip
echo "🔧 Upgrading pip..."
pip3 install --upgrade pip

# Install AKShare
echo "🔧 Installing AKShare..."
pip3 install akshare

# Test AKShare installation
echo "🧪 Testing AKShare installation..."
python3 -c "import akshare; print('✅ AKShare version:', akshare.__version__)"

# Test the akshare_client.py script
echo "🧪 Testing akshare_client.py script..."
cd backend/scripts
python3 akshare_client.py get_index_list > /dev/null 2>&1

if [ $? -eq 0 ]; then
    echo "✅ akshare_client.py works correctly"
else
    echo "❌ akshare_client.py has issues"
fi

cd ../..

echo "🎉 AKShare setup complete!"
echo ""
echo "To use AKShare in the future, activate the virtual environment with:"
echo "source akshare_env/bin/activate"