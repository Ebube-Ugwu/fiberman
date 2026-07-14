import { Routes } from '@angular/router';
import { AppShellComponent } from './layout/app-shell.component';
import { DashboardPageComponent } from './pages/dashboard-page.component';
import { InvoiceBuilderPageComponent } from './pages/invoice-builder-page.component';
import { LogsPageComponent } from './pages/logs-page.component';
import { NetworkTopologyPageComponent } from './pages/network-topology-page.component';
import { PaymentsHistoryPageComponent } from './pages/payments-history-page.component';
import { RpcExplorerPageComponent } from './pages/rpc-explorer-page.component';
import { SettingsPageComponent } from './pages/settings-page.component';

export const routes: Routes = [
  {
    path: '',
    component: AppShellComponent,
    children: [
      { path: '', pathMatch: 'full', redirectTo: 'dashboard' },
      { path: 'dashboard', component: DashboardPageComponent },
      { path: 'explorer', component: RpcExplorerPageComponent },
      { path: 'invoice', component: InvoiceBuilderPageComponent },
      { path: 'topology', component: NetworkTopologyPageComponent },
      { path: 'payments', component: PaymentsHistoryPageComponent },
      { path: 'logs', component: LogsPageComponent },
      { path: 'settings', component: SettingsPageComponent }
    ]
  },
  { path: '**', redirectTo: 'dashboard' }
];
