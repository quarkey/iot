import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatIconModule } from '@angular/material/icon';
import { MatButtonModule } from '@angular/material/button';
import { ToolbarComponent } from './navigation/toolbar/toolbar.component';
import { SidenavComponent } from './navigation/sidenav/sidenav.component';
import { MatListModule } from '@angular/material/list';
import { DashboardComponent } from './views/dashboard/dashboard.component';
import { DevicesComponent } from './views/devices/devices.component';
import { DatasetsComponent } from './views/datasets/datasets.component';
import { HttpClientModule } from '@angular/common/http';
import { DeviceComponent } from './views/devices/device/device.component';
import { DatasetComponent } from './views/datasets/dataset/dataset.component';
import { DetailsComponent } from './views/datasets/dataset/details/details.component';
import { MatTabsModule } from '@angular/material/tabs';
import { DataComponent } from './views/datasets/dataset/data/data.component';
import { NewDeviceDialogComponent } from './dialogs/new-device/new-device.component';
import { MatDialogModule } from '@angular/material/dialog';
import { MatCardModule } from '@angular/material/card';
import { MatInputModule } from '@angular/material/input';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';

@NgModule({
  declarations: [
    AppComponent,
    ToolbarComponent,
    SidenavComponent,
    DashboardComponent,
    DevicesComponent,
    DatasetsComponent,
    DeviceComponent,
    DatasetComponent,
    DetailsComponent,
    DataComponent,
    NewDeviceDialogComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    MatSidenavModule,
    MatToolbarModule,
    MatIconModule,
    MatButtonModule,
    MatListModule,
    HttpClientModule,
    MatTabsModule,
    MatButtonModule,
    MatDialogModule,
    MatCardModule,
    MatInputModule,
    ReactiveFormsModule,
    FormsModule,
  ],
  providers: [NewDeviceDialogComponent],
  bootstrap: [AppComponent],
})
export class AppModule {}
