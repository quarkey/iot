import { Component, OnInit } from "@angular/core";
import { Dashboard } from "src/app/models/general";
import { GeneralService } from "src/app/services/general.service";

@Component({
  selector: "app-dashboard",
  templateUrl: "./dashboard.component.html",
  styleUrls: ["./dashboard.component.scss"],
})
export class DashboardComponent implements OnInit {
  loading = true;
  infoItems: Dashboard;
  constructor(private generalService: GeneralService) {}

  ngOnInit(): void {
    this.generalService.Dashboard().subscribe((res) => {
      if (res) {
        this.infoItems = res;
        this.loading = false;
      }
    });
  }
}
