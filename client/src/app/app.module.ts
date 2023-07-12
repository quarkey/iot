import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { MatSidenavModule } from '@angular/material/sidenav';
import { MatToolbarModule } from '@angular/material/toolbar';
import { MatIconModule } from '@angular/material/icon';
import { ToolbarComponent } from './navigation/toolbar/toolbar.component';
import { SidenavComponent } from './navigation/sidenav/sidenav.component';
import { MatListModule } from '@angular/material/list';
import { DashboardComponent } from './views/dashboard/dashboard.component';
import { DevicesListComponent } from './views/devices-list/devices-list.component';
import { DatasetsListComponent } from './views/dataset-list/dataset-list.component';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';
import { DeviceComponent } from './views/device/device.component';
import { DatasetComponent } from './views/dataset/dataset.component';
import { MatTabsModule } from '@angular/material/tabs';
import { DatasetAreaChartNgxComponent } from './views/dataset/dataset-area-chart-ngx/dataset-area-chart-ngx.component';
import { NewDeviceDialogComponent } from './components/dialogs/new-device/new-device.component';
// import { MatLegacyDialogModule as MatDialogModule } from '@angular/material/legacy-dialog';
import { MatDialogModule } from '@angular/material/dialog';

// import { MatLegacyInputModule as MatInputModule } from '@angular/material/legacy-input';
import { MatInputModule } from '@angular/material/input';

import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { NewDatasetDialogComponent } from './components/dialogs/new-dataset/new-dataset.component';
// import { MatLegacySelectModule as MatSelectModule } from '@angular/material/legacy-select';
import { MatSelectModule } from '@angular/material/select';

import { NgxChartsModule } from '@swimlane/ngx-charts';
// import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatProgressSpinnerModule } from '@angular/material/progress-spinner';

import { DatasetLineChartComponent } from './views/dataset/dataset-line-chart/dataset-line-chart.component';
import { NgChartsModule } from 'ng2-charts';
import { EventsComponent } from './views/events/events.component';

// import { MatLegacyTableModule as MatTableModule } from '@angular/material/legacy-table';
import { MatTableModule } from '@angular/material/table';

import { MatSortModule } from '@angular/material/sort';
import { HttpErrorInterceptor } from './services/httperrorinterceptor.service';
import { ErrorHandlingDialogComponent } from './components/dialogs/error-handling-dialog/error-handling-dialog.component';
import { DatasetDetailsComponent } from './views/dataset/dataset-details/dataset-details.component';
import { DatasetCardComponent } from './components/dataset-card/dataset-card.component';
// import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';

import { DeviceCardComponent } from './components/device-card/device-card.component';
import { DashComponent } from './views/dash/dash.component';
import { MatGridListModule } from '@angular/material/grid-list';
import { MatCardModule } from '@angular/material/card';
import { MatMenuModule } from '@angular/material/menu';
import { LayoutModule } from '@angular/cdk/layout';
import { SystemComponent } from './views/system/system.component';
import { ControllersListComponent } from './views/controllers-list/controllers-list.component';
import { ControllerCardComponent } from './components/controller-card/controller-card.component';
import { ControllerDetailsComponent } from './components/controller-details/controller-details.component';
import { ControllerComponent } from './views/controller/controller.component';
import { NewControllerComponent } from './components/dialogs/new-controller/new-controller.component';
import { TwoColumnDivComponent } from './shared/auxiliary/components/two-column-div/two-column-div.component';
import { StatusDotComponent } from './shared/auxiliary/components/status-dot/status-dot.component';
// import { MatLegacySlideToggleModule as MatSlideToggleModule } from '@angular/material/legacy-slide-toggle';

import { MatTimepickerModule } from 'mat-timepicker';
import { MatDatepickerModule } from '@angular/material/datepicker';
import { ConfirmationDialogComponent } from './components/dialogs/confirmation-dialog/confirmation-dialog.component';
import { SelectDeviceComponent } from './shared/auxiliary/forms/select-device/select-device.component';
import { TipTextComponent } from './shared/auxiliary/components/tip-text/tip-text.component';
import { ControllerTableComponent } from './components/controller-table/controller-table.component';
import { SensorIconComponent } from './shared/auxiliary/components/sensor-icon/sensor-icon.component';
import { DatasetTableComponent } from './components/dataset-table/dataset-table.component';

import { DatasetOverviewComponent } from './views/dataset/dataset-overview/dataset-overview.component';
import { MinMaxAverageComponent } from './shared/auxiliary/components/min-max-average/min-max-average.component';
import { DisplayValueBoxComponent } from './shared/auxiliary/components/display-value-box/display-value-box.component';
import { MatButtonModule } from '@angular/material/button';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatTooltipModule } from '@angular/material/tooltip';

@NgModule({
  declarations: [
    AppComponent,
    ToolbarComponent,
    SidenavComponent,
    DashboardComponent,
    DevicesListComponent,
    DatasetsListComponent,
    DeviceComponent,
    DatasetComponent,
    DatasetDetailsComponent,
    DatasetAreaChartNgxComponent,
    NewDeviceDialogComponent,
    NewDatasetDialogComponent,
    DatasetLineChartComponent,
    EventsComponent,
    ErrorHandlingDialogComponent,
    DatasetCardComponent,
    DeviceCardComponent,
    DashComponent,
    SystemComponent,
    ControllersListComponent,
    ControllerCardComponent,
    ControllerDetailsComponent,
    ControllerComponent,
    NewControllerComponent,
    TwoColumnDivComponent,
    StatusDotComponent,
    ConfirmationDialogComponent,
    SelectDeviceComponent,
    TipTextComponent,
    ControllerTableComponent,
    SensorIconComponent,
    DatasetTableComponent,
    DatasetOverviewComponent,
    MinMaxAverageComponent,
    DisplayValueBoxComponent,
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    // MatTimepickerModule,
    MatSidenavModule,
    MatDatepickerModule,
    MatToolbarModule,
    MatIconModule,
    MatButtonModule,
    // MatGridListModule,
    MatListModule,
    HttpClientModule,
    MatTabsModule,
    MatButtonModule,
    MatDialogModule,
    MatCardModule,
    MatInputModule,
    ReactiveFormsModule,
    FormsModule,
    MatFormFieldModule,
    MatSelectModule,
    NgxChartsModule,
    MatProgressSpinnerModule,
    NgChartsModule,
    MatTableModule,
    MatSortModule,
    // FlexLayoutModule,
    MatTooltipModule,
    MatGridListModule,
    MatMenuModule,
    LayoutModule,
    // MatSlideToggleModule,
  ],
  providers: [
    NewDeviceDialogComponent,
    NewDatasetDialogComponent,
    ErrorHandlingDialogComponent,
    ConfirmationDialogComponent,
    {
      provide: HTTP_INTERCEPTORS,
      useClass: HttpErrorInterceptor,
      multi: true,
    },
  ],
  bootstrap: [AppComponent],
})
export class AppModule {}
