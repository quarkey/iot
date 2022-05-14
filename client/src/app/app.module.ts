import { NgModule } from "@angular/core";
import { BrowserModule } from "@angular/platform-browser";

import { AppRoutingModule } from "./app-routing.module";
import { AppComponent } from "./app.component";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { MatSidenavModule } from "@angular/material/sidenav";
import { MatToolbarModule } from "@angular/material/toolbar";
import { MatIconModule } from "@angular/material/icon";
import { MatButtonModule } from "@angular/material/button";
import { ToolbarComponent } from "./navigation/toolbar/toolbar.component";
import { SidenavComponent } from "./navigation/sidenav/sidenav.component";
import { MatListModule } from "@angular/material/list";
import { DashboardComponent } from "./views/dashboard/dashboard.component";
import { DevicesListComponent } from "./views/devices-list/devices-list.component";
import { DatasetsListComponent } from "./views/dataset-list/dataset-list.component";
import { HttpClientModule, HTTP_INTERCEPTORS } from "@angular/common/http";
import { DeviceComponent } from "./views/device/device.component";
import { DatasetComponent } from "./views/dataset/dataset.component";
import { MatTabsModule } from "@angular/material/tabs";
import { LineChartNgxComponent } from "./components/line-chart-ngx/line-chart-ngx.component";
import { NewDeviceDialogComponent } from "./components/dialogs/new-device/new-device.component";
import { MatDialogModule } from "@angular/material/dialog";
import { MatCardModule } from "@angular/material/card";
import { MatInputModule } from "@angular/material/input";
import { FormsModule, ReactiveFormsModule } from "@angular/forms";
import { NewDatasetDialogComponent } from "./components/dialogs/new-dataset/new-dataset.component";
import { MatSelectModule } from "@angular/material/select";
import { NgxChartsModule } from "@swimlane/ngx-charts";
import { MatProgressSpinnerModule } from "@angular/material/progress-spinner";
import { LineChartComponent } from "./components/line-chart/line-chart.component";
import { NgChartsModule } from "ng2-charts";
import { EventsComponent } from "./views/events/events.component";
import { MatTableModule } from "@angular/material/table";
import { MatSortModule } from "@angular/material/sort";
import { HttpErrorInterceptor } from "./services/httperrorinterceptor.service";
import { ErrorHandlingDialogComponent } from "./components/dialogs/error-handling-dialog/error-handling-dialog.component";
import { FlexLayoutModule } from "@angular/flex-layout";
import { DatasetDetailsComponent } from "./components/dataset-details/dataset-details.component";
import { DatasetCardComponent } from "./components/dataset-card/dataset-card.component";
import { MatTooltipModule } from "@angular/material/tooltip";
import { DeviceCardComponent } from "./components/device-card/device-card.component";
import { DashComponent } from './views/dash/dash.component';
import { MatGridListModule } from '@angular/material/grid-list';
import { MatMenuModule } from '@angular/material/menu';
import { LayoutModule } from '@angular/cdk/layout';
import { SystemComponent } from './views/system/system.component';

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
    LineChartNgxComponent,
    NewDeviceDialogComponent,
    NewDatasetDialogComponent,
    LineChartComponent,
    EventsComponent,
    ErrorHandlingDialogComponent,
    DatasetCardComponent,
    DeviceCardComponent,
    DashComponent,
    SystemComponent,
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
    MatSelectModule,
    NgxChartsModule,
    MatProgressSpinnerModule,
    NgChartsModule,
    MatTableModule,
    MatSortModule,
    FlexLayoutModule,
    MatTooltipModule,
    MatGridListModule,
    MatMenuModule,
    LayoutModule,
  ],
  providers: [
    NewDeviceDialogComponent,
    NewDatasetDialogComponent,
    ErrorHandlingDialogComponent,
    {
      provide: HTTP_INTERCEPTORS,
      useClass: HttpErrorInterceptor,
      multi: true,
    },
  ],
  bootstrap: [AppComponent],
})
export class AppModule {}
