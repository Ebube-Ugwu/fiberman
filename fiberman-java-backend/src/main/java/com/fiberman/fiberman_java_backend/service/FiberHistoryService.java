package com.fiberman.fiberman_java_backend.service;

import com.fiberman.fiberman_java_backend.dto.FiberCallResponse;
import jakarta.servlet.http.HttpSession;
import java.util.ArrayList;
import java.util.List;
import org.springframework.stereotype.Service;

@Service
public class FiberHistoryService {
    private static final String SESSION_KEY = "fiberCallHistory";
    private static final int MAX_HISTORY_SIZE = 20;

    public FiberCallResponse append(HttpSession session, FiberCallResponse entry) {
        List<FiberCallResponse> history = getHistoryInternal(session);
        history.add(0, entry);
        if (history.size() > MAX_HISTORY_SIZE) {
            history.remove(history.size() - 1);
        }
        session.setAttribute(SESSION_KEY, history);
        return entry;
    }

    public List<FiberCallResponse> getHistory(HttpSession session) {
        return List.copyOf(getHistoryInternal(session));
    }

    public void clear(HttpSession session) {
        session.removeAttribute(SESSION_KEY);
    }

    @SuppressWarnings("unchecked")
    private List<FiberCallResponse> getHistoryInternal(HttpSession session) {
        Object attribute = session.getAttribute(SESSION_KEY);
        if (attribute instanceof List<?> history) {
            return (List<FiberCallResponse>) history;
        }
        List<FiberCallResponse> history = new ArrayList<>();
        session.setAttribute(SESSION_KEY, history);
        return history;
    }
}
