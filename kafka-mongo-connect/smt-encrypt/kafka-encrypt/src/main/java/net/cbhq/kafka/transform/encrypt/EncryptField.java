package net.cbhq.kafka.transform.encrypt;

import net.cbhq.kafka.transform.cipher.Encryptor;
import org.apache.kafka.connect.connector.ConnectRecord;
import org.apache.kafka.connect.data.Schema;
import org.apache.kafka.connect.transforms.Transformation;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.util.*;

public class EncryptField<R extends ConnectRecord<R>> implements Transformation<R> {
    private static final Logger log = LoggerFactory.getLogger(EncryptField.class);
    private static final ObjectMapper mapper = new ObjectMapper();

    final private static String FULL_DOCUMENT = "fullDocument";
    final private static Encryptor ENCRYPTOR = new Encryptor("1234567890123456");

    private List<String> fields;

    @Override
    public R apply(R record) {
        try {
            return applyInternal(record);
        } catch (Exception e) {
            log.info("an exception occurred", e);
        }

        return record;
    }

    public R applyInternal(R record) {
        log.info("Applying Encrypt Field: {}", record);

        if (record.value() == null) {
            log.info("record value is null");
            return record;
        }

        // This is the entire change record from Mongo CDC
        Object value = record.value();

        if (!(value instanceof String)) {
            log.info("value is unknown type: {}", value.getClass().getName());
            return record;
        }

        // This means the CDC is encoded as string.
        // Ex. In the document below, the payload is a stringified version of JSON.
        // "{"schema":{"type":"string","optional":false},"payload":"{\"_id\": {\"_data\": \"82681DD279000000012B022C0100296E5A10045F7DD9443EF84E5FA49F2398B31E948046645F69640064681DD279A17CDFD18BA00AC00004\"}, \"operationType\": \"insert\", \"clusterTime\": {\"$timestamp\": {\"t\": 1746784889, \"i\": 1}}, \"fullDocument\": {\"_id\": \"681dd279a17cdfd18ba00ac0\", \"empty\": \"notempty\"}, \"ns\": {\"db\": \"testdb\", \"coll\": \"testcoll\"}, \"documentKey\": {\"_id\": \"681dd279a17cdfd18ba00ac0\"}}"}"
        log.info("value is instance of String");
        Map<String, Object> jsonObject;

        try {
            // This converts the payload into Map<String, Object>.
            jsonObject = mapper.readValue((String) value, Map.class);
        } catch (Exception e) {
            log.info("Failed to parse value as JSON: {}", value, e);
            return record;
        }

        if (!jsonObject.containsKey(FULL_DOCUMENT)) {
            log.info("the record doesn't have fullDocument.");
            return record;
        }

        Object fullDocumentObject = jsonObject.get(FULL_DOCUMENT);
        if (!(fullDocumentObject instanceof Map) || fullDocumentObject == null) {
            log.info("fullDocument is null or not map {}", fullDocumentObject);
            return record;
        }

        // This throws exception if type casting issue
        Map<String, Object> fullDocument = (Map<String, Object>) fullDocumentObject;

        // Each field is encrypted
        for (String field : fields) {
            Object fieldValue = getField(fullDocument, field);
            if (fieldValue == null) {
                log.info("field {} is null", field);
                continue;
            }

            // Base64 encoded cipher text
            String encryptedFieldValue;
            try {
                encryptedFieldValue = encrypt(fieldValue);
                putField(fullDocument, field, encryptedFieldValue);
            } catch (Exception e) {
                log.info("Failed to encrypt field {}", field, e);
                continue;
            }

            try {
                Object decrypted = decryptField(encryptedFieldValue);
                log.info("Encrypted Field: {}; Value {} ;;;;; Decrypted Value: {}", field, encryptedFieldValue, decrypted);
            } catch (Exception e) {
                log.info("error while decrypting field: {}", field, e);
            }
        }

        log.info("successfully applied Encrypt Field: {}", record);
        try {
            String updatedChange = mapper.writeValueAsString(jsonObject);
            log.info("string schema boys");
            return record.newRecord(
                    record.topic(), record.kafkaPartition(),
                    record.keySchema(), record.key(),
                    Schema.STRING_SCHEMA, updatedChange,
                    record.timestamp()
            );
        } catch (Exception e) {
            log.info("error while mapping json object", e);
        }

        return record.newRecord(
                record.topic(), record.kafkaPartition(),
                record.keySchema(), record.key(),
                record.valueSchema(), record,
                record.timestamp()
        );

    }

    // Get nested field by path like "a.b.c"
    public static Object getField(Map<String, Object> jsonValue, String path) {
        String[] parts = path.split("\\.");
        Object current = jsonValue;

        for (String part : parts) {
            if (!(current instanceof Map)) {
                return null; // Path is invalid
            }
            current = ((Map<String, Object>) current).get(part);
            if (current == null) {
                return null; // Missing intermediate field
            }
        }

        return current;
    }

    public static void putField(Map<String, Object> jsonValue, String path, Object value) {
        String[] parts = path.split("\\.");
        Map<String, Object> current = jsonValue;

        for (int i = 0; i < parts.length - 1; i++) {
            String part = parts[i];
            Object next = current.get(part);

            if (!(next instanceof Map)) {
                // Create intermediate map if missing
                next = new LinkedHashMap<String, Object>();
                current.put(part, next);
            }

            current = (Map<String, Object>) next;
        }

        // Set value at the final part
        current.put(parts[parts.length - 1], value);
    }


    // Returns Base64 encoded cipher text string
    private String encrypt(Object o) {
        try {
            String plainText;

            if (o instanceof String || o instanceof Number || o instanceof Boolean) {
                plainText = o.toString();
            } else {
                // Complex object (Map, List, etc.)
                plainText = mapper.writeValueAsString(o);
            }


            String cipherText = ENCRYPTOR.encrypt(plainText);

            return cipherText;
        } catch (Exception e) {
            log.info("Failed to encrypt value: {}", o);
            throw new RuntimeException("Failed to encrypt value", e);
        }
    }

    public static Object decryptField(String encryptedValue) throws Exception {
        try {
            String plainText = ENCRYPTOR.decrypt(encryptedValue);

            // Try to parse as JSON
            try {
                Object parsed = mapper.readValue(plainText, Object.class);
                return parsed;
            } catch (Exception e) {
                // Not JSON, treat as String
                return plainText;
            }
        } catch (Exception e) {
            log.info("Failed to decrypt value: {}", encryptedValue);
            return encryptedValue;
        }
    }

    @Override
    public void configure(Map<String, ?> configs) {
        log.info("Configuring AppendATransform");
        try {
            String v = (String) configs.get("fields");
            if (fields == null) fields = new ArrayList<>();

            fields = Arrays.asList(v.split(","));
            log.info("configured encryption fields: {}", fields);
        } catch (Exception e) {
            log.info("an exception occurred while configuring connector", e);
        }
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
