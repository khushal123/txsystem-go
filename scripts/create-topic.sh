#!/bin/bash
set -e

# Default values for topics, partitions, replication
DEFAULT_TOPICS=${DEFAULT_TOPICS:-"transactions,logs"}
DEFAULT_PARTITIONS=${DEFAULT_PARTITIONS:-1}
DEFAULT_REPLICATION=${DEFAULT_REPLICATION:-1}

# Use env vars if set, else use defaults
TOPICS=${KAFKA_TOPIC_TRANSCATIONSS:-$DEFAULT_TOPICS}
PARTITIONS=${KAFKA_PARTITIONS:-$DEFAULT_PARTITIONS}
REPLICATION=${KAFKA_REPLICATION:-$DEFAULT_REPLICATION}
BOOTSTRAP_SERVER=${BOOTSTRAP_SERVER:-kafka:9092}

echo "Waiting for Kafka at $BOOTSTRAP_SERVER to be ready..."

MAX_RETRIES=30
RETRY_INTERVAL=2
COUNT=0

until /opt/kafka/bin/kafka-topics.sh --bootstrap-server $BOOTSTRAP_SERVER --list > /dev/null 2>&1; do
  ((COUNT++))
  if [ $COUNT -ge $MAX_RETRIES ]; then
    echo "Error: Kafka not ready after $MAX_RETRIES attempts."
    exit 1
  fi
  echo "Kafka not ready yet. Retry $COUNT/$MAX_RETRIES, waiting $RETRY_INTERVAL seconds..."
  sleep $RETRY_INTERVAL
done

echo "Kafka is ready! Creating topics..."

# Split comma-separated topics into array
IFS=',' read -ra TOPIC_ARRAY <<< "$TOPICS"

for TOPIC in "${TOPIC_ARRAY[@]}"; do
  echo "Creating topic: $TOPIC"
  /opt/kafka/bin/kafka-topics.sh --bootstrap-server $BOOTSTRAP_SERVER --create --if-not-exists \
    --topic "$TOPIC" --partitions "$PARTITIONS" --replication-factor "$REPLICATION"

  # Verify topic creation
  if /opt/kafka/bin/kafka-topics.sh --bootstrap-server $BOOTSTRAP_SERVER --list | grep -q "^$TOPIC$"; then
    echo "✅ Topic '$TOPIC' created successfully"
  else
    echo "❌ Failed to create topic '$TOPIC'"
  fi
done

echo -e "\nAll Kafka topics:"
/opt/kafka/bin/kafka-topics.sh --bootstrap-server $BOOTSTRAP_SERVER --list

echo -e "\n✅ Topic setup complete!"
