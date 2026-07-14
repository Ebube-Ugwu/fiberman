import { Injectable, computed, signal } from '@angular/core';
import { PlaygroundSettings, UpdatePlaygroundSettingsPayload } from './fiber-types';
import { FiberApiService } from './fiber-api.service';

const DEFAULT_SETTINGS: PlaygroundSettings = {
  nodeUrl: 'http://127.0.0.1:8227',
  authToken: '',
  timeoutSeconds: 30,
  defaultInvoiceCurrency: '',
  playgroundBaseUrl: 'http://localhost:9020'
};

@Injectable({ providedIn: 'root' })
export class PlaygroundSettingsService {
  private readonly settingsState = signal<PlaygroundSettings>(DEFAULT_SETTINGS);
  protected readonly loadingState = signal(false);
  protected readonly savingState = signal(false);
  protected readonly errorState = signal<string | null>(null);

  constructor(private readonly api: FiberApiService) {}

  readonly settings = computed(() => this.settingsState());
  readonly defaultInvoiceCurrency = computed(() => this.settingsState().defaultInvoiceCurrency.trim());
  readonly loading = computed(() => this.loadingState());
  readonly saving = computed(() => this.savingState());
  readonly error = computed(() => this.errorState());

  load(): void {
    if (this.loadingState()) {
      return;
    }

    this.loadingState.set(true);
    this.errorState.set(null);
    this.api.getSettings().subscribe({
      next: (settings) => {
        this.settingsState.set(settings);
        this.loadingState.set(false);
      },
      error: (error) => {
        this.errorState.set(error.error?.message ?? 'Unable to load settings.');
        this.loadingState.set(false);
      }
    });
  }

  save(payload: UpdatePlaygroundSettingsPayload, onSuccess?: () => void): void {
    this.savingState.set(true);
    this.errorState.set(null);
    this.api.updateSettings(payload).subscribe({
      next: (settings) => {
        this.settingsState.set(settings);
        this.savingState.set(false);
        onSuccess?.();
      },
      error: (error) => {
        this.errorState.set(error.error?.message ?? 'Unable to save settings.');
        this.savingState.set(false);
      }
    });
  }
}
