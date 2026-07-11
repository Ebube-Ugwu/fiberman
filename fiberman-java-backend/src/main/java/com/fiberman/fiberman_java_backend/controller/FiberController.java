package com.fiberman.fiberman_java_backend.controller;

import com.fiberman.fiberman_java_backend.dto.CreateInvoiceApiRequest;
import com.fiberman.fiberman_java_backend.dto.FiberCallResponse;
import com.fiberman.fiberman_java_backend.dto.GetChannelApiRequest;
import com.fiberman.fiberman_java_backend.dto.GetPaymentStatusApiRequest;
import com.fiberman.fiberman_java_backend.dto.InvoiceQrCodeRequest;
import com.fiberman.fiberman_java_backend.dto.SendPaymentApiRequest;
import com.fiberman.fiberman_java_backend.service.FiberGatewayService;
import com.fiberman.fiberman_java_backend.service.FiberHistoryService;
import jakarta.servlet.http.HttpSession;
import jakarta.validation.Valid;
import java.util.List;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/fiber")
public class FiberController {
    private final FiberGatewayService fiberGatewayService;
    private final FiberHistoryService fiberHistoryService;

    public FiberController(FiberGatewayService fiberGatewayService, FiberHistoryService fiberHistoryService) {
        this.fiberGatewayService = fiberGatewayService;
        this.fiberHistoryService = fiberHistoryService;
    }

    @GetMapping("/node-info")
    public FiberCallResponse nodeInfo(HttpSession session) {
        return fiberGatewayService.nodeInfo(session);
    }

    @GetMapping("/channels")
    public FiberCallResponse listChannels(HttpSession session) {
        return fiberGatewayService.listChannels(session);
    }

    @PostMapping("/channels/details")
    public FiberCallResponse getChannel(
            HttpSession session,
            @Valid @RequestBody GetChannelApiRequest request
    ) {
        return fiberGatewayService.getChannel(session, request);
    }

    @GetMapping("/peers")
    public FiberCallResponse listPeers(HttpSession session) {
        return fiberGatewayService.listPeers(session);
    }

    @PostMapping("/invoices")
    public FiberCallResponse createInvoice(
            HttpSession session,
            @Valid @RequestBody CreateInvoiceApiRequest request
    ) {
        return fiberGatewayService.createInvoice(session, request);
    }

    @PostMapping("/invoices/qr")
    public FiberCallResponse invoiceQrCode(@Valid @RequestBody InvoiceQrCodeRequest request) {
        return fiberGatewayService.invoiceQrCode(request.invoice(), request.size());
    }

    @PostMapping("/payments")
    public FiberCallResponse sendPayment(
            HttpSession session,
            @Valid @RequestBody SendPaymentApiRequest request
    ) {
        return fiberGatewayService.sendPayment(session, request);
    }

    @PostMapping("/payments/status")
    public FiberCallResponse paymentStatus(
            HttpSession session,
            @Valid @RequestBody GetPaymentStatusApiRequest request
    ) {
        return fiberGatewayService.paymentStatus(session, request);
    }

    @GetMapping("/history")
    public List<FiberCallResponse> history(HttpSession session) {
        return fiberHistoryService.getHistory(session);
    }

    @DeleteMapping("/history")
    public void clearHistory(HttpSession session) {
        fiberHistoryService.clear(session);
    }
}
