import { CommonModule } from '@angular/common';
import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { FiberApiService } from '../core/fiber-api.service';
import { FiberCallResponse } from '../core/fiber-types';
import { PlaygroundSettingsService } from '../core/playground-settings.service';

type MethodKey =
  | 'node_info'
  | 'list_channels'
  | 'list_peers'
  | 'new_invoice'
  | 'send_payment';

interface MethodOption {
  key: MethodKey;
  label: string;
}

type SnippetLanguage = 'java' | 'golang' | 'typescript' | 'rust';

@Component({
  selector: 'app-rpc-explorer-page',
  imports: [CommonModule, FormsModule],
  templateUrl: './rpc-explorer-page.component.html',
  styleUrl: './rpc-explorer-page.component.css'
})
export class RpcExplorerPageComponent implements OnInit {
  private readonly api = inject(FiberApiService);
  private readonly settingsService = inject(PlaygroundSettingsService);

  protected readonly methods: MethodOption[] = [
    { key: 'new_invoice', label: 'create_invoice' },
    { key: 'node_info', label: 'node_info' },
    { key: 'list_channels', label: 'list_channels' },
    { key: 'list_peers', label: 'list_peers' },
    { key: 'send_payment', label: 'send_payment' }
  ];

  protected readonly selectedMethod = signal<MethodKey>('new_invoice');
  protected readonly currentResponse = signal<FiberCallResponse | null>(null);
  protected readonly history = signal<FiberCallResponse[]>([]);
  protected readonly loading = signal(false);
  protected readonly activePane = signal<'json' | 'snippets' | 'history'>('json');
  protected readonly defaultInvoiceCurrency = computed(() => this.settingsService.defaultInvoiceCurrency());
  protected readonly formError = signal<string | null>(null);
  protected readonly selectedSnippetLanguage = signal<SnippetLanguage>('java');
  protected readonly snippetPopupMessage = signal<string | null>(null);
  protected readonly payloadCopied = signal(false);

  protected invoiceAmount = 25000;
  protected currency = '';
  protected description = 'Project Milestone Payment 1';
  protected expirySeconds = 3600;
  protected invoice = '';
  protected paymentAmount: number | null = null;
  protected timeoutSeconds = 30;
  protected channelId = '';
  protected paymentId = '';

  ngOnInit(): void {
    this.settingsService.load();
    this.loadHistory();
  }

  protected execute(): void {
    this.formError.set(null);
    this.loading.set(true);

    const request = this.selectedMethod();
    if (request === 'new_invoice' && !this.resolvedCurrency()) {
      this.loading.set(false);
      this.formError.set('Set a currency explicitly or configure a default invoice currency in Settings.');
      return;
    }

    const task =
      request === 'node_info'
        ? this.api.getNodeInfo()
        : request === 'list_channels'
          ? this.api.listChannels()
        : request === 'list_peers'
          ? this.api.listPeers()
          : request === 'new_invoice'
            ? this.api.createInvoice({
                amount: this.invoiceAmount,
                currency: this.resolvedCurrency(),
                description: this.description,
                expirySeconds: this.expirySeconds
              })
            : this.api.sendPayment({
                invoice: this.invoice,
                amount: this.paymentAmount,
                timeoutSeconds: this.timeoutSeconds
              });

    task.subscribe({
      next: (response) => {
        this.currentResponse.set(response);
        this.activePane.set('json');
        this.selectedSnippetLanguage.set('java');
        this.loading.set(false);
        this.loadHistory();
      },
      error: (error) => {
        const fallback: FiberCallResponse = {
          historyId: crypto.randomUUID(),
          method: this.selectedMethod(),
          backendPath: '',
          timestamp: new Date().toISOString(),
          success: false,
          params: {},
          result: null,
          error: error.error ?? { message: 'Request failed' },
          codeArtifacts: null,
          invoiceQrCode: null
        };
        this.currentResponse.set(fallback);
        this.activePane.set('json');
        this.loading.set(false);
        this.loadHistory();
      }
    });
  }

  protected selectHistory(entry: FiberCallResponse): void {
    this.currentResponse.set(entry);
    this.activePane.set('json');
  }

  protected selectSnippetLanguage(language: SnippetLanguage): void {
    if (language === 'typescript' || language === 'rust') {
      this.snippetPopupMessage.set(`${this.languageLabel(language)} SDK snippets are planned but not implemented yet.`);
      return;
    }
    this.selectedSnippetLanguage.set(language);
  }

  protected snippetLanguageIs(language: SnippetLanguage): boolean {
    return this.selectedSnippetLanguage() === language;
  }

  protected formattedSnippet(): string {
    const artifacts = this.currentResponse()?.codeArtifacts;
    if (!artifacts) {
      return 'Execute a method to generate runnable SDK code.';
    }

    return this.selectedSnippetLanguage() === 'golang'
      ? artifacts.goSnippet || 'The backend did not return a Go snippet for this call yet.'
      : artifacts.javaSnippet || 'The backend did not return a Java snippet for this call yet.';
  }

  protected copySnippet(language: SnippetLanguage): void {
    if (language === 'typescript' || language === 'rust') {
      this.selectSnippetLanguage(language);
      return;
    }

    const artifacts = this.currentResponse()?.codeArtifacts;
    const value = language === 'golang' ? artifacts?.goSnippet : artifacts?.javaSnippet;
    this.copy(value ?? '');
    this.selectedSnippetLanguage.set(language);
  }

  protected copyCurrentSnippet(): void {
    this.copySnippet(this.selectedSnippetLanguage());
  }

  protected currentSnippetCopyLabel(): string {
    return this.selectedSnippetLanguage() === 'golang' ? 'Copy as Golang' : 'Copy as Java';
  }

  protected closeSnippetPopup(): void {
    this.snippetPopupMessage.set(null);
  }

  protected copy(value: string | null | undefined): void {
    if (!value) {
      return;
    }
    navigator.clipboard.writeText(value).catch(() => undefined);
  }

  protected formattedCurrentPayload(): string {
    return JSON.stringify(this.currentResponse()?.result ?? this.currentResponse()?.error ?? {}, null, 2);
  }

  protected copyFormattedPayload(): void {
    navigator.clipboard.writeText(this.formattedCurrentPayload()).then(() => {
      this.payloadCopied.set(true);
      window.setTimeout(() => this.payloadCopied.set(false), 1200);
    }).catch(() => undefined);
  }

  protected payloadCopyLabel(): string {
    return this.payloadCopied() ? 'Copied' : 'Copy';
  }

  protected supportsField(field: 'invoice' | 'amount' | 'description' | 'expiry'): boolean {
    const method = this.selectedMethod();
    return (
      (field === 'amount' && (method === 'new_invoice' || method === 'send_payment')) ||
      (field === 'description' && method === 'new_invoice') ||
      (field === 'expiry' && method === 'new_invoice') ||
      (field === 'invoice' && method === 'send_payment')
    );
  }

  private loadHistory(): void {
    this.api.getHistory().subscribe({
      next: (history) => this.history.set(history),
      error: () => this.history.set([])
    });
  }

  private resolvedCurrency(): string {
    return this.currency.trim() || this.defaultInvoiceCurrency();
  }

  private languageLabel(language: SnippetLanguage): string {
    return language === 'golang' ? 'Go' : language.charAt(0).toUpperCase() + language.slice(1);
  }
}
