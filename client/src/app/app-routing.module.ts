import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { DashboardComponent } from './views/dashboard/dashboard.component';
import { DatasetComponent } from './views/datasets/dataset/dataset.component';
import { DatasetsComponent } from './views/datasets/datasets.component';
import { DeviceComponent } from './views/devices/device/device.component';
import { DevicesComponent } from './views/devices/devices.component';

const routes: Routes = [
  { path: 'dashboard', component: DashboardComponent },
  { path: 'devices', component: DevicesComponent },
  { path: 'devices/:arduino_key', component: DeviceComponent },
  { path: 'datasets', component: DatasetsComponent },
  { path: 'datasets/:reference', component: DatasetComponent },

  {
    path: '**',
    redirectTo: '/dashboard',
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
