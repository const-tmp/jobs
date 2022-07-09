FROM python:3.10

WORKDIR /bot

ENV PYTHONUNBUFFERED 1

COPY requirements.txt .
RUN pip install -r requirements.txt

COPY jbot jbot
COPY setup.py .
RUN pip install .

CMD jb