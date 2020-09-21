FROM python:3-slim
MAINTAINER Michalski Luc <michalski.luc@gmail.com>

WORKDIR /opt/service

COPY requirements.txt .

RUN apt-get update && \
    apt-get install -y --no-install-recommends git ca-certificates gcc libc6-dev && \
    python3 -m pip install -U pip && \
    pip3 install -r requirements.txt

COPY . .

CMD ["python3", "train.py"]
