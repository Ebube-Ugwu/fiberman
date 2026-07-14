import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable } from 'rxjs';
import {
  CreateInvoicePayload,
  FiberCallResponse,
  GetChannelPayload,
  GetPaymentPayload,
  InvoiceQrPayload,
  PlaygroundSettings,
  SendPaymentPayload,
  UpdatePlaygroundSettingsPayload
} from './fiber-types';

@Injectable({ providedIn: 'root' })
export class FiberApiService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = '/api/fiber';
  private readonly settingsUrl = '/api/settings';

  getNodeInfo(): Observable<FiberCallResponse> {
    return this.http.get<FiberCallResponse>(`${this.baseUrl}/node-info`);
  }

  listChannels(): Observable<FiberCallResponse> {
    return this.http.get<FiberCallResponse>(`${this.baseUrl}/channels`);
  }

  listPeers(): Observable<FiberCallResponse> {
    return this.http.get<FiberCallResponse>(`${this.baseUrl}/peers`);
  }

  createInvoice(payload: CreateInvoicePayload): Observable<FiberCallResponse> {
    return this.http.post<FiberCallResponse>(`${this.baseUrl}/invoices`, payload);
  }

  generateInvoiceQr(payload: InvoiceQrPayload): Observable<FiberCallResponse> {
    return this.http.post<FiberCallResponse>(`${this.baseUrl}/invoices/qr`, payload);
  }

  sendPayment(payload: SendPaymentPayload): Observable<FiberCallResponse> {
    return this.http.post<FiberCallResponse>(`${this.baseUrl}/payments`, payload);
  }

  getChannel(payload: GetChannelPayload): Observable<FiberCallResponse> {
    return this.http.post<FiberCallResponse>(`${this.baseUrl}/channels/details`, payload);
  }

  getPaymentStatus(payload: GetPaymentPayload): Observable<FiberCallResponse> {
    return this.http.post<FiberCallResponse>(`${this.baseUrl}/payments/status`, payload);
  }

  getHistory(): Observable<FiberCallResponse[]> {
    return this.http.get<FiberCallResponse[]>(`${this.baseUrl}/history`);
  }

  clearHistory(): Observable<void> {
    return this.http.delete<void>(`${this.baseUrl}/history`);
  }

  getSettings(): Observable<PlaygroundSettings> {
    return this.http.get<PlaygroundSettings>(this.settingsUrl);
  }

  updateSettings(payload: UpdatePlaygroundSettingsPayload): Observable<PlaygroundSettings> {
    return this.http.put<PlaygroundSettings>(this.settingsUrl, payload);
  }
}
