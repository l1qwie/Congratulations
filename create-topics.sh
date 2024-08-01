#!/bin/bash

BROKER="congratulations-kafka:9092"

kafka-topics.sh --create --bootstrap-server $BROKER --replication-factor 1 --partitions 1 --topic employee-redis
kafka-topics.sh --create --bootstrap-server $BROKER --replication-factor 1 --partitions 1 --topic employee-sub
kafka-topics.sh --create --bootstrap-server $BROKER --replication-factor 1 --partitions 1 --topic employee-other
