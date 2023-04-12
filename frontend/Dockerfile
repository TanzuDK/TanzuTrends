# Base image
FROM python:3.8-slim-buster

# Install system dependencies
RUN apt-get update && \
    apt-get -y install gcc && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install Python dependencies
COPY requirements.txt requirements.txt
RUN pip install --upgrade pip && \
    pip install -r requirements.txt

# Copy app files
WORKDIR /app
COPY main.py /app/main.py

# Expose port
EXPOSE 8501

# Start Streamlit
CMD streamlit run main.py --server.port $PORT
