package com.fiberman.fiberman_java_backend.service;

import com.fiberman.fiberman_java_backend.dto.CodeArtifacts;
import com.fiberman.fiberman_java_backend.dto.FiberCallResponse;
import java.time.Instant;
import java.util.List;
import java.util.Map;
import org.junit.jupiter.api.Test;
import org.springframework.mock.web.MockHttpSession;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertTrue;

class FiberHistoryServiceTest {
    private final FiberHistoryService historyService = new FiberHistoryService();

    @Test
    void appendsEntriesMostRecentFirst() {
        MockHttpSession session = new MockHttpSession();
        FiberCallResponse first = response("1", "node_info");
        FiberCallResponse second = response("2", "list_channels");

        historyService.append(session, first);
        historyService.append(session, second);

        List<FiberCallResponse> history = historyService.getHistory(session);
        assertEquals(2, history.size());
        assertEquals("2", history.get(0).historyId());
        assertEquals("1", history.get(1).historyId());
    }

    @Test
    void keepsHistoryBounded() {
        MockHttpSession session = new MockHttpSession();

        for (int index = 1; index <= 25; index++) {
            historyService.append(session, response(String.valueOf(index), "method_" + index));
        }

        List<FiberCallResponse> history = historyService.getHistory(session);
        assertEquals(20, history.size());
        assertEquals("25", history.get(0).historyId());
        assertEquals("6", history.get(19).historyId());
    }

    @Test
    void clearsHistory() {
        MockHttpSession session = new MockHttpSession();
        historyService.append(session, response("1", "node_info"));

        historyService.clear(session);

        assertTrue(historyService.getHistory(session).isEmpty());
    }

    private FiberCallResponse response(String historyId, String method) {
        return new FiberCallResponse(
                historyId,
                method,
                "/api/fiber/" + method,
                Instant.parse("2026-07-11T00:00:00Z"),
                true,
                Map.of(),
                Map.of("ok", true),
                null,
                new CodeArtifacts("curl", "java"),
                null);
    }
}
