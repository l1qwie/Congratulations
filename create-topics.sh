#!/bin/bash

kafka-topics.sh --create --topic employee --partitions 2 --replication-factor 1 --if-not-exists --bootstrap-server localhost:9092
