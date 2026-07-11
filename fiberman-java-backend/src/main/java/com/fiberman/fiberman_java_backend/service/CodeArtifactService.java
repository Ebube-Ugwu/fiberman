package com.fiberman.fiberman_java_backend.service;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fiberman.fiberman_java_backend.config.FiberPlaygroundProperties;
import com.fiberman.fiberman_java_backend.config.FiberNodeProperties;
import com.fiberman.fiberman_java_backend.dto.CodeArtifacts;
import java.util.Map;
import org.springframework.stereotype.Service;

@Service
public class CodeArtifactService {
    private final ObjectMapper objectMapper;
    private final FiberPlaygroundProperties playgroundProperties;
    private final FiberNodeProperties nodeProperties;

    public CodeArtifactService(
            ObjectMapper objectMapper,
            FiberPlaygroundProperties playgroundProperties,
            FiberNodeProperties nodeProperties
    ) {
        this.objectMapper = objectMapper;
        this.playgroundProperties = playgroundProperties;
        this.nodeProperties = nodeProperties;
    }

    public CodeArtifacts generate(String backendPath, String sdkMethodCall, Object requestBody) {
        return new CodeArtifacts(
                buildCurl(backendPath, requestBody),
                buildJavaSnippet(sdkMethodCall));
    }

    private String buildCurl(String backendPath, Object requestBody) {
        StringBuilder builder = new StringBuilder()
                .append("curl -X ")
                .append(requestBody == null ? "GET" : "POST")
                .append(" \"")
                .append(playgroundProperties.baseUrl())
                .append(backendPath)
                .append("\"");

        if (requestBody != null) {
            builder.append(" \\\n  -H \"Content-Type: application/json\"")
                    .append(" \\\n  -d '")
                    .append(escapeSingleQuotes(toJson(requestBody)))
                    .append("'");
        }

        return builder.toString();
    }

    private String buildJavaSnippet(String sdkMethodCall) {
        String authLine = nodeProperties.authToken() == null || nodeProperties.authToken().isBlank()
                ? ""
                : "    .authToken(\"" + escapeJava(nodeProperties.authToken()) + "\")\n";

        return """
                FiberClient client = FiberClient.builder()
                    .baseUrl("%s")
                %s    .build();

                var response = %s;
                """.formatted(
                escapeJava(nodeProperties.url()),
                authLine,
                sdkMethodCall);
    }

    private String toJson(Object requestBody) {
        try {
            return objectMapper.writeValueAsString(requestBody);
        } catch (JsonProcessingException exception) {
            throw new IllegalStateException("Failed to serialize request body for code generation", exception);
        }
    }

    private String escapeSingleQuotes(String value) {
        return value.replace("'", "'\"'\"'");
    }

    private String escapeJava(String value) {
        return value.replace("\\", "\\\\").replace("\"", "\\\"");
    }
}
