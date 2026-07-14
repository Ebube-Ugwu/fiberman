import { CommonModule, DatePipe } from '@angular/common';
import { Component, OnInit, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { FiberApiService } from '../core/fiber-api.service';
import { FiberCallResponse } from '../core/fiber-types';

@Component({
  selector: 'app-payments-history-page',
  imports: [CommonModule, FormsModule, DatePipe],
  templateUrl: './payments-history-page.component.html',
  styleUrl: './payments-history-page.component.css'
})
export class PaymentsHistoryPageComponent implements OnInit {
  private readonly api = inject(FiberApiService);

  protected readonly history = signal<FiberCallResponse[]>([]);
  protected readonly statusFilter = signal('all');

  ngOnInit(): void {
    this.api.getHistory().subscribe({
      next: (history) => this.history.set(history),
      error: () => this.history.set([])
    });
  }

  protected amount(entry: FiberCallResponse): number {
    const value = entry.params?.['amount'];
    return typeof value === 'number' ? value : Number(value ?? 0);
  }

  protected rows(): FiberCallResponse[] {
    return this.history()
      .filter((entry) => ['send_payment', 'get_payment', 'new_invoice'].includes(entry.method))
      .filter(
        (entry) =>
          this.statusFilter() === 'all' || (entry.success ? 'settled' : 'failed') === this.statusFilter()
      );
  }
}
