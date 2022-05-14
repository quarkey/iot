import { NgModule } from "@angular/core";
import { RouterModule, Routes } from "@angular/router";
import { DashboardComponent } from "./views/dashboard/dashboard.component";
import { DatasetComponent } from "./views/dataset/dataset.component";
import { DatasetsListComponent } from "./views/dataset-list/dataset-list.component";
import { DeviceComponent } from "./views/device/device.component";
import { DevicesListComponent } from "./views/devices-list/devices-list.component";
import { EventsComponent } from "./views/events/events.component";
import { DashComponent } from "./views/dash/dash.component";

const routes: Routes = [
  { path: "dashboard", component: DashboardComponent },
  { path: "dash", component: DashComponent },
  { path: "devices", component: DevicesListComponent },
  { path: "devices/:arduino_key", component: DeviceComponent },
  { path: "datasets", component: DatasetsListComponent },
  { path: "datasets/:reference", component: DatasetComponent },
  { path: "events", component: EventsComponent },

  {
    path: "**",
    redirectTo: "/dashboard",
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
