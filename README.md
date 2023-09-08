# Cowrie Honeypot Kafka Ingestor (Cowrie-Honeypot-Ingest)


Cowrie-Honeypot-Ingest is an OpenSearch ingestor for Cowrie Honeypots. It allows you to read sessions through logs transfered by hpfeeds, send them to Kafka, and then consume and ingest them into OpenSearch after enrichment.

## Data Pipeline

The Cowrie-Honeypot-Ingest program follows the following data pipeline:

1. **Loading Sessions:** The program loads sessions from the specified log file.
2. **Sending to Kafka:** Loaded sessions are sent to a Kafka broker using the specified Kafka brokers and topic.
3. **Consuming in Logstash:** Logstash consumes the sessions from Kafka, performs any necessary enrichment or transformations, and prepares the data for indexing in OpenSearch.
4. **Indexing in OpenSearch:** The enriched sessions are indexed in OpenSearch for further analysis and processing. The index name and other relevant configurations can be customized through the environment variables in the `config.env` file.

Please note that the data pipeline involves the additional step of consuming data in Logstash for enrichment and preparation before indexing in OpenSearch.

## Usage

1. Ensure that `sub.go` is running seperately and piping output to desired path.
2. Ensure that the `config.env` file is correctly set with the required parameters for your setup.

3. Build and run the application using Docker Compose:

   ```shell
   docker-compose up --build
   ```

4. The Cowrie-Honeypot-Ingest program will start loading sessions, sending them to Kafka, consuming them in Logstash, and finally indexing them in OpenSearch.



## References

For more information about OpenSearch, Logstash, Kafka, and other related technologies, please refer to the following resources:

- [OpenSearch Documentation](https://opensearch.org/docs/1.3/)
- [Logstash Documentation](https://opensearch.org/docs/latest/tools/logstash/index/)
- [Kafka Documentation](https://github.com/Shopify/sarama)

## User Manual - Customizing Configuration Variables

The following guide provides an overview of the configuration variables available in the `config.env` file(rename configBackup.env to config.env). These variables allow you to customize the behavior of the application according to your specific requirements.

**Note:** Please make sure to follow the instructions carefully and set the appropriate values for each variable.


1. `ENABLE_POLLING`:
   - Description: This variable determines whether polling for new sessions should be enabled or disabled.
   - Accepted values: `true` or `false`
   - Default value: `false`
   - Instructions: Set this variable to `true` if you want to enable polling for new sessions. Otherwise, set it to `false`.

2. `POLLING_INTERVAL`:
   - Description: This variable sets the frequency at which polling for new sessions is performed.
   - Default value: `5m` (1 minute)
   - Instructions: Specify the desired interval for polling in minutes. For example, if you want to poll every 5 minutes, set the value to `5m`.

3. `KAFKA_BROKERS`, `INDEX`, `OPENSEARCH_URL`, `KAFKA_BOOTSTRAP_SERVERS`, `CODEC_CHARSET`, `KAFKA_TOPIC`, `GROUP_ID`, `AUTO_OFFSET_RESET`, `OPENSEARCH_HOSTS`, `OPENSEARCH_USER`, `OPENSEARCH_PASSWORD`, `OPENSEARCH_SSL`, `OPENSEARCH_INDEX`, `OPENSEARCH_MANAGE_TEMPLATE`, `KAFKA_LISTENERS`, `KAFKA_ADVERTISED_LISTENERS`, `KAFKA_ZOOKEEPER_CONNECT`, `KAFKA_CREATE_TOPICS`:
   - Description: These variables control various configurations related to Kafka, Opensearch, and other dependencies.
   - Instructions: If you need to configure any of these variables, provide the appropriate values based on your Kafka and Opensearch setups. Refer to the specific documentation or consult with your system administrator for the correct values.

4. Save the changes to the `config.env` file once you have customized the variables according to your needs.

By following these instructions and setting the variables in the `config.env` file, you can customize the application to fit your specific environment and requirements. Make sure to restart the application after making any changes to the configuration file for the modifications to take effect.

