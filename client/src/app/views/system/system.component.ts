import { Component, OnInit } from "@angular/core";
import { GeneralService } from "src/app/services/general.service";

@Component({
  selector: "app-system",
  templateUrl: "./system.component.html",
  styleUrls: ["./system.component.scss"],
})
export class SystemComponent implements OnInit {
  stats: any;
  constructor(private generalService: GeneralService) {}

  ngOnInit(): void {
    this.generalService.ServerHealth().subscribe((res) => {
      if (res) {
        this.stats = res;
      }
    });
  }
}
