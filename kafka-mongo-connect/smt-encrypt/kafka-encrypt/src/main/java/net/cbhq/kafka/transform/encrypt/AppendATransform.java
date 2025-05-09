package net.cbhq.kafka.transform.encrypt;

import org.apache.kafka.connect.connector.ConnectRecord;
import org.apache.kafka.connect.data.Field;
import org.apache.kafka.connect.data.Schema;
import org.apache.kafka.connect.data.Struct;
import org.apache.kafka.connect.transforms.Transformation;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;

public class AppendATransform<R extends ConnectRecord<R>> implements Transformation<R> {
    private static final Logger log = LoggerFactory.getLogger(AppendATransform.class);

    public static void main(String[] args) {
        log.info("hello");
        System.out.println("From sout");
    }

    @Override
    public R apply(R record) {
        log.info("Applying append transformeeee: {}", record);

        if (record.value() == null) {
            log.info("record value is null");
            return record;
        }

        Object value = record.value();

        if (value instanceof Struct) {
            log.info("value is instance of Struct");
            Struct struct = (Struct) value;
            Schema schema = struct.schema();

            Struct updatedStruct = new Struct(schema);
            for (Field field : schema.fields()) {
                Object fieldValue = struct.get(field);
                if (fieldValue instanceof String) {
                    updatedStruct.put(field.name(), ((String) fieldValue) + "a");
                } else {
                    updatedStruct.put(field.name(), fieldValue);
                }
            }

            return record.newRecord(
                    record.topic(), record.kafkaPartition(),
                    record.keySchema(), record.key(),
                    schema, updatedStruct,
                    record.timestamp()
            );
        } else if (value instanceof Map) {
            log.info("value is instance of Map");
            @SuppressWarnings("unchecked")
            Map<String, Object> map = (Map<String, Object>) value;

            Map<String, Object> updatedMap = new java.util.HashMap<>();
            for (Map.Entry<String, Object> entry : map.entrySet()) {
                Object fieldValue = entry.getValue();
                if (fieldValue instanceof String) {
                    updatedMap.put(entry.getKey(), ((String) fieldValue) + "a");
                } else {
                    updatedMap.put(entry.getKey(), fieldValue);
                }
            }

            return record.newRecord(
                    record.topic(), record.kafkaPartition(),
                    record.keySchema(), record.key(),
                    null, updatedMap,
                    record.timestamp()
            );
        } else if (value instanceof String) {
            log.info("value is instance of Stringee");
            String updatedValue = ((String) value) + "a";
            return record.newRecord(
                    record.topic(), record.kafkaPartition(),
                    record.keySchema(), record.key(),
                    record.valueSchema(), updatedValue,
                    record.timestamp()
            );
        }

        log.info("value is unknown type: {}", value.getClass().getName());
        return record;
    }

    @Override
    public void configure(Map<String, ?> configs) {
        log.info("Configuring AppendATransform");
    }

    @Override
    public void close() {
        log.info("Closing AppendATransform");
    }

    @Override
    public org.apache.kafka.common.config.ConfigDef config() {
        return new org.apache.kafka.common.config.ConfigDef();
    }
}
