package net.cbhq.kafka.transform.encrypt;

import net.cbhq.kafka.transform.cipher.Encryptor;
import org.apache.kafka.connect.connector.ConnectRecord;
import org.apache.kafka.connect.data.Field;
import org.apache.kafka.connect.data.Schema;
import org.apache.kafka.connect.data.Struct;
import org.apache.kafka.connect.transforms.Transformation;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.Map;

public class EncryptField<R extends ConnectRecord<R>> implements Transformation<R> {
    private static final Logger log = LoggerFactory.getLogger(EncryptField.class);
    private static final ObjectMapper mapper = new ObjectMapper();

    @Override
    public R apply(R record) {
        log.info("Applying Encrypt Field: {}", record);

        if (record.value() == null) {
            log.info("record value is null");
            return record;
        }

        // This is the entire change record from Mongo CDC
        Object value = record.value();

        if (value instanceof String) {
            // This means the CDC is encoded as string.
            // Ex. In the document below, the payload is a stringified version of JSON.
            // "{"schema":{"type":"string","optional":false},"payload":"{\"_id\": {\"_data\": \"82681DD279000000012B022C0100296E5A10045F7DD9443EF84E5FA49F2398B31E948046645F69640064681DD279A17CDFD18BA00AC00004\"}, \"operationType\": \"insert\", \"clusterTime\": {\"$timestamp\": {\"t\": 1746784889, \"i\": 1}}, \"fullDocument\": {\"_id\": \"681dd279a17cdfd18ba00ac0\", \"empty\": \"notempty\"}, \"ns\": {\"db\": \"testdb\", \"coll\": \"testcoll\"}, \"documentKey\": {\"_id\": \"681dd279a17cdfd18ba00ac0\"}}"}"
            log.info("value is instance of String");
            Map<String, Object> jsonValue;

            try {
                // This converts the payload into Map<String, Object>.
                jsonValue = mapper.readValue((String) value, Map.class);
            } catch (Exception e) {
                log.error("Failed to parse value as JSON: {}", value, e);
                return record;
            }

            // You can add additional fields to the CDC
            jsonValue.put("someField", "newValue");

            // You can encrypt the values here
            Map<String, Object> fullDocument = (Map<String, Object>) jsonValue.get("fullDocument");
            String field = "field1";
            Object fieldValue = fullDocument.get(field);

            // Encrypt fieldValue
            fieldValue = encrypt(fieldValue);
            fullDocument.put(field, fieldValue);

            jsonValue.put("fullDocument", fullDocument);

            String updatedValue = ((String) value);
            return record.newRecord(
                    record.topic(), record.kafkaPartition(),
                    record.keySchema(), record.key(),
                    null, jsonValue,
                    record.timestamp()
            );
        }

        log.info("value is unknown type: {}", value.getClass().getName());
        return record;
    }

    private Object encrypt(Object o) {
        if (o instanceof String) {
            String s = (String) o;
            // Secret must be 16 characters long
            Encryptor encryptor = new Encryptor("1234567890123456");
            String cipherText = encryptor.encrypt(s);
            return cipherText;
        }
        return o;
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
