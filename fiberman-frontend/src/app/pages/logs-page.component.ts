import { CommonModule, DatePipe } from '@angular/common';
import { Component, OnInit, computed, inject, signal } from '@angular/core';
import { FiberApiService } from '../core/fiber-api.service';
import { FiberCallResponse } from '../core/fiber-types';

@Component({
  selector: 'app-logs-page',
  imports: [CommonModule, DatePipe],
  templateUrl: './logs-page.component.html',
  styleUrl: './logs-page.component.css'
})
export class LogsPageComponent implements OnInit {
  private readonly api = inject(FiberApiService);

  protected readonly history = signal<FiberCallResponse[]>([]);
  protected readonly logLines = computed(() =>
    this.history().map((entry) => ({
      timestamp: entry.timestamp,
      level: entry.success ? 'INFO' : 'ERROR',
      message: `${entry.method} ${entry.success ? 'completed' : 'failed'}`
    }))
  );

  ngOnInit(): void {
    this.api.getHistory().subscribe({
      next: (history) => this.history.set(history),
      error: () => this.history.set([])
    });
  }
}
