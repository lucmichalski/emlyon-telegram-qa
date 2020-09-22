FROM python:3-slim
MAINTAINER Michalski Luc <michalski.luc@gmail.com>

RUN apt-get update && \
    apt-get install -y --no-install-recommends gcc ca-certificates git libc6-dev && \
    python3 -m pip install -U pip && \
    pip3 install git+https://github.com/deepset-ai/haystack

COPY . .

EXPOSE 8000

CMD ["gunicorn", "application:app",  "-b", "0.0.0.0:8000", "-k", "uvicorn.workers.UvicornWorker", "--workers", "2", "--timeout", "180"]
